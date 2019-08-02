package common

import (
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common/config"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ConfigFileBasename = ".yaml"

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

	data, err := ioutil.ReadFile(configFilePath)
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

func findEmbeddedConfig(dataType string) (string, error) {
	for _, configItem := range config.EmbeddedConfigs {
		if configItem[0] == dataType {
			return configItem[1], nil
		}
	}
	return "", fmt.Errorf("can't find embedded config: %s", dataType)
}

type PathItem struct {
	PathType string `yaml:"type"`
	Path     string `yaml:"path"`
}

type DataConfig struct {
	Default  string     `yaml:"default"`
	FileName string     `yaml:"file_name"`
	Paths    []PathItem `yaml:"paths"`
}
