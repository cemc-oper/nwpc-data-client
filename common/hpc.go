package common

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type HpcPathItem struct {
	PathType string `yaml:"type"`
	Path     string `yaml:"path"`
}

type HpcDataConfig struct {
	Default  string        `yaml:"default"`
	FileName string        `yaml:"file_name"`
	Paths    []HpcPathItem `yaml:"paths"`
}

func LoadHpcConfig(configFilePath string) (HpcDataConfig, error) {
	dataConfig := HpcDataConfig{}

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
