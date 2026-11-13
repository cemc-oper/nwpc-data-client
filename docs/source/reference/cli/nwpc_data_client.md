# nwpc_data_client

Find data files for operational systems in CEMC.

## Synopsis

nwpc_data_client is a data file finder for CEMC operational systems.

It searches for data files across multiple directories and polls a series of
files until they are generated, running commands once they are found.

Common subcommands:
  local         Find a local file.
  check local   Wait for several files and run a command(s) when each file is found.
  version       Print version information.

Use "nwpc_data_client [command] --help" for detailed help on a subcommand.

```
nwpc_data_client [flags]
```

## Options

```
  -h, --help   help for nwpc_data_client
```

## SEE ALSO

* [nwpc_data_client check](nwpc_data_client_check.md)	 - Wait for operational data and act when it is ready.
* [nwpc_data_client config](nwpc_data_client_config.md)	 - Validate and maintain data config files.
* [nwpc_data_client local](nwpc_data_client_local.md)	 - Find the local file path for model output data.
* [nwpc_data_client version](nwpc_data_client_version.md)	 - Print version information.

