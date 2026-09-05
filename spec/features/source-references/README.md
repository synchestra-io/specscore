---
format: https://specscore.md/feature-specification
status: Amending
---

# Feature: Source References

> [SpecScore.**Studio**](https://specscore.studio): | [Explore](https://specscore.studio/app/github.com/specscore/specscore/spec/features/source-references?op=explore) | [Edit](https://specscore.studio/app/github.com/specscore/specscore/spec/features/source-references?op=edit) | [Ask question](https://specscore.studio/app/github.com/specscore/specscore/spec/features/source-references?op=ask) | [Request change](https://specscore.studio/app/github.com/specscore/specscore/spec/features/source-references?op=request-change) |

**Status:** Amending
**Source Ideas:** spec-code-linking-via-codegraph

## Summary

Source references are inline annotations in any source file that link code to SpecScore resources. Untyped references record general context. Typed `implements` and `verifies` directives distinguish code that implements a requirement from executable tests that verify an acceptance criterion or requirement. Every committed target remains a clickable canonical `specscore.org` URL.

## Problem

Code that implements a feature often has no machine-readable link back to the specification that defines it. Developers add ad-hoc comments like `// see spec X` or `// Features implemented: cli/task/claim`, but these are convention-dependent, hard to validate, and invisible to spec-aware tooling.

Two concrete gaps exist:

1. **Discoverability** — spec-aware tools that scan dependency sections in spec files cannot tell you which *source files* implement or depend on a feature. Developers lose the spec-to-code traceability that makes specifications useful.
2. **Specification bridge** — embedding links to spec documentation in source code creates a natural bridge between code and specifications. Every developer who clicks a reference lands on the specification page, reinforcing the spec repository as the source of truth.

## Design Philosophy

- **Language-agnostic semantics** — the notation works in any language's comment syntax. SpecScore owns the reference and relation grammar. A code-intelligence provider may attach a directive to the nearest symbol using language-aware parsing.
- **Strict validation** — following Go's philosophy, references that point to non-existent resources are errors, not warnings. Invalid references are caught by linter, pre-commit hook, or PR check.
- **Single prefix** — `specscore:` covers all resource types (features, plans, docs). One prefix to search, one parser to maintain, one convention to learn.
- **Graceful cross-repo** — same-repo references omit host/org/repo for brevity. Cross-repo references use the URL authority form `specscore://{host}/{org}/{repo}/{reference}` ([decision 0010](../../decisions/0010-references-are-urls.md)). Host, org, and repo for the current context are inferred from git remote and can be overridden in `specscore.yaml`.

## Behavior

### Relation directives

A source comment MAY qualify a reference with one of three relation directives:

```
specscore:implements {target}
specscore:verifies {target}
specscore:references {target}
```

`{target}` accepts the same short or expanded address forms defined below. The canonical committed form preserves the directive and expands only its target:

```go
// specscore:implements https://specscore.org/github.com/acme/orders/spec/features/checkout#req:totals
func calculateTotal(...) { ... }

// specscore:verifies https://specscore.org/github.com/acme/orders/spec/features/checkout#ac:discounted-total
func TestDiscountedTotal(t *testing.T) { ... }
```

| Relation | Source | Allowed target | Meaning |
|---|---|---|---|
| `implements` | implementation symbol | REQ | The symbol participates in implementing the requirement |
| `verifies` | executable test symbol | AC, or a REQ that has no useful AC grouping | The test is executable evidence for the target |
| `references` | any source location or symbol | any SpecScore resource | The source depends on or provides context for the resource without claiming implementation or verification |

An existing unqualified source reference has `references` semantics. `implements` and `verifies` are accepted traceability claims and MUST be explicit in committed source. Inferred or agent-suggested pairs MUST NOT be reported as accepted claims until the corresponding directive is committed.

#### REQ: typed-relations

The relation set MUST be exactly `implements`, `verifies`, and `references`. An omitted relation MUST resolve to `references`. Unknown relation names MUST produce an error.

#### REQ: relation-targets

An `implements` directive MUST target a REQ. A `verifies` directive MUST target an AC or a directly verified REQ. A `references` directive MAY target any valid SpecScore resource.

#### REQ: verification-source

A `verifies` directive MUST attach to an executable test symbol. A scanner that can determine symbol kinds MUST reject a `verifies` directive attached to a non-test symbol. A language-agnostic fallback scanner MUST preserve the directive and report that symbol validation was not performed.

#### REQ: accepted-not-inferred

Tooling MUST distinguish accepted directives found in committed source from inferred candidate links. Suggested pairs MUST include their evidence and confidence and MUST NOT satisfy traceability or verification completeness checks.

### Notation format

```
specscore:{reference}
specscore://{host}/{org}/{repo}/{reference}
```

- **`{reference}`** — either a type-prefixed shortcut or a repo-root-relative path (see [Resolution](#short-notation-resolution))
- **`specscore://{host}/{org}/{repo}/…`** — the cross-repo form; omitted (single-colon opaque form) when referencing resources in the same project. `{host}` is the repository host (e.g., `github.com`, `bitbucket.org`, `gitlab.mycompany.com`); the first two path segments after the authority are `{org}/{repo}`, matching the canonical URL structure.
- **`?ref={git-ref}`** — optional on either form; pins a branch, tag, or commit ([decision 0010](../../decisions/0010-references-are-urls.md)). The fragment (`#`) remains reserved for heading anchors within the referenced document and MUST NOT carry version pins.

Both forms are valid URIs: the single-colon form is an opaque URI, and the `//` form honors RFC 3986 authority semantics, so standard URL libraries parse both. The former `specscore:{reference}@{host}/{org}/{repo}` suffix notation was removed by [decision 0010](../../decisions/0010-references-are-urls.md).

#### REQ: specscore-prefix

Every source reference MUST begin with the `specscore:` prefix followed by a reference string. No other prefix is permitted for SpecScore annotations.

#### REQ: cross-repo-authority

Cross-repo references MUST use the authority form `specscore://{host}/{org}/{repo}/{reference}`. Same-repo references MUST use the single-colon form without an authority. The legacy `@{host}/{org}/{repo}` suffix MUST be reported as an error carrying the exact authority-form rewrite, and lint `--fix` MUST apply it.

### Resource type shortcuts

Known type prefixes provide shorthand for common paths. User-configurable types may be added later via project configuration.

| Type prefix | Expands to repo path | Example shortcut | Resolved path |
|---|---|---|---|
| `feature/` | `spec/features/{path}` | `feature/cli/task/claim` | `spec/features/cli/task/claim` |
| `plan/` | `spec/plans/{path}` | `plan/v2-migration` | `spec/plans/v2-migration` |
| `doc/` | `docs/{path}` | `doc/api/rest` | `docs/api/rest` |

#### REQ: type-prefix-expansion

When a reference begins with a known type prefix (`feature/`, `plan/`, `doc/`), the resolver MUST expand it to the corresponding repo path according to the type prefix table. The type prefix itself MUST NOT appear in the resolved path.

#### REQ: type-prefix-set

The set of recognized type prefixes MUST be `feature/`, `plan/`, and `doc/`. Any other first segment MUST NOT be treated as a type prefix.

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

#### REQ: resolution-order

The resolver MUST first attempt type prefix expansion. If the first segment does not match a known type prefix, the resolver MUST treat the entire reference as a repo-root-relative path. If type prefix expansion produces a path that does not exist, the resolver MUST fall back to treating the full reference as a repo-root-relative path.

### URL mapping

Every short reference expands to a canonical URL on `specscore.org`. The URL uses the **resolved repo-root-relative path** — the `{type}` prefix is not present in the URL.

```
specscore:{reference}
  -> https://specscore.org/{host}/{org}/{repo}/{resolved_path}

specscore://{host}/{org}/{repo}/{reference}
  -> https://specscore.org/{host}/{org}/{repo}/{resolved_path}
```

The cross-repo scheme form is deliberately a pure prefix swap away from its
canonical expansion — `specscore://` ↔ `https://specscore.org/` — so the scheme
form, the canonical URL, and the Studio URL are three projections of one address.
A `?ref={git-ref}` pin, when present, is carried through expansion unchanged.

For same-repo references, `{host}/{org}/{repo}` is resolved at expansion time from git remote or project configuration.

**Examples:**

| Short reference | Expanded URL |
|---|---|
| `specscore:feature/cli/task/claim` | `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim` |
| `specscore:spec/features/cli/task/claim` | `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim` |
| `specscore://github.com/acme/orchestrator/feature/agent-skills` | `https://specscore.org/github.com/acme/orchestrator/spec/features/agent-skills` |
| `specscore:plan/v2-migration` | `https://specscore.org/github.com/acme/myproject/spec/plans/v2-migration` |
| `specscore://bitbucket.org/acme/docs/doc/api/rest` | `https://specscore.org/bitbucket.org/acme/docs/docs/api/rest` |
| `specscore:README.md` | `https://specscore.org/github.com/acme/myproject/README.md` |

#### REQ: url-structure

Expanded URLs MUST follow the pattern `https://specscore.org/{host}/{org}/{repo}/{resolved_path}`. The resolved path MUST NOT contain the type prefix — only the expanded repo-root-relative path.

### Canonical form and auto-expansion

The **expanded URL** is the canonical form stored in source files. The short `specscore:` notation is an **authoring convenience** — developers type the short form, and the linter (or pre-commit hook) auto-expands it to the full URL before commit.

**Rationale:** every `https://specscore.org/...` URL in a codebase is a clickable entry point. Developers can open the feature specification with one click — in any IDE, GitHub diff view, or `grep` output. No tooling is required to resolve the reference.

**Authoring workflow:**

1. Developer writes `specscore:feature/cli/task/claim` in a comment
2. Pre-commit hook (or spec-aware linter with `--fix`) resolves the type prefix and expands it to `https://specscore.org/github.com/acme/myproject/spec/features/cli/task/claim`
3. The expanded URL is what gets committed and stored in the repository

For a typed directive, the relation token remains unchanged and the target is expanded:

```
// specscore:verifies feature/checkout#ac:discounted-total
  -> // specscore:verifies https://specscore.org/github.com/acme/myproject/spec/features/checkout#ac:discounted-total
```

#### REQ: canonical-url-form

The canonical form of a source reference MUST be the fully expanded `https://specscore.org/...` URL. The short `specscore:` notation MUST NOT be persisted in committed source files.

#### REQ: auto-expansion

The linter or pre-commit hook MUST auto-expand short reference targets to the canonical URL form before commit. After expansion, an untyped `specscore:` reference SHOULD NOT remain in committed source.

For typed directives, the directive prefix remains `specscore:{relation}` and only its target is expanded. A canonical typed directive therefore contains both `specscore:{relation}` and a canonical `https://specscore.org/...` target.

### Detection strategy

A valid source reference or directive must be preceded on the same line by a recognized comment prefix followed by optional whitespace. A language-agnostic scanner can discover candidates with a line matcher. CodeGrapher or another code-intelligence provider SHOULD use syntax-aware parsing to attach the directive to the nearest symbol and validate whether it is implementation code or an executable test.

**Detection regex (single line):**

```regex
^\s*(//|#|--|[/*]|%|;)\s*(specscore:(implements|verifies|references)\s+|specscore:|https://specscore\.org/)
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

**Reference forms:**

1. **Short notation** — `specscore:` prefix: either the same-repo opaque form (`specscore:{reference}`) or the cross-repo authority form (`specscore://{host}/{org}/{repo}/{reference}`)
2. **Expanded URLs** — `https://specscore.org/` prefix, then `{host}/{org}/{repo}/{resolved_path}`
3. **Typed directives** — `specscore:{relation}` followed by either a short target or an expanded URL

The linter auto-expands short notation to URLs, so committed code should only contain expanded URLs. The short form is accepted as input for authoring convenience.

#### REQ: comment-prefix-required

A source reference MUST be preceded on the same line by a recognized comment prefix (`//`, `#`, `--`, `*`, `/*`, `%`, `;`) followed by optional whitespace. References not preceded by a comment prefix MUST NOT be detected or processed.

#### REQ: reference-forms-recognized

The detection strategy MUST recognize short notation (`specscore:` prefix), expanded URLs (`https://specscore.org/` prefix), and typed directives (`specscore:{relation} {target}`). Every form MUST match only when preceded by a comment prefix.

### Host/org/repo resolution

When a reference has no authority (the same-repo single-colon form), the current project's host, org, and repo must be inferred:

1. **Git remote** — parse `origin` remote URL to extract `{host}`, `{org}`, and `{repo}`. This is the default. For example, `git@github.com:acme/myproject.git` yields `github.com/acme/myproject`.
2. **Project config override** — `specscore.yaml` may declare explicit values that override git remote inference. This handles forks, mirrors, and non-standard remote names.

```yaml
# specscore.yaml
project:
  host: github.com
  org: acme
  repo: myproject
```

#### REQ: git-remote-default

When a reference carries no authority and no project config override exists, the resolver MUST infer host, org, and repo from the `origin` git remote URL.

#### REQ: config-override

When `specscore.yaml` declares `project.host`, `project.org`, and `project.repo`, those values MUST override git remote inference for same-repo reference expansion.

### Validation

References are validated strictly — a reference to a non-existent resource is an error.

**Validation rules:**

| Check | Error condition |
|---|---|
| Reference resolves | The resolved repo path does not exist in the target repository (after trying type prefix expansion and path fallback) |
| Host/org/repo is resolvable | Same-repo reference but host/org/repo cannot be inferred (no git remote, no config override) |
| Cross-repo is reachable | The `specscore://{host}/{org}/{repo}/…` authority points to a repository that is not accessible (optional — may be deferred to CI) |
| Legacy notation | A `specscore:{ref}@{host}/{org}/{repo}` suffix form is an error carrying the exact authority-form rewrite (`--fix` applies it) |

**Enforcement points:**

- **Linter** — a spec-aware linter scans source files, validates all references, reports errors with file:line locations
- **Pre-commit hook** — runs the linter on staged files before commit
- **PR check** — CI workflow that runs the linter on changed files

#### REQ: nonexistent-is-error

A reference that resolves to a path that does not exist in the target repository MUST produce an error. Non-existent references MUST NOT be treated as warnings.

#### REQ: unresolvable-context-error

When a same-repo reference cannot resolve host/org/repo (no git remote and no project config override), the linter MUST report an error. Expansion MUST NOT proceed with incomplete context.

#### REQ: cross-repo-reachability

When a cross-repo reference names a repository in its authority form, the validator SHOULD verify that the target repository is accessible. This check MAY be deferred to CI.

### Integration with spec-aware tools

Spec-aware tools can use source references as a second data source for dependency analysis, complementing the `## Dependencies` sections in spec files:

1. **Spec references** — features whose `## Dependencies` section lists the target
2. **Source references** — source files containing `specscore:feature/{target}` annotations

This enables bidirectional traceability: spec-to-spec via dependency sections, and code-to-spec via source references.

### Code-intelligence provider boundary

SpecScore owns artifact parsing, stable IDs, reference resolution, relation semantics, and validation rules. A code-intelligence provider owns language parsing, symbol identity, incremental indexing, and graph traversal. Integration MUST use a versioned CLI or service contract; SpecScore MUST NOT depend on a provider's private database format.

The provider SHOULD expose queries for accepted links, dangling targets, missing links, inferred candidate pairs, and exact-revision backlinks. Results MUST state the repository and revision they describe. SpecScore Studio MAY render these dynamic backlinks, but generated file and line locations MUST NOT be committed into the specification as if they were durable identifiers.

#### REQ: provider-boundary

SpecScore integrations with a code-intelligence provider MUST use a versioned public contract and MUST NOT read the provider's private storage format. Accepted and inferred links MUST remain distinguishable in every response.

#### REQ: exact-revision-results

Traceability results MUST identify the repository and exact Git revision they describe. Results from a different revision MUST be labelled stale and MUST NOT silently satisfy current-revision validation.

### Coverage evidence

Coverage is derived execution evidence. A coordinator such as WB MAY collect one language-native coverage artifact at a validation gate. The code-intelligence provider maps covered source ranges to symbols and follows accepted `implements` and `verifies` links to attribute evidence to REQs, ACs, and features.

Reports MUST keep these signals separate:

- **implementation exercised** means linked implementation code was executed by the measured suite;
- **AC or REQ verified** means an explicitly linked executable test passed.

Exercising implementation code is not sufficient evidence that an AC passed. Each coverage evidence record carries the exact commit, executed command, package or project scope, coverage mode, and collection time. Partial, missing, and stale results are visible. Aggregate suite coverage is the default; per-test coverage MAY be collected when its additional attribution value justifies its cost.

#### REQ: coverage-is-derived-evidence

Coverage of implementation code MUST NOT be reported as proof that an AC or REQ passed. Verification status requires an accepted `verifies` link to an executable test and a passing result for that test within the recorded run.

#### REQ: coverage-provenance

Every attributed coverage result MUST record the exact Git commit, command, scope, coverage mode, and collection time. Tooling MUST label evidence from another revision as stale and MUST expose incomplete scope.

#### REQ: coverage-collection-cost

Tooling SHOULD reuse one aggregate coverage artifact across all linked features and ACs. Per-test coverage MUST remain optional and MUST NOT be required on every commit.

## Dependencies

- [feature](../feature/README.md)
- [repo-config](../repo-config/README.md)

## Interaction with Other Features

| Feature | Interaction |
|---|---|
| [Feature](../feature/README.md) | Source references point to features; dependency analysis tools consume them |
| [Repo Config](../repo-config/README.md) | `specscore.yaml` provides `project.host`/`project.org`/`project.repo` overrides for git-remote inference |
| [Plan](../plan/README.md) | Plans are a referenceable resource type |
| [Requirement](../requirement/README.md) | Implementation symbols use `implements` links to REQs |
| [Acceptance Criteria](../acceptance-criteria/README.md) | Executable tests use `verifies` links to ACs or directly verified REQs |

## Acceptance Criteria

### AC: notation-and-resolution

**Requirements:** source-references#req:specscore-prefix, source-references#req:type-prefix-expansion, source-references#req:resolution-order

A `specscore:` reference with a known type prefix is expanded to the correct repo path. A reference without a known type prefix is treated as a repo-root-relative path. Both forms resolve to the same canonical location.

### AC: canonical-expansion

**Requirements:** source-references#req:canonical-url-form, source-references#req:auto-expansion, source-references#req:url-structure

Short `specscore:` notation is auto-expanded to a fully qualified `https://specscore.org/...` URL before commit. The expanded URL contains the resolved repo-root-relative path, not the type prefix. No short-form references remain in committed source.

### AC: detection-accuracy

**Requirements:** source-references#req:comment-prefix-required, source-references#req:reference-forms-recognized

References preceded by a recognized comment prefix are detected. References without a comment prefix (bare text, string literals) are ignored. Short notation, expanded URL, and typed directive forms are recognized.

### AC: strict-validation

**Requirements:** source-references#req:nonexistent-is-error, source-references#req:unresolvable-context-error, source-references#req:cross-repo-reachability

References to non-existent resources produce errors, not warnings. Missing host/org/repo context produces an error. Cross-repo reachability is optionally checked.

### AC: context-resolution

**Requirements:** source-references#req:git-remote-default, source-references#req:config-override, source-references#req:cross-repo-authority

Same-repo references resolve host/org/repo from git remote by default. Project config overrides git remote inference. Cross-repo references use the explicit `specscore://{host}/{org}/{repo}/{reference}` authority form; the legacy `@` suffix is rewritten by `--fix`.

### AC: typed-code-and-test-links

**Requirements:** source-references#req:typed-relations, source-references#req:relation-targets, source-references#req:verification-source, source-references#req:accepted-not-inferred

An implementation symbol can explicitly link to the REQ it implements, and an executable test can explicitly link to the AC or standalone REQ it verifies. Invalid relation targets fail validation. Suggested pairs remain separate from accepted committed links.

### AC: exact-revision-backlinks

**Requirements:** source-references#req:provider-boundary, source-references#req:exact-revision-results

A SpecScore client can query a code-intelligence provider through a versioned contract and receive accepted and suggested backlinks for an exact repository revision. Stale results are visible and cannot satisfy current-revision validation.

### AC: attributed-coverage

**Requirements:** source-references#req:coverage-is-derived-evidence, source-references#req:coverage-provenance, source-references#req:coverage-collection-cost

One aggregate coverage artifact can be attributed through linked symbols to features, REQs, and ACs with complete provenance. Reports distinguish exercised implementation from passing linked tests, expose partial or stale scope, and do not require per-test coverage on every commit.

## Open Questions

- Should the set of recognized comment prefixes be extensible via project configuration, or is the built-in set sufficient?
- How should the linter handle references in files with no recognized comment syntax (e.g., plain text files)?
- Should cross-repo reachability checks be mandatory in CI, or always optional?
- What is the behavior when a type prefix expansion fails but the literal path (with the type prefix) exists as a repo-root-relative path?
- Should canonical fragments migrate to lowercase `#req:` and `#ac:` while accepting uppercase `#REQ:` during a compatibility window?
- How should typed links attach to generated code that cannot carry comments?
- Which coverage formats should the first provider contract support beyond Go cover profiles?

---
*This document follows the https://specscore.md/feature-specification*
