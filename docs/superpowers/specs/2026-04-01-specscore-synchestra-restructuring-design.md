# SpecScore / Synchestra Package Restructuring

## Context

SpecScore and Synchestra have parallel implementations of the same core packages (feature discovery, spec linting, source references, exit codes, project config). Synchestra does not currently import SpecScore as a Go dependency — both projects independently implement overlapping functionality.

This design consolidates shared logic into SpecScore as the library layer, with Synchestra importing and extending it for coordination concerns.

## Boundary Principle

**SpecScore** = specification format, validation, navigation, read/create operations.
**Synchestra** = coordination, lifecycle state machine, sync, execution.

SpecScore is a library that knows nothing about Synchestra. Synchestra imports SpecScore and extends it with coordination capabilities. The `synchestra` CLI is the unified user-facing tool that wraps both layers.

## Migration Order

Migration follows the dependency graph (leaves first). Each step is independently testable. Changes push directly to main on both repos — no PRs during prototyping phase.

### Step 1: exitcode

Synchestra's `pkg/cli/exitcode` duplicates specscore's `pkg/exitcode`.

**Actions:**
- Reconcile any differences; ensure specscore's version is the superset
- Add specscore as a Go dependency in synchestra's `go.mod`
- Replace all `github.com/synchestra-io/synchestra/pkg/cli/exitcode` imports with `github.com/synchestra-io/specscore/pkg/exitcode`
- Delete synchestra's `pkg/cli/exitcode/`

### Step 2: sourceref

Synchestra's `pkg/sourceref` duplicates specscore's `pkg/sourceref`.

**Actions in specscore:**
- Add pluggable prefix registry to `pkg/sourceref`
- Default registered prefix: `specscore`
- Public API: `RegisterPrefix(prefix string)` adds a prefix to detection/parsing
- Detection regex and parsing functions check all registered prefixes

**Actions in synchestra:**
- Delete `pkg/sourceref/`
- Import `specscore/pkg/sourceref`
- Call `sourceref.RegisterPrefix("synchestra")` at CLI startup
- Update `cli/code/` to use specscore's sourceref

### Step 3: projectdef

Specscore's `pkg/projectdef` already defines `SpecConfig`, `CodeConfig`, `StateConfig`, and YAML I/O.

**Ownership split:**
- SpecScore owns: `SpecConfig`, `CodeConfig`, read/write functions, config file constants (renamed to `specscore-spec-repo.yaml`, `specscore-code-repo.yaml`)
- Synchestra owns: `StateConfig`, `EmbeddedStateConfig`, state-specific config files

**Config file structure (single file, namespaced extension):**

Config files are renamed from `synchestra-*` to `specscore-*` (no backward compatibility needed — prototyping phase).

```yaml
# specscore-spec-repo.yaml
title: "My Project"
repos:
  - github.com/org/code-repo

synchestra:
  state_repo: "github.com/org/state-repo"
  sync:
    pull: on_commit
    push: on_commit
```

SpecScore reads/writes `title`, `repos`, `planning`. Unknown YAML keys (including `synchestra:`) are round-tripped via `yaml:",inline"` on a `map[string]any` field. Synchestra reads the full file and manages the `synchestra:` section via its own Go types.

**Actions in specscore:**
- Remove `StateConfig` and `EmbeddedStateConfig` from `pkg/projectdef`
- Rename config file constants: `synchestra-spec-repo.yaml` → `specscore-spec-repo.yaml`, `synchestra-code-repo.yaml` → `specscore-code-repo.yaml`
- Ensure `SpecConfig` has an inline extras field for round-tripping unknown keys
- Keep `ReadSpecConfig` / `WriteSpecConfig` / `ReadCodeConfig` / `WriteCodeConfig`

**Actions in synchestra:**
- Move `StateConfig` and `EmbeddedStateConfig` to synchestra's own package
- Import `specscore/pkg/projectdef` for `SpecConfig` and `CodeConfig`
- Define synchestra-specific extension types that read from the `synchestra:` namespace in the spec config YAML
- Delete duplicated config I/O code

### Step 4: feature

Synchestra's `cli/feature/` (20 files) duplicates specscore's `pkg/feature`.

**Actions in specscore:**
- Reconcile: ensure `pkg/feature` is the superset of both implementations
- Merge any synchestra-only functions (if any) into specscore's package

**Actions in synchestra:**
- Delete `cli/feature/` implementation files (discover.go, tree.go, deps.go, refs.go, info.go, slug.go, transitive.go, fields.go, etc.)
- Replace with thin CLI wrappers that call `specscore/pkg/feature` functions
- Keep the cobra command registration in `cli/feature/feature.go` — it just delegates to specscore
- No feature status-change commands — those are not in scope for either project currently

