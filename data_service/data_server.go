package data_service

import (
	"context"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"os"
	"path/filepath"
	"strings"
	"text/template"
	"time"
)

type NWPCDataServer struct {
	ConfigDir string
}

func (s *NWPCDataServer) FindDataPath(ctx context.Context, req *DataRequest) (*DataPathResponse, error) {
	dataType := req.GetDataType()

	fmt.Printf("find data path for: %s\n", dataType)

	emptyResponse := DataPathResponse{LocationType: "unknown", Location: "unknown"}

	configFilePath, err := common.FindConfig(s.ConfigDir, dataType)
	if err != nil {
		fmt.Fprintf(os.Stderr, "model data type config is not found.\n")
		return &emptyResponse, fmt.Errorf("model data type config is not found")
	}

	hpcDataConfig, err := common.LoadHpcConfig(configFilePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load config failed: %s\n", err)
		return &emptyResponse, fmt.Errorf("load config failed: %v", err)
	}

	startTime, err := common.CheckStartTime(req.GetStartTime())
	if err != nil {
		return &emptyResponse, fmt.Errorf("check StartTime failed: %s", err)
	}

	forecastTime, err := common.CheckForecastTime(req.GetForecastTime())
	if err != nil {
		return &emptyResponse, fmt.Errorf("check ForecastTime failed: %s", err)
	}

	filePath := findFile(hpcDataConfig, startTime, forecastTime)
	fmt.Printf("%s\n", filePath.PathType)
	fmt.Printf("%s\n", filePath.Path)

	return &DataPathResponse{
		LocationType: filePath.PathType,
		Location:     filePath.Path,
	}, nil
}

func findFile(config common.HpcDataConfig, startTime time.Time, forecastTime time.Duration) common.HpcPathItem {
	tpVar := common.GenerateTemplateVariable(startTime, forecastTime)

	fileNameTemplate := template.Must(template.New("fileName").
		Delims("{", "}").Parse(config.FileName))

	var fileNameBuilder strings.Builder
	err := fileNameTemplate.Execute(&fileNameBuilder, tpVar)
	if err != nil {
		fmt.Fprintf(os.Stderr, "file name template execute has error: %s\n", err)
		return common.HpcPathItem{
			Path:     config.Default,
			PathType: config.Default,
		}
	}
	fileName := fileNameBuilder.String()

	for _, pathItem := range config.Paths {
		path := pathItem.Path
		pathType := pathItem.PathType
		dirPathTemplate := template.Must(template.New("dirPath").Delims("{", "}").Parse(path))

		var dirPathBuilder strings.Builder
		err = dirPathTemplate.Execute(&dirPathBuilder, tpVar)
		if err != nil {
			fmt.Fprintf(os.Stderr, "dir path template execute has error: %s\n", err)
			continue
		}
		dirPath := dirPathBuilder.String()
		filePath := filepath.Join(dirPath, fileName)
		//fmt.Printf("%s\n", filePath)

		if common.CheckLocalFile(filePath) {
			return common.HpcPathItem{
				Path:     filePath,
				PathType: pathType,
			}
		}
	}

	return common.HpcPathItem{
		Path:     config.Default,
		PathType: config.Default,
	}
}
