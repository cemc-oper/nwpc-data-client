#!/usr/bin/env bash

function run_bats() {
  for bats_file in $(find "$1" -name \*.bats); do
        echo "=> $bats_file"

        set +e
        bats "$bats_file"
        set -e
    done
}

run_bats tests