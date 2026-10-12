VERSION    := $(shell sh scripts/version.sh 2>/dev/null || echo "dev")
BUILD_TIME := $(shell date -u '+%Y-%m-%dT%H:%M:%SZ' 2>/dev/null || echo "unknown")
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

BIN_PATH := $(shell pwd)/bin

LDFLAGS := -X "github.com/cemc-oper/nwpc-data-client/internal/buildinfo.Version=$(VERSION)" \
           -X "github.com/cemc-oper/nwpc-data-client/internal/buildinfo.BuildTime=$(BUILD_TIME)" \
           -X "github.com/cemc-oper/nwpc-data-client/internal/buildinfo.GitCommit=$(GIT_COMMIT)"

.PHONY: all command generate test test-e2e test-integration clean docs-cli docs
all: command

command: $(BIN_PATH)/nwpc_data_client

GO_SOURCES := $(shell find . -type f -name '*.go' -not -path './bin/*')

$(BIN_PATH)/nwpc_data_client: $(GO_SOURCES)
	@mkdir -p $(BIN_PATH)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $@ .

generate:
	cd internal/embedded/generate && go generate

test:
	go test ./...

test-e2e: $(BIN_PATH)/nwpc_data_client
	go test -count=1 ./tests/e2e/...

test-integration: $(BIN_PATH)/nwpc_data_client
	go test -tags=integration -count=1 -v ./tests/integration/...

# Regenerate the CLI reference pages under docs/reference/cli/.
docs-cli:
	go run ./tools/gen-cli-docs ./docs/reference/cli

# Build the Sphinx documentation (requires: pip install -r docs/requirements.txt).
docs: docs-cli
	sphinx-build -b html docs docs/_build/html

clean:
	rm -rf $(BIN_PATH) docs/_build
