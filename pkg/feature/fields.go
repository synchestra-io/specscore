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
// ChildPaths holds child IDs for flat output; ChildNodes holds enriched
// children for tree representations.
type EnrichedFeature struct {
	Path       string             `yaml:"path" json:"path"`
	Focus      *bool              `yaml:"focus,omitempty" json:"focus,omitempty"`
	Cycle      *bool              `yaml:"cycle,omitempty" json:"cycle,omitempty"`
	Status     string             `yaml:"status,omitempty" json:"status,omitempty"`
	OQ         *int               `yaml:"oq,omitempty" json:"oq,omitempty"`
	Deps       []string           `yaml:"deps,omitempty" json:"deps,omitempty"`
	Refs       []string           `yaml:"refs,omitempty" json:"refs,omitempty"`
	Plans      []string           `yaml:"plans,omitempty" json:"plans,omitempty"`
	Proposals  []string           `yaml:"proposals,omitempty" json:"proposals,omitempty"`
	ChildPaths []string           `yaml:"children,omitempty" json:"children,omitempty"`
	ChildNodes []*EnrichedFeature `yaml:"child_nodes,omitempty" json:"child_nodes,omitempty"`
}

// ResolveFields computes the requested metadata fields for a feature.
// Returns partial results alongside any errors encountered.
func ResolveFields(featuresDir, featureID string, fields []string) (*EnrichedFeature, error) {
	ef := &EnrichedFeature{Path: featureID}
	readmePath := ReadmePath(featuresDir, featureID)
	var errs []string

	for _, f := range fields {
		switch f {
		case "status":
			s, err := ParseFeatureStatus(readmePath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("status: %v", err))
			} else {
				ef.Status = s
			}
		case "oq":
			n, err := CountOutstandingQuestions(readmePath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("oq: %v", err))
			} else {
				ef.OQ = &n
			}
		case "deps":
			d, err := ParseDependencies(readmePath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("deps: %v", err))
			} else {
				ef.Deps = d
			}
		case "refs":
			r, err := FindFeatureRefs(featuresDir, featureID)
			if err != nil {
				errs = append(errs, fmt.Sprintf("refs: %v", err))
			} else {
				ef.Refs = r
			}
		case "children":
			c, err := DiscoverChildFeatures(featuresDir, featureID, readmePath)
			if err != nil {
				errs = append(errs, fmt.Sprintf("children: %v", err))
			} else {
				var paths []string
				for _, ch := range c {
					paths = append(paths, ch.Path)
				}
				ef.ChildPaths = paths
			}
		case "plans":
			specRoot := filepath.Dir(featuresDir)
			p, err := FindLinkedPlans(filepath.Dir(specRoot), featureID)
			if err != nil {
				errs = append(errs, fmt.Sprintf("plans: %v", err))
			} else {
				ef.Plans = p
			}
		case "proposals":
			// Proposals not yet implemented in the spec repo structure.
		}
	}
	if len(errs) > 0 {
		return ef, fmt.Errorf("resolve %s: %s", featureID, strings.Join(errs, "; "))
	}
	return ef, nil
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
