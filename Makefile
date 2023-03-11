# makefile for project, with help from:
# https://gist.github.com/serinth/16391e360692f6a000e5a10382d1148c

# Windows check
ifneq ($(shell echo $$OS),Windows_NT)
    $(error This project only builds on Windows at the moment)
endif

# Executable directory
BINARY_DIR=bin/

# Base Go build flags
GOFLAGS=-ldflags '-linkmode external'

# C compiler
CC=gcc

# Base C build flags
CFLAGS=-Wall -Wextra -pedantic -std=c17

# Where to look for C files
SRC_DIR=./pkg ./internal

# Find all .c files recursively in the source directory
SRC_FILES=$(shell find $(SRC_DIR) -name '*.c')

# Convert .c files to corresponding .o files in the same directory
OBJ_FILES=$(patsubst %.c,%.o,$(SRC_FILES))

.PHONY: help clean clean-all fmt vet update-dependencies test test-bench test-cover test-all build build-debug build-release build-obj build-obj-debug build-obj-release

default: help

# Show this help
help:
	@echo 'usage: make [target] ...'
	@echo 'TODO: add more documentation here'

# Show some variables for debugging
env:
	@echo $(SRC_FILES)
	@echo $(OBJ_FILES)

# Run Go clean and remove C library and object files, also remove binary files
clean:
	go clean
	rm -f $(OBJ_FILES)
	rm -f $(shell find $(SRC_DIR) -name 'lib*.a')
	rm -fr bin/*

# More go clean parameters
clean-all: clean
	go clean -i ./...
	go clean -cache

# Format Go and C files
fmt:
	go fmt ./...
	find . -name '*.c' -o -name '*.h' | xargs clang-format -style=Google -i

# Run go vet on source files
vet:
	go vet ./...

# Update Go dependencies
update-dependencies:
	go mod tidy

# Run short tests
test: build-obj
	go test -v ./... -short

# Run benchmark tests
test-bench: build-obj
	go test -bench ./...

# Generate test coverage and generate html report
test-cover: build-obj
	rm -fr coverage
	mkdir coverage
	go list -f '{{if gt (len .TestGoFiles) 0}}"go test -covermode count -coverprofile {{.Name}}.coverprofile -coverpkg ./... {{.ImportPath}}"{{end}}' ./... | xargs -I {} bash -c {}
	echo "mode: count" > coverage/cover.out
	grep -h -v "^mode:" *.coverprofile >> "coverage/cover.out"
	rm *.coverprofile
	go tool cover -html=coverage/cover.out -o=coverage/cover.html

test-all: test test-bench test-cover

# Build with base flags
build: clean build-obj
	go build $(GOFLAGS) -o $(BINARY_DIR) ./...

# Build with debugging flags
build-debug: GOFLAGS += -gcflags="-N -l" -ldflags="-X 'main.Version=dev'" -tags="debug" -v -race
build-debug: clean build-obj-debug
	go build $(GOFLAGS) -o $(BINARY_DIR) ./...

# Build with release flags
build-release: GOFLAGS += -ldflags="-s -w" -trimpath -mod=readonly -buildmode=pie -tags netgo -installsuffix=static
build-release: clean build-obj-release
	go build $(GOFLAGS) -o $(BINARY_DIR) ./...

# Build C object and library files
build-obj: $(OBJ_FILES)

# Build C object and library files with debug flags
build-obj-debug: CFLAGS += -g -DDEBUG
build-obj-debug: $(OBJ_FILES)

# Build C object and library files with release flags
build-obj-release: CFLAGS += -DNDEBUG -O3
build-obj-release: $(OBJ_FILES)

%.o: %.c
	$(CC) $(CFLAGS) -c $^ -o $@
	ar rcs $(join $(dir $@), $(join lib, $(patsubst %.o,%.a,$(notdir $@)))) $@