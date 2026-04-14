package feature

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Info is the top-level metadata structure for a single feature.
type Info struct {
	Path     string        `yaml:"path" json:"path"`
	Status   string        `yaml:"status" json:"status"`
	Deps     []string      `yaml:"deps" json:"deps"`
	Refs     []string      `yaml:"refs" json:"refs"`
	Children []ChildInfo   `yaml:"children,omitempty" json:"children,omitempty"`
	Plans    []string      `yaml:"plans,omitempty" json:"plans,omitempty"`
	Sections []SectionInfo `yaml:"sections" json:"sections"`
}

// ChildInfo represents a child sub-feature.
type ChildInfo struct {
	Path     string `yaml:"path" json:"path"`
	InReadme bool   `yaml:"in_readme" json:"in_readme"`
}

// SectionInfo represents a heading section in the README.
type SectionInfo struct {
	Title    string        `yaml:"title" json:"title"`
	Lines    string        `yaml:"lines" json:"lines"`
	Items    int           `yaml:"items,omitempty" json:"items,omitempty"`
	Children []SectionInfo `yaml:"children,omitempty" json:"children,omitempty"`
}

// GetInfo builds and returns the full Info for a feature.
func GetInfo(featuresDir, featureID string) (*Info, error) {
	readmePath := ReadmePath(featuresDir, featureID)

	status, err := ParseFeatureStatus(readmePath)
	if err != nil {
		return nil, fmt.Errorf("reading feature status: %w", err)
	}

	deps, err := ParseDependencies(readmePath)
	if err != nil {
		return nil, fmt.Errorf("reading dependencies: %w", err)
	}

	refs, err := FindFeatureRefs(featuresDir, featureID)
	if err != nil {
		return nil, fmt.Errorf("finding references: %w", err)
	}

	children, err := DiscoverChildFeatures(featuresDir, featureID, readmePath)
	if err != nil {
		return nil, fmt.Errorf("discovering children: %w", err)
	}

	specRoot := filepath.Dir(featuresDir) // spec/features/ -> spec/
	plans, err := FindLinkedPlans(filepath.Dir(specRoot), featureID)
	if err != nil {
		return nil, fmt.Errorf("finding linked plans: %w", err)
	}

	sections, err := ParseSections(readmePath)
	if err != nil {
		return nil, fmt.Errorf("parsing sections: %w", err)
	}

	return &Info{
		Path:     featureID,
		Status:   status,
		Deps:     deps,
		Refs:     refs,
		Children: children,
		Plans:    plans,
		Sections: sections,
	}, nil
}

// ParseFeatureStatus extracts the status from a feature README.
// Looks for patterns like "**Status:** Draft" or "Status: Stable".
func ParseFeatureStatus(readmePath string) (string, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return "", err
	}
	defer func() { _ = f.Close() }()

	statusRe := regexp.MustCompile(`^\*?\*?Status:?\*?\*?\s*(.+)$`)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if m := statusRe.FindStringSubmatch(line); m != nil {
			status := strings.TrimSpace(m[1])
			status = strings.Trim(status, "`\"'")
			return status, nil
		}
	}
	if err := scanner.Err(); err != nil {
		return "", err
	}
	return "Unknown", nil
}

// FindFeatureRefs finds all features that reference the given featureID
// as a dependency.
func FindFeatureRefs(featuresDir, featureID string) ([]string, error) {
	allFeatures, err := Discover(featuresDir)
	if err != nil {
		return nil, err
	}

	var refs []string
	for _, f := range allFeatures {
		if f.ID == featureID {
			continue
		}
		readmePath := ReadmePath(featuresDir, f.ID)
		deps, depErr := ParseDependencies(readmePath)
		if depErr != nil {
			continue
		}
		for _, dep := range deps {
			if dep == featureID {
				refs = append(refs, f.ID)
				break
			}
		}
	}
	sort.Strings(refs)
	return refs, nil
}

// DiscoverChildFeatures finds immediate child sub-feature directories and
// checks whether each is listed in the parent's ## Contents table.
func DiscoverChildFeatures(featuresDir, featureID, readmePath string) ([]ChildInfo, error) {
	featureDir := filepath.Join(featuresDir, filepath.FromSlash(featureID))
	entries, err := os.ReadDir(featureDir)
	if err != nil {
		return nil, err
	}

	contentsEntries, err := ParseContentsTable(readmePath)
	if err != nil {
		return nil, err
	}

	var children []ChildInfo
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "_") {
			continue
		}
		childReadme := filepath.Join(featureDir, entry.Name(), "README.md")
		if _, statErr := os.Stat(childReadme); statErr != nil {
			continue
		}
		childPath := featureID + "/" + entry.Name()
		inReadme := contentsEntries[entry.Name()]
		children = append(children, ChildInfo{
			Path:     childPath,
			InReadme: inReadme,
		})
	}
	sort.Slice(children, func(i, j int) bool {
		return children[i].Path < children[j].Path
	})
	return children, nil
}

