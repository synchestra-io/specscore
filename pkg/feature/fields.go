package feature

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ValidFields lists all recognized metadata field names.
var ValidFields = map[string]bool{
	"status":    true,
	"oq":        true,
	"deps":      true,
	"refs":      true,
	"children":  true,
	"plans":     true,
	"proposals": true,
}

// ParseFieldNames validates and returns field names from a comma-separated
// string.
func ParseFieldNames(fieldsStr string) ([]string, error) {
	if fieldsStr == "" {
		return nil, nil
	}
	parts := strings.Split(fieldsStr, ",")
	seen := make(map[string]bool)
	var fields []string
	for _, p := range parts {
		f := strings.TrimSpace(p)
		if f == "" {
			continue
		}
		if !ValidFields[f] {
			return nil, fmt.Errorf("unknown field %q (valid: status, oq, deps, refs, children, plans, proposals)", f)
		}
		if !seen[f] {
			seen[f] = true
			fields = append(fields, f)
		}
	}
	return fields, nil
}

// EnrichedFeature holds a feature ID with optional metadata fields.
// Children can be []string (child paths for flat output) or
// []*EnrichedFeature (tree nesting).
type EnrichedFeature struct {
	Path      string      `yaml:"path" json:"path"`
	Focus     *bool       `yaml:"focus,omitempty" json:"focus,omitempty"`
	Cycle     *bool       `yaml:"cycle,omitempty" json:"cycle,omitempty"`
	Status    string      `yaml:"status,omitempty" json:"status,omitempty"`
	OQ        *int        `yaml:"oq,omitempty" json:"oq,omitempty"`
	Deps      []string    `yaml:"deps,omitempty" json:"deps,omitempty"`
	Refs      []string    `yaml:"refs,omitempty" json:"refs,omitempty"`
	Plans     []string    `yaml:"plans,omitempty" json:"plans,omitempty"`
	Proposals []string    `yaml:"proposals,omitempty" json:"proposals,omitempty"`
	Children  interface{} `yaml:"children,omitempty" json:"children,omitempty"`
}

// ResolveFields computes the requested metadata fields for a feature.
func ResolveFields(featuresDir, featureID string, fields []string) *EnrichedFeature {
	ef := &EnrichedFeature{Path: featureID}
	readmePath := ReadmePath(featuresDir, featureID)

	for _, f := range fields {
		switch f {
		case "status":
			if s, err := ParseFeatureStatus(readmePath); err == nil {
				ef.Status = s
			}
		case "oq":
			if n, err := CountOutstandingQuestions(readmePath); err == nil {
				ef.OQ = &n
			}
		case "deps":
			if d, err := ParseDependencies(readmePath); err == nil {
				ef.Deps = d
			}
		case "refs":
			if r, err := FindFeatureRefs(featuresDir, featureID); err == nil {
				ef.Refs = r
			}
		case "children":
			if c, err := DiscoverChildFeatures(featuresDir, featureID, readmePath); err == nil {
				var paths []string
				for _, ch := range c {
					paths = append(paths, ch.Path)
				}
				if len(paths) > 0 {
					ef.Children = paths
				}
			}
		case "plans":
			specRoot := filepath.Dir(featuresDir)
			if p, err := FindLinkedPlans(filepath.Dir(specRoot), featureID); err == nil {
				ef.Plans = p
			}
		case "proposals":
			// Proposals not yet implemented in the spec repo structure.
		}
	}
	return ef
}

// ValidateFormat checks the format flag value is valid.
func ValidateFormat(format string) error {
	if format != "text" && format != "yaml" && format != "json" {
		return fmt.Errorf("invalid format: %s (valid: text, yaml, json)", format)
	}
	return nil
}

// BoolPtr returns a pointer to a bool value.
func BoolPtr(b bool) *bool {
	return &b
}
