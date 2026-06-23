# NWPC Data Client

A data finder CLI tool for operational systems in CEMC/CMA.

## Features

- Find file paths for operational system data in CMA HPC.
- Watch data files and execute commands when a data file is ready.

## Installing

Download the latest release or build from source.

### Local build (Linux / macOS / WSL / Git Bash)

Use `make` to build the binary (`nwpc_data_client`) into the `bin/` directory:

```bash
make
```

Build a single binary:

```bash
make nwpc_data_client
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

## Quick start

Find a local data path using embedded configs:

```bash
./bin/nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

Output:

```text
/g2/op_post/OPER/WORKDIR/NWP_CMA_GFS_GMF_POST_DATA/2026062212/data/output/grib2_orig/gmf.gra.2026062212024.grb2
```

## Data Client

`nwpc_data_client local` find local files on CMA HPC using YAML configuration files.

### Usage

```bash
nwpc_data_client local \
    --data-type <data/type> \
    --start-time <YYYYMMDDHH> \
    --forecast-time <FFFh> \
    [--member <MMM>]
```

- `--data-type` is a relative path under the config directory, such as `cma_gfs_gmf/current/grib2/modelvar`.
- `--start-time` is `YYYYMMDDHH`.
- `--forecast-time` is `FFFh` or `FFFhMMm`.
- `--member` is the ensemble member, such as `000` or `014`.

All flags for `nwpc_data_client local`:

| Flag | Description | Default |
|------|-------------|---------|
| `--config-dir` | Config directory. | embedded configs |
| `--data-config-file` | Path to a single config file; if set, `--config-dir` and `--data-type` are ignored. | - |
| `--data-type` | Data type used to locate the config file in the config dir. | required* |
| `--location-level` | Comma-separated location levels, such as `runtime,archive`. | all levels |
| `--start-time` | Start time, `YYYYMMDDHH`. | required |
| `--forecast-time` | Forecast time, `FFFh` or `FFFhMMm`. | required |
| `--member` | Ensemble member, `MMM`. | - |
| `--show-types` | List available data types and exit. | `false` |
| `--debug` | Enable debug logging. | `false` |

\* `--data-type` is required only when `--data-config-file` is not set.

### Examples

#### Find with embedded configs

Find CMA-GFS GMF GRIB2 original data for start time `2026/06/22 12` and forecast hour `024`:

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

#### Find with custom configs

Use `--config-dir` to load YAML configs from a custom directory instead of the embedded configs:

```bash
nwpc_data_client local \
  --config-dir=/g1/u/op_post/OPER/UTIL/nwpc_data_client/config/local_hpc2023 \
  --data-type=cma_gfs_gmf/v4.2.3/oper/bin/modelvar \
  --location-level=all \
  --start-time=2026062300 \
  --forecast-time=024h
```

#### Find with a single config file

Use a single data config file directly (this ignores `--config-dir` and `--data-type`):

```bash
nwpc_data_client local \
    --data-config-file /path/to/config.yaml \
    --start-time 2025052900 \
    --forecast-time 000h
```

#### Listing available data types

To list all data types embedded in binary program, run:

```bash
nwpc_data_client local --show-types
```

To list all data types available in the configured directory, run:

```bash
nwpc_data_client local --show-types \
  --config-dir /g1/u/op_post/OPER/UTIL/nwpc_data_client/config/local_hpc2023 
```

## Data Checker

`nwpc_data_client check local` polls local data paths for a list of forecast times and optionally runs commands once each file becomes stable. 
It is typically used in operational workflows to wait for model output and then trigger downstream processing.

`start_time` is passed as a positional argument in the form `YYYYMMDDHH`.

### Polling behavior

For each forecast time the checker will:

1. Wait a short staggered delay (`--delay-time`, default `10s`).
2. Repeatedly check the resolved path until the file exists.
3. Continue polling until the file size stops changing, which indicates the file is no longer being written.
4. Run `--execute-command` if one is configured.

If the file is not found or does not stabilize before `--max-check-count` is reached, the checker exits with an error.

### Forecast times via stdin

When no `--checker-config` is provided, forecast times are read from `stdin`, one per line or as whitespace-separated tokens, in the form `FFFh` or `FFFhMMm` (for example `000h`, `024h`, `003h10m`).

