package cmd

import (
	"errors"
	"fmt"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/cmafs"
	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	"github.com/cemc-oper/nwpc-data-client/internal/fetcher"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().SortFlags = false

	fetchCmd.Flags().StringVar(&fetchDataConfigDir, "data-config-dir", "",
		"Directory holding data config files. Optional; uses embedded configs when empty.")
	fetchCmd.Flags().StringVar(&fetchDataType, "data-type", "",
		"Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).")
	fetchCmd.Flags().StringVar(&fetchDataConfigFile, "data-config-file", "",
		"Path to a single data config file. If set, --data-config-dir and --data-type are ignored.")
	fetchCmd.Flags().StringVar(&fetchBackendConfig, "backend-config", "",
		"Path to machine backend YAML. Required only when a cmafs source is needed.")
	fetchCmd.Flags().StringVar(&fetchLocationLevels, "location-level", "",
		"Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.")
	fetchCmd.Flags().StringVar(&fetchStorageBase, "storage-base", "",
		"Absolute local mount root for {{.StorageBase}} in data config templates.")
	fetchCmd.Flags().StringVar(&fetchStartTimeString, "start-time", "",
		"Start time, YYYYMMDDHH, such as 2026062400. (required)")
	fetchCmd.Flags().StringVar(&fetchForecastTimeString, "forecast-time", "",
		"Forecast time, such as 0h, 24h, 120h. (required)")
	fetchCmd.Flags().StringVar(&fetchMember, "member", "",
		"Ensemble member, MMM, such as 000, 014.")
	fetchCmd.Flags().StringVar(&fetchTargetDir, "target-dir", "",
		"Directory in which the fetched file is published. (required)")
	fetchCmd.Flags().StringVar(&fetchMethod, "method", "copy",
		"Materialization method: copy, symlink, or auto. Auto copies storage-level sources.")
	fetchCmd.Flags().StringVar(&fetchStableInterval, "stable-interval", "1s",
		"Interval between two source size/mtime stability checks, such as 1s or 500ms.")
	fetchCmd.Flags().BoolVar(&fetchDebug, "debug", false, "Enable debug logging.")

	_ = fetchCmd.MarkFlagRequired("start-time")
	_ = fetchCmd.MarkFlagRequired("forecast-time")
	_ = fetchCmd.MarkFlagRequired("target-dir")
}

const fetchCommandDocString = `Find one data file and materialize it into a target directory.

The data config is resolved with the same start time, forecast time, member,
location levels, and candidate ordering as "nwpc_data_client local". The first
existing regular file is selected; missing candidates are skipped, while other
source errors stop the operation. Local candidates are read directly. Cmafs
candidates use the machine backend supplied by --backend-config and direct
SSH/rsync transport.

The source must be non-empty and have unchanged size and modification time
across two observations separated by --stable-interval. The destination name
is the source basename. Copy mode writes a unique temporary file in the target
directory, verifies it, and publishes it without replacing an existing file.
Symlink mode creates an absolute symbolic link to the selected local source.
Auto uses symlink for runtime/archive local sources and copy for storage local
sources and cmafs downloads. Cmafs downloads are verified against remote size
and mtime before and after transfer.

Existing destinations are never overwritten. An identical regular file or an
allowed matching symlink is reused; a conflicting destination is an error.
A per-target interprocess lock prevents concurrent duplicate materialization.

On success, only the final absolute path is written to stdout. Logs and errors
go to stderr. This command does not wait for files that have not appeared yet.`

const fetchCommandExample = `  # Copy a CMA-GFS file into a work directory.
  nwpc_data_client fetch --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h --target-dir=/work/input

  # Link a local runtime file instead of copying it.
  nwpc_data_client fetch --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h \
      --target-dir=/work/input --method=symlink

  # Use a single config file and automatic method selection.
  nwpc_data_client fetch --data-config-file=./my_config.yaml \
      --start-time=2026062400 --forecast-time=3h \
      --target-dir=./input --method=auto`

var fetchCmd = &cobra.Command{
	Use:     "fetch",
	Short:   "Find one data file and copy or link it into a target directory.",
	Long:    fetchCommandDocString,
	Example: fetchCommandExample,
	Args:    cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		cmd.SilenceUsage = true
		if fetchDebug {
			log.SetLevel(log.DebugLevel)
		}
		if fetchDataConfigFile == "" && fetchDataType == "" {
			return errors.New("requires --data-type unless --data-config-file is set")
		}

		startTime, err := dataconfig.ParseStartTime(fetchStartTimeString)
		if err != nil {
			return fmt.Errorf("check startTime failed: %w", err)
		}
		forecastTime, err := dataconfig.ParseForecastTime(fetchForecastTimeString)
		if err != nil {
			return fmt.Errorf("check forecastTime failed: %w", err)
		}
		stableInterval, err := time.ParseDuration(fetchStableInterval)
		if err != nil {
			return fmt.Errorf("invalid --stable-interval %q: %w", fetchStableInterval, err)
		}

		var cmafsBackend fetcher.RemoteBackend
		if fetchBackendConfig != "" {
			backend, err := cmafs.Load(fetchBackendConfig)
			if err != nil {
				return fmt.Errorf("load backend config failed: %w", err)
			}
			cmafsBackend = backend
		}

		var configContent string
		if fetchDataConfigFile != "" {
			configContent, err = dataconfig.LoadConfigContentFromFile(fetchDataConfigFile)
		} else {
			dataType := fetchDataType
			if fetchDataConfigDir == "" {
				dataType = "local/" + dataType
			}
			configContent, err = dataconfig.LoadConfigContent(fetchDataConfigDir, dataType)
		}
		if err != nil {
			return fmt.Errorf("load config failed: %w", err)
		}
		dataConfig, err := dataconfig.ParseConfigContentWithStorageBase(configContent, startTime, forecastTime, fetchMember, fetchStorageBase)
		if err != nil {
			return fmt.Errorf("render data config: %w", err)
		}

		ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
		defer stop()
		result, err := fetcher.Fetch(ctx, fetcher.Config{
			TargetDir:      fetchTargetDir,
			Method:         fetcher.Method(fetchMethod),
			StableInterval: stableInterval,
			CmafsBackend:   cmafsBackend,
		}, dataConfig, strings.Split(fetchLocationLevels, ","))
		if err != nil {
			if ctx.Err() != nil {
				return &exitCodeError{code: 3, err: fmt.Errorf("fetch cancelled: %w", err)}
			}
			return err
		}
		_, err = fmt.Fprintln(cmd.OutOrStdout(), result.FinalPath)
		return err
	},
}
