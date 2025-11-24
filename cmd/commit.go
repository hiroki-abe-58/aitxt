package cmd

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	commitProvider string
	commitLang     string
	commitType     string
	commitShort    bool
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Generate commit message from staged changes",
	Long: `Generate a commit message from staged git changes using AI.

This command reads your staged changes (git diff --staged) and generates
a meaningful commit message following Conventional Commits format.

Commit types:
  auto     - Auto-detect type (default)
  feat     - New feature
  fix      - Bug fix
  docs     - Documentation
  style    - Code style changes
  refactor - Code refactoring
  test     - Adding tests
  chore    - Maintenance tasks

Examples:
  aitxt commit
  aitxt commit --lang ja
  aitxt commit --type feat
  aitxt commit --short`,
	RunE: runCommit,
}

func init() {
	rootCmd.AddCommand(commitCmd)

	commitCmd.Flags().StringVarP(&commitProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	commitCmd.Flags().StringVarP(&commitLang, "lang", "l", "en", "Output language (en, ja, zh, ko)")
	commitCmd.Flags().StringVarP(&commitType, "type", "t", "auto", "Commit type (auto, feat, fix, docs, style, refactor, test, chore)")
	commitCmd.Flags().BoolVarP(&commitShort, "short", "s", false, "Generate short message only (no body)")
}

func runCommit(cmd *cobra.Command, args []string) error {
	// Get staged diff
	diff, err := getStagedDiff()
	if err != nil {
		return fmt.Errorf("failed to get staged diff: %w", err)
	}

	if strings.TrimSpace(diff) == "" {
		return fmt.Errorf("no staged changes found. Use 'git add' to stage changes")
	}

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if commitProvider != "" {
		provider = llm.Provider(commitProvider)
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
	systemMsg := buildCommitSystemMsg()
	prompt := buildCommitPrompt(diff)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 500,
	}

	// Generate commit message
	fmt.Printf("Generating commit message with %s...\n\n", provider)

	resp, err := client.Generate(ctx, req)
	if err != nil {
		return fmt.Errorf("failed to generate commit message: %w", err)
	}

	fmt.Println("--- Suggested Commit Message ---")
	fmt.Println(resp.Text)
	fmt.Println("--------------------------------")
	fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	fmt.Println("\nTo use this message, run:")
	fmt.Println("  git commit -m \"<copy message above>\"")

	return nil
}

func getStagedDiff() (string, error) {
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(output), nil
}

func buildCommitSystemMsg() string {
	langPrompt := getCommitLangPrompt(commitLang)
	typePrompt := ""
	if commitType != "auto" {
		typePrompt = fmt.Sprintf("The commit type MUST be '%s'. ", commitType)
	}

	formatPrompt := ""
	if commitShort {
		formatPrompt = "Generate ONLY a single-line commit message (no body)."
	} else {
		formatPrompt = `Generate a commit message with:
1. A concise subject line (50 chars or less)
2. A blank line
3. A body explaining what and why (wrap at 72 chars)`
	}

	return fmt.Sprintf(`You are a helpful assistant that generates git commit messages.
Follow the Conventional Commits format: <type>(<scope>): <description>

Types: feat, fix, docs, style, refactor, test, chore
%s%s
%s

Output ONLY the commit message, nothing else.`, typePrompt, formatPrompt, langPrompt)
}

func buildCommitPrompt(diff string) string {
	// Truncate diff if too long
	maxLen := 8000
	if len(diff) > maxLen {
		diff = diff[:maxLen] + "\n... (truncated)"
	}

	return fmt.Sprintf("Generate a commit message for the following changes:\n\n```diff\n%s\n```", diff)
}

func getCommitLangPrompt(lang string) string {
	switch lang {
	case "ja":
		return "Write the commit message in Japanese."
	case "zh":
		return "Write the commit message in Chinese."
	case "ko":
		return "Write the commit message in Korean."
	default:
		return "Write the commit message in English."
	}
}
