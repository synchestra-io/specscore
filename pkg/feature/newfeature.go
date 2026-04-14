package feature

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/synchestra-io/specscore/pkg/exitcode"
)

// NewOptions holds the parameters for creating a new feature.
type NewOptions struct {
	Title       string   // Human-readable feature title (required).
	Slug        string   // Feature slug (directory name); auto-generated from Title if empty.
	Parent      string   // Parent feature ID for creating a sub-feature.
	Status      string   // Initial feature status (default "Draft").
	Description string   // Short description for the Summary section.
	DependsOn   []string // Feature IDs this feature depends on.
}

// NewResult describes the outcome of creating a new feature.
type NewResult struct {
	FeatureID    string   // The created feature's full ID.
	FeatureDir   string   // Absolute path to the created feature directory.
	ReadmePath   string   // Absolute path to the created README.md.
	ChangedFiles []string // All files that were created or modified.
	Info         Info     // Metadata for the newly created feature.
}

// New scaffolds a new feature directory with a README template.
// It does NOT perform any git operations — those belong in the CLI layer.
func New(featuresDir string, opts NewOptions) (*NewResult, error) {
	if opts.Title == "" {
		return nil, exitcode.InvalidArgsError("missing required option: Title")
	}

	status := opts.Status
	if status == "" {
		status = "Draft"
	}
	if !IsValidStatus(status) {
		return nil, exitcode.InvalidArgsErrorf("invalid status: %s (supported: Draft, In Progress, Stable, Deprecated)", status)
	}

	// Generate or validate slug.
	slug := opts.Slug
	if slug == "" {
		slug = GenerateSlug(opts.Title)
	} else {
		if err := ValidateSlug(slug); err != nil {
			return nil, exitcode.InvalidArgsErrorf("invalid slug: %v", err)
		}
	}

	// Validate mutual exclusion: --parent vs slash-in-slug.
	if opts.Parent != "" && strings.Contains(slug, "/") {
		return nil, exitcode.InvalidArgsError(
			"cannot use Parent with a slug containing slashes; use one or the other",
		)
	}

	// Validate dependency feature IDs exist.
	for _, dep := range opts.DependsOn {
		if !Exists(featuresDir, dep) {
			return nil, exitcode.InvalidArgsErrorf("dependency feature not found: %s", dep)
		}
	}

	// Resolve full feature path.
	var featureID string
	var parentID string

	switch {
	case opts.Parent != "":
		featureID = opts.Parent + "/" + slug
		parentID = opts.Parent
	case strings.Contains(slug, "/"):
		featureID = slug
		parts := strings.Split(slug, "/")
		parentID = strings.Join(parts[:len(parts)-1], "/")
	default:
		featureID = slug
	}

	featureDir := filepath.Join(featuresDir, filepath.FromSlash(featureID))

	// Validate parent exists (for sub-features).
	if parentID != "" {
		if !Exists(featuresDir, parentID) {
			return nil, exitcode.NotFoundErrorf("parent feature not found: %s", parentID)
		}
	}

	// Verify target doesn't exist.
	if _, err := os.Stat(featureDir); err == nil {
		return nil, exitcode.InvalidStateErrorf("feature already exists at: %s", featureID)
	}

	// Create feature directory and README.
	if err := os.MkdirAll(featureDir, 0o755); err != nil {
		return nil, exitcode.UnexpectedErrorf("creating feature directory: %v", err)
	}

	readme := GenerateReadme(opts.Title, status, opts.Description, opts.DependsOn)
	readmePath := filepath.Join(featureDir, "README.md")
	if err := os.WriteFile(readmePath, []byte(readme), 0o644); err != nil {
		return nil, exitcode.UnexpectedErrorf("writing README.md: %v", err)
	}

	changedFiles := []string{readmePath}

	// Update parent's Contents section (sub-features).
	if parentID != "" {
		parentReadme := ReadmePath(featuresDir, parentID)
		changed, err := UpdateParentContents(parentReadme, filepath.Base(featureDir), opts.Description)
		if err != nil {
			return nil, exitcode.UnexpectedErrorf("updating parent contents: %v", err)
		}
		if changed {
			changedFiles = append(changedFiles, parentReadme)
		}
	}

	// Update feature index for top-level features.
	if parentID == "" {
		indexPath := filepath.Join(featuresDir, "README.md")
		changed, err := UpdateFeatureIndex(indexPath, featureID, opts.Description)
		if err != nil {
			return nil, exitcode.UnexpectedErrorf("updating feature index: %v", err)
		}
		if changed {
			changedFiles = append(changedFiles, indexPath)
		}
	}

	// Build info for the result.
	sections, err := ParseSections(readmePath)
	if err != nil {
		return nil, exitcode.UnexpectedErrorf("parsing sections: %v", err)
	}

	deps := opts.DependsOn
	if deps == nil {
		deps = []string{}
	}

	info := Info{
		Path:     featureID,
		Status:   status,
		Deps:     deps,
		Refs:     []string{},
		Children: nil,
		Plans:    nil,
		Sections: sections,
	}

	return &NewResult{
		FeatureID:    featureID,
		FeatureDir:   featureDir,
		ReadmePath:   readmePath,
		ChangedFiles: changedFiles,
		Info:         info,
	}, nil
}

