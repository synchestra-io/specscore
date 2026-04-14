package lint

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/synchestra-io/specscore/pkg/gitremote"
	"github.com/synchestra-io/specscore/pkg/projectdef"
)

const (
	testHubHost = "https://hub.synchestra.io"
	testOwner   = "synchestra-io"
	testRepo    = "specscore"
)

// setupHubProject creates a minimal project root with:
//   - specscore-spec-repo.yaml containing hub.host
//   - a git repo whose origin points at a GitHub URL
//   - spec/features/ scaffolding
//
// Returns the project root path; callers should use filepath.Join(root, "spec") as specRoot.
func setupHubProject(t *testing.T) string {
	t.Helper()
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in test environment")
	}
	root := t.TempDir()

	// Config opts in with hub.host.
	if err := projectdef.WriteSpecConfig(root, projectdef.SpecConfig{
		Title: "Test",
		Hub:   &projectdef.HubConfig{Host: testHubHost},
	}); err != nil {
		t.Fatal(err)
	}

	// Initialize a git repo with a GitHub origin.
	run := func(args ...string) {
		cmd := exec.Command("git", append([]string{"-C", root}, args...)...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("remote", "add", "origin", "git@github.com:"+testOwner+"/"+testRepo+".git")

	if err := os.MkdirAll(filepath.Join(root, "spec", "features"), 0o755); err != nil {
		t.Fatal(err)
	}
	return root
}

func writeHubFeatureReadme(t *testing.T, root, slug, content string) string {
	t.Helper()
	dir := filepath.Join(root, "spec", "features", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return path
}

func expectedBlockquote(slug string) string {
	r := gitremote.Remote{Owner: testOwner, Repo: testRepo, Host: "github.com"}
	u := BuildHubViewURL(testHubHost, r, "spec/features/"+slug)
	return hubViewLinkMarker + u + hubViewLinkSuffix
}

func runHubCheck(t *testing.T, root string) []Violation {
	t.Helper()
	c := newHubViewLinkChecker()
	v, err := c.check(filepath.Join(root, "spec"))
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	return v
}

func runHubFix(t *testing.T, root string) {
	t.Helper()
	c := newHubViewLinkChecker().(fixer)
	if err := c.fix(filepath.Join(root, "spec")); err != nil {
		t.Fatalf("fix: %v", err)
	}
}

func TestHubViewLink_DisabledWhenNoConfig(t *testing.T) {
	root := t.TempDir()
	if err := os.MkdirAll(filepath.Join(root, "spec", "features", "x"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(root, "spec", "features", "x", "README.md"),
		[]byte("# Feature: X\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	// No specscore-spec-repo.yaml → rule is a no-op.
	c := newHubViewLinkChecker()
	v, err := c.check(filepath.Join(root, "spec"))
	if err != nil {
		t.Fatalf("check: %v", err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations when config absent, got %+v", v)
	}
}

func TestHubViewLink_MissingReported(t *testing.T) {
	root := setupHubProject(t)
	writeHubFeatureReadme(t, root, "auth", "# Feature: Auth\n\n**Status:** Draft\n")
	v := runHubCheck(t, root)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(v), v)
	}
	if v[0].Rule != "hub-view-link" || v[0].Severity != "warning" {
		t.Errorf("bad rule/severity: %+v", v[0])
	}
	if v[0].File != filepath.Join("features", "auth", "README.md") {
		t.Errorf("file = %q", v[0].File)
	}
}

func TestHubViewLink_PresentNoViolation(t *testing.T) {
	root := setupHubProject(t)
	content := "# Feature: Auth\n\n" + expectedBlockquote("auth") + "\n\n**Status:** Draft\n"
	writeHubFeatureReadme(t, root, "auth", content)
	v := runHubCheck(t, root)
	if len(v) != 0 {
		t.Errorf("expected 0 violations, got %+v", v)
	}
}

func TestHubViewLink_StaleReported(t *testing.T) {
	root := setupHubProject(t)
	// Wrong URL — still has the marker prefix, so classified as stale.
	stale := hubViewLinkMarker + "https://hub.example.com/wrong" + hubViewLinkSuffix
	content := "# Feature: Auth\n\n" + stale + "\n"
	writeHubFeatureReadme(t, root, "auth", content)
	v := runHubCheck(t, root)
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %+v", v)
	}
	if !strings.Contains(v[0].Message, "out of date") {
		t.Errorf("expected 'out of date' message, got %q", v[0].Message)
	}
}

func TestHubViewLink_FixInsertsUnderH1(t *testing.T) {
	root := setupHubProject(t)
	path := writeHubFeatureReadme(t, root, "auth", "# Feature: Auth\n\n**Status:** Draft\n")
	runHubFix(t, root)
	out, _ := os.ReadFile(path)
	expected := expectedBlockquote("auth")
	if !strings.Contains(string(out), expected) {
		t.Errorf("fix did not insert expected blockquote.\nGot:\n%s\nExpected line:\n%s", out, expected)
	}
	// H1 must still be the first line.
	if !strings.HasPrefix(string(out), "# Feature: Auth\n") {
		t.Errorf("H1 not preserved at top: %q", string(out))
	}
	// Idempotent: second run should not change anything.
	before, _ := os.ReadFile(path)
	runHubFix(t, root)
	after, _ := os.ReadFile(path)
	if string(before) != string(after) {
		t.Errorf("fix not idempotent:\nbefore:\n%s\nafter:\n%s", before, after)
	}
}

func TestHubViewLink_FixReplacesStale(t *testing.T) {
	root := setupHubProject(t)
	stale := hubViewLinkMarker + "https://hub.example.com/wrong" + hubViewLinkSuffix
	path := writeHubFeatureReadme(t, root, "auth",
		"# Feature: Auth\n\n"+stale+"\n\n**Status:** Draft\n")
	runHubFix(t, root)
	out, _ := os.ReadFile(path)
	if strings.Contains(string(out), "hub.example.com/wrong") {
		t.Errorf("stale URL not replaced: %s", out)
	}
	if !strings.Contains(string(out), expectedBlockquote("auth")) {
		t.Errorf("expected blockquote missing: %s", out)
	}
}

func TestHubViewLink_NonGitHubRemoteSkips(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available")
	}
	root := t.TempDir()
	if err := projectdef.WriteSpecConfig(root, projectdef.SpecConfig{
		Title: "T", Hub: &projectdef.HubConfig{Host: testHubHost},
	}); err != nil {
		t.Fatal(err)
	}
	run := func(args ...string) {
		if out, err := exec.Command("git", append([]string{"-C", root}, args...)...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}
	run("init", "-q")
	run("remote", "add", "origin", "https://gitlab.com/o/r.git")
	writeHubFeatureReadme(t, root, "auth", "# Feature: Auth\n")

	v := runHubCheck(t, root)
	if len(v) != 1 {
		t.Fatalf("expected 1 config-level violation, got %+v", v)
	}
	if v[0].File != projectdef.SpecConfigFile {
		t.Errorf("expected violation on config file, got %q", v[0].File)
	}
}

func TestHubViewLink_UnderscoreDirsIgnored(t *testing.T) {
	root := setupHubProject(t)
	// Valid feature with expected blockquote.
	writeHubFeatureReadme(t, root, "auth",
		"# Feature: Auth\n\n"+expectedBlockquote("auth")+"\n")
	// _tests subtree must be skipped entirely.
	testsDir := filepath.Join(root, "spec", "features", "auth", "_tests")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "README.md"),
		[]byte("# Tests\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	v := runHubCheck(t, root)
	if len(v) != 0 {
		t.Errorf("expected 0 violations, got %+v", v)
	}
}

func TestHubViewLink_RegisteredAsKnownRule(t *testing.T) {
	if !AllRuleNames()["hub-view-link"] {
		t.Error("hub-view-link should be a known rule")
	}
}

func TestBuildHubViewURL(t *testing.T) {
	r := gitremote.Remote{Owner: "synchestra-io", Repo: "specscore", Host: "github.com"}
	got := BuildHubViewURL(testHubHost, r, "spec/features/bots")
	want := testHubHost + "/project/features?id=synchestra-io@specscore@github.com&path=spec%2Ffeatures%2Fbots"
	if got != want {
		t.Errorf("BuildHubViewURL =\n  %q\nwant\n  %q", got, want)
	}
}
