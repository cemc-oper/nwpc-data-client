package common

import (
	log "github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func FindLocalFile(config DataConfig, locationLevels []string, startTime time.Time, forecastTime time.Duration) PathItem {

	tpVar := GenerateTimeTemplateVariable(startTime, forecastTime)

	forecastTimeStr := FormatForecastTimeShort(forecastTime)
	currentLog := log.WithFields(log.Fields{"forecastTime": forecastTimeStr})

	var configFileNames []string
	if len(config.FileNames) > 0 {
		currentLog.Debugf("using file_names...")
		configFileNames = config.FileNames
	}
	if len(config.FileName) > 0 {
		currentLog.Debugf("using file_name...")
		configFileNames = append(configFileNames, config.FileName)
	}

	var fileNames []string
	for _, item := range configFileNames {
		fileNameTemplate := template.Must(template.New("fileName").Funcs(template.FuncMap{
			"generateStartTime":    GenerateStartTime,
			"getYear":              GetYear,
			"getMonth":             GetMonth,
			"getDay":               GetDay,
			"getHour":              GetHour,
			"generateForecastTime": GenerateForecastTime,
			"getForecastHour":      GetForecastHour,
			"getForecastMinute":    GetForecastMinute,
		}).Delims("{", "}").Parse(item))
		var fileNameBuilder strings.Builder
		err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
		if err != nil {
			currentLog.Errorf("file name template execute has error: %s", err)
			return PathItem{
				Path:     config.Default,
				PathType: config.Default,
			}
		}
		fileName := fileNameBuilder.String()
		fileNames = append(fileNames, fileName)
		currentLog.Debugf("find file name: %s", fileName)
	}

	for _, item := range config.Paths {
		path := item.Path
		pathType := item.PathType
		locationLevelType := item.LevelType

		if !inLocationLevelTypes(locationLevels, locationLevelType) {
			currentLog.Debugf("skip location level type %s for %s", locationLevelType, path)
			continue
		}

		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))
		var dirPathBuilder strings.Builder
		err := dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			currentLog.Errorf("dir path template execute has error: %s", err)
			continue
		}

		dirPath := dirPathBuilder.String()
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
