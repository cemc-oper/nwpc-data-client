#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
hour=00
check_date_time=${check_date}${hour}

hour_4dv=21
check_date_4dvar=$(date -d "-2 day" +%Y%m%d)
check_data_time_4dvar=${check_date_4dvar}${hour_4dv}

# modelvar

@test "test grapes_gfs_gda/bin/modelvar runtime" {
  expected_result="/g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/${hour_4dv}/output/modelvar${check_data_time_4dvar}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/${hour_4dv}/output/modelvar${check_date}${hour_4dv}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_gfs_gda/bin/modelvar archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/${check_data_time_4dvar}/modelvar${check_data_time_4dvar}_003"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=archive \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

# modelvar with config

config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_gfs_gda/bin/modelvar runtime with config" {
  expected_result="/g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/${hour_4dv}/output/modelvar${check_data_time_4dvar}_003"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/nwp/GRAPES_GFS/MODEL/data/NWP_GDAS/${hour_4dv}/output/modelvar${check_date}${hour_4dv}_003"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_gfs_gda/bin/modelvar archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_GFS_GDA/Fcst-9h/${check_data_time_4dvar}/modelvar${check_data_time_4dvar}_003"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=archive \
        "${config}" \
        --data-type=grapes_gfs_gda/bin/modelvar \
        --start-time "${check_date_time}" \
        --forecast-time 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}
