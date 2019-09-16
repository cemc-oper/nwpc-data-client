# NWPC Data Client

A cli client for data files in NWPC.

## Features

- Find file path for operation system data in HPC-PI or external storage nodes from HPC-PI.
- Download data file through a data service.

## Installing

Download the latest release and build source code. 

Use `Makefile` to build project on Linux and 
`nwpc_data_client` command will be installed in `bin` directory.

## Getting Started

`nwpc_data_client` has several sub-commands.

### local

`nwpc_data_client local` command finds local files on HPC-PI.

```bash
nwpc_data_client local --data-type some/data/type \
    start_time forecast_time
```

`data-type` is some relative path under config directory. Such as

- `grapes_gfs_gda/grib2/modelvar`
- `gmf_graeps_gfs/bin/modelvar`

`start_time` is `YYYYMMDDHH` and `forecast_time` is `FFFh`.

For example, use the command below to find GMF GRAPES GFS GRIB2 data of 24 forecast hour in start hour 00 on 2018/09/03.

```text
$nwpc_data_client local --data-type=grapes_gfs_gmf/grib2/orig 2018090300 24h
/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_V2.2_POST/togrib2/output_togrib2/2018090300/gmf.gra.2018090300024.grb2
```

To list all data types available in some configure directory, run following command

```bash
nwpc_data_client local --show-types
```

Results may be like:

```text
grapes_gfs_gda/bin/modelvar
grapes_gfs_gda/bin/postvar
grapes_gfs_gda/grib2/modelvar
grapes_gfs_gda/grib2/orig
grapes_gfs_gmf/bin/modelvar
grapes_gfs_gmf/bin/postvar
grapes_gfs_gmf/grib2/modelvar
grapes_gfs_gmf/grib2/ne
grapes_gfs_gmf/grib2/orig
```

Use `--config-dir` to set a custom config file directory.

### hpc

`nwpc_data_client hpc` command find files on HPC-PI or external storage nodes from HPC-PI. 

Data files could be on HPC's local storage nodes (eg. /g2) or 
external storage nodes which are mount to special HPC login nodes.

```bash
nwpc_data_client hpc --config-dir=config_dir --data-type=some/data/type \
    start_time forecast_time
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
$nwpc_data_client hpc --data-type=grapes_gfs_gda/bin/modelvar 2019050200 0h
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

See [README.md](./data_service/README.md) under data_service.

## Test

Run `make test` to run all tests. [bats](https://github.com/neurodebian/bats) is required.

## License

Copyright &copy; 2019 Perilla Roc at nwpc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).