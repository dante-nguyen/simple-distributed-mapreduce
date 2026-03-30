PROTO_FILES ?= ./rpc/v1/*.proto

.PHONY: proto-gen
protoc:
	protoc \
		--go_out=. --go_opt=paths=source_relative \
		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
		$(PROTO_FILES)
