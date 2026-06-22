# NWPC Data Client

A data finder CLI tool for operational systems in CEMC/CMA.

## Features

- Find file paths for operational system data in CMA HPC.
- Watch data files and execute some command when data file is ready.

## Installing

Download the latest release or build from source.

### Local build (Linux / macOS / WSL / Git Bash)

Use `make` to build all binaries (`nwpc_data_client`, `nwpc_data_checker`) into the `bin/` directory:

```bash
make
```

Build a single binary:

```bash
make nwpc_data_client
make nwpc_data_checker
```

### Cross-compile from Windows (PowerShell)

If GoReleaser is not installed, use the fallback script:

```powershell
# Linux amd64
.\build_cross.ps1 -Arch amd64

# Linux arm64
.\build_cross.ps1 -Arch arm64
```

### Release builds (GoReleaser)

GoReleaser is the preferred tool for cross-platform and release builds:

```bash
# Snapshot build (no tag required)
goreleaser build --snapshot --clean

# Release build (requires a git tag)
goreleaser release --clean
```

## Getting Started

`nwpc_data_client` has the following sub-command.

### local

`nwpc_data_client local` command finds local files on CMA HPC.

```bash
nwpc_data_client local \
    --data-type some/data/type \
    --start-time start_time \
    --forecast-time forecast_time
```

`data-type` is some relative path under config directory. Such as

- `cma_gfs_gmf/current/grib2/modelvar`
- `cma_gfs_gmf/current/bin/modelvar`

`start-time` is `YYYYMMDDHH` and `forecast-time` is `FFFh` or `FFFh00m`.

For example, use the command below to find CMA-GFS GMF GRIB2 data of 24 forecast hour in start hour 00 on 2025/05/29.

```
$nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2025052900 \
    --forecast-time 000h
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2025052900/data/output/grib2_orig/gmf.gra.2025052900000.grb2
```

To list all data types available in some configured directory, run the following command:

```bash
nwpc_data_client local --show-types
```

Results may be like:

```text
cma_geps/current/grib2/orig
cma_gfs_gda/current/bin/modelvar
cma_gfs_gda/current/bin/postvar
cma_gfs_gda/current/grib2/modelvar
cma_gfs_gda/current/grib2/orig
cma_gfs_gmf/current/bin/modelvar
cma_gfs_gmf/current/bin/postvar
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
cma_reps/current/grib2/orig
cma_tym/current/bin/modelvar
cma_tym/current/bin/modelvar_ctl
cma_tym/current/bin/postvar
cma_tym/current/bin/postvar_ctl
cma_tym/current/grib2/orig
```

Use `--config-dir` to set a custom config file directory.

### checker

`nwpc_data_checker local` polls local data paths for a list of forecast times and optionally runs commands once each file becomes stable. It is typically used in operational workflows to wait for model output and then trigger downstream processing.

Forecast times can be provided in two ways:

1. From `stdin`, one per line or whitespace-separated token, in the form `FFFh` or `FFFhMMm` (for example `000h`, `024h`, `003h10m`). This is the original behavior and remains supported.
2. From a YAML runtime config file passed with `--checker-config`. This is the recommended approach for complex command pipelines.

```bash
nwpc_data_checker local YYYYMMDDHH \
    --data-type some/data/type \
    --location-level runtime,archive
```

The checker will, for each forecast time:

1. Wait a short staggered delay (`--delay-time`, default `0s`).
2. Repeatedly check the resolved path until the file exists.
3. Continue polling until the file size stops changing, which indicates the file is no longer being written.
4. Run `--execute-command` if one is configured.

If the file is not found or does not stabilize before `--max-check-count` is reached, the checker exits with an error.

A complete example: wait for CMA-GFS GMF GRIB2 original output and print the path once it is stable.

```bash
printf "000h\n024h\n" | nwpc_data_checker local 2025052900 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --location-level runtime \
    --execute-command 'echo found {.FilePath}'
```

