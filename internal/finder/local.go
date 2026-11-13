package finder

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

	for _, item := range CandidatePaths(config, locationLevels) {
		if item.PathType != "local" {
			currentLog.Debugf("skip path type %s for %s", item.PathType, item.Path)
			continue
		}
		currentLog.Debugf("find file path: %s", item.Path)
		if CheckLocalFile(item.Path) {
			currentLog.Debugf("find file path success: %s", item.Path)
			return item
		}
		currentLog.Debugf("find file path failed: %s", item.Path)
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

func InLocationLevelTypes(locationLevelTypes []string, locationLevelType string) bool {
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
