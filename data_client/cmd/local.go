package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/nwpc-oper/nwpc-data-client/common/config"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().SortFlags = false

	localCmd.Flags().StringVar(&configDir, "config-dir", "",
		"Config dir.")

	localCmd.Flags().StringVar(&dataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	localCmd.Flags().StringVar(&locationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	localCmd.Flags().BoolVar(&showTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const localCommandName = "local"

const localCommandDocString = `nwpc_data_client local
Find local data path using config files in config dir.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFFh, such as 0h, 120h`

var localCmd = &cobra.Command{
	Use:   localCommandName,
	Short: "Find local data path.",
	Long:  localCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if showTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		startTime, err = common.ParseStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check startTime failed: %s", err)
		}

		forecastTime, err = common.ParseForecastTime(args[1])
		if err != nil {
			return fmt.Errorf("check forecastTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if showTypes {
			showDataTypes(cmd, args)
		} else {
			findLocalFile(cmd, args)
		}
	},
}

func findLocalFile(cmd *cobra.Command, args []string) {
	if len(configDir) == 0 {
		dataType = localCommandName + "/" + dataType
	}
	config, err2 := common.LoadConfig(configDir, dataType)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %v\n", err2)
		return
	}

	levels := strings.Split(locationLevels, ",")

	pathItem := common.FindLocalFile(config, levels, startTime, forecastTime)
	fmt.Printf("%s\n", pathItem.Path)
}

func showDataTypes(cmd *cobra.Command, args []string) {
	if len(configDir) == 0 {
		showEmbeddedDataTypes()
	} else {
		showLocalDataTypes(configDir)
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
		fmt.Fprintf(os.Stderr, "Walk config directory has error: %s\n", err)
		return
	}
	for _, configPath := range configFilePaths {
		relConfigPath, err2 := filepath.Rel(configDir, configPath)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Get rel path failed: %s\n", err2)
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
