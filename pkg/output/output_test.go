package output

import (
	"strings"
	"testing"
)

func TestNewFormatter(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   Format
	}{
		{"JSON format", "json", FormatJSON},
		{"YAML format", "yaml", FormatYAML},
		{"Text format", "text", FormatText},
		{"Invalid format defaults to text", "invalid", FormatText},
		{"Empty format defaults to text", "", FormatText},
		{"Case insensitive", "JSON", FormatJSON},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(tt.format)
			if formatter.format != tt.want {
				t.Errorf("NewFormatter(%s) format = %v, want %v", tt.format, formatter.format, tt.want)
			}
		})
	}
}

func TestIsStructured(t *testing.T) {
	tests := []struct {
		name   string
		format string
		want   bool
	}{
		{"JSON is structured", "json", true},
		{"YAML is structured", "yaml", true},
		{"Text is not structured", "text", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := NewFormatter(tt.format)
			if got := formatter.IsStructured(); got != tt.want {
				t.Errorf("IsStructured() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFormatJSON(t *testing.T) {
	formatter := NewFormatter("json")
	resp := &Response{
		Success:  true,
		Provider: "openai",
		Model:    "gpt-4",
		Content:  "Hello, world!",
		Tokens:   10,
	}

	output, err := formatter.Format(resp)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	if !strings.Contains(output, `"success": true`) {
		t.Error("JSON should contain success field")
	}
	if !strings.Contains(output, `"provider": "openai"`) {
		t.Error("JSON should contain provider field")
	}
	if !strings.Contains(output, `"content": "Hello, world!"`) {
		t.Error("JSON should contain content field")
	}
}

func TestFormatYAML(t *testing.T) {
	formatter := NewFormatter("yaml")
	resp := &Response{
		Success:  true,
		Provider: "claude",
		Content:  "Test content",
	}

	output, err := formatter.Format(resp)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	if !strings.Contains(output, "success: true") {
		t.Error("YAML should contain success field")
	}
	if !strings.Contains(output, "provider: claude") {
		t.Error("YAML should contain provider field")
	}
}

func TestFormatText(t *testing.T) {
	formatter := NewFormatter("text")
	resp := &Response{
		Content: "This is text content",
	}

	output, err := formatter.Format(resp)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	if output != "This is text content" {
		t.Errorf("Text format should return content only, got: %s", output)
	}
}

func TestPrintResponse(t *testing.T) {
	formatter := NewFormatter("text")
	resp := &Response{
		Success:  true,
		Provider: "openai",
		Content:  "Test response",
		Tokens:   15,
	}

	err := formatter.PrintResponse(resp)
	if err != nil {
		t.Errorf("PrintResponse() failed: %v", err)
	}
}

func TestErrorResponse(t *testing.T) {
	resp := ErrorResponse("openai", &testError{"test error"})

	if resp.Success {
		t.Error("ErrorResponse should have Success=false")
	}
	if resp.Provider != "openai" {
		t.Errorf("Expected Provider='openai', got '%s'", resp.Provider)
	}
	if resp.Error != "test error" {
		t.Errorf("Expected Error='test error', got '%s'", resp.Error)
	}
}

func TestSuccessResponse(t *testing.T) {
	resp := SuccessResponse("claude", "claude-3", "Success content", 20)

	if !resp.Success {
		t.Error("SuccessResponse should have Success=true")
	}
	if resp.Provider != "claude" {
		t.Errorf("Expected Provider='claude', got '%s'", resp.Provider)
	}
	if resp.Model != "claude-3" {
		t.Errorf("Expected Model='claude-3', got '%s'", resp.Model)
	}
	if resp.Content != "Success content" {
		t.Errorf("Expected Content='Success content', got '%s'", resp.Content)
	}
	if resp.Tokens != 20 {
		t.Errorf("Expected Tokens=20, got %d", resp.Tokens)
	}
}

func TestTranslateResponse(t *testing.T) {
	resp := &TranslateResponse{
		Response: Response{
			Success:  true,
			Provider: "openai",
			Content:  "Bonjour",
		},
		SourceLang: "en",
		TargetLang: "fr",
		Original:   "Hello",
	}

	formatter := NewFormatter("json")
	output, err := formatter.Format(resp)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	if !strings.Contains(output, `"source_lang": "en"`) {
		t.Error("JSON should contain source_lang field")
	}
	if !strings.Contains(output, `"target_lang": "fr"`) {
		t.Error("JSON should contain target_lang field")
	}
}

func TestReviewResponse(t *testing.T) {
	resp := &ReviewResponse{
		Response: Response{
			Success:  true,
			Provider: "claude",
			Content:  "Code looks good",
		},
		File:     "main.go",
		Language: "Go",
		Focus:    "security",
		Issues:   []string{"issue1", "issue2"},
	}

	formatter := NewFormatter("yaml")
	output, err := formatter.Format(resp)
	if err != nil {
		t.Fatalf("Format() failed: %v", err)
	}

	if !strings.Contains(output, "file: main.go") {
		t.Error("YAML should contain file field")
	}
	if !strings.Contains(output, "language: Go") {
		t.Error("YAML should contain language field")
	}
}

type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}
