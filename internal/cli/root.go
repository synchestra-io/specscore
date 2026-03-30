package cli

import (
	"context"
	"errors"
	"os"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Run executes the specscore CLI with the given arguments.
func Run(args []string) error {
	rootCmd := &cobra.Command{
		Use:           "specscore",
		Short:         "SpecScore CLI — validate and query specification repositories",
		SilenceUsage:  true,
		SilenceErrors: true,
		RunE: func(cmd *cobra.Command, _ []string) error {
			return cmd.Help()
		},
	}
	rootCmd.SetErr(os.Stderr)

	rootCmd.AddCommand(
		versionCommand(),
	)

	rootCmd.SetArgs(args[1:])
	return rootCmd.ExecuteContext(context.Background())
}

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the specscore version",
		Run: func(cmd *cobra.Command, _ []string) {
			_, _ = cmd.OutOrStdout().Write([]byte("specscore " + version + " (" + commit + ") " + date + "\n"))
		},
	}
}

// Fatal prints the error and exits with the appropriate code.
func Fatal(err error) {
	if err == nil {
		return
	}
	_, _ = os.Stderr.WriteString(err.Error() + "\n")

	type exitCoder interface {
		ExitCode() int
	}
	var ec exitCoder
	if errors.As(err, &ec) {
		os.Exit(ec.ExitCode())
	}
	os.Exit(1)
}
