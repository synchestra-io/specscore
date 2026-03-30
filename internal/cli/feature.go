package cli

// Features implemented: cli/feature

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/feature"
	"gopkg.in/yaml.v3"
)

// featureCommand returns the "feature" command group.
func featureCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "feature",
		Short: "Query features — listing, hierarchy, dependencies, references",
	}
	cmd.AddCommand(
		featureInfoCommand(),
		featureListCommand(),
		featureTreeCommand(),
		featureDepsCommand(),
		featureRefsCommand(),
		featureNewCommand(),
	)
	return cmd
}

// resolveFeaturesDir resolves the features directory from a --project flag or CWD.
func resolveFeaturesDir(projectFlag string) (string, error) {
	var startDir string
	if projectFlag != "" {
		abs, err := filepath.Abs(projectFlag)
		if err != nil {
			return "", exitcode.InvalidArgsErrorf("resolving --project path: %v", err)
		}
		startDir = abs
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return "", exitcode.UnexpectedErrorf("cannot determine working directory: %v", err)
		}
		startDir = cwd
	}

	root, err := feature.FindSpecRepoRoot(startDir)
	if err != nil {
		return "", err
	}

	featDir := filepath.Join(root, "spec", "features")
	info, err := os.Stat(featDir)
	if err != nil || !info.IsDir() {
		return "", exitcode.NotFoundErrorf("features directory not found: %s", featDir)
	}
	return featDir, nil
}

// effectiveFormat determines the output format from flags.
// With --fields and no explicit --format, auto-switches to YAML.
func effectiveFormat(cmd *cobra.Command) string {
	format, _ := cmd.Flags().GetString("format")
	if format != "" {
		return format
	}
	fields, _ := cmd.Flags().GetString("fields")
	if fields != "" {
		return "yaml"
	}
	return "text"
}

// validateFormat checks the format flag value is valid.
func validateFormat(fmt string) error {
	if fmt != "text" && fmt != "yaml" && fmt != "json" {
		return exitcode.InvalidArgsErrorf("invalid --format: %s (valid: text, yaml, json)", fmt)
	}
	return nil
}

// writeEnrichedOutput writes enriched features in the specified format.
func writeEnrichedOutput(w io.Writer, features []*feature.EnrichedFeature, fields []string, format string) error {
	switch format {
	case "yaml":
		return writeEnrichedYAML(w, features)
	case "json":
		return writeEnrichedJSON(w, features)
	default:
		return writeEnrichedText(w, features, fields)
	}
}

func writeEnrichedYAML(w io.Writer, features []*feature.EnrichedFeature) error {
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(features); err != nil {
		return err
	}
	return enc.Close()
}

func writeEnrichedJSON(w io.Writer, features []*feature.EnrichedFeature) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(features)
}

func writeEnrichedText(w io.Writer, features []*feature.EnrichedFeature, fields []string) error {
	bw := bufio.NewWriter(w)
	for _, ef := range features {
		writeEnrichedTextNode(bw, ef, fields, 0)
	}
	return bw.Flush()
}

func writeEnrichedTextNode(w *bufio.Writer, ef *feature.EnrichedFeature, fields []string, depth int) {
	indent := strings.Repeat("\t", depth)
	prefix := ""
	if ef.Focus != nil && *ef.Focus {
		prefix = "* "
	}
	if ef.Cycle != nil && *ef.Cycle {
		_, _ = fmt.Fprintf(w, "%s%s (cycle)\n", indent, ef.Path)
		return
	}

	var meta []string
	for _, f := range fields {
		switch f {
		case "status":
			if ef.Status != "" {
				meta = append(meta, fmt.Sprintf("status=%s", ef.Status))
			}
		case "oq":
			if ef.OQ != nil {
				meta = append(meta, fmt.Sprintf("oq=%d", *ef.OQ))
			}
		case "deps":
			if len(ef.Deps) > 0 {
				meta = append(meta, fmt.Sprintf("deps=[%s]", strings.Join(ef.Deps, ",")))
			}
		case "refs":
			if len(ef.Refs) > 0 {
				meta = append(meta, fmt.Sprintf("refs=[%s]", strings.Join(ef.Refs, ",")))
			}
		case "plans":
			if len(ef.Plans) > 0 {
				meta = append(meta, fmt.Sprintf("plans=[%s]", strings.Join(ef.Plans, ",")))
			}
		case "proposals":
			if len(ef.Proposals) > 0 {
				meta = append(meta, fmt.Sprintf("proposals=[%s]", strings.Join(ef.Proposals, ",")))
			}
		}
	}

	suffix := ""
	if len(meta) > 0 {
		suffix = " " + strings.Join(meta, " ")
	}
	_, _ = fmt.Fprintf(w, "%s%s%s%s\n", indent, prefix, ef.Path, suffix)

	for _, child := range ef.ChildNodes {
		writeEnrichedTextNode(w, child, fields, depth+1)
	}
}

