package idea

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"sort"
	"strings"
)

// Valid statuses for an Idea.
var ValidStatuses = map[string]bool{
	"Draft":        true,
	"Under Review": true,
	"Approved":     true,
	"Specified":    true,
	"Archived":     true,
}

// Valid relationship types for Related Ideas entries.
var ValidRelationships = map[string]bool{
	"depends_on":     true,
	"alternative_to": true,
	"extends":        true,
	"conflicts_with": true,
}

// RequiredSections names the sections that every Idea MUST have, in order.
var RequiredSections = []string{
	"Problem Statement",
	"Context",
	"Recommended Direction",
	"Alternatives Considered",
	"MVP Scope",
	"Not Doing (and Why)",
	"Key Assumptions to Validate",
	"SpecScore Integration",
	"Open Questions",
}

// RequiredHeaderFields lists required **X:** fields in order.
var RequiredHeaderFields = []string{
	"Status",
	"Date",
	"Owner",
	"Promotes To",
	"Supersedes",
	"Related Ideas",
}

// slugRe matches a valid idea slug.
var slugRe = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

// ValidateSlug returns nil if slug matches `[a-z0-9]+(-[a-z0-9]+)*`.
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug must not be empty")
	}
	if !slugRe.MatchString(slug) {
		return fmt.Errorf("slug %q does not match [a-z0-9]+(-[a-z0-9]+)*", slug)
	}
	return nil
}

// HeaderField captures a single **Name:** value line.
type HeaderField struct {
	Name  string
	Value string
	Line  int
}

// Section captures a single ## heading and its body.
type Section struct {
	Title     string
	StartLine int
	EndLine   int
	Body      string
	Items     []string // lines starting with "- "
}

// Table captures a markdown pipe-table within a section.
type Table struct {
	Headers []string
	Rows    [][]string
}

// Idea is a parsed Idea artifact.
type Idea struct {
	Path           string
	Slug           string
	Title          string // full title line (without leading "# ")
	TitleName      string // name portion after "Idea: "
	TitleOK        bool   // true if title matches "# Idea: <Name>"
	TitleLine      int
	HasTitle       bool
	Fields         []HeaderField // in-order header fields encountered
	FieldByName    map[string]HeaderField
	Sections       []Section
	SectionByTitle map[string]*Section
	RawLines       []string
}

// Status returns the Status field value or "" if missing.
func (i *Idea) Status() string {
	return strings.TrimSpace(i.FieldByName["Status"].Value)
}

// PromotesTo returns the parsed Promotes To slugs (comma-separated list).
// A value of "—" or "-" means empty.
func (i *Idea) PromotesTo() []string {
	return splitCSVSlugs(i.FieldByName["Promotes To"].Value)
}

// Supersedes returns the parsed Supersedes slugs.
func (i *Idea) Supersedes() []string {
	return splitCSVSlugs(i.FieldByName["Supersedes"].Value)
}

// RelatedIdeas returns the raw entries (un-split) from the Related Ideas field.
func (i *Idea) RelatedIdeas() []string {
	return splitCSVSlugs(i.FieldByName["Related Ideas"].Value)
}

// ArchiveReason returns the Archive Reason value or "".
func (i *Idea) ArchiveReason() string {
	return strings.TrimSpace(i.FieldByName["Archive Reason"].Value)
}

// splitCSVSlugs splits a comma-separated header value, normalizing common
// empty-placeholders ("—", "-", "") to a nil slice.
func splitCSVSlugs(raw string) []string {
	v := strings.TrimSpace(raw)
	if v == "" || v == "—" || v == "-" {
		return nil
	}
	parts := strings.Split(v, ",")
	var out []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		out = append(out, p)
	}
	return out
}

// fieldLineRe matches "**Name:** value" lines.
var fieldLineRe = regexp.MustCompile(`^\*\*([^*]+?):\*\*\s*(.*)$`)

