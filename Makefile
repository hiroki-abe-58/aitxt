# aitxt Makefile

# Variables
BINARY_NAME=aitxt
VERSION=0.1.0
GIT_COMMIT=$(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags "-X github.com/hiroki-abe-58/aitxt/cmd.Version=$(VERSION) \
                  -X github.com/hiroki-abe-58/aitxt/cmd.GitCommit=$(GIT_COMMIT) \
                  -X github.com/hiroki-abe-58/aitxt/cmd.BuildDate=$(BUILD_DATE)"

# Default target
.PHONY: all
all: build

# Build binary
.PHONY: build
build:
	go build $(LDFLAGS) -o $(BINARY_NAME) main.go

# Build for multiple platforms
.PHONY: build-all
build-all:
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-arm64 main.go
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-darwin-amd64 main.go
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-amd64 main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BINARY_NAME)-linux-arm64 main.go
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BINARY_NAME)-windows-amd64.exe main.go

# Install locally
.PHONY: install
install: build
	mv $(BINARY_NAME) /usr/local/bin/

# Run tests
.PHONY: test
test:
	go test -v ./...

# Clean build artifacts
.PHONY: clean
clean:
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_NAME)-*

# Run the application
.PHONY: run
run:
	go run main.go

# Format code
.PHONY: fmt
fmt:
	go fmt ./...

# Lint code
.PHONY: lint
lint:
	golangci-lint run

# Show help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build      - Build binary for current platform"
	@echo "  build-all  - Build binaries for all platforms"
	@echo "  install    - Install binary to /usr/local/bin"
	@echo "  test       - Run tests"
	@echo "  clean      - Remove build artifacts"
	@echo "  run        - Run the application"
	@echo "  fmt        - Format code"
	@echo "  lint       - Lint code"
	@echo "  help       - Show this help"
