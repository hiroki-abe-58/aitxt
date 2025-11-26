package i18n

func init() {
	registerKorean()
}

// registerKorean registers Korean messages
func registerKorean() {
	Register(LangKorean, &Messages{
		// General
		AppDescription:    "🤖 AI 기반 텍스트 처리 CLI 도구",
		AppLongDesc:       "aitxt는 여러 LLM 제공업체를 사용하여 강력한 텍스트 처리 기능을 제공합니다",
		Version:           "버전",
		
		// Errors
		ErrNoInput:        "입력 텍스트가 없습니다",
		ErrReadFile:       "파일 읽기 실패",
		ErrReadStdin:      "표준 입력 읽기 실패",
		ErrLoadConfig:     "설정 로드 실패",
		ErrCreateClient:   "클라이언트 생성 실패",
		ErrGenerate:       "생성 실패",
		ErrNoStagedChanges: "스테이징된 변경 사항이 없습니다. 'git add'를 사용하여 변경 사항을 스테이징하세요",
		
		// Commands
		Summarizing:       "%s로 요약 중...",
		Translating:       "%s로 %s(으)로 번역 중...",
		Proofreading:      "%s로 교정 중 (%s 스타일)...",
		Generating:        "%s로 커밋 메시지 생성 중...",
		Analyzing:         "%s로 오류 분석 중...",
		TokensUsed:        "[사용된 토큰: %d]",
		
		// Commit command
		SuggestedCommit:   "--- 제안된 커밋 메시지 ---",
		CommitInstruction: "이 메시지를 사용하려면:\n  git commit -m \"<위의 메시지를 복사>\"",
	})
}
