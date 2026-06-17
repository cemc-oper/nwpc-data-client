package common

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"time"
)

func FindLocalFile(config DataConfig, locationLevels []string, startTime time.Time, forecastTime time.Duration) PathItem {
	forecastTimeStr := FormatForecastTimeShort(forecastTime)
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
				return PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			} else {
				currentLog.Debugf("check file path failed: %s", filePath)
			}
		}
	}

	return PathItem{
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
	_ = os.MkdirAll(localFileDir, os.ModeDir)
}
