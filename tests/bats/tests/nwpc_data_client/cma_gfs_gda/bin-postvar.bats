#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}

hour_4dv=21
check_date_4dvar=$(date -d "-2 day" +%Y%m%d)
check_data_time_4dvar=${check_date_4dvar}${hour_4dv}

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}

# postvar

@test "test cma_gfs_gda/current/bin/postvar runtime" {
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_EN4DVAR/${hour}/output/postvar${check_data_time_4dvar}_006"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=cma_gfs_gda/current/bin/postvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result_today="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_EN4DVAR/${hour}/output/postvar${check_date}${hour_4dv}_006"
  if [ -f "${expected_result_today}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
      --data-type=cma_gfs_gda/current/bin/postvar \
      --start-time "${today_check_date}${hour}" \
      --forecast-time 3h
    assert_output ${expected_result_today}
    return
  fi

  skip "data is not available: ${expected_result} ${expected_result_today}"
}

@test "test cma_gfs_gda/current/bin/postvar archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/${check_date_time}/postvar${check_data_time_4dvar}_006"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
    --location-level=archive \
    --data-type=cma_gfs_gda/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
  }

# postvar with config

config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_gfs_gda/current/bin/postvar runtime with config" {
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_EN4DVAR/${hour}/output/postvar${check_data_time_4dvar}_006"
  if [ -f "${expected_result}" ]; then
    run ${NWPC_DATA_CLIENT_PROGRAM} local \
      --location-level=runtime \
      "${config}" \
      --data-type=cma_gfs_gda/current/bin/postvar \
      --start-time "${check_date_time}" \
      --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result_today="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_EN4DVAR/${hour}/output/postvar${check_date}${hour_4dv}_006"
  if [ -f "${expected_result_today}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=cma_gfs_gda/current/bin/postvar \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    assert_output ${expected_result_today}
    return
  fi

  skip "data is not available: ${expected_result} ${expected_result_today}"
}

@test "test cma_gfs_gda/current/bin/postvar archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GDA/GRAPES_EN4DVAR/Fcst-9h/${check_date_time}/postvar${check_data_time_4dvar}_006"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
    --location-level=archive \
    "${config}" \
    --data-type=cma_gfs_gda/current/bin/postvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

