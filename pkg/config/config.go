package config

import (
	"fmt"
	"os"

	"github.com/hiroki-abe-58/aitxt/pkg/llm"
)

// Config holds the application configuration
type Config struct {
	Provider    llm.Provider
	OpenAIKey   string
	ClaudeKey   string
	GeminiKey   string
	Model       string
	MaxTokens   int
	Temperature float64
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Provider:    llm.ProviderOpenAI, // default
		OpenAIKey:   os.Getenv("OPENAI_API_KEY"),
		ClaudeKey:   os.Getenv("ANTHROPIC_API_KEY"),
		GeminiKey:   os.Getenv("GOOGLE_API_KEY"),
		MaxTokens:   2000,
		Temperature: 0.7,
	}

	// Auto-detect provider based on available API keys
	if cfg.OpenAIKey != "" {
		cfg.Provider = llm.ProviderOpenAI
	} else if cfg.ClaudeKey != "" {
		cfg.Provider = llm.ProviderClaude
	} else if cfg.GeminiKey != "" {
		cfg.Provider = llm.ProviderGemini
	}

	return cfg, nil
}

// GetAPIKey returns the API key for the specified provider
func (c *Config) GetAPIKey(provider llm.Provider) (string, error) {
	switch provider {
	case llm.ProviderOpenAI:
		if c.OpenAIKey == "" {
			return "", fmt.Errorf("OPENAI_API_KEY is not set")
		}
		return c.OpenAIKey, nil
	case llm.ProviderClaude:
		if c.ClaudeKey == "" {
			return "", fmt.Errorf("ANTHROPIC_API_KEY is not set")
		}
		return c.ClaudeKey, nil
	case llm.ProviderGemini:
		if c.GeminiKey == "" {
			return "", fmt.Errorf("GOOGLE_API_KEY is not set")
		}
		return c.GeminiKey, nil
	default:
		return "", fmt.Errorf("unknown provider: %s", provider)
	}
}

// ToLLMConfig converts to LLM config
func (c *Config) ToLLMConfig(provider llm.Provider) (*llm.Config, error) {
	apiKey, err := c.GetAPIKey(provider)
	if err != nil {
		return nil, err
	}

	return &llm.Config{
		Provider:    provider,
		APIKey:      apiKey,
		Model:       c.Model,
		MaxTokens:   c.MaxTokens,
		Temperature: c.Temperature,
	}, nil
}
