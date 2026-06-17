// Package data_service provides the gRPC API for NWPC data services.
//
//go:generate protoc --proto_path=. --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative data_service.proto
package data_service
