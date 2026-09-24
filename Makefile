VERSION    := $(shell sh scripts/version.sh 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

BIN_PATH := $(shell pwd)/bin

LDFLAGS := -X "github.com/cemc-oper/nwpc-data-client/common.Version=$(VERSION)" \
           -X "github.com/cemc-oper/nwpc-data-client/common.BuildTime=$(BUILD_TIME)" \
           -X "github.com/cemc-oper/nwpc-data-client/common.GitCommit=$(GIT_COMMIT)"

.PHONY: all command generate test test-integration clean
all: command

command: $(BIN_PATH)/nwpc_data_client

$(BIN_PATH)/nwpc_data_client: main.go
	@mkdir -p $(BIN_PATH)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $@ .

generate:
	cd common/config/generate && go generate

test:
	go test ./...

test-integration: $(BIN_PATH)/nwpc_data_client
	go test -tags=integration -count=1 -v ./tests/integration/...

clean:
	rm -rf $(BIN_PATH)
