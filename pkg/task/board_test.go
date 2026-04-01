package task

import (
	"strings"
	"testing"
)

func TestBoardRoundTrip(t *testing.T) {
	bv := &BoardView{
		Rows: []BoardRow{
			{
				Task:      "setup-db",
				Status:    StatusCompleted,
				DependsOn: nil,
				Branch:    "agent/run-1",
				Agent:     "Sonnet 4.5",
				Requester: "@alice",
				Time:      "2026-03-12 10:15",
			},
			{
				Task:      "implement-api",
				Status:    StatusInProgress,
				DependsOn: []string{"setup-db"},
				Branch:    "agent/run-2",
				Agent:     "Opus 4",
				Requester: "@alex",
				Time:      "2026-03-12 10:22",
			},
			{
				Task:      "write-tests",
				Status:    StatusQueued,
				DependsOn: []string{"implement-api"},
				Branch:    "",
				Agent:     "",
				Requester: "@alex",
				Time:      "",
			},
		},
	}

	rendered := RenderBoard(bv)
	parsed, err := ParseBoard(rendered)
	if err != nil {
		t.Fatalf("ParseBoard failed: %v", err)
	}

	if len(parsed.Rows) != len(bv.Rows) {
		t.Fatalf("row count mismatch: got %d, want %d", len(parsed.Rows), len(bv.Rows))
	}

	for i, want := range bv.Rows {
		got := parsed.Rows[i]
		if got.Task != want.Task {
			t.Errorf("row %d: task = %q, want %q", i, got.Task, want.Task)
		}
		if got.Status != want.Status {
			t.Errorf("row %d: status = %q, want %q", i, got.Status, want.Status)
		}
		if len(got.DependsOn) != len(want.DependsOn) {
			t.Errorf("row %d: deps count = %d, want %d", i, len(got.DependsOn), len(want.DependsOn))
		} else {
			for j := range want.DependsOn {
				if got.DependsOn[j] != want.DependsOn[j] {
					t.Errorf("row %d dep %d: got %q, want %q", i, j, got.DependsOn[j], want.DependsOn[j])
				}
			}
		}
		if got.Branch != want.Branch {
			t.Errorf("row %d: branch = %q, want %q", i, got.Branch, want.Branch)
		}
		if got.Agent != want.Agent {
			t.Errorf("row %d: agent = %q, want %q", i, got.Agent, want.Agent)
		}
		if got.Requester != want.Requester {
			t.Errorf("row %d: requester = %q, want %q", i, got.Requester, want.Requester)
		}
		if got.Time != want.Time {
			t.Errorf("row %d: time = %q, want %q", i, got.Time, want.Time)
		}
	}
}

func TestBoardParseEmpty(t *testing.T) {
	input := []byte("# Tasks\n\n| Task | Status | Depends on | Branch | Agent | Requester | Time |\n|---|---|---|---|---|---|---|\n")
	bv, err := ParseBoard(input)
	if err != nil {
		t.Fatalf("ParseBoard failed: %v", err)
	}
	if len(bv.Rows) != 0 {
		t.Errorf("expected 0 rows, got %d", len(bv.Rows))
	}
}

func TestBoardParseMultipleStatuses(t *testing.T) {
	input := "# Tasks\n\n" +
		"| Task | Status | Depends on | Branch | Agent | Requester | Time |\n" +
		"|---|---|---|---|---|---|---|\n" +
		"| [a](a/) | \U0001f4cb `planning` | \u2014 | \u2014 | \u2014 | \u2014 | \u2014 |\n" +
		"| [b](b/) | \u274c `failed` | a | `br-1` | GPT-5 | @bob | 2026-01-01 |\n" +
		"| [c](c/) | \U0001f7e1 `blocked` | a, b | \u2014 | \u2014 | @carol | \u2014 |\n"

	bv, err := ParseBoard([]byte(input))
	if err != nil {
		t.Fatalf("ParseBoard failed: %v", err)
	}
	if len(bv.Rows) != 3 {
		t.Fatalf("expected 3 rows, got %d", len(bv.Rows))
	}

	if bv.Rows[0].Status != StatusPlanning {
		t.Errorf("row 0: status = %q, want planning", bv.Rows[0].Status)
	}
	if bv.Rows[1].Status != StatusFailed {
		t.Errorf("row 1: status = %q, want failed", bv.Rows[1].Status)
	}
	if bv.Rows[2].Status != StatusBlocked {
		t.Errorf("row 2: status = %q, want blocked", bv.Rows[2].Status)
	}
}

func TestParseStatusCellAllStatuses(t *testing.T) {
	tests := []struct {
		cell string
		want TaskStatus
	}{
		{"\U0001f4cb `planning`", StatusPlanning},
		{"\u23f3 `queued`", StatusQueued},
		{"\U0001f512 `claimed`", StatusClaimed},
		{"\U0001f535 `in_progress`", StatusInProgress},
		{"\u2705 `completed`", StatusCompleted},
		{"\u274c `failed`", StatusFailed},
		{"\U0001f7e1 `blocked`", StatusBlocked},
		{"\u26d4 `aborted`", StatusAborted},
	}

	for _, tt := range tests {
		got, err := ParseStatusCell(tt.cell)
		if err != nil {
			t.Errorf("ParseStatusCell(%q): %v", tt.cell, err)
			continue
		}
		if got != tt.want {
			t.Errorf("ParseStatusCell(%q) = %q, want %q", tt.cell, got, tt.want)
		}
	}
}

func TestParseStatusCellInvalid(t *testing.T) {
	invalid := []string{
		"no backticks",
		"\U0001f535 `unknown_status`",
		"",
		"`",
	}
	for _, cell := range invalid {
		_, err := ParseStatusCell(cell)
		if err == nil {
			t.Errorf("ParseStatusCell(%q): expected error, got nil", cell)
		}
	}
}

func TestRenderBoardStatuses(t *testing.T) {
	bv := &BoardView{
		Rows: []BoardRow{
			{Task: "a", Status: StatusPlanning},
			{Task: "b", Status: StatusAborted},
		},
	}
	out := string(RenderBoard(bv))
	if !strings.Contains(out, "\U0001f4cb `planning`") {
		t.Error("missing planning status in rendered output")
	}
	if !strings.Contains(out, "\u26d4 `aborted`") {
		t.Error("missing aborted status in rendered output")
	}
}

func TestDependsOnParsing(t *testing.T) {
	tests := []struct {
		cell string
		want []string
	}{
		{"\u2014", nil},
		{"setup-db", []string{"setup-db"}},
		{"a, b, c", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		got := ParseDeps(tt.cell)
		if len(got) != len(tt.want) {
			t.Errorf("ParseDeps(%q): got %v, want %v", tt.cell, got, tt.want)
			continue
		}
		for i := range tt.want {
			if got[i] != tt.want[i] {
				t.Errorf("ParseDeps(%q)[%d]: got %q, want %q", tt.cell, i, got[i], tt.want[i])
			}
		}
	}
}

func TestSlugExtraction(t *testing.T) {
	tests := []struct {
		cell string
		want string
	}{
		{"[setup-db](setup-db/)", "setup-db"},
		{"[my-task](my-task/)", "my-task"},
		{"plain-text", "plain-text"},
	}
	for _, tt := range tests {
		got := ExtractSlug(tt.cell)
		if got != tt.want {
			t.Errorf("ExtractSlug(%q) = %q, want %q", tt.cell, got, tt.want)
		}
	}
}
