# Decouple `specscore-cli` from `specscore`

**Date:** 2026-04-22
**Status:** Approved (brainstorm complete; awaiting implementation plan)
**Author:** Alexander Trakhimenok + Claude

## Context

Today the `synchestra-io/specscore` repository carries two concerns:

1. The **SpecScore format and documentation** — `spec/`, `docs/`, `blog/`, `examples/`, the `specscore.md` website (`public/`, `tools/site-generator/`, `firebase.json`), and a single project-def `specscore-spec-repo.yaml`.
2. The **reference CLI implementation** — `cmd/specscore/`, `internal/cli/`, `pkg/`, `go.mod`/`go.sum`, `.goreleaser.yml`, the Go CI and release workflows, and (recently added) the CLI's own spec at `spec/features/cli/`.

These two concerns have different audiences, different release cadences, different licenses-of-fit, and different contributor profiles. Bundling them means:

- The Go module ships unrelated documentation churn.
- The spec repo's release tags are CLI release tags, confusing for spec-only consumers.
- The CLI cannot be installed/imported under a name that signals "the reference implementation."
- A planned third repo (`specscore-ai-plugin`, a Claude Code skill bundle) will need its own license without inheriting the CLI's choice.

This document captures the decisions and migration plan to split the CLI into a dedicated public repository.

## Decisions

All decisions below were made through structured Q&A and locked before the design was written.

### D1 — CLI spec ownership: full move to `specscore-cli`

The entire `spec/features/cli/` tree moves to `specscore-cli/spec/features/cli/`. The spec repo no longer contains the CLI's contract.

**Rejected alternatives:**

- *Spec stays in `specscore`* — keeps the umbrella feature in the format repo but separates contract from code, creating two-PR drift risk for any flag change.
- *Hybrid: umbrella stays, children move* — breaks the SpecScore feature-tree invariant. `specscore feature show cli` against `specscore` would see an umbrella with no children; lint would either weaken or fail. This is a concrete tooling bug, not a stylistic concern.

The spec repo's `README.md` and `docs/ecosystem.md` will prominently link to `specscore-cli` as the reference implementation.

### D2 — Licensing per repo

| Repo | License | Rationale |
|---|---|---|
| `specscore` | **CC-BY-4.0** | Pure prose/spec/docs after the split. CC-BY is the convention for spec repositories (JSON Schema, OpenAPI, W3C drafts, MDN). Has a patent termination clause analogous to Apache-2.0. |
| `specscore-cli` | **Apache-2.0** | Matches the Go CLI ecosystem (kubectl, gh, terraform). Patent grant is harmless even if unneeded. No license change from current state. |
| `specscore-ai-plugin` | **MIT** | (Out of scope for this work.) Matches the AI/agent ecosystem (Claude skills, MCP servers, agent frameworks). Plugin will only shell out to the `specscore` binary, so the Apache→MIT direction-of-incompatibility doesn't bite. |

For the spec repo, code samples embedded in documentation are also covered by CC-BY-4.0 unless this becomes a practical problem; we are explicitly not adding a dual-license arrangement up front.

### D3 — Git history: full preservation via `git filter-repo`

The `specscore-cli` repo is created by running `git filter-repo` over a clone of `specscore`, keeping only paths destined for the CLI repo and rewriting them to their new locations. All commits authoring the CLI code, and all `v0.x` tags, come along.

**Rationale:**
- Engineering provenance (`git blame`) survives.
- `svu`-driven release versioning continues from the most recent tag rather than recomputing from the entire history.
- The cost is one extra command in the migration; the benefit is permanent.

### D4 — Dogfood CI install method: `get-cli` script, version-pinned

`specscore` repo's dogfood CI installs the CLI by piping the existing `get-cli` script (sourced from `https://specscore.md/get-cli`) with `SPECSCORE_VERSION` pinned to a specific release.

**Rationale:**
- Dogfooding the *shipped artifact* is more valuable than dogfooding source.
- Version pinning isolates spec-repo CI from CLI release regressions.
- Bumping the pin becomes a deliberate PR — the right amount of friction.

`specscore-cli` repo's own dogfood CI builds from source (the binary doesn't exist yet at build time).

### D5 — Pre-commit hook: out of scope

Initially considered in scope but explicitly removed. No `.githooks/`, no lefthook, no pre-commit framework added to either repo as part of this change.

### D6 — File movement principle

