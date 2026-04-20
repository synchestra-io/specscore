package lint

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// docTypeTarget describes one document type subject to the adherence-footer
// check. The walk function enumerates the consumer paths defined for that
// Kind by the document types registry.
type docTypeTarget struct {
	description string
	url         string
	severity    string
	walk        func(specRoot string, fn func(path string, content []byte)) error
}

// docTypeTargets mirrors the Document-Kind and Index-Kind rows in the
// document types registry at spec/features/README.md. New Document/Index
// Kinds must be added here.
//
// Severity policy during the MVP rollout: feature READMEs stay at "error"
// (pre-existing behavior), new consumer-layer checks ship at "warn" per
// the adherence-footer-and-doc-type-registry Idea's MVP gate, then
// transition to "error" after a cycle of clean runs.
var docTypeTargets = []docTypeTarget{
	{
		description: "feature README",
		url:         "https://specscore.md/feature-specification",
		severity:    "error",
		walk:        walkFeatureReadmesExcludingIndex,
	},
	{
		description: "features-index README",
		url:         "https://specscore.md/features-index-specification",
		severity:    "warn",
		walk:        walkFeaturesIndex,
	},
	{
		description: "Idea file",
		url:         "https://specscore.md/idea-specification",
		severity:    "warn",
		walk:        walkIdeaFiles,
	},
	{
		description: "plans-index README",
		url:         "https://specscore.md/plans-index-specification",
		severity:    "warn",
		walk:        walkPlansIndex,
	},
	{
		description: "ideas-index README",
		url:         "https://specscore.md/ideas-index-specification",
		severity:    "warn",
		walk:        walkIdeasIndex,
	},
	{
		description: "Plan README",
		url:         "https://specscore.md/plan-specification",
		severity:    "warn",
		walk:        walkPlanReadmes,
	},
	{
		description: "Task README",
		url:         "https://specscore.md/task-specification",
		severity:    "warn",
		walk:        walkTaskReadmes,
	},
	{
		description: "Scenario file",
		url:         "https://specscore.md/scenario-specification",
		severity:    "warn",
		walk:        walkScenarioFiles,
	},
	{
		description: "scenarios-index README",
		url:         "https://specscore.md/scenarios-index-specification",
		severity:    "warn",
		walk:        walkScenariosIndexes,
	},
}

// adherenceFooterChecker verifies that every SpecScore document of a
// Document or Index Kind carries the adherence footer URL corresponding
// to its document type, as required by the Adherence Footer feature.
type adherenceFooterChecker struct{}

func newAdherenceFooterChecker() checker {
	return &adherenceFooterChecker{}
}

func (c *adherenceFooterChecker) name() string     { return "adherence-footer" }
func (c *adherenceFooterChecker) severity() string { return "error" }

func (c *adherenceFooterChecker) check(specRoot string) ([]Violation, error) {
	var violations []Violation
	for _, t := range docTypeTargets {
		target := t
		err := target.walk(specRoot, func(path string, content []byte) {
			if strings.Contains(string(content), target.url) {
				return
			}
			relPath, _ := filepath.Rel(specRoot, path)
			violations = append(violations, Violation{
				File:     relPath,
				Line:     0,
				Severity: target.severity,
				Rule:     "adherence-footer",
				Message:  fmt.Sprintf("missing required adherence footer: URL %q not found in %s", target.url, target.description),
			})
		})
		if err != nil {
			return nil, err
		}
	}
	return violations, nil
}

// fix appends the canonical adherence footer to any document that is missing
// it. A document already containing the URL is left untouched. Wrong-URL
// violations are never auto-rewritten — a wrong URL may indicate a
// mis-classified document and silent rewriting would mask the bug. See
// adherence-footer#req:fix-inserts-only.
func (c *adherenceFooterChecker) fix(specRoot string) error {
	for _, t := range docTypeTargets {
		target := t
		err := target.walk(specRoot, func(path string, content []byte) {
			if strings.Contains(string(content), target.url) {
				return
			}
			s := string(content)
			if !strings.HasSuffix(s, "\n") {
				s += "\n"
			}
			s += "\n---\n*This document follows the " + target.url + "*\n"
			_ = os.WriteFile(path, []byte(s), 0o644)
		})
		if err != nil {
			return err
		}
	}
	return nil
}

// walkIdeaFiles invokes fn for every Idea file under specRoot/ideas/*.md,
// skipping README.md (the ideas index) and the archived/ subtree.
func walkIdeaFiles(specRoot string, fn func(path string, content []byte)) error {
	ideasDir := filepath.Join(specRoot, "ideas")
	return walkMatchingFiles(ideasDir, func(path string, depth int, name string) bool {
		if depth != 1 {
			return false
		}
		if name == "README.md" {
			return false
		}
		return strings.HasSuffix(name, ".md")
	}, fn)
}

// walkPlansIndex invokes fn for specRoot/plans/README.md if present.
func walkPlansIndex(specRoot string, fn func(path string, content []byte)) error {
	path := filepath.Join(specRoot, "plans", "README.md")
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fn(path, content)
	return nil
}

