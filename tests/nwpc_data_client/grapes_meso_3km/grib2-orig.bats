#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)

hour=00
check_date_time=${check_date}${hour}

# grib2 orig
@test "test grapes_meso_3km/grib2/orig runtime" {
  expected_result="/g2/nwp_pd/NWP_GRAPES_MESO_3KM_POST_DATA/${check_date_time}/togrib2/output/rmf.hgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_meso_3km/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/grib2/orig archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Prod-grib/${check_date_time}/ORIG/rmf.hgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        --data-type=grapes_meso_3km/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

# grib2 orig with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_meso_3km/grib2/orig runtime with config" {
  expected_result="/g2/nwp_pd/NWP_GRAPES_MESO_3KM_POST_DATA/${check_date_time}/togrib2/output/rmf.hgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_meso_3km/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/grib2/orig archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Prod-grib/${check_date_time}/ORIG/rmf.hgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        "${config}" \
        --data-type=grapes_meso_3km/grib2/orig \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}



