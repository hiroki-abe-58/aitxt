package output

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

// Format represents output format type
type Format string

const (
	FormatText Format = "text"
	FormatJSON Format = "json"
	FormatYAML Format = "yaml"
)

// Response represents a structured response
type Response struct {
	Success  bool   `json:"success" yaml:"success"`
	Provider string `json:"provider" yaml:"provider"`
	Model    string `json:"model,omitempty" yaml:"model,omitempty"`
	Content  string `json:"content" yaml:"content"`
	Tokens   int    `json:"tokens,omitempty" yaml:"tokens,omitempty"`
	Error    string `json:"error,omitempty" yaml:"error,omitempty"`
}

// SummaryResponse for summarize command
type SummaryResponse struct {
	Response `yaml:",inline"`
	Original string `json:"original,omitempty" yaml:"original,omitempty"`
}

// TranslateResponse for translate command
type TranslateResponse struct {
	Response   `yaml:",inline"`
	SourceLang string `json:"source_lang,omitempty" yaml:"source_lang,omitempty"`
	TargetLang string `json:"target_lang" yaml:"target_lang"`
	Original   string `json:"original" yaml:"original"`
}

// ReviewResponse for review command
type ReviewResponse struct {
	Response `yaml:",inline"`
	File     string   `json:"file,omitempty" yaml:"file,omitempty"`
	Language string   `json:"language,omitempty" yaml:"language,omitempty"`
	Focus    string   `json:"focus,omitempty" yaml:"focus,omitempty"`
	Issues   []string `json:"issues,omitempty" yaml:"issues,omitempty"`
}

// Formatter handles output formatting
type Formatter struct {
	format Format
}

// NewFormatter creates a new formatter
func NewFormatter(format string) *Formatter {
	f := Format(strings.ToLower(format))
	switch f {
	case FormatJSON, FormatYAML:
		return &Formatter{format: f}
	default:
		return &Formatter{format: FormatText}
	}
}

// IsStructured returns true if format is JSON or YAML
func (f *Formatter) IsStructured() bool {
	return f.format == FormatJSON || f.format == FormatYAML
}

// Format formats the data according to the formatter's format
func (f *Formatter) Format(data interface{}) (string, error) {
	switch f.format {
	case FormatJSON:
		return f.toJSON(data)
	case FormatYAML:
		return f.toYAML(data)
	default:
		return f.toText(data)
	}
}

// Print prints the formatted data
func (f *Formatter) Print(data interface{}) error {
	output, err := f.Format(data)
	if err != nil {
		return err
	}
	fmt.Println(output)
	return nil
}

// PrintResponse prints a standard response
func (f *Formatter) PrintResponse(resp *Response) error {
	if f.format == FormatText {
		fmt.Println(resp.Content)
		if resp.Tokens > 0 {
			fmt.Printf("\n[%s | Tokens: %d]\n", resp.Provider, resp.Tokens)
		}
		return nil
	}
	return f.Print(resp)
}

func (f *Formatter) toJSON(data interface{}) (string, error) {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal JSON: %w", err)
	}
	return string(bytes), nil
}

func (f *Formatter) toYAML(data interface{}) (string, error) {
	bytes, err := yaml.Marshal(data)
	if err != nil {
		return "", fmt.Errorf("failed to marshal YAML: %w", err)
	}
	return string(bytes), nil
}

func (f *Formatter) toText(data interface{}) (string, error) {
	switch v := data.(type) {
	case *Response:
		return v.Content, nil
	case string:
		return v, nil
	default:
		return fmt.Sprintf("%v", data), nil
	}
}

// ErrorResponse creates an error response
func ErrorResponse(provider string, err error) *Response {
	return &Response{
		Success:  false,
		Provider: provider,
		Error:    err.Error(),
	}
}

// SuccessResponse creates a success response
func SuccessResponse(provider, model, content string, tokens int) *Response {
	return &Response{
		Success:  true,
		Provider: provider,
		Model:    model,
		Content:  content,
		Tokens:   tokens,
	}
}
