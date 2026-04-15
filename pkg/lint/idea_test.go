package lint

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// validIdeaBody returns a full, lint-clean idea markdown body.
// title and status are parameterized.
func validIdeaBody(title, status string, extraFields map[string]string) string {
	fields := map[string]string{
		"Status":        status,
		"Date":          "2026-04-10",
		"Owner":         "alice",
		"Promotes To":   "—",
		"Supersedes":    "—",
		"Related Ideas": "—",
	}
	for k, v := range extraFields {
		fields[k] = v
	}
	order := []string{"Status", "Date", "Owner", "Promotes To", "Supersedes", "Related Ideas"}
	if _, ok := extraFields["Archive Reason"]; ok {
		order = append(order, "Archive Reason")
	}
	var header strings.Builder
	header.WriteString("# Idea: " + title + "\n\n")
	for _, k := range order {
		header.WriteString("**" + k + ":** " + fields[k] + "\n")
	}
	header.WriteString("\n")
	header.WriteString(`## Problem Statement
How Might We ship faster?

## Context
Triggering observation.

## Recommended Direction
Do it.

## Alternatives Considered
Nope.

## MVP Scope
Small.

## Not Doing (and Why)
- Thing — reason

## Key Assumptions to Validate
| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Users want this | Survey |

## SpecScore Integration
- **New Features this would create:** TBD

## Open Questions
None at this time.
`)
	return header.String()
}

// writeSpec creates a fake spec repo structure under dir and returns specRoot.
func writeSpec(t *testing.T, files map[string]string) string {
	t.Helper()
	dir := t.TempDir()
	specRoot := filepath.Join(dir, "spec")
	for rel, content := range files {
		path := filepath.Join(specRoot, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return specRoot
}

// hasRule returns whether vs contains a violation with the given rule name.
func hasRule(vs []Violation, rule string) bool {
	for _, v := range vs {
		if v.Rule == rule {
			return true
		}
	}
	return false
}

const activeIndex = `# SpecScore Ideas

## Index

| Idea | Status | Date | Owner | Promotes To |
|------|--------|------|-------|-------------|

## Outstanding Questions

None at this time.
`

const archivedIndex = `# Archived Ideas

## Outstanding Questions

None at this time.
`

func TestCheckIdeas_CleanIdea(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("offline-mode"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/offline-mode.md":    validIdeaBody("Offline Mode", "Approved", nil),
	})
	vs, err := CheckIdeas(specRoot, false)
	if err != nil {
		t.Fatal(err)
	}
	if len(vs) > 0 {
		t.Fatalf("expected 0 violations, got %d: %+v", len(vs), vs)
	}
}

func activeIndexWith(slugs ...string) string {
	var b strings.Builder
	b.WriteString("# SpecScore Ideas\n\n## Index\n\n| Idea | Status | Date | Owner | Promotes To |\n|------|--------|------|-------|-------------|\n")
	for _, s := range slugs {
		b.WriteString("| [" + s + "](" + s + ".md) | Approved | 2026-04-10 | alice | — |\n")
	}
	b.WriteString("\n## Outstanding Questions\n\nNone at this time.\n")
	return b.String()
}

func TestCheckIdeas_InvalidSlug(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndex,
		"ideas/archived/README.md": archivedIndex,
		"ideas/Bad_Slug.md":        validIdeaBody("Bad Slug", "Approved", nil),
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-slug-format") {
		t.Errorf("expected idea-slug-format violation, got: %+v", vs)
	}
}

func TestCheckIdeas_IdeaDirectoryRejected(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":              activeIndex,
		"ideas/archived/README.md":     archivedIndex,
		"ideas/offline-mode/README.md": "# Idea: Offline Mode\n",
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-single-file") {
		t.Errorf("expected idea-single-file violation")
	}
}

func TestCheckIdeas_MissingTitle(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("bad"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/bad.md":             "**Status:** Draft\n",
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-title-format") {
		t.Errorf("expected idea-title-format violation: %+v", vs)
	}
}

func TestCheckIdeas_WrongTitlePrefix(t *testing.T) {
	body := "# Feature: Something\n\n**Status:** Draft\n**Date:** 2026-04-10\n**Owner:** alice\n**Promotes To:** —\n**Supersedes:** —\n**Related Ideas:** —\n"
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("bad"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/bad.md":             body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-title-format") {
		t.Errorf("expected idea-title-format violation")
	}
}

func TestCheckIdeas_MissingHeaderField(t *testing.T) {
	body := strings.Replace(validIdeaBody("X", "Draft", nil), "**Owner:** alice\n", "", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-header-fields") {
		t.Errorf("expected idea-header-fields violation")
	}
}

