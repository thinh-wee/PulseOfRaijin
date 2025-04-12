# Makefile for Application
.DEFAULT_GOAL := all
.PHONY: all build clean test test-local
SHELL = /usr/bin/env bash

# Variables
APP_NAME = PULSE-OF-RAIJIN
SRC_DIR = $(shell pwd)
BIN_DIR = $(SRC_DIR)/bin
SRC_FILES = $(wildcard cmd/Main/*.go)
BIN_FILE = $(BIN_DIR)/$(APP_NAME)

REPO_URL = $(shell git remote get-url origin) 
BUILD_USER = $(shell whoami)@$(shell hostname)
BUILD_DATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
BUILD_VERSION = $(shell git describe --tags --always)
BUILD_COMMIT = $(shell git rev-parse HEAD)
BUILD_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
BUILD_TAG = $(shell git tag --points-at HEAD)

ifeq ($(REPO_URL),)
	REPO_URL = "https://github.com/thinh-wee/pulse-of-raijin"
endif

ifeq ($(BUILD_TAG),)
	BUILD_VERSION = Branch: $(BUILD_BRANCH) - Commit: $(BUILD_COMMIT)
else
	BUILD_VERSION = $(BUILD_TAG)
endif

# All target
all: build
	@echo "Creating application directory..."
	@mkdir -p \
		Common/Packages/com.cyber.Model/ \
		Common/Packages/com.cyber.Feature/ \
		Common/Packages/com.cyber.Example/ || { echo "Failed to create directories"; exit 1; }
	@tree --version > /dev/null 2>&1 || { \
			echo "Run 'tree $(shell pwd)/Common' error: command not found."; \
			echo "Please install tree command (command: 'sudo apt-get update && sudo apt install -y tree') to see the directory structure."; \
			exit 1; \
		}
	@tree ./Common
	@echo "Run application..." && exec $(BIN_FILE)

# Build target
build: $(SRC_FILES)
	@echo "Building application..."
	@mkdir -p $(BIN_DIR)
	@export CGO_ENABLED=0;\
		go build -ldflags "-extldflags \"-static\" -s -w -X app.BuildUser=$(BUILD_USER)" \
			-o $(BIN_FILE) -trimpath $(SRC_FILES) || { echo "Failed to build application"; exit 1; }
	@echo "Build application successfully"

# Clean target
clean:
	@echo "Cleaning up..."
	@rm -rf $(BIN_DIR)


test:
	@echo "Running tests..."
	@go test -v ./... || { echo "Failed to run tests"; exit 1; }
	@echo "Tests completed successfully"

test-local:
	@echo "Running tests..."
	@go test -v local/pulseOfFraijin_test.go local/pulseOfFraijin.go -run TestMakePulseOfRaijinEncrypt
	@echo "Tests completed successfully"