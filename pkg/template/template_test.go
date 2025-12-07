package template

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func setupTestStore(t *testing.T) (*Store, func()) {
	tmpDir := t.TempDir()
	store := &Store{dir: tmpDir}
	
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	
	return store, cleanup
}

func TestAdd(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tmpl := &Template{
		Name:        "test-template",
		Description: "Test template",
		Command:     "summarize",
		SystemMsg:   "You are a summarizer",
		Prompt:      "Summarize: {{.Content}}",
	}

	err := store.Add(tmpl)
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	// Verify file was created
	filename := filepath.Join(store.dir, "test-template.json")
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		t.Error("Add() should create template file")
	}

	// Verify variables were extracted
	if len(tmpl.Variables) != 1 || tmpl.Variables[0] != "Content" {
		t.Errorf("Expected Variables=['Content'], got %v", tmpl.Variables)
	}
}

func TestAddEmptyName(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tmpl := &Template{
		Name:   "",
		Prompt: "Test",
	}

	err := store.Add(tmpl)
	if err == nil {
		t.Error("Add() should fail with empty name")
	}
}

func TestGet(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	original := &Template{
		Name:        "my-template",
		Description: "My template",
		Prompt:      "Hello {{.Name}}",
	}
	store.Add(original)

	retrieved, err := store.Get("my-template")
	if err != nil {
		t.Fatalf("Get() failed: %v", err)
	}

	if retrieved.Name != original.Name {
		t.Errorf("Expected Name='%s', got '%s'", original.Name, retrieved.Name)
	}
	if retrieved.Description != original.Description {
		t.Errorf("Expected Description='%s', got '%s'", original.Description, retrieved.Description)
	}
}

func TestGetNonExistent(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	_, err := store.Get("nonexistent")
	if err == nil {
		t.Error("Get() should fail for non-existent template")
	}
}

func TestList(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	templates := []*Template{
		{Name: "c-template", Prompt: "C"},
		{Name: "a-template", Prompt: "A"},
		{Name: "b-template", Prompt: "B"},
	}

	for _, tmpl := range templates {
		store.Add(tmpl)
	}

	list, err := store.List()
	if err != nil {
		t.Fatalf("List() failed: %v", err)
	}

	if len(list) != 3 {
		t.Errorf("Expected 3 templates, got %d", len(list))
	}

	// Check if sorted by name
	if list[0].Name != "a-template" || list[1].Name != "b-template" || list[2].Name != "c-template" {
		t.Error("List() should return templates sorted by name")
	}
}

func TestDelete(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tmpl := &Template{Name: "delete-me", Prompt: "Test"}
	store.Add(tmpl)

	err := store.Delete("delete-me")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	_, err = store.Get("delete-me")
	if err == nil {
		t.Error("Get() should fail after Delete()")
	}
}

func TestRender(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tmpl := &Template{
		Name:      "greeting",
		SystemMsg: "You are helpful",
		Prompt:    "Say hello to {{.Name}} from {{.Location}}",
	}
	store.Add(tmpl)

	vars := map[string]string{
		"Name":     "Alice",
		"Location": "Tokyo",
	}

	prompt, systemMsg, err := store.Render("greeting", vars)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	expectedPrompt := "Say hello to Alice from Tokyo"
	if prompt != expectedPrompt {
		t.Errorf("Expected prompt='%s', got '%s'", expectedPrompt, prompt)
	}

	if systemMsg != "You are helpful" {
		t.Errorf("Expected systemMsg='You are helpful', got '%s'", systemMsg)
	}
}

func TestRenderWithDefaults(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	tmpl := &Template{
		Name:   "with-defaults",
		Prompt: "Hello {{.Name}}, you are {{.Age}} years old",
		Defaults: map[string]string{
			"Age": "30",
		},
	}
	store.Add(tmpl)

	vars := map[string]string{
		"Name": "Bob",
	}

	prompt, _, err := store.Render("with-defaults", vars)
	if err != nil {
		t.Fatalf("Render() failed: %v", err)
	}

	expected := "Hello Bob, you are 30 years old"
	if prompt != expected {
		t.Errorf("Expected prompt='%s', got '%s'", expected, prompt)
	}
}

func TestExtractVariables(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{
			name:  "Single variable",
			input: "Hello {{.Name}}",
			want:  []string{"Name"},
		},
		{
			name:  "Multiple variables",
			input: "{{.First}} {{.Second}} {{.Third}}",
			want:  []string{"First", "Second", "Third"},
		},
		{
			name:  "Duplicate variables",
			input: "{{.Name}} and {{.Name}} again",
			want:  []string{"Name"},
		},
		{
			name:  "No variables",
			input: "No variables here",
			want:  []string{},
		},
		{
			name:  "Variables with spaces",
			input: "{{ .Name }} and {{.Age}}",
			want:  []string{"Name", "Age"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractVariables(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("extractVariables() returned %d variables, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("extractVariables()[%d] = %s, want %s", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestGetBuiltinTemplates(t *testing.T) {
	builtins := GetBuiltinTemplates()

	if len(builtins) == 0 {
		t.Error("GetBuiltinTemplates() should return non-empty list")
	}

	// Check for expected builtins
	foundReview := false
	foundTranslate := false
	for _, tmpl := range builtins {
		if tmpl.Name == "code-review-security" {
			foundReview = true
			if tmpl.Command != "review" {
				t.Errorf("Expected code-review-security command='review', got '%s'", tmpl.Command)
			}
		}
		if tmpl.Name == "translate-formal" {
			foundTranslate = true
		}
	}

	if !foundReview {
		t.Error("Expected to find 'code-review-security' in builtins")
	}
	if !foundTranslate {
		t.Error("Expected to find 'translate-formal' in builtins")
	}
}

func TestTemplateTimestamps(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	before := time.Now()
	tmpl := &Template{Name: "timestamp-test", Prompt: "Test"}
	store.Add(tmpl)
	after := time.Now()

	if tmpl.CreatedAt.Before(before) || tmpl.CreatedAt.After(after) {
		t.Error("CreatedAt should be set during Add()")
	}
	if tmpl.UpdatedAt.Before(before) || tmpl.UpdatedAt.After(after) {
		t.Error("UpdatedAt should be set during Add()")
	}
}

func TestExportImport(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Add templates
	store.Add(&Template{Name: "tmpl1", Prompt: "Test 1"})
	store.Add(&Template{Name: "tmpl2", Prompt: "Test 2"})

	// Export
	exportFile := filepath.Join(store.dir, "export.json")
	err := store.Export(exportFile)
	if err != nil {
		t.Fatalf("Export() failed: %v", err)
	}

	// Create new store and import
	store2, cleanup2 := setupTestStore(t)
	defer cleanup2()

	count, err := store2.Import(exportFile)
	if err != nil {
		t.Fatalf("Import() failed: %v", err)
	}

	if count != 2 {
		t.Errorf("Expected to import 2 templates, got %d", count)
	}

	// Verify imported templates
	_, err = store2.Get("tmpl1")
	if err != nil {
		t.Error("Failed to get imported template 'tmpl1'")
	}
}
