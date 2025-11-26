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
	askProvider    string
	askStream      bool
	askLang        string
	askSystem      string
	askTemperature float64
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask AI any question",
	Long: `Ask AI any question directly from the command line.

A simple and versatile command for general AI interactions.
Supports context from files or stdin.

Examples:
  aitxt ask "What is the capital of France?"
  aitxt ask "Explain quantum computing" --lang ja
  aitxt ask "Summarize this" < document.txt
  cat error.log | aitxt ask "What's wrong here?"
  aitxt ask "Write a haiku about coding" --stream
  aitxt ask "Be creative" --temperature 1.5`,
	Args: cobra.MaximumNArgs(1),
	RunE: runAsk,
}

func init() {
	rootCmd.AddCommand(askCmd)

	askCmd.Flags().StringVarP(&askProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	askCmd.Flags().BoolVarP(&askStream, "stream", "s", false, "Enable streaming output")
	askCmd.Flags().StringVarP(&askLang, "lang", "l", "", "Response language (en, ja, zh, ko, th)")
	askCmd.Flags().StringVarP(&askSystem, "system", "S", "", "Custom system prompt")
	askCmd.Flags().Float64VarP(&askTemperature, "temperature", "t", 0.7, "Temperature (0.0-2.0, higher = more creative)")
}

func runAsk(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if askProvider != "" {
		provider = llm.Provider(askProvider)
	}

	// Get question/prompt
	var question string
	var stdinContent string

	// Check if there's stdin input
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		stdinContent = strings.TrimSpace(string(data))
	}

	// Get question from args or stdin
	if len(args) > 0 {
		question = args[0]
	} else if stdinContent != "" {
		question = stdinContent
		stdinContent = ""
	} else {
		return fmt.Errorf("no question provided. Usage: aitxt ask \"your question\"")
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

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	client, err := factory.CreateClientWithContext(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Build prompt
	systemMsg := askSystem
	if systemMsg == "" {
		systemMsg = "You are a helpful AI assistant. Provide clear, accurate, and concise answers."
	}

	if askLang != "" {
		langPrompt := getAskLangPrompt(askLang)
		systemMsg += " " + langPrompt
	}

	prompt := question
	if stdinContent != "" {
		prompt = fmt.Sprintf("%s\n\nContext:\n```\n%s\n```", question, stdinContent)
	}

	req := &llm.Request{
		Prompt:      prompt,
		SystemMsg:   systemMsg,
		MaxTokens:   2000,
		Temperature: askTemperature,
	}

	// Generate response
	if askStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to get response: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[%s | Tokens: %d]\n", provider, resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to get response: %w", err)
	}

	return nil
}

func getAskLangPrompt(lang string) string {
	switch lang {
	case "ja":
		return "Respond in Japanese."
	case "zh":
		return "Respond in Chinese."
	case "ko":
		return "Respond in Korean."
	case "th":
		return "Respond in Thai."
	case "en":
		return "Respond in English."
	default:
		return ""
	}
}
