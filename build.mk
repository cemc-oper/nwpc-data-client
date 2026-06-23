# build.mk — single source of truth for build configuration

VERSION    := $(shell git describe --tags --always --dirty 2>/dev/null || cat VERSION 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

BIN_PATH := $(shell pwd)/bin

LDFLAGS := -X "github.com/cemc-oper/nwpc-data-client/common.Version=$(VERSION)" \
           -X "github.com/cemc-oper/nwpc-data-client/common.BuildTime=$(BUILD_TIME)" \
           -X "github.com/cemc-oper/nwpc-data-client/common.GitCommit=$(GIT_COMMIT)"

# binary_name:source_dir
BINARIES := nwpc_data_client:.
