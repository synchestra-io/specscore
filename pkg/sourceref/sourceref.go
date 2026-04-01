package sourceref

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// Reference represents a parsed source reference found in source code.
type Reference struct {
	ResolvedPath    string
	CrossRepoSuffix string
	Type            string
}

// SourceRef represents a source file reference (file + line number).
type SourceRef struct {
	FilePath    string
	LineNumber  int
	LineContent string
}

var (
	mu       sync.Mutex
	prefixes = []string{"specscore", "synchestra"}
	domains  = []string{"specscore.io", "synchestra.io"}

	// DetectionRegex is rebuilt when prefixes change.
	DetectionRegex *regexp.Regexp
)

func init() {
	DetectionRegex = buildDetectionRegex()
}

// RegisterPrefix adds a short-notation prefix (e.g. "mytool") so that
// "mytool:feature/foo" is recognized as a source reference.
// Also registers "mytool.io" as an expanded URL domain.
func RegisterPrefix(prefix string) {
	mu.Lock()
	defer mu.Unlock()
	for _, p := range prefixes {
		if p == prefix {
			return
		}
	}
	prefixes = append(prefixes, prefix)
	domains = append(domains, prefix+".io")
	DetectionRegex = buildDetectionRegex()
}

func buildDetectionRegex() *regexp.Regexp {
	var shortParts []string
	var urlParts []string
	for _, p := range prefixes {
		shortParts = append(shortParts, regexp.QuoteMeta(p+":"))
	}
	for _, d := range domains {
		urlParts = append(urlParts, regexp.QuoteMeta("https://"+d+"/"))
	}
	all := append(shortParts, urlParts...)
	pattern := `^\s*(//|#|--|/\*|\*|%|;)\s*(` + strings.Join(all, "|") + `)`
	return regexp.MustCompile(pattern)
}

// DetectReference checks if a line contains a source reference.
func DetectReference(line string) bool {
	return DetectionRegex.MatchString(line)
}

// ExtractReference extracts the reference string from a line.
func ExtractReference(line string) string {
	for _, p := range prefixes {
		prefix := p + ":"
		if idx := strings.Index(line, prefix); idx != -1 {
			extracted := line[idx:]
			if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
				extracted = extracted[:endIdx]
			}
			return extracted
		}
	}
	for _, d := range domains {
		urlPrefix := "https://" + d + "/"
		if idx := strings.Index(line, urlPrefix); idx != -1 {
			extracted := line[idx:]
			if endIdx := strings.IndexAny(extracted, " \t\n\r"); endIdx != -1 {
				extracted = extracted[:endIdx]
			}
			return extracted
		}
	}
	return ""
}

// ParseReference parses an extracted reference string and returns a Reference.
func ParseReference(extracted string) (*Reference, error) {
	if extracted == "" {
		return nil, fmt.Errorf("empty reference")
	}
	for _, d := range domains {
		urlPrefix := "https://" + d + "/"
		if strings.HasPrefix(extracted, urlPrefix) {
			return parseExpandedURL(extracted, urlPrefix)
		}
	}
	for _, p := range prefixes {
		prefix := p + ":"
		if strings.HasPrefix(extracted, prefix) {
			return parseShortNotation(extracted, prefix)
		}
	}
	return nil, fmt.Errorf("unrecognized reference format: %s", extracted)
}

func parseExpandedURL(url, urlPrefix string) (*Reference, error) {
	path := strings.TrimPrefix(url, urlPrefix)
	parts := strings.Split(path, "/")
	if len(parts) < 4 {
		return nil, fmt.Errorf("invalid expanded URL format: too few path segments")
	}
	host := parts[0]
	org := parts[1]
	repo := parts[2]
	resolvedPath := strings.Join(parts[3:], "/")
	crossRepoSuffix := fmt.Sprintf("@%s/%s/%s", host, org, repo)
	refType := inferType(resolvedPath)
	return &Reference{
		ResolvedPath:    resolvedPath,
		CrossRepoSuffix: crossRepoSuffix,
		Type:            refType,
	}, nil
}

func parseShortNotation(notation, prefix string) (*Reference, error) {
	notation = strings.TrimPrefix(notation, prefix)
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
