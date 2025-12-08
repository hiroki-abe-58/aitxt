# Contributing to aitxt

Thank you for your interest in contributing to aitxt!

## Development Setup

1. Clone the repository:
```bash
git clone https://github.com/hiroki-abe-58/aitxt.git
cd aitxt
```

2. Install dependencies:
```bash
go mod download
```

3. Build the project:
```bash
make build
```

## Running Tests
```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/config/
```

## Code Style

- Follow Go best practices and conventions
- Run `go fmt` before committing
- Run `go vet` to check for issues
- Add tests for new features

## Pull Request Process

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'feat: Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Commit Message Convention

We follow [Conventional Commits](https://www.conventionalcommits.org/):

- `feat:` - New feature
- `fix:` - Bug fix
- `docs:` - Documentation changes
- `test:` - Test additions/changes
- `chore:` - Maintenance tasks
- `ci:` - CI/CD changes

## License

By contributing, you agree that your contributions will be licensed under the MIT License.