A complete example: wait for CMA-GFS GMF GRIB2 original output and print the path once it is stable.

```bash
printf "000h\n024h\n" | nwpc_data_client check local 2025052900 \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --location-level runtime \
    --execute-command 'echo found {.FilePath}'
```

The `--execute-command` template supports the same variables as configuration files (see [Template variables](#template-variables)) plus `.FilePath`, which is set to the resolved file path. Use `{` and `}` as delimiters, for example `{.FilePath}` or `{.Year}{.Month}{.Day}{.Hour}`.

### Common flags

| Flag | Description | Default |
|------|-------------|---------|
| `--checker-config` | Path to a YAML runtime config file. CLI flags override file values. | - |
| `--data-type` | Data type used to locate the config file in the config dir. | required* |
| `--data-config-dir` | Custom config directory, same as `nwpc_data_client local`. | embedded configs |
| `--data-config-file` | Path to a single config file; if set, `--data-config-dir` and `--data-type` are ignored. | - |
| `--location-level` | Comma-separated location levels, such as `runtime,archive`. | all levels |
| `--max-check-count` | Maximum number of check rounds for one forecast time. | `2880` |
| `--check-interval` | Polling interval, as a Go duration (`30s`, `1m`, ...). | `5s` |
| `--delay-time` | Stagger delay between forecast times, as a Go duration. | `10s` |
| `--execute-command` | Go template command to run when the file is stable. | - |
| `--debug` | Enable debug logging. | `false` |

\* `--data-type` is required only when `--data-config-file` is not set.

### Checker runtime config

For workflows with long command strings or multiple downstream steps, use a YAML runtime config file via `--checker-config` (or `-c`). 
All flags except `start_time` can be moved into this file, and any flag given on the command line overrides the file value.

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
nwpc_data_client check local 2025052900 --checker-config checker-config.yaml
```

If `forecast_times` is omitted from the config file, forecast times are still read from `stdin` as in the original behavior.

The `execute_command` field accepts a single Go template string. 
The `execute_commands` field accepts a list of template strings, which are executed in order after each file stabilizes. 
If any command fails, the checker exits with an error.

> **NOTICE**
> 
> For any field that exists both in the checker runtime config and as a CLI flag, **CLI flags always override the YAML file**.

> **YAML quoting note**
>
> Some YAML parsers (including the one used by `nwpc_data_client check`) have trouble with double-quoted scalars that contain a colon followed by a space, 
> such as `"echo file path: {.FilePath}"`. 
> To avoid parse errors, wrap command strings in single quotes, for example `'echo "file path: {.FilePath}"'`.

## Configuration

Data types are described by YAML configuration files. When no `--config-dir` / `--data-config-dir` is given, `nwpc_data_client` uses the configs embedded in the binary at build time (see [Regenerating embedded configs](#regenerating-embedded-configs)). A custom directory can be used during development or for site-specific overrides.

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
| `type` | Path type. For `nwpc_data_client local` and `nwpc_data_client check local` this is `local`. |
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

### Template helper functions

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
2. `nwpc_data_client` parses the YAML as a Go template with the start time, forecast time, and optional member.
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

## Development

### Version information

The binary provides a `version` subcommand that prints version, git commit, and build time:

```bash
nwpc_data_client version
```

### Running tests

Run unit tests:

```bash
make test
```

This runs `go test ./...`.

Integration tests live under `tests/integration/` and use the `integration` build tag (`//go:build integration`). They invoke `bin/nwpc_data_client` via `os/exec` and check for real operational files on CMA HPC paths, skipping when data is unavailable. `make test-integration` builds the binary automatically and runs them:

```bash
make test-integration
```

Run a single integration test:

```bash
go test -tags=integration -run TestCMA_GFS_GMF ./tests/integration/...
```

The integration test helper honors `NWPC_DATA_CLIENT_PROGRAM` and `NWPC_DATA_CLIENT_CONFIG_DIR` when set, otherwise defaults to `./bin/nwpc_data_client` and `./conf`.

### Regenerating embedded configs

Configs under `conf/` can be embedded into the binary as `common/config/config.autogen.go`:

```bash
make generate
```

This runs `go generate` in `common/config/generate`. Remember to commit `common/config/config.autogen.go` after adding or editing YAML configs under `conf/`.

## License

Copyright &copy; 2019-2025 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).