// printTransitiveText writes transitive results as indented text.
func printTransitiveText(sb *strings.Builder, nodes []*feature.EnrichedFeature, depth int) {
	for _, node := range nodes {
		for i := 0; i < depth; i++ {
			sb.WriteByte('\t')
		}
		sb.WriteString(node.Path)
		if node.Cycle != nil && *node.Cycle {
			sb.WriteString(" (cycle)")
		}
		sb.WriteByte('\n')
		if len(node.ChildNodes) > 0 {
			printTransitiveText(sb, node.ChildNodes, depth+1)
		}
	}
}

// --- feature info ---

func featureInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "info <feature_id>",
		Short: "Show feature metadata, section TOC, and children",
		Long: `Returns structured metadata and a section table-of-contents with line
ranges for a feature's README.md. Default output is YAML; use --format for JSON or text.`,
		Args: cobra.ExactArgs(1),
		RunE: runFeatureInfo,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("format", "yaml", "output format: yaml, json, text")
	return cmd
}

func runFeatureInfo(cmd *cobra.Command, args []string) error {
	featureID := args[0]
	projectFlag, _ := cmd.Flags().GetString("project")
	formatFlag, _ := cmd.Flags().GetString("format")

	if formatFlag != "yaml" && formatFlag != "json" && formatFlag != "text" {
		return exitcode.InvalidArgsErrorf("invalid format: %s (supported: yaml, json, text)", formatFlag)
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	if !feature.Exists(featuresDir, featureID) {
		return exitcode.NotFoundErrorf("feature not found: %s", featureID)
	}

	info, err := feature.GetInfo(featuresDir, featureID)
	if err != nil {
		return exitcode.UnexpectedErrorf("getting feature info: %v", err)
	}

	return writeFeatureInfo(cmd.OutOrStdout(), formatFlag, info)
}

func writeFeatureInfo(w io.Writer, format string, info *feature.Info) error {
	switch format {
	case "yaml":
		enc := yaml.NewEncoder(w)
		enc.SetIndent(2)
		if err := enc.Encode(info); err != nil {
			return exitcode.UnexpectedErrorf("encoding yaml: %v", err)
		}
		return enc.Close()
	case "json":
		enc := json.NewEncoder(w)
		enc.SetIndent("", "  ")
		return enc.Encode(info)
	case "text":
		return writeTextInfo(w, info)
	}
	return nil
}

func writeTextInfo(w io.Writer, info *feature.Info) error {
	bw := bufio.NewWriter(w)

	_, _ = fmt.Fprintf(bw, "Feature: %s\n", info.Path)
	_, _ = fmt.Fprintf(bw, "Status:  %s\n", info.Status)

	if len(info.Deps) > 0 {
		_, _ = fmt.Fprintf(bw, "Deps:    %s\n", strings.Join(info.Deps, ", "))
	} else {
		_, _ = fmt.Fprintln(bw, "Deps:    (none)")
	}

	if len(info.Refs) > 0 {
		_, _ = fmt.Fprintf(bw, "Refs:    %s\n", strings.Join(info.Refs, ", "))
	} else {
		_, _ = fmt.Fprintln(bw, "Refs:    (none)")
	}

	if len(info.Children) > 0 {
		_, _ = fmt.Fprintln(bw, "\nChildren:")
		for _, c := range info.Children {
			marker := "✓"
			if !c.InReadme {
				marker = "✗"
			}
			_, _ = fmt.Fprintf(bw, "  %s %s (in_readme: %v)\n", marker, c.Path, c.InReadme)
		}
	}

	if len(info.Plans) > 0 {
		_, _ = fmt.Fprintf(bw, "\nPlans:   %s\n", strings.Join(info.Plans, ", "))
	}

	if len(info.Sections) > 0 {
		_, _ = fmt.Fprintln(bw, "\nSections:")
		printTextSections(bw, info.Sections, 1)
	}

	return bw.Flush()
}

func printTextSections(w *bufio.Writer, sections []feature.SectionInfo, depth int) {
	indent := strings.Repeat("  ", depth)
	for _, s := range sections {
		itemsSuffix := ""
		if s.Items > 0 {
			itemsSuffix = fmt.Sprintf(" (%d items)", s.Items)
		}
		_, _ = fmt.Fprintf(w, "%s%s [%s]%s\n", indent, s.Title, s.Lines, itemsSuffix)
		if len(s.Children) > 0 {
			printTextSections(w, s.Children, depth+1)
		}
	}
}

// --- feature list ---

func featureListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List all feature IDs, one per line",
		Long: `Lists all features in a project as full feature IDs, one per line,
sorted alphabetically. Use --fields to include metadata for each feature.`,
		Args: cobra.NoArgs,
		RunE: runFeatureList,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("fields", "", "comma-separated metadata fields to include (e.g., status,oq)")
	cmd.Flags().String("format", "", "output format: yaml, json, text (auto-selects yaml when --fields is set)")
	return cmd
}