The `--execute-command` template supports the same variables as configuration files (see [Configuration file](#configuration-file)) plus `.FilePath`, which is set to the resolved file path. Use `{` and `}` as delimiters, for example `{.FilePath}` or `{.Year}{.Month}{.Day}{.Hour}`.

### Checker runtime config

For workflows with long command strings or multiple downstream steps, use a YAML runtime config file via `--checker-config` (or `-c`). All flags except `start_time` can be moved into this file, and any flag given on the command line overrides the file value.

```yaml
# checker-config.yaml
data_type: cma_gfs_gmf/current/grib2/orig
location_levels: runtime
max_check_count: 2880
check_interval: 5s
delay_time: 0s
debug: false

forecast_times:
  - 000h
  - 024h
  - 048h

# Use either execute_command (single command) or execute_commands (list), not both.
execute_commands:
  - 'echo "found {.FilePath}"'
  - '/app/postprocess.sh {.FilePath}'
```

Run it with:

```bash
nwpc_data_checker local 2025052900 --checker-config checker-config.yaml
```

If `forecast_times` is omitted from the config file, forecast times are still read from `stdin` as in the original behavior.

The `execute_command` field accepts a single Go template string. The `execute_commands` field accepts a list of template strings, which are executed in order after each file stabilizes. If any command fails, the checker exits with an error.

> **Note:** 
> Some YAML parsers (including the one used by `nwpc_data_checker`) have trouble with double-quoted scalars that contain a colon followed by a space, 
> such as `"echo file path: {.FilePath}"`. To avoid parse errors, wrap command strings in single quotes, for example `'echo "file path: {.FilePath}"'`.

Common flags:

| Flag | Description | Default |
|------|-------------|---------|
| `--checker-config` | Path to a YAML runtime config file. CLI flags override file values. | - |
| `--data-type` | Data type used to locate the config file in the config dir. | (required if `--data-config-file` is not set) |
| `--data-config-dir` | Custom config directory, same as `nwpc_data_client local`. | embedded configs |
| `--data-config-file` | Path to a single config file; if set, `--data-config-dir` and `--data-type` are ignored. | - |
| `--location-level` | Comma-separated location levels, such as `runtime,archive`. | all levels |
| `--max-check-count` | Maximum number of check rounds for one forecast time. | `2880` |
| `--check-interval` | Polling interval, as a Go duration (`30s`, `1m`, ...). | `5s` |
| `--delay-time` | Stagger delay between forecast times, as a Go duration. | `0s` |
| `--execute-command` | Go template command to run when the file is stable. | - |
| `--debug` | Enable debug logging. | `false` |

## Configuration file

Data types are described by YAML configuration files. When no `--config-dir`/`--data-config-dir` is given, `nwpc_data_client` and `nwpc_data_checker` use the configs embedded in the binary at build time (see `make generate`). A custom directory can be used during development or for site-specific overrides.

### File format

A configuration file is a single YAML document with the following fields:

```yaml
default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2
file_names: []

paths:
  - type: local
    level: runtime
    path: /some/dir/{.Year}{.Month}{.Day}{.Hour}/

  - type: local
    level: archive
    path: /another/dir/{.Year}{.Month}{.Day}{.Hour}/
```

| Field | Description |
|-------|-------------|
| `default` | Value returned when no file is found. Defaults to `NOTFOUND` in most configs. |
| `file_name` | Single filename template. Either `file_name` or `file_names` must be set. |
| `file_names` | List of filename templates. They are tried in order after the optional `file_name`. Useful when the same data may be stored under several names. |
| `paths` | List of directory entries to search. Each entry is checked in order; the first existing file wins. |

Each entry in `paths` has:

| Field | Description |
|-------|-------------|
| `type` | Path type. For `nwpc_data_client local` and `nwpc_data_checker` this is `local`. |
| `level` | Location level, such as `runtime`, `archive`, or `all`. Use `--location-level` to restrict the search to specific levels. An empty level or `all` matches any filter. |
| `path` | Directory template. The resolved filename is joined to this directory. |

### Template variables

The config file is executed as a Go template with `{` and `}` delimiters. The following variables are available:

| Variable | Type | Example | Description |
|----------|------|---------|-------------|
| `.StartTime` | `time.Time` | - | Parsed start time. |
| `.ForecastTime` | `time.Duration` | - | Parsed forecast time. |
| `.Year` | string | `2025` | Four-digit year of start time. |
| `.Month` | string | `05` | Zero-padded month of start time. |
| `.Day` | string | `29` | Zero-padded day of start time. |
| `.Hour` | string | `00` | Zero-padded hour of start time. |
| `.ForecastHour` | string | `024` | Zero-padded forecast hour, three digits. |
| `.ForecastMinute` | string | `00` | Zero-padded forecast minute, two digits. |
| `.Member` | string | - | Ensemble member, when applicable. |

Template helper functions are also available for use in more complex expressions:

```yaml
file_name: gmf.gra.{generateStartTime .StartTime}{getForecastHour .ForecastTime}.grb2
```

| Function | Description |
|----------|-------------|
| `generateStartTime` | Format start time as `YYYYMMDDHH`. |
| `getYear` | Extract the year. |
| `getMonth` | Extract the month. |
| `getDay` | Extract the day. |
| `getHour` | Extract the hour. |
| `generateForecastTime` | Format forecast time as `FFFhMMm`. |
| `getForecastHour` | Extract the forecast hour. |
| `getForecastMinute` | Extract the forecast minute. |

### How resolution works

1. The YAML content is loaded from the embedded config, custom config directory, or a single config file.
2. `nwpc_data_client`/`nwpc_data_checker` parse the YAML as a Go template with the start time, forecast time, and optional member.
3. The resolved directory and filename templates are joined.
4. Each resulting path is checked in order; the first existing file is returned, or `default` if none are found.

### Example

```yaml
# conf/local/cma_gfs_gmf/current/grib2/orig.yaml
default: NOTFOUND

file_name: gmf.gra.{.Year}{.Month}{.Day}{.Hour}{.ForecastHour}.grb2

paths:
  - type: local
    level: runtime
    path: /g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/{.Year}{.Month}{.Day}{.Hour}/data/output/grib2_orig/

  - type: local
    level: archive
    path: /g3/COMMONDATA/OPER/CEMC/GFS_GMF/Prod-grib/{.Year}{.Month}{.Day}{.Hour}/ORIG
```

For start time `2025052900` and forecast time `024h`, the first path resolves to:

```text
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2025052900/data/output/grib2_orig/gmf.gra.2025052900024.grb2
```

If the file does not exist in `runtime`, the tool falls back to the `archive` path.

## Test

Run unit tests:

```bash
make test
```

This runs `go test ./...`.

Integration tests live under `tests/integration/` and use the `integration` build tag
(`//go:build integration`). They invoke `bin/nwpc_data_client` via `os/exec` and check
for real operational files on CMA HPC paths, skipping when data is unavailable.
`make test-integration` builds the binary automatically and runs them:

```bash
make test-integration
```

Run a single integration test:

```bash
go test -tags=integration -run TestCMA_GFS_GMF ./tests/integration/...
```

The integration test helper honors `NWPC_DATA_CLIENT_PROGRAM` and
`NWPC_DATA_CLIENT_CONFIG_DIR` when set, otherwise defaults to `./bin/nwpc_data_client`
and `./conf`.

## License

Copyright &copy; 2019-2025 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).