#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
check_date_time=${check_date}00

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}


@test "test cma_tym/current/grib2/orig runtime" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_TYM_POST_DATA/${check_date_time}/data/output/grib2_orig/rmf.tcgra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    --data-type=cma_tym/current/grib2/orig \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

@test "test cma_tym/current/grib2/orig archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/TYM/Prod-grib/${check_date_time}/ORIG/rmf.tcgra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    --data-type=cma_tym/current/grib2/orig \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

# postvar with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_tym/current/grib2/orig runtime with config" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_TYM_POST_DATA/${check_date_time}/data/output/grib2_orig/rmf.tcgra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    "${config}" \
    --data-type=cma_tym/current/grib2/orig \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}


@test "test cma_tym/current/grib2/orig archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/TYM/Prod-grib/${check_date_time}/ORIG/rmf.tcgra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    "${config}" \
    --location-level=archive \
    --data-type=cma_tym/current/grib2/orig \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}