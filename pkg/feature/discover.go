// Package feature provides feature discovery, traversal, metadata,
// dependency resolution, and scaffolding. Pure library functions — no cobra,
// no fmt.Print, no os.Exit.
package feature

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/synchestra-io/specscore/pkg/exitcode"
)

// Default directory layout constants.
const (
	DefaultSpecDir    = "spec"
	FeaturesSubDir    = "features"
	specRepoConfigYml = "synchestra-spec-repo.yaml"
)

// Feature holds a discovered feature's identity.
type Feature struct {
	// ID is the slash-separated path relative to the features directory
	// (e.g. "cli/task/claim").
	ID string
}

// FindSpecRepoRoot walks up from startDir looking for
// synchestra-spec-repo.yaml. As a fallback it checks for a spec/features/
// directory. Returns the directory that serves as the spec repo root.
func FindSpecRepoRoot(startDir string) (string, error) {
	current, err := filepath.Abs(startDir)
	if err != nil {
		return "", fmt.Errorf("resolving path: %w", err)
	}

	for {
		configPath := filepath.Join(current, specRepoConfigYml)
		if _, err := os.Stat(configPath); err == nil {
			return current, nil
		}

		featPath := filepath.Join(current, DefaultSpecDir, FeaturesSubDir)
		if info, err := os.Stat(featPath); err == nil && info.IsDir() {
			return current, nil
		}

		parent := filepath.Dir(current)
		if parent == current {
			return "", exitcode.NotFoundError(
				"project not found: no synchestra-spec-repo.yaml or spec/features/ in any parent directory",
			)
		}
		current = parent
	}
}

// Discover walks featuresDir and returns all features sorted alphabetically.
// A feature is any directory containing a README.md. Directories prefixed
// with _ are skipped (reserved for tooling).
func Discover(featuresDir string) ([]Feature, error) {
	var features []Feature

	err := filepath.WalkDir(featuresDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			return nil
		}
		if path == featuresDir {
			return nil
		}
		if strings.HasPrefix(d.Name(), "_") {
			return filepath.SkipDir
		}
		readmePath := filepath.Join(path, "README.md")
		if _, statErr := os.Stat(readmePath); statErr == nil {
			relPath, relErr := filepath.Rel(featuresDir, path)
			if relErr != nil {
				return fmt.Errorf("computing relative path: %w", relErr)
			}
			featureID := filepath.ToSlash(relPath)
			features = append(features, Feature{ID: featureID})
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walking features directory: %w", err)
	}

	sort.Slice(features, func(i, j int) bool {
		return features[i].ID < features[j].ID
	})
	return features, nil
}

// Exists checks if a feature ID corresponds to a valid feature directory
// with README.md.
func Exists(featuresDir, featureID string) bool {
	readmePath := filepath.Join(featuresDir, filepath.FromSlash(featureID), "README.md")
	_, err := os.Stat(readmePath)
	return err == nil
}

// ReadmePath returns the absolute path to a feature's README.md.
func ReadmePath(featuresDir, featureID string) string {
	return filepath.Join(featuresDir, filepath.FromSlash(featureID), "README.md")
}

// FeatureNode represents a feature in a tree structure.
type FeatureNode struct {
	Name     string
	ID       string // full feature ID
	Focus    bool   // marked as target in focused tree
	Children []*FeatureNode
}

// BuildTree builds a tree from a sorted list of feature IDs.
func BuildTree(featureIDs []string) []*FeatureNode {
	roots := make([]*FeatureNode, 0)
	nodeMap := make(map[string]*FeatureNode)

	for _, id := range featureIDs {
		parts := strings.Split(id, "/")
		name := parts[len(parts)-1]
		node := &FeatureNode{Name: name, ID: id}
		nodeMap[id] = node

		if len(parts) == 1 {
			roots = append(roots, node)
		} else {
			parentID := strings.Join(parts[:len(parts)-1], "/")
			if parent, ok := nodeMap[parentID]; ok {
				parent.Children = append(parent.Children, node)
			} else {
				roots = append(roots, node)
			}
		}
	}

	return roots
}

