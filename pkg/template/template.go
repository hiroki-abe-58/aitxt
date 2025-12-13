package template

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"
)

// Template represents a prompt template
type Template struct {
	Name        string            `json:"name"`
	Description string            `json:"description"`
	Command     string            `json:"command"`
	SystemMsg   string            `json:"system_msg"`
	Prompt      string            `json:"prompt"`
	Variables   []string          `json:"variables,omitempty"`
	Defaults    map[string]string `json:"defaults,omitempty"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
}

// Store manages templates
type Store struct {
	dir string
}

// NewStore creates a new template store
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dir := filepath.Join(homeDir, ".aitxt", "templates")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create templates directory: %w", err)
	}

	return &Store{dir: dir}, nil
}

// Add adds a new template
func (s *Store) Add(tmpl *Template) error {
	if tmpl.Name == "" {
		return fmt.Errorf("template name is required")
	}

	now := time.Now()
	if tmpl.CreatedAt.IsZero() {
		tmpl.CreatedAt = now
	}
	tmpl.UpdatedAt = now

	// Extract variables from prompt
	tmpl.Variables = extractVariables(tmpl.Prompt)

	filename := filepath.Join(s.dir, tmpl.Name+".json")
	data, err := json.MarshalIndent(tmpl, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	return os.WriteFile(filename, data, 0644)
}

// Get retrieves a template by name
func (s *Store) Get(name string) (*Template, error) {
	filename := filepath.Join(s.dir, name+".json")
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("template not found: %s", name)
	}

	var tmpl Template
	if err := json.Unmarshal(data, &tmpl); err != nil {
		return nil, fmt.Errorf("failed to parse template: %w", err)
	}

	return &tmpl, nil
}

// List returns all templates
func (s *Store) List() ([]*Template, error) {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read templates directory: %w", err)
	}

	var templates []*Template
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ".json")
		tmpl, err := s.Get(name)
		if err != nil {
			continue
		}
		templates = append(templates, tmpl)
	}

	sort.Slice(templates, func(i, j int) bool {
		return templates[i].Name < templates[j].Name
	})

	return templates, nil
}

// Delete deletes a template
func (s *Store) Delete(name string) error {
	filename := filepath.Join(s.dir, name+".json")
	return os.Remove(filename)
}

// Render renders a template with given variables
func (s *Store) Render(name string, vars map[string]string) (string, string, error) {
	tmpl, err := s.Get(name)
	if err != nil {
		return "", "", err
	}

	// Apply defaults
	mergedVars := make(map[string]string)
	for k, v := range tmpl.Defaults {
		mergedVars[k] = v
	}
	for k, v := range vars {
		mergedVars[k] = v
	}

	// Render prompt
	prompt, err := renderString(tmpl.Prompt, mergedVars)
	if err != nil {
		return "", "", fmt.Errorf("failed to render prompt: %w", err)
	}

	// Render system message
	systemMsg, err := renderString(tmpl.SystemMsg, mergedVars)
	if err != nil {
		return "", "", fmt.Errorf("failed to render system message: %w", err)
	}

	return prompt, systemMsg, nil
}

// Export exports all templates to a file
func (s *Store) Export(filename string) error {
	templates, err := s.List()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(templates, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(filename, data, 0644)
}

// Import imports templates from a file
func (s *Store) Import(filename string) (int, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	var templates []*Template
	if err := json.Unmarshal(data, &templates); err != nil {
		return 0, fmt.Errorf("failed to parse import file: %w", err)
	}

	count := 0
	for _, tmpl := range templates {
		if err := s.Add(tmpl); err == nil {
			count++
		}
	}

	return count, nil
}

// extractVariables extracts {{.VarName}} or {{ .VarName }} variables from a string
func extractVariables(s string) []string {
	var vars []string
	seen := make(map[string]bool)

	for i := 0; i < len(s)-4; i++ {
		if s[i:i+2] == "{{" {
			end := strings.Index(s[i:], "}}")
			if end > 0 {
				content := s[i+2 : i+end]
				content = strings.TrimSpace(content)
				if len(content) > 0 && content[0] == '.' {
					varName := strings.TrimSpace(content[1:])
					if varName != "" && !seen[varName] {
						vars = append(vars, varName)
						seen[varName] = true
					}
				}
			}
		}
	}

	return vars
}

func renderString(s string, vars map[string]string) (string, error) {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	var buf strings.Builder
	if err := tmpl.Execute(&buf, vars); err != nil {
		return "", err
	}

	return buf.String(), nil
}

// GetBuiltinTemplates returns built-in templates
func GetBuiltinTemplates() []*Template {
	return []*Template{
		{
			Name:        "code-review-security",
			Description: "Security-focused code review",
			Command:     "review",
			SystemMsg:   "You are a security expert reviewing code for vulnerabilities.",
			Prompt:      "Review the following {{.Language}} code for security issues:\n\n{{.Code}}",
			Variables:   []string{"Language", "Code"},
			Defaults:    map[string]string{"Language": "code"},
		},
		{
			Name:        "translate-formal",
			Description: "Formal translation",
			Command:     "translate",
			SystemMsg:   "You are a professional translator. Use formal, business-appropriate language.",
			Prompt:      "Translate the following to {{.TargetLang}}:\n\n{{.Text}}",
			Variables:   []string{"TargetLang", "Text"},
			Defaults:    map[string]string{"TargetLang": "English"},
		},
		{
			Name:        "explain-beginner",
			Description: "Explain for beginners",
			Command:     "explain",
			SystemMsg:   "You are a patient teacher explaining concepts to beginners. Use simple language and examples.",
			Prompt:      "Explain this error in simple terms:\n\n{{.Error}}",
			Variables:   []string{"Error"},
		},
		{
			Name:        "commit-conventional",
			Description: "Conventional commit message",
			Command:     "commit",
			SystemMsg:   "Generate commit messages following Conventional Commits format strictly.",
			Prompt:      "Generate a {{.Type}} commit message for:\n\n{{.Changes}}",
			Variables:   []string{"Type", "Changes"},
			Defaults:    map[string]string{"Type": "feat"},
		},
		{
			Name:        "summarize-bullet",
			Description: "Bullet point summary",
			Command:     "summarize",
			SystemMsg:   "Summarize content as concise bullet points. Maximum 5-7 points.",
			Prompt:      "Summarize in bullet points:\n\n{{.Content}}",
			Variables:   []string{"Content"},
		},
	}
}
