package cmd

import (
	"fmt"

	"github.com/cemc-oper/nwpc-data-client/internal/checker"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

const checkLintCommandName = "lint"

const checkLintCommandDocString = `Statically validate checker runtime YAML files.

Each argument may be a checker file or a directory. Files supplied directly
are validated regardless of extension; directories are searched recursively
for *.yaml files. All traversal and validation errors are reported together.

Lint applies the normal default values and validates strict YAML fields,
durations, forecast times, data source combinations, and command templates.
An empty forecast_times list is valid because check local may read forecast
times from stdin at runtime.

This command does not require a start time, load data configs, read stdin,
access data directories, execute commands, or create run reports.`

const checkLintCommandExample = `  # Lint one checker config file.
  nwpc_data_client check lint ./checker.yaml

  # Recursively lint multiple checker directories.
  nwpc_data_client check lint ./checker/system_a ./checker/system_b

  # Combine files and directories in one invocation.
  nwpc_data_client check lint ./extra.yaml ./checker`

var checkLintCmd = &cobra.Command{
	Use:     checkLintCommandName + " <checker-file-or-dir> [more paths...]",
	Short:   "Statically validate checker config files.",
	Long:    checkLintCommandDocString,
	Example: checkLintCommandExample,
	Args:    cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		fileCount, err := checker.LintCheckerPaths(args)
		if err != nil {
			log.Fatalf("check lint failed: %v", err)
		}
		fmt.Printf("check lint passed: %d file(s)\n", fileCount)
	},
}