func runFeatureList(cmd *cobra.Command, _ []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")
	fieldsFlag, _ := cmd.Flags().GetString("fields")

	fields, err := feature.ParseFieldNames(fieldsFlag)
	if err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}

	format := effectiveFormat(cmd)
	if err := validateFormat(format); err != nil {
		return err
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	features, err := feature.Discover(featuresDir)
	if err != nil {
		return exitcode.UnexpectedErrorf("discovering features: %v", err)
	}
	featureIDs := feature.FeatureIDs(features)

	w := cmd.OutOrStdout()

	if len(fields) > 0 || format == "yaml" || format == "json" {
		var enriched []*feature.EnrichedFeature
		for _, id := range featureIDs {
			ef, _ := feature.ResolveFields(featuresDir, id, fields)
			enriched = append(enriched, ef)
		}
		return writeEnrichedOutput(w, enriched, fields, format)
	}

	for _, id := range featureIDs {
		_, _ = fmt.Fprintln(w, id)
	}
	return nil
}

// --- feature tree ---

func featureTreeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "tree [feature_id]",
		Short: "Display feature hierarchy as an indented tree",
		Long: `Displays the feature hierarchy as an indented tree. Without a feature ID,
shows the full project tree. With a feature ID, shows the feature in context —
ancestors (path to root) plus its subtree by default. Use --direction to narrow
to ancestors only (up) or subtree only (down).`,
		Args: cobra.MaximumNArgs(1),
		RunE: runFeatureTree,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("direction", "", "up (ancestors only) or down (subtree only); requires feature_id")
	cmd.Flags().String("fields", "", "comma-separated metadata fields to include (e.g., status,oq)")
	cmd.Flags().String("format", "", "output format: yaml, json, text (auto-selects yaml when --fields is set)")
	return cmd
}

