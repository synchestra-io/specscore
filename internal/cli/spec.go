package cli

// Features implemented: cli/spec

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/feature"
	"github.com/synchestra-io/specscore/pkg/lint"
	"gopkg.in/yaml.v3"
)

// specCommand returns the "spec" command group.
func specCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "spec",
		Short: "Validate and search specification repositories",
	}
	cmd.AddCommand(
		specLintCommand(),
	)
	return cmd
}

func specLintCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "lint",
		Short: "Validate spec tree for structural convention violations",
		Long: `Scans the specification tree and reports violations of structural
conventions (README.md files, Outstanding Questions sections, heading
levels, feature references, internal links, index entries).

Violations are categorized by severity: error (must fix), warning (should
fix), info (advisory). Exit code 0 = valid, 1 = violations found, 2 =
invalid arguments, 10+ = unexpected error.`,
		Args: cobra.NoArgs,
		RunE: runSpecLint,
	}
	cmd.Flags().StringP("project", "p", "", "path to spec repository root (default: auto-discover from CWD)")
	cmd.Flags().String("rules", "", "enable only specified rules (comma-separated)")
	cmd.Flags().String("ignore", "", "disable specified rules (comma-separated)")
	cmd.Flags().String("severity", "error", "minimum severity: error, warning, info")
	cmd.Flags().String("format", "text", "output format: text, json, yaml")
	return cmd
}

func runSpecLint(cmd *cobra.Command, args []string) error {
	projectFlag, _ := cmd.Flags().GetString("project")

	var startDir string
	if projectFlag != "" {
		abs, err := filepath.Abs(projectFlag)
		if err != nil {
			return exitcode.InvalidArgsErrorf("resolving --project path: %v", err)
		}
		startDir = abs
	} else {
		cwd, err := os.Getwd()
		if err != nil {
			return exitcode.UnexpectedErrorf("cannot determine working directory: %v", err)
		}
		startDir = cwd
	}

	root, err := feature.FindSpecRepoRoot(startDir)
	if err != nil {
		return err
	}
	specRoot := filepath.Join(root, "spec")

	rulesStr, _ := cmd.Flags().GetString("rules")
	ignoreStr, _ := cmd.Flags().GetString("ignore")
	severity, _ := cmd.Flags().GetString("severity")
	format, _ := cmd.Flags().GetString("format")

	// Validate mutual exclusion.
	if rulesStr != "" && ignoreStr != "" {
		return exitcode.InvalidArgsError("--rules and --ignore are mutually exclusive")
	}

	// Parse rules.
	var rules []string
	if rulesStr != "" {
		for _, r := range strings.Split(rulesStr, ",") {
			r = strings.TrimSpace(r)
			if r != "" {
				rules = append(rules, r)
			}
		}
	}

	// Parse ignore.
	var ignore []string
	if ignoreStr != "" {
		for _, r := range strings.Split(ignoreStr, ",") {
			r = strings.TrimSpace(r)
			if r != "" {
				ignore = append(ignore, r)
			}
		}
	}

	// Validate severity.
	if severity != "error" && severity != "warning" && severity != "info" {
		return exitcode.InvalidArgsErrorf("invalid severity level: %s", severity)
	}

	// Validate format.
	if format != "text" && format != "json" && format != "yaml" {
		return exitcode.InvalidArgsErrorf("invalid format: %s", format)
	}

	// Validate rule names.
	if err := lint.ValidateRuleNames(rules); err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}
	if err := lint.ValidateRuleNames(ignore); err != nil {
		return exitcode.InvalidArgsError(err.Error())
	}

	opts := lint.Options{
		SpecRoot: specRoot,
		Rules:    rules,
		Ignore:   ignore,
		Severity: severity,
	}

	violations, err := lint.Lint(opts)
	if err != nil {
		return exitcode.UnexpectedErrorf("linting error: %v", err)
	}

	w := cmd.OutOrStdout()
	if err := outputLintViolations(w, violations, format); err != nil {
		return exitcode.UnexpectedErrorf("output error: %v", err)
	}

	if len(violations) > 0 {
		return exitcode.ConflictErrorf("%d violation(s) found", len(violations))
	}
	return nil
}

func outputLintViolations(w io.Writer, violations []lint.Violation, format string) error {
	switch format {
	case "json":
		return outputLintJSON(w, violations)
	case "yaml":
		return outputLintYAML(w, violations)
	default:
		return outputLintText(w, violations)
	}
}

func outputLintJSON(w io.Writer, violations []lint.Violation) error {
	if violations == nil {
		violations = []lint.Violation{}
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(violations)
}

func outputLintYAML(w io.Writer, violations []lint.Violation) error {
	if len(violations) == 0 {
		_, _ = fmt.Fprintln(w, "[]")
		return nil
	}
	enc := yaml.NewEncoder(w)
	enc.SetIndent(2)
	if err := enc.Encode(violations); err != nil {
		return err
	}
	return enc.Close()
}

func outputLintText(w io.Writer, violations []lint.Violation) error {
	// Sort by file then line.
	sort.Slice(violations, func(i, j int) bool {
		if violations[i].File != violations[j].File {
			return violations[i].File < violations[j].File
		}
		return violations[i].Line < violations[j].Line
	})

	for _, v := range violations {
		_, _ = fmt.Fprintf(w, "%s:%d [%s] %s: %s\n", v.File, v.Line, v.Severity, v.Rule, v.Message)
	}

	if len(violations) > 0 {
		errorCount := 0
		warningCount := 0
		infoCount := 0
		for _, v := range violations {
			switch v.Severity {
			case "error":
				errorCount++
			case "warning":
				warningCount++
			case "info":
				infoCount++
			}
		}

		_, _ = fmt.Fprintf(w, "\n%d violations found", len(violations))
		var parts []string
		if errorCount > 0 {
			parts = append(parts, fmt.Sprintf("%d error%s", errorCount, lintPlural(errorCount)))
		}
		if warningCount > 0 {
			parts = append(parts, fmt.Sprintf("%d warning%s", warningCount, lintPlural(warningCount)))
		}
		if infoCount > 0 {
			parts = append(parts, fmt.Sprintf("%d info", infoCount))
		}
		if len(parts) > 0 {
			_, _ = fmt.Fprintf(w, " (%s)", strings.Join(parts, ", "))
		}
		_, _ = fmt.Fprintln(w)
	} else {
		_, _ = fmt.Fprintln(w, "0 violations found")
	}
	return nil
}

func lintPlural(n int) string {
	if n == 1 {
		return ""
	}
	return "s"
}
