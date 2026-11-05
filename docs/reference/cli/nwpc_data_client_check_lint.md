# nwpc_data_client check lint

Statically validate checker config files.

## Synopsis

Statically validate checker runtime YAML files.

Each argument may be a checker file or a directory. Files supplied directly
are validated regardless of extension; directories are searched recursively
for *.yaml files. All traversal and validation errors are reported together.

Lint applies the normal default values and validates strict YAML fields,
durations, forecast times, data source combinations, and command templates.
An empty forecast_times list is valid because check local may read forecast
times from stdin at runtime.

This command does not require a start time, load data configs, read stdin,
access data directories, execute commands, or create run reports.

```
nwpc_data_client check lint <checker-file-or-dir> [more paths...] [flags]
```

## Examples

```
  # Lint one checker config file.
  nwpc_data_client check lint ./checker.yaml

  # Recursively lint multiple checker directories.
  nwpc_data_client check lint ./checker/system_a ./checker/system_b

  # Combine files and directories in one invocation.
  nwpc_data_client check lint ./extra.yaml ./checker
```

## Options

```
  -h, --help   help for lint
```

## SEE ALSO

* [nwpc_data_client check](nwpc_data_client_check.md)	 - Wait for operational data and act when it is ready.

