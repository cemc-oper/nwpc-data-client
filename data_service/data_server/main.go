package main

import (
	"github.com/nwpc-oper/nwpc-data-client/data_service"
	"google.golang.org/grpc"
	"log"
	"net"
)

func runServer() {
	lis, err := net.Listen("tcp", ":30105")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	data_service.RegisterNWPCDataServiceServer(s, &data_service.NWPCDataServer{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func main() {
	runServer()
}
