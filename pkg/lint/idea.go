package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/synchestra-io/specscore/pkg/idea"
)

// ideaChecker is a dispatch checker that runs all idea-* rules in one pass
// per file so parsing is shared.
type ideaChecker struct {
	fix bool
}

func newIdeaChecker() *ideaChecker {
	return &ideaChecker{}
}

// Rule names this checker emits. Keep aligned with allRuleNames.
var ideaRuleNames = []string{
	"idea-location",
	"idea-slug-format",
	"idea-single-file",
	"idea-title-format",
	"idea-header-fields",
	"idea-id-is-slug",
	"idea-required-sections",
	"idea-not-doing-non-empty",
	"idea-must-be-true-present",
	"idea-hmw-framing",
	"idea-status-values",
	"idea-specified-requires-promotion",
	"idea-archived-location",
	"idea-archive-reason",
	"idea-supersedes-target-archived",
	"idea-related-ideas-format",
	"idea-related-ideas-target-exists",
	"idea-sync-lint-strict",
	"idea-feature-cross-reference",
	"idea-index-completeness",
	"idea-archived-index-chronological",
}

func (c *ideaChecker) name() string     { return "idea-location" }
func (c *ideaChecker) severity() string { return "error" }

func (c *ideaChecker) check(specRoot string) ([]Violation, error) {
	return CheckIdeas(specRoot, c.fix)
}

// CheckIdeas runs all idea lint rules. If fix is true, auto-fixable
// violations are repaired on disk and omitted from the returned slice.
func CheckIdeas(specRoot string, fix bool) ([]Violation, error) {
	var violations []Violation

	ideasDir := filepath.Join(specRoot, "ideas")
	if info, err := os.Stat(ideasDir); err != nil || !info.IsDir() {
		return nil, nil
	}

	// 1. Detect directories (REQ: single-file).
	dirs, err := idea.FindIdeaDirectories(specRoot)
	if err != nil {
		return nil, err
	}
	for _, d := range dirs {
		rel, _ := filepath.Rel(specRoot, d)
		violations = append(violations, Violation{
			File:     rel,
			Line:     0,
			Severity: "error",
			Rule:     "idea-single-file",
			Message:  fmt.Sprintf("ideas must be single markdown files; directory found at %s", rel),
		})
	}

	// 2. Check for misplaced idea files (idea-location).
	// We report files inside archived/ subdirs of archived/ (deeper nesting) or
	// in docs/ideas/ but only if inside specRoot (we stay within spec tree).
	misplaced, err := findMisplacedIdeaFiles(specRoot)
	if err != nil {
		return nil, err
	}
	for _, f := range misplaced {
		rel, _ := filepath.Rel(specRoot, f)
		violations = append(violations, Violation{
			File:     rel,
			Line:     0,
			Severity: "error",
			Rule:     "idea-location",
			Message:  fmt.Sprintf("idea files must live at spec/ideas/<slug>.md or spec/ideas/archived/<slug>.md; got %s", rel),
		})
	}

	// 3. Discover and parse idea files.
	discovered, err := idea.Discover(specRoot)
	if err != nil {
		return nil, err
	}

	parsed := make(map[string]*idea.Idea)
	archivedMap := make(map[string]bool)
	for _, d := range discovered {
		archivedMap[d.Slug] = d.Archived
		p, err := idea.Parse(d.Path)
		if err != nil {
			rel, _ := filepath.Rel(specRoot, d.Path)
			violations = append(violations, Violation{
				File:     rel,
				Severity: "error",
				Rule:     "idea-location",
				Message:  fmt.Sprintf("cannot read idea file: %v", err),
			})
			continue
		}
		parsed[d.Slug] = p
	}

	// 4. Feature -> ideas map.
	featureIdeas, err := idea.FeatureSourceIdeas(specRoot)
	if err != nil {
		return nil, err
	}

	// 5. Per-idea rules.
	for _, d := range discovered {
		p, ok := parsed[d.Slug]
		if !ok {
			continue
		}
		rel, _ := filepath.Rel(specRoot, d.Path)
		violations = append(violations, ideaFileRules(p, rel, d.Archived, parsed, archivedMap)...)
	}

	// 6. Cross-artifact: sync-lint-strict + feature-cross-reference.
	syncVs, fixedSync := ideaSyncRules(specRoot, parsed, archivedMap, featureIdeas, fix)
	violations = append(violations, syncVs...)

	// 7. Index completeness + chronological order.
	idxVs, _ := ideaIndexRules(specRoot, discovered, parsed, fix)
	violations = append(violations, idxVs...)

	_ = fixedSync
	return violations, nil
}

