package i18n

func init() {
	registerChinese()
}

// registerChinese registers Chinese messages
func registerChinese() {
	Register(LangChinese, &Messages{
		// General
		AppDescription:    "🤖 AI驱动的文本处理CLI工具",
		AppLongDesc:       "aitxt使用多个LLM提供商提供强大的文本处理功能",
		Version:           "版本",
		
		// Errors
		ErrNoInput:        "没有提供输入文本",
		ErrReadFile:       "读取文件失败",
		ErrReadStdin:      "读取标准输入失败",
		ErrLoadConfig:     "加载配置失败",
		ErrCreateClient:   "创建客户端失败",
		ErrGenerate:       "生成失败",
		ErrNoStagedChanges: "没有暂存的更改。使用 'git add' 暂存更改",
		
		// Commands
		Summarizing:       "使用 %s 进行摘要...",
		Translating:       "使用 %s 翻译成%s...",
		Proofreading:      "使用 %s 进行校对（%s风格）...",
		Generating:        "使用 %s 生成提交信息...",
		Analyzing:         "使用 %s 分析错误...",
		TokensUsed:        "[使用的令牌数: %d]",
		
		// Commit command
		SuggestedCommit:   "--- 建议的提交信息 ---",
		CommitInstruction: "要使用此信息，请运行:\n  git commit -m \"<复制上面的信息>\"",
	})
}
