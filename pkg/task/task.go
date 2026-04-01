// Package task defines the canonical task types, status enum, and related
// structures used by all SpecScore-based coordination tools.
package task

import "time"

// TaskStatus represents the lifecycle state of a task.
type TaskStatus string

const (
	StatusPlanning   TaskStatus = "planning"
	StatusQueued     TaskStatus = "queued"
	StatusClaimed    TaskStatus = "claimed"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusFailed     TaskStatus = "failed"
	StatusBlocked    TaskStatus = "blocked"
	StatusAborted    TaskStatus = "aborted"
)

// Task represents a unit of work. Coordination-only fields (Run, Model, ClaimedAt)
// are NOT included here -- those belong in the coordination layer (synchestra).
type Task struct {
	Slug      string
	Title     string
	Status    TaskStatus
	Parent    string // parent task slug, empty for root tasks
	DependsOn []string
	Requester string
	Reason    string // block/fail/abort reason
	Summary   string // completion summary
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CreateParams holds parameters for creating a new task.
type CreateParams struct {
	Slug      string
	Title     string
	Parent    string
	DependsOn []string
	Requester string
}

// Filter holds optional filters for listing tasks.
// Nil pointer fields mean "don't filter on this field."
type Filter struct {
	Status *TaskStatus
	Parent *string
}

// BoardView represents a rendered task board.
type BoardView struct {
	Rows []BoardRow
}

// BoardRow represents a single row in the task board.
type BoardRow struct {
	Task      string
	Status    TaskStatus
	DependsOn []string
	Branch    string
	Agent     string
	Requester string
	Time      string         // raw time string from board column
	StartedAt *time.Time     // structured start time (populated by callers)
	Duration  *time.Duration // structured duration (populated by callers)
}
