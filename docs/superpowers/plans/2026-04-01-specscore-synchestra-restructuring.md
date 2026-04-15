# SpecScore / Synchestra Package Restructuring — Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Consolidate shared packages from synchestra into specscore as the library layer, with synchestra importing and wrapping them.

**Architecture:** Bottom-up migration following the dependency graph. Each task migrates one package: reconcile specscore's version to be the superset, update synchestra to import it, delete the duplicate. Both repos push directly to main after each task.

**Tech Stack:** Go 1.26.1, `github.com/synchestra-io/specscore`, `github.com/synchestra-io/synchestra`

**Repos:**
- SpecScore: `/home/ai/projects/synchestra-io/specscore/` (`github.com/synchestra-io/specscore`)
- Synchestra: `/home/ai/projects/synchestra-io/synchestra/` (`github.com/synchestra-io/synchestra`)

---

### Task 1: Migrate exitcode

The two `exitcode` packages are functionally identical. Synchestra has 54 files importing `synchestra/pkg/cli/exitcode`. SpecScore's `pkg/exitcode` has the same API. The only difference: synchestra's test file uses an external test package (`exitcode_test`) with more coverage; specscore's uses an internal test package with fewer tests.

**Files:**
- Modify: `specscore/pkg/exitcode/exitcode_test.go` (adopt synchestra's more comprehensive tests)
- Modify: `synchestra/go.mod` (add specscore dependency)
- Modify: 53 files in synchestra that import `synchestra/pkg/cli/exitcode` (change import path)
- Delete: `synchestra/pkg/cli/exitcode/exitcode.go`
- Delete: `synchestra/pkg/cli/exitcode/exitcode_test.go`

- [ ] **Step 1: Adopt synchestra's test coverage in specscore**

Replace `specscore/pkg/exitcode/exitcode_test.go` with synchestra's version, adjusted for the specscore module path:

```go
package exitcode_test

import (
	"errors"
	"testing"

	"github.com/synchestra-io/specscore/pkg/exitcode"
)

func TestErrorSatisfiesInterface(t *testing.T) {
	type exitCoder interface{ ExitCode() int }
	var err error = exitcode.New(1, "test")
	var ec exitCoder
	if !errors.As(err, &ec) {
		t.Fatal("exitcode.Error does not satisfy exitCoder interface")
	}
	if ec.ExitCode() != 1 {
		t.Fatalf("expected exit code 1, got %d", ec.ExitCode())
	}
}

func TestErrorMessage(t *testing.T) {
	err := exitcode.New(2, "bad args")
	if err.Error() != "bad args" {
		t.Fatalf("expected 'bad args', got %q", err.Error())
	}
}

func TestNewf(t *testing.T) {
	err := exitcode.Newf(10, "failed: %s", "disk full")
	if err.Error() != "failed: disk full" {
		t.Fatalf("expected 'failed: disk full', got %q", err.Error())
	}
	if err.ExitCode() != 10 {
		t.Fatalf("expected exit code 10, got %d", err.ExitCode())
	}
}

func TestConvenienceConstructors(t *testing.T) {
	tests := []struct {
		name string
		err  *exitcode.Error
		code int
	}{
		{"Conflict", exitcode.ConflictError("c"), exitcode.Conflict},
		{"ConflictF", exitcode.ConflictErrorf("c %d", 1), exitcode.Conflict},
		{"InvalidArgs", exitcode.InvalidArgsError("a"), exitcode.InvalidArgs},
		{"InvalidArgsF", exitcode.InvalidArgsErrorf("a %d", 2), exitcode.InvalidArgs},
		{"NotFound", exitcode.NotFoundError("n"), exitcode.NotFound},
		{"NotFoundF", exitcode.NotFoundErrorf("n %d", 3), exitcode.NotFound},
		{"InvalidState", exitcode.InvalidStateError("s"), exitcode.InvalidState},
		{"InvalidStateF", exitcode.InvalidStateErrorf("s %d", 4), exitcode.InvalidState},
		{"Unexpected", exitcode.UnexpectedError("u"), exitcode.Unexpected},
		{"UnexpectedF", exitcode.UnexpectedErrorf("u %d", 10), exitcode.Unexpected},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.ExitCode() != tt.code {
				t.Errorf("expected code %d, got %d", tt.code, tt.err.ExitCode())
			}
		})
	}
}

func TestConstants(t *testing.T) {
	if exitcode.Success != 0 {
		t.Error("Success should be 0")
	}
	if exitcode.Conflict != 1 {
		t.Error("Conflict should be 1")
	}
	if exitcode.InvalidArgs != 2 {
		t.Error("InvalidArgs should be 2")
	}
	if exitcode.NotFound != 3 {
		t.Error("NotFound should be 3")
	}
	if exitcode.InvalidState != 4 {
		t.Error("InvalidState should be 4")
	}
	if exitcode.Unexpected != 10 {
		t.Error("Unexpected should be 10")
	}
}
```

- [ ] **Step 2: Run specscore tests to verify**

```bash
cd /home/ai/projects/synchestra-io/specscore && go test ./pkg/exitcode/... -v
```

Expected: All 5 tests PASS.

- [ ] **Step 3: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add pkg/exitcode/exitcode_test.go
git commit -m "exitcode: adopt comprehensive test coverage from synchestra"
git push origin main
```

- [ ] **Step 4: Add specscore dependency to synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
```

- [ ] **Step 5: Replace all exitcode imports in synchestra**

Use `sed` or manual replacement to change all 53 import occurrences:

Old: `"github.com/synchestra-io/synchestra/pkg/cli/exitcode"`
New: `"github.com/synchestra-io/specscore/pkg/exitcode"`

Files to update (all in `synchestra/pkg/cli/`):
- `code/deps.go`
- `feature/deps.go`, `feature/discover.go`, `feature/fields.go`, `feature/info.go`, `feature/list.go`, `feature/new.go`, `feature/refs.go`, `feature/tree.go`
- `feature/flags_integration_test.go`, `feature/new_test.go`
- `project/init.go`, `project/init_test.go`, `project/new.go`, `project/new_test.go`
- `resolve/resolve.go`, `resolve/resolve_test.go`
- `spec/lint.go`
- `state/pull.go`, `state/pull_test.go`, `state/push.go`, `state/push_test.go`, `state/sync.go`, `state/sync_test.go`
- `task/abort.go`, `task/abort_test.go`, `task/aborted.go`, `task/aborted_test.go`
- `task/block.go`, `task/block_test.go`, `task/claim.go`, `task/claim_test.go`
- `task/complete.go`, `task/complete_test.go`, `task/enqueue.go`, `task/enqueue_test.go`
- `task/fail.go`, `task/fail_test.go`, `task/info.go`, `task/info_test.go`
- `task/list_test.go`, `task/new.go`, `task/new_test.go`
- `task/release.go`, `task/release_test.go`, `task/resolve.go`
- `task/start.go`, `task/start_test.go`, `task/status.go`, `task/status_test.go`
- `task/unblock.go`, `task/unblock_test.go`

- [ ] **Step 6: Delete synchestra's exitcode package**

```bash
rm -rf /home/ai/projects/synchestra-io/synchestra/pkg/cli/exitcode/
```

- [ ] **Step 7: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go fmt ./...
go build ./...
go test ./... -count=1
```

Expected: All build and tests PASS.

- [ ] **Step 8: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: replace internal exitcode with specscore/pkg/exitcode"
git push origin main
```

---

### Task 2: Migrate sourceref (with prefix registry)

Synchestra's `pkg/sourceref` is imported by only 1 file (`cli/code/deps.go`). SpecScore's version is nearly identical but needs a pluggable prefix registry. Currently both hardcode `synchestra:` as the prefix. After migration, specscore supports registering multiple prefixes; `specscore:` is default, synchestra registers `synchestra:` at startup.

**Files:**
- Modify: `specscore/pkg/sourceref/sourceref.go` (add prefix registry, update detection/parsing)
- Modify: `specscore/pkg/sourceref/sourceref_test.go` (adopt synchestra's comprehensive tests + prefix registry tests)
- Modify: `synchestra/pkg/cli/code/deps.go` (change import path)
- Modify: `synchestra/pkg/cli/main.go` (register `synchestra:` prefix at startup)
- Delete: `synchestra/pkg/sourceref/sourceref.go`
- Delete: `synchestra/pkg/sourceref/sourceref_test.go`
- Delete: `synchestra/pkg/sourceref/scan.go`

- [ ] **Step 1: Add prefix registry to specscore sourceref**

Replace `specscore/pkg/sourceref/sourceref.go` with the version that supports registerable prefixes. Key changes:
- Add `var prefixes = []string{"specscore", "synchestra"}` (keep synchestra as default for backward compat)
- Add `RegisterPrefix(prefix string)` public function
- Update `DetectionRegex` to be dynamically built from registered prefixes via `buildDetectionRegex()`
- Update `ExtractReference()` to check all registered prefixes
- Update `ParseReference()` to handle any registered prefix (not just `synchestra:`)
- Update `parseShortNotation()` to strip any registered prefix
- Update `parseExpandedURL()` to handle `specscore.io` URLs in addition to `synchestra.io`

```go
package sourceref

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Reference represents a parsed source reference found in source code.
type Reference struct {
	ResolvedPath    string
	CrossRepoSuffix string
	Type            string
}

// SourceRef represents a source file reference (file + line number).
type SourceRef struct {
	FilePath    string
	LineNumber  int
	LineContent string
}

var (
	mu       sync.Mutex
	prefixes = []string{"specscore", "synchestra"}
	domains  = []string{"specscore.io", "synchestra.io"}

	// DetectionRegex is rebuilt when prefixes change.
	DetectionRegex *regexp.Regexp
)

func init() {
	DetectionRegex = buildDetectionRegex()
}

// RegisterPrefix adds a short-notation prefix (e.g. "mytool") so that
// "mytool:feature/foo" is recognized as a source reference.
// Also registers "mytool.io" as an expanded URL domain.
func RegisterPrefix(prefix string) {
	mu.Lock()
	defer mu.Unlock()
	for _, p := range prefixes {
		if p == prefix {
			return // already registered
		}
	}
	prefixes = append(prefixes, prefix)
	domains = append(domains, prefix+".io")
	DetectionRegex = buildDetectionRegex()
}

func buildDetectionRegex() *regexp.Regexp {
	var shortParts []string
	var urlParts []string
	for _, p := range prefixes {
		shortParts = append(shortParts, regexp.QuoteMeta(p+":"))
	}
	for _, d := range domains {
		urlParts = append(urlParts, regexp.QuoteMeta("https://"+d+"/"))
	}
	all := append(shortParts, urlParts...)
	pattern := `^\s*(//|#|--|/\*|\*|%|;)\s*(` + strings.Join(all, "|") + `)`
	return regexp.MustCompile(pattern)
}

// DetectReference checks if a line contains a source reference.
func DetectReference(line string) bool {
	return DetectionRegex.MatchString(line)
}

// ExtractReference extracts the reference string from a line.
func ExtractReference(line string) string {
	// Try short notation prefixes
	for _, p := range prefixes {
		prefix := p + ":"
		if idx := strings.Index(line, prefix); idx != -1 {
			extracted := line[idx:]
			if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
				extracted = extracted[:endIdx]
			}
			return extracted
		}
	}
	// Try expanded URL domains
	for _, d := range domains {
		urlPrefix := "https://" + d + "/"
		if idx := strings.Index(line, urlPrefix); idx != -1 {
			extracted := line[idx:]
			if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
				extracted = extracted[:endIdx]
			}
			return extracted
		}
	}
	return ""
}

