package common

import (
	log "github.com/sirupsen/logrus"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func FindHpcFile(
	config DataConfig, locationLevels []string, startTime time.Time, forecastTime time.Duration,
	storageUser string, storageHost string, privateKeyFilePath string, hostKeyFilePath string,
) PathItem {
	tpVar := GenerateConfigTemplateVariable(startTime, forecastTime, "")

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

	for _, pathItem := range config.Paths {
		path := pathItem.Path
		pathType := pathItem.PathType
		locationLevel := pathItem.LevelType
		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))

		if !inLocationLevelTypes(locationLevels, locationLevel) {
			continue
		}

		var dirPathBuilder strings.Builder
		err = dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			log.Errorf("dir path template execute has error: %s", err)
			continue
		}
		dirPath := dirPathBuilder.String()
		filePath := filepath.Join(dirPath, fileName)
		//fmt.Printf("%s\n", filePath)

		if pathType == "storage" {
			if CheckFileOverSSH(filePath, storageUser, storageHost, privateKeyFilePath, hostKeyFilePath) {
				return PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			}
		} else if pathType == "local" {
			// check if file exists
			if CheckLocalFile(filePath) {
				return PathItem{
					Path:     filePath,
					PathType: pathType,
				}
			}
		} else {
			log.Errorf("path type is not supported: %s", pathType)
		}
	}

	return PathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}
