all: nwpc_data_client nwpc_data_server
.PHONY: nwpc_data_client nwpc_data_server

VERSION := $(shell cat VERSION)
BUILD_TIME := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || true)


nwpc_data_client:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nwpc-data-client/common.Version=${VERSION}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.BuildTime=${BUILD_TIME}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.GitCommit=${GIT_COMMIT}\" " \
		-gcflags "all=-N -l" \
		-o bin/nwpc_data_client \
		main.go

nwpc_data_server:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nwpc-data-client/common.Version=${VERSION}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.BuildTime=${BUILD_TIME}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.GitCommit=${GIT_COMMIT}\" " \
		-gcflags "all=-N -l" \
		-o bin/nwpc_data_server \
		data_service/data_server/main.go