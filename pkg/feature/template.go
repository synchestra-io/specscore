package feature

import "strings"

// ValidStatuses lists the allowed feature lifecycle statuses.
var ValidStatuses = []string{"draft", "approved", "implemented"}

// IsValidStatus reports whether status matches one of the ValidStatuses
// (case-insensitive).
func IsValidStatus(status string) bool {
	for _, s := range ValidStatuses {
		if strings.EqualFold(s, status) {
			return true
		}
	}
	return false
}

// GenerateReadme produces the full README.md content for a new feature.
// If description is empty, a TODO placeholder is used.
// The Dependencies section is only included when deps is non-empty.
func GenerateReadme(title, status, description string, deps []string) string {
	if description == "" {
		description = "TODO: Brief summary of the feature."
	}

	var b strings.Builder

	b.WriteString("# Feature: ")
	b.WriteString(title)
	b.WriteByte('\n')

	b.WriteByte('\n')
	b.WriteString("**Status:** ")
	b.WriteString(status)
	b.WriteByte('\n')

	b.WriteByte('\n')
	b.WriteString("## Summary\n")
	b.WriteByte('\n')
	b.WriteString(description)
	b.WriteByte('\n')

	b.WriteByte('\n')
	b.WriteString("## Problem\n")
	b.WriteByte('\n')
	b.WriteString("TODO: What problem does this feature solve?\n")

	b.WriteByte('\n')
	b.WriteString("## Behavior\n")
	b.WriteByte('\n')
	b.WriteString("TODO: How does this feature work?\n")

	if len(deps) > 0 {
		b.WriteByte('\n')
		b.WriteString("## Dependencies\n")
		b.WriteByte('\n')
		for _, dep := range deps {
			b.WriteString("- ")
			b.WriteString(dep)
			b.WriteByte('\n')
		}
	}

	b.WriteByte('\n')
	b.WriteString("## Acceptance Criteria\n")
	b.WriteByte('\n')
	b.WriteString("TODO: Define acceptance criteria.\n")

	b.WriteByte('\n')
	b.WriteString("## Outstanding Questions\n")
	b.WriteByte('\n')
	b.WriteString("None at this time.\n")

	return b.String()
}
