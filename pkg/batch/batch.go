package batch

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/llm"
)

// Task represents a single batch task
type Task struct {
	ID       int
	FilePath string
	Content  string
	Result   string
	Error    error
	Duration time.Duration
}

// Result represents batch processing result
type Result struct {
	TotalTasks     int
	SuccessCount   int
	FailureCount   int
	TotalDuration  time.Duration
	Tasks          []*Task
}

// Processor handles batch processing
type Processor struct {
	client      llm.Client
	concurrency int
	timeout     time.Duration
}

// NewProcessor creates a new batch processor
func NewProcessor(client llm.Client, concurrency int, timeout time.Duration) *Processor {
	if concurrency <= 0 {
		concurrency = 3
	}
	if timeout <= 0 {
		timeout = 120 * time.Second
	}
	return &Processor{
		client:      client,
		concurrency: concurrency,
		timeout:     timeout,
	}
}

// ProcessFiles processes multiple files with the given prompt template
func (p *Processor) ProcessFiles(ctx context.Context, files []string, promptTemplate string, systemMsg string) *Result {
	startTime := time.Now()
	result := &Result{
		TotalTasks: len(files),
		Tasks:      make([]*Task, len(files)),
	}

	// Create tasks
	tasks := make([]*Task, len(files))
	for i, file := range files {
		content, err := os.ReadFile(file)
		if err != nil {
			tasks[i] = &Task{
				ID:       i,
				FilePath: file,
				Error:    fmt.Errorf("failed to read file: %w", err),
			}
			continue
		}
		tasks[i] = &Task{
			ID:       i,
			FilePath: file,
			Content:  string(content),
		}
	}

	// Process with concurrency control
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, p.concurrency)

	for i, task := range tasks {
		if task.Error != nil {
			result.Tasks[i] = task
			result.FailureCount++
			continue
		}

		wg.Add(1)
		go func(idx int, t *Task) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			taskStart := time.Now()
			prompt := fmt.Sprintf(promptTemplate, t.Content)

			req := &llm.Request{
				Prompt:    prompt,
				SystemMsg: systemMsg,
				MaxTokens: 2000,
			}

			taskCtx, cancel := context.WithTimeout(ctx, p.timeout)
			defer cancel()

			resp, err := p.client.Generate(taskCtx, req)
			t.Duration = time.Since(taskStart)

			if err != nil {
				t.Error = err
			} else {
				t.Result = resp.Text
			}

			result.Tasks[idx] = t
		}(i, task)
	}

	wg.Wait()

	// Calculate results
	for _, task := range result.Tasks {
		if task.Error != nil {
			result.FailureCount++
		} else {
			result.SuccessCount++
		}
	}
	result.TotalDuration = time.Since(startTime)

	return result
}

// FindFiles finds files matching the given pattern
func FindFiles(pattern string) ([]string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %w", err)
	}

	var files []string
	for _, match := range matches {
		info, err := os.Stat(match)
		if err != nil {
			continue
		}
		if !info.IsDir() {
			files = append(files, match)
		}
	}

	return files, nil
}

// FindFilesRecursive finds files matching extensions recursively
func FindFilesRecursive(dir string, extensions []string) ([]string, error) {
	var files []string
	extMap := make(map[string]bool)
	for _, ext := range extensions {
		extMap[ext] = true
	}

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		ext := filepath.Ext(path)
		if len(extensions) == 0 || extMap[ext] {
			files = append(files, path)
		}
		return nil
	})

	return files, err
}
