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
	"github.com/hiroki-abe-58/aitxt/pkg/output"
	"github.com/hiroki-abe-58/aitxt/pkg/progress"
	"github.com/spf13/cobra"
)

var (
	askProvider    string
	askStream      bool
	askLang        string
	askSystem      string
	askTemperature float64
	askNoProgress  bool
	askOutput      string
)

var askCmd = &cobra.Command{
	Use:   "ask [question]",
	Short: "Ask AI any question",
	Long: `Ask AI any question directly from the command line.

A simple and versatile command for general AI interactions.
Supports context from files or stdin.

Output formats: text (default), json, yaml

Examples:
  aitxt ask "What is the capital of France?"
  aitxt ask "Explain quantum computing" --lang ja
  aitxt ask "Summarize this" < document.txt
  cat error.log | aitxt ask "What's wrong here?"
  aitxt ask "Write a haiku about coding" --stream
  aitxt ask "Be creative" --temperature 1.5
  aitxt ask "Hello" --output json`,
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
	askCmd.Flags().BoolVar(&askNoProgress, "no-progress", false, "Disable progress spinner")
	askCmd.Flags().StringVarP(&askOutput, "output", "o", "text", "Output format (text, json, yaml)")
}

func runAsk(cmd *cobra.Command, args []string) error {
	formatter := output.NewFormatter(askOutput)

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		if formatter.IsStructured() {
			return formatter.Print(output.ErrorResponse("", fmt.Errorf("failed to load config: %w", err)))
		}
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
			if formatter.IsStructured() {
				return formatter.Print(output.ErrorResponse(string(provider), err))
			}
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
		err := fmt.Errorf("no question provided. Usage: aitxt ask \"your question\"")
		if formatter.IsStructured() {
			return formatter.Print(output.ErrorResponse(string(provider), err))
		}
		return err
	}

	// Create LLM client
	llmConfig, err := cfg.ToLLMConfig(provider)
	if err != nil {
		if formatter.IsStructured() {
			return formatter.Print(output.ErrorResponse(string(provider), err))
		}
		return err
	}

	factory := llm.NewFactory()
	if err := factory.RegisterConfig(llmConfig); err != nil {
		if formatter.IsStructured() {
			return formatter.Print(output.ErrorResponse(string(provider), err))
		}
		return fmt.Errorf("failed to register config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	var client llm.Client
	showSpinner := !askNoProgress && !askStream && !formatter.IsStructured()

	if showSpinner {
		spinner := progress.NewSpinner(fmt.Sprintf("Connecting to %s...", provider))
		spinner.Start()
		client, err = factory.CreateClientWithContext(ctx, provider)
		if err != nil {
			spinner.Error("Failed to create client")
			return fmt.Errorf("failed to create client: %w", err)
		}
		spinner.Stop()
	} else {
		client, err = factory.CreateClientWithContext(ctx, provider)
		if err != nil {
			if formatter.IsStructured() {
				return formatter.Print(output.ErrorResponse(string(provider), err))
			}
			return fmt.Errorf("failed to create client: %w", err)
		}
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
		if formatter.IsStructured() {
			return fmt.Errorf("streaming not supported with JSON/YAML output")
		}
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
		return err
	}

	var resp *llm.Response
	if showSpinner {
		spinner := progress.NewSpinner("Thinking...")
		spinner.Start()
		resp, err = client.Generate(ctx, req)
		if err != nil {
			spinner.Error("Failed to generate response")
			return fmt.Errorf("failed to get response: %w", err)
		}
		spinner.Stop()
	} else {
		resp, err = client.Generate(ctx, req)
		if err != nil {
			if formatter.IsStructured() {
				return formatter.Print(output.ErrorResponse(string(provider), err))
			}
			return fmt.Errorf("failed to get response: %w", err)
		}
	}

	// Output response
	return formatter.PrintResponse(output.SuccessResponse(
		string(provider),
		resp.Model,
		resp.Text,
		resp.TokensUsed,
	))
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
