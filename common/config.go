package common

import (
	"fmt"
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
