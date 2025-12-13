package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ProviderSettings holds the settings for a specific provider
type ProviderSettings struct {
	Temperature float64 `json:"temperature"`
	TopP        float64 `json:"top_p,omitempty"`
	TopK        int     `json:"top_k,omitempty"`
	MaxTokens   int     `json:"max_tokens"`
	Model       string  `json:"model,omitempty"`
}

// Settings holds all provider settings
type Settings struct {
	OpenAI ProviderSettings `json:"openai"`
	Gemini ProviderSettings `json:"gemini"`
	Claude ProviderSettings `json:"claude"`
}

// DefaultSettings returns the default settings for all providers
func DefaultSettings() *Settings {
	return &Settings{
		OpenAI: ProviderSettings{
			Temperature: 0.7,
			TopP:        1.0,
			TopK:        0, // Not supported by OpenAI
			MaxTokens:   2000,
			Model:       "gpt-4o",
		},
		Gemini: ProviderSettings{
			Temperature: 0.7,
			TopP:        0.95,
			TopK:        40,
			MaxTokens:   2000,
			Model:       "gemini-1.5-flash",
		},
		Claude: ProviderSettings{
			Temperature: 0.7,
			TopP:        0.9,
			TopK:        0, // 0 means disabled
			MaxTokens:   2000,
			Model:       "claude-sonnet-4-20250514",
		},
	}
}

// getSettingsPath returns the path to the settings file
func getSettingsPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".aitxt")
	return filepath.Join(configDir, "settings.json"), nil
}

// ensureConfigDir ensures the config directory exists
func ensureConfigDir() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get home directory: %w", err)
	}

	configDir := filepath.Join(homeDir, ".aitxt")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// LoadSettings loads settings from the config file
func LoadSettings() (*Settings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return DefaultSettings(), nil
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if os.IsNotExist(err) {
			// Return default settings if file doesn't exist
			return DefaultSettings(), nil
		}
		return nil, fmt.Errorf("failed to read settings file: %w", err)
	}

	var settings Settings
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse settings file: %w", err)
	}

	// Apply defaults for any missing values
	defaults := DefaultSettings()
	applyDefaults(&settings, defaults)

	return &settings, nil
}

// SaveSettings saves settings to the config file
func SaveSettings(settings *Settings) error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	settingsPath, err := getSettingsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %w", err)
	}

	if err := os.WriteFile(settingsPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// applyDefaults applies default values to settings with zero values
func applyDefaults(settings *Settings, defaults *Settings) {
	// OpenAI
	if settings.OpenAI.MaxTokens == 0 {
		settings.OpenAI.MaxTokens = defaults.OpenAI.MaxTokens
	}
	if settings.OpenAI.TopP == 0 {
		settings.OpenAI.TopP = defaults.OpenAI.TopP
	}
	if settings.OpenAI.Model == "" {
		settings.OpenAI.Model = defaults.OpenAI.Model
	}

	// Gemini
	if settings.Gemini.MaxTokens == 0 {
		settings.Gemini.MaxTokens = defaults.Gemini.MaxTokens
	}
	if settings.Gemini.TopP == 0 {
		settings.Gemini.TopP = defaults.Gemini.TopP
	}
	if settings.Gemini.TopK == 0 {
		settings.Gemini.TopK = defaults.Gemini.TopK
	}
	if settings.Gemini.Model == "" {
		settings.Gemini.Model = defaults.Gemini.Model
	}

	// Claude
	if settings.Claude.MaxTokens == 0 {
		settings.Claude.MaxTokens = defaults.Claude.MaxTokens
	}
	if settings.Claude.TopP == 0 {
		settings.Claude.TopP = defaults.Claude.TopP
	}
	if settings.Claude.Model == "" {
		settings.Claude.Model = defaults.Claude.Model
	}
}

// GetProviderSettings returns settings for a specific provider
func (s *Settings) GetProviderSettings(provider string) *ProviderSettings {
	switch provider {
	case "openai":
		return &s.OpenAI
	case "gemini":
		return &s.Gemini
	case "claude":
		return &s.Claude
	default:
		return &s.OpenAI
	}
}

// ValidateProviderSettings validates settings for a specific provider
func ValidateProviderSettings(provider string, settings *ProviderSettings) error {
	switch provider {
	case "openai":
		if settings.Temperature < 0 || settings.Temperature > 2 {
			return fmt.Errorf("openai temperature must be between 0 and 2")
		}
		if settings.TopP < 0 || settings.TopP > 1 {
			return fmt.Errorf("openai top_p must be between 0 and 1")
		}
	case "gemini":
		if settings.Temperature < 0 || settings.Temperature > 2 {
			return fmt.Errorf("gemini temperature must be between 0 and 2")
		}
		if settings.TopP < 0 || settings.TopP > 1 {
			return fmt.Errorf("gemini top_p must be between 0 and 1")
		}
		if settings.TopK < 0 || settings.TopK > 100 {
			return fmt.Errorf("gemini top_k must be between 0 and 100")
		}
	case "claude":
		if settings.Temperature < 0 || settings.Temperature > 1 {
			return fmt.Errorf("claude temperature must be between 0 and 1")
		}
		if settings.TopP < 0 || settings.TopP > 1 {
			return fmt.Errorf("claude top_p must be between 0 and 1")
		}
	}

	if settings.MaxTokens < 1 {
		return fmt.Errorf("max_tokens must be at least 1")
	}

	return nil
}
