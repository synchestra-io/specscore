# Contributing to SpecScore

SpecScore is an open-source specification framework. Contributions are welcome — whether that's improving the spec format, fixing bugs, improving documentation, or writing guides.

This repository holds the **specification format** and the public website. The `specscore` CLI (and its linter rules) lives at [`synchestra-io/specscore-cli`](https://github.com/synchestra-io/specscore-cli) — see [Linter Rules](#linter-rules) below.

## Getting Started

1. Fork the repository
2. Clone your fork: `git clone https://github.com/<you>/specscore.git`
3. Install the `specscore` CLI: `curl -fsSL https://specscore.md/get-cli | sh`
4. Lint the spec tree: `specscore spec lint`

If you also intend to work on the website, install [pnpm](https://pnpm.io) and run `pnpm install && pnpm test` inside `tools/site-generator/`.

## What to Contribute

### Specification Format
The spec format lives in `spec/`. Changes to the format should include:
- Updated `README.md` for the affected feature
- Updated acceptance criteria
- A matching update or proposal for the linter rules in [`specscore-cli`](https://github.com/synchestra-io/specscore-cli) when the change introduces a new structural rule

### Linter Rules
Linter rules live in [`synchestra-io/specscore-cli`](https://github.com/synchestra-io/specscore-cli) under `pkg/lint/`. Open changes to rules — including new rules and tests in `pkg/lint/*_test.go` — against that repository, not this one.

### Documentation
Guides live in `docs/`. Role-based guides are in `docs/for/`. Blog posts are in `blog/`.

## Code Standards

- `specscore spec lint` MUST pass on `spec/` and on every `examples/*/spec/`.
- For website changes: `pnpm test` and `pnpm build` MUST pass inside `tools/site-generator/`.

## Pull Requests

- One feature per PR
- Update documentation if behavior changes
- Reference the relevant spec feature if applicable

## License

By contributing, you agree that your contributions will be licensed under [CC BY 4.0](LICENSE) (the same license that covers this repository's specification text and documentation). Linter-rule contributions to [`synchestra-io/specscore-cli`](https://github.com/synchestra-io/specscore-cli) are licensed separately under Apache-2.0.
