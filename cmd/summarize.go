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
	summaryProvider string
	summaryStream   bool
	summaryLang     string
)

var summarizeCmd = &cobra.Command{
	Use:   "summarize [file or text]",
	Short: "Summarize text using AI",
	Long: `Summarize text from a file or direct input using AI.

Examples:
  aitxt summarize document.txt
  aitxt summarize "Long text to summarize..."
  aitxt summarize document.txt --provider claude
  cat document.txt | aitxt summarize -`,
	Args: cobra.MaximumNArgs(1),
	RunE: runSummarize,
}

func init() {
	rootCmd.AddCommand(summarizeCmd)

	summarizeCmd.Flags().StringVarP(&summaryProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	summarizeCmd.Flags().BoolVarP(&summaryStream, "stream", "s", false, "Enable streaming output")
	summarizeCmd.Flags().StringVarP(&summaryLang, "lang", "l", "en", "Output language (en, ja, zh, ko, th)")
}

func runSummarize(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if summaryProvider != "" {
		provider = llm.Provider(summaryProvider)
	}

	// Get input text
	var inputText string
	if len(args) == 0 || args[0] == "-" {
		// Read from stdin
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		inputText = string(data)
	} else {
		// Check if it's a file or direct text
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
	langPrompt := getLanguagePrompt(summaryLang)
	systemMsg := fmt.Sprintf("You are a helpful assistant that summarizes text concisely. %s", langPrompt)
	prompt := fmt.Sprintf("Please summarize the following text:\n\n%s", inputText)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 1000,
	}

	// Generate summary
	fmt.Printf("Summarizing with %s...\n\n", provider)

	if summaryStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to generate summary: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to generate summary: %w", err)
	}

	return nil
}

func getLanguagePrompt(lang string) string {
	switch lang {
	case "ja":
		return "Please respond in Japanese."
	case "zh":
		return "Please respond in Chinese."
	case "ko":
		return "Please respond in Korean."
	case "th":
		return "Please respond in Thai."
	default:
		return "Please respond in English."
	}
}
