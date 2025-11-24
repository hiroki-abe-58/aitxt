package i18n

func init() {
	registerEnglish()
}

// registerEnglish registers English messages
func registerEnglish() {
	Register(LangEnglish, &Messages{
		// General
		AppDescription:    "🤖 AI-powered text processing CLI tool",
		AppLongDesc:       "aitxt provides powerful text processing capabilities using multiple LLM providers",
		Version:           "Version",
		
		// Errors
		ErrNoInput:        "No input text provided",
		ErrReadFile:       "Failed to read file",
		ErrReadStdin:      "Failed to read stdin",
		ErrLoadConfig:     "Failed to load config",
		ErrCreateClient:   "Failed to create client",
		ErrGenerate:       "Failed to generate",
		ErrNoStagedChanges: "No staged changes found. Use 'git add' to stage changes",
		
		// Commands
		Summarizing:       "Summarizing with %s...",
		Translating:       "Translating to %s with %s...",
		Proofreading:      "Proofreading with %s (%s style)...",
		Generating:        "Generating commit message with %s...",
		Analyzing:         "Analyzing error with %s...",
		TokensUsed:        "[Tokens used: %d]",
		
		// Commit command
		SuggestedCommit:   "--- Suggested Commit Message ---",
		CommitInstruction: "To use this message, run:\n  git commit -m \"<copy message above>\"",
	})
}
