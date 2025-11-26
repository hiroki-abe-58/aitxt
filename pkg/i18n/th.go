package i18n

func init() {
	registerThai()
}

// registerThai registers Thai messages
func registerThai() {
	Register(LangThai, &Messages{
		// General
		AppDescription:    "🤖 เครื่องมือ CLI ประมวลผลข้อความด้วย AI",
		AppLongDesc:       "aitxt ให้ความสามารถในการประมวลผลข้อความที่ทรงพลังโดยใช้ผู้ให้บริการ LLM หลายราย",
		Version:           "เวอร์ชัน",
		
		// Errors
		ErrNoInput:        "ไม่มีข้อความที่ป้อน",
		ErrReadFile:       "อ่านไฟล์ไม่สำเร็จ",
		ErrReadStdin:      "อ่าน stdin ไม่สำเร็จ",
		ErrLoadConfig:     "โหลดการตั้งค่าไม่สำเร็จ",
		ErrCreateClient:   "สร้างไคลเอนต์ไม่สำเร็จ",
		ErrGenerate:       "สร้างไม่สำเร็จ",
		ErrNoStagedChanges: "ไม่พบการเปลี่ยนแปลงที่ staged ใช้ 'git add' เพื่อ stage การเปลี่ยนแปลง",
		
		// Commands
		Summarizing:       "กำลังสรุปด้วย %s...",
		Translating:       "กำลังแปลเป็น %s ด้วย %s...",
		Proofreading:      "กำลังตรวจสอบด้วย %s (สไตล์ %s)...",
		Generating:        "กำลังสร้างข้อความ commit ด้วย %s...",
		Analyzing:         "กำลังวิเคราะห์ข้อผิดพลาดด้วย %s...",
		TokensUsed:        "[โทเค็นที่ใช้: %d]",
		
		// Commit command
		SuggestedCommit:   "--- ข้อความ Commit ที่แนะนำ ---",
		CommitInstruction: "หากต้องการใช้ข้อความนี้ ให้รัน:\n  git commit -m \"<คัดลอกข้อความด้านบน>\"",
	})
}
