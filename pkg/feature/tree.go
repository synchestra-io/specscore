package feature

import "strings"

// BuildEnrichedTree builds a tree of EnrichedFeature nodes with resolved
// fields.
func BuildEnrichedTree(featuresDir string, featureIDs []string, fields []string, focusID string) []*EnrichedFeature {
	// Filter out "children" from fields for tree output (tree structure IS children)
	treeFields := make([]string, 0, len(fields))
	for _, f := range fields {
		if f != "children" {
			treeFields = append(treeFields, f)
		}
	}

	nodeMap := make(map[string]*EnrichedFeature)
	var roots []*EnrichedFeature

	for _, id := range featureIDs {
		ef := ResolveFields(featuresDir, id, treeFields)
		if id == focusID && focusID != "" {
			ef.Focus = BoolPtr(true)
		}
		nodeMap[id] = ef

		parts := strings.Split(id, "/")
		if len(parts) == 1 {
			roots = append(roots, ef)
		} else {
			parentID := strings.Join(parts[:len(parts)-1], "/")
			if parent, ok := nodeMap[parentID]; ok {
				if children, ok := parent.Children.([]*EnrichedFeature); ok {
					parent.Children = append(children, ef)
				} else {
					parent.Children = []*EnrichedFeature{ef}
				}
			} else {
				roots = append(roots, ef)
			}
		}
	}

	return roots
}
