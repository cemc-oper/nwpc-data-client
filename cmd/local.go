package cmd

import (
	"errors"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
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
		StartTime, err = checkStartTime(args[0])
		if err != nil {
			return fmt.Errorf("check StartTime failed: %s", err)
		}

		ForecastTime, err = checkForecastTime(args[1])
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
		fmt.Printf("model data type config is not found.")
		return
	}
	localDataConfig, err2 := loadConfig(configFilePath)
	if err2 != nil {
		fmt.Printf("load config failed: %s", err2)
		return
	}
	filePath := findFile(localDataConfig, StartTime, ForecastTime)
	fmt.Printf("%s\n", filePath)
}

func checkStartTime(value string) (time.Time, error) {
	if len(value) != 10 {
		return time.Time{}, fmt.Errorf("length of start_time must be 10")
	}
	s, err := time.Parse("2006010215", value)
	if err != nil {
		return s, err
	}
	return s, nil
}

func checkForecastTime(value string) (string, error) {
	if len(value) > 3 {
		return "", fmt.Errorf("length of forecast time must less or equal to 3")
	}

	intValue, err := strconv.Atoi(value)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%03d", intValue), nil
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
	tpVar := generateTemplateObject(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").Parse(config.FileName))
	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s", err)
		return config.Default
	}
	fileName := fileNameBuilder.String()

	for _, path := range config.Paths {
		dirPathTemplate := template.Must(template.New("dirPath").Parse(path))
		var dirPathBuilder strings.Builder
		err := dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dir path template execute has error: %s", err)
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

type templateVariable struct {
	Year     string
	Month    string
	Day      string
	Hour     string
	Forecast string
	Year4DV  string
	Month4DV string
	Day4DV   string
	Hour4DV  string
}

func generateTemplateObject(startTime time.Time, forecastTime string) templateVariable {
	startTime4DV := startTime.Add(time.Hour * -3)
	tpVariable := templateVariable{
		Year:     startTime.Format("2006"),
		Month:    startTime.Format("01"),
		Day:      startTime.Format("02"),
		Hour:     startTime.Format("15"),
		Forecast: forecastTime,
		Year4DV:  startTime4DV.Format("2006"),
		Month4DV: startTime4DV.Format("01"),
		Day4DV:   startTime4DV.Format("02"),
		Hour4DV:  startTime4DV.Format("15"),
	}
	return tpVariable
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
		fmt.Fprintf(os.Stderr, "Walk config directory has error: %s", err)
		return
	}
	for _, configPath := range configFilePaths {
		relConfigPath, err2 := filepath.Rel(ConfigDir, configPath)
		if err2 != nil {
			fmt.Fprintf(os.Stderr, "Get rel path failed: %s", err2)
			continue
		}
		fmt.Printf("%s\n", relConfigPath[:len(relConfigPath)-5])
	}
}
