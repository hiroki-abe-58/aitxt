# 🤖 aitxt

AI-powered text processing CLI tool for developers.

## ✨ Features

- **Multiple LLM Providers**: OpenAI, Claude (Anthropic), Google Gemini
- **Text Processing**: Summarize, translate, proofread, and transform text
- **Developer Tools**: Code review, documentation generation, commit message creation
- **No Local Storage**: All processing via API (zero disk space required)
- **Fast & Efficient**: Optimized for Apple Silicon (M1/M2/M3)
- **Easy to Use**: Simple command-line interface

## 🚀 Installation
```bash
go install github.com/hiroki-abe-58/aitxt@latest
```

Or build from source:
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
go build -o aitxt main.go
```

## 📋 Usage

### Text Processing
```bash
# Summarize text
aitxt summarize document.txt

# Translate text
aitxt translate "Hello, world!" --to ja

# Proofread and fix grammar
aitxt proofread essay.txt
```

### Developer Tools
```bash
# Code review
aitxt review code.go

# Generate commit message from staged changes
aitxt commit
```

## ⚙️ Configuration

Set your API key as an environment variable:
```bash
export OPENAI_API_KEY="your-key-here"
```

## 📦 Requirements

- Go 1.23 or higher
- API key for at least one LLM provider

## 📄 License

MIT License - see [LICENSE](LICENSE) file for details.

## 🔗 Related Projects

- [imgai](https://github.com/hiroki-abe-58/imgai) - AI-powered image processing CLI tool
- [ghstat](https://github.com/hiroki-abe-58/ghstat) - GitHub statistics visualization tool