func runFeatureTree(cmd *cobra.Command, args []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")
	directionFlag, _ := cmd.Flags().GetString("direction")
	fieldsFlag, _ := cmd.Flags().GetString("fields")

	if directionFlag != "" && directionFlag != "up" && directionFlag != "down" {
		return exitcode.InvalidArgsErrorf("invalid --direction: %s (valid: up, down)", directionFlag)
	}
	if directionFlag != "" && len(args) == 0 {
		return exitcode.InvalidArgsError("--direction requires a feature_id argument")
	}

	fields, err := feature.ParseFieldNames(fieldsFlag)
	if err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}

	format := effectiveFormat(cmd)
	if err := validateFormat(format); err != nil {
		return err
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	features, err := feature.Discover(featuresDir)
	if err != nil {
		return exitcode.UnexpectedErrorf("discovering features: %v", err)
	}
	featureIDs := feature.FeatureIDs(features)

	w := cmd.OutOrStdout()

	var targetID string
	if len(args) > 0 {
		targetID = args[0]
		if !feature.Exists(featuresDir, targetID) {
			return exitcode.NotFoundErrorf("feature not found: %s", targetID)
		}
	}

	if len(fields) > 0 || format == "yaml" || format == "json" {
		var filtered []string
		if targetID == "" {
			filtered = featureIDs
		} else {
			filtered = feature.FilterFocusedFeatures(featureIDs, targetID, directionFlag)
		}
		roots := feature.BuildEnrichedTree(featuresDir, filtered, fields, targetID)
		switch format {
		case "yaml":
			return writeEnrichedYAML(w, roots)
		case "json":
			return writeEnrichedJSON(w, roots)
		default:
			return writeEnrichedText(w, roots, fields)
		}
	}

	if targetID == "" {
		roots := feature.BuildTree(featureIDs)
		var sb strings.Builder
		feature.PrintTree(&sb, roots, 0)
		_, _ = fmt.Fprint(w, sb.String())
		return nil
	}

	filtered := feature.FilterFocusedFeatures(featureIDs, targetID, directionFlag)
	roots := feature.BuildTree(filtered)
	feature.MarkFocus(roots, targetID)
	var sb strings.Builder
	feature.PrintTree(&sb, roots, 0)
	_, _ = fmt.Fprint(w, sb.String())
	return nil
}

// --- feature deps ---

func featureDepsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps <feature_id>",
		Short: "Show features that a given feature depends on",
		Long: `Shows the features that a given feature depends on. Dependencies are
read from the ## Dependencies section in the feature's README.md. Use --transitive
to follow the full dependency chain. Use --fields to include metadata.`,
		Args: cobra.ExactArgs(1),
		RunE: runFeatureDeps,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("fields", "", "comma-separated metadata fields to include (e.g., status,oq)")
	cmd.Flags().String("format", "", "output format: yaml, json, text (auto-selects yaml when --fields is set)")
	cmd.Flags().Bool("transitive", false, "follow dependency chain recursively")
	return cmd
}

func runFeatureDeps(cmd *cobra.Command, args []string) error {
	featureID := args[0]
	projectFlag, _ := cmd.Flags().GetString("project")
	fieldsFlag, _ := cmd.Flags().GetString("fields")
	transitive, _ := cmd.Flags().GetBool("transitive")

	fields, err := feature.ParseFieldNames(fieldsFlag)
	if err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}

	format := effectiveFormat(cmd)
	if err := validateFormat(format); err != nil {
		return err
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	if !feature.Exists(featuresDir, featureID) {
		return exitcode.NotFoundErrorf("feature not found: %s", featureID)
	}

	w := cmd.OutOrStdout()

	if transitive {
		nodes := feature.TransitiveDeps(featuresDir, featureID)
		if len(fields) > 0 {
			feature.EnrichTransitiveNodes(featuresDir, nodes, fields)
		}
		switch format {
		case "yaml":
			return writeEnrichedYAML(w, nodes)
		case "json":
			return writeEnrichedJSON(w, nodes)
		default:
			if len(fields) > 0 {
				return writeEnrichedText(w, nodes, fields)
			}
			var sb strings.Builder
			printTransitiveText(&sb, nodes, 0)
			_, _ = fmt.Fprint(w, sb.String())
		}
		return nil
	}

	readmePath := feature.ReadmePath(featuresDir, featureID)
	deps, err := feature.ParseDependencies(readmePath)
	if err != nil {
		return exitcode.UnexpectedErrorf("reading feature %s: %v", featureID, err)
	}

	if len(fields) > 0 || format == "yaml" || format == "json" {
		var enriched []*feature.EnrichedFeature
		for _, dep := range deps {
			ef, _ := feature.ResolveFields(featuresDir, dep, fields)
			enriched = append(enriched, ef)
		}
		return writeEnrichedOutput(w, enriched, fields, format)
	}

	errW := cmd.ErrOrStderr()
	for _, dep := range deps {
		if !feature.Exists(featuresDir, dep) {
			_, _ = fmt.Fprintf(errW, "%s (not found)\n", dep)
		} else {
			_, _ = fmt.Fprintln(w, dep)
		}
	}
	return nil
}

