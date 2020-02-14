#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}

hour_4dv=21
check_date_4dvar=$(date -d "-2 day" +%Y%m%d)
check_data_time_4dvar=${check_date_4dvar}${hour_4dv}

# grib2 orig
@test "test grapes_gfs_gmf/grib2/orig runtime" {
  expected_result="/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/${check_date_time}/gmf.gra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_gfs_gmf/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available: ${expected_result}"
}

@test "test grapes_gfs_gmf/grib2/orig archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/${check_data_time_4dvar}/ORIG/gmf.gra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        --data-type=grapes_gfs_gmf/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available: ${expected_result}"
}

# grib2 orig with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_gfs_gmf/grib2/orig runtime with config" {
  expected_result="/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_POST/togrib2/output_togrib2/${check_date_time}/gmf.gra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_gfs_gmf/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available: ${expected_result}"
}

@test "test grapes_gfs_gmf/grib2/orig archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GMF/Prod-grib/${check_data_time_4dvar}/ORIG/gmf.gra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        "${config}" \
        --data-type=grapes_gfs_gmf/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available: ${expected_result}"
}



