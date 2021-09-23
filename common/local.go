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
		for _, item := range config.FileNames {
			fileNameTemplate := template.Must(template.New("fileName").
				Delims("{", "}").Parse(item))
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
		}
	}

	if len(config.FileName) > 0 {
		fileNameTemplate := template.Must(template.New("fileName").
			Delims("{", "}").Parse(config.FileName))

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
	}

	for _, item := range config.Paths {
		path := item.Path
		pathType := item.PathType
		locationLevelType := item.LevelType

		if !inLocationLevelTypes(locationLevels, locationLevelType) {
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
			log.Debugf("%s\n", filePath)

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
		return false
	}

	return true
}

func PrepareLocalDir(filePath string) {
	localFileDir := filepath.Dir(filePath)
	_ = os.MkdirAll(localFileDir, os.ModeDir)
}
