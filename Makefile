BINARY_DIR := bin
SERVER_BINARY := $(BINARY_DIR)/server
CLIENT_BINARY := $(BINARY_DIR)/client

.PHONY: help run-server run-client build build-server build-client test fmt vet clean

help:
	@echo "Available commands:"
	@echo "  make run-server    Run the Tic-Tac-Toe server"
	@echo "  make run-client    Run the Tic-Tac-Toe client"
	@echo "  make build         Build both server and client"
	@echo "  make test          Run all tests"
	@echo "  make fmt           Format Go source files"
	@echo "  make vet           Run go vet"
	@echo "  make clean         Remove generated binaries"

run-server:
	go run ./cmd/server

run-client:
	go run ./cmd/client

build: build-server build-client

build-server:
	@mkdir -p $(BINARY_DIR)
	go build -o $(SERVER_BINARY) ./cmd/server

build-client:
	@mkdir -p $(BINARY_DIR)
	go build -o $(CLIENT_BINARY) ./cmd/client

test:
	go test ./...

fmt:
	go fmt ./...

vet:
	go vet ./...

clean:
	rm -rf $(BINARY_DIR)