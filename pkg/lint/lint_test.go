package lint

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFilterBySeverity(t *testing.T) {
	violations := []Violation{
		{Severity: "error", Rule: "r1"},
		{Severity: "warning", Rule: "r2"},
		{Severity: "info", Rule: "r3"},
	}
	filtered := FilterBySeverity(violations, "warning")
	if len(filtered) != 2 {
		t.Errorf("expected 2, got %d", len(filtered))
	}
}

func TestFilterBySeverity_ErrorOnly(t *testing.T) {
	violations := []Violation{
		{Severity: "error", Rule: "readme-exists"},
		{Severity: "warning", Rule: "heading-levels"},
		{Severity: "info", Rule: "diag"},
	}

	errOnly := FilterBySeverity(violations, "error")
	if len(errOnly) != 1 {
		t.Errorf("error filter: expected 1, got %d", len(errOnly))
	}

	all := FilterBySeverity(violations, "info")
	if len(all) != 3 {
		t.Errorf("info filter: expected 3, got %d", len(all))
	}
}

func TestAllRuleNames(t *testing.T) {
	rules := AllRuleNames()
	if len(rules) == 0 {
		t.Error("expected non-empty rule names")
	}
	// Verify known rules are present.
	expected := []string{"readme-exists", "oq-section", "plan-hierarchy", "plan-roi-metadata"}
	for _, name := range expected {
		if !rules[name] {
			t.Errorf("expected rule %q to be present", name)
		}
	}
}

func TestAllRuleNames_ReturnsCopy(t *testing.T) {
	rules := AllRuleNames()
	rules["bogus"] = true
	rules2 := AllRuleNames()
	if rules2["bogus"] {
		t.Error("AllRuleNames should return a copy, not the original map")
	}
}

func TestValidateRuleNames(t *testing.T) {
	if err := ValidateRuleNames(nil); err != nil {
		t.Errorf("unexpected error for nil: %v", err)
	}
	if err := ValidateRuleNames([]string{"readme-exists", "oq-section"}); err != nil {
		t.Errorf("valid rules should not error: %v", err)
	}
	if err := ValidateRuleNames([]string{"nonexistent-rule"}); err == nil {
		t.Error("expected error for unknown rule")
	}
}

func TestViolation_JSONRoundTrip(t *testing.T) {
	v := Violation{
		File:     "features/cli/README.md",
		Line:     42,
		Severity: "error",
		Rule:     "oq-section",
		Message:  "Outstanding Questions section not found",
	}
	data, err := json.Marshal(v)
	if err != nil {
		t.Fatal(err)
	}
	var got Violation
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatal(err)
	}
	if got != v {
		t.Errorf("round-trip mismatch: got %+v, want %+v", got, v)
	}
}

// --- readmeExistsChecker ---

func TestReadmeExists_AllPresent(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "# Root")
	mkdir(t, filepath.Join(root, "child"))
	writeFile(t, filepath.Join(root, "child", "README.md"), "# Child")

	c := newReadmeExistsChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations, got %d: %v", len(v), v)
	}
}

func TestReadmeExists_Missing(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "# Root")
	mkdir(t, filepath.Join(root, "child"))

	c := newReadmeExistsChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(v), v)
	}
	if v[0].Rule != "readme-exists" {
		t.Errorf("expected rule readme-exists, got %s", v[0].Rule)
	}
}

func TestReadmeExists_SkipsHiddenDirs(t *testing.T) {
	root := t.TempDir()
	writeFile(t, filepath.Join(root, "README.md"), "# Root")
	mkdir(t, filepath.Join(root, ".hidden"))

	c := newReadmeExistsChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations for hidden dir, got %d", len(v))
	}
}

// --- oqSectionChecker ---

func TestOQSection_Present(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"features/cli/README.md": "# CLI\n\n## Outstanding Questions\n\n- Should we add X?\n",
	})

	c := newOQSectionChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations, got %d: %v", len(v), v)
	}
}

func TestOQSection_Missing(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"features/cli/README.md": "# CLI\n\n## Summary\n\nSome text.\n",
	})

	c := newOQSectionChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(v), v)
	}
	if v[0].Rule != "oq-section" {
		t.Errorf("expected rule oq-section, got %s", v[0].Rule)
	}
}

func TestOQSection_Empty(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"features/cli/README.md": "# CLI\n\n## Outstanding Questions\n\n## Next Section\n",
	})

	c := newOQSectionChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation (oq-not-empty), got %d: %v", len(v), v)
	}
	if v[0].Rule != "oq-not-empty" {
		t.Errorf("expected rule oq-not-empty, got %s", v[0].Rule)
	}
	if v[0].Severity != "warning" {
		t.Errorf("expected severity warning, got %s", v[0].Severity)
	}
}

// --- indexEntriesChecker ---

func TestIndexEntries_Valid(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"features/cli/README.md":      "# CLI\n\n| Dir | Desc |\n|---|---|\n| [task](task/README.md) | Task mgmt |\n",
		"features/cli/task/README.md": "# Task\n",
	})

	c := newIndexEntriesChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations, got %d: %v", len(v), v)
	}
}

func TestIndexEntries_NonExistentDir(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"features/cli/README.md": "# CLI\n\n| Dir | Desc |\n|---|---|\n| [missing](missing/README.md) | Does not exist |\n",
	})

	c := newIndexEntriesChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(v), v)
	}
	if !strings.Contains(v[0].Message, "missing") {
		t.Errorf("expected message about 'missing', got %s", v[0].Message)
	}
}

// --- linter orchestration ---