// findMisplacedIdeaFiles locates .md files under `spec/ideas/` but outside
// the allowed positions (top-level or archived/).
func findMisplacedIdeaFiles(specRoot string) ([]string, error) {
	ideasDir := filepath.Join(specRoot, "ideas")
	archivedDir := filepath.Join(ideasDir, "archived")
	var out []string
	err := filepath.Walk(ideasDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, ".md") {
			return nil
		}
		dir := filepath.Dir(path)
		if dir == ideasDir || dir == archivedDir {
			return nil
		}
		out = append(out, path)
		return nil
	})
	return out, err
}

// ideaFileRules returns violations for a single parsed idea file.
func ideaFileRules(p *idea.Idea, relPath string, archived bool, all map[string]*idea.Idea, archivedMap map[string]bool) []Violation {
	var vs []Violation

	// idea-slug-format
	if err := idea.ValidateSlug(p.Slug); err != nil {
		vs = append(vs, Violation{
			File: relPath, Line: 0, Severity: "error",
			Rule: "idea-slug-format", Message: err.Error(),
		})
	}

	// idea-title-format
	if !p.HasTitle {
		vs = append(vs, Violation{
			File: relPath, Line: 1, Severity: "error",
			Rule: "idea-title-format", Message: "missing title (expected `# Idea: <Name>`)",
		})
	} else if !p.TitleOK {
		vs = append(vs, Violation{
			File: relPath, Line: p.TitleLine, Severity: "error",
			Rule: "idea-title-format", Message: "title must use `# Idea: <Name>` format",
		})
	} else if strings.TrimSpace(p.TitleName) == "" {
		vs = append(vs, Violation{
			File: relPath, Line: p.TitleLine, Severity: "error",
			Rule: "idea-title-format", Message: "title must use `# Idea: <Name>` format (missing name)",
		})
	}

	// idea-id-is-slug — no **Id:** line.
	if _, ok := p.FieldByName["Id"]; ok {
		vs = append(vs, Violation{
			File: relPath, Line: p.FieldByName["Id"].Line, Severity: "error",
			Rule: "idea-id-is-slug", Message: "ideas must not carry an **Id:** field; filename slug is authoritative",
		})
	}

	// idea-header-fields: required fields present, in order.
	seen := map[string]bool{}
	order := []string{}
	for _, f := range p.Fields {
		if !seen[f.Name] {
			seen[f.Name] = true
			order = append(order, f.Name)
		}
	}
	for _, req := range idea.RequiredHeaderFields {
		if !seen[req] {
			vs = append(vs, Violation{
				File: relPath, Line: 0, Severity: "error",
				Rule: "idea-header-fields", Message: fmt.Sprintf("missing required header field **%s:**", req),
			})
		}
	}
	// Ordering check (ignore unknown fields).
	var filtered []string
	for _, name := range order {
		for _, req := range idea.RequiredHeaderFields {
			if name == req {
				filtered = append(filtered, name)
				break
			}
		}
	}
	// Check that filtered is a prefix of RequiredHeaderFields (subset keeping order).
	j := 0
	orderedOK := true
	for _, name := range filtered {
		for j < len(idea.RequiredHeaderFields) && idea.RequiredHeaderFields[j] != name {
			j++
		}
		if j >= len(idea.RequiredHeaderFields) {
			orderedOK = false
			break
		}
		j++
	}
	if !orderedOK {
		vs = append(vs, Violation{
			File: relPath, Line: 0, Severity: "error",
			Rule: "idea-header-fields", Message: "required header fields are not in the canonical order (Status, Date, Owner, Promotes To, Supersedes, Related Ideas)",
		})
	}

	// idea-status-values
	status := p.Status()
	if status != "" && !idea.ValidStatuses[status] {
		vs = append(vs, Violation{
			File: relPath, Line: p.FieldByName["Status"].Line, Severity: "error",
			Rule:    "idea-status-values",
			Message: fmt.Sprintf("Status %q is not one of Draft, Under Review, Approved, Specified, Archived", status),
		})
	}

	// idea-specified-requires-promotion
	if status == "Specified" {
		if len(p.PromotesTo()) == 0 {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Status"].Line, Severity: "error",
				Rule:    "idea-specified-requires-promotion",
				Message: "Status: Specified requires a non-empty **Promotes To:** list",
			})
		}
	}

	// idea-archived-location
	if status == "Archived" && !archived {
		vs = append(vs, Violation{
			File: relPath, Line: p.FieldByName["Status"].Line, Severity: "error",
			Rule:    "idea-archived-location",
			Message: "Status: Archived requires file to live under spec/ideas/archived/",
		})
	}
	if archived && status != "" && status != "Archived" {
		vs = append(vs, Violation{
			File: relPath, Line: p.FieldByName["Status"].Line, Severity: "error",
			Rule:    "idea-archived-location",
			Message: fmt.Sprintf("files under spec/ideas/archived/ must have Status: Archived; got %q", status),
		})
	}

	// idea-archive-reason
	if status == "Archived" {
		if reason := p.ArchiveReason(); reason == "" || reason == "—" || reason == "-" {
			line := 0
			if f, ok := p.FieldByName["Archive Reason"]; ok {
				line = f.Line
			}
			vs = append(vs, Violation{
				File: relPath, Line: line, Severity: "error",
				Rule:    "idea-archive-reason",
				Message: "Status: Archived requires a non-empty **Archive Reason:**",
			})
		}
	}

	// idea-supersedes-target-archived
	for _, sup := range p.Supersedes() {
		tgt, ok := all[sup]
		if !ok {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Supersedes"].Line, Severity: "error",
				Rule:    "idea-supersedes-target-archived",
				Message: fmt.Sprintf("supersedes target %q not found", sup),
			})
			continue
		}
		if !archivedMap[sup] || tgt.Status() != "Archived" {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Supersedes"].Line, Severity: "error",
				Rule:    "idea-supersedes-target-archived",
				Message: fmt.Sprintf("supersedes target %q must be Archived and live under spec/ideas/archived/", sup),
			})
		}
	}

	// idea-related-ideas-format + target-exists
	for _, entry := range p.RelatedIdeas() {
		rel, slug, ok := strings.Cut(entry, ":")
		rel = strings.TrimSpace(rel)
		slug = strings.TrimSpace(slug)
		if !ok || rel == "" || slug == "" {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Related Ideas"].Line, Severity: "error",
				Rule:    "idea-related-ideas-format",
				Message: fmt.Sprintf("related-ideas entry %q must be `<relationship>:<slug>`", entry),
			})
			continue
		}
		if !idea.ValidRelationships[rel] {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Related Ideas"].Line, Severity: "error",
				Rule:    "idea-related-ideas-format",
				Message: fmt.Sprintf("unknown relationship %q (valid: depends_on, alternative_to, extends, conflicts_with)", rel),
			})
			continue
		}
		if _, ok := all[slug]; !ok {
			vs = append(vs, Violation{
				File: relPath, Line: p.FieldByName["Related Ideas"].Line, Severity: "error",
				Rule:    "idea-related-ideas-target-exists",
				Message: fmt.Sprintf("related-ideas target %q does not resolve to an Idea", slug),
			})
		}
	}

	// idea-required-sections (in order)
	presentOrder := []string{}
	for _, s := range p.Sections {
		presentOrder = append(presentOrder, s.Title)
	}
	var missing []string
	for _, req := range idea.RequiredSections {
		if _, ok := p.SectionByTitle[req]; !ok {
			missing = append(missing, req)
		}
	}
	if len(missing) > 0 {
		vs = append(vs, Violation{
			File: relPath, Line: 0, Severity: "error",
			Rule:    "idea-required-sections",
			Message: fmt.Sprintf("missing required section(s): %s", strings.Join(missing, ", ")),
		})
	} else {
		// Check order.
		idx := 0
		ordered := true
		for _, name := range presentOrder {
			match := -1
			for k, req := range idea.RequiredSections {
				if req == name {
					match = k
					break
				}
			}
			if match == -1 {
				continue
			}
			if match < idx {
				ordered = false
				break
			}
			idx = match
		}
		if !ordered {
			vs = append(vs, Violation{
				File: relPath, Line: 0, Severity: "error",
				Rule:    "idea-required-sections",
				Message: "required sections present but not in canonical order",
			})
		}
	}

	// idea-not-doing-non-empty
	if s, ok := p.SectionByTitle["Not Doing (and Why)"]; ok {
		if len(s.Items) == 0 {
			vs = append(vs, Violation{
				File: relPath, Line: s.StartLine, Severity: "error",
				Rule:    "idea-not-doing-non-empty",
				Message: "Not Doing (and Why) must contain at least one explicit exclusion",
			})
		}
	}

	// idea-must-be-true-present
	if s, ok := p.SectionByTitle["Key Assumptions to Validate"]; ok {
		tab := idea.ParseTable(s.Body)
		found := false
		if tab != nil {
			for _, row := range tab.Rows {
				if len(row) > 0 && strings.EqualFold(row[0], "Must-be-true") {
					// require the assumption column to be non-empty
					if len(row) > 1 && strings.TrimSpace(row[1]) != "" && strings.TrimSpace(row[1]) != "…" {
						found = true
						break
					}
				}
			}
		}
		if !found {
			vs = append(vs, Violation{
				File: relPath, Line: s.StartLine, Severity: "error",
				Rule:    "idea-must-be-true-present",
				Message: "Key Assumptions to Validate must list at least one Must-be-true assumption with content",
			})
		}
	}

	// idea-hmw-framing (warn)
	if s, ok := p.SectionByTitle["Problem Statement"]; ok {
		if !hmwRe.MatchString(s.Body) {
			vs = append(vs, Violation{
				File: relPath, Line: s.StartLine, Severity: "warning",
				Rule:    "idea-hmw-framing",
				Message: "Problem Statement should contain a 'How Might We' framing",
			})
		}
	}

	return vs
}

