package data_service

import (
	"context"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"log"
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
	startTimeString := req.GetStartTime()
	forecastTimeString := req.GetForecastTime()

	log.Printf("find data path for type %s: %s %s\n", dataType, startTimeString, forecastTimeString)

	emptyResponse := DataPathResponse{LocationType: "NOTFOUND", Location: "NOTFOUND"}

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

	startTime, err := common.CheckStartTime(startTimeString)
	if err != nil {
		return &emptyResponse, fmt.Errorf("check StartTime failed: %s", err)
	}

	forecastTime, err := common.CheckForecastTime(forecastTimeString)
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

func (s *NWPCDataServer) GetDataFileInfo(ctx context.Context, req *DataFileRequest) (*DataFileResponse, error) {
	filePath := req.FilePath
	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return &DataFileResponse{
			Status:       StatusCode_Failed,
			ErrorMessage: fmt.Sprintf("check file error: %v", err),
			FilePath:     "",
			FileSize:     -1,
		}, nil
	}

	fileUnixPermission := fileInfo.Mode().Perm()

	if fileUnixPermission&004 == 0 {
		return &DataFileResponse{
			Status:       StatusCode_Failed,
			ErrorMessage: "don't have read permission",
			FilePath:     "",
			FileSize:     -1,
		}, nil
	}

	return &DataFileResponse{
		Status:       StatusCode_Success,
		ErrorMessage: "",
		FilePath:     filePath,
		FileSize:     fileInfo.Size(),
	}, nil
}

func (*NWPCDataServer) DownloadFile(req *FileContentRequest, stream NWPCDataService_DownloadFileServer) error {
	filePath := req.FilePath

	fileInfo, err := os.Stat(filePath)

	if err != nil {
		return fmt.Errorf("file error: %v", err)
	}

	fileUnixPermission := fileInfo.Mode().Perm()

	if fileUnixPermission&004 == 0 {
		return fmt.Errorf("don't have read permission")
	}

	const chunkSize = 64 * 1024

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	buffer := make([]byte, chunkSize)

	for {
		n, err := file.Read(buffer)

		if err != nil {
			break
		}

		_ = stream.Send(&FileContentResponse{
			ChunkLength: int64(n),
			Chunk:       buffer[:n],
		})
	}

	return nil
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
