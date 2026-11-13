package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const rootCommandDocString = `nwpc_data_client is a data file finder for CEMC operational systems.

It searches for data files across multiple directories and polls a series of
files until they are generated, running commands once they are found.

Common subcommands:
  local         Find a local file.
  fetch         Copy or link one local file into a target directory.
  check local   Wait for several files and run a command(s) when each file is found.
  version       Print version information.

Use "nwpc_data_client [command] --help" for detailed help on a subcommand.`

var rootCmd = &cobra.Command{
	Use:           "nwpc_data_client",
	Short:         "Find data files for operational systems in CEMC.",
	Long:          rootCommandDocString,
	SilenceErrors: true,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

type exitCodeError struct {
	code int
	err  error
}

func (e *exitCodeError) Error() string {
	return e.err.Error()
}

func (e *exitCodeError) Unwrap() error {
	return e.err
}

func (e *exitCodeError) ExitCode() int {
	return e.code
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(rootCmd.ErrOrStderr(), err)
		exitCode := 1
		var coded interface{ ExitCode() int }
		if errors.As(err, &coded) {
			exitCode = coded.ExitCode()
		}
		os.Exit(exitCode)
	}
}

// RootCommand returns the root cobra command. It is used by auxiliary
// tooling such as the CLI documentation generator in tools/gen-cli-docs.
func RootCommand() *cobra.Command {
	return rootCmd
}
