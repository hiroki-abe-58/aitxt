package history

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Entry struct {
	ID        string    `json:"id"`
	Timestamp time.Time `json:"timestamp"`
	Provider  string    `json:"provider"`
	Command   string    `json:"command"`
	Input     string    `json:"input"`
	Output    string    `json:"output"`
	Tokens    int       `json:"tokens,omitempty"`
}

type Store struct {
	dir string
}

func NewStore() (*Store, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}
	dir := filepath.Join(homeDir, ".aitxt", "history")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create history directory: %w", err)
	}
	return &Store{dir: dir}, nil
}

func (s *Store) Add(entry *Entry) error {
	if entry.ID == "" {
		entry.ID = fmt.Sprintf("%d", time.Now().UnixNano())
	}
	if entry.Timestamp.IsZero() {
		entry.Timestamp = time.Now()
	}
	filename := filepath.Join(s.dir, entry.ID+".json")
	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (s *Store) List(limit int) ([]*Entry, error) {
	files, err := os.ReadDir(s.dir)
	if err != nil {
		return nil, err
	}
	var entries []*Entry
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}
		entry, err := s.loadEntry(filepath.Join(s.dir, file.Name()))
		if err != nil {
			continue
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Timestamp.After(entries[j].Timestamp)
	})
	if limit > 0 && len(entries) > limit {
		entries = entries[:limit]
	}
	return entries, nil
}

func (s *Store) Get(id string) (*Entry, error) {
	return s.loadEntry(filepath.Join(s.dir, id+".json"))
}

func (s *Store) Search(keyword string, limit int) ([]*Entry, error) {
	entries, err := s.List(0)
	if err != nil {
		return nil, err
	}
	keyword = strings.ToLower(keyword)
	var results []*Entry
	for _, entry := range entries {
		if strings.Contains(strings.ToLower(entry.Input), keyword) ||
			strings.Contains(strings.ToLower(entry.Output), keyword) {
			results = append(results, entry)
			if limit > 0 && len(results) >= limit {
				break
			}
		}
	}
	return results, nil
}

func (s *Store) Delete(id string) error {
	return os.Remove(filepath.Join(s.dir, id+".json"))
}

func (s *Store) Clear() error {
	files, _ := os.ReadDir(s.dir)
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			os.Remove(filepath.Join(s.dir, file.Name()))
		}
	}
	return nil
}

func (s *Store) Export(filename string, format string) error {
	entries, err := s.List(0)
	if err != nil {
		return err
	}
	data, _ := json.MarshalIndent(entries, "", "  ")
	return os.WriteFile(filename, data, 0644)
}

func (s *Store) loadEntry(filename string) (*Entry, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var entry Entry
	json.Unmarshal(data, &entry)
	return &entry, nil
}

func (s *Store) GetStats() (map[string]interface{}, error) {
	entries, _ := s.List(0)
	stats := map[string]interface{}{"total": len(entries), "by_command": map[string]int{}, "by_provider": map[string]int{}, "total_tokens": 0}
	return stats, nil
}
