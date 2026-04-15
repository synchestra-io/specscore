package idea

import (
	"bufio"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
)

// Discovered is a summary of an Idea file found during discovery.
type Discovered struct {
	Slug     string
	Path     string // absolute or relative path to the .md file
	Archived bool   // true if located under archived/
}

// Discover walks `<specRoot>/ideas` and returns every idea file found.
// Returns ([], nil) if the directory does not exist.
func Discover(specRoot string) ([]Discovered, error) {
	ideasDir := filepath.Join(specRoot, "ideas")
	info, err := os.Stat(ideasDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	var out []Discovered
	// Active: direct children *.md (exclude README.md).
	entries, err := os.ReadDir(ideasDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "README.md" || !strings.HasSuffix(name, ".md") {
			continue
		}
		out = append(out, Discovered{
			Slug:     strings.TrimSuffix(name, ".md"),
			Path:     filepath.Join(ideasDir, name),
			Archived: false,
		})
	}

	archivedDir := filepath.Join(ideasDir, "archived")
	if ai, err := os.Stat(archivedDir); err == nil && ai.IsDir() {
		aEntries, err := os.ReadDir(archivedDir)
		if err != nil {
			return nil, err
		}
		for _, e := range aEntries {
			if e.IsDir() {
				continue
			}
			name := e.Name()
			if name == "README.md" || !strings.HasSuffix(name, ".md") {
				continue
			}
			out = append(out, Discovered{
				Slug:     strings.TrimSuffix(name, ".md"),
				Path:     filepath.Join(archivedDir, name),
				Archived: true,
			})
		}
	}

	sort.Slice(out, func(i, j int) bool {
		if out[i].Archived != out[j].Archived {
			return !out[i].Archived
		}
		return out[i].Slug < out[j].Slug
	})
	return out, nil
}

// FindIdeaDirectories returns directories that exist at `spec/ideas/<slug>/`
// (violation per REQ: single-file). Ignores the reserved `archived/` dir.
func FindIdeaDirectories(specRoot string) ([]string, error) {
	ideasDir := filepath.Join(specRoot, "ideas")
	info, err := os.Stat(ideasDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}
	entries, err := os.ReadDir(ideasDir)
	if err != nil {
		return nil, err
	}
	var out []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if e.Name() == "archived" {
			continue
		}
		out = append(out, filepath.Join(ideasDir, e.Name()))
	}
	return out, nil
}

// sourceIdeasRe extracts the value portion of the **Source Ideas:** line.
var sourceIdeasRe = regexp.MustCompile(`^\*\*Source Ideas:\*\*\s*(.*)$`)

// FeatureSourceIdeas scans every `spec/features/**/README.md` and returns a
// map of feature slug -> []idea-slug based on the **Source Ideas:** header.
// Features without the field are omitted. Only top-level feature dirs are
// returned (filepath suffix: "features/<slug>/README.md").
func FeatureSourceIdeas(specRoot string) (map[string][]string, error) {
	featuresDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featuresDir)
	if err != nil || !info.IsDir() {
		return map[string][]string{}, nil
	}
	out := make(map[string][]string)
	entries, err := os.ReadDir(featuresDir)
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		readme := filepath.Join(featuresDir, e.Name(), "README.md")
		if _, err := os.Stat(readme); err != nil {
			continue
		}
		ideas, err := parseSourceIdeas(readme)
		if err != nil {
			continue
		}
		if len(ideas) > 0 {
			out[e.Name()] = ideas
		}
	}
	return out, nil
}

func parseSourceIdeas(readmePath string) ([]string, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Stop at first ## heading.
		if strings.HasPrefix(line, "## ") {
			break
		}
		if m := sourceIdeasRe.FindStringSubmatch(line); m != nil {
			return splitCSVSlugs(m[1]), scanner.Err()
		}
	}
	return nil, scanner.Err()
}
