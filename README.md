# 🤖 aitxt - AI-Powered Text Processing CLI Tool

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/hiroki-abe-58/aitxt/actions/workflows/build.yml/badge.svg)](https://github.com/hiroki-abe-58/aitxt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiroki-abe-58/aitxt)](https://goreportcard.com/report/github.com/hiroki-abe-58/aitxt)

English | [日本語](README_ja.md) | [中文](README_zh.md)

A powerful command-line tool for AI-powered text processing with support for multiple LLM providers (OpenAI, Claude, Gemini).

## ✨ Features

### 🎯 Core Commands
- **ask** - General AI questions and assistance
- **summarize** - Text summarization with length control
- **translate** - Multi-language translation (13 languages)
- **proofread** - Grammar and spelling correction with style options
- **style** - Transform text style (formal, casual, technical, etc.)
- **explain** - Error message explanations with solutions
- **review** - Code review with security, performance, and style analysis
- **doc** - Generate documentation from code

### 💼 Developer Tools
- **commit** - Generate conventional commit messages from git diff
- **chat** - Interactive conversation mode with history
- **batch** - Process multiple files with parallel execution

### ⚙️ Productivity Features
- **template** - Reusable prompt templates with variables
- **alias** - Command shortcuts for frequent operations
- **history** - Search and manage conversation history
- **config** - View configuration and API key status

### 🌍 Multi-Language Support
- **UI Languages**: English, Japanese, Chinese, Korean, Thai
- **Translation**: 13 languages including Spanish, French, German, Portuguese, Russian, Arabic, Hindi, Vietnamese

### 🔧 Technical Features
- **3 LLM Providers**: OpenAI (GPT-4o), Claude (Sonnet 4), Gemini (1.5 Flash)
- **Parallel Processing**: Batch operations with configurable concurrency
- **Output Formats**: text, JSON, YAML
- **Shell Completion**: bash, zsh, fish, powershell
- **Progress Indicators**: Spinners and progress bars
- **Streaming Output**: Real-time response display

## 📦 Installation

### From Release (Recommended)

Download the latest release for your platform:
```bash
# macOS (Apple Silicon)
curl -L https://github.com/hiroki-abe-58/aitxt/releases/latest/download/aitxt-darwin-arm64.tar.gz | tar xz
sudo mv aitxt-darwin-arm64 /usr/local/bin/aitxt

# macOS (Intel)
curl -L https://github.com/hiroki-abe-58/aitxt/releases/latest/download/aitxt-darwin-amd64.tar.gz | tar xz
sudo mv aitxt-darwin-amd64 /usr/local/bin/aitxt

# Linux (x86_64)
curl -L https://github.com/hiroki-abe-58/aitxt/releases/latest/download/aitxt-linux-amd64.tar.gz | tar xz
sudo mv aitxt-linux-amd64 /usr/local/bin/aitxt
```

### From Source
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
make build
sudo make install
```

## 🚀 Quick Start

### 1. Set up API Keys
```bash
# Choose your preferred provider
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-claude-key"
export GOOGLE_API_KEY="your-gemini-key"
```

### 2. Verify Installation
```bash
aitxt config
aitxt version
```

### 3. Try Basic Commands
```bash
# Ask a question
aitxt ask "What is quantum computing?"

# Summarize a file
aitxt summarize document.txt

# Translate text
aitxt translate "Hello, world!" --to ja

# Generate commit message
git add .
aitxt commit
```

## 📚 Usage Examples

### Text Processing
```bash
# Summarize with custom length
aitxt summarize article.txt --max-length 200

# Translate to multiple languages
aitxt translate "Good morning" --to ja,zh,ko

# Proofread with formal style
aitxt proofread email.txt --style formal

# Transform writing style
aitxt style letter.txt --style professional
```

### Developer Tools
```bash
# Code review with focus
aitxt review main.go --focus security
aitxt review api.py --focus performance

# Generate documentation
aitxt doc src/ --type readme --format markdown

# Explain error messages
aitxt explain "TypeError: undefined is not a function"

# Generate commit messages
git add .
aitxt commit --lang ja --type feat
```

### Batch Processing
```bash
# Translate multiple files
aitxt batch *.md --operation translate --lang ja --concurrency 5

# Review all Python files
aitxt batch **/*.py --operation review --focus security

# Summarize documents
aitxt batch docs/*.txt --operation summarize --output summaries/
```

### Interactive Mode
```bash
# Start chat session
aitxt chat

# Chat with specific provider
aitxt chat --provider claude --model claude-sonnet-4-20250514

# Available chat commands:
# /help      - Show help
# /clear     - Clear history
# /system    - Change system message
# /provider  - Switch provider
# /history   - Show conversation history
# /save      - Save conversation
# /exit      - Exit chat
```

### Templates
```bash
# Initialize built-in templates
aitxt template init

# List templates
aitxt template list

# Use a template
aitxt template use code-review-security --var Language=Go --var Code="$(cat main.go)"

# Create custom template
aitxt template add --name my-review \
  --prompt "Review {{.Code}} for {{.Focus}}" \
  --system "You are a code reviewer"

# Export/import templates
aitxt template export my-templates.json
aitxt template import shared-templates.json
```

### Aliases
```bash
# Initialize suggested aliases
aitxt alias init

# Use aliases (after init)
aitxt s document.txt          # summarize
aitxt tj "Hello"               # translate to Japanese
aitxt rs main.go              # security review

# Create custom aliases
aitxt alias add fix "review --focus bugs"
aitxt fix broken-code.js

# List all aliases
aitxt alias list
```

### History Management
```bash
# List recent history
aitxt history list --limit 10

# Search history
aitxt history search "translation"

# Show specific entry
aitxt history show <id>

# Export history
aitxt history export history.json --format json

# View statistics
aitxt history stats
```

## 🔧 Configuration

### Environment Variables
```bash
# API Keys (required)
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# Optional Settings
export AITXT_PROVIDER="openai"        # Default provider
export AITXT_MODEL="gpt-4o"           # Default model
export AITXT_MAX_TOKENS="1000"        # Token limit
export AITXT_TEMPERATURE="0.7"        # Creativity (0.0-2.0)
export LANG="ja_JP.UTF-8"             # UI language
```

### Shell Completion
```bash
# bash
aitxt completion bash > /etc/bash_completion.d/aitxt

# zsh
aitxt completion zsh > /usr/local/share/zsh/site-functions/_aitxt

# fish
aitxt completion fish > ~/.config/fish/completions/aitxt.fish
```

## 🏗️ Project Structure
```
aitxt/
├── cmd/                    # CLI commands
│   ├── ask.go
│   ├── summarize.go
│   ├── translate.go
│   └── ...
├── pkg/                    # Core packages
│   ├── llm/               # LLM provider implementations
│   ├── config/            # Configuration management
│   ├── i18n/              # Internationalization
│   ├── template/          # Template system
│   ├── alias/             # Alias management
│   ├── history/           # History tracking
│   ├── batch/             # Batch processing
│   ├── progress/          # Progress indicators
│   └── output/            # Output formatting
├── .github/workflows/     # CI/CD pipelines
├── completions/           # Shell completion scripts
├── Makefile              # Build automation
└── README.md
```

## 🧪 Development

### Requirements

- Go 1.23 or higher
- Make (optional, for using Makefile)

### Build
```bash
# Build for current platform
make build

# Build for all platforms
make build-all

# Run tests
make test

# Run linter
golangci-lint run
```

### Testing
```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/config/
```

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for guidelines.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Run tests and linter
6. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) CLI framework
- OpenAI, Anthropic, and Google for LLM APIs
- Go community for excellent libraries

## 📊 Statistics

- **17 Commands**: Comprehensive text processing toolkit
- **3 LLM Providers**: OpenAI, Claude, Gemini
- **5 UI Languages**: en, ja, zh, ko, th
- **13 Translation Languages**: Multi-language support
- **85%+ Test Coverage**: Well-tested codebase
- **CI/CD**: Automated build and release

## 🔗 Links

- [GitHub Repository](https://github.com/hiroki-abe-58/aitxt)
- [Issue Tracker](https://github.com/hiroki-abe-58/aitxt/issues)
- [Releases](https://github.com/hiroki-abe-58/aitxt/releases)

## 💡 Use Cases

- **Content Creation**: Summarize articles, translate content, proofread writing
- **Development**: Generate commit messages, review code, create documentation
- **Learning**: Explain errors, ask questions, get code examples
- **Automation**: Batch processing, scripting, CI/CD integration
- **Research**: Process documents, extract information, analyze text

---

Made with ❤️ by [hiroki-abe-58](https://github.com/hiroki-abe-58)
