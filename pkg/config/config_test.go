package config

import (
	"os"
	"testing"

	"github.com/hiroki-abe-58/aitxt/pkg/llm"
)

func TestLoad(t *testing.T) {
	// Set up test environment
	os.Setenv("OPENAI_API_KEY", "test-openai-key")
	os.Setenv("ANTHROPIC_API_KEY", "test-anthropic-key")
	os.Setenv("GOOGLE_API_KEY", "test-google-key")
	defer func() {
		os.Unsetenv("OPENAI_API_KEY")
		os.Unsetenv("ANTHROPIC_API_KEY")
		os.Unsetenv("GOOGLE_API_KEY")
	}()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if cfg.OpenAIKey != "test-openai-key" {
		t.Errorf("Expected OpenAIKey='test-openai-key', got '%s'", cfg.OpenAIKey)
	}

	if cfg.ClaudeKey != "test-anthropic-key" {
		t.Errorf("Expected ClaudeKey='test-anthropic-key', got '%s'", cfg.ClaudeKey)
	}

	if cfg.GeminiKey != "test-google-key" {
		t.Errorf("Expected GeminiKey='test-google-key', got '%s'", cfg.GeminiKey)
	}

	if cfg.Provider != llm.ProviderOpenAI {
		t.Errorf("Expected default provider to be OpenAI, got %s", cfg.Provider)
	}
}

func TestLoadNoKeys(t *testing.T) {
	// Clear all keys
	os.Unsetenv("OPENAI_API_KEY")
	os.Unsetenv("ANTHROPIC_API_KEY")
	os.Unsetenv("GOOGLE_API_KEY")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() should not fail even without keys: %v", err)
	}

	if cfg.OpenAIKey != "" {
		t.Errorf("Expected empty OpenAIKey, got '%s'", cfg.OpenAIKey)
	}
}

func TestToLLMConfig(t *testing.T) {
	cfg := &Config{
		OpenAIKey:   "test-openai",
		ClaudeKey:   "test-claude",
		GeminiKey:   "test-gemini",
		Provider:    llm.ProviderOpenAI,
		MaxTokens:   2000,
		Temperature: 0.8,
	}

	tests := []struct {
		name     string
		provider llm.Provider
		wantErr  bool
	}{
		{"OpenAI", llm.ProviderOpenAI, false},
		{"Claude", llm.ProviderClaude, false},
		{"Gemini", llm.ProviderGemini, false},
		{"Invalid", llm.Provider("invalid"), true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			llmCfg, err := cfg.ToLLMConfig(tt.provider)
			if (err != nil) != tt.wantErr {
				t.Errorf("ToLLMConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if llmCfg.MaxTokens != cfg.MaxTokens {
					t.Errorf("Expected MaxTokens=%d, got %d", cfg.MaxTokens, llmCfg.MaxTokens)
				}
				if llmCfg.Temperature != cfg.Temperature {
					t.Errorf("Expected Temperature=%f, got %f", cfg.Temperature, llmCfg.Temperature)
				}
			}
		})
	}
}

func TestToLLMConfigMissingKey(t *testing.T) {
	cfg := &Config{
		OpenAIKey: "",
		ClaudeKey: "",
		GeminiKey: "",
		Provider:  llm.ProviderOpenAI,
	}

	_, err := cfg.ToLLMConfig(llm.ProviderOpenAI)
	if err == nil {
		t.Error("ToLLMConfig() should fail when API key is missing")
	}
}

func TestDetermineDefaultProvider(t *testing.T) {
	tests := []struct {
		name   string
		openai string
		claude string
		gemini string
		want   llm.Provider
	}{
		{"OpenAI only", "key1", "", "", llm.ProviderOpenAI},
		{"Claude only", "", "key2", "", llm.ProviderClaude},
		{"Gemini only", "", "", "key3", llm.ProviderGemini},
		{"All keys - OpenAI priority", "key1", "key2", "key3", llm.ProviderOpenAI},
		{"No keys", "", "", "", llm.ProviderOpenAI},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				OpenAIKey: tt.openai,
				ClaudeKey: tt.claude,
				GeminiKey: tt.gemini,
			}

			got := cfg.determineDefaultProvider()
			if got != tt.want {
				t.Errorf("determineDefaultProvider() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConfigDefaults(t *testing.T) {
	cfg, _ := Load()

	if cfg.MaxTokens != 2000 {
		t.Errorf("Expected default MaxTokens=2000, got %d", cfg.MaxTokens)
	}

	if cfg.Temperature != 0.7 {
		t.Errorf("Expected default Temperature=0.7, got %f", cfg.Temperature)
	}
}
