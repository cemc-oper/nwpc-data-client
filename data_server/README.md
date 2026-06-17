# NWPC Data Server

`nwpc_data_server` is a gRPC server that exposes the same data-path resolution as `nwpc_data_client`, plus file info and chunked download.

## Build

Build `nwpc_data_server` from the project root:

```bash
cd repo/nwpc-data-client
make nwpc_data_server
```

Or build all binaries at once:

```bash
make
```

## Server

Start the data server on port 33483:

```bash
./bin/nwpc_data_server serve --address ":33483" --config-dir=some/config/dir
```

## Client

Use `nwpc_data_client service` to connect to `nwpc_data_server`.

Use `--action` to specify the remote action:

- `findDataPath`
- `getDataFileInfo`
- `downloadDataFile`

### findDataPath

Get the resolved data path from the server:

```bash
nwpc_data_client service --address=data-service-address \
    --action findDataPath \
    --data-type=some/data/type \
    --start-time=start_time \
    --forecast-time=forecast_time
```

Output is the same as `nwpc_data_client hpc`:

```text
local
/g2/op_post/OPER/.../gmf.gra.2025052900000.grb2
```

### getDataFileInfo

Get data file information on the server:

```bash
nwpc_data_client service --address=data-service-address \
    --action getDataFileInfo \
    --data-type=some/data/type \
    --start-time=start_time \
    --forecast-time=forecast_time
```

When the data file is found, two lines are printed: file path and file size.

```text
/g2/nwp/OPER_ARCH_TEST/nwp/GRAPES_GFS/GMF_GRAPES_GFS/Prod-grib/2019061621/ORIG/gmf.gra.2019061700000.grb2
303997137
```

When the data file is not found or an error occurs, the error message is printed to stderr:

```text
check file error: stat NOTFOUND: no such file or directory
```

### downloadDataFile

Download a data file from the server:

```bash
nwpc_data_client service --address=data-service-address \
    --action downloadDataFile \
    --output-dir=output/dir \
    --data-type=some/data/type \
    --start-time=start_time \
    --forecast-time=forecast_time
```

The file is saved to `--output-dir` with its original remote file name.

## Develop

### Prerequisites

Install the modern Go protobuf and gRPC plugins:

```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.31.0
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.3.0
```

### Regenerate gRPC code

After changing `api/data_service/data_service.proto`, regenerate the Go code:

```bash
cd repo/nwpc-data-client
make generate
```

Or run `go generate` directly in the proto package:

```bash
cd api/data_service
go generate
```

This generates two files:

- `api/data_service/data_service.pb.go` — messages
- `api/data_service/data_service_grpc.pb.go` — gRPC client/server interfaces

## License

Copyright © 2019-2025 developers at cemc-oper.

`nwpc-data-client` is licensed under [The MIT License](https://opensource.org/licenses/MIT).
