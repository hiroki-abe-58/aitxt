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
	AppDescription    string
	AppLongDesc       string
	Version           string
	
	// Errors
	ErrNoInput        string
	ErrReadFile       string
	ErrReadStdin      string
	ErrLoadConfig     string
	ErrCreateClient   string
	ErrGenerate       string
	ErrNoStagedChanges string
	
	// Commands
	Summarizing       string
	Translating       string
	Proofreading      string
	Generating        string
	Analyzing         string
	TokensUsed        string
	
	// Commit command
	SuggestedCommit   string
	CommitInstruction string
}

var currentLang Language = LangEnglish
var messages = make(map[Language]*Messages)

func init() {
	// Register Japanese messages
	registerJapanese()
}

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
	// Fallback to English (will be added later)
	return messages[LangJapanese]
}

// Register registers messages for a language
func Register(lang Language, msg *Messages) {
	messages[lang] = msg
}

// registerJapanese registers Japanese messages
func registerJapanese() {
	Register(LangJapanese, &Messages{
		// General
		AppDescription:    "🤖 AI搭載テキスト処理CLIツール",
		AppLongDesc:       "aitxtは複数のLLMプロバイダーを使用した強力なテキスト処理機能を提供します",
		Version:           "バージョン",
		
		// Errors
		ErrNoInput:        "入力テキストがありません",
		ErrReadFile:       "ファイルの読み込みに失敗しました",
		ErrReadStdin:      "標準入力の読み込みに失敗しました",
		ErrLoadConfig:     "設定の読み込みに失敗しました",
		ErrCreateClient:   "クライアントの作成に失敗しました",
		ErrGenerate:       "生成に失敗しました",
		ErrNoStagedChanges: "ステージされた変更がありません。'git add'で変更をステージしてください",
		
		// Commands
		Summarizing:       "%sで要約中...",
		Translating:       "%sで%sに翻訳中...",
		Proofreading:      "%sで校正中（%sスタイル）...",
		Generating:        "%sでコミットメッセージを生成中...",
		Analyzing:         "%sでエラーを分析中...",
		TokensUsed:        "[使用トークン数: %d]",
		
		// Commit command
		SuggestedCommit:   "--- 提案されたコミットメッセージ ---",
		CommitInstruction: "このメッセージを使用するには:\n  git commit -m \"<上記のメッセージをコピー>\"",
	})
}
