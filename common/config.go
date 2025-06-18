package common

import (
	"fmt"
	"github.com/cemc-oper/nwpc-data-client/common/config"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

const ConfigFileBasename = ".yaml"

type PathItem struct {
	PathType  string `yaml:"type"`
	LevelType string `yaml:"level"`
	Path      string `yaml:"path"`
}

type DataConfig struct {
	Default   string     `yaml:"default"`
	FileName  string     `yaml:"file_name"`
	FileNames []string   `yaml:"file_names"`
	Paths     []PathItem `yaml:"paths"`
}

// LoadConfigContent read config content from a local file or from embedded config strings.
func LoadConfigContent(configDir string, dataType string) (string, error) {
	if len(configDir) == 0 {
		content, err := findEmbeddedConfig(dataType)
		if err != nil {
			return "", fmt.Errorf("find embedded config for data type %s has error: %v", dataType, err)
		}
		return content, nil
	}

	content, err := loadContentFromLocal(configDir, dataType)
	if err != nil {
		return "", err
	}

	return content, nil
}

// loadContentFromLocal read config content from a local config file.
func loadContentFromLocal(configDir string, dataType string) (string, error) {
	configFilePath, err := findLocalConfig(configDir, dataType)
	if err != nil {
		return "", fmt.Errorf("model data type config is not found: %v", err)
	}

	data, err := os.ReadFile(configFilePath)
	content := string(data)
	if err != nil {
		return "", err
	}
	return content, nil
}

// ParseConfigContent parse config content template with start time and forecast time.
func ParseConfigContent(content string, startTime time.Time, forecastTime time.Duration, member string) (DataConfig, error) {
	dataConfig := DataConfig{}
	tpVar := GenerateConfigTemplateVariable(startTime, forecastTime, member)

	forecastTimeStr := FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecastTime": forecastTimeStr})

	contentTemplate := template.Must(template.New("dataConfig").Funcs(template.FuncMap{
		"generateStartTime":    GenerateStartTime,
		"getYear":              GetYear,
		"getMonth":             GetMonth,
		"getDay":               GetDay,
		"getHour":              GetHour,
		"generateForecastTime": GenerateForecastTime,
		"getForecastHour":      GetForecastHour,
		"getForecastMinute":    GetForecastMinute,
	}).Delims("{", "}").Parse(content))
	var configBuilder strings.Builder
	err := contentTemplate.Execute(&configBuilder, tpVar)
	if err != nil {
		currentLog.Errorf("file name template execute has error: %s", err)
		return dataConfig, err
	}
	parsedContent := configBuilder.String()

	dataConfig, err = LoadDataConfigFromContent(parsedContent)
	if err != nil {
		currentLog.Errorf("load dataConfig from content has error: %s", err)
		return dataConfig, err
	}
	return dataConfig, nil
}

// findEmbeddedConfig find data config string in embedded configs.
func findEmbeddedConfig(dataType string) (string, error) {
	for _, configItem := range config.EmbeddedConfigs {
		if configItem[0] == dataType {
			return configItem[1], nil
		}
	}
	return "", fmt.Errorf("can't find embedded config: %s", dataType)
}

// LoadDataConfigFromContent parse DataConfig from string content.
func LoadDataConfigFromContent(content string) (DataConfig, error) {
	dataConfig := DataConfig{}
	err := yaml.Unmarshal([]byte(content), &dataConfig)
	if err != nil {
		return dataConfig, err
	}
	return dataConfig, nil
}

func inLocationLevelTypes(locationLevelTypes []string, locationLevelType string) bool {
	if len(locationLevelTypes) == 0 {
		return true
	}
	for _, item := range locationLevelTypes {
		if item == "" || item == "all" {
			return true
		}
		if item == locationLevelType {
			return true
		}
	}
	return false
}

func LoadConfig(configDir string, dataType string) (DataConfig, error) {
	dataConfig := DataConfig{}

	var configObject DataConfig
	var err error

	if len(configDir) == 0 {
		configObject, err = loadEmbeddedConfig(dataType)
	} else {
		configObject, err = loadLocalConfig(configDir, dataType)
	}

	if err != nil {
		return dataConfig, fmt.Errorf("load config failed: %v", err)
	}
	return configObject, nil
}

func loadLocalConfig(configDir string, dataType string) (DataConfig, error) {
	dataConfig := DataConfig{}

	configFilePath, err := findLocalConfig(configDir, dataType)
	if err != nil {
		return dataConfig, fmt.Errorf("model data type config is not found: %v", err)
	}

	data, err := os.ReadFile(configFilePath)
	if err != nil {
		return dataConfig, err
	}

	err = yaml.Unmarshal(data, &dataConfig)
	if err != nil {
		return dataConfig, err
	}

	return dataConfig, nil
}

func findLocalConfig(configDir string, dataType string) (string, error) {
	configDirPath, _ := filepath.Abs(configDir)
	configFilePath := filepath.Join(configDirPath, dataType+ConfigFileBasename)
	relPath, err := filepath.Rel(configDirPath, configFilePath)
	if err != nil {
		return "", fmt.Errorf("check directory match has error: %v", err)
	}
	if relPath[0:2] == ".." {
		return "", fmt.Errorf("check directory failed: dataType (%s) should be under configDir (%s)",
			dataType, configDirPath)
	}
	if _, err := os.Stat(configFilePath); os.IsNotExist(err) {
		return configFilePath, fmt.Errorf("file is not exist")
	}
	return configFilePath, nil
}

func loadEmbeddedConfig(dataType string) (DataConfig, error) {
	dataConfig := DataConfig{}

	configContent, err := findEmbeddedConfig(dataType)
	if err != nil {
		return dataConfig, fmt.Errorf("model data type config is not found: %v", err)
	}

	err = yaml.Unmarshal([]byte(configContent), &dataConfig)
	if err != nil {
		return dataConfig, err
	}

	return dataConfig, nil
}
