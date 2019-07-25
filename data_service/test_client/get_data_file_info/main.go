package main

import (
	"context"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("10.40.140.43:33383", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := data_service.NewNWPCDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	stream, err := c.DownloadFile(ctx, &data_service.FileContentRequest{
		FilePath: "/sstorage1/COMMONDATA/OPER/old/nwp/GMFS_GRIB2_GRAPES/CMACAST/GRAPES_GFS_forCAST_2019061418/ne_gmf.gra.2019061418000.grb2",
	})

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.DownloadFile(_) = _, %v", c, err)
		}
		log.Println(chunk.ChunkLength)
	}
}