> Anything related to the CLI binary, its build, its release, its source, or its spec → moves to `specscore-cli`.
> Anything related to the spec format, spec docs, the website, or example consumers → stays in `specscore`.

Detailed mapping in **§ Final-state architecture** below.

### D7 — Installer single source of truth: `specscore-cli/scripts/install.sh`

The current duplication (`tools/site-generator/get-cli.sh` and `public/get-cli` are byte-identical copies) is a latent bug. After the split:

- The script's source-of-truth is `specscore-cli/scripts/install.sh`.
- The `specscore` site build fetches it from `specscore-cli` at build time and writes `public/get-cli`.
- The user-facing install URL (`https://specscore.md/get-cli`) is unchanged.

### D8 — Sequencing: phased cutover

Five-phase plan with verification gates between phases. Documented in **§ Phased migration plan**.

### D9 — Tag handling: delete `v0.x` tags from `specscore` after the split

Once `specscore-cli` carries the tags (Phase 1) and the CLI code is removed from `specscore` (Phase 4), the existing `v0.x` tags in `specscore` no longer point at meaningful code. They are deleted in Phase 5. A `README.md` note records the migration date.

A backup of the tag list is taken before deletion to permit rollback.

## Final-state architecture

### `synchestra-io/specscore`

**Purpose:** Canonical SpecScore format and documentation.
**License:** CC-BY-4.0.

```
specscore/
├── spec/                       # Format definitions (no features/cli/)
├── docs/                       # Narrative documentation
├── blog/                       # Blog posts
├── examples/todo-app/          # Example consumer of the format
├── public/                     # Static site (get-cli regenerated at build)
├── tools/site-generator/       # Site build tooling
├── firebase.json, .firebaserc  # Hosting config
├── specscore-spec-repo.yaml    # Project def for this repo
├── README.md                   # Format-focused; links to specscore-cli
├── LICENSE                     # CC-BY-4.0
└── .github/workflows/
    ├── site-ci.yml             # Existing
    └── dogfood.yml             # NEW: lints spec/ with pinned specscore binary
```

No `.go` files, no `go.mod`, no Go release pipeline.

### `synchestra-io/specscore-cli`

**Purpose:** Reference CLI implementation.
**License:** Apache-2.0.

```
specscore-cli/
├── cmd/specscore/main.go
├── internal/cli/               # Cobra commands
├── pkg/{exitcode,feature,gitremote,idea,lint,projectdef,sourceref,task}/
├── go.mod                      # module github.com/synchestra-io/specscore-cli
├── go.sum
├── .goreleaser.yml             # ldflags updated to new module path
├── spec/features/cli/          # The CLI's own contract (per D1)
├── scripts/install.sh          # Single source of truth (per D7)
├── README.md                   # CLI-focused; links to specscore for format
├── LICENSE                     # Apache-2.0
└── .github/workflows/
    ├── go-ci.yml               # Existing
    ├── release.yml             # Existing
    └── dogfood.yml             # NEW: builds CLI from source, lints own spec
```

Full git history; existing `v0.x` tags carried over.

### Cross-repo relationship

```
                    User runs:
        curl -fsSL https://specscore.md/get-cli | sh
                          │
                          ▼
        ┌─────────────────────────────────┐
        │  specscore (CC-BY-4.0)          │
        │  ─ specscore.md static site     │
        │  ─ public/get-cli  (built from  │
        │     specscore-cli/scripts/install.sh
        │  ─ dogfood CI uses specscore    │◄────┐
        │    binary (pinned version)      │     │ binary
        └────────────┬────────────────────┘     │ install
                     │ links to                 │
                     ▼                          │
        ┌─────────────────────────────────┐     │
        │  specscore-cli (Apache-2.0)     │     │
        │  ─ cmd/, internal/, pkg/        │─────┘
        │  ─ spec/features/cli/           │
        │  ─ scripts/install.sh (source)  │
        │  ─ goreleaser → GitHub Releases │
        │  ─ dogfoods own spec in CI      │
        └─────────────────────────────────┘
```

## Mechanics

### Module rename

`github.com/synchestra-io/specscore` → `github.com/synchestra-io/specscore-cli`

Affects:

| What | Change |
|---|---|
| `go.mod` `module` directive | `github.com/synchestra-io/specscore-cli` |
| Internal imports across `internal/cli/*.go` and `pkg/**/*.go` | Find/replace |
| `.goreleaser.yml` ldflags `-X` paths | `github.com/synchestra-io/specscore-cli/internal/cli.{version,commit,date}` |

