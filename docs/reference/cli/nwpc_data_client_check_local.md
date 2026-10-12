# nwpc_data_client check local

Wait for local files, then run a command for each one.

## Synopsis

Wait for a list of data files with different forecast times and
run command(s) when each file is found.

For each forecast time, this command repeatedly looks up the data file (using
the same config resolution as "nwpc_data_client local"). Once a file exists and
its size stops changing between checks, it is considered stable and the
configured command(s) are run. Each forecast time is checked concurrently after
a staggered --delay-time.

Config resolution mirrors the local command:

  * --data-type            Use an embedded config named "local/<data-type>".
  * --data-type combined with --data-config-dir
		Load a user defined config from "<data-config-dir>/<data-type>.yaml".
  * --data-config-file     Load a user defined config directly from "<data-config-file>".

The list of forecast times comes from --checker-config (the forecast_times
field) or, when absent, from whitespace-separated tokens on stdin.

CLI flags override values from --checker-config. The command(s) to run support
Go templates with the usual config variables (Year, Month, Day, Hour,
ForecastHour, ...) plus FilePath. Use --execute-command for a single command;
multiple commands (execute_commands) can be supplied via --checker-config.

Args:
    start_time   Start time, YYYYMMDDHH, such as 2026062400.

```
nwpc_data_client check local <start_time> [flags]
```

## Examples

```
  # Check via a checker config file (recommended).
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml

  # Provide forecast times on stdin and run a command for each ready file.
  echo "0h 3h 6h" | nwpc_data_client check local 2026062400 \
      --data-type=cma_gfs_gmf/grib2/orig \
      --execute-command="echo ready: {{.FilePath}}"

  # Poll faster and limit the number of attempts per forecast time.
  nwpc_data_client check local 2026062400 --checker-config=./checker.yaml \
      --check-interval=10s --max-check-count=360
```

## Options

```
      --data-config-dir string             Directory holding data config files. Optional; uses embedded configs when empty.
      --data-config-file string            Path to a single data config file. If set, --data-config-dir and --data-type are ignored.
      --checker-config string              Path to the checker runtime config YAML. CLI flags override values in this file.
      --data-type string                   Data type used to locate the config (embedded as local/<data-type>, or <dir>/<data-type>.yaml).
      --location-level string              Location levels to search, split by ',', such as 'runtime,archive'. Empty searches all.
      --max-check-count int                Maximum number of existence/size checks per forecast time. (default 2880)
      --check-interval string              Interval between checks, a duration such as 5s, 30s, or 1m. (default "5s")
      --execute-command string             Go template for a command to run once a file becomes stable. Mutually exclusive with execute_commands.
      --delay-time string                  Stagger delay between consecutive forecast times (duration such as 10s). (default "10s")
      --check-interval-jitter float        Jitter ratio applied to the polling interval: each round sleeps interval×(1±j). 0 disables jitter. Range [0, 0.5]. (default 0.1)
      --max-concurrent-stable-checks int   Maximum number of files in size-stability observation at once (batched check mode only). (default 4)
      --check-mode string                  Polling mode: per_time (one goroutine per forecast time) or batched (single readdir-based loop). (default "per_time")
      --debug                              Enable debug logging.
  -h, --help                               help for local
```

## SEE ALSO

* [nwpc_data_client check](nwpc_data_client_check.md)	 - Wait for operational data and act when it is ready.

