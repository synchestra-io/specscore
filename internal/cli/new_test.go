package cli

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/synchestra-io/specscore/pkg/idea"
)

// setupSpecRoot stages a temp spec repo with empty indexes and returns the
// repo root (the parent of `spec/`). The test's CWD is set to root so that
// FindSpecRepoRoot picks it up.
func setupSpecRoot(t *testing.T) string {
	t.Helper()
	root := t.TempDir()
	specDir := filepath.Join(root, "spec")
	featDir := filepath.Join(specDir, "features")
	if err := os.MkdirAll(featDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	ideasDir := filepath.Join(specDir, "ideas")
	if err := os.MkdirAll(filepath.Join(ideasDir, "archived"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	idx := "# Ideas\n\n## Index\n\n| Idea | Status | Date | Owner | Promotes To |\n|------|--------|------|-------|-------------|\n\n_No active ideas yet._\n\n## Outstanding Questions\n\nNone at this time.\n"
	_ = os.WriteFile(filepath.Join(ideasDir, "README.md"), []byte(idx), 0o644)
	arch := "# Archived\n\n_No archived ideas yet._\n\n## Outstanding Questions\n\nNone at this time.\n"
	_ = os.WriteFile(filepath.Join(ideasDir, "archived", "README.md"), []byte(arch), 0o644)
	return root
}

func withCwd(t *testing.T, dir string) {
	t.Helper()
	old, _ := os.Getwd()
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(old) })
}

func runNew(t *testing.T, args ...string) (string, string, error) {
	t.Helper()
	cmd := newCommand()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return out.String(), errOut.String(), err
}

func TestNewIdea_BareInvocationLintClean(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	_, _, err := runNew(t, "idea", "demo-bare")
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	path := filepath.Join(root, "spec", "ideas", "demo-bare.md")
	if _, err := os.Stat(path); err != nil {
		t.Fatalf("expected idea file at %s: %v", path, err)
	}
	// Index should now list demo-bare.
	idx, _ := os.ReadFile(filepath.Join(root, "spec", "ideas", "README.md"))
	if !strings.Contains(string(idx), "[demo-bare](demo-bare.md)") {
		t.Errorf("index not updated: %s", idx)
	}
}

func TestNewIdea_FlagsInject(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	_, _, err := runNew(t, "idea", "demo-flag",
		"--title", "Demo Flag",
		"--owner", "alice",
		"--hmw", "How might we flag this?",
		"--not-doing", "thing one — reason",
	)
	if err != nil {
		t.Fatalf("command failed: %v", err)
	}
	body, _ := os.ReadFile(filepath.Join(root, "spec", "ideas", "demo-flag.md"))
	s := string(body)
	for _, want := range []string{
		"# Idea: Demo Flag",
		"**Owner:** alice",
		"How might we flag this?",
		"- thing one — reason",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("missing %q in generated body:\n%s", want, s)
		}
	}
}

func TestNewIdea_DuplicateRefused(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	if _, _, err := runNew(t, "idea", "dup"); err != nil {
		t.Fatalf("first run: %v", err)
	}
	_, _, err := runNew(t, "idea", "dup")
	if err == nil {
		t.Fatal("expected second run to fail without --force")
	}
	if !strings.Contains(err.Error(), "already exists") {
		t.Errorf("unexpected error: %v", err)
	}

	// With --force, succeeds.
	if _, _, err := runNew(t, "idea", "dup", "--force", "--title", "After Force"); err != nil {
		t.Fatalf("--force run: %v", err)
	}
	body, _ := os.ReadFile(filepath.Join(root, "spec", "ideas", "dup.md"))
	if !strings.Contains(string(body), "# Idea: After Force") {
		t.Errorf("expected overwrite, got:\n%s", body)
	}
}

func TestNewIdea_InvalidSlug(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	_, _, err := runNew(t, "idea", "BadSlug")
	if err == nil {
		t.Fatal("expected error for invalid slug")
	}
}

func TestNewIdea_Interactive(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	cmd := newCommand()
	var out, errOut bytes.Buffer
	cmd.SetOut(&out)
	cmd.SetErr(&errOut)
	// Provide: title, owner, hmw, then blank-through the rest, then a
	// single not-doing bullet, then blank to end.
	cmd.SetIn(strings.NewReader("Interactive Demo\nzoe\nHow might we interact?\n\n\n\nthing interactive — reason\n\n"))
	cmd.SetArgs([]string{"idea", "demo-inter", "-i"})
	if err := cmd.Execute(); err != nil {
		t.Fatalf("interactive run failed: %v\nerr: %s", err, errOut.String())
	}
	body, _ := os.ReadFile(filepath.Join(root, "spec", "ideas", "demo-inter.md"))
	s := string(body)
	if !strings.Contains(s, "# Idea: Interactive Demo") {
		t.Errorf("title not captured from interactive input:\n%s", s)
	}
	if !strings.Contains(s, "**Owner:** zoe") {
		t.Errorf("owner not captured:\n%s", s)
	}
	if !strings.Contains(s, "How might we interact?") {
		t.Errorf("hmw not captured:\n%s", s)
	}
	if !strings.Contains(s, "- thing interactive — reason") {
		t.Errorf("not-doing not captured:\n%s", s)
	}
}

func TestNewIdea_IndexRowInserted(t *testing.T) {
	root := setupSpecRoot(t)
	withCwd(t, root)

	if _, _, err := runNew(t, "idea", "idx-check", "--owner", "pat"); err != nil {
		t.Fatalf("run: %v", err)
	}
	idx, _ := os.ReadFile(filepath.Join(root, "spec", "ideas", "README.md"))
	s := string(idx)
	if !strings.Contains(s, "[idx-check](idx-check.md)") {
		t.Errorf("missing idea link in index:\n%s", s)
	}
	if !strings.Contains(s, "pat") {
		t.Errorf("missing owner in index:\n%s", s)
	}
}

// Ensure the exported idea.Scaffold still produces a valid file for a
// hand-chosen set of options.
func TestNewIdea_ScaffoldExported(t *testing.T) {
	body, err := idea.Scaffold(idea.ScaffoldOptions{Slug: "export-check"})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if !strings.Contains(string(body), "# Idea: Export Check") {
		t.Errorf("unexpected title: %s", body)
	}
}
