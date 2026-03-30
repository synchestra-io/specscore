package sourceref

import (
	"testing"
)

func TestDetectReference(t *testing.T) {
	tests := []struct {
		line string
		want bool
	}{
		{"// synchestra:feature/cli/task", true},
		{"# synchestra:plan/chat-feature", true},
		{"// https://synchestra.io/github.com/synchestra-io/synchestra/spec/features/cli", true},
		{"no reference here", false},
		{"synchestra: not a comment", false},
	}
	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			got := DetectReference(tt.line)
			if got != tt.want {
				t.Errorf("DetectReference(%q) = %v, want %v", tt.line, got, tt.want)
			}
		})
	}
}

func TestParseReference(t *testing.T) {
	tests := []struct {
		input    string
		wantPath string
		wantType string
	}{
		{"synchestra:feature/cli/task", "spec/features/cli/task", "feature"},
		{"synchestra:plan/chat-feature", "spec/plans/chat-feature", "plan"},
		{"synchestra:doc/api", "docs/api", "doc"},
	}
	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			ref, err := ParseReference(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ref.ResolvedPath != tt.wantPath {
				t.Errorf("ResolvedPath = %q, want %q", ref.ResolvedPath, tt.wantPath)
			}
			if ref.Type != tt.wantType {
				t.Errorf("Type = %q, want %q", ref.Type, tt.wantType)
			}
		})
	}
}

func TestScanLine(t *testing.T) {
	ref := ScanLine("// synchestra:feature/cli/task/claim")
	if ref == nil {
		t.Fatal("expected reference, got nil")
	}
	if ref.ResolvedPath != "spec/features/cli/task/claim" {
		t.Errorf("ResolvedPath = %q, want %q", ref.ResolvedPath, "spec/features/cli/task/claim")
	}
}
