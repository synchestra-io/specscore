# Feature: Source References

**Status:** Conceptual

## Summary

Source references are inline annotations in any source file that link code to SpecScore resources (features, plans, documents). A single prefix — `specscore:` — lets any tool, linter, or pre-commit hook discover references by binary search, resolve them against the project's spec repository, and transform them into clickable URLs pointing to `specscore.org`.

## Problem

Code that implements a feature often has no machine-readable link back to the specification that defines it. Developers add ad-hoc comments like `// see spec X` or `// Features implemented: cli/task/claim`, but these are convention-dependent, hard to validate, and invisible to spec-aware tooling.

Two concrete gaps exist:

1. **Discoverability** — spec-aware tools that scan dependency sections in spec files cannot tell you which *source files* implement or depend on a feature. Developers lose the spec-to-code traceability that makes specifications useful.
2. **Specification bridge** — embedding links to spec documentation in source code creates a natural bridge between code and specifications. Every developer who clicks a reference lands on the specification page, reinforcing the spec repository as the source of truth.

## Design Philosophy

- **Language-agnostic** — the notation must work in any language's comment syntax. Detection requires a recognized comment prefix on the same line — no AST parsing, just a single-line regex match.
- **Strict validation** — following Go's philosophy, references that point to non-existent resources are errors, not warnings. Invalid references are caught by linter, pre-commit hook, or PR check.
- **Single prefix** — `specscore:` covers all resource types (features, plans, docs). One prefix to search, one parser to maintain, one convention to learn.
- **Graceful cross-repo** — same-repo references omit host/org/repo for brevity. Cross-repo references append `@{host}/{org}/{repo}`. Host, org, and repo for the current context are inferred from git remote and can be overridden in project config.

## Behavior

### Notation format

```
specscore:{reference}
specscore:{reference}@{host}/{org}/{repo}
```

