package idea_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/synchestra-io/specscore/pkg/idea"
	"github.com/synchestra-io/specscore/pkg/lint"
)

// writeSpecTree stages a minimal spec tree with a single idea file at
// spec/ideas/<slug>.md and an empty active index. Returns the spec root
// (i.e. the directory containing "ideas/").
func writeSpecTree(t *testing.T, body []byte, slug string) string {
	t.Helper()
	root := t.TempDir()
	ideasDir := filepath.Join(root, "ideas")
	if err := os.MkdirAll(filepath.Join(ideasDir, "archived"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	index := "# Ideas\n\n## Index\n\n| Idea | Status | Date | Owner | Promotes To |\n|------|--------|------|-------|-------------|\n\n_No active ideas yet._\n\n## Outstanding Questions\n\nNone at this time.\n"
	if err := os.WriteFile(filepath.Join(ideasDir, "README.md"), []byte(index), 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	archIdx := "# Archived Ideas\n\n_No archived ideas yet._\n\n## Outstanding Questions\n\nNone at this time.\n"
	if err := os.WriteFile(filepath.Join(ideasDir, "archived", "README.md"), []byte(archIdx), 0o644); err != nil {
		t.Fatalf("write arch idx: %v", err)
	}
	if err := os.WriteFile(filepath.Join(ideasDir, slug+".md"), body, 0o644); err != nil {
		t.Fatalf("write idea: %v", err)
	}
	// After scaffolding, run lint --fix once to reconcile the index.
	if _, err := lint.Lint(lint.Options{SpecRoot: root, Fix: true}); err != nil {
		t.Fatalf("lint fix: %v", err)
	}
	return root
}

func ideaErrorsFor(t *testing.T, specRoot, slug string) []lint.Violation {
	t.Helper()
	vs, err := lint.Lint(lint.Options{SpecRoot: specRoot})
	if err != nil {
		t.Fatalf("lint: %v", err)
	}
	var errs []lint.Violation
	for _, v := range vs {
		if v.Severity == "error" && strings.Contains(v.File, slug) {
			errs = append(errs, v)
		}
	}
	return errs
}

func TestScaffold_BareIsLintClean(t *testing.T) {
	body, err := idea.Scaffold(idea.ScaffoldOptions{Slug: "demo-bare"})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	if !strings.Contains(string(body), "# Idea: Demo Bare") {
		t.Errorf("expected title `# Idea: Demo Bare`, body:\n%s", body)
	}
	root := writeSpecTree(t, body, "demo-bare")
	if errs := ideaErrorsFor(t, root, "demo-bare"); len(errs) > 0 {
		t.Errorf("bare scaffold failed lint: %+v", errs)
	}
}

func TestScaffold_FlagInjection(t *testing.T) {
	body, err := idea.Scaffold(idea.ScaffoldOptions{
		Slug:     "demo-flags",
		Title:    "Demo Flags",
		Owner:    "alice",
		HMW:      "How might we demo flag injection?",
		NotDoing: []string{"thing one — reason", "thing two — reason"},
	})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	s := string(body)
	for _, want := range []string{
		"# Idea: Demo Flags",
		"**Owner:** alice",
		"How might we demo flag injection?",
		"- thing one — reason",
		"- thing two — reason",
	} {
		if !strings.Contains(s, want) {
			t.Errorf("generated body missing %q\n%s", want, s)
		}
	}
	root := writeSpecTree(t, body, "demo-flags")
	if errs := ideaErrorsFor(t, root, "demo-flags"); len(errs) > 0 {
		t.Errorf("flagged scaffold failed lint: %+v", errs)
	}
}

func TestScaffold_LintCleanMatrix(t *testing.T) {
	cases := []struct {
		name string
		opts idea.ScaffoldOptions
	}{
		{"bare", idea.ScaffoldOptions{Slug: "m-bare"}},
		{"title-only", idea.ScaffoldOptions{Slug: "m-title", Title: "Custom Title"}},
		{"owner-date", idea.ScaffoldOptions{Slug: "m-od", Owner: "bob", Date: "2025-06-01"}},
		{"full", idea.ScaffoldOptions{
			Slug:                 "m-full",
			Title:                "Full Example",
			Owner:                "carol",
			HMW:                  "How might we fully exercise the scaffolder?",
			Context:              "This is context.",
			RecommendedDirection: "We should do X because Y.",
			Alternatives:         []string{"A: did not work", "B: too expensive"},
			MVP:                  "Single-week prototype.",
			NotDoing:             []string{"not A — reason", "not B — reason"},
		}},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			body, err := idea.Scaffold(tc.opts)
			if err != nil {
				t.Fatalf("Scaffold: %v", err)
			}
			root := writeSpecTree(t, body, tc.opts.Slug)
			if errs := ideaErrorsFor(t, root, tc.opts.Slug); len(errs) > 0 {
				t.Errorf("%s variant failed lint: %+v\nbody:\n%s", tc.name, errs, body)
			}
		})
	}
}

func TestScaffold_InvalidSlug(t *testing.T) {
	if _, err := idea.Scaffold(idea.ScaffoldOptions{Slug: "BadSlug"}); err == nil {
		t.Error("expected error for invalid slug, got nil")
	}
	if _, err := idea.Scaffold(idea.ScaffoldOptions{Slug: ""}); err == nil {
		t.Error("expected error for empty slug, got nil")
	}
}

func TestScaffold_DefaultsApplied(t *testing.T) {
	body, err := idea.Scaffold(idea.ScaffoldOptions{Slug: "defaults-test"})
	if err != nil {
		t.Fatalf("Scaffold: %v", err)
	}
	s := string(body)
	if !strings.Contains(s, "**Status:** Draft") {
		t.Error("default status should be Draft")
	}
	if !strings.Contains(s, "# Idea: Defaults Test") {
		t.Errorf("default title should be title-cased slug; got:\n%s", s)
	}
	if !strings.Contains(s, "**Promotes To:** —") {
		t.Error("Promotes To should default to em-dash")
	}
	if !strings.Contains(s, "<!-- One \"How Might We") {
		t.Error("HMW prompt should be HTML comment")
	}
}
