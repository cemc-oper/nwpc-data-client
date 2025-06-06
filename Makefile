all: command
.PHONY: nwpc_data_client nwpc_data_server nwpc_data_checker generate test

export VERSION := $(shell cat VERSION)
export BUILD_TIME := $(shell date --utc --rfc-3339 ns 2> /dev/null | sed -e 's/ /T/')
export GIT_COMMIT := $(shell git rev-parse --short HEAD 2> /dev/null || true)

export BIN_PATH := $(shell pwd)/bin


command: nwpc_data_client nwpc_data_server nwpc_data_checker

nwpc_data_client:
	$(MAKE) -C data_client

nwpc_data_server:
	$(MAKE) -C data_service

nwpc_data_checker:
	$(MAKE) -C data_checker

generate:
	cd common/config/generate && go generate

test:
	cd tests/bats && ./run_bats.sh