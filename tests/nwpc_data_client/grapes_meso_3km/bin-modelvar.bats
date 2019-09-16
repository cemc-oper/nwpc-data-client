#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
year=$(echo ${check_date} | cut -b1-4)
month=$(echo ${check_date} | cut -b5-6)
day=$(echo ${check_date} | cut -b7-8)

hour=00
check_date_time=${check_date}${hour}

# modelvar

@test "test grapes_meso_3km/bin/modelvar runtime" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${today_check_date}${hour}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${today_check_date}${hour}" 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/bin/modelvar runtime/archive" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/DATABAK/cold/${check_date_time}/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime/archive \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/bin/modelvar archive" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/${check_date_time}/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=archive \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

# modelvar with config

config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test grapes_meso_3km/bin/modelvar runtime with config" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${today_check_date}${hour}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        "${config}" \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${today_check_date}${hour}" 3h
    [ "x${output}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/bin/modelvar runtime/archive with config" {
  expected_result="/g2/nwp_qu/NWP_RMFS_DATA/grapes_meso_3km/DATABAK/cold/${check_date_time}/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=runtime/archive \
        "${config}" \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}

@test "test grapes_meso_3km/bin/modelvar archive with config" {
  expected_result="/g1/COMMONDATA/OPER/NWPC/GRAPES_MESO_3KM/Fcst-main/${check_date_time}/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    result=$(${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=archive \
        "${config}" \
        --data-type=grapes_meso_3km/bin/modelvar \
        "${check_date_time}" 3h)
    [ "x${result}" = "x${expected_result}" ]
    return
  fi

  skip "data is not available"
}
