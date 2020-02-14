package common

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

func FindHpcFile(
	config DataConfig, locationLevels []string, startTime time.Time, forecastTime time.Duration,
	storageUser string, storageHost string, privateKeyFilePath string, hostKeyFilePath string,
) PathItem {
	tpVar := GenerateTimeTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
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
			fmt.Fprintf(os.Stderr, "dir path template execute has error: %s\n", err)
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
			fmt.Fprintf(os.Stderr, "path type is not supported: %s", pathType)
		}
	}

	return PathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}
