package lint

import (
	"os"
	"path/filepath"
	"testing"
)

const validFooterURL = "https://specscore.md/feature-specification"

func TestAdherenceFooter_MissingURL_Reports(t *testing.T) {
	tmp := t.TempDir()
	writeFeatureReadme(t, tmp, "bad", "# Feature: Bad\n\n**Status:** Draft\n\n## Summary\nNo footer here.\n")

	violations := runAdherenceFooterCheck(t, tmp)

	if len(violations) != 1 {
		t.Fatalf("expected 1 violation, got %d: %+v", len(violations), violations)
	}
	v := violations[0]
	if v.Rule != "adherence-footer" {
		t.Errorf("Rule = %q, want %q", v.Rule, "adherence-footer")
	}
	if v.Severity != "error" {
		t.Errorf("Severity = %q, want %q", v.Severity, "error")
	}
	if v.File != filepath.Join("features", "bad", "README.md") {
		t.Errorf("File = %q, want features/bad/README.md", v.File)
	}
}

func TestAdherenceFooter_URLPresent_NoViolation(t *testing.T) {
	tmp := t.TempDir()
	content := "# Feature: Good\n\n**Status:** Draft\n\n## Summary\nHas footer.\n\n---\n*This document follows the " + validFooterURL + "*\n"
	writeFeatureReadme(t, tmp, "good", content)

	violations := runAdherenceFooterCheck(t, tmp)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestAdherenceFooter_TrailingSlashTolerated(t *testing.T) {
	tmp := t.TempDir()
	content := "# Feature: Slash\n\n**Status:** Draft\n\n## Summary\n\n" + validFooterURL + "/\n"
	writeFeatureReadme(t, tmp, "slash", content)

	violations := runAdherenceFooterCheck(t, tmp)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestAdherenceFooter_URLAnywhereInDoc_NoViolation(t *testing.T) {
	tmp := t.TempDir()
	// URL cited in the middle of the body rather than as a footer — still counts.
	content := "# Feature: Middle\n\nSee " + validFooterURL + " for the spec.\n\n## Summary\nBody.\n"
	writeFeatureReadme(t, tmp, "middle", content)

	violations := runAdherenceFooterCheck(t, tmp)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations, got %d: %+v", len(violations), violations)
	}
}

func TestAdherenceFooter_NonFeatureReadmesIgnored(t *testing.T) {
	tmp := t.TempDir()
	// plans/ README missing the URL — must NOT be reported (rule is feature-scoped).
	plansDir := filepath.Join(tmp, "plans", "some-plan")
	if err := os.MkdirAll(plansDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(plansDir, "README.md"),
		[]byte("# Plan: Some Plan\n\nNo footer needed here.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	violations := runAdherenceFooterCheck(t, tmp)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations (plan readme out of scope), got %d: %+v", len(violations), violations)
	}
}

func TestAdherenceFooter_UnderscoreReservedDirsIgnored(t *testing.T) {
	tmp := t.TempDir()
	// A valid feature with the footer.
	writeFeatureReadme(t, tmp, "auth", "# Feature: Auth\n\n**Status:** Draft\n\n"+validFooterURL+"\n")
	// A README inside _tests/ without the footer — must NOT be flagged.
	testsDir := filepath.Join(tmp, "features", "auth", "_tests")
	if err := os.MkdirAll(testsDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(testsDir, "README.md"),
		[]byte("# Tests for Auth\n\nNo footer needed.\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	violations := runAdherenceFooterCheck(t, tmp)
	if len(violations) != 0 {
		t.Errorf("expected 0 violations (_tests ignored), got %d: %+v", len(violations), violations)
	}
}

func TestAdherenceFooter_RegisteredAsKnownRule(t *testing.T) {
	rules := AllRuleNames()
	if !rules["adherence-footer"] {
		t.Error("expected adherence-footer to be a known rule")
	}
}

// writeFeatureReadme writes a feature README under specRoot/features/<slug>/README.md.
func writeFeatureReadme(t *testing.T, specRoot, slug, content string) {
	t.Helper()
	dir := filepath.Join(specRoot, "features", slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "README.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// runAdherenceFooterCheck runs only the adherence-footer checker against specRoot.
func runAdherenceFooterCheck(t *testing.T, specRoot string) []Violation {
	t.Helper()
	c := newAdherenceFooterChecker()
	violations, err := c.check(specRoot)
	if err != nil {
		t.Fatalf("check returned error: %v", err)
	}
	return violations
}
