package main

import (
	"context"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

func main() {
	const address = "10.40.140.43:33383"
	const remoteFilePath = "/sstorage1/COMMONDATA/OPER/old/nwp/" +
		"GMFS_GRIB2_GRAPES/CMACAST/GRAPES_GFS_forCAST_2019061418/" +
		"ne_gmf.gra.2019061418000.grb2"

	localFilePath := "./dist/ne_gmf.gra.2019061418000.grb2"

	prepareLocalDir(localFilePath)

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := data_service.NewNWPCDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	stream, err := c.DownloadFile(ctx, &data_service.FileContentRequest{
		FilePath: remoteFilePath,
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
