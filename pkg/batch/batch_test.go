package batch

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/llm"
)

// Mock LLM client for testing
type mockClient struct {
	generateFunc func(ctx context.Context, req *llm.Request) (*llm.Response, error)
}

func (m *mockClient) Generate(ctx context.Context, req *llm.Request) (*llm.Response, error) {
	if m.generateFunc != nil {
		return m.generateFunc(ctx, req)
	}
	return &llm.Response{
		Text:       "Mock response",
		Provider:   llm.ProviderOpenAI,
		Model:      "mock-model",
		TokensUsed: 10,
	}, nil
}

func (m *mockClient) Stream(ctx context.Context, req *llm.Request, callback func(chunk string) error) error {
	return callback("Mock chunk")
}

func (m *mockClient) GetProvider() llm.Provider {
	return llm.ProviderOpenAI
}

func (m *mockClient) GetModel() string {
	return "mock-model"
}

func (m *mockClient) Validate() error {
	return nil
}

func TestNewProcessor(t *testing.T) {
	client := &mockClient{}

	tests := []struct {
		name        string
		concurrency int
		timeout     time.Duration
		wantConc    int
		wantTimeout time.Duration
	}{
		{"Default concurrency", 0, 60 * time.Second, 3, 60 * time.Second},
		{"Default timeout", 5, 0, 5, 120 * time.Second},
		{"Custom values", 10, 30 * time.Second, 10, 30 * time.Second},
		{"Negative concurrency", -1, 0, 3, 120 * time.Second},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := NewProcessor(client, tt.concurrency, tt.timeout)
			if p.concurrency != tt.wantConc {
				t.Errorf("Expected concurrency=%d, got %d", tt.wantConc, p.concurrency)
			}
			if p.timeout != tt.wantTimeout {
				t.Errorf("Expected timeout=%v, got %v", tt.wantTimeout, p.timeout)
			}
		})
	}
}

func TestProcessFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	file1 := filepath.Join(tmpDir, "test1.txt")
	file2 := filepath.Join(tmpDir, "test2.txt")
	os.WriteFile(file1, []byte("Content 1"), 0644)
	os.WriteFile(file2, []byte("Content 2"), 0644)

	client := &mockClient{}
	processor := NewProcessor(client, 2, 10*time.Second)

	files := []string{file1, file2}
	result := processor.ProcessFiles(context.Background(), files, "Process: %s", "System message")

	if result.TotalTasks != 2 {
		t.Errorf("Expected TotalTasks=2, got %d", result.TotalTasks)
	}

	if result.SuccessCount != 2 {
		t.Errorf("Expected SuccessCount=2, got %d", result.SuccessCount)
	}

	if result.FailureCount != 0 {
		t.Errorf("Expected FailureCount=0, got %d", result.FailureCount)
	}

	if len(result.Tasks) != 2 {
		t.Errorf("Expected 2 tasks, got %d", len(result.Tasks))
	}

	for _, task := range result.Tasks {
		if task.Error != nil {
			t.Errorf("Task %d should not have error: %v", task.ID, task.Error)
		}
		if task.Result == "" {
			t.Error("Task should have result")
		}
	}
}

func TestProcessFilesWithError(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(file1, []byte("Content"), 0644)

	client := &mockClient{
		generateFunc: func(ctx context.Context, req *llm.Request) (*llm.Response, error) {
			return nil, os.ErrNotExist
		},
	}
	processor := NewProcessor(client, 1, 10*time.Second)

	files := []string{file1}
	result := processor.ProcessFiles(context.Background(), files, "Process: %s", "System")

	if result.FailureCount != 2 {
		t.Errorf("Expected FailureCount=1, got %d", result.FailureCount)
	}

	if result.Tasks[0].Error == nil {
		t.Error("Task should have error")
	}
}

func TestProcessFilesMissingFile(t *testing.T) {
	client := &mockClient{}
	processor := NewProcessor(client, 1, 10*time.Second)

	files := []string{"/nonexistent/file.txt"}
	result := processor.ProcessFiles(context.Background(), files, "Process: %s", "System")

	if result.FailureCount != 2 {
		t.Errorf("Expected FailureCount=1, got %d", result.FailureCount)
	}

	if result.Tasks[0].Error == nil {
		t.Error("Task should have error for missing file")
	}
}

func TestFindFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	os.WriteFile(filepath.Join(tmpDir, "test1.txt"), []byte("1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "test2.txt"), []byte("2"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "test.md"), []byte("3"), 0644)
	os.Mkdir(filepath.Join(tmpDir, "subdir"), 0755)

	pattern := filepath.Join(tmpDir, "*.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("FindFiles() failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 .txt files, got %d", len(files))
	}
}

func TestFindFilesInvalidPattern(t *testing.T) {
	_, err := FindFiles("[invalid")
	if err == nil {
		t.Error("FindFiles() should fail with invalid pattern")
	}
}

func TestFindFilesRecursive(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test structure
	os.WriteFile(filepath.Join(tmpDir, "root.go"), []byte("1"), 0644)
	os.WriteFile(filepath.Join(tmpDir, "root.txt"), []byte("2"), 0644)
	
	subdir := filepath.Join(tmpDir, "sub")
	os.Mkdir(subdir, 0755)
	os.WriteFile(filepath.Join(subdir, "sub.go"), []byte("3"), 0644)
	os.WriteFile(filepath.Join(subdir, "sub.py"), []byte("4"), 0644)

	// Find all Go files
	files, err := FindFilesRecursive(tmpDir, []string{".go"})
	if err != nil {
		t.Fatalf("FindFilesRecursive() failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("Expected 2 .go files, got %d", len(files))
	}

	// Find all files
	files, err = FindFilesRecursive(tmpDir, []string{})
	if err != nil {
		t.Fatalf("FindFilesRecursive() failed: %v", err)
	}

	if len(files) != 4 {
		t.Errorf("Expected 4 total files, got %d", len(files))
	}
}

func TestFindFilesRecursiveNonExistent(t *testing.T) {
	_, err := FindFilesRecursive("/nonexistent/path", []string{".go"})
	if err == nil {
		t.Error("FindFilesRecursive() should fail for non-existent directory")
	}
}

func TestResultStats(t *testing.T) {
	tmpDir := t.TempDir()
	file1 := filepath.Join(tmpDir, "test1.txt")
	file2 := filepath.Join(tmpDir, "test2.txt")
	os.WriteFile(file1, []byte("Content 1"), 0644)
	os.WriteFile(file2, []byte("Content 2"), 0644)

	client := &mockClient{}
	processor := NewProcessor(client, 2, 10*time.Second)

	result := processor.ProcessFiles(context.Background(), []string{file1, file2}, "Process: %s", "System")

	if result.TotalDuration == 0 {
		t.Error("TotalDuration should be non-zero")
	}

	for _, task := range result.Tasks {
		if task.Duration == 0 {
			t.Error("Task duration should be non-zero")
		}
	}
}

func TestConcurrency(t *testing.T) {
	tmpDir := t.TempDir()

	// Create multiple test files
	var files []string
	for i := 0; i < 10; i++ {
		file := filepath.Join(tmpDir, "test"+string(rune('0'+i))+".txt")
		os.WriteFile(file, []byte("Content"), 0644)
		files = append(files, file)
	}

	// Track concurrent executions
	executing := 0
	maxConcurrent := 0

	client := &mockClient{
		generateFunc: func(ctx context.Context, req *llm.Request) (*llm.Response, error) {
			executing++
			if executing > maxConcurrent {
				maxConcurrent = executing
			}
			time.Sleep(10 * time.Millisecond)
			executing--
			return &llm.Response{Text: "Response", Provider: llm.ProviderOpenAI, Model: "mock", TokensUsed: 10}, nil
		},
	}

	processor := NewProcessor(client, 3, 10*time.Second)
	processor.ProcessFiles(context.Background(), files, "Process: %s", "System")

	// Max concurrent should not exceed concurrency limit
	if maxConcurrent > 3 {
		t.Errorf("Max concurrent executions was %d, should not exceed 3", maxConcurrent)
	}
}
