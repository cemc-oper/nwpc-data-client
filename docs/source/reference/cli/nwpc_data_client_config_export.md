# nwpc_data_client config export

Export the embedded data configs into a directory tree.

## Synopsis

Export the embedded data configs into a directory tree.

The exported tree mirrors the conf/ directory of the nwpc-data-client
repository: each embedded config "local/<system>/<version>/<format>/<name>"
is written to "<output-dir>/local/<system>/<version>/<format>/<name>.yaml",
byte-identical to the source file in conf/.

This is the transitional answer to the "--data-config-dir replaces the
embedded configs entirely" limitation: a deployment pipeline can merge an
exported baseline with an incremental config repository before distribution,
instead of copying the embedded configs by hand.

Existing files are not overwritten unless --force is given.

```
nwpc_data_client config export <output-dir> [flags]
```

## Examples

```
  # Export the embedded configs into ./exported_conf.
  nwpc_data_client config export ./exported_conf
```

## Options

```
      --force   Overwrite existing files in the output directory.
  -h, --help    help for export
```

## SEE ALSO

* [nwpc_data_client config](nwpc_data_client_config.md)	 - Validate and maintain data config files.

