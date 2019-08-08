package main

import (
	"context"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"google.golang.org/grpc"
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

	r, err := c.FindDataPath(ctx, &data_service.DataRequest{
		DataType:       "grapes_gfs_gmf/grib2/orig",
		LocationLevels: []string{"archive"},
		StartTime:      "2019080700",
		ForecastTime:   "120h",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Response: %v", r)
}
