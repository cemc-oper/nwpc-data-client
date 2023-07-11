.PHONY: data_checker
all: data_checker

data_checker:
	go build \
		-ldflags "-X \"github.com/cemc-oper/nwpc-data-client/common.Version=${VERSION}\" \
		-X \"github.com/cemc-oper/nwpc-data-client/common.BuildTime=${BUILD_TIME}\" \
		-X \"github.com/cemc-oper/nwpc-data-client/common.GitCommit=${GIT_COMMIT}\" " \
		-o ${BIN_PATH}/nwpc_data_checker \
		main.go