package sourceref

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// ScanResult represents the references found in a set of files.
type ScanResult struct {
	// FileRefs maps file path to list of references found in that file
	FileRefs map[string][]*Reference
}

// ScanFiles scans a list of files for source references.
// Returns a ScanResult with all references grouped by file.
func ScanFiles(filePaths []string) (*ScanResult, error) {
	result := &ScanResult{
		FileRefs: make(map[string][]*Reference),
	}
	for _, filePath := range filePaths {
		refs, err := scanFile(filePath)
		if err != nil {
			continue
		}
		if len(refs) > 0 {
			result.FileRefs[filePath] = refs
		}
	}
	return result, nil
}

func scanFile(filePath string) ([]*Reference, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	seen := make(map[string]bool)
	var refs []*Reference

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		ref := ScanLine(line)
		if ref != nil {
			key := ref.ResolvedPath + ref.CrossRepoSuffix
			if !seen[key] {
				seen[key] = true
				refs = append(refs, ref)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}

	sort.Slice(refs, func(i, j int) bool {
		return refs[i].ResolvedPath+refs[i].CrossRepoSuffix < refs[j].ResolvedPath+refs[j].CrossRepoSuffix
	})
	return refs, nil
}

// ExpandGlobPattern expands a glob pattern to a list of file paths.
// Returns sorted file paths.
func ExpandGlobPattern(pattern string) ([]string, error) {
	if pattern == "" {
		pattern = "**/*"
	}

	// Validate glob pattern first
	if _, err := filepath.Match(pattern, "test"); err != nil && pattern != "**" && pattern != "**/*" {
		_, err := filepath.Match(pattern, "")
		if err != nil {
			return nil, err
		}
	}

	var matches []string
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			return nil
		}

		normalPath := filepath.ToSlash(path)
		if normalPath == "."+string(filepath.Separator) {
			normalPath = normalPath[2:]
		} else if strings.HasPrefix(normalPath, "./") {
			normalPath = normalPath[2:]
		}

		ok, matchErr := matchGlobPattern(normalPath, pattern)
		if matchErr != nil {
			return nil
		}
		if ok {
			matches = append(matches, normalPath)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	sort.Strings(matches)
	return matches, nil
}

// matchGlobPattern matches a file path against a glob pattern.
// Supports * (matches within a path segment) and ** (matches across segments).
func matchGlobPattern(path string, pattern string) (bool, error) {
	if pattern == "**/*" || pattern == "**" {
		return true, nil
	}

	matched, err := filepath.Match(pattern, path)
	if err != nil {
		return false, err
	}
	return matched, nil
}

// GetUniqueReferences extracts unique references from a ScanResult, optionally filtered by type.
// Returns references sorted by (resolved_path, cross_repo_suffix).
func GetUniqueReferences(result *ScanResult, typeFilter string) []*Reference {
	seen := make(map[string]*Reference)

	for _, refs := range result.FileRefs {
		for _, ref := range refs {
			if typeFilter != "" && ref.Type != typeFilter {
				continue
			}
			key := ref.ResolvedPath + ref.CrossRepoSuffix
			if _, exists := seen[key]; !exists {
				seen[key] = ref
			}
		}
	}

	var unique []*Reference
	for _, ref := range seen {
		unique = append(unique, ref)
	}

	sort.Slice(unique, func(i, j int) bool {
		keyI := unique[i].ResolvedPath + unique[i].CrossRepoSuffix
		keyJ := unique[j].ResolvedPath + unique[j].CrossRepoSuffix
		return keyI < keyJ
	})

	return unique
}

// FormatOutput formats the scan results for output.
// If singleFile is true, returns a flat list. Otherwise, groups by file with headers.
func FormatOutput(result *ScanResult, singleFile bool, typeFilter string) string {
	if len(result.FileRefs) == 0 {
		return ""
	}

	var output []string

	if singleFile {
		refs := GetUniqueReferences(result, typeFilter)
		for _, ref := range refs {
			output = append(output, ref.ResolvedPath+ref.CrossRepoSuffix)
		}
	} else {
		fileNames := make([]string, 0, len(result.FileRefs))
		for fname := range result.FileRefs {
			fileNames = append(fileNames, fname)
		}
		sort.Strings(fileNames)

		for i, fname := range fileNames {
			if i > 0 {
				output = append(output, "")
			}
			output = append(output, fname)

			refs := result.FileRefs[fname]
			filtered := refs
			if typeFilter != "" {
				filtered = nil
				for _, ref := range refs {
					if ref.Type == typeFilter {
						filtered = append(filtered, ref)
					}
				}
			}

			for _, ref := range filtered {
				output = append(output, "  "+ref.ResolvedPath+ref.CrossRepoSuffix)
			}
		}
	}

	if len(output) == 0 {
		return ""
	}

	return fmt.Sprintf("%s\n", joinStrings(output))
}

// joinStrings joins strings with newlines.
func joinStrings(strs []string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += "\n"
		}
		result += s
	}
	return result
}
