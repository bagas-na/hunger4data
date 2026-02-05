PROTO_DIR = ./proto/

.PHONY: proto
proto:
	@echo "Generating protobuf files..."
	protoc --proto_path=$(PROTO_DIR) \
	       --go_out=$(PROTO_DIR) --go_opt=paths=source_relative \
	       --go-grpc_out=$(PROTO_DIR) --go-grpc_opt=paths=source_relative \
	       $(shell find $(PROTO_DIR) -name "*.proto")

.PHONY: tests
tests:
	@echo "Testing REST microservice..."
	go test -count=5 ./apps/api-gateway/...
	@echo "Testing Auth gRPC microservice..."
	go test -count=5 ./apps/authenticator/...
	@echo "Testing Notification microservice..."
	go test -count=5 ./apps/notification/...
	@echo "Testing Payment microservice..."
	go test -count=5 ./apps/payment/...
	@echo "Testing Subscription microservice..."
	go test -count=5 ./apps/subscription/...
