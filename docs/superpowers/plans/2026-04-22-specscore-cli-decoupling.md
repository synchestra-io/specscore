# specscore-cli Decoupling Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Split the SpecScore reference CLI out of `synchestra-io/specscore` into a new public Apache-2.0 repo `synchestra-io/specscore-cli`, leaving `synchestra-io/specscore` as a CC-BY-4.0 documentation/format repo with dogfood CI that lints itself using the released CLI binary.

**Architecture:** Five-phase migration with verification gates between phases. Phase 1 uses `git filter-repo` to clone the CLI's history into a new repo with the module path rewritten. Phase 2 cuts the first `specscore-cli` release. Phase 3 turns on dogfood CI in `specscore` *before* deleting the source. Phase 4 deletes the CLI from `specscore`, relicenses, rewrites docs, and switches the installer to fetch-at-build. Phase 5 cleans up old tags.

**Tech Stack:** Go 1.26, Cobra, GoReleaser, `svu` for semver bumps, `git filter-repo` for history surgery, GitHub Actions, Firebase Hosting, Node.js site generator.

**Source design doc:** `docs/superpowers/specs/2026-04-22-specscore-cli-decoupling-design.md`.

---

## Prerequisites

- [ ] **P0: Confirm tools installed**

Run:
```bash
git --version
go version
svu --version
gh --version
which git-filter-repo || echo MISSING
```

Expected: all print versions. If `git-filter-repo` prints `MISSING`, install it:
```bash
brew install git-filter-repo
# or: pip3 install --user git-filter-repo
```

Re-verify: `git filter-repo --version` prints a version string.

- [ ] **P1: Confirm GitHub CLI auth and org access**

Run:
```bash
gh auth status
gh api user/orgs --jq '.[].login' | grep synchestra-io
```

Expected: authenticated; `synchestra-io` listed in your orgs (or you have admin via membership).

- [ ] **P2: Use a worktree for `specscore`-side changes**

All Phase 3, 4, 5 work in `synchestra-io/specscore` should happen in a dedicated worktree to isolate from any other in-flight work:

```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git worktree add ../specscore.worktrees/cli-decoupling -b feat/cli-decoupling
cd ../specscore.worktrees/cli-decoupling
```

Phase 1 work happens in a *separate* temp directory (a fresh clone, not a worktree, because `filter-repo` rewrites history destructively).

---

## File Structure

### New repo: `synchestra-io/specscore-cli`

After Phase 1, the repo contains (paths relative to repo root):

```
specscore-cli/
├── cmd/specscore/main.go              # carried from specscore (history preserved)
├── internal/cli/                      # carried from specscore
│   ├── code.go
│   ├── feature.go
│   ├── new.go
│   ├── new_test.go
│   ├── root.go
│   ├── spec.go
│   └── task.go
├── pkg/                               # carried from specscore
│   ├── exitcode/
│   ├── feature/
│   ├── gitremote/
│   ├── idea/
│   ├── lint/
│   ├── projectdef/
│   ├── sourceref/
│   ├── task/
│   └── README.md
├── go.mod                             # module github.com/synchestra-io/specscore-cli
├── go.sum
├── .goreleaser.yml                    # ldflags rewritten to new module path
├── spec/features/cli/                 # CLI's own spec (carried)
├── scripts/install.sh                 # renamed from tools/site-generator/get-cli.sh
├── .github/workflows/
│   ├── go-ci.yml                      # carried
│   ├── release.yml                    # carried
│   └── dogfood.yml                    # NEW (Phase 1)
├── README.md                          # NEW (Phase 1) — CLI-focused
├── LICENSE                            # NEW (Phase 1) — Apache-2.0
└── .gitignore                         # NEW (Phase 1) — Go-focused
```

### Existing repo: `synchestra-io/specscore` after Phase 4

Files **deleted**:
- `cmd/`, `internal/`, `pkg/`
- `go.mod`, `go.sum`
- `.goreleaser.yml`
- `spec/features/cli/` (entire directory)
- `.github/workflows/go-ci.yml`
- `.github/workflows/release.yml`
- `tools/site-generator/get-cli.sh`
- `public/get-cli` (regenerated at site-build time)

Files **created**:
- `.github/workflows/dogfood.yml` (Phase 3)

Files **modified**:
- `LICENSE` (Apache-2.0 → CC-BY-4.0, Phase 4)
- `README.md` (rewrite to format-focused, Phase 4)
- `docs/installation.md` (Phase 4)
- `docs/ecosystem.md` (Phase 4)
- `tools/site-generator/build.js` (replace local `cp` of `get-cli.sh` with fetch from `specscore-cli`, Phase 4)
- `tools/site-generator/package.json` (add `node-fetch` if needed for fetch step)
- `spec/features/source-references/README.md` and `_tests/*.md` (cross-refs to `spec/features/cli`, Phase 4)

---

## Phase 1: Stand up `specscore-cli`

### Task 1.1: Create empty GitHub repo

**Files:** None (GitHub-side action)

- [ ] **Step 1: Create the repo**

Run:
```bash
gh repo create synchestra-io/specscore-cli \
  --public \
  --description "Reference CLI for SpecScore — lint, query, and scaffold SpecScore specifications." \
  --homepage "https://specscore.md"
```

Expected: repo URL printed, e.g. `https://github.com/synchestra-io/specscore-cli`.

- [ ] **Step 2: Verify repo is empty**

Run:
```bash
gh api repos/synchestra-io/specscore-cli --jq '{name, size, default_branch}'
```

Expected: `size: 0`, no default branch yet (fresh repo).

### Task 1.2: Clone `specscore` into a temp working directory

**Files:** None (clone operation)

- [ ] **Step 1: Fresh clone for filter-repo**

Run:
```bash
cd /tmp
rm -rf specscore-cli-migration
git clone https://github.com/synchestra-io/specscore.git specscore-cli-migration
cd specscore-cli-migration
```