- **`{reference}`** — either a type-prefixed shortcut or a repo-root-relative path (see [Resolution](#short-notation-resolution))
- **`@{host}/{org}/{repo}`** — optional; omitted when referencing resources in the same project. `{host}` is the repository host (e.g., `github.com`, `bitbucket.org`, `gitlab.mycompany.com`)

### Resource type shortcuts

Known type prefixes provide shorthand for common paths. User-configurable types may be added later via project configuration.

| Type prefix | Expands to repo path | Example shortcut | Resolved path |
|---|---|---|---|
| `feature/` | `spec/features/{path}` | `feature/cli/task/claim` | `spec/features/cli/task/claim` |
| `plan/` | `spec/plans/{path}` | `plan/v2-migration` | `spec/plans/v2-migration` |
| `doc/` | `docs/{path}` | `doc/api/rest` | `docs/api/rest` |

### Short notation resolution

When resolving a `specscore:` reference, the following order is used:

1. **Type prefix** — if the first segment matches a known type prefix (`feature`, `plan`, `doc`), expand it to the corresponding repo path
2. **Fallback to path** — if the first segment is not a known prefix, or if type-based resolution fails (path does not exist), treat the entire value as a repo-root-relative path

**Examples:**

| Short notation | Resolution | Resolved repo path |
|---|---|---|
| `specscore:feature/cli/task/claim` | Type prefix `feature/` | `spec/features/cli/task/claim` |
| `specscore:plan/v2-migration` | Type prefix `plan/` | `spec/plans/v2-migration` |
| `specscore:doc/api/rest` | Type prefix `doc/` | `docs/api/rest` |
| `specscore:spec/features/cli/task/claim` | Not a known prefix — path | `spec/features/cli/task/claim` |
| `specscore:docs/api/rest` | Not a known prefix — path | `docs/api/rest` |
| `specscore:README.md` | Not a known prefix — path | `README.md` |

### Examples

What developers type (authoring — using type shortcuts or full paths):

```go
// specscore:feature/cli/task/claim
// specscore:feature/agent-skills@github.com/acme/orchestrator
// specscore:spec/features/cli/task/claim    (full path — equivalent to above)
```

What gets committed after lint/pre-commit expansion:

```go
// https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim
// https://specscore.org/github.com/acme/orchestrator/spec/features/agent-skills
```

```python
# https://specscore.org/github.com/acme/myproject/spec/features/model-selection
```

```yaml
# https://specscore.org/github.com/acme/myproject/spec/features/project-definition
```

```sql
-- https://specscore.org/bitbucket.org/acme/data-pipeline/spec/features/etl-config
```

### URL mapping

Every short reference expands to a canonical URL on `specscore.org`. The URL uses the **resolved repo-root-relative path** — the `{type}` prefix is not present in the URL.

```
specscore:{reference}
  → https://specscore.org/{host}/{org}/{repo}/{resolved_path}

specscore:{reference}@{host}/{org}/{repo}
  → https://specscore.org/{host}/{org}/{repo}/{resolved_path}
```

For same-repo references, `{host}/{org}/{repo}` is resolved at expansion time from git remote or project configuration.

**Examples:**

| Short reference | Expanded URL |
|---|---|
| `specscore:feature/cli/task/claim` | `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim` |
| `specscore:spec/features/cli/task/claim` | `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim` |
| `specscore:feature/agent-skills@github.com/acme/orchestrator` | `https://specscore.org/github.com/acme/orchestrator/spec/features/agent-skills` |
| `specscore:plan/v2-migration` | `https://specscore.org/github.com/acme/myproject/spec/plans/v2-migration` |
| `specscore:doc/api/rest@bitbucket.org/acme/docs` | `https://specscore.org/bitbucket.org/acme/docs/docs/api/rest` |
| `specscore:README.md` | `https://specscore.org/github.com/acme/myproject/README.md` |

### Canonical form and auto-expansion

The **expanded URL** is the canonical form stored in source files. The short `specscore:` notation is an **authoring convenience** — developers type the short form, and the linter (or pre-commit hook) auto-expands it to the full URL before commit.

**Rationale:** every `https://specscore.org/...` URL in a codebase is a clickable entry point. Developers can open the feature specification with one click — in any IDE, GitHub diff view, or `grep` output. No tooling is required to resolve the reference.

**Authoring workflow:**

1. Developer writes `specscore:feature/cli/task/claim` in a comment
2. Pre-commit hook (or spec-aware linter with `--fix`) resolves the type prefix and expands it to `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim`
3. The expanded URL is what gets committed and stored in the repository

The short form is never persisted in committed source — it exists only between authoring and the next lint/commit cycle.

### Detection strategy

A valid source reference must be preceded on the same line by a recognized comment prefix followed by optional whitespace. This eliminates false positives from string literals and non-comment code without requiring AST parsing.

**Detection regex (single line):**

```regex
^\s*(//|#|--|[/*]|%|;)\s*(specscore:|https://specscore\.org/)
```

**Recognized comment prefixes:**

| Prefix | Languages |
|---|---|
| `//` | Go, JS, TS, Java, C, C++, Rust, Swift, Kotlin |
| `#` | Python, Ruby, YAML, Shell, Perl, Elixir |
| `--` | SQL, Lua, Haskell |
| `*` or `/*` | Block comments in C-family languages |
| `%` | LaTeX, Erlang |
| `;` | Lisp, Clojure, INI files |

**Valid examples:**

```
// specscore:feature/cli/task/claim          (Go, JS)
//specscore:feature/cli/task/claim           (no space)
#  specscore:feature/model-selection         (Python, YAML)
-- https://specscore.org/github.com/org/repo/spec/features/x   (SQL)
; specscore:plan/v2-migration                (Lisp)
```

**Invalid examples (not detected):**

```
specscore:feature/cli/task/claim             (no comment prefix)
fmt.Println("specscore:feature/x")          (inside string literal)
var x = "https://specscore.org/github.com/org/repo/..." (inside string literal)
```

Users with uncommon comment syntax can open an issue to expand the prefix set, or override it in project configuration (future).

**Two reference forms are recognized:**

1. **Short notation** — `specscore:` prefix, then `{reference}[@{host}/{org}/{repo}]`
2. **Expanded URLs** — `https://specscore.org/` prefix, then `{host}/{org}/{repo}/{resolved_path}`

The linter auto-expands short notation to URLs, so committed code should only contain expanded URLs. The short form is accepted as input for authoring convenience.

### Host/org/repo resolution

When a reference omits `@{host}/{org}/{repo}`, the current project's host, org, and repo must be inferred:

1. **Git remote** — parse `origin` remote URL to extract `{host}`, `{org}`, and `{repo}`. This is the default. For example, `git@github.com:acme/myproject.git` yields `github.com/acme/myproject`.
2. **Project config override** — `specscore-project.yaml` may declare explicit values that override git remote inference. This handles forks, mirrors, and non-standard remote names.

```yaml
# specscore-project.yaml
project:
  host: github.com
  org: acme
  repo: myproject
```

### Validation

References are validated strictly — a reference to a non-existent resource is an error.

**Validation rules:**

| Check | Error condition |
|---|---|
| Reference resolves | The resolved repo path does not exist in the target repository (after trying type prefix expansion and path fallback) |
| Host/org/repo is resolvable | Same-repo reference but host/org/repo cannot be inferred (no git remote, no config override) |
| Cross-repo is reachable | `@{host}/{org}/{repo}` points to a repository that is not accessible (optional — may be deferred to CI) |

**Enforcement points:**

- **Linter** — a spec-aware linter scans source files, validates all references, reports errors with file:line locations
- **Pre-commit hook** — runs the linter on staged files before commit
- **PR check** — CI workflow that runs the linter on changed files

### Integration with spec-aware tools

Spec-aware tools can use source references as a second data source for dependency analysis, complementing the `## Dependencies` sections in spec files:

1. **Spec references** — features whose `## Dependencies` section lists the target
2. **Source references** — source files containing `specscore:feature/{target}` annotations

This enables bidirectional traceability: spec-to-spec via dependency sections, and code-to-spec via source references.

## Dependencies

- [feature](../feature/README.md)
- [project-definition](../project-definition/README.md)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Source references point to features; dependency analysis tools consume them |
| [Project Definition](../project-definition/README.md) | `specscore-project.yaml` provides org/repo override for resolution |
| [Development Plan](../development-plan/README.md) | Plans are a referenceable resource type |

## Acceptance Criteria

Not defined yet.

## Outstanding Questions

None at this time.
