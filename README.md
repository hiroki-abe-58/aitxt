# 🎨 imgai

AI-powered image processing CLI tool for modern workflows.

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?logo=go)](https://go.dev/)

English | [日本語](README_ja.md)

A fast, efficient, and user-friendly command-line tool for image processing, built with Go and optimized for Apple Silicon.

## ✨ Features

### 🖼️ Image Processing
- **Resize** - Maintain aspect ratio or specify exact dimensions
- **Convert** - Transform between JPEG, PNG, and WebP formats
- **Quality Control** - Adjust compression for optimal file size

### 📊 Batch Operations
- **Parallel Processing** - Leverage goroutines for maximum performance
- **Progress Bar** - Visual feedback for batch operations
- **Glob Patterns** - Process multiple files with `*.jpg` patterns

### 🔒 Privacy & Metadata
- **EXIF Reading** - View camera settings, GPS, and metadata
- **EXIF Removal** - Strip all metadata for privacy protection

### 🛡️ Safety Features
- **Dry-Run Mode** - Preview operations before execution
- **Error Handling** - Detailed error messages and recovery
- **Automatic Naming** - Smart output filename generation

## 🚀 Installation

### Prerequisites
- Go 1.21 or higher
- macOS (Apple Silicon optimized) / Linux / Windows

### From Source
```bash
git clone https://github.com/hiroki-abe-58/imgai.git
cd imgai
go build -o imgai
```

### Quick Start
```bash
# Make it executable and move to PATH
chmod +x imgai
sudo mv imgai /usr/local/bin/

# Verify installation
imgai --version
```

## 🌐 Language Settings
```bash
# Use in Japanese
export IMGAI_LANG=ja
imgai resize photo.jpg --width 800

# Use in English (default)
export IMGAI_LANG=en
imgai resize photo.jpg --width 800
```

## 📖 Usage

### Resize Images
```bash
# Resize to 800px width (maintain aspect ratio)
imgai resize photo.jpg --width 800

# Resize to 600px height
imgai resize photo.jpg --height 600

# Resize to exact dimensions
imgai resize photo.jpg --width 1920 --height 1080

# Batch resize with progress bar
imgai resize *.jpg --width 800

# Preview before resizing
imgai resize *.jpg --width 800 --dry-run

# Use 8 parallel workers for faster processing
imgai resize *.jpg --width 800 --workers 8
```

### Convert Formats
```bash
# Convert to PNG
imgai convert photo.jpg --format png

# Convert to JPEG with custom quality
imgai convert image.png --format jpg --quality 85

# Convert to WebP (modern format)
imgai convert photo.jpg --format webp

# Batch convert all PNGs to JPEGs
imgai convert *.png --format jpg --quality 90

# Preview before converting
imgai convert *.png --format jpg --dry-run
```

### Manage Metadata
```bash
# View EXIF data
imgai exif photo.jpg

# Remove all metadata (privacy mode)
imgai strip photo.jpg

# Strip metadata and save to new file
imgai strip photo.jpg --output clean.jpg

# Batch strip metadata
imgai strip *.jpg --workers 8

# Preview metadata removal
imgai strip *.jpg --dry-run
```

## 🏗️ Architecture
```
imgai/
├── cmd/              # CLI commands
│   ├── root.go       # Root command with Cobra
│   ├── resize.go     # Resize command
│   ├── convert.go    # Format conversion
│   ├── exif.go       # EXIF reading
│   └── strip.go      # EXIF removal
├── pkg/              # Core packages
│   ├── image/        # Image processing logic
│   ├── batch/        # Batch processing with goroutines
│   ├── metadata/     # EXIF handling
│   └── i18n/         # Internationalization
└── main.go           # Entry point
```

## 🎯 Key Design Principles

- **DRY** - Don't Repeat Yourself
- **SOLID** - Object-oriented design principles
- **Separation of Concerns** - Clear module boundaries
- **Error Handling** - Comprehensive error messages
- **Cross-Platform** - Works on macOS, Linux, and Windows

## 🔧 Development

### Building
```bash
go build -o imgai
```

### Testing
```bash
go test ./...
```

### Adding New Features
```bash
# Create feature branch
git checkout -b feature/new-feature

# Make changes and commit
git add .
git commit -m "feat: add new feature"

# Push and create PR
git push origin feature/new-feature
```

## 📊 Performance

- **Parallel Processing** - Up to 8x faster with multiple workers
- **Memory Efficient** - Streaming image processing
- **Apple Silicon Optimized** - Native ARM64 support
- **Goroutines** - Concurrent processing for batch operations

## 🗺️ Roadmap

- [x] Project initialization
- [x] CLI framework (Cobra)
- [x] Image resize functionality
- [x] Format conversion (JPEG/PNG/WebP)
- [x] Batch processing with goroutines
- [x] EXIF metadata reading
- [x] Progress bar for batch operations
- [x] Dry-run mode
- [x] EXIF metadata removal
- [x] Japanese language support (i18n)
- [ ] AI-powered features (future)
- [ ] Cross-platform binaries (releases)

## 📝 License

MIT License - see [LICENSE](LICENSE) for details

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 👤 Author

**Hiroki Abe**
- GitHub: [@hiroki-abe-58](https://github.com/hiroki-abe-58)
- Repository: [imgai](https://github.com/hiroki-abe-58/imgai)

## 🙏 Acknowledgments

Built with:
- [Cobra](https://github.com/spf13/cobra) - CLI framework
- [imaging](https://github.com/disintegration/imaging) - Image processing
- [goexif](https://github.com/rwcarlsen/goexif) - EXIF handling
- [progressbar](https://github.com/schollz/progressbar) - Progress visualization

---

⭐ If you find this project useful, please consider giving it a star!