Binary name remains `specscore`. Users see no change.

External-importer impact: no public Go-library promise was made for `pkg/`; the project is at v0.x; we accept the break and document it in the `specscore-cli` v-next release notes.

### Dogfood CI flow in `specscore`

`.github/workflows/dogfood.yml`:

```yaml
name: Dogfood lint
on:
  push:
    paths: ['spec/**', 'examples/**', '.github/workflows/dogfood.yml']
  pull_request:
    paths: ['spec/**', 'examples/**', '.github/workflows/dogfood.yml']
  workflow_dispatch:

env:
  SPECSCORE_VERSION: v0.6.0  # pinned — bump intentionally via PR

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - name: Install specscore CLI
        run: |
          curl -fsSL https://specscore.md/get-cli | \
            SPECSCORE_VERSION="${SPECSCORE_VERSION}" \
            SPECSCORE_INSTALL_DIR="$HOME/.local/bin" sh
          echo "$HOME/.local/bin" >> "$GITHUB_PATH"
      - name: Lint spec tree
        run: specscore spec lint
      - name: Lint example projects
        run: |
          for dir in examples/*/; do
            (cd "$dir" && specscore spec lint)
          done
```

### Dogfood CI flow in `specscore-cli`

```yaml
name: Dogfood lint
on:
  push:
    paths: ['spec/**', 'cmd/**', 'internal/**', 'pkg/**', '.github/workflows/dogfood.yml']
  pull_request:
    paths: ['spec/**', 'cmd/**', 'internal/**', 'pkg/**', '.github/workflows/dogfood.yml']

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/setup-go@v6
        with: { go-version-file: go.mod, cache: true }
      - name: Build CLI
        run: go build -o /tmp/specscore ./cmd/specscore
      - name: Lint own spec tree
        run: /tmp/specscore spec lint
```

### Installer fetch-at-build flow

The `specscore` site build fetches `install.sh` from `specscore-cli` at build time:

```js
// In tools/site-generator/build.js (or a pre-build step)
const SPECSCORE_CLI_REF = process.env.SPECSCORE_CLI_REF || 'main';
const url = `https://raw.githubusercontent.com/synchestra-io/specscore-cli/${SPECSCORE_CLI_REF}/scripts/install.sh`;
// fetch → write to public/get-cli → fail build on fetch error
```

`SPECSCORE_VERSION` is pinned in dogfood CI; `SPECSCORE_CLI_REF` defaults to `main` for the installer (the script changes rarely and is forward-compatible — it reads `SPECSCORE_VERSION` from env at install time). Override `SPECSCORE_CLI_REF` to pin during a known-bad period.

## Phased migration plan

### Phase 1 — Stand up `specscore-cli`

**Actions:**

1. Create `synchestra-io/specscore-cli` GitHub repo (empty, public, no init files).
2. In a temp clone of `specscore`, run `git filter-repo` keeping:
   - `cmd/`, `internal/`, `pkg/`, `go.mod`, `go.sum`, `.goreleaser.yml`
   - `spec/features/cli/`
   - `.github/workflows/go-ci.yml`, `.github/workflows/release.yml`
   - `tools/site-generator/get-cli.sh` → renamed to `scripts/install.sh`
3. Rewrite Go module path (`github.com/synchestra-io/specscore` → `github.com/synchestra-io/specscore-cli`) in `go.mod` and all `.go` files.
4. Update `.goreleaser.yml` ldflags `-X` paths.
5. Add `LICENSE` (Apache-2.0, content unchanged from current `specscore/LICENSE`).
6. Add `README.md` (CLI-focused, links back to `specscore`).
7. Push to `synchestra-io/specscore-cli` including all preserved tags.

**Verification gate:**

- `go build ./...` and `go test ./...` pass on a fresh clone.
- GitHub Actions Go CI passes on `main`.
- `git tag -l` shows the existing `v0.x` series.
- `git log spec/features/cli/README.md` shows the original commit `dce973a`.

**Rollback:** delete the GitHub repo and restart; `specscore` is untouched.

### Phase 2 — Cut first `specscore-cli` release

**Actions:**

1. Trigger `release.yml` (manual dispatch, `auto` bump → `svu` computes next version from conventional commits since last tag).
2. Confirm GoReleaser produces archives + checksums and uploads to GitHub Releases.

**Verification gate:**

- Release exists at `https://github.com/synchestra-io/specscore-cli/releases/tag/<v-next>`.
- Artifacts present for all platform/arch combinations.
- Manual install test: `curl -fsSL .../install.sh | SPECSCORE_VERSION=<v-next> sh` produces a working binary; `specscore --version` prints `<v-next>`.

