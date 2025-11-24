package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// ClaudeClient implements the Client interface for Anthropic Claude
type ClaudeClient struct {
	client *anthropic.Client
	config *Config
}

// NewClaudeClient creates a new Claude client
func NewClaudeClient(config *Config) (*ClaudeClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if config.Model == "" {
		config.Model = "claude-sonnet-4-20250514"
	}

	client := anthropic.NewClient(option.WithAPIKey(config.APIKey))

	return &ClaudeClient{
		client: client,
		config: config,
	}, nil
}

// Generate sends a text generation request to Claude
func (c *ClaudeClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int64(c.config.MaxTokens)
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.F(c.config.Model),
		MaxTokens: anthropic.F(maxTokens),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		}),
	}

	if req.SystemMsg != "" {
		params.System = anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(req.SystemMsg),
		})
	}

	message, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Claude API error: %w", err)
	}

	var text string
	for _, block := range message.Content {
		if block.Type == anthropic.ContentBlockTypeText {
			text += block.Text
		}
	}

	return &Response{
		Text:         text,
		Provider:     ProviderClaude,
		Model:        string(message.Model),
		TokensUsed:   int(message.Usage.InputTokens + message.Usage.OutputTokens),
		FinishReason: string(message.StopReason),
	}, nil
}

// Stream sends a streaming text generation request to Claude
func (c *ClaudeClient) Stream(ctx context.Context, req *Request, callback func(chunk string) error) error {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int64(c.config.MaxTokens)
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.F(c.config.Model),
		MaxTokens: anthropic.F(maxTokens),
		Messages: anthropic.F([]anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		}),
	}

	if req.SystemMsg != "" {
		params.System = anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(req.SystemMsg),
		})
	}

	stream := c.client.Messages.NewStreaming(ctx, params)

	for stream.Next() {
		event := stream.Current()
		switch delta := event.Delta.(type) {
		case anthropic.ContentBlockDeltaEventDelta:
			if delta.Text != "" {
				if err := callback(delta.Text); err != nil {
					return err
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return fmt.Errorf("Claude streaming error: %w", err)
	}

	return nil
}

// GetProvider returns the provider name
func (c *ClaudeClient) GetProvider() Provider {
	return ProviderClaude
}

// GetModel returns the model name
func (c *ClaudeClient) GetModel() string {
	return c.config.Model
}

// Validate checks if the client is properly configured
func (c *ClaudeClient) Validate() error {
	return c.config.Validate()
}
