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
	docProvider string
	docStream   bool
	docLang     string
	docFormat   string
	docType     string
)

var docCmd = &cobra.Command{
	Use:   "doc [file]",
	Short: "Generate documentation from code using AI",
	Long: `Generate documentation from source code using AI.

Documentation types:
  auto      - Auto-detect appropriate documentation (default)
  readme    - README.md style documentation
  api       - API reference documentation
  inline    - Inline code comments
  tutorial  - Tutorial/guide style
  changelog - Changelog entry

Output formats:
  markdown  - Markdown format (default)
  plain     - Plain text
  html      - HTML format

Examples:
  aitxt doc main.go
  aitxt doc src/api.py --type api
  aitxt doc handler.js --type readme --lang ja
  aitxt doc lib/ --type tutorial
  cat code.go | aitxt doc - --format markdown`,
	Args: cobra.MaximumNArgs(1),
	RunE: runDoc,
}

func init() {
	rootCmd.AddCommand(docCmd)

	docCmd.Flags().StringVarP(&docProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	docCmd.Flags().BoolVarP(&docStream, "stream", "s", false, "Enable streaming output")
	docCmd.Flags().StringVarP(&docLang, "lang", "l", "en", "Output language (en, ja, zh, ko, th)")
	docCmd.Flags().StringVarP(&docFormat, "format", "f", "markdown", "Output format (markdown, plain, html)")
	docCmd.Flags().StringVarP(&docType, "type", "t", "auto", "Documentation type (auto, readme, api, inline, tutorial, changelog)")
}

func runDoc(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if docProvider != "" {
		provider = llm.Provider(docProvider)
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
		if info, err := os.Stat(args[0]); err == nil {
			if info.IsDir() {
				// Read all code files in directory
				inputCode, err = readDirectoryCode(args[0])
				if err != nil {
					return fmt.Errorf("failed to read directory: %w", err)
				}
				fileName = args[0]
			} else {
				data, err := os.ReadFile(args[0])
				if err != nil {
					return fmt.Errorf("failed to read file: %w", err)
				}
				inputCode = string(data)
				fileName = args[0]
			}
		} else {
			return fmt.Errorf("file or directory not found: %s", args[0])
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
	progLang := detectLanguageFromFile(fileName)
	typePrompt := getDocTypePrompt(docType)
	formatPrompt := getDocFormatPrompt(docFormat)
	langPrompt := getDocLangPrompt(docLang)

	systemMsg := fmt.Sprintf(`You are an expert technical writer who creates clear, comprehensive documentation.
%s
%s
%s

Create well-structured documentation that is:
- Clear and easy to understand
- Accurate and complete
- Well-organized with appropriate sections
- Includes examples where helpful`, typePrompt, formatPrompt, langPrompt)

	prompt := fmt.Sprintf("Generate documentation for the following %s code:\n\n```%s\n%s\n```", progLang, progLang, inputCode)

	req := &llm.Request{
		Prompt:    prompt,
		SystemMsg: systemMsg,
		MaxTokens: 3000,
	}

	// Generate documentation
	fmt.Printf("Generating %s documentation with %s...\n\n", docType, provider)

	if docStream {
		err = client.Stream(ctx, req, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		fmt.Println()
	} else {
		resp, err := client.Generate(ctx, req)
		if err != nil {
			return fmt.Errorf("failed to generate documentation: %w", err)
		}
		fmt.Println(resp.Text)
		fmt.Printf("\n[Tokens used: %d]\n", resp.TokensUsed)
	}

	if err != nil {
		return fmt.Errorf("failed to generate documentation: %w", err)
	}

	return nil
}

func readDirectoryCode(dir string) (string, error) {
	var codeFiles []string
	extensions := []string{".go", ".py", ".js", ".ts", ".java", ".rb", ".rs", ".cpp", ".c", ".cs"}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := strings.ToLower(filepath.Ext(path))
		for _, validExt := range extensions {
			if ext == validExt {
				data, err := os.ReadFile(path)
				if err != nil {
					return err
				}
				relPath, _ := filepath.Rel(dir, path)
				codeFiles = append(codeFiles, fmt.Sprintf("// File: %s\n%s", relPath, string(data)))
				break
			}
		}
		return nil
	})

	if err != nil {
		return "", err
	}

	return strings.Join(codeFiles, "\n\n---\n\n"), nil
}

func getDocTypePrompt(docType string) string {
	types := map[string]string{
		"readme":    "Create a README.md style documentation with project overview, installation, usage, and examples.",
		"api":       "Create API reference documentation with function signatures, parameters, return values, and examples.",
		"inline":    "Create inline code comments explaining the logic, with JSDoc/GoDoc/docstring style annotations.",
		"tutorial":  "Create a tutorial/guide style documentation that walks through the code step by step.",
		"changelog": "Create a changelog entry describing the changes, additions, and fixes in this code.",
	}

	if prompt, ok := types[docType]; ok {
		return prompt
	}
	return "Create appropriate documentation based on the code structure and purpose."
}

func getDocFormatPrompt(format string) string {
	switch format {
	case "plain":
		return "Output in plain text format without markup."
	case "html":
		return "Output in HTML format with appropriate tags."
	default:
		return "Output in Markdown format with proper headers, code blocks, and formatting."
	}
}

func getDocLangPrompt(lang string) string {
	switch lang {
	case "ja":
		return "Write the documentation in Japanese."
	case "zh":
		return "Write the documentation in Chinese."
	case "ko":
		return "Write the documentation in Korean."
	case "th":
		return "Write the documentation in Thai."
	default:
		return "Write the documentation in English."
	}
}
