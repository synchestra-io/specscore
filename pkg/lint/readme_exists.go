package lint

import (
	"os"
	"path/filepath"
)

// readmeExistsChecker verifies that every spec directory has a README.md file.
type readmeExistsChecker struct{}

func newReadmeExistsChecker() checker {
	return &readmeExistsChecker{}
}

func (c *readmeExistsChecker) name() string     { return "readme-exists" }
func (c *readmeExistsChecker) severity() string { return "error" }

func (c *readmeExistsChecker) check(specRoot string) ([]Violation, error) {
	var violations []Violation

	err := walkSpecDirs(specRoot, func(dirPath, relPath string) error {
		readmePath := filepath.Join(dirPath, "README.md")
		if _, err := os.Stat(readmePath); err != nil {
			violations = append(violations, Violation{
				File:     relPath,
				Line:     0,
				Severity: "error",
				Rule:     "readme-exists",
				Message:  "README.md not found",
			})
		}
		return nil
	})

	return violations, err
}
