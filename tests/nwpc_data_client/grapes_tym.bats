#!/usr/bin/env bats

check_date=$(date -d "-1 day" +%Y%m%d)
check_date_time=${check_date}00

nwpc_data_client_bin=bin/nwpc_data_client

@test "test grapes_tym/bin/postvar" {
  result=$(${nwpc_data_client_bin} local --data-type=grapes_tym/bin/postvar "${check_date_time}" 3h)
  [ "$result" = "/g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat/postvar${check_date_time}00300" ]
}

@test "test grapes_tym/bin/modelvar" {
  result=$(${nwpc_data_client_bin} local --data-type=grapes_tym/bin/modelvar "${check_date_time}" 3h)
  [ "$result" = "/g2/nwp_qu/NWP_RMFS_DATA/grapes_tym/grapes_d01/dat/modelvar${check_date_time}00300" ]
}