// ParseReference parses an extracted reference string and returns a Reference.
func ParseReference(extracted string) (*Reference, error) {
	if extracted == "" {
		return nil, fmt.Errorf("empty reference")
	}
	// Check expanded URL format
	for _, d := range domains {
		urlPrefix := "https://" + d + "/"
		if strings.HasPrefix(extracted, urlPrefix) {
			return parseExpandedURL(extracted, urlPrefix)
		}
	}
	// Check short notation
	for _, p := range prefixes {
		prefix := p + ":"
		if strings.HasPrefix(extracted, prefix) {
			return parseShortNotation(extracted, prefix)
		}
	}
	return nil, fmt.Errorf("unrecognized reference format: %s", extracted)
}

func parseExpandedURL(url, urlPrefix string) (*Reference, error) {
	path := strings.TrimPrefix(url, urlPrefix)
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid expanded URL format: too few path segments")
	}
	host := parts[0]
	org := parts[1]
	repo := parts[2]
	resolvedPath := strings.Join(parts[3:], "/")
	crossRepoSuffix := ""
	// Expanded URLs always include host/org/repo — treat as cross-repo
	// unless it matches the current project (caller can decide).
	crossRepoSuffix = fmt.Sprintf("@%s/%s/%s", host, org, repo)
	refType := inferType(resolvedPath)
	return &Reference{
		ResolvedPath:    resolvedPath,
		CrossRepoSuffix: crossRepoSuffix,
		Type:            refType,
	}, nil
}