var hmwRe = regexp.MustCompile(`(?i)how might we`)

// ideaSyncRules enforces that Feature **Source Ideas:** entries agree with
// each idea's **Promotes To:** / **Status:**. When fix is true, rewrites the
// idea header in place and omits the fixed violations.
func ideaSyncRules(specRoot string, parsed map[string]*idea.Idea, archivedMap map[string]bool, featureIdeas map[string][]string, fix bool) ([]Violation, bool) {
	var vs []Violation
	fixed := false

	// Build reverse index: idea slug -> []feature slug.
	reverse := make(map[string][]string)
	for feature, ideas := range featureIdeas {
		for _, slug := range ideas {
			reverse[slug] = append(reverse[slug], feature)
		}
	}
	for _, list := range reverse {
		sort.Strings(list)
	}

	// idea-feature-cross-reference: each feature->idea reference must resolve
	// and the idea must be Approved or Specified (not Draft/Under Review/Archived).
	for feature, ideas := range featureIdeas {
		for _, slug := range ideas {
			p, ok := parsed[slug]
			if !ok {
				rel := filepath.Join("features", feature, "README.md")
				vs = append(vs, Violation{
					File: rel, Line: 0, Severity: "error",
					Rule:    "idea-feature-cross-reference",
					Message: fmt.Sprintf("Source Ideas references %q which does not resolve to an Idea", slug),
				})
				continue
			}
			st := p.Status()
			if st != "Approved" && st != "Specified" {
				rel := filepath.Join("features", feature, "README.md")
				vs = append(vs, Violation{
					File: rel, Line: 0, Severity: "error",
					Rule:    "idea-feature-cross-reference",
					Message: fmt.Sprintf("Source Ideas references idea %q with Status %q (must be Approved or Specified)", slug, st),
				})
			}
		}
	}

	// idea-sync-lint-strict per idea.
	for slug, p := range parsed {
		refs := reverse[slug]
		expectedPromotes := append([]string{}, refs...)
		sort.Strings(expectedPromotes)
		actualPromotes := p.PromotesTo()
		sort.Strings(actualPromotes)

		var expectedStatus string
		if len(refs) > 0 {
			expectedStatus = "Specified"
		} else {
			// Not Specified. But the idea may legitimately be Draft/Approved/etc.
			// Drift only if current status == Specified.
			if p.Status() == "Specified" {
				expectedStatus = "Approved" // drop back
			} else {
				expectedStatus = "" // no change expected
			}
		}

		promotesDrift := !stringSliceEq(expectedPromotes, actualPromotes)
		statusDrift := expectedStatus != "" && p.Status() != expectedStatus

		if !promotesDrift && !statusDrift {
			continue
		}

		rel, _ := filepath.Rel(specRoot, p.Path)
		if fix {
			newStatus := p.Status()
			if statusDrift {
				newStatus = expectedStatus
			}
			newPromotes := strings.Join(expectedPromotes, ", ")
			if newPromotes == "" {
				newPromotes = "—"
			}
			if err := rewriteIdeaHeader(p.Path, map[string]string{
				"Status":      newStatus,
				"Promotes To": newPromotes,
			}); err == nil {
				fixed = true
				continue
			}
		}

		vs = append(vs, Violation{
			File: rel, Line: 0, Severity: "error",
			Rule:    "idea-sync-lint-strict",
			Message: fmt.Sprintf("idea %q drift: Promotes To / Status disagree with referencing features (run `specscore lint --fix`)", slug),
		})
	}

	return vs, fixed
}

func stringSliceEq(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// rewriteIdeaHeader updates named **Field:** lines in place. Only fields
// already present are updated; no new fields are inserted.
func rewriteIdeaHeader(path string, updates map[string]string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "## ") {
			break
		}
		m := fieldRewriteRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		name := m[2]
		if v, ok := updates[name]; ok {
			lines[i] = fmt.Sprintf("%s**%s:** %s", m[1], name, v)
		}
	}
	return os.WriteFile(path, []byte(strings.Join(lines, "\n")), 0o644)
}

var fieldRewriteRe = regexp.MustCompile(`^(\s*)\*\*([^*]+?):\*\*\s*(.*)$`)
