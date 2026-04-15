package lint

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/synchestra-io/specscore/pkg/idea"
)

// ideaIndexRules enforces:
//   - idea-index-completeness: spec/ideas/README.md lists every active idea;
//     spec/ideas/archived/README.md lists every archived idea.
//   - idea-archived-index-chronological: archived entries ordered by Date asc.
//
// When fix is true, missing/misordered entries are regenerated.
func ideaIndexRules(specRoot string, discovered []idea.Discovered, parsed map[string]*idea.Idea, fix bool) ([]Violation, bool) {
	var vs []Violation
	fixed := false

	ideasDir := filepath.Join(specRoot, "ideas")
	archivedDir := filepath.Join(ideasDir, "archived")

	var active []idea.Discovered
	var archived []idea.Discovered
	for _, d := range discovered {
		if d.Archived {
			archived = append(archived, d)
		} else {
			active = append(active, d)
		}
	}

	// Active index.
	activeIdx := filepath.Join(ideasDir, "README.md")
	if _, err := os.Stat(activeIdx); err == nil {
		listed, err := readIndexSlugs(activeIdx)
		if err == nil {
			listedSet := toSet(listed)
			var missing []string
			for _, d := range active {
				if !listedSet[d.Slug] {
					missing = append(missing, d.Slug)
				}
			}
			sort.Strings(missing)
			if len(missing) > 0 {
				rel, _ := filepath.Rel(specRoot, activeIdx)
				if fix {
					if err := rewriteActiveIndex(activeIdx, active, parsed); err == nil {
						fixed = true
					} else {
						vs = append(vs, Violation{
							File: rel, Line: 0, Severity: "error",
							Rule:    "idea-index-completeness",
							Message: fmt.Sprintf("active idea index missing entries: %s (fix failed: %v)", strings.Join(missing, ", "), err),
						})
					}
				} else {
					vs = append(vs, Violation{
						File: rel, Line: 0, Severity: "error",
						Rule:    "idea-index-completeness",
						Message: fmt.Sprintf("active idea index missing entries: %s", strings.Join(missing, ", ")),
					})
				}
			}
		}
	}

	// Archived index.
	archivedIdx := filepath.Join(archivedDir, "README.md")
	if _, err := os.Stat(archivedIdx); err == nil {
		entries, err := readArchivedEntries(archivedIdx)
		if err == nil {
			listedSet := make(map[string]bool)
			for _, e := range entries {
				listedSet[e.slug] = true
			}
			var missing []string
			for _, d := range archived {
				if !listedSet[d.Slug] {
					missing = append(missing, d.Slug)
				}
			}
			sort.Strings(missing)

			rel, _ := filepath.Rel(specRoot, archivedIdx)

			// Chronological check based on listed entries.
			chronoErr := false
			for i := 1; i < len(entries); i++ {
				if entries[i-1].date > entries[i].date {
					chronoErr = true
					break
				}
			}

			if len(missing) > 0 || chronoErr {
				if fix {
					if err := rewriteArchivedIndex(archivedIdx, archived, parsed); err == nil {
						fixed = true
					} else {
						if len(missing) > 0 {
							vs = append(vs, Violation{
								File: rel, Line: 0, Severity: "error",
								Rule:    "idea-index-completeness",
								Message: fmt.Sprintf("archived idea index missing entries: %s (fix failed: %v)", strings.Join(missing, ", "), err),
							})
						}
						if chronoErr {
							vs = append(vs, Violation{
								File: rel, Line: 0, Severity: "error",
								Rule:    "idea-archived-index-chronological",
								Message: fmt.Sprintf("archived idea index entries must appear in chronological order by Date (oldest first) (fix failed: %v)", err),
							})
						}
					}
				} else {
					if len(missing) > 0 {
						vs = append(vs, Violation{
							File: rel, Line: 0, Severity: "error",
							Rule:    "idea-index-completeness",
							Message: fmt.Sprintf("archived idea index missing entries: %s", strings.Join(missing, ", ")),
						})
					}
					if chronoErr {
						vs = append(vs, Violation{
							File: rel, Line: 0, Severity: "error",
							Rule:    "idea-archived-index-chronological",
							Message: "archived idea index entries must appear in chronological order by Date (oldest first)",
						})
					}
				}
			}
		}
	}

	return vs, fixed
}

func toSet(s []string) map[string]bool {
	m := make(map[string]bool, len(s))
	for _, v := range s {
		m[v] = true
	}
	return m
}

