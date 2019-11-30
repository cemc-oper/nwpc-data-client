.PHONY: data_server
all: data_server

data_server:
	go build \
		-ldflags "-X \"github.com/nwpc-oper/nwpc-data-client/common.Version=${VERSION}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.BuildTime=${BUILD_TIME}\" \
		-X \"github.com/nwpc-oper/nwpc-data-client/common.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nwpc_data_server \
		data_server/main.go