package data_service

import (
	"context"
)

type NWPCDataServer struct {
}

func (s *NWPCDataServer) FindDataPath(ctx context.Context, req *DataRequest) (*DataPathResponse, error) {
	return &DataPathResponse{LocationType: "unknown", Location: "unknown"}, nil
}
