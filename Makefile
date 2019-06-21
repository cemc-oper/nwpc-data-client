all:
	go build \
	    -gcflags "all=-N -l" \
        -o bin/nwpc_data_client \
		main.go