Expected: clone completes; `git log --oneline | head` shows recent commits.

- [ ] **Step 2: Confirm working tree is clean**

Run:
```bash
git status
git tag -l 'v*' | head
```

Expected: clean working tree; existing `v0.x` tags listed.

### Task 1.3: Run `git filter-repo` to keep only CLI paths

**Files:** Rewrites all history in `/tmp/specscore-cli-migration`

- [ ] **Step 1: Run filter-repo**

Run (from `/tmp/specscore-cli-migration`):
```bash
git filter-repo \
  --path cmd/ \
  --path internal/ \
  --path pkg/ \
  --path go.mod \
  --path go.sum \
  --path .goreleaser.yml \
  --path spec/features/cli/ \
  --path .github/workflows/go-ci.yml \
  --path .github/workflows/release.yml \
  --path tools/site-generator/get-cli.sh \
  --path-rename tools/site-generator/get-cli.sh:scripts/install.sh
```

Expected: filter-repo runs; reports the number of commits processed and a "New history written" line.

- [ ] **Step 2: Verify the resulting tree**

Run:
```bash
ls -la
ls cmd/ internal/ pkg/ scripts/ spec/features/cli/ .github/workflows/
```

Expected:
- Top-level: `cmd/`, `internal/`, `pkg/`, `scripts/`, `spec/`, `.github/`, `go.mod`, `go.sum`, `.goreleaser.yml`
- No `docs/`, `blog/`, `public/`, `examples/`, `firebase.json`, `tools/`
- `scripts/install.sh` exists (renamed from old path)
- `spec/features/cli/` intact

- [ ] **Step 3: Verify tag preservation**

Run:
```bash
git tag -l 'v*'
git log --oneline spec/features/cli/README.md | head -3
```

Expected: same `v0.x` tags as in source repo; `spec/features/cli/README.md` log shows commit `dce973a` (or its rewritten SHA).

### Task 1.4: Rewrite Go module path

**Files:** Modify `/tmp/specscore-cli-migration/go.mod`

- [ ] **Step 1: Update module directive**

Run:
```bash
go mod edit -module github.com/synchestra-io/specscore-cli
```

- [ ] **Step 2: Verify**

Run:
```bash
head -3 go.mod
```

Expected:
```
module github.com/synchestra-io/specscore-cli

go 1.26.1
```

### Task 1.5: Rewrite import paths in `.go` files

**Files:** Modify all `.go` files under `/tmp/specscore-cli-migration`

- [ ] **Step 1: Run sed across all .go files**

Run (BSD/macOS sed-compatible form):
```bash
find . -name '*.go' -type f -exec sed -i.bak \
  's|github.com/synchestra-io/specscore|github.com/synchestra-io/specscore-cli|g' {} +
find . -name '*.go.bak' -type f -delete
```

- [ ] **Step 2: Verify no stale references**

Run:
```bash
grep -rn 'github.com/synchestra-io/specscore"' --include='*.go' .
grep -rn 'github.com/synchestra-io/specscore/' --include='*.go' . | grep -v 'specscore-cli'
```

Expected: both commands print nothing (no remaining bare `synchestra-io/specscore` references in Go code).

- [ ] **Step 3: Verify imports in a sample file**

Run:
```bash
grep -n 'synchestra-io' pkg/lint/idea.go
```

Expected: line 9 shows `"github.com/synchestra-io/specscore-cli/pkg/idea"`.

### Task 1.6: Update `.goreleaser.yml` ldflags

**Files:** Modify `/tmp/specscore-cli-migration/.goreleaser.yml`

- [ ] **Step 1: Replace ldflags paths**

Run:
```bash
sed -i.bak \
  's|github.com/synchestra-io/specscore/internal|github.com/synchestra-io/specscore-cli/internal|g' \
  .goreleaser.yml
rm .goreleaser.yml.bak
```

- [ ] **Step 2: Verify**

Run:
```bash
grep -n 'github.com/synchestra-io' .goreleaser.yml
```

Expected: 3 lines, all containing `github.com/synchestra-io/specscore-cli/internal/cli.{version,commit,date}`.

### Task 1.7: Update spec doc references in `spec/features/cli/`

**Files:** Modify markdown under `/tmp/specscore-cli-migration/spec/features/cli/`

The CLI's own spec contains references to the old module path and Synchestra Hub URLs that need rewriting.

- [ ] **Step 1: Update Hub URLs (project id changes from `specscore` to `specscore-cli`)**

Run:
```bash
find spec/features/cli -name '*.md' -type f -exec sed -i.bak \
  's|id=specscore@synchestra-io@github.com|id=specscore-cli@synchestra-io@github.com|g' {} +
find spec/features/cli -name '*.md.bak' -type f -delete
```

- [ ] **Step 2: Update bare module-path references in spec text**

Run:
```bash
find spec/features/cli -name '*.md' -type f -exec sed -i.bak \
  's|github.com/synchestra-io/specscore/internal|github.com/synchestra-io/specscore-cli/internal|g' {} +
find spec/features/cli -name '*.md.bak' -type f -delete
```

- [ ] **Step 3: Verify**

Run:
```bash
grep -rn 'id=specscore@synchestra-io' spec/features/cli/
grep -rn 'github.com/synchestra-io/specscore/internal' spec/features/cli/
```

Expected: both print nothing.

Run:
```bash
grep -n 'id=specscore-cli' spec/features/cli/README.md
```

Expected: line 3 shows the rewritten Hub URL.

### Task 1.8: Add Apache-2.0 LICENSE

**Files:** Create `/tmp/specscore-cli-migration/LICENSE`

- [ ] **Step 1: Copy the existing Apache-2.0 LICENSE from specscore source repo**

Run:
```bash
cp /Users/alexandertrakhimenok/projects/synchestra-io/specscore/LICENSE LICENSE
```

- [ ] **Step 2: Verify**

Run:
```bash
head -3 LICENSE
wc -l LICENSE
```

