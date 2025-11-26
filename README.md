# 🤖 aitxt

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg)](https://github.com/hiroki-abe-58/aitxt/pulls)

AI-powered text processing CLI tool for developers.

## ✨ Features

- **Multiple LLM Providers** - OpenAI, Claude (Anthropic), Google Gemini
- **Text Processing** - Summarize, translate, proofread, and transform text
- **Developer Tools** - Code review, documentation generation, commit message creation
- **Error Analysis** - Explain error messages with solutions
- **Multi-language** - Support for English, Japanese, Chinese, Korean, Thai
- **Streaming Output** - Real-time response streaming
- **No Local Storage** - All processing via API (zero disk space for models)

## 📦 Installation

### From Source
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
go build -o aitxt main.go

# Optional: Install globally
sudo mv aitxt /usr/local/bin/
```

### Using Go Install
```bash
go install github.com/hiroki-abe-58/aitxt@latest
```

## ⚙️ Configuration

Set your API key as an environment variable:
```bash
# OpenAI (default)
export OPENAI_API_KEY="sk-..."

# Anthropic Claude
export ANTHROPIC_API_KEY="sk-ant-..."

# Google Gemini
export GOOGLE_API_KEY="..."
```

The tool auto-detects available providers based on set API keys.

## 🚀 Commands

### Text Processing

#### `summarize` - Summarize text
```bash
aitxt summarize document.txt
aitxt summarize article.txt --lang ja
cat long_text.txt | aitxt summarize -
```

#### `translate` - Translate text
```bash
aitxt translate "Hello, world!" --to ja
aitxt translate document.txt --to zh
aitxt translate japanese.txt --from ja --to en
```

Supported languages: `en`, `ja`, `zh`, `ko`, `th`, `es`, `fr`, `de`, `pt`, `ru`, `ar`, `hi`, `vi`

#### `proofread` - Proofread and correct text
```bash
aitxt proofread essay.txt
aitxt proofread email.txt --style formal
aitxt proofread report.txt --style academic --lang ja
```

Styles: `standard`, `formal`, `casual`, `academic`, `technical`

#### `style` - Transform writing style
```bash
aitxt style "your text" --to formal
aitxt style document.txt --to casual
aitxt style email.txt --to professional
```

Styles: `formal`, `casual`, `academic`, `technical`, `creative`, `simple`, `persuasive`, `humorous`, `poetic`, `journalistic`

### Developer Tools

#### `review` - AI code review
```bash
aitxt review main.go
aitxt review src/app.py --focus security
aitxt review handler.js --focus performance
git diff | aitxt review - --focus bugs
```

Focus areas: `all`, `security`, `performance`, `style`, `bugs`, `test`

#### `commit` - Generate commit messages
```bash
# Stage your changes first
git add .

# Generate commit message
aitxt commit
aitxt commit --lang ja
aitxt commit --type feat
```

#### `doc` - Generate documentation
```bash
aitxt doc main.go
aitxt doc src/api.py --type api
aitxt doc lib/ --type readme
aitxt doc handler.js --type tutorial --lang ja
```

Types: `auto`, `readme`, `api`, `inline`, `tutorial`, `changelog`

#### `explain` - Explain error messages
```bash
aitxt explain "undefined: fmt.Prinln"
aitxt explain error.log
aitxt explain "ECONNREFUSED 127.0.0.1:5432" --context "PostgreSQL"
```

### General

#### `ask` - Ask AI anything
```bash
aitxt ask "What is the capital of France?"
aitxt ask "Explain quantum computing" --lang ja
cat data.json | aitxt ask "Analyze this data"
aitxt ask "Be creative" --temperature 1.5
```

## 🌐 Multi-language Support

All commands support multiple output languages:
```bash
aitxt summarize doc.txt --lang ja    # Japanese
aitxt summarize doc.txt --lang zh    # Chinese
aitxt summarize doc.txt --lang ko    # Korean
aitxt summarize doc.txt --lang th    # Thai
```

## 🔄 Streaming Output

Enable real-time streaming for faster feedback:
```bash
aitxt ask "Tell me a story" --stream
aitxt summarize long_doc.txt --stream
aitxt review large_file.go --stream
```

## 🤝 Provider Selection

Specify which LLM provider to use:
```bash
aitxt ask "Hello" --provider openai
aitxt ask "Hello" --provider claude
aitxt ask "Hello" --provider gemini
```

## 📋 Examples

### Daily Workflow
```bash
# Morning: Summarize overnight emails
aitxt summarize emails.txt --lang ja

# Coding: Review your changes before commit
git diff | aitxt review - --focus bugs

# Commit: Generate meaningful commit messages
git add . && aitxt commit

# Debug: Understand error messages
aitxt explain "panic: runtime error: index out of range"

# Documentation: Generate docs for new code
aitxt doc new_feature.go --type api
```

### Writing Assistance
```bash
# Proofread your blog post
aitxt proofread blog_draft.md --style casual

# Translate for international audience
aitxt translate announcement.txt --to zh

# Make technical docs more accessible
aitxt style technical_doc.txt --to simple
```

## 🛠️ Development
```bash
# Clone repository
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt

# Build
go build -o aitxt main.go

# Run tests
go test ./...

# Run specific command
go run main.go ask "Hello"
```

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [imgai](https://github.com/hiroki-abe-58/imgai) - AI-powered image processing CLI tool
- [ghstat](https://github.com/hiroki-abe-58/ghstat) - GitHub statistics visualization tool

## 🙏 Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [OpenAI Go](https://github.com/sashabaranov/go-openai) - OpenAI client
- [Anthropic Go](https://github.com/anthropics/anthropic-sdk-go) - Claude client
- [Google Generative AI](https://github.com/google/generative-ai-go) - Gemini client

---

Made with ❤️ for developers who love the command line.
