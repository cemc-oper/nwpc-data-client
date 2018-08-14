# NWPC Data Clinet

A cli client for data files in NWPC.

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
nwpc_find_data_path --config-dir=config_dir --data-type some/data/type \
    start_time forecast_time
```

Use `--config-dir` to set config file direcotry.

`data-type` is some relative path under config directory. Such as

- `gda_grapes_gfs/grib2_modelvar`
- `gmf_graeps_gfs/modelvar`
