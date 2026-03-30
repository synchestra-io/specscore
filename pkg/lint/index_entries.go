package lint

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// indexEntriesChecker verifies that feature README indices match actual child directories.
type indexEntriesChecker struct{}

func newIndexEntriesChecker() checker {
	return &indexEntriesChecker{}
}

func (c *indexEntriesChecker) name() string     { return "index-entries" }
func (c *indexEntriesChecker) severity() string { return "error" }

func (c *indexEntriesChecker) check(specRoot string) ([]Violation, error) {
	var violations []Violation

	featureDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featureDir)
	if err != nil || !info.IsDir() {
		return violations, nil
	}

	err = filepath.Walk(featureDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}

		readmePath := filepath.Join(path, "README.md")
		if _, statErr := os.Stat(readmePath); statErr != nil {
			return nil
		}

		// Get actual child directories (excluding hidden and _args convention dirs).
		entries, readErr := os.ReadDir(path)
		if readErr != nil {
			return nil
		}

		var actualChildren []string
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") && !strings.HasPrefix(entry.Name(), "_") {
				actualChildren = append(actualChildren, entry.Name())
			}
		}

		mentioned, parseErr := extractChildRefsFromReadme(readmePath)
		if parseErr != nil || mentioned == nil {
			return nil
		}

		relPath, _ := filepath.Rel(specRoot, readmePath)

		// Flag index entries that reference non-existent directories.
		actualSet := make(map[string]bool, len(actualChildren))
		for _, a := range actualChildren {
			actualSet[a] = true
		}
		for _, m := range mentioned {
			if !actualSet[m] {
				violations = append(violations, Violation{
					File:     relPath,
					Line:     0,
					Severity: "error",
					Rule:     "index-entries",
					Message:  "Index mentions non-existent directory: " + m,
				})
			}
		}

		return nil
	})

	return violations, err
}

// extractChildRefsFromReadme scans a README for markdown links pointing to
// child directories (e.g. `[name](dirname/README.md)`).
func extractChildRefsFromReadme(readmePath string) ([]string, error) {
	file, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	var children []string
	seen := make(map[string]bool)
	inCodeBlock := false

	for scanner.Scan() {
		line := scanner.Text()

		// Skip fenced code blocks.
		if strings.HasPrefix(strings.TrimSpace(line), "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock {
			continue
		}

		// Look for links to child README.md: [text](dirname/README.md)
		for {
			idx := strings.Index(line, "](")
			if idx < 0 {
				break
			}
			rest := line[idx+2:]
			end := strings.Index(rest, ")")
			if end < 0 {
				break
			}
			linkTarget := rest[:end]
			line = rest[end+1:] // advance past this link

			// Only consider links ending in README.md and pointing to a direct child.
			if !strings.HasSuffix(linkTarget, "README.md") && !strings.HasSuffix(linkTarget, "README.md)") {
				continue
			}
			parts := strings.Split(strings.TrimPrefix(linkTarget, "./"), "/")
			if len(parts) == 2 {
				dirname := parts[0]
				if dirname != "." && dirname != ".." && !strings.HasPrefix(dirname, "_") && !seen[dirname] {
					seen[dirname] = true
					children = append(children, dirname)
				}
			}
		}
	}

	if len(children) == 0 {
		return nil, nil
	}
	return children, scanner.Err()
}
