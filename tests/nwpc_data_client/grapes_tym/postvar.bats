#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
check_date_time=${check_date}00

@test "test grapes_tym/bin/postvar runtime" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat/postvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_tym/bin/postvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_tym/bin/postvar archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/${check_date_time}/postvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        --data-type=grapes_tym/bin/postvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}

# postvar with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_tym/bin/postvar runtime with config" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat/postvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        "${config}" \
        --data-type=grapes_tym/bin/postvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}


@test "test grapes_tym/bin/postvar archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Fcst-main/${check_date_time}/postvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        "${config}" \
        --location-level=archive \
        --data-type=grapes_tym/bin/postvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}