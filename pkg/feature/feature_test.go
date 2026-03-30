package feature

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

// setupTestFeatures creates a temporary features directory with the given structure.
// features is a map of feature ID -> README.md content.
func setupTestFeatures(t *testing.T, features map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	for id, content := range features {
		featureDir := filepath.Join(dir, filepath.FromSlash(id))
		if err := os.MkdirAll(featureDir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(featureDir, "README.md"), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func TestDiscover(t *testing.T) {
	featDir := setupTestFeatures(t, map[string]string{
		"alpha":           "# Alpha",
		"beta":            "# Beta",
		"alpha/child-one": "# Child One",
		"alpha/child-two": "# Child Two",
	})

	// Create a reserved _args directory that should be skipped.
	argsDir := filepath.Join(featDir, "alpha", "_args")
	if err := os.MkdirAll(argsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(argsDir, "README.md"), []byte("# Args"), 0o644); err != nil {
		t.Fatal(err)
	}

	features, err := Discover(featDir)
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"alpha", "alpha/child-one", "alpha/child-two", "beta"}
	if len(features) != len(expected) {
		t.Fatalf("got %d features %v, want %d %v", len(features), features, len(expected), expected)
	}
	for i, f := range features {
		if f.ID != expected[i] {
			t.Errorf("feature[%d] = %q, want %q", i, f.ID, expected[i])
		}
	}
}

func TestDiscover_Empty(t *testing.T) {
	dir := t.TempDir()
	features, err := Discover(dir)
	if err != nil {
		t.Fatal(err)
	}
	if len(features) != 0 {
		t.Errorf("got %d features, want 0", len(features))
	}
}

func TestBuildTree(t *testing.T) {
	ids := []string{"alpha", "alpha/child", "alpha/child/deep", "beta", "gamma"}
	roots := BuildTree(ids)

	if len(roots) != 3 {
		t.Fatalf("got %d roots, want 3", len(roots))
	}
	if roots[0].Name != "alpha" {
		t.Errorf("root[0] = %q, want alpha", roots[0].Name)
	}
	if len(roots[0].Children) != 1 || roots[0].Children[0].Name != "child" {
		t.Errorf("alpha should have 1 child named 'child'")
	}
	if len(roots[0].Children[0].Children) != 1 || roots[0].Children[0].Children[0].Name != "deep" {
		t.Errorf("alpha/child should have 1 child named 'deep'")
	}
}

func TestPrintTree(t *testing.T) {
	ids := []string{"alpha", "alpha/child", "beta"}
	roots := BuildTree(ids)
	var sb strings.Builder
	PrintTree(&sb, roots, 0)

	got := sb.String()
	want := "alpha\n\tchild\nbeta\n"
	if got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestValidateSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		slug    string
		wantErr bool
	}{
		{"valid multi-word slug", "task-status-board", false},
		{"valid single word", "cli", false},
		{"valid nested path", "cli/task/claim", false},
		{"valid single character", "a", false},
		{"invalid empty", "", true},
		{"invalid uppercase", "Task", true},
		{"invalid consecutive hyphens", "foo--bar", true},
		{"invalid leading hyphen", "-foo", true},
		{"invalid trailing hyphen", "foo-", true},
		{"invalid trailing slash", "foo/bar/", true},
		{"invalid spaces", "foo bar", true},
		{"invalid underscores", "foo_bar", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := ValidateSlug(tt.slug)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateSlug(%q) error = %v, wantErr %v", tt.slug, err, tt.wantErr)
			}
		})
	}
}

func TestGenerateSlug(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name  string
		title string
		want  string
	}{
		{"lowercase with spaces", "Task Status Board", "task-status-board"},
		{"all caps acronym", "CLI", "cli"},
		{"already hyphenated", "Cross-Repo Sync", "cross-repo-sync"},
		{"parenthetical stripped", "Outstanding Questions (OQ)", "outstanding-questions-oq"},
		{"extra spaces", "  Extra   Spaces  ", "extra-spaces"},
		{"underscores to hyphens", "Hello_World", "hello-world"},
		{"leading trailing hyphens stripped", "---Leading---Trailing---", "leading-trailing"},
		{"empty string", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := GenerateSlug(tt.title)
			if got != tt.want {
				t.Errorf("GenerateSlug(%q) = %q, want %q", tt.title, got, tt.want)
			}
		})
	}
}

