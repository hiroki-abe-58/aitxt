function app() {
    return {
        // State
        currentTab: 'ask',
        darkMode: false,
        loading: false,
        output: '',
        error: '',
        tokenCount: 0,
        streamOutput: '',
        
        // Config
        config: {
            provider: 'openai',
            availableProviders: [],
            openaiConfigured: false,
            claudeConfigured: false,
            geminiConfigured: false
        },
        selectedProvider: '',
        streamEnabled: false,
        
        // Settings
        settingsTab: 'openai',
        settingsSaving: false,
        settingsSaved: false,
        providerSettings: {
            openai: { temperature: 0.7, top_p: 1.0, top_k: 0, max_tokens: 2000, model: 'gpt-4o' },
            gemini: { temperature: 0.7, top_p: 0.95, top_k: 40, max_tokens: 2000, model: 'gemini-1.5-flash' },
            claude: { temperature: 0.7, top_p: 0.9, top_k: 0, max_tokens: 2000, model: 'claude-sonnet-4-20250514' }
        },
        
        // Menu Items
        menuItems: [
            { id: 'ask', label: 'Ask', icon: 'chat', description: 'Ask AI any question' },
            { id: 'translate', label: 'Translate', icon: 'translate', description: 'Translate text between languages' },
            { id: 'summarize', label: 'Summarize', icon: 'summarize', description: 'Summarize long texts' },
            { id: 'proofread', label: 'Proofread', icon: 'spellcheck', description: 'Check grammar and spelling' },
            { id: 'style', label: 'Style', icon: 'style', description: 'Transform writing style' },
            { id: 'explain', label: 'Explain', icon: 'help_outline', description: 'Explain error messages' },
            { id: 'review', label: 'Review', icon: 'rate_review', description: 'Review code for issues' },
            { id: 'doc', label: 'Document', icon: 'description', description: 'Generate documentation' },
            { id: 'chat', label: 'Chat', icon: 'forum', description: 'Interactive conversation' },
            { id: 'settings', label: 'Settings', icon: 'settings', description: 'Configure providers' }
        ],
        
        // Form inputs
        askInput: '',
        askLang: '',
        askTemperature: 0.7,
        
        translateInput: '',
        translateFromLang: '',
        translateToLang: 'ja',
        
        summarizeInput: '',
        summarizeMaxLength: 200,
        
        proofreadInput: '',
        proofreadStyle: 'standard',
        
        styleInput: '',
        styleTarget: 'professional',
        
        explainInput: '',
        
        reviewInput: '',
        reviewFocus: 'general',
        
        docInput: '',
        docType: 'inline',
        
        chatInput: '',
        chatHistory: [],
        
        // Computed
        get currentMenuItem() {
            return this.menuItems.find(item => item.id === this.currentTab) || this.menuItems[0];
        },
        
        // Methods
        async init() {
            // Load dark mode preference
            this.darkMode = localStorage.getItem('darkMode') === 'true' || 
                           window.matchMedia('(prefers-color-scheme: dark)').matches;
            this.updateDarkMode();
            
            // Load config
            await this.loadConfig();
            
            // Load settings
            await this.loadProviderSettings();
        },
        
        toggleDarkMode() {
            this.darkMode = !this.darkMode;
            localStorage.setItem('darkMode', this.darkMode);
            this.updateDarkMode();
        },
        
        updateDarkMode() {
            if (this.darkMode) {
                document.documentElement.classList.add('dark');
            } else {
                document.documentElement.classList.remove('dark');
            }
        },
        
        async loadConfig() {
            try {
                const response = await fetch('/api/config');
                if (response.ok) {
                    this.config = await response.json();
                    if (this.config.availableProviders.length > 0) {
                        this.selectedProvider = this.config.provider || this.config.availableProviders[0];
                    }
                }
            } catch (err) {
                console.error('Failed to load config:', err);
            }
        },
        
        async submitRequest(endpoint, body) {
            this.loading = true;
            this.output = '';
            this.error = '';
            this.tokenCount = 0;
            this.streamOutput = '';
            
            try {
                if (this.streamEnabled) {
                    await this.submitStreamRequest(endpoint, body);
                } else {
                    const response = await fetch(`/api/${endpoint}`, {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            ...body,
                            provider: this.selectedProvider,
                            stream: false
                        })
                    });
                    
                    const data = await response.json();
                    
                    if (data.success) {
                        this.output = data.text;
                        this.tokenCount = data.tokensUsed || 0;
                    } else {
                        this.error = data.error || 'An error occurred';
                    }
                }
            } catch (err) {
                this.error = err.message || 'Network error';
            } finally {
                this.loading = false;
            }
        },
        
        async submitStreamRequest(endpoint, body) {
            const response = await fetch(`/api/${endpoint}`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({
                    ...body,
                    provider: this.selectedProvider,
                    stream: true
                })
            });
            
            const reader = response.body.getReader();
            const decoder = new TextDecoder();
            let buffer = '';
            
            while (true) {
                const { done, value } = await reader.read();
                if (done) break;
                
                buffer += decoder.decode(value, { stream: true });
                const lines = buffer.split('\n');
                buffer = lines.pop() || '';
                
                for (const line of lines) {
                    if (line.startsWith('data: ')) {
                        try {
                            const data = JSON.parse(line.slice(6));
                            if (data.error) {
                                this.error = data.error;
                                return;
                            }
                            if (data.done) {
                                return;
                            }
                            if (data.text) {
                                this.output += data.text;
                                this.streamOutput = this.output;
                            }
                        } catch (e) {
                            // Ignore parse errors
                        }
                    }
                }
            }
        },
        
        async submitAsk() {
            await this.submitRequest('ask', {
                text: this.askInput,
                language: this.askLang,
                temperature: parseFloat(this.askTemperature) || 0.7
            });
        },
        
        async submitTranslate() {
            await this.submitRequest('translate', {
                text: this.translateInput,
                fromLang: this.translateFromLang,
                toLang: this.translateToLang
            });
        },
        
        async submitSummarize() {
            await this.submitRequest('summarize', {
                text: this.summarizeInput,
                maxLength: parseInt(this.summarizeMaxLength) || 200
            });
        },
        
        async submitProofread() {
            await this.submitRequest('proofread', {
                text: this.proofreadInput,
                style: this.proofreadStyle
            });
        },
        
        async submitStyle() {
            await this.submitRequest('style', {
                text: this.styleInput,
                style: this.styleTarget
            });
        },
        
        async submitExplain() {
            await this.submitRequest('explain', {
                text: this.explainInput
            });
        },
        
        async submitReview() {
            await this.submitRequest('review', {
                text: this.reviewInput,
                focus: this.reviewFocus
            });
        },
        
        async submitDoc() {
            await this.submitRequest('doc', {
                text: this.docInput,
                docType: this.docType
            });
        },
        
        async submitChat() {
            if (!this.chatInput.trim()) return;
            
            const userMessage = this.chatInput.trim();
            this.chatHistory.push({ role: 'user', content: userMessage });
            this.chatInput = '';
            this.loading = true;
            this.streamOutput = '';
            
            // Scroll to bottom
            this.$nextTick(() => {
                const chatEl = document.getElementById('chatMessages');
                if (chatEl) chatEl.scrollTop = chatEl.scrollHeight;
            });
            
            try {
                if (this.streamEnabled) {
                    let assistantMessage = '';
                    
                    const response = await fetch('/api/chat', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            text: userMessage,
                            provider: this.selectedProvider,
                            stream: true
                        })
                    });
                    
                    const reader = response.body.getReader();
                    const decoder = new TextDecoder();
                    let buffer = '';
                    
                    while (true) {
                        const { done, value } = await reader.read();
                        if (done) break;
                        
                        buffer += decoder.decode(value, { stream: true });
                        const lines = buffer.split('\n');
                        buffer = lines.pop() || '';
                        
                        for (const line of lines) {
                            if (line.startsWith('data: ')) {
                                try {
                                    const data = JSON.parse(line.slice(6));
                                    if (data.error) {
                                        this.error = data.error;
                                        return;
                                    }
                                    if (data.done) {
                                        this.chatHistory.push({ role: 'assistant', content: assistantMessage });
                                        return;
                                    }
                                    if (data.text) {
                                        assistantMessage += data.text;
                                        this.streamOutput = assistantMessage;
                                    }
                                } catch (e) {
                                    // Ignore parse errors
                                }
                            }
                        }
                    }
                    
                    if (assistantMessage) {
                        this.chatHistory.push({ role: 'assistant', content: assistantMessage });
                    }
                } else {
                    const response = await fetch('/api/chat', {
                        method: 'POST',
                        headers: { 'Content-Type': 'application/json' },
                        body: JSON.stringify({
                            text: userMessage,
                            provider: this.selectedProvider,
                            stream: false
                        })
                    });
                    
                    const data = await response.json();
                    
                    if (data.success) {
                        this.chatHistory.push({ role: 'assistant', content: data.text });
                    } else {
                        this.error = data.error || 'An error occurred';
                    }
                }
            } catch (err) {
                this.error = err.message || 'Network error';
            } finally {
                this.loading = false;
                this.streamOutput = '';
                
                // Scroll to bottom
                this.$nextTick(() => {
                    const chatEl = document.getElementById('chatMessages');
                    if (chatEl) chatEl.scrollTop = chatEl.scrollHeight;
                });
            }
        },
        
        clearChat() {
            this.chatHistory = [];
        },
        
        async loadProviderSettings() {
            try {
                const response = await fetch('/api/settings');
                if (response.ok) {
                    const data = await response.json();
                    this.providerSettings = {
                        openai: {
                            temperature: parseFloat(data.openai.temperature) || 0.7,
                            top_p: parseFloat(data.openai.top_p) || 1.0,
                            top_k: parseInt(data.openai.top_k) || 0,
                            max_tokens: parseInt(data.openai.max_tokens) || 2000,
                            model: data.openai.model || 'gpt-4o'
                        },
                        gemini: {
                            temperature: parseFloat(data.gemini.temperature) || 0.7,
                            top_p: parseFloat(data.gemini.top_p) || 0.95,
                            top_k: parseInt(data.gemini.top_k) || 40,
                            max_tokens: parseInt(data.gemini.max_tokens) || 2000,
                            model: data.gemini.model || 'gemini-1.5-flash'
                        },
                        claude: {
                            temperature: parseFloat(data.claude.temperature) || 0.7,
                            top_p: parseFloat(data.claude.top_p) || 0.9,
                            top_k: parseInt(data.claude.top_k) || 0,
                            max_tokens: parseInt(data.claude.max_tokens) || 2000,
                            model: data.claude.model || 'claude-sonnet-4-20250514'
                        }
                    };
                }
            } catch (err) {
                console.error('Failed to load settings:', err);
            }
        },
        
        async saveProviderSettings() {
            this.settingsSaving = true;
            this.settingsSaved = false;
            
            try {
                const response = await fetch('/api/settings', {
                    method: 'POST',
                    headers: { 'Content-Type': 'application/json' },
                    body: JSON.stringify({
                        openai: {
                            temperature: parseFloat(this.providerSettings.openai.temperature),
                            top_p: parseFloat(this.providerSettings.openai.top_p),
                            top_k: parseInt(this.providerSettings.openai.top_k),
                            max_tokens: parseInt(this.providerSettings.openai.max_tokens),
                            model: this.providerSettings.openai.model
                        },
                        gemini: {
                            temperature: parseFloat(this.providerSettings.gemini.temperature),
                            top_p: parseFloat(this.providerSettings.gemini.top_p),
                            top_k: parseInt(this.providerSettings.gemini.top_k),
                            max_tokens: parseInt(this.providerSettings.gemini.max_tokens),
                            model: this.providerSettings.gemini.model
                        },
                        claude: {
                            temperature: parseFloat(this.providerSettings.claude.temperature),
                            top_p: parseFloat(this.providerSettings.claude.top_p),
                            top_k: parseInt(this.providerSettings.claude.top_k),
                            max_tokens: parseInt(this.providerSettings.claude.max_tokens),
                            model: this.providerSettings.claude.model
                        }
                    })
                });
                
                const data = await response.json();
                
                if (data.success) {
                    this.settingsSaved = true;
                    setTimeout(() => { this.settingsSaved = false; }, 3000);
                } else {
                    this.error = data.error || 'Failed to save settings';
                }
            } catch (err) {
                this.error = err.message || 'Failed to save settings';
            } finally {
                this.settingsSaving = false;
            }
        },
        
        resetProviderSettings(provider) {
            const defaults = {
                openai: { temperature: 0.7, top_p: 1.0, top_k: 0, max_tokens: 2000, model: 'gpt-4o' },
                gemini: { temperature: 0.7, top_p: 0.95, top_k: 40, max_tokens: 2000, model: 'gemini-1.5-flash' },
                claude: { temperature: 0.7, top_p: 0.9, top_k: 0, max_tokens: 2000, model: 'claude-sonnet-4-20250514' }
            };
            
            if (defaults[provider]) {
                this.providerSettings[provider] = { ...defaults[provider] };
            }
        }
    };
}

