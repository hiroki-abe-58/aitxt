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
	translateProvider string
	translateStream   bool
	translateTo       string
	translateFrom     string
)

var translateCmd = &cobra.Command{
	Use:   "translate [file or text]",
	Short: "Translate text using AI",
	Long: `Translate text from a file or direct input using AI.

Supported languages: en (English), ja (Japanese), zh (Chinese), 
                     ko (Korean), th (Thai), es (Spanish), 
                     fr (French), de (German)

Examples:
  aitxt translate "Hello, world!" --to ja
  aitxt translate document.txt --to zh
  aitxt translate japanese.txt --from ja --to en
  cat document.txt | aitxt translate - --to ko`,
	Args: cobra.MaximumNArgs(1),
	RunE: runTranslate,
}

func init() {
	rootCmd.AddCommand(translateCmd)

	translateCmd.Flags().StringVarP(&translateProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	translateCmd.Flags().BoolVarP(&translateStream, "stream", "s", false, "Enable streaming output")
	translateCmd.Flags().StringVarP(&translateTo, "to", "t", "en", "Target language")
	translateCmd.Flags().StringVarP(&translateFrom, "from", "f", "", "Source language (auto-detect if not specified)")
}

func runTranslate(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if translateProvider != "" {
		provider = llm.Provider(translateProvider)
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
	targetLang := getLanguageName(translateTo)
	var prompt string
	var systemMsg string

	if translateFrom != "" {
		sourceLang := getLanguageName(translateFrom)
		systemMsg = fmt.Sprintf("You are a professional translator. Translate the text from %s to %s. Output only the translation, no explanations.", sourceLang, targetLang)
		prompt = inputText
	} else {
		systemMsg = fmt.Sprintf("You are a professional translator. Detect the source language and translate the text to %s. Output only the translation, no explanations.", targetLang)
		prompt = inputText
	}

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	// Generate translation
	fmt.Printf("Translating to %s with %s...\n\n", targetLang, provider)

	if translateStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to translate: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to translate: %w", err)
	}

	return nil
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en": "English",
		"ja": "Japanese",
		"zh": "Chinese",
		"ko": "Korean",
		"th": "Thai",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"pt": "Portuguese",
		"ru": "Russian",
		"ar": "Arabic",
		"hi": "Hindi",
		"vi": "Vietnamese",
	}

	if name, ok := languages[code]; ok {
		return name
	}
	return code
}
