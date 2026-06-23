include build.mk

BIN_NAMES := $(foreach bin,$(BINARIES),$(word 1,$(subst :, ,$(bin))))
BIN_TARGETS := $(addprefix $(BIN_PATH)/,$(BIN_NAMES))

.PHONY: all command generate test test-integration clean
all: command

command: $(BIN_TARGETS)

$(BIN_PATH)/nwpc_data_client: main.go

$(BIN_TARGETS):
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
