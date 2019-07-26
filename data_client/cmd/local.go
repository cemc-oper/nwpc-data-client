package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func init() {
	rootCmd.AddCommand(localCmd)

	localCmd.Flags().SortFlags = false

	localCmd.Flags().StringVar(&ConfigDir, "config-dir", "",
		"Config dir")

	localCmd.Flags().StringVar(&DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")

	localCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

const localCommandDocString = `nwpc_data_client local
Find local data path using config files in config dir.

Args:
    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFF, such as 000`

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Find local data path.",
	Long:  localCommandDocString,
	Args: func(cmd *cobra.Command, args []string) error {
		if ShowTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		StartTime, err = common.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = common.CheckForecastHour(args[1])
		if err != nil {
			return fmt.Errorf("check ForecastTime failed: %s", err)
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		if ShowTypes {
			showDataTypes(cmd, args)
		} else {
			findLocalFile(cmd, args)
		}
	},
}

func findLocalFile(cmd *cobra.Command, args []string) {
	configFilePath, err := common.FindConfig(ConfigDir, DataType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "model data type config is not found.\n")
		return
	}
	config, err2 := common.LoadConfig(configFilePath)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err2)
		return
	}
	filePath := findFile(config, StartTime, ForecastTime)
	fmt.Printf("%s\n", filePath)
}

func findFile(config common.DataConfig, startTime time.Time, forecastTime time.Duration) string {
	tpVar := common.GenerateTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
		return config.Default
	}
	fileName := fileNameBuilder.String()

	for _, item := range config.Paths {
		path := item.Path
		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))

		var dirPathBuilder strings.Builder
		err = dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dir path template execute has error: %s\n", err)
			continue
		}
		dirPath := dirPathBuilder.String()
		filePath := filepath.Join(dirPath, fileName)
		//fmt.Printf("%s\n", filePath)

		if common.CheckLocalFile(filePath) {
			return filePath
		}
	}

	return config.Default
}

func showDataTypes(cmd *cobra.Command, args []string) {
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

	err := filepath.Walk(ConfigDir, walkConfigDirectory)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Walk config directory has error: %s\n", err)
		return
	}
	for _, configPath := range configFilePaths {
		relConfigPath, err2 := filepath.Rel(ConfigDir, configPath)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Get rel path failed: %s\n", err2)
			continue
		}
		fmt.Printf("%s\n", relConfigPath[:len(relConfigPath)-5])
	}
}
