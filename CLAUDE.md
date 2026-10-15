# CLAUDE.md

Guidance for Claude Code when working in this repository.

## What this is

`nwpc-data-client` — a Go CLI for operational NWPC/CEMC/CMA weather models. It resolves model output file paths by rendering YAML config templates with start time and forecast time, and can poll for files and run commands once they become stable. Builds a single binary `nwpc_data_client` into `bin/`.

## Commands

```bash
make            # build into ./bin (version via scripts/version.sh, injected with ldflags into internal/buildinfo)
make test       # unit tests
make generate   # regenerate internal/embedded/config.autogen.go from conf/ — required after any conf/ edit, commit the result
make test-integration   # integration tests (build tag `integration`, need real CMA HPC data, skip when absent)
goreleaser build --snapshot --clean --single-target   # fast local release-style build
```

`go vet` and `gofmt` are enforced in CI (`.github/workflows/ci.yml`); releases use GoReleaser on `v*` tags.

## Gotchas

- `--data-config-dir` is a **full replacement** of the embedded configs, not an overlay — a missing file is an error, never a fallback. `config export` exists to materialize the embedded baseline for merging.
- Config format versions: v3 is current (whole-file `{{...}}` templates). Loaders reject detectable older versions at path-resolution time; `config lint`/`config migrate` accept all versions.
- Data configs reserve `_`-prefixed top-level fields: `_config_version` (format version) and `_arrival_order` (sequential/unordered arrival, gates batched polling).
- `check local` has two polling modes: `per_time` (default, one goroutine per forecast time) and `batched` (`--check-mode`, single readdir-based loop) with sequential arrival gating (`_arrival_order` in the data config, overridable via `--arrival-order`).

## Architecture pointers

- `cmd/` is Cobra wiring only; implementation lives in `internal/` (`dataconfig`, `finder`, `checker`, `embedded`, `buildinfo`).
- Resolution flow: `dataconfig.LoadConfigContent` → `ParseConfigContent` (template execution) → `finder.FindLocalFile`.
- Full user/developer documentation is in `docs/` (e.g. `docs/reference/data-config-format.md`, `docs/guide/checking-files.md`, `docs/reference/checker-config.md`) — prefer reading/pointing there over duplicating details here.