// walkIdeasIndex invokes fn for specRoot/ideas/README.md if present.
func walkIdeasIndex(specRoot string, fn func(path string, content []byte)) error {
	path := filepath.Join(specRoot, "ideas", "README.md")
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fn(path, content)
	return nil
}

// walkFeaturesIndex invokes fn for specRoot/features/README.md if present.
// This file is an Index-Kind document (features-index), not a Feature README,
// and is checked against the features-index-specification URL.
func walkFeaturesIndex(specRoot string, fn func(path string, content []byte)) error {
	path := filepath.Join(specRoot, "features", "README.md")
	content, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	fn(path, content)
	return nil
}

// walkFeatureReadmesExcludingIndex walks feature READMEs but skips the root
// spec/features/README.md (which is an Index-Kind document, not a Feature).
// Used by the feature-specification adherence-footer check so the root index
// isn't flagged for missing feature-specification — its footer is
// features-index-specification, handled by walkFeaturesIndex.
func walkFeatureReadmesExcludingIndex(specRoot string, fn func(path string, content []byte)) error {
	rootIndex := filepath.Join(specRoot, "features", "README.md")
	return walkFeatureReadmes(specRoot, func(path string, content []byte) {
		if path == rootIndex {
			return
		}
		fn(path, content)
	})
}

// walkPlanReadmes invokes fn for every Plan README under specRoot/plans/**/README.md,
// excluding specRoot/plans/README.md (which is the plans-index, walked separately)
// and any README.md inside a reserved _-prefixed directory.
func walkPlanReadmes(specRoot string, fn func(path string, content []byte)) error {
	plansDir := filepath.Join(specRoot, "plans")
	info, err := os.Stat(plansDir)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(plansDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != plansDir && strings.HasPrefix(info.Name(), "_") {
				return filepath.SkipDir
			}
			return nil
		}
		if info.Name() != "README.md" {
			return nil
		}
		// Skip the plans-index itself (handled by walkPlansIndex).
		if path == filepath.Join(plansDir, "README.md") {
			return nil
		}
		// Skip task READMEs — tasks/<slug>/README.md is a Task, not a Plan.
		if filepath.Base(filepath.Dir(filepath.Dir(path))) == "tasks" {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		fn(path, content)
		return nil
	})
}

// walkTaskReadmes invokes fn for every Task README under
// specRoot/plans/**/tasks/*/README.md.
func walkTaskReadmes(specRoot string, fn func(path string, content []byte)) error {
	plansDir := filepath.Join(specRoot, "plans")
	info, err := os.Stat(plansDir)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(plansDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != "README.md" {
			return nil
		}
		// A Task README lives at plans/**/tasks/<slug>/README.md. Check that
		// the grandparent directory is "tasks".
		if filepath.Base(filepath.Dir(filepath.Dir(path))) != "tasks" {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		fn(path, content)
		return nil
	})
}

// walkScenariosIndexes invokes fn for every scenarios index at
// specRoot/features/**/_tests/README.md.
func walkScenariosIndexes(specRoot string, fn func(path string, content []byte)) error {
	featuresDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featuresDir)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(featuresDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() || info.Name() != "README.md" {
			return nil
		}
		// README.md must be directly inside a `_tests` directory.
		if filepath.Base(filepath.Dir(path)) != "_tests" {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		fn(path, content)
		return nil
	})
}

// walkScenarioFiles invokes fn for every scenario markdown file under
// specRoot/features/**/_tests/*.md, excluding README.md (which is an index,
// not a scenario).
func walkScenarioFiles(specRoot string, fn func(path string, content []byte)) error {
	featuresDir := filepath.Join(specRoot, "features")
	info, err := os.Stat(featuresDir)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(featuresDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(info.Name(), ".md") {
			return nil
		}
		if info.Name() == "README.md" {
			return nil
		}
		// Path must include a segment named _tests.
		parent := filepath.Base(filepath.Dir(path))
		if parent != "_tests" {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		fn(path, content)
		return nil
	})
}

// walkMatchingFiles enumerates files under root, invoking fn for each file
// where match returns true. depth is measured relative to root (1 = direct
// child). Subdirectories are walked but match is only consulted for files.
func walkMatchingFiles(root string, match func(path string, depth int, name string) bool, fn func(path string, content []byte)) error {
	info, err := os.Stat(root)
	if err != nil || !info.IsDir() {
		return nil
	}
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if path != root && strings.HasPrefix(info.Name(), "_") {
				return filepath.SkipDir
			}
			// Also skip archived/ for ideas.
			if info.Name() == "archived" && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		rel, _ := filepath.Rel(root, path)
		depth := strings.Count(rel, string(filepath.Separator)) + 1
		if !match(path, depth, info.Name()) {
			return nil
		}
		content, readErr := os.ReadFile(path)
		if readErr != nil {
			return nil
		}
		fn(path, content)
		return nil
	})
}