func parseShortNotation(notation, prefix string) (*Reference, error) {
	notation = strings.TrimPrefix(notation, prefix)
	crossRepoSuffix := ""
	reference := notation
	if idx := strings.LastIndex(notation, "@"); idx != -1 {
		crossRepoSuffix = notation[idx:]
		reference = notation[:idx]
	}
	resolvedPath, err := resolveReference(reference)
	if err != nil {
		return nil, err
	}
	refType := inferType(resolvedPath)
	return &Reference{
		ResolvedPath:    resolvedPath,
		CrossRepoSuffix: crossRepoSuffix,
		Type:            refType,
	}, nil
}

func resolveReference(ref string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("empty reference")
	}
	if strings.HasPrefix(ref, "feature/") {
		return "spec/features/" + strings.TrimPrefix(ref, "feature/"), nil
	}
	if strings.HasPrefix(ref, "plan/") {
		return "spec/plans/" + strings.TrimPrefix(ref, "plan/"), nil
	}
	if strings.HasPrefix(ref, "doc/") {
		return "docs/" + strings.TrimPrefix(ref, "doc/"), nil
	}
	return ref, nil
}

func inferType(resolvedPath string) string {
	if strings.HasPrefix(resolvedPath, "spec/features/") {
		return "feature"
	}
	if strings.HasPrefix(resolvedPath, "spec/plans/") {
		return "plan"
	}
	if strings.HasPrefix(resolvedPath, "docs/") {
		return "doc"
	}
	return ""
}

// ScanLine scans a single line for references. Returns nil if none found.
func ScanLine(line string) *Reference {
	if !DetectReference(line) {
		return nil
	}
	extracted := ExtractReference(line)
	if extracted == "" {
		return nil
	}
	ref, err := ParseReference(extracted)
	if err != nil {
		return nil
	}
	return ref
}
```

- [ ] **Step 2: Adopt synchestra's scan.go in specscore (if not already equivalent)**

Ensure `specscore/pkg/sourceref/scan.go` includes all functions from synchestra's version: `ScanFiles`, `ExpandGlobPattern`, `GetUniqueReferences`, `FormatOutput`. These are already present — verify parity and merge any missing edge-case handling.

- [ ] **Step 3: Update specscore sourceref tests**

Add prefix registry tests to `specscore/pkg/sourceref/sourceref_test.go`:

```go
func TestRegisterPrefix(t *testing.T) {
	RegisterPrefix("mytool")
	line := "// mytool:feature/auth"
	if !DetectReference(line) {
		t.Error("should detect registered prefix 'mytool:'")
	}
	ref := ScanLine(line)
	if ref == nil {
		t.Fatal("ScanLine returned nil for registered prefix")
	}
	if ref.ResolvedPath != "spec/features/auth" {
		t.Errorf("expected spec/features/auth, got %s", ref.ResolvedPath)
	}
}

func TestMultiplePrefixes(t *testing.T) {
	tests := []struct {
		line     string
		detected bool
	}{
		{"// specscore:feature/cli", true},
		{"// synchestra:feature/cli", true},
		{"# specscore:plan/v2", true},
		{"// unknown:feature/cli", false},
	}
	for _, tt := range tests {
		if DetectReference(tt.line) != tt.detected {
			t.Errorf("DetectReference(%q) = %v, want %v", tt.line, !tt.detected, tt.detected)
		}
	}
}
```

Also adopt synchestra's more comprehensive test cases for `TestDetectReference`, `TestExtractReference`, `TestParseReference`, `TestResolveReference`, `TestInferType`.

- [ ] **Step 4: Run specscore tests**

```bash
cd /home/ai/projects/synchestra-io/specscore && go test ./pkg/sourceref/... -v
```

Expected: All tests PASS.

- [ ] **Step 5: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add pkg/sourceref/
git commit -m "sourceref: add pluggable prefix registry, specscore: and synchestra: supported by default"
git push origin main
```

- [ ] **Step 6: Update synchestra to import specscore's sourceref**

In `synchestra/pkg/cli/code/deps.go`, change:
```go
// Old:
"github.com/synchestra-io/synchestra/pkg/sourceref"
// New:
"github.com/synchestra-io/specscore/pkg/sourceref"
```

No startup registration needed since `synchestra:` is a default prefix.

- [ ] **Step 7: Delete synchestra's sourceref package**

```bash
rm -rf /home/ai/projects/synchestra-io/synchestra/pkg/sourceref/
```

- [ ] **Step 8: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
go fmt ./...
go build ./...
go test ./... -count=1
```

Expected: All build and tests PASS.

- [ ] **Step 9: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: replace internal sourceref with specscore/pkg/sourceref"
git push origin main
```

---

### Task 3: Migrate projectdef (with ownership split)

SpecScore's `pkg/projectdef` and synchestra's `pkg/cli/project/configfiles.go` are identical. After migration:
- SpecScore owns: `SpecConfig`, `CodeConfig`, `PlanningConfig`, read/write functions, file constants
- Synchestra owns: `StateConfig`, `EmbeddedStateConfig`, `EmbeddedSyncCfg`
- SpecScore's `SpecConfig` gets an `Extras map[string]any` field for round-tripping the `synchestra:` namespace

