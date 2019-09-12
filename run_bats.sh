#!/usr/bin/env bash

function run_bats() {
  for bats_file in $(find "$1" -name \*.bats); do
        echo "=> $bats_file"

        set +e
        bats "$bats_file"
        set -e
    done
}

export NWPC_DATA_CLIENT_BIN_DIR=$(pwd)/bin
export NWPC_DATA_CLIENT_PROGRAM="${NWPC_DATA_CLIENT_BIN_DIR}/nwpc_data_client"

export TESTS_BASE_DIR=tests

run_bats ${TESTS_BASE_DIR}