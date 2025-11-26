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
	styleProvider string
	styleStream   bool
	styleTo       string
	styleLang     string
)

var styleCmd = &cobra.Command{
	Use:   "style [file or text]",
	Short: "Transform text style using AI",
	Long: `Transform the writing style of text using AI.

Available styles:
  formal      - Professional, business-appropriate language
  casual      - Friendly, conversational tone
  academic    - Scholarly, research-paper style
  technical   - Clear, precise technical documentation
  creative    - Expressive, engaging prose
  simple      - Easy-to-understand, plain language
  persuasive  - Compelling, convincing tone
  humorous    - Light-hearted, witty style
  poetic      - Artistic, lyrical expression
  journalistic - News article style

Examples:
  aitxt style "This is my text" --to formal
  aitxt style document.txt --to casual
  aitxt style email.txt --to professional --lang ja
  cat essay.txt | aitxt style - --to academic`,
	Args: cobra.MaximumNArgs(1),
	RunE: runStyle,
}

func init() {
	rootCmd.AddCommand(styleCmd)

	styleCmd.Flags().StringVarP(&styleProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	styleCmd.Flags().BoolVarP(&styleStream, "stream", "s", false, "Enable streaming output")
	styleCmd.Flags().StringVarP(&styleTo, "to", "t", "formal", "Target style")
	styleCmd.Flags().StringVarP(&styleLang, "lang", "l", "", "Output language (en, ja, zh, ko, th)")
}

func runStyle(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if styleProvider != "" {
		provider = llm.Provider(styleProvider)
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
	styleDesc := getStyleDescription(styleTo)
	langPrompt := ""
	if styleLang != "" {
		langPrompt = getStyleLangPrompt(styleLang)
	}

	systemMsg := fmt.Sprintf(`You are an expert writer who can transform text into different writing styles.
Transform the given text to: %s

Guidelines:
- Maintain the original meaning and key information
- Adjust vocabulary, sentence structure, and tone appropriately
- Keep the same language unless specified otherwise
- Output only the transformed text, no explanations
%s`, styleDesc, langPrompt)

	prompt := fmt.Sprintf("Transform the following text:\n\n%s", inputText)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	// Generate styled text
	fmt.Printf("Transforming to %s style with %s...\n\n", styleTo, provider)

	if styleStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to transform style: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to transform style: %w", err)
	}

	return nil
}

func getStyleDescription(style string) string {
	styles := map[string]string{
		"formal":       "Formal style - Professional, business-appropriate language with proper grammar and sophisticated vocabulary",
		"casual":       "Casual style - Friendly, conversational tone as if talking to a friend",
		"academic":     "Academic style - Scholarly language suitable for research papers and academic publications",
		"technical":    "Technical style - Clear, precise documentation style with accurate terminology",
		"creative":     "Creative style - Expressive, engaging prose with vivid descriptions and varied sentence structures",
		"simple":       "Simple style - Plain, easy-to-understand language accessible to all readers",
		"persuasive":   "Persuasive style - Compelling, convincing tone designed to influence the reader",
		"humorous":     "Humorous style - Light-hearted, witty writing with appropriate humor",
		"poetic":       "Poetic style - Artistic, lyrical expression with rhythm and imagery",
		"journalistic": "Journalistic style - Objective, news-article format with clear facts and quotes",
	}

	if desc, ok := styles[style]; ok {
		return desc
	}
	return fmt.Sprintf("%s style", style)
}

func getStyleLangPrompt(lang string) string {
	switch lang {
	case "ja":
		return "Output the result in Japanese."
	case "zh":
		return "Output the result in Chinese."
	case "ko":
		return "Output the result in Korean."
	case "th":
		return "Output the result in Thai."
	case "en":
		return "Output the result in English."
	default:
		return ""
	}
}
