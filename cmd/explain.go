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
	explainProvider string
	explainStream   bool
	explainLang     string
	explainContext  string
)

var explainCmd = &cobra.Command{
	Use:   "explain [error message or file]",
	Short: "Explain error messages using AI",
	Long: `Explain error messages and provide solutions using AI.

Supports various types of errors:
  - Compiler errors (Go, Python, JavaScript, etc.)
  - Runtime errors and stack traces
  - Build/deployment errors
  - System errors

Examples:
  aitxt explain "undefined: fmt.Prinln"
  aitxt explain error.log
  aitxt explain "ECONNREFUSED 127.0.0.1:5432" --context "PostgreSQL connection"
  cat error.log | aitxt explain -
  aitxt explain "panic: runtime error: index out of range" --lang ja`,
	Args: cobra.MaximumNArgs(1),
	RunE: runExplain,
}

func init() {
	rootCmd.AddCommand(explainCmd)

	explainCmd.Flags().StringVarP(&explainProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	explainCmd.Flags().BoolVarP(&explainStream, "stream", "s", false, "Enable streaming output")
	explainCmd.Flags().StringVarP(&explainLang, "lang", "l", "en", "Output language (en, ja, zh, ko, th)")
	explainCmd.Flags().StringVarP(&explainContext, "context", "c", "", "Additional context about the error")
}

func runExplain(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if explainProvider != "" {
		provider = llm.Provider(explainProvider)
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
		return fmt.Errorf("no error message provided")
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
	langPrompt := getExplainLangPrompt(explainLang)
	contextPrompt := ""
	if explainContext != "" {
		contextPrompt = fmt.Sprintf("\nAdditional context: %s", explainContext)
	}

	systemMsg := fmt.Sprintf(`You are an expert developer who helps explain error messages.
%s

When explaining an error, provide:
1. **What the error means** - A clear explanation of the error
2. **Why it happened** - Common causes of this error
3. **How to fix it** - Step-by-step solutions
4. **Example** - Code example if applicable

Be concise but thorough. Focus on practical solutions.`, langPrompt)

	prompt := fmt.Sprintf("Please explain this error and how to fix it:%s\n\n```\n%s\n```", contextPrompt, inputText)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 1500,
	}

	// Generate explanation
	fmt.Printf("Analyzing error with %s...\n\n", provider)

	if explainStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to explain error: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to explain error: %w", err)
	}

	return nil
}

func getExplainLangPrompt(lang string) string {
	switch lang {
	case "ja":
		return "Respond in Japanese (日本語で回答してください)."
	case "zh":
		return "Respond in Chinese (请用中文回答)."
	case "ko":
		return "Respond in Korean (한국어로 답변해주세요)."
	case "th":
		return "Respond in Thai (กรุณาตอบเป็นภาษาไทย)."
	default:
		return "Respond in English."
	}
}