func TestCheckIdeas_InvalidStatus(t *testing.T) {
	body := strings.Replace(validIdeaBody("X", "Draft", nil), "**Status:** Draft", "**Status:** WIP", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-status-values") {
		t.Errorf("expected idea-status-values violation")
	}
}

func TestCheckIdeas_SpecifiedWithoutPromotion(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               validIdeaBody("X", "Specified", nil),
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-specified-requires-promotion") {
		t.Errorf("expected idea-specified-requires-promotion violation")
	}
}

func TestCheckIdeas_ArchivedOutsideArchivedDir(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               validIdeaBody("X", "Archived", map[string]string{"Archive Reason": "abandoned"}),
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-archived-location") {
		t.Errorf("expected idea-archived-location violation")
	}
}

func TestCheckIdeas_ArchivedMissingReason(t *testing.T) {
	body := validIdeaBody("X", "Archived", nil)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndex,
		"ideas/archived/README.md": archivedIndex,
		"ideas/archived/x.md":      body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-archive-reason") {
		t.Errorf("expected idea-archive-reason violation: %+v", vs)
	}
}

func TestCheckIdeas_SupersedesNonArchived(t *testing.T) {
	x := validIdeaBody("X", "Approved", nil)
	y := strings.Replace(validIdeaBody("Y", "Approved", nil), "**Supersedes:** —", "**Supersedes:** x", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x", "y"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               x,
		"ideas/y.md":               y,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-supersedes-target-archived") {
		t.Errorf("expected idea-supersedes-target-archived violation")
	}
}

func TestCheckIdeas_InvalidRelatedIdeasSyntax(t *testing.T) {
	body := strings.Replace(validIdeaBody("X", "Approved", nil), "**Related Ideas:** —", "**Related Ideas:** bogus_rel:other", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-related-ideas-format") {
		t.Errorf("expected idea-related-ideas-format violation")
	}
}

func TestCheckIdeas_BrokenRelatedSlug(t *testing.T) {
	body := strings.Replace(validIdeaBody("X", "Approved", nil), "**Related Ideas:** —", "**Related Ideas:** depends_on:ghost", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-related-ideas-target-exists") {
		t.Errorf("expected idea-related-ideas-target-exists violation")
	}
}

func TestCheckIdeas_EmptyNotDoing(t *testing.T) {
	body := validIdeaBody("X", "Approved", nil)
	body = strings.Replace(body, "- Thing — reason\n", "", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-not-doing-non-empty") {
		t.Errorf("expected idea-not-doing-non-empty violation")
	}
}

func TestCheckIdeas_MissingMustBeTrue(t *testing.T) {
	body := validIdeaBody("X", "Approved", nil)
	body = strings.Replace(body, "| Must-be-true | Users want this | Survey |\n", "", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-must-be-true-present") {
		t.Errorf("expected idea-must-be-true-present violation")
	}
}

func TestCheckIdeas_MissingHMWFraming(t *testing.T) {
	body := validIdeaBody("X", "Approved", nil)
	body = strings.Replace(body, "How Might We ship faster?", "Some other framing.", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-hmw-framing") {
		t.Errorf("expected idea-hmw-framing warning")
	}
}

func TestCheckIdeas_MissingRequiredSection(t *testing.T) {
	body := validIdeaBody("X", "Approved", nil)
	body = strings.Replace(body, "## MVP Scope\nSmall.\n\n", "", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-required-sections") {
		t.Errorf("expected idea-required-sections violation")
	}
}

func TestCheckIdeas_SyncDriftDetectedAndFixed(t *testing.T) {
	feature := `# Feature: Offline Sync

**Status:** Draft
**Source Ideas:** offline-mode

## Summary
Stuff.

## Outstanding Questions

None at this time.
`
	body := validIdeaBody("Offline Mode", "Approved", nil)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":                 activeIndexWith("offline-mode"),
		"ideas/archived/README.md":        archivedIndex,
		"ideas/offline-mode.md":           body,
		"features/offline-sync/README.md": feature,
	})

	// Without --fix: drift reported.
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-sync-lint-strict") {
		t.Fatalf("expected idea-sync-lint-strict: %+v", vs)
	}

	// With --fix: drift repaired.
	vs2, _ := CheckIdeas(specRoot, true)
	if hasRule(vs2, "idea-sync-lint-strict") {
		t.Errorf("fix did not repair drift: %+v", vs2)
	}

	// Verify idea file was rewritten.
	data, err := os.ReadFile(filepath.Join(specRoot, "ideas", "offline-mode.md"))
	if err != nil {
		t.Fatal(err)
	}
	s := string(data)
	if !strings.Contains(s, "**Status:** Specified") {
		t.Errorf("status not updated: %s", s)
	}
	if !strings.Contains(s, "**Promotes To:** offline-sync") {
		t.Errorf("promotes-to not updated: %s", s)
	}

	// Subsequent lint passes.
	vs3, _ := CheckIdeas(specRoot, false)
	if hasRule(vs3, "idea-sync-lint-strict") {
		t.Errorf("subsequent lint still reports drift: %+v", vs3)
	}
}

