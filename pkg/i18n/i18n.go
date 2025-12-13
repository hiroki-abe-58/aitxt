package i18n

import (
	"os"
	"strings"
)

// Language represents a supported language
type Language string

const (
	LangEnglish  Language = "en"
	LangJapanese Language = "ja"
	LangChinese  Language = "zh"
	LangKorean   Language = "ko"
	LangThai     Language = "th"
)

// Messages holds all translatable messages
type Messages struct {
	// General
	AppDescription string
	AppLongDesc    string
	Version        string

	// Errors
	ErrNoInput         string
	ErrReadFile        string
	ErrReadStdin       string
	ErrLoadConfig      string
	ErrCreateClient    string
	ErrGenerate        string
	ErrNoStagedChanges string

	// Commands
	Summarizing  string
	Translating  string
	Proofreading string
	Generating   string
	Analyzing    string
	TokensUsed   string

	// Commit command
	SuggestedCommit   string
	CommitInstruction string
}

var currentLang Language = LangEnglish
var messages = make(map[Language]*Messages)

// SetLanguage sets the current language
func SetLanguage(lang Language) {
	if _, ok := messages[lang]; ok {
		currentLang = lang
	}
}

// DetectLanguage detects the language from environment
func DetectLanguage() Language {
	// Check LANG environment variable
	lang := os.Getenv("LANG")
	lang = strings.ToLower(lang)

	switch {
	case strings.HasPrefix(lang, "ja"):
		return LangJapanese
	case strings.HasPrefix(lang, "zh"):
		return LangChinese
	case strings.HasPrefix(lang, "ko"):
		return LangKorean
	case strings.HasPrefix(lang, "th"):
		return LangThai
	default:
		return LangEnglish
	}
}

// Get returns the current messages
func Get() *Messages {
	if msg, ok := messages[currentLang]; ok {
		return msg
	}
	// Fallback to English
	return messages[LangEnglish]
}

// Register registers messages for a language
func Register(lang Language, msg *Messages) {
	messages[lang] = msg
}
