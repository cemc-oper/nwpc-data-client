package main

import (
	"context"
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"google.golang.org/grpc"
	"log"
	"time"
)

func main() {
	conn, err := grpc.Dial("localhost:30105", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	defer conn.Close()

	c := data_service.NewNWPCDataServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Minute*2)
	defer cancel()

	r, err := c.FindDataPath(ctx, &data_service.DataRequest{
		DataType:     "gmf_grapes_gfs/bin/modelvar",
		StartTime:    "2019072200",
		ForecastTime: "000h",
	})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}

	log.Printf("Response: %v", r)
}