func TestGenerateReadme(t *testing.T) {
	t.Parallel()

	requiredSections := []string{
		"## Summary",
		"## Problem",
		"## Behavior",
		"## Acceptance Criteria",
		"## Outstanding Questions",
	}

	tests := []struct {
		name        string
		title       string
		status      string
		description string
		deps        []string
		wantContain []string
		wantAbsent  []string
	}{
		{
			name:        "basic feature with description no deps",
			title:       "My Feature",
			status:      "draft",
			description: "A cool feature.",
			deps:        nil,
			wantContain: []string{
				"# Feature: My Feature",
				"**Status:** draft",
				"A cool feature.",
			},
			wantAbsent: []string{
				"## Dependencies",
			},
		},
		{
			name:        "feature with deps no description",
			title:       "Task Board",
			status:      "approved",
			description: "",
			deps:        []string{"state-store", "cli"},
			wantContain: []string{
				"## Dependencies",
				"- state-store",
				"- cli",
				"TODO: Brief summary of the feature.",
			},
		},
		{
			name:        "feature with description and deps",
			title:       "Test",
			status:      "implemented",
			description: "Test desc.",
			deps:        []string{"dep-a"},
			wantContain: []string{
				"Test desc.",
				"## Dependencies",
				"- dep-a",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got := GenerateReadme(tt.title, tt.status, tt.description, tt.deps)

			for _, want := range tt.wantContain {
				if !strings.Contains(got, want) {
					t.Errorf("GenerateReadme() missing expected string %q\ngot:\n%s", want, got)
				}
			}
			for _, absent := range tt.wantAbsent {
				if strings.Contains(got, absent) {
					t.Errorf("GenerateReadme() should not contain %q\ngot:\n%s", absent, got)
				}
			}
			for _, section := range requiredSections {
				if !strings.Contains(got, section) {
					t.Errorf("GenerateReadme() missing required section %q", section)
				}
			}
			if !strings.HasSuffix(got, "None at this time.\n") {
				t.Error("GenerateReadme() should end with 'None at this time.\\n'")
			}
		})
	}
}

func TestIsValidStatus(t *testing.T) {
	t.Parallel()

	tests := []struct {
		status string
		want   bool
	}{
		{"draft", true},
		{"approved", true},
		{"implemented", true},
		{"Draft", true},
		{"APPROVED", true},
		{"conceptual", false},
		{"not-started", false},
		{"", false},
		{"in_progress", false},
	}

	for _, tt := range tests {
		t.Run("status_"+tt.status, func(t *testing.T) {
			t.Parallel()
			got := IsValidStatus(tt.status)
			if got != tt.want {
				t.Errorf("IsValidStatus(%q) = %v, want %v", tt.status, got, tt.want)
			}
		})
	}
}

func TestParseDependencies_BareIDs(t *testing.T) {
	content := `# Feature: Test

## Dependencies

- claim-and-push
- conflict-resolution

## Outstanding Questions

None at this time.
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseDependencies(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"claim-and-push", "conflict-resolution"}
	if len(deps) != len(expected) {
		t.Fatalf("got %d deps %v, want %d %v", len(deps), deps, len(expected), expected)
	}
	for i, d := range deps {
		if d != expected[i] {
			t.Errorf("dep[%d] = %q, want %q", i, d, expected[i])
		}
	}
}

func TestParseDependencies_MarkdownLinks(t *testing.T) {
	content := `# Feature: GitHub App

## Dependencies

- [API](../api/README.md) — callback endpoint
- [Project Definition](../project-definition/README.md) — config format

## Outstanding Questions
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseDependencies(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}

	expected := []string{"api", "project-definition"}
	sort.Strings(expected)
	if len(deps) != len(expected) {
		t.Fatalf("got %d deps %v, want %d %v", len(deps), deps, len(expected), expected)
	}
	for i, d := range deps {
		if d != expected[i] {
			t.Errorf("dep[%d] = %q, want %q", i, d, expected[i])
		}
	}
}

