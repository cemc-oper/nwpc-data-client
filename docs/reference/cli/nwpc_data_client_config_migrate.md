## nwpc_data_client config migrate

Migrate data config files to a newer format version.

### Synopsis

Migrate data config files to a newer config format version.

Supported migrations:

  v1 -> v2: rename template variables ({.Forecast} -> {.ForecastHour},
            the 4DV/1HR offset variables -> generateStartTime /
            generateForecastTime pipelines)
  v2 -> v3: replace "{...}" action delimiters with "{{...}}" (only actions
            starting with a variable, a known function or a template keyword
            are rewritten; literal braces are kept)
  v1 -> v3: both steps in sequence

The migrated content gets a "_config_version" field and is validated with
the config lint machinery before being written. A file whose source does not
pass lint (e.g. a v1 config using functions in paths[].path, where no
functions were registered) or whose migrated result does not pass lint is
reported and left untouched.

Files are overwritten in place unless --output-dir is given. Files already
at the target version are skipped.

```
nwpc_data_client config migrate [config-dir] [flags]
```

### Examples

```
  # Migrate a single config file to v3, in place.
  nwpc_data_client config migrate --data-config-file=./my_config.yaml

  # Migrate all configs under a directory to v3, writing to another tree.
  nwpc_data_client config migrate ./old_conf/local --output-dir=./new_conf/local

  # Migrate to v2 only.
  nwpc_data_client config migrate --data-config-file=./my_config.yaml --to-version=2
```

### Options

```
      --data-config-file string   Path to a single data config file. If set, the config directory argument is ignored.
      --to-version int            Target config format version (2 or 3). (default 3)
      --output-dir string         Write migrated files under this directory instead of overwriting the source files. The directory structure of the config directory is mirrored.
  -h, --help                      help for migrate
```

### SEE ALSO

* [nwpc_data_client config](nwpc_data_client_config.md)	 - Validate and maintain data config files.

