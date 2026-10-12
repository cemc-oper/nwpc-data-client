## nwpc_data_client config lint

Validate data config files.

### Synopsis

Lint data config files of any config format version (v1/v2/v3).

Each file is validated in two stages:

  1. template validation with the real text/template parser, using the
     delimiters, functions and variables matching the detected config format
     version (the "_config_version" field wins; heuristic sniffing is the
     fallback). Unknown functions fail at parse, unknown variables at execute.
  2. structural validation of the rendered result against the data config
     structure (paths items require type/level/path, and so on).

Pass a config directory to lint every *.yaml file under it recursively, or
use --data-config-file to lint a single file.

```
nwpc_data_client config lint [config-dir] [flags]
```

### Examples

```
  # Lint all configs under conf/local recursively.
  nwpc_data_client config lint conf/local

  # Lint a single config file.
  nwpc_data_client config lint --data-config-file=./my_config.yaml
```

### Options

```
      --data-config-file string   Path to a single data config file. If set, the config directory argument is ignored.
  -h, --help                      help for lint
```

### SEE ALSO

* [nwpc_data_client config](nwpc_data_client_config.md)	 - Validate and maintain data config files.

