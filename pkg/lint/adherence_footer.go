package lint

import (
	"os"
	"path/filepath"
	"strings"
)

// adherenceFooterURL is the canonical SpecScore Feature Specification URL
// that every feature README must reference.
const adherenceFooterURL = "https://specscore.md/feature-specification"

// adherenceFooterChecker verifies that every feature README contains a link
// to the SpecScore Feature Specification, as required by
// [REQ: adherence-footer](spec/features/feature/README.md#req-adherence-footer).
type adherenceFooterChecker struct{}

func newAdherenceFooterChecker() checker {
	return &adherenceFooterChecker{}
}

func (c *adherenceFooterChecker) name() string     { return "adherence-footer" }
func (c *adherenceFooterChecker) severity() string { return "error" }

func (c *adherenceFooterChecker) check(specRoot string) ([]Violation, error) {
	var violations []Violation

	featuresDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featuresDir)
	if err != nil || !info.IsDir() {
		return nil, nil
	}

	err = filepath.Walk(featuresDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		// Skip reserved directories (per REQ: underscore-reserved).
		// Do not skip featuresDir itself if its name happens to start with '_'.
		if path != featuresDir && strings.HasPrefix(info.Name(), "_") {
			return filepath.SkipDir
		}

		readmePath := filepath.Join(path, "README.md")
		readmeInfo, statErr := os.Stat(readmePath)
		if statErr != nil || readmeInfo.IsDir() {
			return nil
		}

		content, readErr := os.ReadFile(readmePath)
		if readErr != nil {
			return nil
		}

		if strings.Contains(string(content), adherenceFooterURL) {
			return nil
		}

		relPath, _ := filepath.Rel(specRoot, readmePath)
		violations = append(violations, Violation{
			File:     relPath,
			Line:     0,
			Severity: "error",
			Rule:     "adherence-footer",
			Message:  "missing required adherence footer: URL '" + adherenceFooterURL + "' not found in feature README",
		})
		return nil
	})

	if err != nil {
		return nil, err
	}

	return violations, nil
}
