## nwpc_data_client local

Find the local file path for model output data.

### Synopsis

Find a local data file through multiple directories and print the file path when found.

The path is computed by applying --start-time and --forecast-time to a YAML
config template. The config is chosen in one of three ways:

  * --data-type            Use an embedded config named "local/<data-type>".
  * --data-type combined with --data-config-dir  Load a user defined config from "<data-config-dir>/<data-type>.yaml".
  * --data-config-file     Load a user defined config directly from "<data-config-file>".

If a matching file exists it is printed; otherwise the config's "default" value
(usually NOTFOUND) is printed. Use --location-level to restrict which path
entries are searched (for example "runtime" or "runtime,archive").

Run with --show-types to list the available data types (embedded, or from
--data-config-dir) and exit.

```
nwpc_data_client local [flags]
```

### Examples

```
  # Use an embedded config for CMA-GFS grib2 data.
  nwpc_data_client local --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h

  # Use a config directory and only search runtime locations.
  nwpc_data_client local --data-config-dir=./config --data-type=cma_meso_1km/bin/modelvar \
      --start-time=2026062400 --forecast-time=24h --location-level=runtime

  # Use a single config file directly.
  nwpc_data_client local --data-config-file=./my_config.yaml \
      --start-time=2026062400 --forecast-time=3h

  # List available embedded data types.
  nwpc_data_client local --show-types
```

### Options

```
      --data-config-dir string    Directory holding data config files. Optional; uses embedded configs when empty.
      --data-type string          Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).
      --data-config-file string   Path to a single data config file. If set, --data-config-dir and --data-type are ignored.
      --location-level string     Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.
      --start-time string         Start time, YYYYMMDDHH, such as 2026062400. (required)
      --forecast-time string      Forecast time, such as 0h, 24h, 120h. (required)
      --member string             Ensemble member, MMM, such as 000, 014.
      --show-types                List the supported data types and exit. Honors --data-config-dir when set.
      --debug                     Enable debug logging.
  -h, --help                      help for local
```

### SEE ALSO

* [nwpc_data_client](nwpc_data_client.md)	 - Find data files for operational systems in CEMC.