// ParseContentsTable reads a README and extracts entries from the
// ## Contents section. Returns a map of directory names found in the table.
func ParseContentsTable(readmePath string) (map[string]bool, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	entries := make(map[string]bool)
	inContents := false
	scanner := bufio.NewScanner(f)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "## Contents" {
			inContents = true
			continue
		}
		if inContents && strings.HasPrefix(line, "## ") {
			break
		}
		if inContents && strings.HasPrefix(line, "|") {
			matches := linkRe.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				linkPath := m[2]
				dir := strings.TrimSuffix(linkPath, "/README.md")
				dir = strings.TrimSuffix(dir, "/readme.md")
				dir = strings.TrimPrefix(dir, "./")
				if parts := strings.SplitN(dir, "/", 2); len(parts) > 0 && parts[0] != "" {
					entries[parts[0]] = true
				}
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// FindLinkedPlans scans spec/plans/*/README.md for plans that reference
// the given feature.
func FindLinkedPlans(repoRoot, featureID string) ([]string, error) {
	plansDir := filepath.Join(repoRoot, "spec", "plans")
	if _, err := os.Stat(plansDir); err != nil {
		return nil, nil
	}

	var plans []string
	err := filepath.WalkDir(plansDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if d.Name() != "README.md" {
			return nil
		}
		planDir := filepath.Dir(path)
		if planDir == plansDir {
			return nil
		}
		planName := filepath.Base(planDir)

		if planReferencesFeature(path, featureID) {
			plans = append(plans, planName)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(plans)
	return plans, nil
}

// planReferencesFeature checks if a plan README references the given feature
// in its **Features:** section.
func planReferencesFeature(planReadmePath, featureID string) bool {
	f, err := os.Open(planReadmePath)
	if err != nil {
		return false
	}
	defer func() { _ = f.Close() }()

	inFeatures := false
	scanner := bufio.NewScanner(f)
	linkRe := regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	featureSuffix := "features/" + featureID + "/README.md"

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if strings.HasPrefix(line, "**Features:**") {
			inFeatures = true
			continue
		}
		if inFeatures && !strings.HasPrefix(line, "-") && !strings.HasPrefix(line, " ") && line != "" {
			break
		}
		if inFeatures {
			matches := linkRe.FindAllStringSubmatch(line, -1)
			for _, m := range matches {
				linkPath := m[2]
				if strings.HasSuffix(linkPath, featureSuffix) {
					return true
				}
			}
		}
	}
	return false
}

// ParseSections reads a README and builds a section TOC from markdown
// headings. Supports h2 and h3 nesting. Counts list items (lines starting
// with "- ") within each section.
func ParseSections(readmePath string) ([]SectionInfo, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	type rawSection struct {
		title     string
		level     int
		startLine int
		endLine   int
		items     int
	}

	var raw []rawSection
	scanner := bufio.NewScanner(f)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if strings.HasPrefix(trimmed, "## ") || strings.HasPrefix(trimmed, "### ") {
			level := 2
			title := strings.TrimPrefix(trimmed, "## ")
			if strings.HasPrefix(trimmed, "### ") {
				level = 3
				title = strings.TrimPrefix(trimmed, "### ")
			}
			title = strings.TrimSpace(title)

			if len(raw) > 0 {
				raw[len(raw)-1].endLine = lineNum - 1
			}

			raw = append(raw, rawSection{
				title:     title,
				level:     level,
				startLine: lineNum,
			})
			continue
		}

		if len(raw) > 0 && strings.HasPrefix(trimmed, "- ") {
			raw[len(raw)-1].items++
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if len(raw) > 0 {
		raw[len(raw)-1].endLine = lineNum
	}

	var sections []SectionInfo
	for i := 0; i < len(raw); i++ {
		s := raw[i]
		if s.level == 2 {
			section := SectionInfo{
				Title: s.title,
				Lines: fmt.Sprintf("%d-%d", s.startLine, s.endLine),
				Items: s.items,
			}
			for j := i + 1; j < len(raw) && raw[j].level == 3; j++ {
				child := SectionInfo{
					Title: raw[j].title,
					Lines: fmt.Sprintf("%d-%d", raw[j].startLine, raw[j].endLine),
					Items: raw[j].items,
				}
				section.Children = append(section.Children, child)
			}
			sections = append(sections, section)
		}
	}

	return sections, nil
}

// CountOutstandingQuestions counts list items in the
// ## Outstanding Questions section.
func CountOutstandingQuestions(readmePath string) (int, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return 0, err
	}
	defer func() { _ = f.Close() }()

	inOQ := false
	count := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "## Outstanding Questions" {
			inOQ = true
			continue
		}
		if inOQ && strings.HasPrefix(line, "## ") {
			break
		}
		if inOQ && strings.HasPrefix(line, "- ") {
			count++
		}
	}
	return count, scanner.Err()
}
