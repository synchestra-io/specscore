package feature

import (
	"fmt"
	"strings"
	"unicode"
)

// GenerateSlug converts a human-readable title to a URL-safe directory name.
//
// It lowercases the title, replaces spaces and underscores with hyphens,
// removes non-alphanumeric/non-hyphen characters, collapses consecutive
// hyphens, and trims leading/trailing hyphens.
func GenerateSlug(title string) string {
	s := strings.ToLower(title)

	var b strings.Builder
	b.Grow(len(s))
	for _, r := range s {
		switch {
		case r == ' ' || r == '_':
			b.WriteRune('-')
		case unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-':
			b.WriteRune(r)
		}
	}
	s = b.String()

	for strings.Contains(s, "--") {
		s = strings.ReplaceAll(s, "--", "-")
	}

	s = strings.Trim(s, "-")
	return s
}

// ValidateSlug checks that slug is a valid feature slug.
//
// A valid slug is non-empty, lowercase, contains only alphanumeric chars,
// hyphens, and forward slashes (for nested paths). Consecutive hyphens are
// not allowed. Each segment (split by "/") must not have leading or trailing
// hyphens.
func ValidateSlug(slug string) error {
	if slug == "" {
		return fmt.Errorf("slug must not be empty")
	}

	if slug != strings.ToLower(slug) {
		return fmt.Errorf("slug must be lowercase, got %q", slug)
	}

	for _, r := range slug {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '-' && r != '/' {
			return fmt.Errorf("slug contains invalid character %q", r)
		}
	}

	if strings.Contains(slug, "--") {
		return fmt.Errorf("slug must not contain consecutive hyphens")
	}

	segments := strings.Split(slug, "/")
	for _, seg := range segments {
		if seg == "" {
			return fmt.Errorf("slug contains an empty segment (double slash or leading/trailing slash)")
		}
		if strings.HasPrefix(seg, "-") || strings.HasSuffix(seg, "-") {
			return fmt.Errorf("slug segment %q must not have leading or trailing hyphens", seg)
		}
	}

	return nil
}
