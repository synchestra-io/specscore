package idea

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateSlug(t *testing.T) {
	cases := []struct {
		slug    string
		wantErr bool
	}{
		{"offline-mode", false},
		{"payment-fraud-signals", false},
		{"a", false},
		{"a1-b2", false},
		{"", true},
		{"Offline-Mode", true},
		{"offline_mode", true},
		{"-offline", true},
		{"offline-", true},
		{"offline--mode", true},
	}
	for _, tc := range cases {
		t.Run(tc.slug, func(t *testing.T) {
			err := ValidateSlug(tc.slug)
			if (err != nil) != tc.wantErr {
				t.Fatalf("ValidateSlug(%q) err=%v, wantErr=%v", tc.slug, err, tc.wantErr)
			}
		})
	}
}

func TestParse_FullIdea(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "offline-mode.md")
	content := `# Idea: Offline Mode

**Status:** Approved
**Date:** 2026-04-10
**Owner:** alice
**Promotes To:** —
**Supersedes:** —
**Related Ideas:** depends_on:payment-rails-audit

## Problem Statement
How Might We let users work without a connection?

## Context
Triggering observation.

## Recommended Direction
Do it.

## Alternatives Considered
Did not do it.

## MVP Scope
Single job.

## Not Doing (and Why)
- Sync conflicts — out of scope

## Key Assumptions to Validate
| Tier | Assumption | How to validate |
|------|------------|-----------------|
| Must-be-true | Users want this | Survey |

## SpecScore Integration
- **New Features this would create:** TBD

## Open Questions
None at this time.
`
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if p.Slug != "offline-mode" {
		t.Errorf("slug = %q", p.Slug)
	}
	if !p.HasTitle || !p.TitleOK || p.TitleName != "Offline Mode" {
		t.Errorf("title parse failed: %+v", p)
	}
	if p.Status() != "Approved" {
		t.Errorf("status = %q", p.Status())
	}
	if got := p.FieldByName["Date"].Value; got != "2026-04-10" {
		t.Errorf("date = %q", got)
	}
	if len(p.PromotesTo()) != 0 {
		t.Errorf("expected empty promotes, got %v", p.PromotesTo())
	}
	if got := p.RelatedIdeas(); len(got) != 1 || got[0] != "depends_on:payment-rails-audit" {
		t.Errorf("related = %v", got)
	}
	// Required sections all present.
	for _, s := range RequiredSections {
		if _, ok := p.SectionByTitle[s]; !ok {
			t.Errorf("missing section %q", s)
		}
	}
	// Not Doing has items.
	if len(p.SectionByTitle["Not Doing (and Why)"].Items) == 0 {
		t.Errorf("not-doing items empty")
	}
	// Key Assumptions table.
	tab := ParseTable(p.SectionByTitle["Key Assumptions to Validate"].Body)
	if tab == nil || len(tab.Rows) == 0 {
		t.Errorf("assumptions table parse failed: %+v", tab)
	}
}

func TestParse_MalformedTitle(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "x.md")
	if err := os.WriteFile(path, []byte("# Something Else\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, err := Parse(path)
	if err != nil {
		t.Fatal(err)
	}
	if !p.HasTitle || p.TitleOK {
		t.Fatalf("expected title without Idea: prefix; got HasTitle=%v TitleOK=%v", p.HasTitle, p.TitleOK)
	}
}

func TestSplitCSVSlugs(t *testing.T) {
	cases := []struct {
		in   string
		want []string
	}{
		{"—", nil},
		{"-", nil},
		{"", nil},
		{"a", []string{"a"}},
		{"a, b", []string{"a", "b"}},
		{"a ,  b  ,c", []string{"a", "b", "c"}},
	}
	for _, tc := range cases {
		got := splitCSVSlugs(tc.in)
		if len(got) != len(tc.want) {
			t.Errorf("splitCSVSlugs(%q) = %v, want %v", tc.in, got, tc.want)
			continue
		}
		for i := range got {
			if got[i] != tc.want[i] {
				t.Errorf("splitCSVSlugs(%q)[%d] = %q, want %q", tc.in, i, got[i], tc.want[i])
			}
		}
	}
}
