# cdcloud-io Makefile for GO.
# @: only show the output of the command, not the command itself

# build variables
BIN_DIR := bin
MODULE_NAME := cliDreamStream
URL_PATH := 

.DEFAULT_GOAL := all

# .PHONY as targets do not represent files.
.PHONY: all initapi build test test-with-cover generate-mocks clean run deps mod prod-build asm lint

# Print header information
all:
	@echo "**********************************************************"
	@echo "**          cdcloud-io GO build tool                    **"
	@echo "**********************************************************"

# Ensure BIN_DIR exists
$(BIN_DIR):
	@mkdir -p $@

# Lint the code
lint:
	@golangci-lint run --enable-all

# Build the cli application
build-cli: | $(BIN_DIR)
	@go build -v -o ${BIN_DIR}/cli/cli-dreamstream ./cmd/cli

# Build the qagent service
build-qagent: | $(BIN_DIR)
	@go build -v -o ${BIN_DIR}/qagent/qagent ./cmd/qagent

# Run the application
run-cli:
	@go run cmd/cli/*.go

# Run the qagent
run-qagent:
	@go run cmd/qagent/*.go

# Run tests
test:
	@go test -v $(shell go list ./... | grep -v /test/)

# Run tests with coverage
test-with-cover:
	@go test -v -coverprofile=cover.out $(shell go list ./... | grep -v /test/)
	@go tool cover -html=cover.out -o cover.html

# Generate mocks
generate-mocks:
	@mockery --all --with-expecter --keeptree

# Clean the build artifacts
clean:
	@go clean
	@rm -rf ${BIN_DIR}/*
	@rm -rf vendor

# Get dependencies
deps:
	@go get ./...

# Manage Go modules
mod: 
	@go mod tidy || (echo "Tidying go.mod failed with exit code $$?"; exit 1)
	@go mod download || (echo "Downloading modules failed with exit code $$?"; exit 1)
	@go mod vendor || (echo "Vendoring failed with exit code $$?"; exit 1)

# Production build
prod-build: mod | $(BIN_DIR)
	@go build -mod=vendor -ldflags="-s -w" -o ${BIN_DIR}/${MODULE_NAME}/cli-dreamstream ./cmd/${MODULE_NAME}/main.go || (echo "Build failed with exit code $$?"; exit 1)

# Generate assembly code
asm:
	@go tool compile -S cmd/${MODULE_NAME}/main.go > ${MODULE_NAME}.asm