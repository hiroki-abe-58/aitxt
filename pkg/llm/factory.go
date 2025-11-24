package llm

import (
	"context"
	"fmt"
)

// Factory creates LLM clients based on the provider
type Factory struct {
	configs map[Provider]*Config
}

// NewFactory creates a new client factory
func NewFactory() *Factory {
	return &Factory{
		configs: make(map[Provider]*Config),
	}
}

// RegisterConfig registers a configuration for a provider
func (f *Factory) RegisterConfig(config *Config) error {
	if err := config.Validate(); err != nil {
		return fmt.Errorf("invalid config: %w", err)
	}
	f.configs[config.Provider] = config
	return nil
}

// CreateClient creates a client for the specified provider
func (f *Factory) CreateClient(provider Provider) (Client, error) {
	return f.CreateClientWithContext(context.Background(), provider)
}

// CreateClientWithContext creates a client with context for the specified provider
func (f *Factory) CreateClientWithContext(ctx context.Context, provider Provider) (Client, error) {
	config, exists := f.configs[provider]
	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, provider)
	}

	switch provider {
	case ProviderOpenAI:
		return NewOpenAIClient(config)
	case ProviderClaude:
		return NewClaudeClient(config)
	case ProviderGemini:
		return NewGeminiClient(ctx, config)
	default:
		return nil, fmt.Errorf("%w: %s", ErrInvalidProvider, provider)
	}
}

// GetRegisteredProviders returns all registered providers
func (f *Factory) GetRegisteredProviders() []Provider {
	providers := make([]Provider, 0, len(f.configs))
	for provider := range f.configs {
		providers = append(providers, provider)
	}
	return providers
}
