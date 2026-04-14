---
title: Linting Specs Like You Lint Code
description: How SpecScore's CLI catches ambiguity in your specifications the same way ESLint catches bugs in your JavaScript.
date: 2026-04-14
---

# Linting Specs Like You Lint Code

You lint your code. You lint your CSS. You lint your commit messages. Why aren't you linting your specifications?

Specifications are the upstream input to everything an AI agent builds. If the spec is ambiguous, the output is unpredictable. If the spec has missing fields, the agent fills in blanks with assumptions. If the spec structure is inconsistent, every agent interaction starts with "what format is this in?"

SpecScore brings the same discipline to specifications that ESLint brought to JavaScript.

## What the Linter Catches

The `specscore lint` command validates your specification files against the SpecScore schema:

```bash
$ specscore lint ./spec

spec/features/auth/README.md
  ✗ missing required field: status
  ✗ acceptance criteria AC3 has no testable condition

spec/features/search/README.md
  ✗ requirement R2 references undefined feature: indexing
  ⚠ no source references found for this feature

2 features, 3 errors, 1 warning
```

The linter checks for:

- **Structural completeness** — required fields, proper hierarchy, valid metadata
- **Referential integrity** — features, requirements, and acceptance criteria reference things that exist
- **Testability** — acceptance criteria contain verifiable conditions, not vague statements
- **Consistency** — naming conventions, status lifecycle, and format conformance

## Pluggable Rules

The linting engine is designed for extension. The core rules ship with SpecScore. Teams can add custom rules for their own conventions — naming patterns, required metadata fields, domain-specific validation.

## In Your Editor

The SpecScore LSP (coming soon) brings linting into your editor. Write a spec, see violations inline, fix them before you commit. The same feedback loop you have for code, applied to specifications.

## In CI/CD

Add `specscore lint` to your CI pipeline:

```yaml
# .github/workflows/specs.yml
- name: Lint specifications
  run: specscore lint ./spec
```

Specs that don't pass the linter don't get merged. The same gate you have for tests and code quality, applied to specifications.

## Try It

```bash
go install github.com/synchestra-io/specscore/cmd/specscore@latest
specscore lint ./spec
```

[Read the specification format](/specifications) to understand what the linter validates.
