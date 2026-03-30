# SpecScore — AI Agent Rules

## Build, test, and lint commands

- `go build ./...`
- `go vet ./...`
- `go test ./...`
- `golangci-lint run ./...`

These are also run by the CI workflow in `.github/workflows/go-ci.yml`.

## High-level architecture

SpecScore is both the specification format definition and its reference tooling:

- `spec/` is the technical source of truth for the SpecScore format (features, acceptance criteria, development plans, project definition, source references).
- `pkg/` contains the Go library packages — the public API importable by other tools (synchestra, rehearse, etc.).
- `internal/cli/` contains thin cobra command wrappers around `pkg/` functions.
- `cmd/specscore/` is the CLI entry point.
- `docs/` contains user-facing explanations.

Key packages:

- `pkg/exitcode` — shared exit code constants and error type
- `pkg/feature` — feature discovery, traversal, metadata, dependency resolution
- `pkg/lint` — specification linting engine with pluggable rules
- `pkg/sourceref` — `specscore:` annotation parsing and source-to-spec linking
- `pkg/projectdef` — `specscore-project.yaml` schema and read/write

## Directory structure

- Every directory MUST have a `README.md` file.
- Every `README.md` MUST have an "Outstanding Questions" section. If there are none, it explicitly states "None at this time."
- Every `README.md` that has child directories MUST include a brief summary for each immediate child.

## Go error handling requirements

Every function call that returns an error must have its error value handled explicitly:

- Capture error returns: `result, err := someFunction()`
- Check or explicitly ignore errors:
  - **Check**: `if err != nil { return err }`
  - **Explicitly ignore**: `_, _ = someFunction()`
- Do not silently drop error returns

## Go validation after code changes

After any change to `.go` files, agents must run the full Go validation sequence:

- `gofmt -w <changed-go-files>`
- `golangci-lint run ./...`
- `go test ./...`
- `go build ./...`
- `go vet ./...`
