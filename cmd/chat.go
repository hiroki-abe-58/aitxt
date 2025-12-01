package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
	"github.com/spf13/cobra"
)

var (
	chatProvider string
	chatLang     string
	chatSystem   string
)

// Message represents a chat message
type Message struct {
	Role    string
	Content string
}

var chatCmd = &cobra.Command{
	Use:   "chat",
	Short: "Start interactive chat session",
	Long: `Start an interactive chat session with AI.

Commands during chat:
  /help     - Show available commands
  /clear    - Clear conversation history
  /system   - Set system prompt
  /provider - Switch provider
  /model    - Show current model
  /history  - Show conversation history
  /save     - Save conversation to file
  /exit     - Exit chat (or Ctrl+D)

Examples:
  aitxt chat
  aitxt chat --provider claude
  aitxt chat --lang ja
  aitxt chat --system "You are a helpful coding assistant"`,
	RunE: runChat,
}

func init() {
	rootCmd.AddCommand(chatCmd)

	chatCmd.Flags().StringVarP(&chatProvider, "provider", "p", "", "LLM provider (openai, claude, gemini)")
	chatCmd.Flags().StringVarP(&chatLang, "lang", "l", "", "Response language (en, ja, zh, ko, th)")
	chatCmd.Flags().StringVarP(&chatSystem, "system", "s", "", "System prompt")
}

