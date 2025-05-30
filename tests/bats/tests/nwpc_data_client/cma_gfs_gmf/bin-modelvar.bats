#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}

@test "test cma_gfs_gmf/current/bin/modelvar runtime" {
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_GMFS/${hour}/output/modelvar${check_date_time}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
        --data-type=cma_gfs_gmf/current/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_GMFS/${hour}/output/modelvar${today_check_date}${hour}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=cma_gfs_gmf/current/bin/modelvar \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  skip "data is not available for ${check_date_time} and ${today_check_date}${hour}"
}

@test "test cma_gfs_gmf/current/bin/modelvar archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/${check_date_time}/modelvar${check_date_time}_003"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    --data-type=cma_gfs_gmf/current/bin/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

# modelvar with config

config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_gfs_gmf/current/bin/modelvar runtime with config" {
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_GMFS/${hour}/output/modelvar${check_date_time}_003"
  if [ -f "${expected_result}" ]; then
    run ${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=runtime \
        "${config}" \
        --data-type=cma_gfs_gmf/current/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/op_gfs/CMA-GFS/CMA-GFS4.2_DATA/MODEL/GRAPES_GMFS/${hour}/output/modelvar${today_check_date}${hour}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=cma_gfs_gmf/current/bin/modelvar \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  skip "data is not available for ${check_date_time} and ${today_check_date}"
}

@test "test cma_gfs_gmf/current/bin/modelvar archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Fcst-long/${check_date_time}/modelvar${check_date_time}_003"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
    --location-level=archive \
    "${config}" \
    --data-type=cma_gfs_gmf/current/bin/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}
