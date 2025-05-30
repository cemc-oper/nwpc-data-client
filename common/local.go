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

	var fileNames []string

	if len(config.FileNames) > 0 {
		log.Debugf("using file_names...\n")
		for _, item := range config.FileNames {
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
				log.Errorf("file name template execute has error: %s", err)
				return PathItem{
					Path:     config.Default,
					PathType: config.Default,
				}
			}
			fileName := fileNameBuilder.String()
			fileNames = append(fileNames, fileName)
			log.Debugf("find file name %s\n", fileName)
		}
	}

	if len(config.FileName) > 0 {
		log.Debugf("using file_name...\n")
		fileNameTemplate := template.Must(template.New("fileName").Funcs(template.FuncMap{
			"generateStartTime":    GenerateStartTime,
			"getYear":              GetYear,
			"getMonth":             GetMonth,
			"getDay":               GetDay,
			"getHour":              GetHour,
			"generateForecastTime": GenerateForecastTime,
			"getForecastHour":      GetForecastHour,
			"getForecastMinute":    GetForecastMinute,
		}).Delims("{", "}").Parse(config.FileName))
		var fileNameBuilder strings.Builder
		err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
		if err != nil {
			log.Errorf("file name template execute has error: %s", err)
			return PathItem{
				Path:     config.Default,
				PathType: config.Default,
			}
		}
		fileName := fileNameBuilder.String()
		fileNames = append(fileNames, fileName)
		log.Debugf("find file name %s\n", fileName)
	}

	for _, item := range config.Paths {
		path := item.Path
		pathType := item.PathType
		locationLevelType := item.LevelType

		if !inLocationLevelTypes(locationLevels, locationLevelType) {
			log.Debugf("skip location level type %s\n", locationLevelType)
			continue
		}

		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))

		var dirPathBuilder strings.Builder
		err := dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			log.Errorf("dir path template execute has error: %s", err)
			continue
		}
		dirPath := dirPathBuilder.String()
		for _, fileName := range fileNames {
			filePath := filepath.Join(dirPath, fileName)
			//log.Infof("%s\n", filePath)

			if CheckLocalFile(filePath) {
				return PathItem{
					Path:     filePath,
					PathType: pathType,
				}
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
		log.Debugf("check file path failed: %s", filePath)
		return false
	}
	log.Debugf("check file path success: %s", filePath)
	return true
}

func PrepareLocalDir(filePath string) {
	localFileDir := filepath.Dir(filePath)
	_ = os.MkdirAll(localFileDir, os.ModeDir)
}
