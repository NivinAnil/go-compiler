# Define paths
PROTO_DIR=./proto
OUT_DIR=./generated

# Define proto files
PROTO_FILES=$(wildcard $(PROTO_DIR)/*.proto)

# Default target
all: generate

# Generate Go code from .proto files
generate:
	protoc -I=$(PROTO_DIR) \
        --go_out=Mindex.proto=go-compiler/generated/requestpb:$(OUT_DIR) \
        --go-grpc_out=Mindex.proto=go-compiler/generated/requestpb:$(OUT_DIR) \
        $(PROTO_FILES)

# Clean generated files
clean:
	rm -rf $(OUT_DIR)

# Install required tools
install-tools:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Run the application
run:
	go run main.go

.PHONY: all generate clean install-tools run
