package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	reviewProvider string
	reviewStream   bool
	reviewLang     string
	reviewFocus    string
)

var reviewCmd = &cobra.Command{
	Use:   "review [file]",
	Short: "Review code using AI",
	Long: `Review code from a file or stdin using AI.

Focus areas:
  all        - Comprehensive review (default)
  security   - Security vulnerabilities
  performance - Performance issues
  style      - Code style and readability
  bugs       - Potential bugs and errors
  test       - Test coverage suggestions

Examples:
  aitxt review main.go
  aitxt review src/app.py --focus security
  aitxt review handler.js --focus performance --lang ja
  cat code.go | aitxt review -
  git diff | aitxt review - --focus bugs`,
	Args: cobra.MaximumNArgs(1),
	RunE: runReview,
}

func init() {
	rootCmd.AddCommand(reviewCmd)

	reviewCmd.Flags().StringVarP(&reviewProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	reviewCmd.Flags().BoolVarP(&reviewStream, "stream", "s", false, "Enable streaming output")
	reviewCmd.Flags().StringVarP(&reviewLang, "lang", "l", "en", "Output language (en, ja, zh, ko, th)")
	reviewCmd.Flags().StringVarP(&reviewFocus, "focus", "f", "all", "Focus area (all, security, performance, style, bugs, test)")
}

func runReview(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if reviewProvider != "" {
		provider = llm.Provider(reviewProvider)
	}

	// Get input code
	var inputCode string
	var fileName string
	if len(args) == 0 || args[0] == "-" {
		data, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		inputCode = string(data)
		fileName = "stdin"
	} else {
		if _, err := os.Stat(args[0]); err == nil {
			data, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			inputCode = string(data)
			fileName = args[0]
		} else {
			return fmt.Errorf("file not found: %s", args[0])
		}
	}

	inputCode = strings.TrimSpace(inputCode)
	if inputCode == "" {
		return fmt.Errorf("no code provided")
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
	langName := detectLanguageFromFile(fileName)
	focusPrompt := getFocusPrompt(reviewFocus)
	langPrompt := getReviewLangPrompt(reviewLang)

	systemMsg := fmt.Sprintf(`You are an expert code reviewer with deep knowledge of software engineering best practices.
%s
%s

Provide a structured code review with:
1. **Summary** - Brief overview of the code
2. **Issues Found** - List of problems (severity: high/medium/low)
3. **Suggestions** - Specific improvements with code examples
4. **Good Practices** - What the code does well

Be constructive and specific. Include line numbers when relevant.`, focusPrompt, langPrompt)

	prompt := fmt.Sprintf("Please review the following %s code:\n\n```%s\n%s\n```", langName, langName, inputCode)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 2500,
	}

	// Generate review
	fmt.Printf("Reviewing %s with %s (focus: %s)...\n\n", fileName, provider, reviewFocus)

	if reviewStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to review code: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to review code: %w", err)
	}

	return nil
}

func detectLanguageFromFile(fileName string) string {
	ext := strings.ToLower(filepath.Ext(fileName))
	languages := map[string]string{
		".go":    "Go",
		".py":    "Python",
		".js":    "JavaScript",
		".ts":    "TypeScript",
		".jsx":   "React JSX",
		".tsx":   "React TSX",
		".java":  "Java",
		".rb":    "Ruby",
		".rs":    "Rust",
		".cpp":   "C++",
		".c":     "C",
		".cs":    "C#",
		".php":   "PHP",
		".swift": "Swift",
		".kt":    "Kotlin",
		".sh":    "Shell",
		".sql":   "SQL",
		".html":  "HTML",
		".css":   "CSS",
	}

	if lang, ok := languages[ext]; ok {
		return lang
	}
	return "code"
}

func getFocusPrompt(focus string) string {
	switch focus {
	case "security":
		return "Focus specifically on security vulnerabilities: SQL injection, XSS, authentication issues, data exposure, input validation, etc."
	case "performance":
		return "Focus specifically on performance issues: algorithmic complexity, memory usage, unnecessary operations, caching opportunities, etc."
	case "style":
		return "Focus specifically on code style and readability: naming conventions, code organization, comments, formatting, etc."
	case "bugs":
		return "Focus specifically on potential bugs: edge cases, null checks, error handling, race conditions, logic errors, etc."
	case "test":
		return "Focus specifically on test coverage: missing tests, test quality, edge cases to test, mocking strategies, etc."
	default:
		return "Provide a comprehensive review covering security, performance, style, bugs, and maintainability."
	}
}

func getReviewLangPrompt(lang string) string {
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