// UpdateParentContents adds or updates the ## Contents section in the
// parent's README. Returns true if the file was modified.
func UpdateParentContents(parentReadmePath, childSlug, description string) (bool, error) {
	content, err := os.ReadFile(parentReadmePath)
	if err != nil {
		return false, err
	}

	lines := strings.Split(string(content), "\n")
	desc := description
	if desc == "" {
		desc = "TODO: Add description."
	}

	newRow := fmt.Sprintf("| [%s](%s/README.md) | %s |", childSlug, childSlug, desc)

	contentsIdx := -1
	for i, line := range lines {
		if strings.TrimSpace(line) == "## Contents" {
			contentsIdx = i
			break
		}
	}

	if contentsIdx >= 0 {
		insertIdx := contentsIdx + 1
		for insertIdx < len(lines) {
			trimmed := strings.TrimSpace(lines[insertIdx])
			if strings.HasPrefix(trimmed, "## ") && trimmed != "## Contents" {
				break
			}
			insertIdx++
		}
		for insertIdx > contentsIdx+1 && strings.TrimSpace(lines[insertIdx-1]) == "" {
			insertIdx--
		}
		lines = append(lines[:insertIdx+1], lines[insertIdx:]...)
		lines[insertIdx] = newRow
	} else {
		summaryIdx := -1
		for i, line := range lines {
			if strings.TrimSpace(line) == "## Summary" {
				summaryIdx = i
				break
			}
		}

		insertAfter := 0
		if summaryIdx >= 0 {
			insertAfter = summaryIdx + 1
			for insertAfter < len(lines) {
				trimmed := strings.TrimSpace(lines[insertAfter])
				if strings.HasPrefix(trimmed, "## ") {
					break
				}
				insertAfter++
			}
		}

		contentsBlock := []string{
			"## Contents",
			"",
			"| Child | Description |",
			"|---|---|",
			newRow,
			"",
		}

		newLines := make([]string, 0, len(lines)+len(contentsBlock))
		newLines = append(newLines, lines[:insertAfter]...)
		newLines = append(newLines, contentsBlock...)
		newLines = append(newLines, lines[insertAfter:]...)
		lines = newLines
	}

	result := strings.Join(lines, "\n")
	if err := os.WriteFile(parentReadmePath, []byte(result), 0o644); err != nil {
		return false, err
	}
	return true, nil
}

// UpdateFeatureIndex adds a new row to the feature index at
// spec/features/README.md. Returns true if the file was modified.
func UpdateFeatureIndex(indexPath, featureID, description string) (bool, error) {
	content, err := os.ReadFile(indexPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}

	lines := strings.Split(string(content), "\n")
	desc := description
	if desc == "" {
		desc = "TODO: Add description."
	}

	newRow := fmt.Sprintf("| [%s](%s/README.md) | %s |", featureID, featureID, desc)

	tableEnd := -1
	inTable := false
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") && strings.Contains(trimmed, "|") {
			inTable = true
			tableEnd = i
		} else if inTable && trimmed == "" {
			break
		} else if inTable && strings.HasPrefix(trimmed, "## ") {
			break
		}
	}

	if tableEnd >= 0 {
		insertIdx := tableEnd + 1
		lines = append(lines[:insertIdx+1], lines[insertIdx:]...)
		lines[insertIdx] = newRow
	} else {
		lines = append(lines, "", newRow)
	}

	result := strings.Join(lines, "\n")
	if err := os.WriteFile(indexPath, []byte(result), 0o644); err != nil {
		return false, err
	}
	return true, nil
}
