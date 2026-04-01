package task

import "testing"

func TestTaskStatusConstants(t *testing.T) {
	statuses := []struct {
		status TaskStatus
		want   string
	}{
		{StatusPlanning, "planning"},
		{StatusQueued, "queued"},
		{StatusClaimed, "claimed"},
		{StatusInProgress, "in_progress"},
		{StatusCompleted, "completed"},
		{StatusFailed, "failed"},
		{StatusBlocked, "blocked"},
		{StatusAborted, "aborted"},
	}
	for _, tt := range statuses {
		if string(tt.status) != tt.want {
			t.Errorf("status %q != %q", tt.status, tt.want)
		}
	}
}

func TestStatusEmojiKnown(t *testing.T) {
	for status := range StatusEmojis {
		emoji := StatusEmoji(status)
		if emoji == "\u2753" {
			t.Errorf("StatusEmoji(%q) returned unknown emoji", status)
		}
	}
}

func TestStatusEmojiUnknown(t *testing.T) {
	got := StatusEmoji(TaskStatus("nonexistent"))
	if got != "\u2753" {
		t.Errorf("StatusEmoji(unknown) = %q, want question mark", got)
	}
}
