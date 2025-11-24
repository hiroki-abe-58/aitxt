package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	proofreadProvider string
	proofreadStream   bool
	proofreadLang     string
	proofreadStyle    string
)

var proofreadCmd = &cobra.Command{
	Use:   "proofread [file or text]",
	Short: "Proofread and correct text using AI",
	Long: `Proofread text from a file or direct input using AI.
Corrects grammar, spelling, punctuation, and improves clarity.

Styles:
  standard  - General proofreading (default)
  formal    - Formal/business writing
  casual    - Casual/friendly tone
  academic  - Academic writing style
  technical - Technical documentation

Examples:
  aitxt proofread document.txt
  aitxt proofread "I has a apple" 
  aitxt proofread essay.txt --style academic
  aitxt proofread report.txt --style formal --lang ja
  cat document.txt | aitxt proofread -`,
	Args: cobra.MaximumNArgs(1),
	RunE: runProofread,
}

func init() {
	rootCmd.AddCommand(proofreadCmd)

	proofreadCmd.Flags().StringVarP(&proofreadProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	proofreadCmd.Flags().BoolVarP(&proofreadStream, "stream", "s", false, "Enable streaming output")
	proofreadCmd.Flags().StringVarP(&proofreadLang, "lang", "l", "", "Text language (auto-detect if not specified)")
	proofreadCmd.Flags().StringVarP(&proofreadStyle, "style", "t", "standard", "Writing style (standard, formal, casual, academic, technical)")
}

func runProofread(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if proofreadProvider != "" {
		provider = llm.Provider(proofreadProvider)
	}

	// Get input text
	var inputText string
	if len(args) == 0 || args[0] == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		inputText = string(data)
	} else {
		if _, err := os.Stat(args[0]); err == nil {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			inputText = string(data)
		} else {
			inputText = args[0]
		}
	}

	inputText = strings.TrimSpace(inputText)
	if inputText == "" {
		return fmt.Errorf("no input text provided")
	}

	// Create LLM client
	llmConfig, err := cfg.ToLLMConfig(provider)
	if err != nil {
		return err
	}

	factory := llm.NewFactory()
	if err := factory.RegisterConfig(llmConfig); err != nil {
		return fmt.Errorf("failed to register config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	client, err := factory.CreateClientWithContext(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Build prompt
	stylePrompt := getStylePrompt(proofreadStyle)
	langPrompt := ""
	if proofreadLang != "" {
		langPrompt = fmt.Sprintf("The text is in %s. ", getLanguageName(proofreadLang))
	}

	systemMsg := fmt.Sprintf(`You are a professional proofreader and editor. %s%s
Your task is to:
1. Fix grammar, spelling, and punctuation errors
2. Improve clarity and readability
3. Maintain the original meaning and tone

Output format:
- First, provide the corrected text
- Then, list the changes made (if any) in a brief summary`, langPrompt, stylePrompt)

	prompt := fmt.Sprintf("Please proofread and correct the following text:\n\n%s", inputText)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	// Generate proofread text
	fmt.Printf("Proofreading with %s (%s style)...\n\n", provider, proofreadStyle)

	if proofreadStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to proofread: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to proofread: %w", err)
	}

	return nil
}

func getStylePrompt(style string) string {
	switch style {
	case "formal":
		return "Use formal, professional language suitable for business communication."
	case "casual":
		return "Use casual, friendly language while maintaining correctness."
	case "academic":
		return "Use academic writing conventions with precise, scholarly language."
	case "technical":
		return "Use clear, precise technical writing style suitable for documentation."
	default:
		return "Use clear, standard language."
	}
}
