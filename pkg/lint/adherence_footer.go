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
	err := walkFeatureReadmes(specRoot, func(readmePath string, content []byte) {
		if strings.Contains(string(content), adherenceFooterURL) {
			return
		}
		relPath, _ := filepath.Rel(specRoot, readmePath)
		violations = append(violations, Violation{
			File:     relPath,
			Line:     0,
			Severity: "error",
			Rule:     "adherence-footer",
			Message:  "missing required adherence footer: URL '" + adherenceFooterURL + "' not found in feature README",
		})
	})
	if err != nil {
		return nil, err
	}
	return violations, nil
}

// fix appends the canonical adherence footer to any feature README that is
// missing it. A README already containing the URL is left untouched.
func (c *adherenceFooterChecker) fix(specRoot string) error {
	return walkFeatureReadmes(specRoot, func(readmePath string, content []byte) {
		s := string(content)
		if strings.Contains(s, adherenceFooterURL) {
			return
		}
		if !strings.HasSuffix(s, "\n") {
			s += "\n"
		}
		s += "\n---\n*This document follows the " + adherenceFooterURL + "*\n"
		_ = os.WriteFile(readmePath, []byte(s), 0o644)
	})
}
