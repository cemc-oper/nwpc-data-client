# NWPC Data Client

A data finder CLI tool for operational systems in CEMC/CMA.

## Features

- Find file paths for operational system data in CMA HPC.

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

`nwpc_data_client` has several sub-commands.

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

## Test

Run `make test` to run all tests. 

> NOTE:
>
> Tests in `tests/bats` should run in CMA-HPC2023.
>

Tests in `tests/bats` use the following projects embedded as git submodules:

- [bats-core](https://github.com/bats-core/bats-core)
- [bats-support](https://github.com/bats-core/bats-support)
- [bats-assert](https://github.com/bats-core/bats-assert)

## License

Copyright &copy; 2019-2025 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).