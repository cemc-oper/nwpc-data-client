.PHONY: data_client
all: data_client

data_client:
	go build \
		-ldflags "-X \"github.com/cemc-oper/nwpc-data-client/common.Version=${VERSION}\" \
		-X \"github.com/cemc-oper/nwpc-data-client/common.BuildTime=${BUILD_TIME}\" \
		-X \"github.com/cemc-oper/nwpc-data-client/common.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nwpc_data_client \
		main.go