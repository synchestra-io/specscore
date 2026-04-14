package gitremote

import "testing"

func TestParse(t *testing.T) {
	tests := []struct {
		in        string
		wantOK    bool
		wantOwner string
		wantRepo  string
		wantHost  string
	}{
		{"https://github.com/synchestra-io/specscore.git", true, "synchestra-io", "specscore", "github.com"},
		{"https://github.com/synchestra-io/specscore", true, "synchestra-io", "specscore", "github.com"},
		{"http://github.com/o/r.git", true, "o", "r", "github.com"},
		{"https://GITHUB.COM/O/R.git", true, "O", "R", "github.com"},
		{"git@github.com:synchestra-io/specscore.git", true, "synchestra-io", "specscore", "github.com"},
		{"git@github.com:o/r", true, "o", "r", "github.com"},
		{"ssh://git@github.com/synchestra-io/specscore.git", true, "synchestra-io", "specscore", "github.com"},
		{"ssh://git@github.com/o/r", true, "o", "r", "github.com"},
		// Non-GitHub: rejected in MVP.
		{"https://gitlab.com/o/r.git", false, "", "", ""},
		{"git@gitlab.com:o/r.git", false, "", "", ""},
		{"https://bitbucket.org/o/r", false, "", "", ""},
		// Malformed.
		{"", false, "", "", ""},
		{"not-a-url", false, "", "", ""},
		{"https://github.com/only-owner", false, "", "", ""},
	}
	for _, tt := range tests {
		got, ok := Parse(tt.in)
		if ok != tt.wantOK {
			t.Errorf("Parse(%q) ok = %v, want %v", tt.in, ok, tt.wantOK)
			continue
		}
		if !ok {
			continue
		}
		if got.Owner != tt.wantOwner || got.Repo != tt.wantRepo || got.Host != tt.wantHost {
			t.Errorf("Parse(%q) = %+v, want owner=%q repo=%q host=%q",
				tt.in, got, tt.wantOwner, tt.wantRepo, tt.wantHost)
		}
	}
}