**Files:**
- Modify: `specscore/pkg/projectdef/projectdef.go` (remove StateConfig/EmbeddedStateConfig, add Extras field to SpecConfig)
- Modify: `specscore/pkg/projectdef/projectdef_test.go` (update tests for Extras round-trip, remove state config tests)
- Modify: `synchestra/pkg/cli/project/configfiles.go` (keep only StateConfig/EmbeddedStateConfig, import specscore for the rest)
- Modify: `synchestra/pkg/cli/project/init.go` (update imports)
- Modify: `synchestra/pkg/cli/project/new.go` (update imports)
- Modify: `synchestra/pkg/cli/project/init_test.go`, `new_test.go`, `configfiles_test.go` (update imports)

- [ ] **Step 1: Update specscore's SpecConfig with Extras field**

In `specscore/pkg/projectdef/projectdef.go`:

1. Remove `StateConfig`, `EmbeddedStateConfig`, `EmbeddedSyncCfg` types
2. Remove `ReadStateConfig`, `WriteStateConfig`, `ReadEmbeddedStateConfig`, `WriteEmbeddedStateConfig` functions
3. Remove `StateConfigFile` and `EmbeddedStateFile` constants
4. Add `Extras map[string]any` field to `SpecConfig` with `yaml:",inline"` tag
5. Keep `writeYAML` as an exported `WriteYAML` helper for synchestra to use for its own configs

Updated `SpecConfig`:
```go
type SpecConfig struct {
	Title     string          `yaml:"title"`
	StateRepo string          `yaml:"state_repo"`
	Repos     []string        `yaml:"repos"`
	Planning  *PlanningConfig `yaml:"planning,omitempty"`
	Extras    map[string]any  `yaml:",inline"`
}
```

Rename config file constants to specscore:
```go
const (
	SpecConfigFile = "specscore-spec-repo.yaml"
	CodeConfigFile = "specscore-code-repo.yaml"
)
```

Add exported YAML helper:
```go
func WriteYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
```

- [ ] **Step 2: Add Extras round-trip test**

Add test in `specscore/pkg/projectdef/projectdef_test.go`:

```go
func TestSpecConfigExtrasRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := SpecConfig{
		Title: "Test",
		Repos: []string{"github.com/org/code"},
		Extras: map[string]any{
			"synchestra": map[string]any{
				"state_repo": "github.com/org/state",
			},
		},
	}
	if err := WriteSpecConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	got, err := ReadSpecConfig(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got.Title != "Test" {
		t.Errorf("title: got %q, want %q", got.Title, "Test")
	}
	synExt, ok := got.Extras["synchestra"].(map[string]any)
	if !ok {
		t.Fatal("synchestra extension not round-tripped")
	}
	if synExt["state_repo"] != "github.com/org/state" {
		t.Errorf("state_repo: got %v", synExt["state_repo"])
	}
}
```

- [ ] **Step 3: Run specscore tests**

```bash
cd /home/ai/projects/synchestra-io/specscore && go test ./pkg/projectdef/... -v
```

Expected: All tests PASS.

- [ ] **Step 4: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add pkg/projectdef/
git commit -m "projectdef: remove state config types, add Extras for namespace extension"
git push origin main
```

- [ ] **Step 5: Rename existing config files in any test fixtures or examples**

Any test fixtures, example projects, or real projects that reference `synchestra-spec-repo.yaml` or `synchestra-code-repo.yaml` need to be renamed to `specscore-spec-repo.yaml` and `specscore-code-repo.yaml`. Search both repos:

```bash
grep -r "synchestra-spec-repo.yaml" /home/ai/projects/synchestra-io/specscore/ --include="*.go" -l
grep -r "synchestra-code-repo.yaml" /home/ai/projects/synchestra-io/specscore/ --include="*.go" -l
grep -r "synchestra-spec-repo.yaml" /home/ai/projects/synchestra-io/synchestra/ --include="*.go" -l
grep -r "synchestra-code-repo.yaml" /home/ai/projects/synchestra-io/synchestra/ --include="*.go" -l
```

Update all references. Also rename any actual YAML files on disk in test fixtures.

- [ ] **Step 6: Update synchestra's configfiles.go**

Replace `synchestra/pkg/cli/project/configfiles.go` contents. Keep only synchestra-owned types and import specscore for the rest:

```go
package project

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/synchestra-io/specscore/pkg/projectdef"
	"gopkg.in/yaml.v3"
)

// Re-export specscore constants for convenience.
const (
	SpecConfigFile = projectdef.SpecConfigFile // "specscore-spec-repo.yaml"
	CodeConfigFile = projectdef.CodeConfigFile // "specscore-code-repo.yaml"
)

// Synchestra-owned config file constants.
const (
	StateConfigFile   = "synchestra-state-repo.yaml"
	EmbeddedStateFile = "synchestra-state.yaml"
)

// StateConfig represents the contents of synchestra-state-repo.yaml.
type StateConfig struct {
	Title     string   `yaml:"title"`
	MainRepo  string   `yaml:"main_repo"`
	SpecRepos []string `yaml:"spec_repos"`
	CodeRepos []string `yaml:"code_repos,omitempty"`
}

// EmbeddedStateConfig lives on the orphan branch (inside the worktree).
type EmbeddedStateConfig struct {
	Title        string           `yaml:"title"`
	Mode         string           `yaml:"mode"`
	SourceBranch string           `yaml:"source_branch"`
	Sync         *EmbeddedSyncCfg `yaml:"sync,omitempty"`
}

// EmbeddedSyncCfg controls sync policy for embedded state.
type EmbeddedSyncCfg struct {
	Pull string `yaml:"pull"`
	Push string `yaml:"push"`
}

func WriteStateConfig(dir string, cfg StateConfig) error {
	return writeYAML(filepath.Join(dir, StateConfigFile), cfg)
}

func ReadStateConfig(dir string) (StateConfig, error) {
	var cfg StateConfig
	data, err := os.ReadFile(filepath.Join(dir, StateConfigFile))
	if err != nil {
		return cfg, fmt.Errorf("reading state config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing state config: %w", err)
	}
	return cfg, nil
}

func WriteEmbeddedStateConfig(dir string, cfg EmbeddedStateConfig) error {
	return writeYAML(filepath.Join(dir, EmbeddedStateFile), cfg)
}

func ReadEmbeddedStateConfig(dir string) (EmbeddedStateConfig, error) {
	var cfg EmbeddedStateConfig
	data, err := os.ReadFile(filepath.Join(dir, EmbeddedStateFile))
	if err != nil {
		return cfg, fmt.Errorf("reading embedded state config: %w", err)
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("parsing embedded state config: %w", err)
	}
	return cfg, nil
}

