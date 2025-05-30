#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)

hour=00
check_date_time=${check_date}${hour}

setup() {
  load '../../../tool/bats-support/load'
  load '../../../tool/bats-assert/load'
}

# modelvar

@test "test cma_meso_3km/current/bin/modelvar.cold runtime" {
  expected_result="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
        --data-type=cma_meso_3km/current/bin/modelvar.cold \
        --start-time "${check_date_time}" \
        --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result_today="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${today_check_date}${hour}00300"
  if [ -f "${expected_result_today}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
        --location-level=runtime \
        --data-type=cma_meso_3km/current/bin/modelvar.cold \
        --start-time "${today_check_date}${hour}" \
        --forecast-time 3h
    assert_output ${expected_result_today}
    return
  fi

  skip "data is not available: ${expected_result} and ${expected_result_today}"
}

@test "test cma_meso_3km/current/bin/modelvar.cold runtime/archive" {
  expected_result="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/${check_date_time}/modelvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run "${NWPC_DATA_CLIENT_PROGRAM}" local \
    --location-level=runtime/archive \
    --data-type=cma_meso_3km/current/bin/modelvar.cold \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

@test "test cma_meso_3km/current/bin/modelvar.cold archive" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/${check_date_time}/modelvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
    --location-level=archive \
    --data-type=cma_meso_3km/current/bin/modelvar.cold \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

# modelvar with config

config="--config-dir=${NWPC_DATA_CLIENT_CONFIG_DIR}/local"

@test "test cma_meso_3km/bin/modelvar.cold runtime with config" {
  expected_result="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${check_date_time}00300"
  if [ -f "${expected_result}" ]; then
    run ${NWPC_DATA_CLIENT_PROGRAM} local \
      --location-level=runtime \
      "${config}" \
      --data-type=cma_meso_3km/current/bin/modelvar.cold \
      --start-time "${check_date_time}" \
      --forecast-time 3h
    assert_output ${expected_result}
    return
  fi

  today_check_date=$(date +%Y%m%d)
  expected_result_today="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/cold/${hour}/fcst/grapes_model/run/modelvar${today_check_date}${hour}00300"
  if [ -f "${expected_result_today}" ]; then
    run "${NWPC_DATA_CLIENT_PROGRAM}" local \
      --location-level=runtime \
      "${config}" \
      --data-type=cma_meso_3km/current/bin/modelvar.cold \
      --start-time "${today_check_date}${hour}" \
      --forecast-time 3h
    assert_output ${expected_result_today}
    return
  fi

  skip "data is not available: ${expected_result} and ${expected_result_today}"
}

@test "test cma_meso_3km/bin/modelvar.cold runtime/archive with config" {
  expected_result="/g2/op_meso/OPER/WORKDIR/cma_meso_3km/DATABAK/cold/${check_date_time}/modelvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
    --location-level=runtime/archive \
    "${config}" \
    --data-type=cma_meso_3km/current/bin/modelvar.cold \
    --start-time "${check_date_time}" \
    --forecast-time 3h
  assert_output ${expected_result}
}

@test "test cma_meso_3km/bin/modelvar.cold archive with config" {
  expected_result="/g3/COMMONDATA/OPER/CEMC/MESO_3KM/Fcst-main/${check_date_time}/modelvar${check_date_time}00300"
  if [ ! -f "${expected_result}" ]; then
    skip "data is not available: ${expected_result}"
  fi
  run ${NWPC_DATA_CLIENT_PROGRAM} local \
        --location-level=archive \
        "${config}" \
        --data-type=cma_meso_3km/current/bin/modelvar.cold \
        --start-time "${check_date_time}" \
        --forecast-time 3h
  assert_output ${expected_result}
}
