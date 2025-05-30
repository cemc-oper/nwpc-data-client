#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
check_date_time=${check_date}00

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}

@test "test cma_tym/current/bin/postvar runtime" {
  expected_result="/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_DATA/grapes_d01/dat/postvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    --data-type=cma_tym/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

@test "test cma_tym/current/bin/postvar archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/${check_date_time}/postvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    --data-type=cma_tym/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

# postvar with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_tym/bin/postvar runtime with config" {
  expected_result="/g2/op_meso/OPER/WORKDIR/NWP_CMA_TYM_DATA/grapes_d01/dat/postvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    "${config}" \
    --data-type=cma_tym/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}


@test "test cma_tym/current/bin/postvar archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/TYM/Fcst-main/${check_date_time}/postvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    "${config}" \
    --data-type=cma_tym/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}