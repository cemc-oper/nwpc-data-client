package cmd

import (
	"errors"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/data_client"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

var (
	ConfigDir    = ""
	DataType     = ""
	ShowTypes    = false
	StartTime    time.Time
	ForecastTime = ""
)

func init() {
	rootCmd.AddCommand(localCmd)
	localCmd.Flags().StringVar(&ConfigDir, "config-dir", "",
		"Config dir")
	localCmd.Flags().StringVar(&DataType, "data-type", "",
		"Data type used to locate config file path in config dir.")
	localCmd.Flags().BoolVar(&ShowTypes, "show-types", false,
		"Show supported data types defined in config dir and exit.")
}

var localCmd = &cobra.Command{
	Use:   "local",
	Short: "Find local data path.",
	Long: `Find local data path using config files in config dir.

    start_time: YYYYMMDDHH, such as 2018080100
    forecast_time: FFF, such as 000`,
	Args: func(cmd *cobra.Command, args []string) error {
		if ShowTypes {
			return nil
		}

		cmd.MarkFlagRequired("data-type")

		if len(args) != 2 {
			return errors.New("requires two arguments")
		}
		var err error
		StartTime, err = data_client.CheckStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = data_client.CheckForecastTime(args[1])
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
	configFilePath, err := findConfig(ConfigDir, DataType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "model data type config is not found.\n")
		return
	}
	localDataConfig, err2 := loadConfig(configFilePath)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err2)
		return
	}
	filePath := findFile(localDataConfig, StartTime, ForecastTime)
	fmt.Printf("%s\n", filePath)
}

func findConfig(configDir string, dataType string) (string, error) {
	configFilePath := filepath.Join(configDir, dataType+".yaml")
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return configFilePath, fmt.Errorf("file is not exist")
	}
	return configFilePath, nil
}

type LocalDataConfig struct {
	Default  string   `yaml:"default"`
	FileName string   `yaml:"file_name"`
	Paths    []string `yaml:"paths"`
}

func loadConfig(configFilePath string) (LocalDataConfig, error) {
	localDataConfig := LocalDataConfig{}

	data, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return localDataConfig, err
	}

	err = yaml.Unmarshal(data, &localDataConfig)
	if err != nil {
		return localDataConfig, err
	}

	return localDataConfig, nil
}

func findFile(config LocalDataConfig, startTime time.Time, forecastTime string) string {
	tpVar := data_client.GenerateTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
		return config.Default
	}
	fileName := fileNameBuilder.String()

	for _, path := range config.Paths {
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
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			continue
		}
		return filePath
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
