PROTO_DIR = ../../proto/notification

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	       $(shell find $(PROTO_DIR) -name "*.proto")
