package cmd

import (
	"fmt"
	"github.com/cemc-oper/nwpc-data-client/common"
	"github.com/cemc-oper/nwpc-data-client/common/config"
	log "github.com/sirupsen/logrus"
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

	localCmd.Flags().StringVar(&startTimeSting, "start-time", "",
		"start time, YYYYMMDDHH, such as 2020021400")
	localCmd.Flags().StringVar(&forecastTimeString, "forecast-time", "",
		"forecast time, FFFh, such as 0h, 120h")

	localCmd.Flags().BoolVar(&showTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const localCommandName = "local"

const localCommandDocString = `nwpc_data_client local
Find local data path using config files in config dir.
`

var localCmd = &cobra.Command{
	Use:   localCommandName,
	Short: "Find local data path.",
	Long:  localCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if showTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")
		cmd.MarkFlagRequired("start-time")
		cmd.MarkFlagRequired("forecast-time")

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
	startTime, err := common.ParseStartTime(startTimeSting)
	if err != nil {
		log.Errorf("check startTime failed: %s", err)
		return
	}

	forecastTime, err = common.ParseForecastTime(forecastTimeString)
	if err != nil {
		log.Errorf("check forecastTime failed: %s", err)
		return
	}

	if len(configDir) == 0 {
		dataType = localCommandName + "/" + dataType
	}

	localConfig, err2 := common.LoadConfig(configDir, dataType)
	if err2 != nil {
		log.Errorf("load config failed: %v", err2)
		return
	}

	levels := strings.Split(locationLevels, ",")

	pathItem := common.FindLocalFile(localConfig, levels, startTime, forecastTime)
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