// Parse reads an Idea file and returns its parsed representation.
// Parse is resilient — it returns a partial Idea even if the file is
// malformed (missing title, missing sections). Callers (lint rules) decide
// what is an error.
func Parse(path string) (*Idea, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	idea := &Idea{
		Path:           path,
		FieldByName:    make(map[string]HeaderField),
		SectionByTitle: make(map[string]*Section),
	}

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	idea.RawLines = lines

	// First pass: find title, header fields, and section boundaries.
	type rawSection struct {
		title string
		start int
		end   int
	}
	var raws []rawSection

	inHeader := true
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Title: first "# " line
		if !idea.HasTitle && strings.HasPrefix(trimmed, "# ") {
			idea.HasTitle = true
			idea.TitleLine = i + 1
			idea.Title = strings.TrimPrefix(trimmed, "# ")
			if name, ok := strings.CutPrefix(idea.Title, "Idea: "); ok {
				idea.TitleOK = true
				idea.TitleName = strings.TrimSpace(name)
			}
			continue
		}

		// Section heading stops header parsing.
		if strings.HasPrefix(trimmed, "## ") {
			inHeader = false
			title := strings.TrimSpace(strings.TrimPrefix(trimmed, "## "))
			if len(raws) > 0 {
				raws[len(raws)-1].end = i - 1
			}
			raws = append(raws, rawSection{title: title, start: i})
			continue
		}

		if inHeader {
			if m := fieldLineRe.FindStringSubmatch(trimmed); m != nil {
				name := strings.TrimSpace(m[1])
				val := strings.TrimSpace(m[2])
				hf := HeaderField{Name: name, Value: val, Line: i + 1}
				idea.Fields = append(idea.Fields, hf)
				idea.FieldByName[name] = hf
			}
		}
	}
	if len(raws) > 0 {
		raws[len(raws)-1].end = len(lines) - 1
	}

	for _, r := range raws {
		body := ""
		if r.end >= r.start {
			bodyLines := lines[r.start+1 : r.end+1]
			body = strings.Join(bodyLines, "\n")
		}
		var items []string
		for _, bl := range strings.Split(body, "\n") {
			t := strings.TrimSpace(bl)
			if strings.HasPrefix(t, "- ") {
				items = append(items, strings.TrimPrefix(t, "- "))
			}
		}
		sec := Section{
			Title:     r.title,
			StartLine: r.start + 1,
			EndLine:   r.end + 1,
			Body:      body,
			Items:     items,
		}
		idea.Sections = append(idea.Sections, sec)
	}
	for i := range idea.Sections {
		s := &idea.Sections[i]
		idea.SectionByTitle[s.Title] = s
	}

	// Slug from filename (without .md).
	base := pathBase(path)
	base = strings.TrimSuffix(base, ".md")
	idea.Slug = base

	return idea, nil
}

// ParseTable extracts a markdown pipe-table from the given section body.
// Returns nil if no table is present. Only the first table is returned.
func ParseTable(body string) *Table {
	lines := strings.Split(body, "\n")
	var headerLine string
	var dataStart int
	for i, ln := range lines {
		t := strings.TrimSpace(ln)
		if !strings.HasPrefix(t, "|") {
			continue
		}
		// Found first table row. Next line should be separator.
		if i+1 >= len(lines) {
			return nil
		}
		sep := strings.TrimSpace(lines[i+1])
		if !strings.HasPrefix(sep, "|") || !strings.Contains(sep, "-") {
			return nil
		}
		headerLine = t
		dataStart = i + 2
		break
	}
	if headerLine == "" {
		return nil
	}
	tab := &Table{Headers: splitTableRow(headerLine)}
	for j := dataStart; j < len(lines); j++ {
		t := strings.TrimSpace(lines[j])
		if !strings.HasPrefix(t, "|") {
			break
		}
		tab.Rows = append(tab.Rows, splitTableRow(t))
	}
	return tab
}

func splitTableRow(row string) []string {
	row = strings.Trim(row, "|")
	parts := strings.Split(row, "|")
	out := make([]string, len(parts))
	for i, p := range parts {
		out[i] = strings.TrimSpace(p)
	}
	return out
}

// pathBase mirrors filepath.Base without importing filepath just here.
func pathBase(p string) string {
	for i := len(p) - 1; i >= 0; i-- {
		if p[i] == '/' || p[i] == '\\' {
			return p[i+1:]
		}
	}
	return p
}

// Sort returns a slice of statuses sorted alphabetically (helper for tests).
func SortedStatuses() []string {
	var out []string
	for k := range ValidStatuses {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}