// --- feature refs ---

func featureRefsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "refs <feature_id>",
		Short: "Show features that reference a given feature as a dependency",
		Long: `Shows features that reference (depend on) a given feature. This is the
inverse of deps — it scans all features' ## Dependencies sections to find
those that list the given feature ID. Use --transitive to follow the full chain.`,
		Args: cobra.ExactArgs(1),
		RunE: runFeatureRefs,
	}
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("fields", "", "comma-separated metadata fields to include (e.g., status,oq)")
	cmd.Flags().String("format", "", "output format: yaml, json, text (auto-selects yaml when --fields is set)")
	cmd.Flags().Bool("transitive", false, "follow reference chain recursively")
	return cmd
}

func runFeatureRefs(cmd *cobra.Command, args []string) error {
	featureID := args[0]
	projectFlag, _ := cmd.Flags().GetString("project")
	fieldsFlag, _ := cmd.Flags().GetString("fields")
	transitive, _ := cmd.Flags().GetBool("transitive")

	fields, err := feature.ParseFieldNames(fieldsFlag)
	if err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}

	format := effectiveFormat(cmd)
	if err := validateFormat(format); err != nil {
		return err
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	if !feature.Exists(featuresDir, featureID) {
		return exitcode.NotFoundErrorf("feature not found: %s", featureID)
	}

	w := cmd.OutOrStdout()

	if transitive {
		nodes := feature.TransitiveRefs(featuresDir, featureID)
		if len(fields) > 0 {
			feature.EnrichTransitiveNodes(featuresDir, nodes, fields)
		}
		switch format {
		case "yaml":
			return writeEnrichedYAML(w, nodes)
		case "json":
			return writeEnrichedJSON(w, nodes)
		default:
			if len(fields) > 0 {
				return writeEnrichedText(w, nodes, fields)
			}
			var sb strings.Builder
			printTransitiveText(&sb, nodes, 0)
			_, _ = fmt.Fprint(w, sb.String())
		}
		return nil
	}

	refs, err := feature.FindFeatureRefs(featuresDir, featureID)
	if err != nil {
		return exitcode.UnexpectedErrorf("finding references: %v", err)
	}

	if len(fields) > 0 || format == "yaml" || format == "json" {
		var enriched []*feature.EnrichedFeature
		for _, ref := range refs {
			ef, _ := feature.ResolveFields(featuresDir, ref, fields)
			enriched = append(enriched, ef)
		}
		return writeEnrichedOutput(w, enriched, fields, format)
	}

	for _, ref := range refs {
		_, _ = fmt.Fprintln(w, ref)
	}
	return nil
}

// --- feature new ---

func featureNewCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Scaffold a new feature directory with a README template",
		Long: `Creates a new feature directory with a README containing all required
sections. Changes are local by default; use --commit to create a git
commit, or --push for atomic commit-and-push.`,
		RunE: runFeatureNew,
	}
	cmd.Flags().String("title", "", "human-readable feature title (required)")
	cmd.Flags().String("slug", "", "feature slug (directory name); auto-generated from title if omitted")
	cmd.Flags().String("parent", "", "parent feature ID for creating a sub-feature")
	cmd.Flags().String("status", "draft", "initial feature status: draft, approved, implemented")
	cmd.Flags().String("description", "", "short description placed in the Summary section")
	cmd.Flags().String("depends-on", "", "comma-separated list of feature IDs this feature depends on")
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().String("format", "yaml", "output format: yaml, json, text")
	cmd.Flags().Bool("commit", false, "create a git commit with the changes")
	cmd.Flags().Bool("push", false, "commit and push atomically (implies --commit)")
	return cmd
}

