package gui

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/hiroki-abe-58/aitxt/pkg/config"
	"github.com/hiroki-abe-58/aitxt/pkg/llm"
)

// Handlers contains all API handlers
type Handlers struct{}

// NewHandlers creates new handlers
func NewHandlers() *Handlers {
	return &Handlers{}
}

// Register registers all API routes
func (h *Handlers) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/config", h.handleConfig)
	mux.HandleFunc("/api/ask", h.handleAsk)
	mux.HandleFunc("/api/translate", h.handleTranslate)
	mux.HandleFunc("/api/summarize", h.handleSummarize)
	mux.HandleFunc("/api/proofread", h.handleProofread)
	mux.HandleFunc("/api/style", h.handleStyle)
	mux.HandleFunc("/api/explain", h.handleExplain)
	mux.HandleFunc("/api/review", h.handleReview)
	mux.HandleFunc("/api/doc", h.handleDoc)
	mux.HandleFunc("/api/chat", h.handleChat)
}

// APIRequest represents a generic API request
type APIRequest struct {
	Text        string  `json:"text"`
	Provider    string  `json:"provider"`
	Stream      bool    `json:"stream"`
	Language    string  `json:"language"`
	FromLang    string  `json:"fromLang"`
	ToLang      string  `json:"toLang"`
	MaxLength   int     `json:"maxLength"`
	Style       string  `json:"style"`
	Focus       string  `json:"focus"`
	DocType     string  `json:"docType"`
	Temperature float64 `json:"temperature"`
	SystemMsg   string  `json:"systemMsg"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Success    bool   `json:"success"`
	Text       string `json:"text,omitempty"`
	Error      string `json:"error,omitempty"`
	Provider   string `json:"provider,omitempty"`
	Model      string `json:"model,omitempty"`
	TokensUsed int    `json:"tokensUsed,omitempty"`
}

// ConfigResponse represents config API response
type ConfigResponse struct {
	Provider           string   `json:"provider"`
	AvailableProviders []string `json:"availableProviders"`
	OpenAIConfigured   bool     `json:"openaiConfigured"`
	ClaudeConfigured   bool     `json:"claudeConfigured"`
	GeminiConfigured   bool     `json:"geminiConfigured"`
}

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, APIResponse{
		Success: false,
		Error:   err.Error(),
	})
}

func (h *Handlers) handleConfig(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	cfg, err := config.Load()
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	resp := ConfigResponse{
		Provider:           string(cfg.Provider),
		AvailableProviders: []string{},
	}

	// Check which providers are configured
	if cfg.OpenAIKey != "" {
		resp.OpenAIConfigured = true
		resp.AvailableProviders = append(resp.AvailableProviders, "openai")
	}
	if cfg.ClaudeKey != "" {
		resp.ClaudeConfigured = true
		resp.AvailableProviders = append(resp.AvailableProviders, "claude")
	}
	if cfg.GeminiKey != "" {
		resp.GeminiConfigured = true
		resp.AvailableProviders = append(resp.AvailableProviders, "gemini")
	}

	writeJSON(w, http.StatusOK, resp)
}

func (h *Handlers) createClient(provider string) (llm.Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	p := cfg.Provider
	if provider != "" {
		p = llm.Provider(provider)
	}

	llmConfig, err := cfg.ToLLMConfig(p)
	if err != nil {
		return nil, err
	}

	factory := llm.NewFactory()
	if err := factory.RegisterConfig(llmConfig); err != nil {
		return nil, fmt.Errorf("failed to register config: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return factory.CreateClientWithContext(ctx, p)
}

func (h *Handlers) handleAsk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	if req.Stream {
		h.handleStreamRequest(w, r, req, "You are a helpful AI assistant. Provide clear, accurate, and concise answers.")
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	systemMsg := "You are a helpful AI assistant. Provide clear, accurate, and concise answers."
	if req.SystemMsg != "" {
		systemMsg = req.SystemMsg
	}
	if req.Language != "" {
		systemMsg += " " + getLanguagePrompt(req.Language)
	}

	temperature := 0.7
	if req.Temperature > 0 {
		temperature = req.Temperature
	}

	llmReq := &llm.Request{
		Prompt:      req.Text,
		SystemMsg:   systemMsg,
		MaxTokens:   2000,
		Temperature: temperature,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleTranslate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	toLang := req.ToLang
	if toLang == "" {
		toLang = "en"
	}

	if req.Stream {
		systemMsg := buildTranslateSystemMsg(req.FromLang, toLang)
		h.handleStreamRequest(w, r, APIRequest{
			Text:     req.Text,
			Provider: req.Provider,
			Stream:   true,
		}, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	systemMsg := buildTranslateSystemMsg(req.FromLang, toLang)

	llmReq := &llm.Request{
		Prompt:    req.Text,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleSummarize(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	maxLength := req.MaxLength
	if maxLength <= 0 {
		maxLength = 200
	}

	systemMsg := fmt.Sprintf("You are a professional summarizer. Summarize the following text concisely in about %d words. Preserve the key points and main ideas.", maxLength)

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    req.Text,
		SystemMsg: systemMsg,
		MaxTokens: 1000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleProofread(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	style := req.Style
	if style == "" {
		style = "standard"
	}

	systemMsg := fmt.Sprintf("You are a professional proofreader. Check and correct the following text for grammar, spelling, and punctuation errors. Use a %s tone. Return only the corrected text without explanations.", style)

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    req.Text,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleStyle(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	style := req.Style
	if style == "" {
		style = "professional"
	}

	systemMsg := fmt.Sprintf("You are a professional writer. Rewrite the following text in a %s style while preserving the original meaning. Return only the rewritten text.", style)

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    req.Text,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleExplain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	systemMsg := "You are an expert programmer and technical support specialist. Explain the following error message in simple terms, identify the likely cause, and provide practical solutions to fix it."

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    req.Text,
		SystemMsg: systemMsg,
		MaxTokens: 2000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleReview(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	focus := req.Focus
	if focus == "" {
		focus = "general"
	}

	var systemMsg string
	switch focus {
	case "security":
		systemMsg = "You are a security expert. Review the following code for security vulnerabilities, potential exploits, and security best practices. Provide specific recommendations."
	case "performance":
		systemMsg = "You are a performance optimization expert. Review the following code for performance issues, bottlenecks, and optimization opportunities. Provide specific recommendations."
	case "style":
		systemMsg = "You are a code style expert. Review the following code for coding style, readability, and maintainability. Suggest improvements following best practices."
	case "bugs":
		systemMsg = "You are a debugging expert. Review the following code for potential bugs, logic errors, and edge cases. Identify issues and suggest fixes."
	default:
		systemMsg = "You are an expert code reviewer. Review the following code for bugs, security issues, performance problems, and style improvements. Provide a comprehensive review with specific recommendations."
	}

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    fmt.Sprintf("```\n%s\n```", req.Text),
		SystemMsg: systemMsg,
		MaxTokens: 3000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleDoc(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	docType := req.DocType
	if docType == "" {
		docType = "inline"
	}

	var systemMsg string
	switch docType {
	case "readme":
		systemMsg = "You are a technical writer. Generate a comprehensive README.md file for the following code. Include sections for description, installation, usage, API reference, and examples."
	case "api":
		systemMsg = "You are a technical writer. Generate API documentation for the following code. Include function signatures, parameters, return values, and usage examples."
	default:
		systemMsg = "You are a technical writer. Add comprehensive inline documentation (comments) to the following code. Explain the purpose of functions, parameters, and complex logic."
	}

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	llmReq := &llm.Request{
		Prompt:    fmt.Sprintf("```\n%s\n```", req.Text),
		SystemMsg: systemMsg,
		MaxTokens: 4000,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

func (h *Handlers) handleChat(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, fmt.Errorf("method not allowed"))
		return
	}

	var req APIRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if req.Text == "" {
		writeError(w, http.StatusBadRequest, fmt.Errorf("text is required"))
		return
	}

	systemMsg := "You are a helpful AI assistant. Engage in natural conversation and provide helpful, accurate responses."
	if req.SystemMsg != "" {
		systemMsg = req.SystemMsg
	}

	if req.Stream {
		h.handleStreamRequest(w, r, req, systemMsg)
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	temperature := 0.7
	if req.Temperature > 0 {
		temperature = req.Temperature
	}

	llmReq := &llm.Request{
		Prompt:      req.Text,
		SystemMsg:   systemMsg,
		MaxTokens:   2000,
		Temperature: temperature,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	resp, err := client.Generate(ctx, llmReq)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}

	writeJSON(w, http.StatusOK, APIResponse{
		Success:    true,
		Text:       resp.Text,
		Provider:   string(resp.Provider),
		Model:      resp.Model,
		TokensUsed: resp.TokensUsed,
	})
}

// SSE streaming handler
func (h *Handlers) handleStreamRequest(w http.ResponseWriter, r *http.Request, req APIRequest, systemMsg string) {
	// Set SSE headers
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		writeError(w, http.StatusInternalServerError, fmt.Errorf("streaming not supported"))
		return
	}

	client, err := h.createClient(req.Provider)
	if err != nil {
		fmt.Fprintf(w, "data: {\"error\": %q}\n\n", err.Error())
		flusher.Flush()
		return
	}

	temperature := 0.7
	if req.Temperature > 0 {
		temperature = req.Temperature
	}

	llmReq := &llm.Request{
		Prompt:      req.Text,
		SystemMsg:   systemMsg,
		MaxTokens:   2000,
		Temperature: temperature,
		Stream:      true,
	}

	ctx, cancel := context.WithTimeout(r.Context(), 120*time.Second)
	defer cancel()

	err = client.Stream(ctx, llmReq, func(chunk string) error {
		data, _ := json.Marshal(map[string]string{"text": chunk})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return nil
	})

	if err != nil {
		data, _ := json.Marshal(map[string]string{"error": err.Error()})
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
		return
	}

	// Send done event
	fmt.Fprintf(w, "data: {\"done\": true}\n\n")
	flusher.Flush()
}

// Helper functions
func getLanguagePrompt(lang string) string {
	switch lang {
	case "ja":
		return "Respond in Japanese."
	case "zh":
		return "Respond in Chinese."
	case "ko":
		return "Respond in Korean."
	case "th":
		return "Respond in Thai."
	case "en":
		return "Respond in English."
	case "es":
		return "Respond in Spanish."
	case "fr":
		return "Respond in French."
	case "de":
		return "Respond in German."
	default:
		return ""
	}
}

func buildTranslateSystemMsg(fromLang, toLang string) string {
	targetLang := getLanguageName(toLang)
	if fromLang != "" {
		sourceLang := getLanguageName(fromLang)
		return fmt.Sprintf("You are a professional translator. Translate the text from %s to %s. Output only the translation, no explanations.", sourceLang, targetLang)
	}
	return fmt.Sprintf("You are a professional translator. Detect the source language and translate the text to %s. Output only the translation, no explanations.", targetLang)
}

func getLanguageName(code string) string {
	languages := map[string]string{
		"en": "English",
		"ja": "Japanese",
		"zh": "Chinese",
		"ko": "Korean",
		"th": "Thai",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
		"pt": "Portuguese",
		"ru": "Russian",
		"ar": "Arabic",
		"hi": "Hindi",
		"vi": "Vietnamese",
	}
	if name, ok := languages[code]; ok {
		return name
	}
	return code
}
