package common

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func FindConfig(configDir string, dataType string) (string, error) {
	configFilePath := filepath.Join(configDir, dataType+".yaml")
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