func runFeatureNew(cmd *cobra.Command, _ []string) error {
	title, _ := cmd.Flags().GetString("title")
	slugFlag, _ := cmd.Flags().GetString("slug")
	parentFlag, _ := cmd.Flags().GetString("parent")
	statusFlag, _ := cmd.Flags().GetString("status")
	descFlag, _ := cmd.Flags().GetString("description")
	depsFlag, _ := cmd.Flags().GetString("depends-on")
	projectFlag, _ := cmd.Flags().GetString("project")
	formatFlag, _ := cmd.Flags().GetString("format")
	commitFlag, _ := cmd.Flags().GetBool("commit")
	pushFlag, _ := cmd.Flags().GetBool("push")

	if title == "" {
		return exitcode.InvalidArgsError("missing required flag: --title")
	}
	if formatFlag != "yaml" && formatFlag != "json" && formatFlag != "text" {
		return exitcode.InvalidArgsErrorf("invalid format: %s (supported: yaml, json, text)", formatFlag)
	}
	if pushFlag {
		commitFlag = true
	}

	// Parse --depends-on
	var deps []string
	if depsFlag != "" {
		for _, d := range strings.Split(depsFlag, ",") {
			d = strings.TrimSpace(d)
			if d != "" {
				deps = append(deps, d)
			}
		}
	}

	featuresDir, err := resolveFeaturesDir(projectFlag)
	if err != nil {
		return err
	}

	opts := feature.NewOptions{
		Title:       title,
		Slug:        slugFlag,
		Parent:      parentFlag,
		Status:      statusFlag,
		Description: descFlag,
		DependsOn:   deps,
	}

	result, err := feature.New(featuresDir, opts)
	if err != nil {
		return err
	}

	if commitFlag {
		repoRoot := filepath.Dir(filepath.Dir(featuresDir)) // spec/features/ → repo root
		if !isGitRepo(repoRoot) {
			return exitcode.UnexpectedError("not a git repository; cannot commit")
		}

		relFiles := make([]string, 0, len(result.ChangedFiles))
		for _, f := range result.ChangedFiles {
			rel, relErr := filepath.Rel(repoRoot, f)
			if relErr != nil {
				rel = f
			}
			relFiles = append(relFiles, rel)
		}

		commitMsg := fmt.Sprintf("feat(spec): add feature %s", result.FeatureID)

		if pushFlag {
			if err := gitCommitAndPush(repoRoot, relFiles, commitMsg); err != nil {
				return exitcode.ConflictErrorf("commit and push failed: %v", err)
			}
		} else {
			if err := gitCommitOnly(repoRoot, relFiles, commitMsg); err != nil {
				return exitcode.UnexpectedErrorf("commit failed: %v", err)
			}
		}
	}

	return writeFeatureInfo(cmd.OutOrStdout(), formatFlag, &result.Info)
}

// isGitRepo returns true if dir is inside a git repository.
func isGitRepo(dir string) bool {
	cmd := exec.Command("git", "-C", dir, "rev-parse", "--git-dir")
	return cmd.Run() == nil
}

// gitCommitOnly stages files and creates a commit without pushing.
func gitCommitOnly(repoDir string, files []string, message string) error {
	addArgs := append([]string{"-C", repoDir, "add"}, files...)
	addCmd := exec.Command("git", addArgs...)
	if out, err := addCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git add: %w\n%s", err, out)
	}

	commitCmd := exec.Command("git", "-C", repoDir, "commit", "-m", message)
	if out, err := commitCmd.CombinedOutput(); err != nil {
		return fmt.Errorf("git commit: %w\n%s", err, out)
	}
	return nil
}

// gitCommitAndPush stages files, commits, and pushes. On push conflict it
// pulls, re-stages, and retries once.
func gitCommitAndPush(repoDir string, files []string, message string) error {
	if err := gitCommitOnly(repoDir, files, message); err != nil {
		return err
	}

	pushCmd := exec.Command("git", "-C", repoDir, "push")
	if out, err := pushCmd.CombinedOutput(); err != nil {
		// Try pull + re-push once.
		pullCmd := exec.Command("git", "-C", repoDir, "pull", "--rebase")
		if pullOut, pullErr := pullCmd.CombinedOutput(); pullErr != nil {
			return fmt.Errorf("git push: %w\n%s\ngit pull: %v\n%s", err, out, pullErr, pullOut)
		}
		retryCmd := exec.Command("git", "-C", repoDir, "push")
		if retryOut, retryErr := retryCmd.CombinedOutput(); retryErr != nil {
			return fmt.Errorf("git push (retry): %w\n%s", retryErr, retryOut)
		}
	}
	return nil
}
