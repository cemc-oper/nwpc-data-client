package data_service

import (
	"context"
	"fmt"
	"github.com/nwpc-oper/nwpc-data-client/common"
	"log"
	"os"
)

type NWPCDataServer struct {
	ConfigDir string
}

func (s *NWPCDataServer) FindDataPath(ctx context.Context, req *DataRequest) (*DataPathResponse, error) {
	dataType := req.GetDataType()
	startTimeString := req.GetStartTime()
	forecastTimeString := req.GetForecastTime()

	log.Printf("FindDataPath for type %s: %s %s\n", dataType, startTimeString, forecastTimeString)

	response, err := s.findDataPath(req)

	log.Printf("Find data path type: %s\n", response.LocationType)
	log.Printf("Find data path: %s\n", response.Location)

	return response, err
}

func (s *NWPCDataServer) GetDataFileInfo(ctx context.Context, req *DataRequest) (*DataFileResponse, error) {
	log.Printf("GetDataFileInfo for %s", req)

	dataResponse, err := s.findDataPath(req)
	if err != nil {
		return &DataFileResponse{
			Status:       StatusCode_Failed,
			ErrorMessage: fmt.Sprintf("%v", err),
			FilePath:     "",
			FileSize:     -1,
		}, nil
	}

	filePath := dataResponse.Location

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

func (s *NWPCDataServer) DownloadDataFile(req *DataRequest, stream NWPCDataService_DownloadDataFileServer) error {

	log.Printf("DownloadFile for %s", req)

	dataResponse, err := s.findDataPath(req)
	if err != nil {
		return err
	}

	filePath := dataResponse.Location
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
	defer file.Close()

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

func (s *NWPCDataServer) findDataPath(req *DataRequest) (*DataPathResponse, error) {
	dataType := req.GetDataType()
	startTimeString := req.GetStartTime()
	forecastTimeString := req.GetForecastTime()

	emptyResponse := DataPathResponse{LocationType: "NOTFOUND", Location: "NOTFOUND"}

	hpcDataConfig, err := common.LoadConfig(s.ConfigDir, dataType)
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

	filePath := common.FindLocalFile(hpcDataConfig, startTime, forecastTime)

	return &DataPathResponse{
		LocationType: filePath.PathType,
		Location:     filePath.Path,
	}, nil
}
