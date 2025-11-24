package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	appVersion string
)

var rootCmd = &cobra.Command{
	Use:   "aitxt",
	Short: "🤖 AI-powered text processing CLI tool",
	Long: `🤖 aitxt - AI-powered text processing CLI tool

aitxt provides powerful text processing capabilities using multiple LLM providers:
  • Text summarization
  • Translation (multiple languages)
  • Proofreading and grammar correction
  • Style transformation
  • Code review and documentation generation
  • Developer productivity tools

Supported providers: OpenAI, Claude (Anthropic), Google Gemini

Examples:
  aitxt summarize document.txt
  aitxt translate "Hello world" --to ja
  aitxt review code.go
  aitxt commit`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute runs the root command
func Execute(version string) error {
	appVersion = version
	rootCmd.Version = version
	return rootCmd.Execute()
}

func init() {
	rootCmd.SetVersionTemplate(fmt.Sprintf("aitxt version %s\n", appVersion))
}