func writeYAML(path string, v any) error {
	data, err := yaml.Marshal(v)
	if err != nil {
		return fmt.Errorf("marshalling YAML: %w", err)
	}
	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("writing %s: %w", path, err)
	}
	return nil
}
```

- [ ] **Step 7: Update synchestra files that use SpecConfig/CodeConfig**

Update `init.go`, `new.go`, and their tests to use `projectdef.SpecConfig`, `projectdef.ReadSpecConfig`, `projectdef.WriteSpecConfig`, `projectdef.CodeConfig`, `projectdef.ReadCodeConfig`, `projectdef.WriteCodeConfig` instead of the local versions (which are now deleted).

Add import `"github.com/synchestra-io/specscore/pkg/projectdef"` where needed. The `PlanningConfig` type is now `projectdef.PlanningConfig`. The `WhatsNextMode()` method is on `projectdef.SpecConfig`. The `ParseStateRepo()` method is on `projectdef.SpecConfig`.

- [ ] **Step 8: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
go fmt ./...
go build ./...
go test ./pkg/cli/project/... -v
go test ./... -count=1
```

Expected: All tests PASS.

- [ ] **Step 9: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: use specscore/pkg/projectdef for SpecConfig and CodeConfig, rename config files to specscore-*"
git push origin main
```

---

### Task 4: Migrate feature package

Synchestra's `cli/feature/` has 20 files with all-private implementations. SpecScore's `pkg/feature/` has equivalent exported functions for every operation. Only `cli/main.go` imports synchestra's feature package (via `feature.Command()`). After migration, synchestra's `cli/feature/` becomes thin CLI wrappers calling `specscore/pkg/feature`.

**Files:**
- Verify: `specscore/pkg/feature/` is the superset (compare function lists)
- Rewrite: `synchestra/pkg/cli/feature/feature.go` (keep Command(), delegate to specscore)
- Rewrite: `synchestra/pkg/cli/feature/list.go` (thin wrapper)
- Rewrite: `synchestra/pkg/cli/feature/tree.go` (thin wrapper)
- Rewrite: `synchestra/pkg/cli/feature/info.go` (thin wrapper)
- Rewrite: `synchestra/pkg/cli/feature/deps.go` (thin wrapper)
- Rewrite: `synchestra/pkg/cli/feature/refs.go` (thin wrapper)
- Rewrite: `synchestra/pkg/cli/feature/new.go` (thin wrapper)
- Delete: `synchestra/pkg/cli/feature/discover.go`, `fields.go`, `slug.go`, `template.go`, `transitive.go`
- Delete: All test files for deleted implementations (keep integration tests that test CLI behavior)

- [ ] **Step 1: Verify specscore feature package is superset**

Compare every function in synchestra's `cli/feature/` against specscore's `pkg/feature/`. Key mappings:

| Synchestra (private) | SpecScore (exported) |
|---|---|
| `discoverFeatures()` | `Discover()` + `FeatureIDs()` |
| `findSpecRepoRoot()` | `FindSpecRepoRoot()` |
| `buildTree()` | `BuildTree()` |
| `printTree()` | `PrintTree()` |
| `parseDependencies()` | `ParseDependencies()` |
| `extractFeatureID()` | `ExtractFeatureID()` |
| `featureExists()` | `Exists()` |
| `featureReadmePath()` | `ReadmePath()` |
| `parseFeatureStatus()` | `ParseFeatureStatus()` |
| `findFeatureRefs()` | `FindFeatureRefs()` |
| `discoverChildFeatures()` | `DiscoverChildFeatures()` |
| `parseSections()` | `ParseSections()` |
| `countOutstandingQuestions()` | `CountOutstandingQuestions()` |
| `findLinkedPlans()` | `FindLinkedPlans()` |
| `resolveTransitiveDeps()` | `TransitiveDeps()` |
| `resolveTransitiveRefs()` | `TransitiveRefs()` |
| `buildEnrichedTree()` | `BuildEnrichedTree()` |
| `filterFocusedFeatures()` | `FilterFocusedFeatures()` |
| `markFocus()` | `MarkFocus()` |
| `generateSlug()` | `GenerateSlug()` |
| `validateSlug()` | `ValidateSlug()` |
| `generateReadme()` | `GenerateReadme()` |
| `isValidStatus()` | `IsValidStatus()` |
| `parseFieldNames()` | `ParseFieldNames()` |
| `resolveFields()` | `ResolveFields()` |
| `validateFormat()` | `ValidateFormat()` |
| `enrichedFeature` | `EnrichedFeature` |
| `featureNode` | `FeatureNode` |

Output formatting functions (`writeEnrichedYAML`, `writeEnrichedJSON`, `writeEnrichedText`, `writeTextInfo`) have NO specscore equivalent — these are CLI-only concerns and should stay in synchestra's wrapper files.

- [ ] **Step 2: Rewrite synchestra's feature command files as thin wrappers**

Each command file (list.go, tree.go, info.go, deps.go, refs.go, new.go) should:
1. Keep the cobra command definition and flag parsing
2. Replace calls to private functions with calls to `specscore/pkg/feature` exported functions
3. Keep output formatting logic (YAML/JSON/text rendering) since that's CLI-specific

Example pattern for `list.go`:
```go
package feature

import (
	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/feature"
)

func listCommand() *cobra.Command { /* same cobra definition */ }

func runList(cmd *cobra.Command, _ []string) error {
	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}
	features, err := feature.Discover(featuresDir)
	if err != nil {
		return exitcode.UnexpectedErrorf("discovering features: %v", err)
	}
	ids := feature.FeatureIDs(features)
	// ... output formatting stays here ...
}
```

The `resolveFeaturesDir` helper should stay in synchestra (it calls `feature.FindSpecRepoRoot` internally).

- [ ] **Step 3: Delete synchestra's implementation files**

Delete files that are now fully replaced by specscore imports:
```bash
rm synchestra/pkg/cli/feature/discover.go
rm synchestra/pkg/cli/feature/fields.go
rm synchestra/pkg/cli/feature/slug.go
rm synchestra/pkg/cli/feature/template.go
rm synchestra/pkg/cli/feature/transitive.go
```

Keep: `feature.go`, `list.go`, `tree.go`, `info.go`, `deps.go`, `refs.go`, `new.go` (rewritten as wrappers).

- [ ] **Step 4: Update or delete test files**

Delete unit tests for deleted private functions. Keep integration tests that test CLI behavior (e.g., `flags_integration_test.go`). Update remaining tests to use specscore types where needed.

- [ ] **Step 5: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
go fmt ./...
go build ./...
go test ./pkg/cli/feature/... -v
go test ./... -count=1
```