Expected: header `Apache License / Version 2.0, January 2004`; ~202 lines.

### Task 1.9: Add `README.md`

**Files:** Create `/tmp/specscore-cli-migration/README.md`

- [ ] **Step 1: Write a minimal CLI-focused README**

Create `/tmp/specscore-cli-migration/README.md` with:

````markdown
# specscore-cli

Reference CLI for [SpecScore](https://specscore.md) — lint, query, and scaffold SpecScore specifications.

## Install

```bash
curl -fsSL https://specscore.md/get-cli | sh
```

See [installation docs](https://specscore.md/installation) for options (version pinning, custom install dir).

## Usage

```bash
specscore spec lint              # lint the current spec tree
specscore feature list           # list features
specscore feature show <slug>    # inspect a feature
specscore task list              # show the task board
specscore version                # full build identity
specscore --version              # bare semver
```

Full command reference: see [`spec/features/cli/`](spec/features/cli/).

## License

Apache License 2.0 — see [LICENSE](LICENSE).

## Related

- [synchestra-io/specscore](https://github.com/synchestra-io/specscore) — the SpecScore format and documentation
````

- [ ] **Step 2: Verify**

Run:
```bash
head -5 README.md
```

Expected: shows the title and tagline.

### Task 1.10: Add `.gitignore`

**Files:** Create `/tmp/specscore-cli-migration/.gitignore`

- [ ] **Step 1: Write a Go-focused .gitignore**

Create `/tmp/specscore-cli-migration/.gitignore` with:

```
# Binaries
/specscore
/dist/

# Go workspace
go.work
go.work.sum

# IDE
.idea/
.vscode/

# OS
.DS_Store
```

### Task 1.11: Verify build & tests pass locally

**Files:** None (verification)

- [ ] **Step 1: Tidy modules and build**

Run (from `/tmp/specscore-cli-migration`):
```bash
go mod tidy
go build ./...
```

Expected: both succeed with no errors.

- [ ] **Step 2: Run tests**

Run:
```bash
go test ./...
```

Expected: all tests pass.

- [ ] **Step 3: Verify binary builds and runs**

Run:
```bash
go build -o /tmp/specscore-test ./cmd/specscore
/tmp/specscore-test --version
/tmp/specscore-test --help | head -10
rm /tmp/specscore-test
```

Expected: `--version` prints something (may be empty since ldflags weren't injected); `--help` prints usage.

### Task 1.12: Add dogfood workflow for specscore-cli

**Files:** Create `/tmp/specscore-cli-migration/.github/workflows/dogfood.yml`

- [ ] **Step 1: Write the workflow**

Create `/tmp/specscore-cli-migration/.github/workflows/dogfood.yml` with:

```yaml
name: Dogfood lint

on:
  push:
    paths:
      - 'spec/**'
      - 'cmd/**'
      - 'internal/**'
      - 'pkg/**'
      - '.github/workflows/dogfood.yml'
      - 'go.mod'
      - 'go.sum'
  pull_request:
    paths:
      - 'spec/**'
      - 'cmd/**'
      - 'internal/**'
      - 'pkg/**'
      - '.github/workflows/dogfood.yml'
      - 'go.mod'
      - 'go.sum'
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: dogfood-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v6
      - uses: actions/setup-go@v6
        with:
          go-version-file: go.mod
          cache: true
      - name: Build CLI
        run: go build -o /tmp/specscore ./cmd/specscore
      - name: Lint own spec tree
        run: /tmp/specscore spec lint
```

### Task 1.13: Initial commit, push, and tag verification

**Files:** None (git operations)

- [ ] **Step 1: Commit the new files**

Run:
```bash
git add LICENSE README.md .gitignore .github/workflows/dogfood.yml go.mod go.sum
git status  # confirm only these are added; no spurious deletions
git commit -m "$(cat <<'EOF'
chore: bootstrap specscore-cli repo

Initial bootstrap commit on top of filter-repo'd history. Adds:
- Apache-2.0 LICENSE
- CLI-focused README pointing back to specscore for the format
- Go-focused .gitignore
- Dogfood workflow that builds the CLI from source and lints its own spec
- go.mod/go.sum updated for new module path: github.com/synchestra-io/specscore-cli

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 2: Add origin and push**

Run:
```bash
git remote remove origin 2>/dev/null || true
git remote add origin https://github.com/synchestra-io/specscore-cli.git
git push -u origin main
git push origin --tags
```

Expected: push completes; tags pushed.

- [ ] **Step 3: Verify on GitHub**

Run:
```bash
gh api repos/synchestra-io/specscore-cli --jq '{default_branch, size}'
gh api repos/synchestra-io/specscore-cli/tags --jq '.[].name'
```

Expected: `default_branch: main`; tags include the existing `v0.x` series.

### Task 1.14: Verify GitHub Actions Go CI passes

**Files:** None (CI verification)

- [ ] **Step 1: Watch the workflow run**

Run:
```bash
gh run list --repo synchestra-io/specscore-cli --limit 5
gh run watch --repo synchestra-io/specscore-cli
```

Expected: the most recent workflow run completes successfully (Go CI + dogfood lint).

- [ ] **Step 2: If failures occur**

Investigate via:
```bash
gh run view --repo synchestra-io/specscore-cli --log-failed
```

Common issue: `spec lint` fails because of stale cross-references missed in Task 1.7 — fix in a follow-up commit, push, re-run.

---

## Phase 2: First `specscore-cli` release

### Task 2.1: Trigger release

**Files:** None (CI dispatch)

- [ ] **Step 1: Check what svu would compute**

Run (from `/tmp/specscore-cli-migration`):
```bash
git fetch --tags
svu next --v0
```

Expected: prints the next version, e.g. `v0.6.0` (depends on current tag and conventional commits since).

- [ ] **Step 2: Dispatch the release workflow**

Run:
```bash
gh workflow run release.yml \
  --repo synchestra-io/specscore-cli \
  --ref main \
  -f release_tag=auto
```

- [ ] **Step 3: Watch the run**

Run:
```bash
gh run watch --repo synchestra-io/specscore-cli
```

Expected: workflow completes successfully; new release created.

### Task 2.2: Verify release artifacts

**Files:** None (verification)

- [ ] **Step 1: List the release**

Run:
```bash
gh release list --repo synchestra-io/specscore-cli
gh release view --repo synchestra-io/specscore-cli
```

Expected: shows the new release tag and asset list.

- [ ] **Step 2: Confirm asset coverage**

Expected assets:
- `specscore_<v>_linux_amd64.tar.gz`
- `specscore_<v>_linux_arm64.tar.gz`
- `specscore_<v>_darwin_amd64.tar.gz`
- `specscore_<v>_darwin_arm64.tar.gz`
- `specscore_<v>_windows_amd64.zip`
- `specscore_<v>_checksums.txt`

If anything is missing, re-check `.goreleaser.yml` `goos`/`goarch` matrix — must match Phase 1's preserved config.

### Task 2.3: Manual install smoke test

**Files:** None (verification)

- [ ] **Step 1: Capture the new tag**

Run:
```bash
NEW_TAG=$(gh release view --repo synchestra-io/specscore-cli --json tagName --jq .tagName)
echo "Released tag: $NEW_TAG"
```

- [ ] **Step 2: Test install via the install script (still served from old location until Phase 4)**

Run (using the script from the new repo since the public URL still serves the old version until Phase 4 ships):
```bash
mkdir -p /tmp/install-test
cd /tmp/install-test
curl -fsSL "https://raw.githubusercontent.com/synchestra-io/specscore-cli/main/scripts/install.sh" | \
  SPECSCORE_VERSION="$NEW_TAG" SPECSCORE_INSTALL_DIR="$PWD" sh
```

Expected: downloads the matching archive, verifies checksum, installs `specscore` binary.

- [ ] **Step 3: Verify the binary**

Run:
```bash
./specscore --version
./specscore spec lint --help | head -5
rm -rf /tmp/install-test
```

Expected: `--version` prints the bare semver matching `$NEW_TAG`; `--help` prints usage.

- [ ] **Step 4: Note the version for Phase 3**

Record `$NEW_TAG` — Phase 3 pins `SPECSCORE_VERSION` to it.

---

## Phase 3: Turn on dogfood CI in `specscore` (CLI source still present)

All Phase 3 work happens in the worktree created in P2:
`/Users/alexandertrakhimenok/projects/synchestra-io/specscore.worktrees/cli-decoupling`

### Task 3.1: Add dogfood workflow

**Files:** Create `.github/workflows/dogfood.yml` in the worktree.

- [ ] **Step 1: Write the workflow**

Create `.github/workflows/dogfood.yml` with (replace `<v0.x.y>` with the tag from Task 2.3):

```yaml
name: Dogfood lint

on:
  push:
    paths:
      - 'spec/**'
      - 'examples/**'
      - '.github/workflows/dogfood.yml'
  pull_request:
    paths:
      - 'spec/**'
      - 'examples/**'
      - '.github/workflows/dogfood.yml'
  workflow_dispatch:

permissions:
  contents: read

concurrency:
  group: dogfood-${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

env:
  SPECSCORE_VERSION: <v0.x.y>  # bump intentionally via PR

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
            if [ -d "$dir/spec" ]; then
              (cd "$dir" && specscore spec lint)
            fi
          done
```

- [ ] **Step 2: Local sanity check**

Run (from the worktree):
```bash
specscore spec lint
for dir in examples/*/; do
  if [ -d "$dir/spec" ]; then
    (cd "$dir" && specscore spec lint)
  fi
done
```

Expected: both lint runs exit 0. If failures, fix the spec tree (or the lint rules) before pushing.

### Task 3.2: Commit, push, open PR

**Files:** None (git/PR ops)

- [ ] **Step 1: Commit**

Run:
```bash
git add .github/workflows/dogfood.yml
git commit -m "$(cat <<'EOF'
ci: add dogfood lint workflow using released specscore CLI

Installs the pinned specscore binary via get-cli script and runs
specscore spec lint over the spec/ tree and any examples/*/spec/
trees. Pin (SPECSCORE_VERSION) is bumped intentionally via PR to
isolate this repo from upstream CLI regressions.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 2: Push and open PR**

Run:
```bash
git push -u origin feat/cli-decoupling
gh pr create --title "ci: add dogfood lint workflow using released specscore CLI" --body "$(cat <<'EOF'
## Summary

- Adds .github/workflows/dogfood.yml that installs the pinned specscore binary and runs spec lint on this repo
- Pin: SPECSCORE_VERSION=<v0.x.y> (the first specscore-cli release, see https://github.com/synchestra-io/specscore-cli/releases)
- This is Phase 3 of the CLI decoupling per docs/superpowers/specs/2026-04-22-specscore-cli-decoupling-design.md

## Test plan

- [ ] Dogfood workflow run is green on this PR
- [ ] Existing site-ci and go-ci runs are still green

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### Task 3.3: Verify PR CI green

**Files:** None (verification)

- [ ] **Step 1: Watch PR checks**

Run:
```bash
gh pr checks --watch
```

Expected: dogfood workflow + existing checks all green.

- [ ] **Step 2: If dogfood fails**

Investigate:
```bash
gh run view --log-failed
```

Common failure modes:
- `get-cli` script can't reach GitHub releases → check the install URL works
- `specscore spec lint` reports violations against the current spec tree → fix the violations or update lint rules in `specscore-cli` and cut a new patch release

### Task 3.4: Merge PR

**Files:** None (git op)

- [ ] **Step 1: Merge**

Run:
```bash
gh pr merge --squash --delete-branch
```

- [ ] **Step 2: Clean up worktree**

Run:
```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git worktree remove ../specscore.worktrees/cli-decoupling
git pull origin main
```

---

## Phase 4: Delete CLI from `specscore`, relicense, rewrite docs

Re-create the worktree for the next batch of changes.

### Task 4.1: Set up worktree

**Files:** None (git op)

- [ ] **Step 1: Create worktree on a new branch**

Run:
```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git worktree add ../specscore.worktrees/cli-removal -b feat/cli-removal
cd ../specscore.worktrees/cli-removal
```

### Task 4.2: Delete Go code and build files

**Files:** Delete `cmd/`, `internal/`, `pkg/`, `go.mod`, `go.sum`, `.goreleaser.yml`

- [ ] **Step 1: Delete**

Run:
```bash
git rm -r cmd internal pkg
git rm go.mod go.sum
git rm .goreleaser.yml
```

- [ ] **Step 2: Verify no .go files remain**

Run:
```bash
find . -name '*.go' -not -path './.git/*' | wc -l
```

Expected: `0`.

### Task 4.3: Delete CLI spec directory

**Files:** Delete `spec/features/cli/`

- [ ] **Step 1: Delete**

Run:
```bash
git rm -r spec/features/cli
```

- [ ] **Step 2: Verify**

Run:
```bash
ls spec/features/ | grep -c '^cli$' || echo "OK: cli removed"
```

Expected: `OK: cli removed`.

### Task 4.4: Delete Go workflows

**Files:** Delete `.github/workflows/go-ci.yml` and `.github/workflows/release.yml`

- [ ] **Step 1: Delete**

Run:
```bash
git rm .github/workflows/go-ci.yml
git rm .github/workflows/release.yml
```

- [ ] **Step 2: Verify remaining workflows**

Run:
```bash
ls .github/workflows/
```

Expected: `dogfood.yml`, `site-ci.yml`, `README.md` (no go-ci.yml or release.yml).

### Task 4.5: Delete the installer source

**Files:** Delete `tools/site-generator/get-cli.sh`

The local `get-cli.sh` source is no longer canonical (lives in `specscore-cli/scripts/install.sh` now). The corresponding `public/get-cli` is a build artifact in a gitignored directory (`public/` is in `.gitignore`), so nothing needs to be done about it in git — it will simply be regenerated by Task 4.10's fetch logic on the next site build.

- [ ] **Step 1: Delete the source**

Run:
```bash
git rm tools/site-generator/get-cli.sh
```

- [ ] **Step 2: Verify public/ remains gitignored**

Run:
```bash
grep -n '^public/$' .gitignore
```

Expected: matches (no change needed; `public/get-cli` will not appear in `git status` after rebuild).

### Task 4.6: Replace LICENSE with CC-BY-4.0

**Files:** Modify `LICENSE`

- [ ] **Step 1: Fetch the canonical CC-BY-4.0 text**

Run:
```bash
curl -fsSL "https://creativecommons.org/licenses/by/4.0/legalcode.txt" -o LICENSE
```

- [ ] **Step 2: Verify**

Run:
```bash
head -5 LICENSE
```

Expected: starts with `Creative Commons Attribution 4.0 International` or similar canonical heading.

- [ ] **Step 3: Stage**

Run:
```bash
git add LICENSE
```

### Task 4.7: Rewrite `README.md`

**Files:** Modify `README.md`

- [ ] **Step 1: Read the existing README to understand what content to preserve**

Run:
```bash
cat README.md | head -100
```

- [ ] **Step 2: Rewrite to format-focused content**

Replace the entire `README.md` content. The new file should:
- Lead with what SpecScore *is* (the format)
- Have a prominent **Reference implementation** section linking to `synchestra-io/specscore-cli` with the install one-liner
- Link to `https://specscore.md` for full docs
- Update the license badge to CC-BY-4.0
- Remove any Go install / `go install` / `cmd/specscore` references
- Preserve high-level "what is SpecScore" framing from the existing README

A representative skeleton (the implementer should expand this against the existing README's tone and structure):

````markdown
# SpecScore

> A specification format for software products — features, requirements, scenarios, plans, tasks, decisions, and the source-code links that bind them together.

[![License: CC BY 4.0](https://img.shields.io/badge/License-CC_BY_4.0-lightgrey.svg)](https://creativecommons.org/licenses/by/4.0/)
[![Site](https://img.shields.io/badge/site-specscore.md-blue)](https://specscore.md)

## What is SpecScore?

SpecScore is a markdown-first format for writing software specifications that humans, AI agents, and tooling can all read. It defines a small set of artifact types (Feature, Requirement, Scenario, Plan, Task, Decision, Idea) with conventions for layout, cross-referencing, and source-code linkage.

[Read the full documentation →](https://specscore.md)

## Reference implementation

The reference CLI for working with SpecScore repositories is [`specscore-cli`](https://github.com/synchestra-io/specscore-cli):

```bash
curl -fsSL https://specscore.md/get-cli | sh
```

It validates spec trees, queries features, manages task boards, and scaffolds new artifacts. Source: <https://github.com/synchestra-io/specscore-cli>.

## Repository contents

- [`spec/`](spec/) — the SpecScore format itself, written as a SpecScore spec
- [`docs/`](docs/) — narrative documentation
- [`blog/`](blog/) — articles
- [`examples/`](examples/) — example consumer projects (e.g. [`todo-app`](examples/todo-app/))

## License

The contents of this repository are licensed under [CC BY 4.0](LICENSE).

The reference CLI [`specscore-cli`](https://github.com/synchestra-io/specscore-cli) is licensed separately under Apache-2.0.
````

(Adapt to match the existing README's actual voice/structure — preserve sections like project structure, contributing pointers, etc., minus anything CLI-specific.)

- [ ] **Step 3: Stage**

Run:
```bash
git add README.md
```

### Task 4.8: Update `docs/installation.md`

**Files:** Modify `docs/installation.md`

The existing doc references `synchestra-io/specscore` releases. Update to reference `synchestra-io/specscore-cli`.

- [ ] **Step 1: Update the GitHub release URL reference**

Run:
```bash
sed -i.bak \
  's|github.com/synchestra-io/specscore/releases|github.com/synchestra-io/specscore-cli/releases|g' \
  docs/installation.md
rm docs/installation.md.bak
```

- [ ] **Step 2: Verify**

Run:
```bash
grep -n 'synchestra-io/specscore' docs/installation.md
```

Expected: any remaining matches reference `specscore-cli` (no bare `specscore`). If bare matches remain that should *stay* as `specscore` (e.g., the format spec link), confirm those are intentional.

- [ ] **Step 3: Stage**

Run:
```bash
git add docs/installation.md
```

### Task 4.9: Update `docs/ecosystem.md` and other docs

**Files:** Modify `docs/ecosystem.md`, `docs/for/*.md`, possibly `docs/superpowers/`

- [ ] **Step 1: Find all references to the CLI repo path**

Run:
```bash
grep -rn 'github.com/synchestra-io/specscore' docs/ blog/ \
  --include='*.md' | grep -v 'specscore-cli'
```

- [ ] **Step 2: Review each match and decide intent**

For each line:
- If it refers to *the CLI repo* (e.g., release URL, source link, install instruction, Go import path) → rewrite to `specscore-cli`
- If it refers to *this spec repo itself* (e.g., "this repository is at github.com/synchestra-io/specscore") → leave as-is

- [ ] **Step 3: Apply targeted edits**

Use `Edit` (or `sed` on individual files) for each match identified in Step 2.

- [ ] **Step 4: Verify**

Run:
```bash
grep -rn 'github.com/synchestra-io/specscore' docs/ blog/ \
  --include='*.md' | grep -v 'specscore-cli'
```

Expected: only intentional self-references remain.

- [ ] **Step 5: Stage**

Run:
```bash
git add docs/ blog/
```

### Task 4.10: Update site generator to fetch installer at build time

**Files:** Modify `tools/site-generator/build.js`

Current behavior (around the comment `// CLI installer script — served at /get-cli`):
```js
await cp(
  join(__dirname, 'get-cli.sh'),
  join(OUTPUT, 'get-cli')
);
```

New behavior: fetch from `specscore-cli` at build time.

- [ ] **Step 1: Read current build.js**

Run:
```bash
grep -n 'get-cli' tools/site-generator/build.js
```

Note the line numbers of the existing `cp` block (typically around lines 31-34).

- [ ] **Step 2: Replace with fetch logic**

Replace the existing `cp` block with:

```js
// CLI installer script — served at /get-cli
// Source of truth: synchestra-io/specscore-cli/scripts/install.sh
{
  const cliRef = process.env.SPECSCORE_CLI_REF || 'main';
  const url = `https://raw.githubusercontent.com/synchestra-io/specscore-cli/${cliRef}/scripts/install.sh`;
  const res = await fetch(url);
  if (!res.ok) {
    throw new Error(`Failed to fetch install.sh from ${url}: ${res.status} ${res.statusText}`);
  }
  const installScript = await res.text();
  await writeFile(join(OUTPUT, 'get-cli'), installScript, { mode: 0o755 });
  console.log(`  get-cli (fetched from specscore-cli@${cliRef})`);
}
```

Use `Edit` to replace exactly the existing 4-line `await cp(...)` block plus its comment.

- [ ] **Step 3: Verify Node version supports fetch**

Run:
```bash
grep -E '"node"|"engines"' tools/site-generator/package.json
node --version
```

Native `fetch` requires Node 18+. If `package.json` declares an older minimum, bump it:
```bash
# in tools/site-generator/package.json, ensure "engines": { "node": ">=18" }
```

- [ ] **Step 4: Test the build locally**

Run:
```bash
cd tools/site-generator
pnpm install
pnpm run build  # or whatever the build script is — check package.json
ls -la ../../public/get-cli
head -5 ../../public/get-cli
```

Expected: `public/get-cli` exists and starts with `#!/bin/sh\n# SpecScore CLI installer`.

- [ ] **Step 5: Stage**

Run (from worktree root):
```bash
git add tools/site-generator/build.js tools/site-generator/package.json
```

Note: `public/` is already gitignored (verified in Task 4.5), so the regenerated `public/get-cli` artifact will not appear in `git status`. No `.gitignore` change needed.

### Task 4.11: Update `spec/features/source-references` cross-references

**Files:** Modify markdown files referencing `spec/features/cli`

- [ ] **Step 1: Find references**

Run:
```bash
grep -rn 'spec/features/cli' spec/ --include='*.md'
```

- [ ] **Step 2: For each match, decide:**

- If it's a path-style example used to *illustrate* source-references (not a real link), leave the path text as-is — it's documentation showing what a path *looks like*, not a hyperlink.
- If it's a hyperlink that resolves at render time, repoint to `https://github.com/synchestra-io/specscore-cli/tree/main/spec/features/cli` or remove the link.

- [ ] **Step 3: Apply edits**

Use `Edit` per file. Most likely no functional changes are needed (the references are illustrative path examples in lint test fixtures), but verify by reading each match in context.

- [ ] **Step 4: Stage**

Run:
```bash
git add spec/features/source-references/
```

### Task 4.12: Final sweep for stale `synchestra-io/specscore` references

**Files:** Possibly modify `.github/`, `tools/`, `examples/`, root files

- [ ] **Step 1: Sweep**

Run:
```bash
grep -rln 'github.com/synchestra-io/specscore' . \
  --include='*.md' --include='*.yml' --include='*.yaml' --include='*.json' \
  --include='*.js' --include='*.html' \
  --exclude-dir=node_modules --exclude-dir=.git \
  | grep -v 'specscore-cli'
```

- [ ] **Step 2: Triage each file**

For each file printed:
- If the reference *should* be `specscore-cli`, rewrite (use `Edit`)
- If the reference correctly points to the spec repo itself, leave it

- [ ] **Step 3: Final check**

Repeat the sweep — confirm only intentional self-references remain.

### Task 4.13: Local verification before pushing

**Files:** None (verification)

- [ ] **Step 1: Site build succeeds**

Run:
```bash
cd tools/site-generator
pnpm run build
cd ../..
ls -la public/get-cli
head -3 public/get-cli
```

Expected: `public/get-cli` exists, starts with `#!/bin/sh`.

- [ ] **Step 2: Dogfood lint passes locally**

Run (using locally installed CLI):
```bash
specscore spec lint
for dir in examples/*/; do
  if [ -d "$dir/spec" ]; then
    (cd "$dir" && specscore spec lint)
  fi
done
```

Expected: all lint runs exit 0.

- [ ] **Step 3: No .go files remain**

Run:
```bash
find . -name '*.go' -not -path './.git/*' -not -path './node_modules/*' | wc -l
```

Expected: `0`.

### Task 4.14: Commit, push, open PR

**Files:** None (git/PR op)

- [ ] **Step 1: Review staged changes**

Run:
```bash
git status
git diff --cached --stat | tail -20
```

Expected: large deletion (cmd/, internal/, pkg/, spec/features/cli/, .goreleaser.yml, etc.) + targeted modifications (LICENSE, README.md, docs, build.js).

- [ ] **Step 2: Commit**

Run:
```bash
git commit -m "$(cat <<'EOF'
feat!: extract CLI to synchestra-io/specscore-cli

BREAKING CHANGE: The specscore CLI source code has been moved to
https://github.com/synchestra-io/specscore-cli. This repository now
contains only the SpecScore format specification, documentation,
website, and example consumer projects.

- Delete cmd/, internal/, pkg/, go.mod, go.sum, .goreleaser.yml
- Delete spec/features/cli/ (moved with full history to specscore-cli)
- Delete .github/workflows/{go-ci,release}.yml (CLI build/release moved)
- Delete tools/site-generator/get-cli.sh and public/get-cli (regenerated
  at site build time by fetching from specscore-cli/scripts/install.sh)
- Relicense from Apache-2.0 to CC-BY-4.0 (spec/docs are pure prose)
- Rewrite README.md to format-focused content with prominent links to
  the reference CLI implementation
- Update docs/installation.md and docs/ecosystem.md to reference the
  new CLI repo for releases and source

Install URL https://specscore.md/get-cli is unchanged for users.

See docs/superpowers/specs/2026-04-22-specscore-cli-decoupling-design.md
for the full migration design.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
```

- [ ] **Step 3: Push and open PR**

Run:
```bash
git push -u origin feat/cli-removal
gh pr create --title "feat!: extract CLI to synchestra-io/specscore-cli" --body "$(cat <<'EOF'
## Summary

- Removes all CLI Go code and build infrastructure (now lives in [synchestra-io/specscore-cli](https://github.com/synchestra-io/specscore-cli))
- Removes spec/features/cli/ (moved with full history to specscore-cli)
- Relicenses from Apache-2.0 to CC-BY-4.0 (this repo is now pure documentation)
- Rewrites README.md to focus on the SpecScore format with prominent links to the reference CLI
- Switches site build to fetch the install script from specscore-cli at build time

This is Phase 4 of the CLI decoupling per `docs/superpowers/specs/2026-04-22-specscore-cli-decoupling-design.md`. Phases 1-3 already complete (specscore-cli repo created, first release cut, dogfood CI live).

## Test plan

- [ ] Site CI build succeeds and produces `public/get-cli`
- [ ] Dogfood lint workflow stays green against the slimmed spec tree
- [ ] Manual smoke test after merge: `curl -fsSL https://specscore.md/get-cli | sh` installs working specscore
- [ ] No `.go` files remain in the repo
- [ ] License badge in README reflects CC-BY-4.0

🤖 Generated with [Claude Code](https://claude.com/claude-code)
EOF
)"
```

### Task 4.15: Verify PR CI green

**Files:** None (CI verification)

- [ ] **Step 1: Watch checks**

Run:
```bash
gh pr checks --watch
```

Expected: site-ci and dogfood workflows both green.

- [ ] **Step 2: If failures**

Investigate and fix in additional commits on the same branch.

### Task 4.16: Merge and verify production

**Files:** None (verification)

- [ ] **Step 1: Merge**

Run:
```bash
gh pr merge --squash --delete-branch
```

- [ ] **Step 2: Wait for site deploy**

Run:
```bash
gh run list --workflow=site-ci.yml --limit 1
gh run watch
```

Expected: site-ci runs and deploys successfully on `main`.

- [ ] **Step 3: Smoke test public install URL**

Run:
```bash
mkdir -p /tmp/install-prod-test
cd /tmp/install-prod-test
curl -fsSL https://specscore.md/get-cli | \
  SPECSCORE_INSTALL_DIR="$PWD" sh
./specscore --version
rm -rf /tmp/install-prod-test
```

Expected: install succeeds; `--version` prints the latest released semver.

- [ ] **Step 4: Clean up worktree**

Run:
```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git worktree remove ../specscore.worktrees/cli-removal
git pull origin main
```

---

## Phase 5: Tag cleanup on `specscore`

### Task 5.1: Backup tag list

**Files:** Create `/tmp/specscore-tags-backup.txt`

- [ ] **Step 1: Snapshot current tags with their commit SHAs**

Run (from `/Users/alexandertrakhimenok/projects/synchestra-io/specscore`):
```bash
git fetch --tags
git for-each-ref --format='%(refname) %(objectname)' refs/tags > /tmp/specscore-tags-backup.txt
cat /tmp/specscore-tags-backup.txt
```

Expected: lists all `refs/tags/v*` with SHAs. **Save this file** (out of repo) in case rollback is needed.

### Task 5.2: Delete tags locally

**Files:** None (git op)

- [ ] **Step 1: Delete all v-prefixed tags locally**

Run:
```bash
git tag -l 'v*' | xargs git tag -d
```

- [ ] **Step 2: Verify**

Run:
```bash
git tag -l 'v*'
```

Expected: prints nothing.

### Task 5.3: Delete tags from remote

**Files:** None (git op)

- [ ] **Step 1: Get list of remote v-tags**

Run:
```bash
git ls-remote --tags origin | awk '{print $2}' | grep '^refs/tags/v' > /tmp/specscore-remote-tags.txt
cat /tmp/specscore-remote-tags.txt | wc -l
```

Expected: count matches the local tag count from Task 5.1.

- [ ] **Step 2: Push delete refspecs in batches**

Run:
```bash
xargs -n 20 git push origin --delete < /tmp/specscore-remote-tags.txt
```

(Or one-at-a-time if `xargs` batching is awkward:
```bash
while read -r ref; do
  git push origin --delete "$ref"
done < /tmp/specscore-remote-tags.txt
```
)

- [ ] **Step 3: Verify**

Run:
```bash
git ls-remote --tags origin | grep -c 'refs/tags/v' || echo "0 remote v-tags remain"
```

Expected: `0 remote v-tags remain`.

### Task 5.4: Add migration note to README

**Files:** Modify `README.md` in a worktree

- [ ] **Step 1: Set up worktree**

Run:
```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git worktree add ../specscore.worktrees/tag-note -b chore/tag-cleanup-note
cd ../specscore.worktrees/tag-note
```

- [ ] **Step 2: Add a one-paragraph note**

Append to `README.md` (or insert near the bottom under a new `## History` heading), replacing `<DATE>` with today's date in `YYYY-MM-DD` format:

```markdown
## History

The reference CLI `specscore` was extracted from this repository on <DATE>; its source code, releases, and `v0.x` release tags now live at [`synchestra-io/specscore-cli`](https://github.com/synchestra-io/specscore-cli). This repository's `v*` tags were removed at the same time, since they tagged CLI releases that no longer reside here. Engineering history (commits) for the extracted code is preserved in `specscore-cli` via `git filter-repo`.
```

Use `Edit` to add this section. Place it after the existing top-level sections.

- [ ] **Step 3: Commit, push, open PR**

Run:
```bash
git add README.md
git commit -m "$(cat <<'EOF'
docs: add migration note for v0.x tag removal

The v0.x tags previously on this repo tagged CLI releases that have
been moved to synchestra-io/specscore-cli. The tags were removed from
this repo on <DATE>. Engineering history is preserved in specscore-cli.

Co-Authored-By: Claude Opus 4.7 (1M context) <noreply@anthropic.com>
EOF
)"
git push -u origin chore/tag-cleanup-note
gh pr create --title "docs: add migration note for v0.x tag removal" --body "Adds a History section to README explaining where the v0.x tags went and why they were removed. Phase 5 wrap-up of the CLI decoupling."
```

- [ ] **Step 4: Merge**

Run:
```bash
gh pr checks --watch
gh pr merge --squash --delete-branch
```

### Task 5.5: Final verification

**Files:** None (verification)

- [ ] **Step 1: Confirm specscore has no Go code, no v-tags**

Run:
```bash
cd /Users/alexandertrakhimenok/projects/synchestra-io/specscore
git pull origin main
find . -name '*.go' -not -path './.git/*' -not -path './node_modules/*' | wc -l
git ls-remote --tags origin | grep -c 'refs/tags/v' || echo "0 v-tags"
```

Expected: `0` go files, `0 v-tags`.

- [ ] **Step 2: Confirm specscore-cli is healthy**

Run:
```bash
gh run list --repo synchestra-io/specscore-cli --limit 5
gh release list --repo synchestra-io/specscore-cli --limit 3
```

Expected: recent runs green; releases present.

- [ ] **Step 3: Confirm install URL still works**

Run:
```bash
mkdir -p /tmp/final-smoke
cd /tmp/final-smoke
curl -fsSL https://specscore.md/get-cli | \
  SPECSCORE_INSTALL_DIR="$PWD" sh
./specscore --version
rm -rf /tmp/final-smoke
```

Expected: install succeeds; binary works.

- [ ] **Step 4: Clean up worktree**

Run:
```bash
git worktree remove ../specscore.worktrees/tag-note
```

---

## Self-Review Checklist

Before declaring this plan complete, the implementer should confirm:

- [ ] All 5 phases' verification gates passed
- [ ] `synchestra-io/specscore`: 0 `.go` files, 0 `v*` tags, CC-BY-4.0 LICENSE
- [ ] `synchestra-io/specscore-cli`: builds, tests pass, dogfood CI green, has `v0.x` tag history
- [ ] Public install URL `https://specscore.md/get-cli` installs working CLI
- [ ] README in `specscore` prominently links to `specscore-cli`
- [ ] Tag backup file (`/tmp/specscore-tags-backup.txt`) preserved off-machine until confidence is high

## Rollback notes

- **Phase 1-2 only completed:** delete `synchestra-io/specscore-cli` repo via `gh repo delete`. `specscore` is untouched.
- **Phase 3 completed:** revert the dogfood-workflow PR; `specscore` returns to pre-Phase-3 state.
- **Phase 4 completed:** revert the Phase 4 PR; `specscore` returns to pre-Phase-4 state. The `specscore-cli` repo remains independent.
- **Phase 5 completed (tags deleted):** re-push tags from `/tmp/specscore-tags-backup.txt`:
  ```bash
  while read -r ref sha; do
    git push origin "$sha:$ref"
  done < /tmp/specscore-tags-backup.txt
  ```

## Out of scope

- Pre-commit hooks in either repo
- Renaming `synchestra-io/ai-plugin-specscore` → `synchestra-io/specscore-ai-plugin` (separate work; **do not create a new repo**)
- Reusable `setup-specscore` GitHub Action (defer until external consumers ask)
- `pkg/` package restructuring or stable-API promises
