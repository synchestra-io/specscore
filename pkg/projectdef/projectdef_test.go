package projectdef

import (
	"os"
	"path/filepath"
	"testing"
)

func TestSpecConfigRoundTrip(t *testing.T) {
	dir := t.TempDir()
	cfg := SpecConfig{
		Title:     "Test Project",
		StateRepo: "https://github.com/test/state.git",
		Repos:     []string{"https://github.com/test/code.git"},
	}
	if err := WriteSpecConfig(dir, cfg); err != nil {
		t.Fatalf("WriteSpecConfig: %v", err)
	}
	got, err := ReadSpecConfig(dir)
	if err != nil {
		t.Fatalf("ReadSpecConfig: %v", err)
	}
	if got.Title != cfg.Title {
		t.Errorf("Title = %q, want %q", got.Title, cfg.Title)
	}
	if got.StateRepo != cfg.StateRepo {
		t.Errorf("StateRepo = %q, want %q", got.StateRepo, cfg.StateRepo)
	}
}

func TestParseStateRepo(t *testing.T) {
	tests := []struct {
		stateRepo  string
		wantMode   string
		wantBranch string
	}{
		{"worktree://synchestra-state", "worktree", "synchestra-state"},
		{"https://github.com/test/state.git", "repo", ""},
		{"", "", ""},
	}
	for _, tt := range tests {
		cfg := SpecConfig{StateRepo: tt.stateRepo}
		mode, branch := cfg.ParseStateRepo()
		if mode != tt.wantMode || branch != tt.wantBranch {
			t.Errorf("ParseStateRepo(%q) = (%q, %q), want (%q, %q)",
				tt.stateRepo, mode, branch, tt.wantMode, tt.wantBranch)
		}
	}
}

func TestSpecConfigFileExists(t *testing.T) {
	dir := t.TempDir()
	cfg := SpecConfig{Title: "t"}
	if err := WriteSpecConfig(dir, cfg); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, SpecConfigFile)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("expected %s to exist", path)
	}
}