Expected: All tests PASS.

- [ ] **Step 6: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: replace feature implementation with specscore/pkg/feature wrappers"
git push origin main
```

---

### Task 5: Migrate lint package (with pluggable checkers)

Synchestra's `cli/spec/` and specscore's `pkg/lint/` have identical checker implementations. After migration, specscore's linter supports registering custom checkers so synchestra can inject coordination-specific rules.

**Files:**
- Modify: `specscore/pkg/lint/lint.go` (add `RegisterChecker` public function)
- Modify: `specscore/pkg/lint/linter.go` (support external checkers)
- Rewrite: `synchestra/pkg/cli/spec/spec.go` (keep Command())
- Rewrite: `synchestra/pkg/cli/spec/lint.go` (thin wrapper calling `lint.Lint()`)
- Delete: `synchestra/pkg/cli/spec/linter.go`, all checker files

- [ ] **Step 1: Add pluggable checker registration to specscore lint**

Add to `specscore/pkg/lint/lint.go`:

```go
// Checker is the public interface for custom rule implementations.
// External tools can implement this to add custom rules to the linter.
type Checker interface {
	Name() string
	Severity() string
	Check(specRoot string) ([]Violation, error)
}

var customCheckers []Checker

// RegisterChecker registers a custom checker that will run alongside
// built-in checkers during Lint().
func RegisterChecker(c Checker) {
	customCheckers = append(customCheckers, c)
}

// ResetCustomCheckers clears all registered custom checkers (for testing).
func ResetCustomCheckers() {
	customCheckers = nil
}
```

Update `newLinter()` in `linter.go` to register custom checkers:
```go
func newLinter(opts Options) *linter {
	l := &linter{
		opts:    opts,
		ruleSet: make(map[string]checker),
	}
	// ... existing built-in registrations ...

	// Register custom checkers
	for _, c := range customCheckers {
		l.ruleSet[c.Name()] = &customCheckerAdapter{c}
	}

	return l
}

// customCheckerAdapter adapts the public Checker interface to the internal checker interface.
type customCheckerAdapter struct {
	c Checker
}

func (a *customCheckerAdapter) name() string                          { return a.c.Name() }
func (a *customCheckerAdapter) severity() string                      { return a.c.Severity() }
func (a *customCheckerAdapter) check(specRoot string) ([]Violation, error) { return a.c.Check(specRoot) }
```

- [ ] **Step 2: Add custom checker test**

```go
func TestRegisterChecker(t *testing.T) {
	defer ResetCustomCheckers()

	RegisterChecker(&testChecker{
		n: "custom-rule",
		s: "warning",
		violations: []Violation{{
			File: "test.md", Line: 1, Severity: "warning",
			Rule: "custom-rule", Message: "custom violation",
		}},
	})

	dir := t.TempDir()
	// Create minimal spec structure
	os.MkdirAll(filepath.Join(dir, "spec", "features"), 0755)
	os.WriteFile(filepath.Join(dir, "spec", "features", "README.md"), []byte("# Features\n\n## Outstanding Questions\n\nNone.\n"), 0644)

	violations, err := Lint(Options{SpecRoot: dir, Rules: []string{"custom-rule"}})
	if err != nil {
		t.Fatal(err)
	}
	if len(violations) != 1 || violations[0].Rule != "custom-rule" {
		t.Errorf("expected 1 custom violation, got %d", len(violations))
	}
}

type testChecker struct {
	n          string
	s          string
	violations []Violation
}

func (c *testChecker) Name() string                             { return c.n }
func (c *testChecker) Severity() string                         { return c.s }
func (c *testChecker) Check(string) ([]Violation, error) { return c.violations, nil }
```

- [ ] **Step 3: Run specscore tests**

```bash
cd /home/ai/projects/synchestra-io/specscore && go test ./pkg/lint/... -v
```

Expected: All tests PASS.

- [ ] **Step 4: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add pkg/lint/
git commit -m "lint: add pluggable checker registration for external tools"
git push origin main
```

- [ ] **Step 5: Rewrite synchestra's spec command as thin wrapper**

Replace `synchestra/pkg/cli/spec/lint.go` with:

```go
package spec

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/lint"
	"gopkg.in/yaml.v3"
)

func lintCommand() *cobra.Command {
	// Same cobra definition as before — flags: --project, --rules, --ignore, --severity, --format
}

func runLint(cmd *cobra.Command, _ []string) error {
	// Parse flags same as before
	// Build lint.Options from flags
	opts := lint.Options{
		SpecRoot: specRoot,
		Rules:    rules,
		Ignore:   ignore,
		Severity: severity,
	}

	violations, err := lint.Lint(opts)
	if err != nil {
		return exitcode.UnexpectedErrorf("linting: %v", err)
	}

	if severity != "" {
		violations = lint.FilterBySeverity(violations, severity)
	}

	// Output formatting (text/json/yaml) stays here — same logic as before
	// ...

	if hasErrors(violations) {
		return exitcode.New(exitcode.Conflict, "lint violations found")
	}
	return nil
}
```

- [ ] **Step 6: Delete synchestra's lint implementation files**

```bash
rm synchestra/pkg/cli/spec/linter.go
rm synchestra/pkg/cli/spec/readme_exists.go
rm synchestra/pkg/cli/spec/oq_section.go
rm synchestra/pkg/cli/spec/index_entries.go
rm synchestra/pkg/cli/spec/plan_roi.go
rm synchestra/pkg/cli/spec/plan_hierarchy.go
rm synchestra/pkg/cli/spec/checkers_extended.go
```

Keep: `spec.go` (Command entry point), `lint.go` (rewritten wrapper).

- [ ] **Step 7: Delete or update test files**