**Rollback:** delete the release and tag; iterate.

### Phase 3 — Turn on dogfood CI in `specscore` (CLI code still present)

**Actions in `specscore`:**

1. Add `.github/workflows/dogfood.yml` with `SPECSCORE_VERSION` pinned to the Phase 2 release.
2. Push to a branch, open PR.

**Verification gate:**

- Dogfood workflow run is green: `specscore spec lint` passes against the spec tree (still including `features/cli/`) and against `examples/todo-app/spec/`.
- Existing Go CI and site CI also green.
- Merge.

**Why before deletion:** if dogfood lint fails, the CLI source is still locally available to debug.

### Phase 4 — Delete CLI from `specscore`, rewrite docs, switch license

**Actions in a single PR on `specscore`:**

1. Delete: `cmd/`, `internal/`, `pkg/`, `go.mod`, `go.sum`, `.goreleaser.yml`, `spec/features/cli/`, `.github/workflows/go-ci.yml`, `.github/workflows/release.yml`, `tools/site-generator/get-cli.sh`, `public/get-cli`.
2. Replace `LICENSE` with CC-BY-4.0 text. Update the license badge in `README.md`.
3. Rewrite `README.md`: format-focused; add prominent "Reference implementation" section linking to `specscore-cli`.
4. Update `docs/installation.md`: same install URL; text updated to reference `specscore-cli` releases.
5. Update `docs/ecosystem.md` and any other docs referencing the CLI by repo path or Go import path.
6. Update `tools/site-generator/build.js` to fetch `install.sh` from `specscore-cli` at build time.
7. Update the cross-references in `spec/features/source-references/README.md` and `_tests/*.md` that mention `spec/features/cli` (broken link → either remove or repoint to `specscore-cli`).
8. Sweep for any remaining `synchestra-io/specscore` references in `.github/`, `docs/`, `blog/`, `README.md` that should now be `synchestra-io/specscore-cli`.

**Verification gate:**

- `find . -name '*.go' -not -path './.git/*' | wc -l` returns 0.
- Site build (`site-ci.yml`) succeeds and produces `public/get-cli` from the fetched `install.sh`.
- Dogfood CI green against the slimmed spec tree.
- Manual smoke test: `https://specscore.md/get-cli | sh` still installs working `specscore`.
- Markdown links in `README.md` and `docs/` resolve.

**Rollback:** revert the PR; CLI code returns; dogfood CI keeps working because the released binary still exists.

### Phase 5 — Tag cleanup on `specscore`

**Actions:**

1. Take a backup: `git for-each-ref --format='%(refname) %(objectname)' refs/tags > /tmp/specscore-tags-backup.txt`.
2. Delete locally: `git tag | grep '^v' | xargs git tag -d`.
3. Delete remotely: `git push origin --delete <each tag>` (batched).
4. Add a one-paragraph note to `specscore/README.md` (e.g., "v0.x tags moved to [`specscore-cli`](https://github.com/synchestra-io/specscore-cli/tags) on YYYY-MM-DD").

**Verification gate:**

- `git ls-remote --tags origin` on `specscore` returns no `v*` tags.
- `git ls-remote --tags origin` on `specscore-cli` is unchanged from Phase 1.

**Rollback:** re-push tags from the backup.

## Risks & open items

**Known risks:**

- **External `pkg/` importers** — unlikely but unprobed. We document the module-path break in `specscore-cli` v-next release notes and accept it.
- **`specscore.md` DNS / Firebase config** — out of scope; the migration assumes site infra continues to function.
- **Conventional-commit history in `specscore`** — `feat(cli): …` commits remain in `specscore` history. `svu` no longer runs there post-split, so this is purely cosmetic.

**Out of scope (explicitly):**

- Pre-commit hook setup in either repo.
- `synchestra-io/specscore-ai-plugin` (separate work item).
- Renaming or restructuring `pkg/` packages.
- Reusable GitHub Action wrapper (`setup-specscore`) — defer until external consumers ask for it.
