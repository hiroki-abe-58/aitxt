package llm

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

// GeminiClient implements the Client interface for Google Gemini
type GeminiClient struct {
	client *genai.Client
	model  *genai.GenerativeModel
	config *Config
}

// NewGeminiClient creates a new Gemini client
func NewGeminiClient(ctx context.Context, config *Config) (*GeminiClient, error) {
	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	if config.Model == "" {
		config.Model = "gemini-1.5-flash"
	}

	client, err := genai.NewClient(ctx, option.WithAPIKey(config.APIKey))
	if err != nil {
		return nil, fmt.Errorf("failed to create Gemini client: %w", err)
	}

	model := client.GenerativeModel(config.Model)
	model.SetTemperature(float32(config.Temperature))
	model.SetMaxOutputTokens(int32(config.MaxTokens))
	if config.TopP > 0 {
		model.SetTopP(float32(config.TopP))
	}
	if config.TopK > 0 {
		model.SetTopK(int32(config.TopK))
	}

	return &GeminiClient{
		client: client,
		model:  model,
		config: config,
	}, nil
}

// Generate sends a text generation request to Gemini
func (c *GeminiClient) Generate(ctx context.Context, req *Request) (*Response, error) {
	if req.SystemMsg != "" {
		c.model.SystemInstruction = genai.NewUserContent(genai.Text(req.SystemMsg))
	}

	if req.Temperature != 0 {
		c.model.SetTemperature(float32(req.Temperature))
	}
	if req.MaxTokens != 0 {
		c.model.SetMaxOutputTokens(int32(req.MaxTokens))
	}

	resp, err := c.model.GenerateContent(ctx, genai.Text(req.Prompt))
	if err != nil {
		return nil, fmt.Errorf("Gemini API error: %w", err)
	}

	if len(resp.Candidates) == 0 || len(resp.Candidates[0].Content.Parts) == 0 {
		return nil, ErrInvalidResponse
	}

	var text string
	for _, part := range resp.Candidates[0].Content.Parts {
		if t, ok := part.(genai.Text); ok {
			text += string(t)
		}
	}

	tokensUsed := 0
	if resp.UsageMetadata != nil {
		tokensUsed = int(resp.UsageMetadata.TotalTokenCount)
	}

	finishReason := ""
	if resp.Candidates[0].FinishReason > 0 {
		finishReason = resp.Candidates[0].FinishReason.String()
	}

	return &Response{
		Text:         text,
		Provider:     ProviderGemini,
		Model:        c.config.Model,
		TokensUsed:   tokensUsed,
		FinishReason: finishReason,
	}, nil
}

// Stream sends a streaming text generation request to Gemini
func (c *GeminiClient) Stream(ctx context.Context, req *Request, callback func(chunk string) error) error {
	if req.SystemMsg != "" {
		c.model.SystemInstruction = genai.NewUserContent(genai.Text(req.SystemMsg))
	}

	if req.Temperature != 0 {
		c.model.SetTemperature(float32(req.Temperature))
	}
	if req.MaxTokens != 0 {
		c.model.SetMaxOutputTokens(int32(req.MaxTokens))
	}

	iter := c.model.GenerateContentStream(ctx, genai.Text(req.Prompt))

	for {
		resp, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return fmt.Errorf("Gemini streaming error: %w", err)
		}

		for _, candidate := range resp.Candidates {
			for _, part := range candidate.Content.Parts {
				if t, ok := part.(genai.Text); ok {
					if err := callback(string(t)); err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

// GetProvider returns the provider name
func (c *GeminiClient) GetProvider() Provider {
	return ProviderGemini
}

// GetModel returns the model name
func (c *GeminiClient) GetModel() string {
	return c.config.Model
}

// Validate checks if the client is properly configured
func (c *GeminiClient) Validate() error {
	return c.config.Validate()
}

// Close closes the Gemini client
func (c *GeminiClient) Close() error {
	return c.client.Close()
}
