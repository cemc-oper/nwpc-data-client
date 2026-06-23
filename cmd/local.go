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

	localCmd.Flags().StringVar(&localConfigDir, "config-dir", "",
		"Config dir.")

	localCmd.Flags().StringVar(&localConfigFile, "data-config-file", "",
		"Data config file path. If set, --config-dir and --data-type are ignored.")

	localCmd.Flags().StringVar(&localDataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	localCmd.Flags().StringVar(&localLocationLevels, "location-level", "",
		"Location levels, split by ',', such as 'runtime,archive'.")

	localCmd.Flags().StringVar(&localStartTimeString, "start-time", "",
		"start time, YYYYMMDDHH, such as 2020021400")
	localCmd.Flags().StringVar(&localForecastTimeString, "forecast-time", "",
		"forecast time, FFFh, such as 0h, 120h")
	localCmd.Flags().StringVar(&localMember, "member", "",
		"ensemble member, MMM, such as 000, 014")

	localCmd.Flags().BoolVar(&localShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")

	localCmd.Flags().BoolVar(&localDebugMode, "debug", false, "debug mode")
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
		if localShowTypes {
			return nil
		}

		if localConfigFile == "" {
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
	if localConfigFile != "" {
		configContent, err = common.LoadConfigContentFromFile(localConfigFile)
	} else {
		dataType := localDataType
		if len(localConfigDir) == 0 {
			dataType = localCommandName + "/" + dataType
		}
		configContent, err = common.LoadConfigContent(localConfigDir, dataType)
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
	if len(localConfigDir) == 0 {
		showEmbeddedDataTypes()
	} else {
		showLocalDataTypes(localConfigDir)
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
