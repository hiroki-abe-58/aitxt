package alias

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Alias represents a command alias
type Alias struct {
	Name    string `json:"name"`
	Command string `json:"command"`
	Desc    string `json:"description,omitempty"`
}

// Store manages aliases
type Store struct {
	file    string
	aliases map[string]*Alias
}

// NewStore creates a new alias store
func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	dir := filepath.Join(homeDir, ".aitxt")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	store := &Store{
		file:    filepath.Join(dir, "aliases.json"),
		aliases: make(map[string]*Alias),
	}

	store.load()
	return store, nil
}

// load loads aliases from file
func (s *Store) load() error {
	data, err := os.ReadFile(s.file)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	var aliases []*Alias
	if err := json.Unmarshal(data, &aliases); err != nil {
		return err
	}

	for _, a := range aliases {
		s.aliases[a.Name] = a
	}

	return nil
}

// save saves aliases to file
func (s *Store) save() error {
	aliases := s.List()
	data, err := json.MarshalIndent(aliases, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.file, data, 0644)
}

// Add adds a new alias
func (s *Store) Add(name, command, desc string) error {
	if name == "" {
		return fmt.Errorf("alias name is required")
	}
	if command == "" {
		return fmt.Errorf("command is required")
	}

	// Check for reserved names
	reserved := []string{"ask", "batch", "chat", "commit", "completion", "config",
		"doc", "explain", "history", "proofread", "review", "style",
		"summarize", "template", "translate", "version", "alias", "help"}
	for _, r := range reserved {
		if name == r {
			return fmt.Errorf("'%s' is a reserved command name", name)
		}
	}

	s.aliases[name] = &Alias{
		Name:    name,
		Command: command,
		Desc:    desc,
	}

	return s.save()
}

// Get retrieves an alias by name
func (s *Store) Get(name string) (*Alias, bool) {
	alias, ok := s.aliases[name]
	return alias, ok
}

// Delete removes an alias
func (s *Store) Delete(name string) error {
	if _, ok := s.aliases[name]; !ok {
		return fmt.Errorf("alias not found: %s", name)
	}

	delete(s.aliases, name)
	return s.save()
}

// List returns all aliases sorted by name
func (s *Store) List() []*Alias {
	var aliases []*Alias
	for _, a := range s.aliases {
		aliases = append(aliases, a)
	}

	sort.Slice(aliases, func(i, j int) bool {
		return aliases[i].Name < aliases[j].Name
	})

	return aliases
}

// Resolve expands an alias to its full command
func (s *Store) Resolve(args []string) []string {
	if len(args) == 0 {
		return args
	}

	alias, ok := s.aliases[args[0]]
	if !ok {
		return args
	}

	// Parse the alias command into parts
	parts := strings.Fields(alias.Command)
	// Append remaining original args
	parts = append(parts, args[1:]...)

	return parts
}

// GetBuiltinAliases returns suggested built-in aliases
func GetBuiltinAliases() []*Alias {
	return []*Alias{
		{Name: "s", Command: "summarize", Desc: "Shortcut for summarize"},
		{Name: "t", Command: "translate", Desc: "Shortcut for translate"},
		{Name: "tj", Command: "translate --to ja", Desc: "Translate to Japanese"},
		{Name: "tc", Command: "translate --to zh", Desc: "Translate to Chinese"},
		{Name: "p", Command: "proofread", Desc: "Shortcut for proofread"},
		{Name: "r", Command: "review", Desc: "Shortcut for review"},
		{Name: "rs", Command: "review --focus security", Desc: "Security review"},
		{Name: "rp", Command: "review --focus performance", Desc: "Performance review"},
		{Name: "e", Command: "explain", Desc: "Shortcut for explain"},
		{Name: "ej", Command: "explain --lang ja", Desc: "Explain in Japanese"},
		{Name: "c", Command: "commit", Desc: "Shortcut for commit"},
		{Name: "cj", Command: "commit --lang ja", Desc: "Commit message in Japanese"},
		{Name: "d", Command: "doc", Desc: "Shortcut for doc"},
		{Name: "a", Command: "ask", Desc: "Shortcut for ask"},
		{Name: "aj", Command: "ask --lang ja", Desc: "Ask in Japanese"},
	}
}
