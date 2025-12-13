package llm

import (
	"context"
	"fmt"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type ClaudeClient struct {
	client anthropic.Client
	config *Config
}

func NewClaudeClient(config *Config) (*ClaudeClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}
	if config.Model == "" {
		config.Model = "claude-sonnet-4-20250514"
	}
	client := anthropic.NewClient(option.WithAPIKey(config.APIKey))
	return &ClaudeClient{client: client, config: config}, nil
}

func (c *ClaudeClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int64(c.config.MaxTokens)
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = c.config.Temperature
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.config.Model),
		MaxTokens: maxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		},
	}

	// Set temperature if specified
	if temperature > 0 {
		params.Temperature = anthropic.Float(temperature)
	}

	// Set top_p if specified
	if c.config.TopP > 0 {
		params.TopP = anthropic.Float(c.config.TopP)
	}

	// Set top_k if specified
	if c.config.TopK > 0 {
		params.TopK = anthropic.Int(int64(c.config.TopK))
	}

	if req.SystemMsg != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemMsg},
		}
	}

	message, err := c.client.Messages.New(ctx, params)
	if err != nil {
		return nil, fmt.Errorf("Claude API error: %w", err)
	}

	var text string
	for _, block := range message.Content {
		if block.Type == "text" {
			text += block.Text
		}
	}

	return &Response{
		Text:       text,
		Provider:   ProviderClaude,
		Model:      string(message.Model),
		TokensUsed: int(message.Usage.InputTokens + message.Usage.OutputTokens),
	}, nil
}

func (c *ClaudeClient) Stream(ctx context.Context, req *Request, callback func(chunk string) error) error {
	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = int64(c.config.MaxTokens)
	}

	temperature := req.Temperature
	if temperature == 0 {
		temperature = c.config.Temperature
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.Model(c.config.Model),
		MaxTokens: maxTokens,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(req.Prompt)),
		},
	}

	// Set temperature if specified
	if temperature > 0 {
		params.Temperature = anthropic.Float(temperature)
	}

	// Set top_p if specified
	if c.config.TopP > 0 {
		params.TopP = anthropic.Float(c.config.TopP)
	}

	// Set top_k if specified
	if c.config.TopK > 0 {
		params.TopK = anthropic.Int(int64(c.config.TopK))
	}

	if req.SystemMsg != "" {
		params.System = []anthropic.TextBlockParam{
			{Text: req.SystemMsg},
		}
	}

	stream := c.client.Messages.NewStreaming(ctx, params)
	for stream.Next() {
		event := stream.Current()
		if event.Delta.Text != "" {
			if err := callback(event.Delta.Text); err != nil {
				return err
			}
		}
	}
	return stream.Err()
}

func (c *ClaudeClient) GetProvider() Provider { return ProviderClaude }
func (c *ClaudeClient) GetModel() string      { return c.config.Model }
func (c *ClaudeClient) Validate() error       { return c.config.Validate() }