Delete `lint_test.go`, `plan_roi_test.go`, `plan_hierarchy_test.go` — these test checker internals that now live in specscore. If there are integration tests for the CLI command, keep and update them.

- [ ] **Step 8: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
go fmt ./...
go build ./...
go test ./... -count=1
```

Expected: All tests PASS.

- [ ] **Step 9: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: replace lint implementation with specscore/pkg/lint wrappers"
git push origin main
```

---

### Task 6: Extract task types and read operations to specscore

This is the most complex step. Create a new `specscore/pkg/task` package containing:
- Task types: `Task`, `TaskStatus`, `TaskCreateParams`, `TaskFilter`, `BoardView`, `BoardRow`
- Task file format: YAML frontmatter parsing/serialization
- Board format: markdown table rendering/parsing
- Read operations and CLI commands: `task list`, `task info`, `task new`

Synchestra's `pkg/state` then imports these types and keeps the lifecycle state machine.

**Files:**
- Create: `specscore/pkg/task/task.go` (types and status enum)
- Create: `specscore/pkg/task/taskfile.go` (README format parsing/rendering)
- Create: `specscore/pkg/task/board.go` (markdown board format)
- Create: `specscore/pkg/task/task_test.go`
- Create: `specscore/pkg/task/taskfile_test.go`
- Create: `specscore/pkg/task/board_test.go`
- Create: `specscore/internal/cli/task.go` (task list/info/new commands for specscore CLI)
- Modify: `synchestra/pkg/state/types.go` (import task types from specscore, keep coordination-only types)
- Modify: `synchestra/pkg/state/task.go` (TaskStore uses specscore task types)
- Modify: `synchestra/pkg/state/gitstore/taskfile.go` (import specscore task format)
- Modify: `synchestra/pkg/state/gitstore/board.go` (import specscore board format)
- Modify: `synchestra/pkg/state/gitstore/task.go` (use specscore task types)

- [ ] **Step 1: Create specscore/pkg/task/task.go**

Extract types from `synchestra/pkg/state/types.go`:

```go
package task

import "time"

// TaskStatus represents the lifecycle state of a task.
type TaskStatus string

const (
	StatusPlanning   TaskStatus = "planning"
	StatusQueued     TaskStatus = "queued"
	StatusClaimed    TaskStatus = "claimed"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusBlocked    TaskStatus = "blocked"
	StatusAborted    TaskStatus = "aborted"
)

// Task represents a unit of work.
type Task struct {
	Slug      string
	Title     string
	Status    TaskStatus
	Parent    string
	DependsOn []string
	Requester string
	Reason    string
	Summary   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateParams holds parameters for creating a new task.
type CreateParams struct {
	Slug      string
	Title     string
	Parent    string
	DependsOn []string
	Requester string
}

// Filter holds optional filters for listing tasks.
type Filter struct {
	Status *TaskStatus
	Parent *string
}

// BoardView represents a rendered task board.
type BoardView struct {
	Rows []BoardRow
}

// BoardRow represents a single row in the task board.
type BoardRow struct {
	Task      string
	Status    TaskStatus
	DependsOn []string
	Branch    string
	Agent     string
	Requester string
	StartedAt *time.Time
	Duration  *time.Duration
}
```

- [ ] **Step 2: Create specscore/pkg/task/taskfile.go**

Extract format parsing from `synchestra/pkg/state/gitstore/taskfile.go`:

```go
package task

import (
	"bufio"
	"fmt"
	"strings"
)

// TaskFileData holds the parsed contents of a task README.md.
type TaskFileData struct {
	Title       string
	Description string
	DependsOn   []string
	Summary     string
}

// ParseTaskFile parses a task README.md content string.
// The format is: # Title\n\nDescription\n\n## Dependencies\n\n- dep1\n- dep2\n\n## Summary\n\nSummary text
// Read the full implementation from synchestra/pkg/state/gitstore/taskfile.go:parseTaskFile()
// and export it as a public function. The logic scans lines for the title heading, collects
// description lines, parses the Dependencies section as a bullet list, and captures the Summary
// section. No changes to the parsing logic — just make it public and move it here.
func ParseTaskFile(content string) (*TaskFileData, error) {
	scanner := bufio.NewScanner(strings.NewReader(content))
	data := &TaskFileData{}
	section := "description"
	var descLines []string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "# ") && data.Title == "" {
			data.Title = strings.TrimPrefix(trimmed, "# ")
			continue
		}
		if trimmed == "## Dependencies" {
			section = "dependencies"
			continue
		}
		if trimmed == "## Summary" {
			section = "summary"
			continue
		}
		if strings.HasPrefix(trimmed, "## ") {
			section = "other"
			continue
		}

		switch section {
		case "description":
			descLines = append(descLines, line)
		case "dependencies":
			if strings.HasPrefix(trimmed, "- ") {
				dep := strings.TrimPrefix(trimmed, "- ")
				data.DependsOn = append(data.DependsOn, dep)
			}
		case "summary":
			if data.Summary == "" {
				data.Summary = trimmed
			} else if trimmed != "" {
				data.Summary += "\n" + line
			}
		}
	}

	data.Description = strings.TrimSpace(strings.Join(descLines, "\n"))
	data.Summary = strings.TrimSpace(data.Summary)
	return data, nil
}

// RenderTaskFile renders TaskFileData back to markdown.
func RenderTaskFile(data *TaskFileData) string {
	var sb strings.Builder
	fmt.Fprintf(&sb, "# %s\n\n", data.Title)
	if data.Description != "" {
		sb.WriteString(data.Description)
		sb.WriteString("\n\n")
	}
	if len(data.DependsOn) > 0 {
		sb.WriteString("## Dependencies\n\n")
		for _, dep := range data.DependsOn {
			fmt.Fprintf(&sb, "- %s\n", dep)
		}
		sb.WriteString("\n")
	}
	if data.Summary != "" {
		sb.WriteString("## Summary\n\n")
		sb.WriteString(data.Summary)
		sb.WriteString("\n")
	}
	return sb.String()
}
```

- [ ] **Step 3: Create specscore/pkg/task/board.go**

Extract board format from `synchestra/pkg/state/gitstore/board.go`:

