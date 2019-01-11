# NWPC Data Clinet

A cli client for data files in NWPC.

## Features

-   Find operation system data in HPC PI.

## Installing

Download the lastest release and install.

```bash
python setup.py install
```

`nwpc_find_data_path` command will be installed in python's binary path.

## Getting Started

### local

`nwpc_find_data_path local` command finds local files using config files.

```bash
nwpc_find_data_path local --config-dir=config_dir --data-type some/data/type \
    start_time forecast_time
```

Use `--config-dir` to set config file direcotry.

`data-type` is some relative path under config directory. Such as

-   `gda_grapes_gfs/grib2_modelvar`
-   `gmf_graeps_gfs/modelvar`

`start_time` is `YYYYMMDDHH` and `forecast_time` is `FFF`.

For example, use the command below to find GMF GRAPES GFS GRIB2 data of 24 forecast hour in start hour 00 on 2018/09/03.

```text
$nwpc_find_data_path local --data-type=gmf_grapes_gfs/grib2/orig 2018090300 24
/g2/nwp_pd/NWP_PST_DATA/GMF_GRAPES_GFS_V2.2_POST/togrib2/output_togrib2/2018090300/gmf.gra.2018090300024.grb2
```
