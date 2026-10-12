package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/cemc-oper/nwpc-data-client/internal/checker"
	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
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

	checkLocalCmd.Flags().IntVar(&checkFlagConfig.MaxCheckCount, "max-check-count", 2880,
		"Maximum number of existence/size checks per forecast time.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CheckInterval, "check-interval", "5s",
		"Interval between checks, a duration such as 5s, 30s, or 1m.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.ExecuteCommand, "execute-command", "",
		"Go template for a command to run once a file becomes stable. Mutually exclusive with execute_commands.")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.DelayTime, "delay-time", "10s",
		"Stagger delay between consecutive forecast times (duration such as 10s).")

	checkLocalCmd.Flags().Float64Var(&checkFlagConfig.CheckIntervalJitter, "check-interval-jitter", 0.1,
		"Jitter ratio applied to the polling interval: each round sleeps interval×(1±j). 0 disables jitter. Range [0, 0.5].")

	checkLocalCmd.Flags().IntVar(&checkFlagConfig.MaxConcurrentStableChecks, "max-concurrent-stable-checks", 4,
		"Maximum number of files in size-stability observation at once (batched check mode only).")

	checkLocalCmd.Flags().StringVar(&checkFlagConfig.CheckMode, "check-mode", "per_time",
		"Polling mode: per_time (one goroutine per forecast time) or batched (single readdir-based loop).")

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

Args:
    start_time   Start time, YYYYMMDDHH, such as 2026062400.`

const checkLocalCommandExample = `  # Check via a checker config file (recommended).
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml

  # Provide forecast times on stdin and run a command for each ready file.
  echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
      --data-type=cma_gfs_gmf/grib2/orig \
      --execute-command="echo ready: {{.FilePath}}"

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
	Run: func(cmd *cobra.Command, args []string) {
		// Args already validated the format, so this parse cannot fail.
		startTime, _ := dataconfig.ParseStartTime(args[0])

		// Load checker runtime config from YAML if provided, otherwise start empty.
		fileConfig := checker.CheckerConfig{}
		if checkConfigFile != "" {
			var err error
			fileConfig, err = checker.LoadCheckerConfig(checkConfigFile)
			if err != nil {
				log.Fatalf("load checker config failed: %v", err)
			}
		}

		// Merge YAML config with CLI flags into a single config object.
		config, err := checker.MergeCheckerConfig(fileConfig, checkFlagConfig, cmd)
		if err != nil {
			log.Fatalf("invalid checker config: %v", err)
		}

		checker.Run(config, startTime, os.Stdin)
	},
}
