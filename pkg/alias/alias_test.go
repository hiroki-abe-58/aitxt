package alias

import (
	"os"
	"path/filepath"
	"testing"
)

func setupTestStore(t *testing.T) (*Store, func()) {
	tmpDir := t.TempDir()
	store := &Store{
		file:    filepath.Join(tmpDir, "aliases.json"),
		aliases: make(map[string]*Alias),
	}
	
	cleanup := func() {
		os.RemoveAll(tmpDir)
	}
	
	return store, cleanup
}

func TestAdd(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.Add("s", "summarize", "Shortcut for summarize")
	if err != nil {
		t.Fatalf("Add() failed: %v", err)
	}

	alias, ok := store.Get("s")
	if !ok {
		t.Fatal("Get() failed to retrieve added alias")
	}

	if alias.Name != "s" {
		t.Errorf("Expected Name='s', got '%s'", alias.Name)
	}
	if alias.Command != "summarize" {
		t.Errorf("Expected Command='summarize', got '%s'", alias.Command)
	}
	if alias.Desc != "Shortcut for summarize" {
		t.Errorf("Expected Desc='Shortcut for summarize', got '%s'", alias.Desc)
	}
}

func TestAddEmptyName(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.Add("", "summarize", "")
	if err == nil {
		t.Error("Add() should fail with empty name")
	}
}

func TestAddEmptyCommand(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.Add("s", "", "")
	if err == nil {
		t.Error("Add() should fail with empty command")
	}
}

func TestAddReservedName(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	reserved := []string{"ask", "summarize", "translate", "config", "version"}
	for _, name := range reserved {
		err := store.Add(name, "test", "")
		if err == nil {
			t.Errorf("Add() should fail for reserved name '%s'", name)
		}
	}
}

func TestGet(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.Add("t", "translate", "")

	alias, ok := store.Get("t")
	if !ok {
		t.Fatal("Get() should find existing alias")
	}
	if alias.Command != "translate" {
		t.Errorf("Expected Command='translate', got '%s'", alias.Command)
	}

	_, ok = store.Get("nonexistent")
	if ok {
		t.Error("Get() should return false for non-existent alias")
	}
}

func TestDelete(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.Add("test", "summarize", "")

	err := store.Delete("test")
	if err != nil {
		t.Fatalf("Delete() failed: %v", err)
	}

	_, ok := store.Get("test")
	if ok {
		t.Error("Get() should not find deleted alias")
	}
}

func TestDeleteNonExistent(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	err := store.Delete("nonexistent")
	if err == nil {
		t.Error("Delete() should fail for non-existent alias")
	}
}

func TestList(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.Add("s", "summarize", "")
	store.Add("t", "translate", "")
	store.Add("a", "ask", "")

	aliases := store.List()
	if len(aliases) != 3 {
		t.Errorf("Expected 3 aliases, got %d", len(aliases))
	}

	// Check if sorted by name
	if aliases[0].Name != "a" || aliases[1].Name != "s" || aliases[2].Name != "t" {
		t.Error("List() should return aliases sorted by name")
	}
}

func TestResolve(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	store.Add("s", "summarize", "")
	store.Add("tj", "translate --to ja", "")

	tests := []struct {
		name  string
		input []string
		want  []string
	}{
		{
			name:  "Simple alias",
			input: []string{"s", "file.txt"},
			want:  []string{"summarize", "file.txt"},
		},
		{
			name:  "Alias with flags",
			input: []string{"tj", "hello"},
			want:  []string{"translate", "--to", "ja", "hello"},
		},
		{
			name:  "No alias",
			input: []string{"ask", "question"},
			want:  []string{"ask", "question"},
		},
		{
			name:  "Empty args",
			input: []string{},
			want:  []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := store.Resolve(tt.input)
			if len(got) != len(tt.want) {
				t.Errorf("Resolve() returned %d args, want %d", len(got), len(tt.want))
				return
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Errorf("Resolve()[%d] = %s, want %s", i, got[i], tt.want[i])
				}
			}
		})
	}
}

func TestSaveAndLoad(t *testing.T) {
	store, cleanup := setupTestStore(t)
	defer cleanup()

	// Add aliases
	store.Add("s", "summarize", "Summary")
	store.Add("t", "translate", "Translate")

	// Create new store with same file
	store2 := &Store{
		file:    store.file,
		aliases: make(map[string]*Alias),
	}
	store2.load()

	// Verify aliases were loaded
	if len(store2.aliases) != 2 {
		t.Errorf("Expected 2 loaded aliases, got %d", len(store2.aliases))
	}

	alias, ok := store2.Get("s")
	if !ok || alias.Command != "summarize" {
		t.Error("Failed to load saved alias")
	}
}

func TestGetBuiltinAliases(t *testing.T) {
	builtins := GetBuiltinAliases()

	if len(builtins) == 0 {
		t.Error("GetBuiltinAliases() should return non-empty list")
	}

	// Check a few expected aliases
	foundS := false
	foundTj := false
	for _, a := range builtins {
		if a.Name == "s" && a.Command == "summarize" {
			foundS = true
		}
		if a.Name == "tj" && a.Command == "translate --to ja" {
			foundTj = true
		}
	}

	if !foundS {
		t.Error("Expected to find 's' alias in builtins")
	}
	if !foundTj {
		t.Error("Expected to find 'tj' alias in builtins")
	}
}
