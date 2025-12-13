package llm

import (
	"context"
	"fmt"
	"io"

	"github.com/sashabaranov/go-openai"
)

// OpenAIClient implements the Client interface for OpenAI
type OpenAIClient struct {
	client *openai.Client
	config *Config
}

// NewOpenAIClient creates a new OpenAI client
func NewOpenAIClient(config *Config) (*OpenAIClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if config.Model == "" {
		config.Model = openai.GPT4o
	}

	client := openai.NewClient(config.APIKey)

	return &OpenAIClient{
		client: client,
		config: config,
	}, nil
}

// Generate sends a text generation request to OpenAI
func (c *OpenAIClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Prompt,
		},
	}

	if req.SystemMsg != "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: req.SystemMsg,
			},
		}, messages...)
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = c.config.MaxTokens
	}

	temperature := float32(req.Temperature)
	if temperature == 0 {
		temperature = float32(c.config.Temperature)
	}

	topP := float32(c.config.TopP)
	if topP == 0 {
		topP = 1.0
	}

	chatReq := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
	}

	resp, err := c.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, fmt.Errorf("OpenAI API error: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, ErrInvalidResponse
	}

	return &Response{
		Text:         resp.Choices[0].Message.Content,
		Provider:     ProviderOpenAI,
		Model:        resp.Model,
		TokensUsed:   resp.Usage.TotalTokens,
		FinishReason: string(resp.Choices[0].FinishReason),
	}, nil
}

// Stream sends a streaming text generation request to OpenAI
func (c *OpenAIClient) Stream(ctx context.Context, req *Request, callback func(chunk string) error) error {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: req.Prompt,
		},
	}

	if req.SystemMsg != "" {
		messages = append([]openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: req.SystemMsg,
			},
		}, messages...)
	}

	maxTokens := req.MaxTokens
	if maxTokens == 0 {
		maxTokens = c.config.MaxTokens
	}

	temperature := float32(req.Temperature)
	if temperature == 0 {
		temperature = float32(c.config.Temperature)
	}

	topP := float32(c.config.TopP)
	if topP == 0 {
		topP = 1.0
	}

	chatReq := openai.ChatCompletionRequest{
		Model:       c.config.Model,
		Messages:    messages,
		MaxTokens:   maxTokens,
		Temperature: temperature,
		TopP:        topP,
		Stream:      true,
	}

	stream, err := c.client.CreateChatCompletionStream(ctx, chatReq)
	if err != nil {
		return fmt.Errorf("OpenAI streaming error: %w", err)
	}
	defer stream.Close()

	for {
		response, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("stream receive error: %w", err)
		}

		if len(response.Choices) > 0 {
			chunk := response.Choices[0].Delta.Content
			if chunk != "" {
				if err := callback(chunk); err != nil {
					return err
				}
			}
		}
	}

	return nil
}

// GetProvider returns the provider name
func (c *OpenAIClient) GetProvider() Provider {
	return ProviderOpenAI
}

// GetModel returns the model name
func (c *OpenAIClient) GetModel() string {
	return c.config.Model
}

// Validate checks if the client is properly configured
func (c *OpenAIClient) Validate() error {
	return c.config.Validate()
}
