package llm

import (
	"context"
	"errors"
)

// Provider represents the LLM provider type
type Provider string

const (
	ProviderOpenAI   Provider = "openai"
	ProviderClaude   Provider = "claude"
	ProviderGemini   Provider = "gemini"
	ProviderUnknown  Provider = "unknown"
)

// Common errors
var (
	ErrInvalidProvider    = errors.New("invalid provider")
	ErrAPIKeyMissing      = errors.New("API key is missing")
	ErrRateLimitExceeded  = errors.New("rate limit exceeded")
	ErrInvalidResponse    = errors.New("invalid response from API")
	ErrContextCanceled    = errors.New("context canceled")
)

// Request represents a text generation request
type Request struct {
	Prompt      string
	MaxTokens   int
	Temperature float64
	SystemMsg   string
	Stream      bool
}

// Response represents a text generation response
type Response struct {
	Text         string
	Provider     Provider
	Model        string
	TokensUsed   int
	FinishReason string
}

// Client is the interface that all LLM providers must implement
type Client interface {
	// Generate sends a text generation request and returns the response
	Generate(ctx context.Context, req *Request) (*Response, error)

	// Stream sends a streaming text generation request
	Stream(ctx context.Context, req *Request, callback func(chunk string) error) error

	// GetProvider returns the provider name
	GetProvider() Provider

	// GetModel returns the model name being used
	GetModel() string

	// Validate checks if the client is properly configured
	Validate() error
}

// Config represents the configuration for an LLM client
type Config struct {
	Provider    Provider
	APIKey      string
	Model       string
	MaxTokens   int
	Temperature float64
	BaseURL     string
}

// DefaultConfig returns a default configuration
func DefaultConfig(provider Provider) *Config {
	return &Config{
		Provider:    provider,
		MaxTokens:   2000,
		Temperature: 0.7,
	}
}

// Validate checks if the configuration is valid
func (c *Config) Validate() error {
	if c.Provider == ProviderUnknown || c.Provider == "" {
		return ErrInvalidProvider
	}
	if c.APIKey == "" {
		return ErrAPIKeyMissing
	}
	if c.MaxTokens <= 0 {
		c.MaxTokens = 2000
	}
	if c.Temperature < 0 || c.Temperature > 2 {
		c.Temperature = 0.7
	}
	return nil
}
