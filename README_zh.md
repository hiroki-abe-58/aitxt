# 🤖 aitxt - AI 驱动的文本处理 CLI 工具

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/hiroki-abe-58/aitxt/actions/workflows/build.yml/badge.svg)](https://github.com/hiroki-abe-58/aitxt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiroki-abe-58/aitxt)](https://goreportcard.com/report/github.com/hiroki-abe-58/aitxt)

[English](README.md) | [日本語](README_ja.md) | 中文

一款强大的命令行工具，支持多个 LLM 提供商（OpenAI、Claude、Gemini）进行 AI 驱动的文本处理。

## ✨ 功能特性

### 🎯 核心命令
- **ask** - 通用 AI 问答和辅助
- **summarize** - 可控长度的文本摘要
- **translate** - 多语言翻译（13 种语言）
- **proofread** - 语法和拼写校正，支持多种风格选项
- **style** - 文本风格转换（正式、随意、技术性等）
- **explain** - 错误信息解释与解决方案
- **review** - 代码审查，包含安全性、性能和风格分析
- **doc** - 从代码生成文档

### 💼 开发者工具
- **commit** - 从 git diff 生成规范的提交信息
- **chat** - 带历史记录的交互式对话模式
- **batch** - 并行处理多个文件

### ⚙️ 效率功能
- **template** - 支持变量的可复用提示词模板
- **alias** - 常用操作的命令快捷方式
- **history** - 搜索和管理对话历史
- **config** - 查看配置和 API 密钥状态

### 🌍 多语言支持
- **界面语言**: 英语、日语、中文、韩语、泰语
- **翻译语言**: 13 种语言，包括西班牙语、法语、德语、葡萄牙语、俄语、阿拉伯语、印地语、越南语

### 🔧 技术特性
- **3 个 LLM 提供商**: OpenAI (GPT-4o)、Claude (Sonnet 4)、Gemini (1.5 Flash)
- **并行处理**: 可配置并发数的批量操作
- **输出格式**: text、JSON、YAML
- **Shell 自动补全**: bash、zsh、fish、powershell
- **进度指示器**: 旋转动画和进度条
- **流式输出**: 实时响应显示

## 📦 安装

### 从发布版安装（推荐）

下载适合您平台的最新版本：
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

### 从源码安装
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
make build
sudo make install
```

## 🚀 快速开始

### 1. 设置 API 密钥
```bash
# 选择您喜欢的提供商
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-claude-key"
export GOOGLE_API_KEY="your-gemini-key"
```

### 2. 验证安装
```bash
aitxt config
aitxt version
```

### 3. 尝试基本命令
```bash
# 提问
aitxt ask "什么是量子计算？"

# 摘要文件
aitxt summarize document.txt

# 翻译文本
aitxt translate "你好，世界！" --to en

# 生成提交信息
git add .
aitxt commit
```

## 📚 使用示例

### 文本处理
```bash
# 自定义长度的摘要
aitxt summarize article.txt --max-length 200

# 翻译到多种语言
aitxt translate "早上好" --to en,ja,ko

# 正式风格的校对
aitxt proofread email.txt --style formal

# 转换写作风格
aitxt style letter.txt --style professional
```

### 开发者工具
```bash
# 专注特定方面的代码审查
aitxt review main.go --focus security
aitxt review api.py --focus performance

# 生成文档
aitxt doc src/ --type readme --format markdown

# 解释错误信息
aitxt explain "TypeError: undefined is not a function"

# 生成提交信息
git add .
aitxt commit --lang zh --type feat
```

### 批量处理
```bash
# 翻译多个文件
aitxt batch *.md --operation translate --lang zh --concurrency 5

# 审查所有 Python 文件
aitxt batch **/*.py --operation review --focus security

# 摘要文档
aitxt batch docs/*.txt --operation summarize --output summaries/
```

### 交互模式
```bash
# 启动聊天会话
aitxt chat

# 使用特定提供商聊天
aitxt chat --provider claude --model claude-sonnet-4-20250514

# 可用的聊天命令：
# /help      - 显示帮助
# /clear     - 清除历史
# /system    - 更改系统消息
# /provider  - 切换提供商
# /history   - 显示对话历史
# /save      - 保存对话
# /exit      - 退出聊天
```

### 模板
```bash
# 初始化内置模板
aitxt template init

# 列出模板
aitxt template list

# 使用模板
aitxt template use code-review-security --var Language=Go --var Code="$(cat main.go)"

# 创建自定义模板
aitxt template add --name my-review \
  --prompt "审查 {{.Code}} 的 {{.Focus}}" \
  --system "你是一个代码审查员"

# 导出/导入模板
aitxt template export my-templates.json
aitxt template import shared-templates.json
```

### 别名
```bash
# 初始化建议的别名
aitxt alias init

# 使用别名（初始化后）
aitxt s document.txt          # summarize
aitxt tj "你好"                # 翻译到日语
aitxt rs main.go              # 安全审查

# 创建自定义别名
aitxt alias add fix "review --focus bugs"
aitxt fix broken-code.js

# 列出所有别名
aitxt alias list
```

### 历史管理
```bash
# 列出最近的历史
aitxt history list --limit 10

# 搜索历史
aitxt history search "翻译"

# 显示特定条目
aitxt history show <id>

# 导出历史
aitxt history export history.json --format json

# 查看统计
aitxt history stats
```

## 🔧 配置

### 环境变量
```bash
# API 密钥（必需）
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# 可选设置
export AITXT_PROVIDER="openai"        # 默认提供商
export AITXT_MODEL="gpt-4o"           # 默认模型
export AITXT_MAX_TOKENS="1000"        # Token 限制
export AITXT_TEMPERATURE="0.7"        # 创造性 (0.0-2.0)
export LANG="zh_CN.UTF-8"             # 界面语言
```

### Shell 自动补全
```bash
# bash
aitxt completion bash > /etc/bash_completion.d/aitxt

# zsh
aitxt completion zsh > /usr/local/share/zsh/site-functions/_aitxt

# fish
aitxt completion fish > ~/.config/fish/completions/aitxt.fish
```

## 🏗️ 项目结构
```
aitxt/
├── cmd/                    # CLI 命令
│   ├── ask.go
│   ├── summarize.go
│   ├── translate.go
│   └── ...
├── pkg/                    # 核心包
│   ├── llm/               # LLM 提供商实现
│   ├── config/            # 配置管理
│   ├── i18n/              # 国际化
│   ├── template/          # 模板系统
│   ├── alias/             # 别名管理
│   ├── history/           # 历史跟踪
│   ├── batch/             # 批量处理
│   ├── progress/          # 进度指示器
│   └── output/            # 输出格式化
├── .github/workflows/     # CI/CD 流水线
├── completions/           # Shell 补全脚本
├── Makefile              # 构建自动化
└── README.md
```

## 🧪 开发

### 环境要求

- Go 1.23 或更高版本
- Make（可选，用于使用 Makefile）

### 构建
```bash
# 为当前平台构建
make build

# 为所有平台构建
make build-all

# 运行测试
make test

# 运行代码检查
golangci-lint run
```

### 测试
```bash
# 运行所有测试
go test ./...

# 带覆盖率运行
go test -cover ./...

# 运行特定包的测试
go test ./pkg/config/
```

## 🤝 贡献

欢迎贡献！请阅读 [CONTRIBUTING.md](CONTRIBUTING.md) 了解指南。

### 开发流程

1. Fork 仓库
2. 创建功能分支
3. 进行修改
4. 添加测试
5. 运行测试和代码检查
6. 提交 Pull Request

## 📝 许可证

本项目采用 MIT 许可证 - 详见 [LICENSE](LICENSE) 文件。

## 🙏 致谢

- 使用 [Cobra](https://github.com/spf13/cobra) CLI 框架构建
- 感谢 OpenAI、Anthropic 和 Google 提供的 LLM API
- 感谢 Go 社区提供的优秀库

## 📊 统计

- **17 个命令**: 全面的文本处理工具包
- **3 个 LLM 提供商**: OpenAI、Claude、Gemini
- **5 种界面语言**: en、ja、zh、ko、th
- **13 种翻译语言**: 多语言支持
- **85%+ 测试覆盖率**: 经过充分测试的代码库
- **CI/CD**: 自动化构建和发布

## 🔗 链接

- [GitHub 仓库](https://github.com/hiroki-abe-58/aitxt)
- [问题跟踪](https://github.com/hiroki-abe-58/aitxt/issues)
- [发布版本](https://github.com/hiroki-abe-58/aitxt/releases)

## 💡 使用场景

- **内容创作**: 摘要文章、翻译内容、校对写作
- **开发**: 生成提交信息、审查代码、创建文档
- **学习**: 解释错误、提问、获取代码示例
- **自动化**: 批量处理、脚本编写、CI/CD 集成
- **研究**: 处理文档、提取信息、分析文本

---

Made with ❤️ by [hiroki-abe-58](https://github.com/hiroki-abe-58)