func TestLinter_RulesFilter(t *testing.T) {
	opts := Options{
		SpecRoot: t.TempDir(),
		Rules:    []string{"oq-section"},
	}
	l := newLinter(opts)
	if l.isRuleEnabled("readme-exists") {
		t.Error("readme-exists should be disabled when Rules=oq-section")
	}
	if !l.isRuleEnabled("oq-section") {
		t.Error("oq-section should be enabled")
	}
}

func TestLinter_IgnoreFilter(t *testing.T) {
	opts := Options{
		SpecRoot: t.TempDir(),
		Ignore:   []string{"code-annotations"},
	}
	l := newLinter(opts)
	if l.isRuleEnabled("code-annotations") {
		t.Error("code-annotations should be disabled when Ignore=code-annotations")
	}
	if !l.isRuleEnabled("readme-exists") {
		t.Error("readme-exists should be enabled")
	}
}

// --- planHierarchyChecker ---

func TestPlanHierarchyChecker_RoadmapWithSteps(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/roadmap-a/README.md":            "# Roadmap A\n\n## Steps\n\n- Step 1\n- Step 2\n",
		"plans/roadmap-a/child-plan/README.md": "# Child Plan\n\n## Steps\n\n- Do something\n",
	})

	c := newPlanHierarchyChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}

	var stepsViolations []Violation
	for _, viol := range v {
		if strings.Contains(viol.Message, "Steps") {
			stepsViolations = append(stepsViolations, viol)
		}
	}
	if len(stepsViolations) != 1 {
		t.Fatalf("expected 1 Steps violation, got %d: %v", len(stepsViolations), v)
	}
	if stepsViolations[0].Rule != "plan-hierarchy" {
		t.Errorf("expected rule plan-hierarchy, got %s", stepsViolations[0].Rule)
	}
}

func TestPlanHierarchyChecker_ThreeLevelNesting(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/roadmap-a/README.md":                       "# Roadmap A\n\n## Child Plans\n\n- child-plan\n",
		"plans/roadmap-a/child-plan/README.md":            "# Child Plan\n\n## Child Plans\n\n- grandchild\n",
		"plans/roadmap-a/child-plan/grandchild/README.md": "# Grandchild\n\n## Steps\n\n- Do something\n",
	})

	c := newPlanHierarchyChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}

	var nestingViolations []Violation
	for _, viol := range v {
		if strings.Contains(viol.Message, "nesting") || strings.Contains(viol.Message, "depth") {
			nestingViolations = append(nestingViolations, viol)
		}
	}
	if len(nestingViolations) == 0 {
		t.Fatalf("expected nesting violation, got none; all violations: %v", v)
	}
}

func TestPlanHierarchyChecker_ValidHierarchy(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/roadmap-a/README.md":            "# Roadmap A\n\n## Child Plans\n\n- child-plan\n",
		"plans/roadmap-a/child-plan/README.md": "# Child Plan\n\n## Steps\n\n- Do something\n",
		"plans/standalone/README.md":           "# Standalone Plan\n\n## Steps\n\n- Step 1\n",
	})

	c := newPlanHierarchyChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations for valid hierarchy, got %d: %v", len(v), v)
	}
}

// --- planROIChecker ---

func TestPlanROIChecker_InvalidEffort(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/my-plan/README.md": "# My Plan\n\n**Effort:** huge\n**Impact:** high\n\n## Steps\n\n- Step 1\n",
	})

	c := newPlanROIChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 1 {
		t.Fatalf("expected 1 violation, got %d: %v", len(v), v)
	}
	if !strings.Contains(v[0].Message, "Effort") {
		t.Errorf("expected violation to mention Effort, got: %s", v[0].Message)
	}
}

func TestPlanROIChecker_ValidMetadata(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/my-plan/README.md": "# My Plan\n\n**Effort:** M\n**Impact:** high\n\n## Steps\n\n- Step 1\n",
	})

	c := newPlanROIChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations for valid metadata, got %d: %v", len(v), v)
	}
}

func TestPlanROIChecker_NoMetadata(t *testing.T) {
	root := setupSpecTree(t, map[string]string{
		"plans/my-plan/README.md": "# My Plan\n\n## Steps\n\n- Step 1\n",
	})

	c := newPlanROIChecker()
	v, err := c.check(root)
	if err != nil {
		t.Fatal(err)
	}
	if len(v) != 0 {
		t.Errorf("expected 0 violations when metadata absent, got %d: %v", len(v), v)
	}
}

// --- Lint integration ---

func TestLint_InvalidSpecRoot(t *testing.T) {
	_, err := Lint(Options{SpecRoot: "/nonexistent/path"})
	if err == nil {
		t.Error("expected error for nonexistent spec root")
	}
}

func TestLint_EmptyDir(t *testing.T) {
	root := t.TempDir()
	violations, err := Lint(Options{SpecRoot: root, Severity: "info"})
	if err != nil {
		t.Fatal(err)
	}
	// An empty dir with no README.md should produce a readme-exists violation.
	found := false
	for _, v := range violations {
		if v.Rule == "readme-exists" {
			found = true
		}
	}
	if !found {
		t.Error("expected readme-exists violation for empty dir")
	}
}

// --- helpers ---

func setupSpecTree(t *testing.T, files map[string]string) string {
	t.Helper()
	root := t.TempDir()
	for relPath, content := range files {
		fullPath := filepath.Join(root, relPath)
		mkdir(t, filepath.Dir(fullPath))
		writeFile(t, fullPath, content)
	}
	return root
}

func mkdir(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatal(err)
	}
}

func writeFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}
