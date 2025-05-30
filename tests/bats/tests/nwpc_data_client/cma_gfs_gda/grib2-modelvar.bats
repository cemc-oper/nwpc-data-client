#!/usr/bin/env bats

check_date=$(date -d "-2 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}


setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}


# grib2 orig
@test "test cma_gfs_gda/current/grib2/orig runtime" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/${check_date_time}/data/output/grib2_orig/modelvar${check_date_time}000.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    --data-type=cma_gfs_gda/current/grib2/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 0h
  assert_output ${expected_result}
}

@test "test cma_gfs_gda/current/grib2/orig archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/${check_date_time}/MODELVAR/modelvar${check_date_time}000.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    --data-type=cma_gfs_gda/current/grib2/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 0h
  assert_output ${expected_result}
}

# grib2 orig with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_gfs_gda/current/grib2/orig runtime with config" {
  expected_result="/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GDA_POST_DATA/${check_date_time}/data/output/grib2_orig/modelvar${check_date_time}000.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime \
    "${config}" \
    --data-type=cma_gfs_gda/current/grib2/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 0h
  assert_output ${expected_result}
}

@test "test cma_gfs_gda/grib2/orig archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/GFS_GDA/Prod-grib/${check_date_time}/MODELVAR/modelvar${check_date_time}000.grb2"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=archive \
    "${config}" \
    --data-type=cma_gfs_gda/current/grib2/modelvar \
    --start-time "${check_date_time}" \
    --forecast-time 0h
  assert_output ${expected_result}
}