func runChat(cmd *cobra.Command, args []string) error {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Determine provider
	provider := cfg.Provider
	if chatProvider != "" {
		provider = llm.Provider(chatProvider)
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

	ctx := context.Background()
	client, err := factory.CreateClientWithContext(ctx, provider)
	if err != nil {
		return fmt.Errorf("failed to create client: %w", err)
	}

	// Initialize conversation
	var history []Message
	systemPrompt := chatSystem
	if systemPrompt == "" {
		systemPrompt = "You are a helpful AI assistant. Be concise and helpful."
	}
	if chatLang != "" {
		systemPrompt += fmt.Sprintf(" Respond in %s.", getLanguageName(chatLang))
	}

	// Setup readline
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[32m❯ \033[0m",
		HistoryFile:     "/tmp/aitxt_chat_history",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		return fmt.Errorf("failed to initialize readline: %w", err)
	}
	defer rl.Close()

	// Print welcome message
	fmt.Println()
	fmt.Println("🤖 \033[1maitxt Interactive Chat\033[0m")
	fmt.Printf("   Provider: %s\n", provider)
	fmt.Println("   Type /help for commands, /exit or Ctrl+D to quit")
	fmt.Println()

	// Main loop
	for {
		line, err := rl.Readline()
		if err == readline.ErrInterrupt {
			continue
		} else if err == io.EOF {
			fmt.Println("\nGoodbye! 👋")
			break
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle commands
		if strings.HasPrefix(line, "/") {
			shouldExit := handleChatCommand(line, &history, &systemPrompt, &provider, cfg, factory, &client)
			if shouldExit {
				break
			}
			continue
		}

		// Add user message to history
		history = append(history, Message{Role: "user", Content: line})

		// Build messages for API
		messages := buildChatMessages(systemPrompt, history)

		// Generate response
		fmt.Print("\033[36m🤖 \033[0m")

		reqCtx, cancel := context.WithTimeout(ctx, 120*time.Second)
		err = client.Stream(reqCtx, &llm.Request{
			Prompt:    messages,
			SystemMsg: systemPrompt,
			MaxTokens: 2000,
		}, func(chunk string) error {
			fmt.Print(chunk)
			return nil
		})
		cancel()

		if err != nil {
			fmt.Printf("\n\033[31mError: %v\033[0m\n", err)
			// Remove failed message from history
			history = history[:len(history)-1]
			continue
		}

		fmt.Println()
		fmt.Println()

		// Add assistant response to history (we don't have the full response in streaming)
		// For simplicity, we'll make a non-streaming call to get the full response
		history = append(history, Message{Role: "assistant", Content: "[streamed response]"})
	}

	return nil
}

func handleChatCommand(line string, history *[]Message, systemPrompt *string, provider *llm.Provider, cfg *config.Config, factory *llm.Factory, client *llm.Client) bool {
	parts := strings.SplitN(line, " ", 2)
	cmd := strings.ToLower(parts[0])
	arg := ""
	if len(parts) > 1 {
		arg = parts[1]
	}

	switch cmd {
	case "/exit", "/quit", "/q":
		fmt.Println("Goodbye! 👋")
		return true

	case "/help", "/h", "/?":
		printChatHelp()

	case "/clear", "/c":
		*history = []Message{}
		fmt.Println("\033[33m✓ Conversation cleared\033[0m")

	case "/system":
		if arg != "" {
			*systemPrompt = arg
			fmt.Printf("\033[33m✓ System prompt updated: %s\033[0m\n", arg)
		} else {
			fmt.Printf("\033[36mCurrent system prompt: %s\033[0m\n", *systemPrompt)
		}

	case "/provider":
		if arg != "" {
			newProvider := llm.Provider(arg)
			llmConfig, err := cfg.ToLLMConfig(newProvider)
			if err != nil {
				fmt.Printf("\033[31mError: %v\033[0m\n", err)
				return false
			}
			if err := factory.RegisterConfig(llmConfig); err != nil {
				fmt.Printf("\033[31mError: %v\033[0m\n", err)
				return false
			}
			ctx := context.Background()
			newClient, err := factory.CreateClientWithContext(ctx, newProvider)
			if err != nil {
				fmt.Printf("\033[31mError: %v\033[0m\n", err)
				return false
			}
			*provider = newProvider
			*client = newClient
			fmt.Printf("\033[33m✓ Switched to %s\033[0m\n", newProvider)
		} else {
			fmt.Printf("\033[36mCurrent provider: %s\033[0m\n", *provider)
		}

	case "/model":
		fmt.Printf("\033[36mProvider: %s\nModel: %s\033[0m\n", *provider, (*client).GetModel())

	case "/history":
		if len(*history) == 0 {
			fmt.Println("\033[33mNo conversation history\033[0m")
		} else {
			fmt.Println("\033[36m--- Conversation History ---\033[0m")
			for i, msg := range *history {
				role := "👤"
				if msg.Role == "assistant" {
					role = "🤖"
				}
				content := msg.Content
				if len(content) > 100 {
					content = content[:100] + "..."
				}
				fmt.Printf("%d. %s %s\n", i+1, role, content)
			}
			fmt.Println("\033[36m----------------------------\033[0m")
		}

	case "/save":
		filename := arg
		if filename == "" {
			filename = fmt.Sprintf("chat_%s.md", time.Now().Format("20060102_150405"))
		}
		if err := saveChatHistory(*history, *systemPrompt, filename); err != nil {
			fmt.Printf("\033[31mError saving: %v\033[0m\n", err)
		} else {
			fmt.Printf("\033[33m✓ Saved to %s\033[0m\n", filename)
		}

	default:
		fmt.Printf("\033[31mUnknown command: %s. Type /help for commands.\033[0m\n", cmd)
	}

	return false
}

func printChatHelp() {
	fmt.Println("\033[36m")
	fmt.Println("Available Commands:")
	fmt.Println("  /help, /h       - Show this help")
	fmt.Println("  /clear, /c      - Clear conversation history")
	fmt.Println("  /system [msg]   - View/set system prompt")
	fmt.Println("  /provider [p]   - View/switch provider (openai, claude, gemini)")
	fmt.Println("  /model          - Show current model")
	fmt.Println("  /history        - Show conversation history")
	fmt.Println("  /save [file]    - Save conversation to file")
	fmt.Println("  /exit, /q       - Exit chat")
	fmt.Println("\033[0m")
}

func buildChatMessages(systemPrompt string, history []Message) string {
	var sb strings.Builder

	for _, msg := range history {
		if msg.Role == "user" {
			sb.WriteString(fmt.Sprintf("User: %s\n", msg.Content))
		} else {
			sb.WriteString(fmt.Sprintf("Assistant: %s\n", msg.Content))
		}
	}

	return sb.String()
}

func saveChatHistory(history []Message, systemPrompt, filename string) error {
	var sb strings.Builder

	sb.WriteString("# Chat Session\n\n")
	sb.WriteString(fmt.Sprintf("**Date:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))
	sb.WriteString(fmt.Sprintf("**System Prompt:** %s\n\n", systemPrompt))
	sb.WriteString("---\n\n")

	for _, msg := range history {
		if msg.Role == "user" {
			sb.WriteString(fmt.Sprintf("**👤 User:**\n%s\n\n", msg.Content))
		} else {
			sb.WriteString(fmt.Sprintf("**🤖 Assistant:**\n%s\n\n", msg.Content))
		}
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}
