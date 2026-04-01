package sourceref

import (
	"testing"
)

func TestDetectReference(t *testing.T) {
	tests := []struct {
		line string
		want bool
		name string
	}{
		// Valid detection (different comment prefixes)
		{"// synchestra:feature/cli/task/claim", true, "Go comment with space"},
		{"//synchestra:feature/cli/task/claim", true, "Go comment no space"},
		{"# synchestra:feature/model-selection", true, "Python comment"},
		{"-- https://synchestra.io/github.com/org/repo/spec/features/x", true, "SQL comment with URL"},
		{"; synchestra:plan/v2-migration", true, "Lisp comment"},
		{"%synchestra:doc/api", true, "LaTeX comment no space"},
		{"/* synchestra:feature/x", true, "Block comment start"},
		{"  // synchestra:feature/x", true, "Go comment with leading whitespace"},
		{"// specscore:feature/cli", true, "specscore prefix"},
		{"# specscore:plan/v2", true, "specscore plan prefix"},

		// Invalid (no comment prefix)
		{"synchestra:feature/cli/task/claim", false, "No comment prefix"},
		{"fmt.Println(\"synchestra:feature/x\")", false, "Inside string literal"},
		{"var x = \"https://synchestra.io/...\"", false, "URL in string literal"},
		{"", false, "Empty line"},
		{"// just a comment", false, "Comment without reference"},
		{"// unknown:feature/cli", false, "Unregistered prefix"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := DetectReference(tt.line)
			if got != tt.want {
				t.Errorf("DetectReference(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestExtractReference(t *testing.T) {
	tests := []struct {
		line string
		want string
		name string
	}{
		{"// synchestra:feature/cli/task/claim", "synchestra:feature/cli/task/claim", "Go comment"},
		{"# synchestra:plan/v2", "synchestra:plan/v2", "Python comment"},
		{"-- https://synchestra.io/github.com/org/repo/spec/features/x", "https://synchestra.io/github.com/org/repo/spec/features/x", "SQL with URL"},
		{"// synchestra:feature/cli/task/claim more text", "synchestra:feature/cli/task/claim", "With trailing text"},
		{"  //  synchestra:doc/api/rest  ", "synchestra:doc/api/rest", "Whitespace handling"},
		{"// specscore:feature/cli", "specscore:feature/cli", "specscore prefix"},
		{"", "", "Empty line"},
		{"// just a comment", "", "No reference marker"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ExtractReference(tt.line)
			if got != tt.want {
				t.Errorf("ExtractReference(%q) = %q, want %q", tt.line, got, tt.want)
			}
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := []struct {
		extracted   string
		wantPath    string
		wantType    string
		wantCrossRe string
		wantErr     bool
		name        string
	}{
		// Short notation with type prefixes — synchestra:
		{
			"synchestra:feature/cli/task/claim",
			"spec/features/cli/task/claim",
			"feature",
			"",
			false,
			"Feature shortcut synchestra",
		},
		{
			"synchestra:plan/v2-migration",
			"spec/plans/v2-migration",
			"plan",
			"",
			false,
			"Plan shortcut synchestra",
		},
		{
			"synchestra:doc/api/rest",
			"docs/api/rest",
			"doc",
			"",
			false,
			"Doc shortcut synchestra",
		},

		// Short notation with type prefixes — specscore:
		{
			"specscore:feature/cli/task/claim",
			"spec/features/cli/task/claim",
			"feature",
			"",
			false,
			"Feature shortcut specscore",
		},
		{
			"specscore:plan/v2-migration",
			"spec/plans/v2-migration",
			"plan",
			"",
			false,
			"Plan shortcut specscore",
		},

		// Short notation with full paths (no prefix)
		{
			"synchestra:spec/features/cli/task/claim",
			"spec/features/cli/task/claim",
			"feature",
			"",
			false,
			"Full path to feature",
		},
		{
			"synchestra:README.md",
			"README.md",
			"",
			"",
			false,
			"Full path to file",
		},

		// Cross-repo references
		{
			"synchestra:feature/agent-skills@github.com/acme/orchestrator",
			"spec/features/agent-skills",
			"feature",
			"@github.com/acme/orchestrator",
			false,
			"Feature with cross-repo",
		},
		{
			"synchestra:doc/api/rest@bitbucket.org/acme/docs",
			"docs/api/rest",
			"doc",
			"@bitbucket.org/acme/docs",
			false,
			"Doc with cross-repo",
		},

		// Edge cases
		{
			"synchestra:",
			"",
			"",
			"",
			true,
			"Empty reference",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseReference(tt.extracted)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseReference(%q) error = %v, wantErr %v", tt.extracted, err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if got.ResolvedPath != tt.wantPath {
				t.Errorf("ParseReference(%q).ResolvedPath = %q, want %q", tt.extracted, got.ResolvedPath, tt.wantPath)
			}
			if got.Type != tt.wantType {
				t.Errorf("ParseReference(%q).Type = %q, want %q", tt.extracted, got.Type, tt.wantType)
			}
			if got.CrossRepoSuffix != tt.wantCrossRe {
				t.Errorf("ParseReference(%q).CrossRepoSuffix = %q, want %q", tt.extracted, got.CrossRepoSuffix, tt.wantCrossRe)
			}
		})
	}
}

func TestResolveReference(t *testing.T) {
	tests := []struct {
		ref  string
		want string
		name string
	}{
		{"feature/cli/task/claim", "spec/features/cli/task/claim", "Feature prefix"},
		{"plan/v2", "spec/plans/v2", "Plan prefix"},
		{"doc/api", "docs/api", "Doc prefix"},
		{"spec/features/cli/task", "spec/features/cli/task", "Full path to feature"},
		{"README.md", "README.md", "File path"},
		{"some/arbitrary/path", "some/arbitrary/path", "Arbitrary path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveReference(tt.ref)
			if err != nil {
				t.Errorf("resolveReference(%q) error = %v", tt.ref, err)
				return
			}
			if got != tt.want {
				t.Errorf("resolveReference(%q) = %q, want %q", tt.ref, got, tt.want)
			}
		})
	}
}

func TestInferType(t *testing.T) {
	tests := []struct {
		path string
		want string
		name string
	}{
		{"spec/features/cli/task/claim", "feature", "Feature path"},
		{"spec/features/model-selection", "feature", "Feature path variant"},
		{"spec/plans/v2-migration", "plan", "Plan path"},
		{"docs/api/rest", "doc", "Doc path"},
		{"README.md", "", "File without type"},
		{"arbitrary/path", "", "Arbitrary path"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := inferType(tt.path)
			if got != tt.want {
				t.Errorf("inferType(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}

func TestScanLine(t *testing.T) {
	tests := []struct {
		line string
		want *Reference
		name string
	}{
		{
			"// synchestra:feature/cli/task/claim",
			&Reference{
				ResolvedPath:    "spec/features/cli/task/claim",
				CrossRepoSuffix: "",
				Type:            "feature",
			},
			"Valid feature reference",
		},
		{
			"# synchestra:plan/v2",
			&Reference{
				ResolvedPath:    "spec/plans/v2",
				CrossRepoSuffix: "",
				Type:            "plan",
			},
			"Valid plan reference",
		},
		{
			"// specscore:feature/cli",
			&Reference{
				ResolvedPath:    "spec/features/cli",
				CrossRepoSuffix: "",
				Type:            "feature",
			},
			"specscore prefix feature",
		},
		{
			"// just a comment",
			nil,
			"No reference",
		},
		{
			"",
			nil,
			"Empty line",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ScanLine(tt.line)
			if tt.want == nil {
				if got != nil {
					t.Errorf("ScanLine(%q) = %v, want nil", tt.line, got)
				}
				return
			}
			if got == nil {
				t.Errorf("ScanLine(%q) = nil, want %v", tt.line, tt.want)
				return
			}
			if got.ResolvedPath != tt.want.ResolvedPath {
				t.Errorf("ScanLine(%q).ResolvedPath = %q, want %q", tt.line, got.ResolvedPath, tt.want.ResolvedPath)
			}
			if got.Type != tt.want.Type {
				t.Errorf("ScanLine(%q).Type = %q, want %q", tt.line, got.Type, tt.want.Type)
			}
		})
	}
}

func TestRegisterPrefix(t *testing.T) {
	RegisterPrefix("mytool")
	line := "// mytool:feature/auth"
	if !DetectReference(line) {
		t.Error("should detect registered prefix 'mytool:'")
	}
	ref := ScanLine(line)
	if ref == nil {
		t.Fatal("ScanLine returned nil for registered prefix")
	}
	if ref.ResolvedPath != "spec/features/auth" {
		t.Errorf("expected spec/features/auth, got %s", ref.ResolvedPath)
	}
}

func TestMultiplePrefixes(t *testing.T) {
	tests := []struct {
		line     string
		detected bool
	}{
		{"// specscore:feature/cli", true},
		{"// synchestra:feature/cli", true},
		{"# specscore:plan/v2", true},
		{"// unknown:feature/cli", false},
	}
	for _, tt := range tests {
		if DetectReference(tt.line) != tt.detected {
			t.Errorf("DetectReference(%q) = %v, want %v", tt.line, !tt.detected, tt.detected)
		}
	}
}