func TestParseDependencies_NoDependencies(t *testing.T) {
	content := `# Feature: Independent

## Summary

Does its own thing.

## Outstanding Questions

None at this time.
`
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	deps, err := ParseDependencies(filepath.Join(dir, "README.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(deps) != 0 {
		t.Errorf("got %d deps, want 0", len(deps))
	}
}

func TestExtractFeatureID(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"claim-and-push", "claim-and-push"},
		{"cli/task", "cli/task"},
		{"[API](../api/README.md)", "api"},
		{"[API](../api/README.md) — description", "api"},
		{"[Project Definition](../project-definition/README.md) — config format", "project-definition"},
		{"[CLI](../cli/README.md) — entry point", "cli"},
		{"bare-id — some description", "bare-id"},
	}
	for _, tt := range tests {
		got := ExtractFeatureID(tt.input)
		if got != tt.want {
			t.Errorf("ExtractFeatureID(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFeatureIDFromRelativePath(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"../api/README.md", "api"},
		{"../cli/README.md", "cli"},
		{"../project-definition/README.md", "project-definition"},
		{"../../some/nested/README.md", "some/nested"},
		{"./local/README.md", "local"},
	}
	for _, tt := range tests {
		got := FeatureIDFromRelativePath(tt.input)
		if got != tt.want {
			t.Errorf("FeatureIDFromRelativePath(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFindSpecRepoRoot(t *testing.T) {
	// Create a temp dir with spec/features/ structure.
	root := t.TempDir()
	featDir := filepath.Join(root, "spec", "features")
	if err := os.MkdirAll(featDir, 0o755); err != nil {
		t.Fatal(err)
	}

	// Start from the features dir itself.
	got, err := FindSpecRepoRoot(featDir)
	if err != nil {
		t.Fatal(err)
	}
	if got != root {
		t.Errorf("FindSpecRepoRoot(%q) = %q, want %q", featDir, got, root)
	}
}

func TestNew(t *testing.T) {
	featDir := t.TempDir()

	result, err := New(featDir, NewOptions{
		Title:       "My New Feature",
		Status:      "draft",
		Description: "A test feature.",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.FeatureID != "my-new-feature" {
		t.Errorf("FeatureID = %q, want %q", result.FeatureID, "my-new-feature")
	}

	// Verify the README was created.
	content, err := os.ReadFile(result.ReadmePath)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(content), "# Feature: My New Feature") {
		t.Error("README missing title")
	}
	if !strings.Contains(string(content), "**Status:** draft") {
		t.Error("README missing status")
	}
}

func TestNew_SubFeature(t *testing.T) {
	featDir := setupTestFeatures(t, map[string]string{
		"parent": "# Parent\n\n## Summary\n\nParent feature.\n",
	})

	result, err := New(featDir, NewOptions{
		Title:  "Child Feature",
		Parent: "parent",
		Status: "draft",
	})
	if err != nil {
		t.Fatal(err)
	}

	if result.FeatureID != "parent/child-feature" {
		t.Errorf("FeatureID = %q, want %q", result.FeatureID, "parent/child-feature")
	}

	// Verify changed files include parent README.
	if len(result.ChangedFiles) < 2 {
		t.Errorf("expected at least 2 changed files, got %d", len(result.ChangedFiles))
	}
}

func TestParseFeatureStatus(t *testing.T) {
	content := `# Feature: Test

**Status:** Conceptual

## Summary
`
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	status, err := ParseFeatureStatus(path)
	if err != nil {
		t.Fatal(err)
	}
	if status != "Conceptual" {
		t.Errorf("status = %q, want %q", status, "Conceptual")
	}
}

func TestParseSections(t *testing.T) {
	content := `# Feature: Test

## Summary

Some text.

## Dependencies

- dep-a
- dep-b

## Outstanding Questions

None at this time.
`
	dir := t.TempDir()
	path := filepath.Join(dir, "README.md")
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	sections, err := ParseSections(path)
	if err != nil {
		t.Fatal(err)
	}

	if len(sections) != 3 {
		t.Fatalf("got %d sections, want 3", len(sections))
	}
	if sections[0].Title != "Summary" {
		t.Errorf("sections[0].Title = %q, want %q", sections[0].Title, "Summary")
	}
	if sections[1].Title != "Dependencies" {
		t.Errorf("sections[1].Title = %q, want %q", sections[1].Title, "Dependencies")
	}
	if sections[1].Items != 2 {
		t.Errorf("sections[1].Items = %d, want 2", sections[1].Items)
	}
}

func TestFilterFocusedFeatures(t *testing.T) {
	all := []string{"cli", "cli/task", "cli/task/claim", "cli/feature", "api"}

	// Both directions (default).
	got := FilterFocusedFeatures(all, "cli/task", "")
	expected := map[string]bool{"cli": true, "cli/task": true, "cli/task/claim": true}
	for _, f := range got {
		if !expected[f] {
			t.Errorf("unexpected feature in filtered: %q", f)
		}
	}

	// Down only.
	got = FilterFocusedFeatures(all, "cli/task", "down")
	expectedDown := map[string]bool{"cli/task": true, "cli/task/claim": true}
	for _, f := range got {
		if !expectedDown[f] {
			t.Errorf("unexpected feature in filtered (down): %q", f)
		}
	}

	// Up only.
	got = FilterFocusedFeatures(all, "cli/task", "up")
	expectedUp := map[string]bool{"cli": true, "cli/task": true}
	for _, f := range got {
		if !expectedUp[f] {
			t.Errorf("unexpected feature in filtered (up): %q", f)
		}
	}
}

func TestTransitiveDeps(t *testing.T) {
	// Create features: A depends on B, B depends on C.
	featDir := setupTestFeatures(t, map[string]string{
		"a": "# A\n\n## Dependencies\n\n- b\n",
		"b": "# B\n\n## Dependencies\n\n- c\n",
		"c": "# C\n",
	})

	nodes := TransitiveDeps(featDir, "a")
	if len(nodes) != 1 {
		t.Fatalf("got %d nodes, want 1", len(nodes))
	}
	if nodes[0].Path != "b" {
		t.Errorf("nodes[0].Path = %q, want %q", nodes[0].Path, "b")
	}
	children, ok := nodes[0].Children.([]*EnrichedFeature)
	if !ok || len(children) != 1 {
		t.Fatalf("expected 1 child of b, got %v", nodes[0].Children)
	}
	if children[0].Path != "c" {
		t.Errorf("child path = %q, want %q", children[0].Path, "c")
	}
}

func TestFeatureIDs(t *testing.T) {
	features := []Feature{{ID: "a"}, {ID: "b/c"}, {ID: "d"}}
	ids := FeatureIDs(features)
	if len(ids) != 3 || ids[0] != "a" || ids[1] != "b/c" || ids[2] != "d" {
		t.Errorf("FeatureIDs = %v, want [a b/c d]", ids)
	}
}

func TestParseFieldNames(t *testing.T) {
	t.Parallel()

	fields, err := ParseFieldNames("status,deps,oq")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 3 {
		t.Fatalf("got %d fields, want 3", len(fields))
	}

	// Duplicate removal.
	fields, err = ParseFieldNames("status,status")
	if err != nil {
		t.Fatal(err)
	}
	if len(fields) != 1 {
		t.Errorf("got %d fields, want 1 (dedup)", len(fields))
	}

	// Invalid field.
	_, err = ParseFieldNames("invalid")
	if err == nil {
		t.Error("expected error for invalid field name")
	}

	// Empty string.
	fields, err = ParseFieldNames("")
	if err != nil {
		t.Fatal(err)
	}
	if fields != nil {
		t.Errorf("got %v, want nil for empty string", fields)
	}
}
