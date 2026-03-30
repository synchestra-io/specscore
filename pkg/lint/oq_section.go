package lint

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// oqSectionChecker verifies that feature/plan READMEs have an Outstanding Questions section.
type oqSectionChecker struct{}

func newOQSectionChecker() checker {
	return &oqSectionChecker{}
}

func (c *oqSectionChecker) name() string     { return "oq-section" }
func (c *oqSectionChecker) severity() string { return "error" }

func (c *oqSectionChecker) check(specRoot string) ([]Violation, error) {
	var violations []Violation

	specSubDirs := []string{
		filepath.Join(specRoot, "features"),
		filepath.Join(specRoot, "plans"),
	}

	for _, dir := range specSubDirs {
		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}

		err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
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

			result, parseErr := parseOQSection(readmePath)
			if parseErr != nil {
				return nil
			}

			relPath, _ := filepath.Rel(specRoot, readmePath)

			if !result.found {
				violations = append(violations, Violation{
					File:     relPath,
					Line:     0,
					Severity: "error",
					Rule:     "oq-section",
					Message:  "Outstanding Questions section not found",
				})
			} else if result.empty {
				violations = append(violations, Violation{
					File:     relPath,
					Line:     result.line,
					Severity: "warning",
					Rule:     "oq-not-empty",
					Message:  "Outstanding Questions section appears empty",
				})
			}

			return nil
		})

		if err != nil {
			return nil, err
		}
	}

	return violations, nil
}

type oqResult struct {
	found bool
	empty bool
	line  int
}

// parseOQSection scans a README for "## Outstanding Questions" and determines
// whether it exists and whether it has content.
func parseOQSection(readmePath string) (oqResult, error) {
	file, err := os.Open(readmePath)
	if err != nil {
		return oqResult{}, err
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()

		if !strings.HasPrefix(line, "## Outstanding Questions") {
			continue
		}

		oqLine := lineNum

		// Scan forward to see if the section has content.
		for scanner.Scan() {
			lineNum++
			next := strings.TrimSpace(scanner.Text())
			if next == "" {
				continue
			}
			// A new heading means the OQ section was empty.
			if strings.HasPrefix(next, "#") {
				return oqResult{found: true, empty: true, line: oqLine}, nil
			}
			// Any non-blank, non-heading content means it's populated.
			return oqResult{found: true, empty: false, line: oqLine}, nil
		}

		// OQ heading was the last thing in the file with no content after it.
		return oqResult{found: true, empty: true, line: oqLine}, nil
	}

	return oqResult{found: false}, scanner.Err()
}
