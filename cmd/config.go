package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display configuration and API status",
	Long: `Display current configuration and check API key status.

Shows:
  - Configured API keys (masked)
  - Default provider
  - Available providers
  - Environment variables

Examples:
  aitxt config
  aitxt config --check`,
	RunE: runConfig,
}

var configCheck bool

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().BoolVarP(&configCheck, "check", "c", false, "Check API connectivity")
}

func runConfig(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("🤖 aitxt Configuration")
	fmt.Println("======================")
	fmt.Println()

	// Provider Status
	fmt.Println("📦 Providers:")
	fmt.Println()

	providers := []struct {
		name    string
		envKey  string
		apiKey  string
		provider llm.Provider
	}{
		{"OpenAI", "OPENAI_API_KEY", cfg.OpenAIKey, llm.ProviderOpenAI},
		{"Claude", "ANTHROPIC_API_KEY", cfg.ClaudeKey, llm.ProviderClaude},
		{"Gemini", "GOOGLE_API_KEY", cfg.GeminiKey, llm.ProviderGemini},
	}

	for _, p := range providers {
		status := "❌ Not configured"
		keyDisplay := ""
		if p.apiKey != "" {
			status = "✅ Configured"
			keyDisplay = maskAPIKey(p.apiKey)
		}
		
		defaultMark := ""
		if p.provider == cfg.Provider {
			defaultMark = " (default)"
		}
		
		fmt.Printf("  %s%s\n", p.name, defaultMark)
		fmt.Printf("    Status: %s\n", status)
		if keyDisplay != "" {
			fmt.Printf("    Key:    %s\n", keyDisplay)
		}
		fmt.Printf("    Env:    %s\n", p.envKey)
		fmt.Println()
	}

	// Settings
	fmt.Println("⚙️  Settings:")
	fmt.Printf("  Default Provider: %s\n", cfg.Provider)
	fmt.Printf("  Max Tokens:       %d\n", cfg.MaxTokens)
	fmt.Printf("  Temperature:      %.1f\n", cfg.Temperature)
	fmt.Println()

	// Environment
	fmt.Println("🌐 Environment:")
	fmt.Printf("  LANG: %s\n", getEnvOrDefault("LANG", "(not set)"))
	fmt.Printf("  TERM: %s\n", getEnvOrDefault("TERM", "(not set)"))
	fmt.Println()

	// Available Commands
	fmt.Println("📝 Available Commands:")
	commands := []string{
		"summarize  - Summarize text",
		"translate  - Translate text",
		"proofread  - Proofread text",
		"style      - Transform writing style",
		"review     - Code review",
		"commit     - Generate commit message",
		"doc        - Generate documentation",
		"explain    - Explain errors",
		"ask        - Ask AI anything",
	}
	for _, c := range commands {
		fmt.Printf("  %s\n", c)
	}
	fmt.Println()

	// Quick start
	if cfg.OpenAIKey == "" && cfg.ClaudeKey == "" && cfg.GeminiKey == "" {
		fmt.Println("⚠️  No API keys configured!")
		fmt.Println()
		fmt.Println("Quick Start:")
		fmt.Println("  export OPENAI_API_KEY=\"sk-...\"")
		fmt.Println("  # or")
		fmt.Println("  export ANTHROPIC_API_KEY=\"sk-ant-...\"")
		fmt.Println("  # or")
		fmt.Println("  export GOOGLE_API_KEY=\"...\"")
		fmt.Println()
	} else {
		fmt.Println("✅ Ready to use!")
		fmt.Println()
		fmt.Println("Try:")
		fmt.Printf("  aitxt ask \"Hello, world!\"\n")
	}

	return nil
}

func maskAPIKey(key string) string {
	if len(key) <= 8 {
		return "****"
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func getEnvOrDefault(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}
