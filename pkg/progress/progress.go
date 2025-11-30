package progress

import (
	"fmt"
	"time"

	"github.com/briandowns/spinner"
)

// Spinner wraps the spinner library for consistent usage
type Spinner struct {
	spinner *spinner.Spinner
	message string
}

// SpinnerStyle defines available spinner styles
type SpinnerStyle int

const (
	StyleDots SpinnerStyle = iota
	StyleLine
	StyleCircle
	StyleArrow
	StyleBounce
)

var spinnerCharSets = map[SpinnerStyle][]string{
	StyleDots:   {"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
	StyleLine:   {"-", "\\", "|", "/"},
	StyleCircle: {"◐", "◓", "◑", "◒"},
	StyleArrow:  {"←", "↖", "↑", "↗", "→", "↘", "↓", "↙"},
	StyleBounce: {"⠁", "⠂", "⠄", "⠂"},
}

// NewSpinner creates a new spinner with a message
func NewSpinner(message string) *Spinner {
	s := spinner.New(spinnerCharSets[StyleDots], 100*time.Millisecond)
	s.Suffix = " " + message
	return &Spinner{
		spinner: s,
		message: message,
	}
}

// NewSpinnerWithStyle creates a spinner with a specific style
func NewSpinnerWithStyle(message string, style SpinnerStyle) *Spinner {
	charSet, ok := spinnerCharSets[style]
	if !ok {
		charSet = spinnerCharSets[StyleDots]
	}
	s := spinner.New(charSet, 100*time.Millisecond)
	s.Suffix = " " + message
	return &Spinner{
		spinner: s,
		message: message,
	}
}

// Start starts the spinner
func (s *Spinner) Start() {
	s.spinner.Start()
}

// Stop stops the spinner
func (s *Spinner) Stop() {
	s.spinner.Stop()
}

// Success stops with a success message
func (s *Spinner) Success(message string) {
	s.spinner.Stop()
	fmt.Printf("✅ %s\n", message)
}

// Error stops with an error message
func (s *Spinner) Error(message string) {
	s.spinner.Stop()
	fmt.Printf("❌ %s\n", message)
}

// UpdateMessage updates the spinner message
func (s *Spinner) UpdateMessage(message string) {
	s.spinner.Suffix = " " + message
}

// Processing shows a processing spinner for a function
func Processing(message string, fn func() error) error {
	s := NewSpinner(message)
	s.Start()
	err := fn()
	if err != nil {
		s.Error(err.Error())
		return err
	}
	s.Stop()
	return nil
}

// ProcessingWithResult shows a spinner and returns a result
func ProcessingWithResult[T any](message string, fn func() (T, error)) (T, error) {
	s := NewSpinner(message)
	s.Start()
	result, err := fn()
	if err != nil {
		s.Error(err.Error())
		return result, err
	}
	s.Stop()
	return result, nil
}
