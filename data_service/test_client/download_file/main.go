package main

import (
	"context"
	"github.com/cemc-oper/nwpc-data-client/data_service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"io"
	"os"
	"path/filepath"
	"time"
)

func main() {
	const address = "10.40.139.28:33483"
	localFilePath := "./dist/gmf.gra.2019080700120.grb2"

	prepareLocalDir(localFilePath)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := data_service.NewNWPCDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	stream, err := c.DownloadDataFile(ctx, &data_service.DataRequest{
		DataType:       "cma_gfs_gmf/grib2/orig",
		LocationLevels: []string{"all"},
		StartTime:      "2019080700",
		ForecastTime:   "120h",
	})

	f, err := os.OpenFile(localFilePath, os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.DownloadFile(_) = _, %v", c, err)
		}
		f.Write(chunk.Chunk)
		log.Println(chunk.ChunkLength)
	}

}

func prepareLocalDir(filePath string) {
	localFileDir := filepath.Dir(filePath)
	_ = os.MkdirAll(localFileDir, os.ModeDir)
}
