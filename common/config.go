package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

const ConfigFileBasename = ".yaml"

func FindConfig(configDir string, dataType string) (string, error) {
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

type PathItem struct {
	PathType string `yaml:"type"`
	Path     string `yaml:"path"`
}

type DataConfig struct {
	Default  string     `yaml:"default"`
	FileName string     `yaml:"file_name"`
	Paths    []PathItem `yaml:"paths"`
}

func LoadConfig(configFilePath string) (DataConfig, error) {
	dataConfig := DataConfig{}

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
