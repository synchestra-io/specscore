package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

// linter orchestrates rule checking across the spec tree.
type linter struct {
	opts    Options
	ruleSet map[string]checker
}

// checker is the interface for individual rule implementations.
type checker interface {
	check(specRoot string) ([]Violation, error)
	name() string
	severity() string
}

func newLinter(opts Options) *linter {
	l := &linter{
		opts:    opts,
		ruleSet: make(map[string]checker),
	}

	l.registerChecker(newReadmeExistsChecker())
	oqChecker := newOQSectionChecker()
	l.registerChecker(oqChecker)
	l.ruleSet["oq-not-empty"] = oqChecker
	l.registerChecker(newIndexEntriesChecker())
	l.registerChecker(newHeadingLevelsChecker())
	l.registerChecker(newFeatureRefSyntaxChecker())
	l.registerChecker(newInternalLinksChecker())
	l.registerChecker(newForwardRefsChecker())
	l.registerChecker(newCodeAnnotationsChecker())
	l.registerChecker(newPlanHierarchyChecker())
	l.registerChecker(newPlanROIChecker())
	l.registerChecker(newAdherenceFooterChecker())

	// Register custom checkers
	for _, c := range customCheckers {
		l.ruleSet[c.Name()] = &customCheckerAdapter{c}
	}

	return l
}

func (l *linter) registerChecker(c checker) {
	l.ruleSet[c.name()] = c
}

func (l *linter) isRuleEnabled(ruleName string) bool {
	if len(l.opts.Rules) > 0 {
		return slices.Contains(l.opts.Rules, ruleName)
	}

	if len(l.opts.Ignore) > 0 {
		return !slices.Contains(l.opts.Ignore, ruleName)
	}

	return true
}

// lint runs all enabled checkers and returns violations.
// A checker runs if any of the rule names it is registered under are
// enabled.  Individual violations are then filtered so only violations
// whose Rule field matches an enabled rule are returned.
func (l *linter) lint() ([]Violation, error) {
	var violations []Violation

	// Deduplicate checkers (the same checker may be registered under
	// multiple rule names).
	seen := make(map[checker]bool)
	for _, c := range l.ruleSet {
		if seen[c] {
			continue
		}
		seen[c] = true

		// Run checker if any of its registered rule names are enabled.
		enabled := false
		for ruleName, rc := range l.ruleSet {
			if rc == c && l.isRuleEnabled(ruleName) {
				enabled = true
				break
			}
		}
		if !enabled {
			continue
		}

		v, err := c.check(l.opts.SpecRoot)
		if err != nil {
			return nil, fmt.Errorf("checker %s: %v", c.name(), err)
		}

		// Keep only violations whose rule is enabled.
		for _, vi := range v {
			if l.isRuleEnabled(vi.Rule) {
				violations = append(violations, vi)
			}
		}
	}

	return violations, nil
}

// customCheckerAdapter adapts the public Checker interface to the internal checker interface.
type customCheckerAdapter struct {
	c Checker
}

func (a *customCheckerAdapter) name() string     { return a.c.Name() }
func (a *customCheckerAdapter) severity() string { return a.c.Severity() }
func (a *customCheckerAdapter) check(specRoot string) ([]Violation, error) {
	return a.c.Check(specRoot)
}

// walkSpecDirs returns subdirectory paths under specRoot, skipping hidden dirs
// except .github (whose children are traversed but .github itself is skipped).
func walkSpecDirs(specRoot string, fn func(dirPath, relPath string) error) error {
	return filepath.Walk(specRoot, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			return nil
		}
		if strings.HasPrefix(info.Name(), ".") {
			if info.Name() == ".github" {
				return nil // traverse children but skip .github itself
			}
			return filepath.SkipDir
		}
		relPath, relErr := filepath.Rel(specRoot, path)
		if relErr != nil {
			return fmt.Errorf("computing relative path: %w", relErr)
		}
		if relPath == "." {
			relPath = filepath.Base(specRoot)
		}
		return fn(path, relPath)
	})
}
