# nwpc data service


## Build

Prepare build environment.

Download and build gPRC library.

Get gRPC plugin for golang.

```cmd
go get -u google.golang.org/grpc
go get -u github.com/golang/protobuf/protoc-gen-go
set PATH=%PATH%;%GOPATH%/bin
```

Generate gRPC codes.

```cmd
protoc.exe ^
    -I data_service ^
    data_service/data_service.proto ^
    --go_out=plugins=grpc:.
```