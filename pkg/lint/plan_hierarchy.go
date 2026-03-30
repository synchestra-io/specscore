package lint

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// planHierarchyChecker validates hierarchical plan conventions:
// - Roadmaps (plans with child plan subdirectories) must not have Steps sections
// - Nesting is limited to 2 levels (roadmap -> plan)
type planHierarchyChecker struct{}

func newPlanHierarchyChecker() checker {
	return &planHierarchyChecker{}
}

func (c *planHierarchyChecker) name() string     { return "plan-hierarchy" }
func (c *planHierarchyChecker) severity() string { return "error" }

func (c *planHierarchyChecker) check(specRoot string) ([]Violation, error) {
	plansDir := filepath.Join(specRoot, "plans")
	info, err := os.Stat(plansDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	var violations []Violation

	// Walk plan directories at all levels under plans/
	err = filepath.Walk(plansDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		// Skip hidden dirs, acs, and reports subdirs.
		if strings.HasPrefix(fi.Name(), ".") {
			return filepath.SkipDir
		}
		if fi.Name() == "acs" || fi.Name() == "reports" {
			return filepath.SkipDir
		}
		// Skip the plans/ directory itself.
		if path == plansDir {
			return nil
		}

		readmePath := filepath.Join(path, "README.md")
		if _, statErr := os.Stat(readmePath); statErr != nil {
			return nil
		}

		relReadme, _ := filepath.Rel(specRoot, readmePath)

		// Determine nesting depth relative to plans/
		relToPlans, _ := filepath.Rel(plansDir, path)
		depth := len(strings.Split(relToPlans, string(os.PathSeparator)))

		// Check nesting depth: max 2 levels (roadmap -> plan)
		if depth > 2 {
			violations = append(violations, Violation{
				File:     relReadme,
				Line:     0,
				Severity: "error",
				Rule:     "plan-hierarchy",
				Message:  "Plan nesting depth exceeds 2 levels; maximum is roadmap -> plan",
			})
			return nil
		}

		// Detect if this is a roadmap (has child plan subdirectories containing README.md)
		isRoadmap := hasChildPlanDirs(path)

		// Parse the README for section headings
		hasSteps, stepsLine := hasSection(readmePath, "## Steps")

		if isRoadmap && hasSteps {
			violations = append(violations, Violation{
				File:     relReadme,
				Line:     stepsLine,
				Severity: "error",
				Rule:     "plan-hierarchy",
				Message:  "Roadmap must not have a Steps section; use Child Plans instead",
			})
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return violations, nil
}

// hasChildPlanDirs checks whether the directory contains subdirectories with README.md files
// (indicating child plans), skipping hidden dirs, acs, and reports.
func hasChildPlanDirs(dir string) bool {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "acs" || name == "reports" {
			continue
		}
		childReadme := filepath.Join(dir, name, "README.md")
		if _, statErr := os.Stat(childReadme); statErr == nil {
			return true
		}
	}
	return false
}

// hasSection scans a markdown file for a specific heading and returns whether it was found
// and the line number.
func hasSection(readmePath, heading string) (bool, int) {
	file, err := os.Open(readmePath)
	if err != nil {
		return false, 0
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := strings.TrimSpace(scanner.Text())
		if line == heading || strings.HasPrefix(line, heading+" ") {
			return true, lineNum
		}
	}
	return false, 0
}
