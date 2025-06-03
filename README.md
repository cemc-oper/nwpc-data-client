# NWPC Data Client

A data finder CLI tool for operational systems in CEMC/CMA.

## Features

- Find file paths for operational system data in CMA HPC.
- :thumbsdown: Find file paths in external storage nodes from CMA HPC.
- :thumbsdown: Download data file through a data service.

## Installing

Download the latest release and build source code. 

Use `Makefile` to build the project on Linux and `nwpc_data_client` command will be installed in `bin` directory.

## Getting Started

`nwpc_data_client` has several sub-commands.

### local

`nwpc_data_client local` command finds local files on CMA HPC.

```bash
nwpc_data_client local --data-type some/data/type \
    --start-time start_time \
    --forecast-time forecast_time
```

`data-type` is some relative path under config directory. Such as

- `cma_gfs_gmf/current/grib2/modelvar`
- `cma_gfs_gmf/current/bin/modelvar`

`start-time` is `YYYYMMDDHH` and `forecast-time` is `FFFh`.

For example, use the command below to find CMA-GFS GMF GRIB2 data of 24 forecast hour in start hour 00 on 2018/09/03.

```
$nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2025052900 \
    --forecast-time 000h
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2025052900/data/output/grib2_orig/gmf.gra.2025052900000.grb2
```

To list all data types available in some configured directory, run the following command

```bash
nwpc_data_client local --show-types
```

Results may be like:

```text
cma_gfs_gda/current/bin/modelvar
cma_gfs_gda/current/bin/postvar
cma_gfs_gda/current/grib2/modelvar
cma_gfs_gda/current/grib2/orig
cma_gfs_gmf/current/bin/modelvar
cma_gfs_gmf/current/grib2/modelvar
cma_gfs_gmf/current/grib2/ne
cma_gfs_gmf/current/grib2/orig
cma_meso_1km/current/bin/modelvar.cold
cma_meso_1km/current/bin/modelvar
cma_meso_1km/current/bin/modelvar_ctl.cold
cma_meso_1km/current/bin/modelvar_ctl
cma_meso_1km/current/grib2/orig.cold
cma_meso_1km/current/grib2/orig
cma_meso_3km/current/bin/modelvar.cold
cma_meso_3km/current/bin/modelvar
cma_meso_3km/current/bin/modelvar_ctl.cold
cma_meso_3km/current/bin/modelvar_ctl
cma_meso_3km/current/bin/postvar.cold
cma_meso_3km/current/bin/postvar
cma_meso_3km/current/bin/postvar_ctl.cold
cma_meso_3km/current/bin/postvar_ctl
cma_meso_3km/current/grib2/orig.cold
cma_meso_3km/current/grib2/orig
cma_tym/current/bin/modelvar
cma_tym/current/bin/modelvar_ctl
cma_tym/current/bin/postvar
cma_tym/current/bin/postvar_ctl
cma_tym/current/grib2/orig
```

Use `--config-dir` to set a custom config file directory.

### hpc

> :warning: This command is not supported yet.

`nwpc_data_client hpc` command find files on HPC-PI or external storage nodes from HPC-PI. 

Data files could be on HPC's local storage nodes (eg. /g2) or 
external storage nodes which are mount to special HPC login nodes.

```bash
nwpc_data_client hpc --config-dir=config_dir --data-type=some/data/type \
    --start-time start_time \
    --forecast-time forecast_time
```

`hpc` support all options of `local`, and has more options to access external storage nodes using ssh protocol.

- `--storage-user`: user to login, default is environment variable `USER`
- `--storage-host`: host to login, default is `10.40.140.44:22`
- `--private-key`: private key file path, default is `$HOME/.ssh/id_rsa`
- `--host-key`: host key file path,  default is `$HOME/.ssh/known_hosts`

Currently only no password private key is supported, 
and user should test to access remote host manually before using this command.

`paths` field in `hpc`'s config file has two types: 

- `local` for local files
- `storage` for files on external storage.

For example, use the command below to find GDA GRAPES GFS modelvar data of 000 forecast hour in start hour 00 on 2019/05/20.

```text
$nwpc_data_client hpc --data-type=grapes_gfs_gda/bin/modelvar --start-time 2019050200 --forecast-time 0h
storage
/sstorage1/COMMONDATA/OPER/nwp/GRAPES_GFS/GDA_GRAPES_GFS/Fcst-9h/2019050121/modelvar2019050121_000
```

The command return two lines: 

1. first line is path type: `local` or `storage`.
2. second line is file path.

If no file is found, both lines will be value of `default` field in config file. 

`paths` section of a config file may like this:

```yaml
paths:
  - type: local
    level: archive
    path: /g2/nwp/OPER_ARCH_TEST/nwp/GRAPES_GFS/GDA_GRAPES_GFS/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
  - type: storage
    level: storage
    path: /sstorage1/COMMONDATA/OPER/nwp/GRAPES_GFS/GDA_GRAPES_GFS/Fcst-9h/{.Year4DV}{.Month4DV}{.Day4DV}{.Hour4DV}
```

## NWPC Data Service

> :warning: This command is not supported yet.

See [README.md](./data_service/README.md) under data_service.

## Test

Run `make test` to run all tests. [bats](https://github.com/neurodebian/bats) is required.

## License

Copyright &copy; 2019-2025 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).