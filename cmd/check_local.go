package cmd

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/cemc-oper/nwpc-data-client/internal/checker"
	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/spf13/cobra"
)

func init() {
	checkLocalCmd.Flags().SortFlags = false

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataConfigDir, "data-config-dir", "",
		"Directory holding data config files. Optional; uses embedded configs when empty.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataConfigFile, "data-config-file", "",
		"Path to a single data config file. If set, --data-config-dir and --data-type are ignored.")

	checkLocalCmd.Flags().StringVar(&checkConfigFile, "checker-config", "",
		"Path to the checker runtime config YAML. CLI flags override values in this file.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DataType, "data-type", "",
		"Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.LocationLevels, "location-level", "",
		"Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.")
	checkLocalCmd.Flags().StringVar(&checkFlagConfig.StorageBase, "storage-base", "",
		"Absolute local mount root for {{.StorageBase}} in data config and command templates.")

	checkLocalCmd.Flags().IntVar(&checkFlagConfig.MaxCheckCount, "max-check-count", 2880,
		"Maximum number of existence/size checks per forecast time.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CheckInterval, "check-interval", "5s",
		"Interval between checks, a duration such as 5s, 30s, or 1m.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.ExecuteCommand, "execute-command", "",
		"Go template for a command to run once a file becomes stable. Mutually exclusive with execute_commands.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CommandTimeout, "command-timeout", "30m",
		"Timeout for each executed command, such as 30m or 1h. 0 disables the timeout.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.ReportFile, "report-file", "",
		"Path for the optional JSON run report. Empty does not generate a report.")

	checkLocalCmd.Flags().BoolVar(&checkDryRun, "dry-run", false,
		"Print the resolved run plan without polling files or executing commands.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DelayTime, "delay-time", "10s",
		"Stagger delay between consecutive forecast times (duration such as 10s).")

	checkLocalCmd.Flags().Float64Var(&checkFlagConfig.CheckIntervalJitter, "check-interval-jitter", 0.1,
		"Jitter ratio applied to the polling interval: each round sleeps interval×(1±j). 0 disables jitter. Range [0, 0.5].")

	checkLocalCmd.Flags().IntVar(&checkFlagConfig.MaxConcurrentStableChecks, "max-concurrent-stable-checks", 4,
		"Maximum number of files in size-stability observation at once (batched check mode only).")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CheckMode, "check-mode", "per_time",
		"Polling mode: per_time (one goroutine per forecast time) or batched (single readdir-based loop).")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.ArrivalOrder, "arrival-order", "auto",
		"Arrival order override for batched mode gating: auto (follow the data config), sequential, or unordered.")

	checkLocalCmd.Flags().BoolVar(&checkFlagConfig.Debug, "debug", false, "Enable debug logging.")
}

const checkLocalCommandName = "local"

const checkLocalCommandDocString = `Wait for a list of data files with different forecast times and
run command(s) when each file is found.

For each forecast time, this command repeatedly looks up the data file (using
the same config resolution as "nwpc_data_client local"). Once a file exists and
its size stops changing between checks, it is considered stable and the
configured command(s) are run. Each forecast time is checked concurrently after
a staggered --delay-time.

Config resolution mirrors the local command:

  * --data-type            Use an embedded config named "local/<data-type>".
  * --data-type combined with --data-config-dir
		Load a user defined config from "<data-config-dir>/<data-type>.yaml".
  * --data-config-file     Load a user defined config directly from "<data-config-file>".

The list of forecast times comes from --checker-config (the forecast_times
field) or, when absent, from whitespace-separated tokens on stdin.

CLI flags override values from --checker-config. The command(s) to run support
Go templates with the usual config variables (Year, Month, Day, Hour,
ForecastHour, ...) plus FilePath. Use --execute-command for a single command;
multiple commands (execute_commands) can be supplied via --checker-config.

--dry-run resolves the same run plan and prints the final mode, effective
arrival order, candidate paths, and command templates. Candidate paths are
not checked for existence or selection. Dry-run does not poll, sleep, execute
commands, or write a run report.

Args:
    start_time   Start time, YYYYMMDDHH, such as 2026062400.`

const checkLocalCommandExample = `  # Check via a checker config file (recommended).
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml

  # Provide forecast times on stdin and run a command for each ready file.
  echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
      --data-type=cma_gfs_gmf/grib2/orig \
      --execute-command="echo ready: {{.FilePath}}"

  # Print the resolved plan without accessing candidate files or running commands.
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml --dry-run

  # Poll faster and limit the number of attempts per forecast time.
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml \
      --check-interval=10s --max-check-count=360`

var checkLocalCmd = &cobra.Command{
	Use:     checkLocalCommandName + " <start_time>",
	Short:   "Wait for local files, then run a command for each one.",
	Long:    checkLocalCommandDocString,
	Example: checkLocalCommandExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires one arguments")
		}
		if _, err := dataconfig.ParseStartTime(args[0]); err != nil {
			return fmt.Errorf("check startTime failed: %s", err)
		}
		return nil
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Argument errors still show usage; runtime failures return a single
		// diagnostic to the root command without a usage page.
		cmd.SilenceUsage = true

		// Args already validated the format, so this parse cannot fail.
		startTime, _ := dataconfig.ParseStartTime(args[0])

		// Load checker runtime config from YAML if provided, otherwise start empty.
		fileConfig := checker.CheckerConfig{}
		if checkConfigFile != "" {
			var err error
			fileConfig, err = checker.LoadCheckerConfig(checkConfigFile)
			if err != nil {
				return fmt.Errorf("load checker config failed: %w", err)
			}
		}

		// Merge YAML config with CLI flags into a single config object.
		config, err := checker.MergeCheckerConfig(fileConfig, checkFlagConfig, cmd)
		if err != nil {
			return fmt.Errorf("invalid checker config: %w", err)
		}

		if checkDryRun {
			if err := checker.DryRun(config, startTime, os.Stdin, cmd.OutOrStdout()); err != nil {
				return fmt.Errorf("dry run failed: %w", err)
			}
			return nil
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		defer stop()
		result, runErr := checker.Run(ctx, config, startTime, os.Stdin)
		var reportErr error
		if result != nil {
			reportErr = checker.WriteRunReport(config.ReportFile, result)
		}
		if err := errors.Join(runErr, reportErr); err != nil {
			if ctx.Err() != nil {
				return &exitCodeError{code: 3, err: fmt.Errorf("check local cancelled: %w", err)}
			}
			return err
		}
		return nil
	},
}