```go
package task

import (
	"fmt"
	"strings"
)

// StatusEmojis maps task statuses to their emoji indicators.
var StatusEmojis = map[TaskStatus]string{
	StatusPlanning:   "📋",
	StatusQueued:     "⏳",
	StatusClaimed:    "🔒",
	StatusInProgress: "🔵",
	StatusCompleted:  "✅",
	StatusFailed:     "❌",
	StatusBlocked:    "🟡",
	StatusAborted:    "⛔",
}

// ParseBoard parses a tasks/README.md markdown table into BoardView.
// The table has 7 columns: Task | Status | Depends on | Branch | Agent | Requester | Time
// Read the full implementation from synchestra/pkg/state/gitstore/board.go:parseBoard()
// and export it as a public function. The logic:
// 1. Skips header and separator rows
// 2. Splits each row by | into 7 columns
// 3. Parses task name from markdown link [slug](./slug/)
// 4. Parses status from emoji+backtick format (e.g., "🔵 `in_progress`")
// 5. Parses depends-on as comma-separated list
// 6. Strips backticks from branch name
// 7. Parses time as duration
// No changes to parsing logic — just make it public.
func ParseBoard(content string) (*BoardView, error) { /* extract from synchestra */ }

// RenderBoard renders a BoardView back to markdown table.
// Renders the 7-column markdown table with emoji status indicators.
// Read the full implementation from synchestra/pkg/state/gitstore/board.go:renderBoard()
// and export it as a public function. No changes to rendering logic.
func RenderBoard(board *BoardView) string { /* extract from synchestra */ }
```

- [ ] **Step 4: Create tests for specscore task package**

Copy and adapt tests from synchestra's `gitstore/taskfile_test.go` and `gitstore/board_test.go`. These test the format parsing/rendering in isolation — no git or store dependencies.

- [ ] **Step 5: Run specscore tests**

```bash
cd /home/ai/projects/synchestra-io/specscore && go test ./pkg/task/... -v
```

Expected: All tests PASS.

- [ ] **Step 6: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add pkg/task/
git commit -m "task: add task types, file format parsing, and board rendering"
git push origin main
```

- [ ] **Step 7: Update synchestra's state types to use specscore task types**

In `synchestra/pkg/state/types.go`:
- Remove `Task`, `TaskStatus`, `TaskCreateParams`, `TaskFilter`, `BoardView`, `BoardRow` (now in specscore)
- Keep `ClaimParams`, `Chat`, `ChatStatus`, `ChatMessage`, `ArtifactRef`, `ProjectConfig`, `StoreOptions`, `SyncConfig`, `SyncPolicy`
- Add `CoordinatedTask` that embeds `specscore/pkg/task.Task` with coordination fields:

```go
package state

import (
	"time"

	"github.com/synchestra-io/specscore/pkg/task"
)

// CoordinatedTask extends specscore's Task with coordination fields.
type CoordinatedTask struct {
	task.Task
	Run       string
	Model     string
	ClaimedAt *time.Time
}

// ClaimParams holds parameters for claiming a task.
type ClaimParams struct {
	Run   string
	Model string
}
```

- [ ] **Step 8: Update synchestra's TaskStore interface**

In `synchestra/pkg/state/task.go`, update method signatures to use `task.Task`, `task.CreateParams`, `task.Filter`, etc. from specscore. Methods that return tasks should return `CoordinatedTask` (which embeds `task.Task`).

- [ ] **Step 9: Update synchestra's gitstore to use specscore format parsing**

In `synchestra/pkg/state/gitstore/taskfile.go` — replace internal `parseTaskFile`/`renderTaskFile` with calls to `task.ParseTaskFile`/`task.RenderTaskFile` from specscore.

In `synchestra/pkg/state/gitstore/board.go` — replace internal `parseBoard`/`renderBoard` with calls to `task.ParseBoard`/`task.RenderBoard`. Keep the optimistic locking logic that wraps these format operations.

In `synchestra/pkg/state/gitstore/task.go` — update all type references.

- [ ] **Step 10: Update synchestra CLI task commands**

Update all files in `synchestra/pkg/cli/task/` to use `specscore/pkg/task` types where they reference `state.Task`, `state.TaskStatus`, etc. The `state.CoordinatedTask` type is used when coordination fields are needed.

- [ ] **Step 11: Run synchestra build and tests**

```bash
cd /home/ai/projects/synchestra-io/synchestra
go get github.com/synchestra-io/specscore@latest
go fmt ./...
go build ./...
go test ./... -count=1
```

Expected: All tests PASS.

- [ ] **Step 12: Commit and push synchestra**

```bash
cd /home/ai/projects/synchestra-io/synchestra
git add -A
git commit -m "refactor: use specscore/pkg/task for task types and format parsing"
git push origin main
```

---

### Task 7: Add task CLI commands to specscore

After Task 6 establishes the specscore task package, add read-only CLI commands to specscore: `task list`, `task info`, `task new`.

**Files:**
- Create: `specscore/internal/cli/task.go` (cobra commands)
- Modify: `specscore/internal/cli/root.go` (register task subcommand)

- [ ] **Step 1: Create specscore task CLI commands**

In `specscore/internal/cli/task.go`:

```go
package cli

import (
	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/task"
)

func taskCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "task",
		Short: "Task management (read operations)",
	}
	cmd.AddCommand(
		taskListCommand(),
		taskInfoCommand(),
		taskNewCommand(),
	)
	return cmd
}

func taskListCommand() *cobra.Command { /* list tasks using task.ParseBoard */ }
func taskInfoCommand() *cobra.Command { /* show task details using task.ParseTaskFile */ }
func taskNewCommand() *cobra.Command  { /* create task using task.RenderTaskFile */ }
```

- [ ] **Step 2: Register task command in root**

In `specscore/internal/cli/root.go`, add `taskCommand()` to the command tree.

- [ ] **Step 3: Run specscore build and tests**

```bash
cd /home/ai/projects/synchestra-io/specscore
go fmt ./...
go build ./...
go test ./... -count=1
```

Expected: All tests PASS.

- [ ] **Step 4: Commit and push specscore**

```bash
cd /home/ai/projects/synchestra-io/specscore
git add internal/cli/task.go internal/cli/root.go
git commit -m "cli: add task list/info/new read-only commands"
git push origin main
```
