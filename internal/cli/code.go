package cli

// Features implemented: cli/code

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/sourceref"
)

// codeCommand returns the "code" command group.
func codeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "code",
		Short: "Query source code relationships to SpecScore resources",
		Long: `Commands for querying source code relationships to SpecScore resources.
Scans source files for specscore: annotations and URLs embedded in comments,
showing the resources (features, plans, docs) that code depends on.`,
	}
	cmd.AddCommand(
		codeDepsCommand(),
	)
	return cmd
}

func codeDepsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deps [flags]",
		Short: "Show SpecScore resources that source files depend on",
		Long: `Shows the SpecScore resources (features, plans, docs) that source files
depend on. Scans source files for specscore: annotations and URLs in comments
and lists the referenced resources.

This is a read-only command that scans the working tree and does not mutate anything.`,
		RunE: runCodeDeps,
	}
	cmd.Flags().String("path", "**/*", "glob pattern to select files (e.g., pkg/**/*.go). Defaults to **/* (all files)")
	cmd.Flags().String("type", "", "filter results to a specific resource type: feature, plan, or doc")
	return cmd
}

func runCodeDeps(cmd *cobra.Command, _ []string) error {
	pathPattern, _ := cmd.Flags().GetString("path")
	typeFilter, _ := cmd.Flags().GetString("type")

	if typeFilter != "" && typeFilter != "feature" && typeFilter != "plan" && typeFilter != "doc" {
		return exitcode.InvalidArgsErrorf("invalid --type value: %s (must be feature, plan, or doc)", typeFilter)
	}

	files, err := sourceref.ExpandGlobPattern(pathPattern)
	if err != nil {
		return exitcode.InvalidArgsErrorf("invalid glob pattern %q: %v", pathPattern, err)
	}

	if len(files) == 0 {
		return nil
	}

	result, err := sourceref.ScanFiles(files)
	if err != nil {
		return exitcode.UnexpectedErrorf("scanning files: %v", err)
	}

	singleFile := len(result.FileRefs) == 1

	w := cmd.OutOrStdout()
	output := sourceref.FormatOutput(result, singleFile, typeFilter)
	if output != "" {
		_, _ = fmt.Fprint(w, output)
	}

	return nil
}
