# Contributing to SpecScore

SpecScore is an open-source specification framework. Contributions are welcome — whether that's improving the spec format, adding linter rules, fixing bugs, improving documentation, or writing guides.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<you>/specscore.git`
3. Install Go (1.22+): see [go.dev/dl](https://go.dev/dl)
4. Run tests: `go test ./...`
5. Run linter: `golangci-lint run ./...`

## What to Contribute

### Specification Format
The spec format lives in `spec/`. Changes to the format should include:
- Updated `README.md` for the affected feature
- Updated acceptance criteria
- Updated linter rules if applicable

### Linter Rules
Linter rules live in `pkg/lint/`. Each rule should:
- Have a clear, descriptive name
- Include tests in `pkg/lint/*_test.go`
- Be documented in the rule's error message

### Documentation
Guides live in `docs/`. Role-based guides are in `docs/for/`. Blog posts are in `blog/`.

## Code Standards

- Format with `gofmt`
- Pass `golangci-lint run ./...`
- Pass `go test ./...`
- Pass `go vet ./...`
- Handle all errors explicitly (see AGENTS.md)

## Pull Requests

- One feature per PR
- Include tests for new functionality
- Update documentation if behavior changes
- Reference the relevant spec feature if applicable

## License

By contributing, you agree that your contributions will be licensed under the Apache License 2.0.