func TestCheckIdeas_FeatureReferencesDraftIdea(t *testing.T) {
	feature := `# Feature: Offline Sync

**Status:** Draft
**Source Ideas:** offline-mode

## Outstanding Questions

None at this time.
`
	body := validIdeaBody("Offline Mode", "Draft", nil)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":                 activeIndexWith("offline-mode"),
		"ideas/archived/README.md":        archivedIndex,
		"ideas/offline-mode.md":           body,
		"features/offline-sync/README.md": feature,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-feature-cross-reference") {
		t.Errorf("expected idea-feature-cross-reference violation")
	}
}

func TestCheckIdeas_UnlistedIdeaInIndex(t *testing.T) {
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndex, // empty index
		"ideas/archived/README.md": archivedIndex,
		"ideas/offline-mode.md":    validIdeaBody("Offline Mode", "Approved", nil),
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-index-completeness") {
		t.Errorf("expected idea-index-completeness violation: %+v", vs)
	}

	// With --fix, index is regenerated.
	vs2, _ := CheckIdeas(specRoot, true)
	if hasRule(vs2, "idea-index-completeness") {
		t.Errorf("--fix did not regenerate index: %+v", vs2)
	}
	data, _ := os.ReadFile(filepath.Join(specRoot, "ideas", "README.md"))
	if !strings.Contains(string(data), "offline-mode") {
		t.Errorf("index not updated: %s", string(data))
	}
}

func TestCheckIdeas_ArchivedIndexOutOfOrderAndFixed(t *testing.T) {
	older := validIdeaBody("Older", "Archived", map[string]string{"Archive Reason": "pivoted", "Date": "2024-11-02"})
	newer := validIdeaBody("Newer", "Archived", map[string]string{"Archive Reason": "pivoted", "Date": "2025-03-10"})
	// Date inside body: need to set Date via extraFields; validIdeaBody uses fixed default.
	// To override we rewrite the Date line.
	older = strings.Replace(older, "**Date:** 2026-04-10", "**Date:** 2024-11-02", 1)
	newer = strings.Replace(newer, "**Date:** 2026-04-10", "**Date:** 2025-03-10", 1)

	// Index with wrong order (newer first).
	badArchIndex := `# Archived Ideas

- 2025-03-10 — [newer](newer.md) — pivoted
- 2024-11-02 — [older](older.md) — pivoted

## Outstanding Questions

None at this time.
`
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndex,
		"ideas/archived/README.md": badArchIndex,
		"ideas/archived/older.md":  older,
		"ideas/archived/newer.md":  newer,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-archived-index-chronological") {
		t.Errorf("expected idea-archived-index-chronological: %+v", vs)
	}

	vs2, _ := CheckIdeas(specRoot, true)
	if hasRule(vs2, "idea-archived-index-chronological") {
		t.Errorf("--fix did not reorder: %+v", vs2)
	}
	data, _ := os.ReadFile(filepath.Join(specRoot, "ideas", "archived", "README.md"))
	body := string(data)
	iOlder := strings.Index(body, "older.md")
	iNewer := strings.Index(body, "newer.md")
	if iOlder == -1 || iNewer == -1 || iOlder > iNewer {
		t.Errorf("entries not reordered: %s", body)
	}
}

func TestCheckIdeas_IdeaWithIdField(t *testing.T) {
	body := validIdeaBody("X", "Draft", nil)
	body = strings.Replace(body, "**Status:** Draft\n", "**Status:** Draft\n**Id:** x\n", 1)
	specRoot := writeSpec(t, map[string]string{
		"ideas/README.md":          activeIndexWith("x"),
		"ideas/archived/README.md": archivedIndex,
		"ideas/x.md":               body,
	})
	vs, _ := CheckIdeas(specRoot, false)
	if !hasRule(vs, "idea-id-is-slug") {
		t.Errorf("expected idea-id-is-slug violation")
	}
}
