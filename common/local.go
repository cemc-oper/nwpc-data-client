package common

import (
	"os"
	"path/filepath"
	"time"

	"github.com/cemc-oper/nwpc-data-client/internal/dataconfig"
	log "github.com/sirupsen/logrus"
)

func FindLocalFile(config dataconfig.DataConfig, locationLevels []string, startTime time.Time, forecastTime time.Duration) dataconfig.PathItem {
	forecastTimeStr := dataconfig.FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecastTime": forecastTimeStr})

	var fileNames []string
	if len(config.FileNames) > 0 {
		currentLog.Debugf("using file_names...")
		fileNames = config.FileNames
	}
	if len(config.FileName) > 0 {
		currentLog.Debugf("using file_name...")
		fileNames = append(fileNames, config.FileName)
	}

	for _, item := range config.Paths {
		dirPath := item.Path
		pathType := item.PathType
		locationLevelType := item.LevelType

		if !inLocationLevelTypes(locationLevels, locationLevelType) {
			currentLog.Debugf("skip location level type %s for %s", locationLevelType, dirPath)
			continue
		}

		for _, fileName := range fileNames {
			filePath := filepath.Join(dirPath, fileName)
			currentLog.Debugf("find file path: %s", filePath)

			if CheckLocalFile(filePath) {
				currentLog.Debugf("check file path success: %s", filePath)
				return dataconfig.PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			} else {
				currentLog.Debugf("check file path failed: %s", filePath)
			}
		}
	}

	return dataconfig.PathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}

func CheckLocalFile(filePath string) bool {
	if _, err := os.Stat(filePath); err != nil {
		return false
	}
	return true
}

func PrepareLocalDir(filePath string) {
	localFileDir := filepath.Dir(filePath)
	_ = os.MkdirAll(localFileDir, 0755)
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