### Step 5: lint

Synchestra's `cli/spec/` (12 files) duplicates specscore's `pkg/lint`.

**Actions in specscore:**
- Add pluggable checker registration to `pkg/lint`:
  ```go
  type Checker interface {
      Name() string
      Check(specRoot string) ([]Violation, error)
  }
  func RegisterChecker(c Checker)
  ```
- `Lint()` runs all built-in checkers plus any registered custom checkers
- Merge any synchestra-only rules that are spec-format concerns (not coordination-specific)

**Actions in synchestra:**
- Delete `cli/spec/` implementation files (linter.go, checkers_*.go)
- `spec lint` command becomes a thin wrapper: registers any coordination-specific custom checkers, then calls `lint.Lint(opts)`
- If no custom checkers exist, it's a pure pass-through

### Step 6: Task Types and Read Operations

This is the most complex step — splitting synchestra's `pkg/state` between the two projects.

**Moves to specscore (`pkg/task`):**
- Types: `Task` struct, `TaskStatus` enum and constants, `TaskCreateParams`, `TaskFilter`, `BoardView`, `BoardRow`
- Task file format: YAML frontmatter parsing/serialization (from `gitstore/taskfile.go`)
- Board format: markdown table rendering/parsing — format only, not locking mechanics (from `gitstore/board.go`)
- Read operations: `Get`, `List`, `Create` (create in `planning` status only)
- CLI commands: `task list`, `task info`, `task new`

**Stays in synchestra:**
- `Store` / `TaskStore` interfaces (with lifecycle methods: Claim, Start, Complete, Fail, Block, etc.)
- `gitstore` backend — coordination mechanics (optimistic locking, conflict detection, atomic commit-and-push)
- `ChatStore`, `StateSync`, `ProjectStore` interfaces and implementations
- `SyncConfig`, `SyncPolicy` types
- CLI lifecycle commands: `task enqueue/claim/start/status/complete/fail/block/unblock/release/abort/aborted`

**How synchestra uses specscore's task types:**
- Synchestra's `TaskStore` interface methods accept/return `specscore/pkg/task.Task`, `specscore/pkg/task.TaskStatus`, etc.
- Synchestra's `gitstore` imports `specscore/pkg/task` for file parsing and board rendering
- Synchestra embeds `specscore/pkg/task.Task` in its own `CoordinatedTask` type, adding coordination-only fields (Run, Model, ClaimedAt). The gitstore backend reads/writes the specscore base fields via `pkg/task` and manages coordination fields separately

## Post-Migration CLI Command Tree

```
synchestra
|-- project new/init/info/set       <- synchestra (orchestration) using specscore/pkg/projectdef
|-- project code add/remove         <- synchestra only
|-- feature list/tree/info/deps/refs/new  <- thin wrappers around specscore/pkg/feature
|-- spec lint                       <- wrapper around specscore/pkg/lint + custom checkers
|-- code deps                       <- wrapper around specscore/pkg/sourceref
|-- task new/list/info              <- wrappers around specscore/pkg/task
|-- task enqueue/claim/start/status <- synchestra (lifecycle via state.Store)
|-- task complete/fail/block/unblock <- synchestra (lifecycle via state.Store)
|-- task release/abort/aborted      <- synchestra (lifecycle via state.Store)
|-- state pull/push/sync            <- synchestra only
|-- config show/set/clear           <- synchestra only
+-- test run                        <- synchestra only
```

## What Stays in Synchestra (No Migration)

These packages are purely coordination/execution:

- `pkg/state/store.go` — Store, TaskStore, ChatStore, StateSync interfaces
- `pkg/state/gitstore/` — git-backed implementation (minus format parsing)
- `pkg/state/` — sync config, error types, chat types
- `cli/task/` — lifecycle commands (claim, start, complete, fail, block, unblock, release, abort)
- `cli/state/` — pull, push, sync commands
- `cli/gitops/` — git operations
- `cli/globalconfig/` — `~/.synchestra.yaml`
- `cli/reporef/`, `cli/resolve/` — repo resolution
- `cli/test/` — Rehearse integration
- `cli/project/` — multi-repo orchestration (imports specscore/pkg/projectdef for config I/O)

## Future Migration Path (Config Files)

If the `synchestra:` namespace in specscore's config file becomes unwieldy, migration to a dedicated `synchestra.yaml` companion file is straightforward:
1. Add `synchestra.yaml` reader alongside specscore config reader
2. If `synchestra.yaml` exists, prefer it over the inline namespace
3. Provide a `synchestra config migrate` command to extract the namespace to a dedicated file
