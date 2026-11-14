# nwpc_data_client fetch

Find one data file and copy or link it into a target directory.

## Synopsis

Find one data file and materialize it into a target directory.

The data config is resolved with the same start time, forecast time, member,
location levels, and candidate ordering as "nwpc_data_client local". The first
existing regular file is selected; missing candidates are skipped, while other
source errors stop the operation. Local candidates are read directly. Cmafs
candidates use the machine backend supplied by --backend-config and direct
SSH/rsync transport.

The source must be non-empty and have unchanged size and modification time
across two observations separated by --stable-interval. The destination name
is the source basename. Copy mode writes a unique temporary file in the target
directory, verifies it, and publishes it without replacing an existing file.
Symlink mode creates an absolute symbolic link to the selected local source.
Auto uses symlink for runtime/archive local sources and copy for storage local
sources and cmafs downloads. Cmafs downloads are verified against remote size
and mtime before and after transfer.

Existing destinations are never overwritten. An identical regular file or an
allowed matching symlink is reused; a conflicting destination is an error.
A per-target interprocess lock prevents concurrent duplicate materialization.

On success, only the final absolute path is written to stdout. Logs and errors
go to stderr. This command does not wait for files that have not appeared yet.

```
nwpc_data_client fetch [flags]
```

## Examples

```
  # Copy a CMA-GFS file into a work directory.
  nwpc_data_client fetch --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h --target-dir=/work/input

  # Link a local runtime file instead of copying it.
  nwpc_data_client fetch --data-type=cma_gfs_gmf/current/grib2/orig \
      --start-time=2026062400 --forecast-time=24h \
      --target-dir=/work/input --method=symlink

  # Use a single config file and automatic method selection.
  nwpc_data_client fetch --data-config-file=./my_config.yaml \
      --start-time=2026062400 --forecast-time=3h \
      --target-dir=./input --method=auto
```

## Options

```
      --data-config-dir string    Directory holding data config files. Optional; uses embedded configs when empty.
      --data-type string          Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).
      --data-config-file string   Path to a single data config file. If set, --data-config-dir and --data-type are ignored.
      --backend-config string     Path to machine backend YAML. Required only when a cmafs source is needed.
      --location-level string     Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.
      --storage-base string       Absolute local mount root for {{.StorageBase}} in data config templates.
      --start-time string         Start time, YYYYMMDDHH, such as 2026062400. (required)
      --forecast-time string      Forecast time, such as 0h, 24h, 120h. (required)
      --member string             Ensemble member, MMM, such as 000, 014.
      --target-dir string         Directory in which the fetched file is published. (required)
      --method string             Materialization method: copy, symlink, or auto. Auto copies storage-level sources. (default "copy")
      --stable-interval string    Interval between two source size/mtime stability checks, such as 1s or 500ms. (default "1s")
      --debug                     Enable debug logging.
  -h, --help                      help for fetch
```

## SEE ALSO

* [nwpc_data_client](nwpc_data_client.md)	 - Find data files for operational systems in CEMC.

