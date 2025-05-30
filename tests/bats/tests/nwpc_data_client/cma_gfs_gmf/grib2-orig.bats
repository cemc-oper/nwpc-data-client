#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}

@test "test cma_gfs_gmf/grib2/orig runtime" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/${check_date_time}/data/output/grib2_orig/gmf.gra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
  assert_output ${expected_result}
}

@test "test cma_gfs_gmf/current/grib2/orig archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/${check_date_time}/ORIG/gmf.gra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi

  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
  assert_output ${expected_result}
}

# grib2 orig with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_gfs_gmf/current/grib2/orig runtime with config" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/${check_date_time}/data/output/grib2_orig/gmf.gra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
      "${config}" \
      --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time "${check_date_time}" \
      --forecast-time 3h
    assert_output ${expected_result}
}

@test "test cma_gfs_gmf/current/grib2/orig archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/${check_date_time}/ORIG/gmf.gra.${check_date_time}003.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    "${config}" \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}



