#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
check_date_time=${check_date}00

@test "test grapes_tym/grib2/orig runtime" {
  expected_result="/g2/nwp_pd/NWP_GRAPES_TYM_POST_DATA/${check_date_time}/rundir/output/orig_grib2/rmf.tcgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_tym/grib2/orig \
        "${check_date_time}" 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_tym/grib2/orig archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Prod-grib/${check_date_time}/rmf.tcgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=archive \
        --data-type=grapes_tym/grib2/orig \
        "${check_date_time}" 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}

# postvar with config
config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_tym/grib2/orig runtime with config" {
  expected_result="/g2/nwp_pd/NWP_GRAPES_TYM_POST_DATA/${check_date_time}/rundir/output/orig_grib2/rmf.tcgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        "${config}" \
        --data-type=grapes_tym/grib2/orig \
        "${check_date_time}" 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}


@test "test grapes_tym/grib2/orig archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_TYM/Prod-grib/${check_date_time}/rmf.tcgra.${check_date_time}003.grb2"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        "${config}" \
        --location-level=archive \
        --data-type=grapes_tym/grib2/orig \
        "${check_date_time}" 3h
    [ "${output}" = "${expected_result}" ]
    return
  fi

  skip "data is not available"
}