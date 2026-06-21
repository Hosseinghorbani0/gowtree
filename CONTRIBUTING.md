# Contributing to gowtree

Thank you for your interest in contributing to gowtree! Here's how you can help:

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/YOUR_USERNAME/gowtree.git`
3. Create a branch: `git checkout -b feature/your-feature`

## Building

```powershell
go build -o gowtree.exe ./cmd/gowtree
```

## Testing

```bash
go test ./...
```

## Submitting Changes

1. Commit your changes: `git commit -am 'Add feature'`
2. Push to your fork: `git push origin feature/your-feature`
3. Create a Pull Request with a clear description

## Code Style

- Follow Go conventions
- Run `go fmt` before committing
- Add tests for new features
- Update README if needed

## Reporting Issues

- Use GitHub Issues to report bugs
- Include OS, Go version, and steps to reproduce
- Attach output and error messages

## Ideas for Contribution

- Add more file type icons
- Write tests for `--files-first` and `--time` flags
- Add support for `.gitignore` patterns
- Create installers for macOS and Linux
- Translate documentation

Thank you!
