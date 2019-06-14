# NWPC Data Clinet

A cli client for data files in NWPC.

## Features

-   Find operation system data in HPC-PI.

## Installing

Download the latest release and build source code. 

Use `Makefile` to build project on Linux.

```bash
make
```

`nwpc_data_client` command will be installed in `bin` directory.

## Getting Started

### local

`nwpc_data_client local` command finds local files on HPC-PI using config files.

```bash
nwpc_data_client local --config-dir=config_dir --data-type some/data/type \
    start_time forecast_time
```

Use `--config-dir` to set config file directory.

`--data-type` is some relative path under config directory. Such as

-   `gda_grapes_gfs/grib2/modelvar`
-   `gmf_graeps_gfs/bin/modelvar`

`start_time` is `YYYYMMDDHH` and `forecast_time` is `FFF`.

For example, use the command below to find GMF GRAPES GFS GRIB2 data of 24 forecast hour in start hour 00 on 2018/09/03.

```text
$nwpc_data_client local --data-type=gmf_grapes_gfs/grib2/orig 2018090300 24
/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_V2.2_POST/togrib2/output_togrib2/2018090300/gmf.gra.2018090300024.grb2
```

To list all data types available in some configure directory, run following command

```bash
nwpc_data_client local --config-dir=some/config/dir --show-types
```

Results may be like:

```text
gda_grapes_gfs/bin/modelvar
gda_grapes_gfs/bin/postvar
gda_grapes_gfs/grib2/modelvar
gda_grapes_gfs/grib2/orig
gmf_grapes_gfs/bin/modelvar
gmf_grapes_gfs/bin/postvar
gmf_grapes_gfs/grib2/modelvar
gmf_grapes_gfs/grib2/ne
gmf_grapes_gfs/grib2/orig
```