package feature

import "strings"

// RelationResolver is a function that returns related feature IDs for a
// given feature.
type RelationResolver func(featuresDir, featureID string) ([]string, error)

// DepsResolver returns the direct dependencies for a feature.
func DepsResolver(featuresDir, featureID string) ([]string, error) {
	return ParseDependencies(ReadmePath(featuresDir, featureID))
}

// RefsResolver returns features that depend on the given feature.
func RefsResolver(featuresDir, featureID string) ([]string, error) {
	return FindFeatureRefs(featuresDir, featureID)
}

// TransitiveDeps follows dependency chains recursively with cycle detection.
func TransitiveDeps(featuresDir, startID string) []*EnrichedFeature {
	visited := map[string]bool{startID: true}
	return walkTransitive(featuresDir, startID, visited, DepsResolver)
}

// TransitiveRefs follows reference chains recursively with cycle detection.
func TransitiveRefs(featuresDir, startID string) []*EnrichedFeature {
	visited := map[string]bool{startID: true}
	return walkTransitive(featuresDir, startID, visited, RefsResolver)
}

func walkTransitive(featuresDir, featureID string, visited map[string]bool, resolver RelationResolver) []*EnrichedFeature {
	related, err := resolver(featuresDir, featureID)
	if err != nil {
		return nil
	}

	var nodes []*EnrichedFeature
	for _, r := range related {
		if visited[r] {
			nodes = append(nodes, &EnrichedFeature{Path: r, Cycle: BoolPtr(true)})
			continue
		}
		visited[r] = true
		node := &EnrichedFeature{Path: r}
		children := walkTransitive(featuresDir, r, visited, resolver)
		if len(children) > 0 {
			node.Children = children
		}
		nodes = append(nodes, node)
	}
	return nodes
}

// EnrichTransitiveNodes adds field metadata to a transitive tree.
func EnrichTransitiveNodes(featuresDir string, nodes []*EnrichedFeature, fields []string) {
	for _, node := range nodes {
		if node.Cycle != nil && *node.Cycle {
			continue
		}
		enrichNodeFields(featuresDir, node, fields)
		if children, ok := node.Children.([]*EnrichedFeature); ok {
			EnrichTransitiveNodes(featuresDir, children, fields)
		}
	}
}

// enrichNodeFields copies resolved field values into a node without
// overwriting its tree children.
func enrichNodeFields(featuresDir string, node *EnrichedFeature, fields []string) {
	resolved := ResolveFields(featuresDir, node.Path, fields)
	node.Status = resolved.Status
	node.OQ = resolved.OQ
	node.Deps = resolved.Deps
	node.Refs = resolved.Refs
	node.Plans = resolved.Plans
	node.Proposals = resolved.Proposals
}

// PrintTransitiveText writes transitive results as indented text.
func PrintTransitiveText(sb *strings.Builder, nodes []*EnrichedFeature, depth int) {
	for _, node := range nodes {
		for i := 0; i < depth; i++ {
			sb.WriteByte('\t')
		}
		sb.WriteString(node.Path)
		if node.Cycle != nil && *node.Cycle {
			sb.WriteString(" (cycle)")
		}
		sb.WriteByte('\n')
		if children, ok := node.Children.([]*EnrichedFeature); ok {
			PrintTransitiveText(sb, children, depth+1)
		}
	}
}
