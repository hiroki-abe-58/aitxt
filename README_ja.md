# 🤖 aitxt - AI 駆動テキスト処理 CLI ツール

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Build Status](https://github.com/hiroki-abe-58/aitxt/actions/workflows/build.yml/badge.svg)](https://github.com/hiroki-abe-58/aitxt/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/hiroki-abe-58/aitxt)](https://goreportcard.com/report/github.com/hiroki-abe-58/aitxt)

[English](README.md) | 日本語 | [中文](README_zh.md)

複数の LLM プロバイダー（OpenAI、Claude、Gemini）をサポートする、AI 駆動のテキスト処理のための強力なコマンドラインツールです。

## ✨ 機能

### 🎯 コアコマンド
- **ask** - 一般的な AI 質問と支援
- **summarize** - 長さ制御可能なテキスト要約
- **translate** - 多言語翻訳（13言語）
- **proofread** - スタイルオプション付きの文法・スペルチェック
- **style** - テキストスタイルの変換（フォーマル、カジュアル、技術的など）
- **explain** - エラーメッセージの説明と解決策
- **review** - セキュリティ、パフォーマンス、スタイル分析を含むコードレビュー
- **doc** - コードからドキュメント生成

### 💼 開発者ツール
- **commit** - git diff から規約に準拠したコミットメッセージを生成
- **chat** - 履歴付きの対話モード
- **batch** - 並列実行による複数ファイル処理

### ⚙️ 生産性機能
- **template** - 変数をサポートする再利用可能なプロンプトテンプレート
- **alias** - よく使う操作のコマンドショートカット
- **history** - 会話履歴の検索と管理
- **config** - 設定と API キーの状態を表示

### 🌍 多言語サポート
- **UI 言語**: 英語、日本語、中国語、韓国語、タイ語
- **翻訳言語**: スペイン語、フランス語、ドイツ語、ポルトガル語、ロシア語、アラビア語、ヒンディー語、ベトナム語を含む13言語

### 🔧 技術的特徴
- **3つの LLM プロバイダー**: OpenAI (GPT-4o)、Claude (Sonnet 4)、Gemini (1.5 Flash)
- **並列処理**: 並行数設定可能なバッチ操作
- **出力フォーマット**: text、JSON、YAML
- **シェル補完**: bash、zsh、fish、powershell
- **進捗インジケーター**: スピナーとプログレスバー
- **ストリーミング出力**: リアルタイムレスポンス表示

## 📦 インストール

### リリース版からインストール（推奨）

お使いのプラットフォーム向けの最新リリースをダウンロード：
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

### ソースからインストール
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
make build
sudo make install
```

## 🚀 クイックスタート

### 1. API キーの設定
```bash
# お好みのプロバイダーを選択
export OPENAI_API_KEY="your-openai-key"
export ANTHROPIC_API_KEY="your-claude-key"
export GOOGLE_API_KEY="your-gemini-key"
```

### 2. インストールの確認
```bash
aitxt config
aitxt version
```

### 3. 基本コマンドを試す
```bash
# 質問する
aitxt ask "量子コンピューティングとは何ですか？"

# ファイルを要約
aitxt summarize document.txt

# テキストを翻訳
aitxt translate "こんにちは、世界！" --to en

# コミットメッセージを生成
git add .
aitxt commit
```

## 📚 使用例

### テキスト処理
```bash
# カスタム長さで要約
aitxt summarize article.txt --max-length 200

# 複数言語に翻訳
aitxt translate "おはようございます" --to en,zh,ko

# フォーマルスタイルで校正
aitxt proofread email.txt --style formal

# 文体を変換
aitxt style letter.txt --style professional
```

### 開発者ツール
```bash
# 特定の観点でコードレビュー
aitxt review main.go --focus security
aitxt review api.py --focus performance

# ドキュメント生成
aitxt doc src/ --type readme --format markdown

# エラーメッセージを説明
aitxt explain "TypeError: undefined is not a function"

# コミットメッセージ生成
git add .
aitxt commit --lang ja --type feat
```

### バッチ処理
```bash
# 複数ファイルを翻訳
aitxt batch *.md --operation translate --lang ja --concurrency 5

# すべての Python ファイルをレビュー
aitxt batch **/*.py --operation review --focus security

# ドキュメントを要約
aitxt batch docs/*.txt --operation summarize --output summaries/
```

### インタラクティブモード
```bash
# チャットセッションを開始
aitxt chat

# 特定のプロバイダーでチャット
aitxt chat --provider claude --model claude-sonnet-4-20250514

# 利用可能なチャットコマンド：
# /help      - ヘルプを表示
# /clear     - 履歴をクリア
# /system    - システムメッセージを変更
# /provider  - プロバイダーを切り替え
# /history   - 会話履歴を表示
# /save      - 会話を保存
# /exit      - チャットを終了
```

### テンプレート
```bash
# 組み込みテンプレートを初期化
aitxt template init

# テンプレートを一覧表示
aitxt template list

# テンプレートを使用
aitxt template use code-review-security --var Language=Go --var Code="$(cat main.go)"

# カスタムテンプレートを作成
aitxt template add --name my-review \
  --prompt "{{.Code}} を {{.Focus}} の観点でレビュー" \
  --system "あなたはコードレビュアーです"

# テンプレートをエクスポート/インポート
aitxt template export my-templates.json
aitxt template import shared-templates.json
```

### エイリアス
```bash
# 推奨エイリアスを初期化
aitxt alias init

# エイリアスを使用（初期化後）
aitxt s document.txt          # summarize
aitxt tj "Hello"               # 日本語に翻訳
aitxt rs main.go              # セキュリティレビュー

# カスタムエイリアスを作成
aitxt alias add fix "review --focus bugs"
aitxt fix broken-code.js

# すべてのエイリアスを一覧表示
aitxt alias list
```

### 履歴管理
```bash
# 最近の履歴を表示
aitxt history list --limit 10

# 履歴を検索
aitxt history search "翻訳"

# 特定のエントリを表示
aitxt history show <id>

# 履歴をエクスポート
aitxt history export history.json --format json

# 統計を表示
aitxt history stats
```

## 🔧 設定

### 環境変数
```bash
# API キー（必須）
export OPENAI_API_KEY="sk-..."
export ANTHROPIC_API_KEY="sk-ant-..."
export GOOGLE_API_KEY="..."

# オプション設定
export AITXT_PROVIDER="openai"        # デフォルトプロバイダー
export AITXT_MODEL="gpt-4o"           # デフォルトモデル
export AITXT_MAX_TOKENS="1000"        # トークン制限
export AITXT_TEMPERATURE="0.7"        # 創造性 (0.0-2.0)
export LANG="ja_JP.UTF-8"             # UI 言語
```

### シェル補完
```bash
# bash
aitxt completion bash > /etc/bash_completion.d/aitxt

# zsh
aitxt completion zsh > /usr/local/share/zsh/site-functions/_aitxt

# fish
aitxt completion fish > ~/.config/fish/completions/aitxt.fish
```

## 🏗️ プロジェクト構成
```
aitxt/
├── cmd/                    # CLI コマンド
│   ├── ask.go
│   ├── summarize.go
│   ├── translate.go
│   └── ...
├── pkg/                    # コアパッケージ
│   ├── llm/               # LLM プロバイダー実装
│   ├── config/            # 設定管理
│   ├── i18n/              # 国際化
│   ├── template/          # テンプレートシステム
│   ├── alias/             # エイリアス管理
│   ├── history/           # 履歴追跡
│   ├── batch/             # バッチ処理
│   ├── progress/          # 進捗インジケーター
│   └── output/            # 出力フォーマット
├── .github/workflows/     # CI/CD パイプライン
├── completions/           # シェル補完スクリプト
├── Makefile              # ビルド自動化
└── README.md
```

## 🧪 開発

### 必要要件

- Go 1.23 以上
- Make（オプション、Makefile 使用時）

### ビルド
```bash
# 現在のプラットフォーム向けにビルド
make build

# 全プラットフォーム向けにビルド
make build-all

# テスト実行
make test

# リンター実行
golangci-lint run
```

### テスト
```bash
# 全テスト実行
go test ./...

# カバレッジ付きで実行
go test -cover ./...

# 特定パッケージのテスト
go test ./pkg/config/
```

## 🤝 コントリビューション

コントリビューションを歓迎します！ガイドラインについては [CONTRIBUTING.md](CONTRIBUTING.md) をお読みください。

### 開発ワークフロー

1. リポジトリをフォーク
2. フィーチャーブランチを作成
3. 変更を実施
4. テストを追加
5. テストとリンターを実行
6. プルリクエストを提出

## 📝 ライセンス

このプロジェクトは MIT ライセンスの下で公開されています - 詳細は [LICENSE](LICENSE) ファイルをご覧ください。

## 🙏 謝辞

- [Cobra](https://github.com/spf13/cobra) CLI フレームワークを使用
- LLM API を提供する OpenAI、Anthropic、Google に感謝
- 優れたライブラリを提供する Go コミュニティに感謝

## 📊 統計

- **17 コマンド**: 包括的なテキスト処理ツールキット
- **3 LLM プロバイダー**: OpenAI、Claude、Gemini
- **5 UI 言語**: en、ja、zh、ko、th
- **13 翻訳言語**: 多言語サポート
- **85%+ テストカバレッジ**: 十分にテストされたコードベース
- **CI/CD**: 自動ビルドとリリース

## 🔗 リンク

- [GitHub リポジトリ](https://github.com/hiroki-abe-58/aitxt)
- [Issue トラッカー](https://github.com/hiroki-abe-58/aitxt/issues)
- [リリース](https://github.com/hiroki-abe-58/aitxt/releases)

## 💡 ユースケース

- **コンテンツ作成**: 記事の要約、コンテンツの翻訳、文章の校正
- **開発**: コミットメッセージ生成、コードレビュー、ドキュメント作成
- **学習**: エラー説明、質問、コード例の取得
- **自動化**: バッチ処理、スクリプト作成、CI/CD 統合
- **研究**: ドキュメント処理、情報抽出、テキスト分析

---

Made with ❤️ by [hiroki-abe-58](https://github.com/hiroki-abe-58)