// readIndexSlugs scans an index README and returns the slugs mentioned in
// the ## Index table (via links of form [text](<slug>.md)).
var linkRe = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\.md\)`)

func readIndexSlugs(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var slugs []string
	inIndex := false
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if strings.HasPrefix(line, "## Index") {
			inIndex = true
			continue
		}
		if inIndex && strings.HasPrefix(line, "## ") {
			break
		}
		if !inIndex {
			continue
		}
		matches := linkRe.FindAllStringSubmatch(line, -1)
		for _, m := range matches {
			p := m[2]
			if strings.Contains(p, "/") {
				continue
			}
			slugs = append(slugs, p)
		}
	}
	return slugs, scanner.Err()
}

// archivedEntryRe matches lines like:
//
//   - 2024-11-02 — [offline-mode](offline-mode.md) — pivoted
var archivedEntryRe = regexp.MustCompile(`^-\s+(\d{4}-\d{2}-\d{2})\s+—\s+\[([^\]]+)\]\(([^)]+)\.md\)\s+—\s+(.+)$`)

type archivedEntry struct {
	date   string
	slug   string
	reason string
}

func readArchivedEntries(path string) ([]archivedEntry, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	var out []archivedEntry
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		m := archivedEntryRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		out = append(out, archivedEntry{date: m[1], slug: m[2], reason: m[4]})
	}
	return out, scanner.Err()
}

// rewriteActiveIndex regenerates `spec/ideas/README.md` preserving the
// prologue before "## Index" and replacing the table body.
func rewriteActiveIndex(path string, active []idea.Discovered, parsed map[string]*idea.Idea) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	// Locate "## Index" and "## Outstanding Questions".
	indexStart := -1
	oqStart := -1
	for i, ln := range lines {
		t := strings.TrimSpace(ln)
		if strings.HasPrefix(t, "## Index") && indexStart == -1 {
			indexStart = i
		} else if strings.HasPrefix(t, "## ") && indexStart != -1 && oqStart == -1 {
			oqStart = i
			break
		}
	}
	if indexStart == -1 {
		return fmt.Errorf("cannot locate ## Index heading")
	}

	// Build new Index section.
	var tbl strings.Builder
	tbl.WriteString("## Index\n\n")
	tbl.WriteString("| Idea | Status | Date | Owner | Promotes To |\n")
	tbl.WriteString("|------|--------|------|-------|-------------|\n")

	sorted := append([]idea.Discovered{}, active...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Slug < sorted[j].Slug })

	if len(sorted) == 0 {
		tbl.WriteString("\n_No active ideas yet._\n\n")
	} else {
		for _, d := range sorted {
			p := parsed[d.Slug]
			status, date, owner, promotes := "", "", "", "—"
			if p != nil {
				status = p.Status()
				date = strings.TrimSpace(p.FieldByName["Date"].Value)
				owner = strings.TrimSpace(p.FieldByName["Owner"].Value)
				if pt := p.PromotesTo(); len(pt) > 0 {
					promotes = strings.Join(pt, ", ")
				}
			}
			fmt.Fprintf(&tbl, "| [%s](%s.md) | %s | %s | %s | %s |\n",
				d.Slug, d.Slug, status, date, owner, promotes)
		}
		tbl.WriteString("\n")
	}

	var newLines []string
	newLines = append(newLines, lines[:indexStart]...)
	newLines = append(newLines, strings.Split(strings.TrimRight(tbl.String(), "\n"), "\n")...)
	newLines = append(newLines, "")
	if oqStart != -1 {
		newLines = append(newLines, lines[oqStart:]...)
	}
	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0o644)
}

// rewriteArchivedIndex regenerates the chronological list. Preserves the
// prologue up to (but not including) the first `- YYYY-MM-DD` entry or the
// `_No archived ideas yet._` marker, and the "## Outstanding Questions"
// section if present.
func rewriteArchivedIndex(path string, archived []idea.Discovered, parsed map[string]*idea.Idea) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	lines := strings.Split(string(data), "\n")

	// Find the insertion region: after the code fence line closing "```"
	// (format description), and before "## Outstanding Questions".
	oqStart := -1
	for i, ln := range lines {
		if strings.HasPrefix(strings.TrimSpace(ln), "## Outstanding Questions") {
			oqStart = i
			break
		}
	}

	// Find end of prologue: the last non-entry, non-"no archived" line
	// before we hit entries or the OQ heading. Simpler approach: keep
	// everything up to the first line matching archivedEntryRe or
	// "_No archived ideas yet._", whichever comes first.
	prologueEnd := len(lines)
	if oqStart != -1 {
		prologueEnd = oqStart
	}
	for i := 0; i < prologueEnd; i++ {
		t := strings.TrimSpace(lines[i])
		if archivedEntryRe.MatchString(t) || strings.HasPrefix(t, "_No archived ideas yet._") {
			prologueEnd = i
			break
		}
	}
	// Trim trailing blank lines in prologue.
	for prologueEnd > 0 && strings.TrimSpace(lines[prologueEnd-1]) == "" {
		prologueEnd--
	}

	// Build entries sorted by Date asc.
	type rowT struct {
		date, slug, reason string
	}
	var rows []rowT
	for _, d := range archived {
		p := parsed[d.Slug]
		if p == nil {
			continue
		}
		rows = append(rows, rowT{
			date:   strings.TrimSpace(p.FieldByName["Date"].Value),
			slug:   d.Slug,
			reason: p.ArchiveReason(),
		})
	}
	sort.Slice(rows, func(i, j int) bool {
		if rows[i].date != rows[j].date {
			return rows[i].date < rows[j].date
		}
		return rows[i].slug < rows[j].slug
	})

	var body strings.Builder
	if len(rows) == 0 {
		body.WriteString("_No archived ideas yet._\n")
	} else {
		for _, r := range rows {
			reason := r.reason
			if reason == "" {
				reason = "—"
			}
			fmt.Fprintf(&body, "- %s — [%s](%s.md) — %s\n", r.date, r.slug, r.slug, reason)
		}
	}

	var newLines []string
	newLines = append(newLines, lines[:prologueEnd]...)
	newLines = append(newLines, "", "")
	newLines = append(newLines, strings.Split(strings.TrimRight(body.String(), "\n"), "\n")...)
	if oqStart != -1 {
		newLines = append(newLines, "", "")
		newLines = append(newLines, lines[oqStart:]...)
	}
	return os.WriteFile(path, []byte(strings.Join(newLines, "\n")), 0o644)
}
