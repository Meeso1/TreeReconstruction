.PHONY: build clean test run

# Variables
BINARY_NAME=treereconstruction
BINARY_DIR=bin
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT=$(shell git rev-parse HEAD 2>/dev/null || echo "none")
BUILD_DATE=$(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS=-ldflags="-X 'treereconstruction/cmd.Version=$(VERSION)' -X 'treereconstruction/cmd.Commit=$(COMMIT)' -X 'treereconstruction/cmd.BuildDate=$(BUILD_DATE)'"

# Tasks
build:
	mkdir -p $(BINARY_DIR)
	go build $(LDFLAGS) -o $(BINARY_DIR)/$(BINARY_NAME)

clean:
	rm -rf $(BINARY_DIR)
	go clean

test:
	go test ./... -v

run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

install: build
	go install $(LDFLAGS)

# Example of running with arguments
run-process: build
	./$(BINARY_DIR)/$(BINARY_NAME) process -i sample.txt -o output.txt 