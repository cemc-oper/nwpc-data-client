include build.mk

BIN_NAMES := $(foreach bin,$(BINARIES),$(word 1,$(subst :, ,$(bin))))
BIN_TARGETS := $(addprefix $(BIN_PATH)/,$(BIN_NAMES))

.PHONY: all command generate test clean
all: command

command: $(BIN_TARGETS)

$(BIN_PATH)/nwpc_data_client: data_client/main.go
$(BIN_PATH)/nwpc_data_checker: data_checker/main.go
$(BIN_PATH)/nwpc_data_server: data_service/data_server/main.go

$(BIN_TARGETS):
	@mkdir -p $(BIN_PATH)
	CGO_ENABLED=0 go build -ldflags "$(LDFLAGS)" -o $@ $<

generate:
	cd common/config/generate && go generate

test:
	cd tests/bats && ./run_bats.sh

clean:
	rm -rf $(BIN_PATH)
