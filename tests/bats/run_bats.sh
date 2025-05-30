#!/usr/bin/env bash

function run_bats() {
  for bats_file in $(find "$1" -name \*.bats); do
    relpath=$(realpath --relative-to="${TESTS_BASE_DIR}" "$bats_file")
    echo "=> ${relpath}"

    set +e
    ${BATS_PROGRAM} "$bats_file"
    set -e
  done
}

REPO_BASE_DIR="$( cd "$(dirname ${BASH_SOURCE[0]})/../.." && pwd )"

export NWPC_DATA_CLIENT_BIN_DIR=${REPO_BASE_DIR}/bin
export NWPC_DATA_CLIENT_PROGRAM="${NWPC_DATA_CLIENT_BIN_DIR}/nwpc_data_client"

export NWPC_DATA_CLIENT_CONFIG_DIR=${REPO_BASE_DIR}/conf

export TESTS_BASE_DIR=${REPO_BASE_DIR}/tests/bats/tests
export BATS_TOOL_BASE_DIR=${REPO_BASE_DIR}/tests/bats/tool
export BATS_PROGRAM=${BATS_TOOL_BASE_DIR}/bats/bin/bats

run_bats ${TESTS_BASE_DIR}