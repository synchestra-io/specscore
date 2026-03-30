package sourceref

import (
	"fmt"
	"regexp"
	"strings"
)

// Reference represents a parsed source reference found in source code.
type Reference struct {
	// ResolvedPath is the repo-root-relative path after type-prefix expansion
	ResolvedPath string
	// CrossRepoSuffix is the optional @host/org/repo (empty string if same-repo)
	CrossRepoSuffix string
	// Type is the inferred resource type: "feature", "plan", "doc", or "" if unknown
	Type string
}

// SourceRef represents a source file reference (file + line number).
type SourceRef struct {
	FilePath    string
	LineNumber  int
	LineContent string
}

// DetectionRegex matches source references preceded by recognized comment prefixes.
var DetectionRegex = regexp.MustCompile(`^\s*(//|#|--|/\*|\*|%|;)\s*(synchestra:|https://synchestra\.io/)`)

// DetectReference checks if a line contains a source reference.
func DetectReference(line string) bool {
	return DetectionRegex.MatchString(line)
}

// ExtractReference extracts the reference string from a line.
func ExtractReference(line string) string {
	idx := strings.Index(line, "synchestra:")
	if idx == -1 {
		idx = strings.Index(line, "https://synchestra.io/")
	}
	if idx == -1 {
		return ""
	}
	extracted := line[idx:]
	if strings.HasPrefix(extracted, "https://") {
		if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
			extracted = extracted[:endIdx]
		}
	} else if strings.HasPrefix(extracted, "synchestra:") {
		if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
			extracted = extracted[:endIdx]
		}
	}
	return extracted
}

// ParseReference parses an extracted reference string and returns a Reference.
func ParseReference(extracted string) (*Reference, error) {
	if extracted == "" {
		return nil, fmt.Errorf("empty reference")
	}
	if strings.HasPrefix(extracted, "https://synchestra.io/") {
		return parseExpandedURL(extracted)
	}
	if strings.HasPrefix(extracted, "synchestra:") {
		return parseShortNotation(extracted)
	}
	return nil, fmt.Errorf("unrecognized reference format: %s", extracted)
}

func parseExpandedURL(url string) (*Reference, error) {
	url = strings.TrimPrefix(url, "https://synchestra.io/")
	parts := strings.Split(url, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid expanded URL format: too few path segments")
	}
	host := parts[0]
	org := parts[1]
	repo := parts[2]
	resolvedPath := strings.Join(parts[3:], "/")
	currentHost, currentOrg, currentRepo := "github.com", "synchestra-io", "synchestra"
	crossRepoSuffix := ""
	if host != currentHost || org != currentOrg || repo != currentRepo {
		crossRepoSuffix = fmt.Sprintf("@%s/%s/%s", host, org, repo)
	}
	refType := inferType(resolvedPath)
	return &Reference{
		ResolvedPath:    resolvedPath,
		CrossRepoSuffix: crossRepoSuffix,
		Type:            refType,
	}, nil
}

func parseShortNotation(notation string) (*Reference, error) {
	notation = strings.TrimPrefix(notation, "synchestra:")
	crossRepoSuffix := ""
	reference := notation
	if idx := strings.LastIndex(notation, "@"); idx != -1 {
		crossRepoSuffix = notation[idx:]
		reference = notation[:idx]
	}
	resolvedPath, err := resolveReference(reference)
	if err != nil {
		return nil, err
	}
	refType := inferType(resolvedPath)
	return &Reference{
		ResolvedPath:    resolvedPath,
		CrossRepoSuffix: crossRepoSuffix,
		Type:            refType,
	}, nil
}

func resolveReference(ref string) (string, error) {
	if ref == "" {
		return "", fmt.Errorf("empty reference")
	}
	if strings.HasPrefix(ref, "feature/") {
		return "spec/features/" + strings.TrimPrefix(ref, "feature/"), nil
	}
	if strings.HasPrefix(ref, "plan/") {
		return "spec/plans/" + strings.TrimPrefix(ref, "plan/"), nil
	}
	if strings.HasPrefix(ref, "doc/") {
		return "docs/" + strings.TrimPrefix(ref, "doc/"), nil
	}
	return ref, nil
}

func inferType(resolvedPath string) string {
	if strings.HasPrefix(resolvedPath, "spec/features/") {
		return "feature"
	}
	if strings.HasPrefix(resolvedPath, "spec/plans/") {
		return "plan"
	}
	if strings.HasPrefix(resolvedPath, "docs/") {
		return "doc"
	}
	return ""
}

// ScanLine scans a single line for references. Returns nil if none found.
func ScanLine(line string) *Reference {
	if !DetectReference(line) {
		return nil
	}
	extracted := ExtractReference(line)
	if extracted == "" {
		return nil
	}
	ref, err := ParseReference(extracted)
	if err != nil {
		return nil
	}
	return ref
}
