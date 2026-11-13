# NWPC Data Client

**English** | [简体中文](README.zh-CN.md)

[![Maturity-Incubating](https://img.shields.io/badge/Maturity-Incubating-56B4E9)](https://github.com/cemc-oper/nwpc-data-client)
[![GitHub Release](https://img.shields.io/github/v/release/cemc-oper/nwpc-data-client)](https://github.com/cemc-oper/nwpc-data-client/releases)
[![CI](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml/badge.svg)](https://github.com/cemc-oper/nwpc-data-client/actions/workflows/ci.yml)
[![Documentation Status](https://readthedocs.org/projects/nwpc-data-client/badge/?version=latest)](https://nwpc-data-client.readthedocs.io/)

A command-line tool for locating data files from CEMC/CMA operational numerical weather prediction systems on CMA supercomputing platforms.

Documentation (Chinese): <https://nwpc-data-client.readthedocs.io/>

## Features

- Locate data files from operational numerical weather prediction systems across multiple directories
- Monitor batches of data files and execute commands once the files are ready

## Quick Start

Download a binary for your platform from [Releases](https://github.com/cemc-oper/nwpc-data-client/releases), or build from source with `make` (the output is written to `bin/`).

Find the base GRIB2 file for the 24-hour forecast from the CMA-GFS 2026-06-22 12 cycle:

```bash
nwpc_data_client local \
    --data-type=cma_gfs_gmf/current/grib2/orig \
    --start-time 2026062212 \
    --forecast-time 24h
```

Monitor a batch of data files and execute a command after each file has been generated and its size has stabilized:

```bash
echo "000h 024h" |
    nwpc_data_client check local 2026062312 \
        --data-type=cma_gfs_gmf/current/grib2/orig \
        --location-level runtime \
        --execute-command 'cp {{.FilePath}} /path/to/archive/'
```

## Build and Development

```bash
make            # Build into ./bin
make test       # Run unit tests
make generate   # Regenerate embedded configs after modifying conf/
```

Cross-compilation and local release builds use GoReleaser snapshot mode:

```bash
goreleaser release --snapshot --clean               # Build for all platforms; output goes to dist/
goreleaser build --snapshot --clean --single-target # Build only for the current platform (faster)
```

See the Development section of the documentation site for versioning, cross-compilation, integration testing, and the release process.

## License

Copyright &copy; 2019-2026 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).
