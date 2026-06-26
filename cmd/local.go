package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/cemc-oper/nwpc-data-client/common"
	"github.com/cemc-oper/nwpc-data-client/common/config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().SortFlags = false

	localCmd.Flags().StringVar(&localDataConfigDir, "data-config-dir", "",
		"Directory holding data config files. Optional; uses embedded configs when empty.")
	localCmd.Flags().StringVar(&localDataType, "data-type", "",
		"Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).")
	localCmd.Flags().StringVar(&localDataConfigFile, "data-config-file", "",
		"Path to a single data config file. If set, --data-config-dir and --data-type are ignored.")

	localCmd.Flags().StringVar(&localLocationLevels, "location-level", "",
		"Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.")
	localCmd.Flags().StringVar(&localStartTimeString, "start-time", "",
		"Start time, YYYYMMDDHH, such as 2026062400. (required)")
	localCmd.Flags().StringVar(&localForecastTimeString, "forecast-time", "",
		"Forecast time, such as 0h, 24h, 120h. (required)")
	localCmd.Flags().StringVar(&localMember, "member", "",
		"Ensemble member, MMM, such as 000, 014.")

	localCmd.Flags().BoolVar(&localShowTypes, "show-types", false,
		"List the supported data types and exit. Honors --data-config-dir when set.")
	localCmd.Flags().BoolVar(&localDebugMode, "debug", false, "Enable debug logging.")
}

const localCommandName = "local"

const localCommandDocString = `Find a local data file through multiple directories and print the file path when found.

The path is computed by applying --start-time and --forecast-time to a YAML
config template. The config is chosen in one of three ways:

  * --data-type            Use an embedded config named "local/<data-type>".
  * --data-type combined with --data-config-dir  Load a user defined config from "<data-config-dir>/<data-type>.yaml".
  * --data-config-file     Load a user defined config directly from "<data-config-file>".

If a matching file exists it is printed; otherwise the config's "default" value
(usually NOTFOUND) is printed. Use --location-level to restrict which path
entries are searched (for example "runtime" or "runtime,archive").

Run with --show-types to list the available data types (embedded, or from
--data-config-dir) and exit.`

const localCommandExample = `  # Use an embedded config for CMA-GFS grib2 data.
  nwpc_data_client local --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h

  # Use a config directory and only search runtime locations.
  nwpc_data_client local --data-config-dir=./config --data-type=cma_meso_1km/bin/modelvar \
      --start-time=2026062400 --forecast-time=24h --location-level=runtime

  # Use a single config file directly.
  nwpc_data_client local --data-config-file=./my_config.yaml \
      --start-time=2026062400 --forecast-time=3h

  # List available embedded data types.
  nwpc_data_client local --show-types`

var localCmd = &cobra.Command{
	Use:     localCommandName,
	Short:   "Find the local file path for model output data.",
	Long:    localCommandDocString,
	Example: localCommandExample,
	Args: func(cmd *cobra.Command, args []string) error {
		if localShowTypes {
			return nil
		}

		if localDataConfigFile == "" {
			cmd.MarkFlagRequired("data-type")
		}
		cmd.MarkFlagRequired("start-time")
		cmd.MarkFlagRequired("forecast-time")

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if localDebugMode {
			log.Info("Running in debug mode")
			log.SetLevel(log.DebugLevel)
		}
		if localShowTypes {
			showDataTypes(cmd, args)
		} else {
			findLocalFile(cmd, args)
		}
	},
}

func findLocalFile(cmd *cobra.Command, args []string) {
	startTime, err := common.ParseStartTime(localStartTimeString)
	if err != nil {
		log.Errorf("check startTime failed: %s", err)
		return
	}
	localStartTime = startTime

	forecastTime, err := common.ParseForecastTime(localForecastTimeString)
	if err != nil {
		log.Errorf("check forecastTime failed: %s", err)
		return
	}
	localForecastTime = forecastTime

	levels := strings.Split(localLocationLevels, ",")

	var configContent string
	if localDataConfigFile != "" {
		configContent, err = common.LoadConfigContentFromFile(localDataConfigFile)
	} else {
		dataType := localDataType
		if len(localDataConfigDir) == 0 {
			dataType = localCommandName + "/" + dataType
		}
		configContent, err = common.LoadConfigContent(localDataConfigDir, dataType)
	}
	if err != nil {
		log.Fatalf("load config failed: %v", err)
		return
	}

	dataConfig, err := common.ParseConfigContent(configContent, localStartTime, localForecastTime, localMember)
	if err != nil {
		log.Fatalf("load config from content has error: %s", err)
		return
	}

	pathItem := common.FindLocalFile(dataConfig, levels, localStartTime, localForecastTime)
	fmt.Printf("%s\n", pathItem.Path)
}

func showDataTypes(cmd *cobra.Command, args []string) {
	if len(localDataConfigDir) == 0 {
		showEmbeddedDataTypes()
	} else {
		showLocalDataTypes(localDataConfigDir)
	}
}

func showLocalDataTypes(configDir string) {
	var configFilePaths []string
	walkConfigDirectory := func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if (!info.IsDir()) && (filepath.Ext(path) == ".yaml") {
			configFilePaths = append(configFilePaths, path)
		}
		return nil
	}

	err := filepath.Walk(configDir, walkConfigDirectory)
	if err != nil {
		log.Errorf("Walk config directory has error: %s", err)
		return
	}
	for _, configPath := range configFilePaths {
		relConfigPath, err2 := filepath.Rel(configDir, configPath)
		if err2 != nil {
			log.Errorf("Get rel path failed: %s", err2)
			continue
		}
		fmt.Printf("%s\n", relConfigPath[:len(relConfigPath)-5])
	}
}

func showEmbeddedDataTypes() {
	for _, item := range config.EmbeddedConfigs {
		name := item[0]
		if strings.HasPrefix(name, "local/") {
			fmt.Printf("%s\n", name[6:])
		}
	}
}
