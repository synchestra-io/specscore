package lint

import (
	"fmt"
	"os"
)

// Violation represents a single linting violation.
type Violation struct {
	File     string `json:"file" yaml:"file"`
	Line     int    `json:"line" yaml:"line"`
	Severity string `json:"severity" yaml:"severity"`
	Rule     string `json:"rule" yaml:"rule"`
	Message  string `json:"message" yaml:"message"`
}

// Options holds linting options.
type Options struct {
	SpecRoot string
	Rules    []string // enabled rules; nil = all
	Ignore   []string // disabled rules
	Severity string   // minimum severity: error, warning, info
}

// Lint runs all enabled lint rules against the spec tree.
func Lint(opts Options) ([]Violation, error) {
	// Check spec root exists.
	info, err := os.Stat(opts.SpecRoot)
	if err != nil {
		return nil, fmt.Errorf("spec root not found: %s", opts.SpecRoot)
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("spec root is not a directory: %s", opts.SpecRoot)
	}

	l := newLinter(opts)
	violations, err := l.lint()
	if err != nil {
		return nil, fmt.Errorf("linting error: %w", err)
	}

	// Filter by severity if specified.
	if opts.Severity != "" {
		violations = FilterBySeverity(violations, opts.Severity)
	}

	return violations, nil
}

// FilterBySeverity filters violations to those at or above the minimum severity.
func FilterBySeverity(violations []Violation, minSeverity string) []Violation {
	severityOrder := map[string]int{"error": 0, "warning": 1, "info": 2}
	minLevel, ok := severityOrder[minSeverity]
	if !ok {
		return violations
	}

	var filtered []Violation
	for _, v := range violations {
		if severityOrder[v.Severity] <= minLevel {
			filtered = append(filtered, v)
		}
	}
	return filtered
}

// Checker is the public interface for custom rule implementations.
// External tools can implement this to add custom rules to the linter.
type Checker interface {
	Name() string
	Severity() string
	Check(specRoot string) ([]Violation, error)
}

var customCheckers []Checker

// RegisterChecker registers a custom checker that will run alongside
// built-in checkers during Lint().
func RegisterChecker(c Checker) {
	customCheckers = append(customCheckers, c)
}

// ResetCustomCheckers clears all registered custom checkers (for testing).
func ResetCustomCheckers() {
	customCheckers = nil
}

// allRuleNames is the canonical list of known rule names.
var allRuleNames = map[string]bool{
	"readme-exists":      true,
	"oq-section":         true,
	"oq-not-empty":       true,
	"index-entries":      true,
	"heading-levels":     true,
	"feature-ref-syntax": true,
	"internal-links":     true,
	"forward-refs":       true,
	"code-annotations":   true,
	"plan-hierarchy":     true,
	"plan-roi-metadata":  true,
}

// AllRuleNames returns the canonical set of known rule names.
func AllRuleNames() map[string]bool {
	// Return a copy to prevent mutation.
	result := make(map[string]bool, len(allRuleNames))
	for k, v := range allRuleNames {
		result[k] = v
	}
	return result
}

// ValidateRuleNames checks that all names are known rules.
func ValidateRuleNames(names []string) error {
	for _, name := range names {
		if !allRuleNames[name] {
			return fmt.Errorf("unknown rule %q", name)
		}
	}
	return nil
}