// PrintTree writes the tree to a strings.Builder with tab indentation.
// Nodes with Focus set are prefixed with "* ".
func PrintTree(w *strings.Builder, nodes []*FeatureNode, depth int) {
	for _, node := range nodes {
		for i := 0; i < depth; i++ {
			w.WriteByte('\t')
		}
		if node.Focus {
			w.WriteString("* ")
		}
		w.WriteString(node.Name)
		w.WriteByte('\n')
		PrintTree(w, node.Children, depth+1)
	}
}

// FilterFocusedFeatures returns features relevant to a focused tree view.
func FilterFocusedFeatures(allFeatures []string, targetID, direction string) []string {
	include := make(map[string]bool)
	include[targetID] = true

	if direction != "down" {
		parts := strings.Split(targetID, "/")
		for i := 1; i < len(parts); i++ {
			ancestor := strings.Join(parts[:i], "/")
			include[ancestor] = true
		}
	}

	if direction != "up" {
		prefix := targetID + "/"
		for _, f := range allFeatures {
			if strings.HasPrefix(f, prefix) {
				include[f] = true
			}
		}
	}

	var filtered []string
	for _, f := range allFeatures {
		if include[f] {
			filtered = append(filtered, f)
		}
	}
	return filtered
}

// MarkFocus sets the focus flag on the target node in the tree.
func MarkFocus(nodes []*FeatureNode, targetID string) {
	for _, node := range nodes {
		if node.ID == targetID {
			node.Focus = true
			return
		}
		MarkFocus(node.Children, targetID)
	}
}

// ParseDependencies reads a feature's README.md and extracts the
// ## Dependencies section. Returns a sorted list of feature IDs listed as
// bullet items.
//
// Supports two formats:
//   - bare ID:       "- claim-and-push"
//   - markdown link: "- [Name](../path/README.md) -- optional description"
func ParseDependencies(readmePath string) ([]string, error) {
	f, err := os.Open(readmePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = f.Close() }()

	var deps []string
	inDeps := false
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)

		if trimmed == "## Dependencies" {
			inDeps = true
			continue
		}
		if inDeps && strings.HasPrefix(trimmed, "## ") {
			break
		}
		if inDeps && strings.HasPrefix(trimmed, "- ") {
			item := strings.TrimSpace(strings.TrimPrefix(trimmed, "- "))
			if item == "" {
				continue
			}
			dep := ExtractFeatureID(item)
			if dep != "" {
				deps = append(deps, dep)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Strings(deps)
	return deps, nil
}

// ExtractFeatureID extracts a feature ID from a dependency list item.
func ExtractFeatureID(item string) string {
	if strings.HasPrefix(item, "[") {
		closeBracket := strings.Index(item, "](")
		if closeBracket == -1 {
			return item
		}
		rest := item[closeBracket+2:]
		closeParen := strings.Index(rest, ")")
		if closeParen == -1 {
			return item
		}
		linkPath := rest[:closeParen]
		return FeatureIDFromRelativePath(linkPath)
	}
	if idx := strings.Index(item, " —"); idx != -1 {
		item = strings.TrimSpace(item[:idx])
	} else if idx := strings.Index(item, " - "); idx != -1 {
		item = strings.TrimSpace(item[:idx])
	}
	return item
}

// FeatureIDFromRelativePath converts a relative path like "../cli/README.md"
// to a feature ID like "cli".
func FeatureIDFromRelativePath(relPath string) string {
	relPath = strings.ReplaceAll(relPath, "\\", "/")
	relPath = strings.TrimSuffix(relPath, "/README.md")
	relPath = strings.TrimSuffix(relPath, "/readme.md")
	parts := strings.Split(relPath, "/")
	var clean []string
	for _, p := range parts {
		if p == ".." || p == "." || p == "" {
			continue
		}
		clean = append(clean, p)
	}
	return strings.Join(clean, "/")
}

// FeatureIDs is a convenience helper that extracts just the ID strings
// from a slice of Feature values.
func FeatureIDs(features []Feature) []string {
	ids := make([]string, len(features))
	for i, f := range features {
		ids[i] = f.ID
	}
	return ids
}
