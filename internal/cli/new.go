package cli

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/synchestra-io/specscore/pkg/exitcode"
	"github.com/synchestra-io/specscore/pkg/feature"
	"github.com/synchestra-io/specscore/pkg/idea"
	"github.com/synchestra-io/specscore/pkg/lint"
)

// newCommand returns the "new" command group — scaffolders for spec artifacts.
func newCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "new",
		Short: "Scaffold new spec artifacts (ideas, features, ...)",
	}
	cmd.AddCommand(newIdeaCommand())
	return cmd
}

// newIdeaCommand scaffolds a lint-clean Idea artifact at spec/ideas/<slug>.md.
func newIdeaCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "idea <slug>",
		Short: "Scaffold a new Idea artifact",
		Long: `Creates a lint-clean Idea skeleton at spec/ideas/<slug>.md.

Each required section is emitted with an HTML-comment prompt describing
what belongs there. Supply content via flags (--title, --owner, --hmw,
--not-doing) or run with -i to be prompted interactively. The generated
file is always lint-clean — running ` + "`specscore lint`" + ` immediately
afterwards passes.`,
		Args: cobra.ExactArgs(1),
		RunE: runNewIdea,
	}
	cmd.Flags().String("title", "", "Idea title (defaults to title-cased slug)")
	cmd.Flags().String("owner", "", "author identifier (defaults to $USER)")
	cmd.Flags().String("hmw", "", "Problem Statement (How Might We…) sentence")
	cmd.Flags().String("context", "", "Context section content")
	cmd.Flags().String("recommended-direction", "", "Recommended Direction content")
	cmd.Flags().String("mvp", "", "MVP Scope content")
	cmd.Flags().StringArray("not-doing", nil, "Not Doing exclusion (repeatable). Format: `<thing> — <reason>`")
	cmd.Flags().String("project", "", "project root (autodetected from current directory if omitted)")
	cmd.Flags().BoolP("interactive", "i", false, "prompt for each field on stdin")
	cmd.Flags().Bool("force", false, "overwrite an existing idea file at that slug")
	return cmd
}

func runNewIdea(cmd *cobra.Command, args []string) error {
	slug := args[0]
	if err := idea.ValidateSlug(slug); err != nil {
		return exitcode.InvalidArgsErrorf("invalid slug %q: %v", slug, err)
	}

	title, _ := cmd.Flags().GetString("title")
	owner, _ := cmd.Flags().GetString("owner")
	hmw, _ := cmd.Flags().GetString("hmw")
	ctx, _ := cmd.Flags().GetString("context")
	direction, _ := cmd.Flags().GetString("recommended-direction")
	mvp, _ := cmd.Flags().GetString("mvp")
	notDoing, _ := cmd.Flags().GetStringArray("not-doing")
	projectFlag, _ := cmd.Flags().GetString("project")
	interactive, _ := cmd.Flags().GetBool("interactive")
	force, _ := cmd.Flags().GetBool("force")

	if owner == "" {
		if u := os.Getenv("USER"); u != "" {
			owner = u
		}
	}

	opts := idea.ScaffoldOptions{
		Slug:                 slug,
		Title:                title,
		Owner:                owner,
		HMW:                  hmw,
		Context:              ctx,
		RecommendedDirection: direction,
		MVP:                  mvp,
		NotDoing:             notDoing,
	}

	if interactive {
		if err := runInteractivePrompts(cmd.InOrStdin(), cmd.OutOrStdout(), &opts); err != nil {
			return err
		}
	}

	specRoot, err := resolveSpecRoot(projectFlag)
	if err != nil {
		return err
	}
	ideasDir := filepath.Join(specRoot, "spec", "ideas")
	if err := os.MkdirAll(ideasDir, 0o755); err != nil {
		return exitcode.UnexpectedErrorf("creating %s: %v", ideasDir, err)
	}
	target := filepath.Join(ideasDir, slug+".md")
	if _, err := os.Stat(target); err == nil && !force {
		return exitcode.ConflictErrorf("idea already exists: %s (pass --force to overwrite)", target)
	}

	body, err := idea.Scaffold(opts)
	if err != nil {
		return exitcode.UnexpectedErrorf("scaffolding idea: %v", err)
	}
	if err := os.WriteFile(target, body, 0o644); err != nil {
		return exitcode.UnexpectedErrorf("writing %s: %v", target, err)
	}

	// Run lint in --fix mode to update the active index, then re-run
	// without fix to surface any remaining errors touching this file.
	specSub := filepath.Join(specRoot, "spec")
	if _, err := lint.Lint(lint.Options{SpecRoot: specSub, Fix: true}); err != nil {
		// Remove the partial file so re-runs don't trip over conflict.
		_ = os.Remove(target)
		return exitcode.UnexpectedErrorf("running lint fix: %v", err)
	}
	violations, err := lint.Lint(lint.Options{SpecRoot: specSub})
	if err != nil {
		return exitcode.UnexpectedErrorf("running lint: %v", err)
	}
	relTarget, _ := filepath.Rel(specSub, target)
	var own []lint.Violation
	for _, v := range violations {
		if v.Severity == "error" && (v.File == relTarget || strings.HasPrefix(v.File, "ideas/")) {
			own = append(own, v)
		}
	}
	if len(own) > 0 {
		var sb strings.Builder
		sb.WriteString("generated idea failed lint:\n")
		for _, v := range own {
			fmt.Fprintf(&sb, "  %s:%d [%s] %s\n", v.File, v.Line, v.Rule, v.Message)
		}
		return exitcode.UnexpectedError(sb.String())
	}

	_, _ = fmt.Fprintf(cmd.OutOrStdout(), "%s\n", target)
	return nil
}

// resolveSpecRoot resolves the project root (repo root, not spec/ itself)
// from --project or cwd, using the same heuristic as feature commands.
func resolveSpecRoot(projectFlag string) (string, error) {
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
	return feature.FindSpecRepoRoot(startDir)
}

// runInteractivePrompts fills unset fields in opts by reading line-delimited
// input from r. An empty line keeps the current default.
func runInteractivePrompts(r io.Reader, w io.Writer, opts *idea.ScaffoldOptions) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)

	type prompt struct {
		label string
		dst   *string
	}
	prompts := []prompt{
		{"Title", &opts.Title},
		{"Owner", &opts.Owner},
		{"Problem Statement (How Might We…)", &opts.HMW},
		{"Context", &opts.Context},
		{"Recommended Direction", &opts.RecommendedDirection},
		{"MVP Scope", &opts.MVP},
	}
	for _, p := range prompts {
		cur := *p.dst
		if cur != "" {
			_, _ = fmt.Fprintf(w, "%s [%s]: ", p.label, cur)
		} else {
			_, _ = fmt.Fprintf(w, "%s: ", p.label)
		}
		if !scanner.Scan() {
			break
		}
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			*p.dst = line
		}
	}

	// Not Doing — accept multiple lines until blank.
	if len(opts.NotDoing) == 0 {
		_, _ = fmt.Fprintln(w, "Not Doing exclusions (one per line, blank to finish):")
		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if line == "" {
				break
			}
			opts.NotDoing = append(opts.NotDoing, line)
		}
	}

	return scanner.Err()
}
