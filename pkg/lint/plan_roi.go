package lint

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// planROIChecker validates ROI metadata values in plan README headers.
// When present, Effort must be S/M/L/XL and Impact must be low/medium/high/critical.
type planROIChecker struct{}

func newPlanROIChecker() checker {
	return &planROIChecker{}
}

func (c *planROIChecker) name() string     { return "plan-roi-metadata" }
func (c *planROIChecker) severity() string { return "warning" }

var validEffort = map[string]bool{
	"S": true, "M": true, "L": true, "XL": true,
}

var validImpact = map[string]bool{
	"low": true, "medium": true, "high": true, "critical": true,
}

func (c *planROIChecker) check(specRoot string) ([]Violation, error) {
	plansDir := filepath.Join(specRoot, "plans")
	info, err := os.Stat(plansDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	var violations []Violation

	err = filepath.Walk(plansDir, func(path string, fi os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !fi.IsDir() {
			return nil
		}
		if strings.HasPrefix(fi.Name(), ".") {
			return filepath.SkipDir
		}
		// Skip the plans/ directory itself (its README is the index).
		if path == plansDir {
			return nil
		}

		readmePath := filepath.Join(path, "README.md")
		if _, statErr := os.Stat(readmePath); statErr != nil {
			return nil
		}

		relReadme, _ := filepath.Rel(specRoot, readmePath)

		v, scanErr := scanROIMetadata(readmePath, relReadme)
		if scanErr != nil {
			return scanErr
		}
		violations = append(violations, v...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return violations, nil
}

// scanROIMetadata reads the header of a plan README (lines before the first ## heading)
// and validates Effort/Impact values if present.
func scanROIMetadata(readmePath, relPath string) ([]Violation, error) {
	file, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	var violations []Violation
	scanner := bufio.NewScanner(file)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		// Stop at the first ## heading — we only check the header.
		if strings.HasPrefix(trimmed, "## ") {
			break
		}

		if value, ok := strings.CutPrefix(trimmed, "**Effort:**"); ok {
			value = strings.TrimSpace(value)
			if !validEffort[value] {
				violations = append(violations, Violation{
					File:     relPath,
					Line:     lineNum,
					Severity: "warning",
					Rule:     "plan-roi-metadata",
					Message:  "Effort value must be one of S, M, L, XL; got " + value,
				})
			}
		}

		if value, ok := strings.CutPrefix(trimmed, "**Impact:**"); ok {
			value = strings.TrimSpace(value)
			if !validImpact[value] {
				violations = append(violations, Violation{
					File:     relPath,
					Line:     lineNum,
					Severity: "warning",
					Rule:     "plan-roi-metadata",
					Message:  "Impact value must be one of low, medium, high, critical; got " + value,
				})
			}
		}
	}

	return violations, nil
}
