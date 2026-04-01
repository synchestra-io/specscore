package task

import (
	"bytes"
	"fmt"
	"strings"
)

// StatusEmojis maps each task status to its emoji prefix for board rendering.
var StatusEmojis = map[TaskStatus]string{
	StatusPlanning:   "\U0001f4cb",
	StatusQueued:     "\u23f3",
	StatusClaimed:    "\U0001f512",
	StatusInProgress: "\U0001f535",
	StatusCompleted:  "\u2705",
	StatusFailed:     "\u274c",
	StatusBlocked:    "\U0001f7e1",
	StatusAborted:    "\u26d4",
}

// StatusEmoji returns the emoji prefix for a task status.
// Returns a question-mark emoji for unknown statuses.
func StatusEmoji(s TaskStatus) string {
	if e, ok := StatusEmojis[s]; ok {
		return e
	}
	return "\u2753"
}

// ParseBoard parses a board markdown file into structured data.
func ParseBoard(data []byte) (*BoardView, error) {
	lines := strings.Split(string(data), "\n")
	var bv BoardView

	// Find the separator line (|---|...).
	sepIdx := -1
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "|") && strings.Contains(trimmed, "---") {
			sepIdx = i
			break
		}
	}
	if sepIdx < 0 {
		return &bv, fmt.Errorf("no table separator found")
	}

	// Parse data rows after the separator.
	for _, line := range lines[sepIdx+1:] {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || !strings.HasPrefix(trimmed, "|") {
			continue
		}
		row, err := parseBoardRow(trimmed)
		if err != nil {
			return &bv, err
		}
		bv.Rows = append(bv.Rows, row)
	}
	return &bv, nil
}

// parseBoardRow parses a single |-delimited table row.
func parseBoardRow(line string) (BoardRow, error) {
	// Split by | and trim; leading/trailing empty cells from outer pipes.
	parts := strings.Split(line, "|")
	var cells []string
	for _, p := range parts {
		cells = append(cells, strings.TrimSpace(p))
	}
	// Remove leading and trailing empty strings from outer pipes.
	if len(cells) > 0 && cells[0] == "" {
		cells = cells[1:]
	}
	if len(cells) > 0 && cells[len(cells)-1] == "" {
		cells = cells[:len(cells)-1]
	}
	if len(cells) != 7 {
		return BoardRow{}, fmt.Errorf("expected 7 columns, got %d", len(cells))
	}

	slug := ExtractSlug(cells[0])

	status, err := ParseStatusCell(cells[1])
	if err != nil {
		return BoardRow{}, err
	}

	return BoardRow{
		Task:      slug,
		Status:    status,
		DependsOn: ParseDeps(cells[2]),
		Branch:    ParseDash(cells[3]),
		Agent:     ParseDash(cells[4]),
		Requester: ParseDash(cells[5]),
		Time:      ParseDashKeep(cells[6]),
	}, nil
}

// ParseStatusCell parses a status cell like "check-mark `completed`" into TaskStatus.
func ParseStatusCell(cell string) (TaskStatus, error) {
	cell = strings.TrimSpace(cell)
	// Extract the status name from between backticks.
	start := strings.IndexByte(cell, '`')
	if start < 0 {
		return "", fmt.Errorf("invalid status cell: %q", cell)
	}
	end := strings.IndexByte(cell[start+1:], '`')
	if end < 0 {
		return "", fmt.Errorf("invalid status cell: %q", cell)
	}
	name := cell[start+1 : start+1+end]
	status := TaskStatus(name)
	if _, ok := StatusEmojis[status]; !ok {
		return "", fmt.Errorf("unknown status: %q", name)
	}
	return status, nil
}

// ExtractSlug extracts the link text from "[slug](slug/)".
func ExtractSlug(cell string) string {
	cell = strings.TrimSpace(cell)
	if start := strings.IndexByte(cell, '['); start >= 0 {
		if end := strings.IndexByte(cell[start:], ']'); end > 0 {
			return cell[start+1 : start+end]
		}
	}
	return cell
}

// ParseDeps parses a comma-separated dependency list, returning nil for em-dash.
func ParseDeps(cell string) []string {
	cell = strings.TrimSpace(cell)
	if cell == "\u2014" || cell == "" {
		return nil
	}
	parts := strings.Split(cell, ",")
	deps := make([]string, 0, len(parts))
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			deps = append(deps, p)
		}
	}
	return deps
}

// ParseDash returns empty string for em-dash, otherwise the trimmed cell.
// Strips surrounding backticks if present (e.g. `agent/run-1`).
func ParseDash(cell string) string {
	cell = strings.TrimSpace(cell)
	// Strip surrounding backticks if present.
	if len(cell) >= 2 && cell[0] == '`' && cell[len(cell)-1] == '`' {
		cell = cell[1 : len(cell)-1]
	}
	if cell == "\u2014" {
		return ""
	}
	return cell
}

// ParseDashKeep returns empty for em-dash, otherwise keeps raw value.
func ParseDashKeep(cell string) string {
	cell = strings.TrimSpace(cell)
	if cell == "\u2014" {
		return ""
	}
	return cell
}

// RenderBoard renders a BoardView to markdown bytes.
func RenderBoard(bv *BoardView) []byte {
	var buf bytes.Buffer

	_, _ = buf.WriteString("# Tasks\n\n")
	_, _ = buf.WriteString("| Task | Status | Depends on | Branch | Agent | Requester | Time |\n")
	_, _ = buf.WriteString("|---|---|---|---|---|---|---|\n")

	for _, r := range bv.Rows {
		taskCell := fmt.Sprintf("[%s](%s/)", r.Task, r.Task)
		status := fmt.Sprintf("%s `%s`", StatusEmoji(r.Status), string(r.Status))
		deps := RenderDash(strings.Join(r.DependsOn, ", "))
		branch := RenderDashBacktick(r.Branch)
		agent := RenderDash(r.Agent)
		requester := RenderDash(r.Requester)
		tm := RenderDash(r.Time)

		_, _ = fmt.Fprintf(&buf, "| %s | %s | %s | %s | %s | %s | %s |\n",
			taskCell, status, deps, branch, agent, requester, tm)
	}

	return buf.Bytes()
}

// RenderDash returns em-dash for empty strings.
func RenderDash(s string) string {
	if s == "" {
		return "\u2014"
	}
	return s
}

// RenderDashBacktick wraps non-empty values in backticks, empty becomes em-dash.
func RenderDashBacktick(s string) string {
	if s == "" {
		return "\u2014"
	}
	return "`" + s + "`"
}
