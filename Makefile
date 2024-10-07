# Makefile for go-compiler services

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOINSTALL=$(GOCMD) install

# Binary names
REQUEST_SERVICE_BINARY=request-service
PYTHON_WORKER_BINARY=python-worker

# Protoc parameters
PROTOC=protoc
PROTO_DIR=request-service/proto
GO_OUT_DIR=request-service/proto

# Main build target
all: deps proto build

# Build all services
build: build-request-service build-python-worker

# Build request-service
build-request-service:
	cd request-service && $(GOBUILD) -o $(REQUEST_SERVICE_BINARY) -v

# Build python-worker
build-python-worker:
	cd workers/python-worker && $(GOBUILD) -o $(PYTHON_WORKER_BINARY) -v

# Clean build files
clean:
	$(GOCLEAN)
	rm -f request-service/$(REQUEST_SERVICE_BINARY)
	rm -f workers/python-worker/$(PYTHON_WORKER_BINARY)

# Run tests
test:
	$(GOTEST) -v ./...

# Get dependencies
deps:
	$(GOGET) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOGET) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GOMOD) tidy
	$(GOINSTALL) google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GOINSTALL) google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Generate protobuf files
proto:
	PATH="$(shell go env GOPATH)/bin:$$PATH" $(PROTOC) --proto_path=$(PROTO_DIR) --go_out=$(GO_OUT_DIR) --go_opt=paths=source_relative \
		--go-grpc_out=$(GO_OUT_DIR) --go-grpc_opt=paths=source_relative \
		$(PROTO_DIR)/*.proto

# Run the request-service
run-request-service: build-request-service
	./request-service/$(REQUEST_SERVICE_BINARY)

# Run the python-worker
run-python-worker: build-python-worker
	./workers/python-worker/$(PYTHON_WORKER_BINARY)

# Help command to display available targets
help:
	@echo "Available targets:"
	@echo "  all                   - Generate proto files, get dependencies, and build all services"
	@echo "  build                 - Build all services"
	@echo "  build-request-service - Build the request-service"
	@echo "  build-python-worker   - Build the python-worker"
	@echo "  clean                 - Clean build files"
	@echo "  test                  - Run tests"
	@echo "  deps                  - Get dependencies"
	@echo "  proto                 - Generate protobuf files"
	@echo "  run-request-service   - Build and run the request-service"
	@echo "  run-python-worker     - Build and run the python-worker"

.PHONY: all build build-request-service build-python-worker clean test deps proto run-request-service run-python-worker help