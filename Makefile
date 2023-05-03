#!/usr/bin/make -f

# Project variables
NAME?=elys
BINARY?=$(NAME)d
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags --match 'v*' --abbrev=8 | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')
GOFLAGS:=""
GOTAGS:=ledger

BUILDDIR ?= $(CURDIR)/build
DOCKER := $(shell which docker)
DOCKER_BUF := $(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace bufbuild/buf:1.7.0
DOCKERNET_HOME=./dockernet
DOCKERNET_COMPOSE_FILE=$(DOCKERNET_HOME)/docker-compose.yml

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)
BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(GOTAGS)'
BUILD_FOLDER = ./build

## install: Install elysd binary in $GOBIN
install: check-version go.sum
	@echo Installing Elysd binary...
	@GOFLAGS=$(GOFLAGS) go install $(BUILD_FLAGS) ./cmd/$(BINARY)
	@elysd version

## build: Build the binary
build: check-version go.sum
	@echo Building Elysd binary...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@GOFLAGS=$(GOFLAGS) go build $(BUILD_FLAGS) -o $(BUILD_FOLDER) ./cmd/$(BINARY)

## build-all: Build binaries for all platforms
build-all:
	@echo Building Elysd binaries for all platforms...
	@-mkdir -p $(BUILD_FOLDER) 2> /dev/null
	@GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_FOLDER)/$(BINARY)-linux-amd64 ./cmd/$(BINARY)
	@GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILD_FOLDER)/$(BINARY)-linux-arm64 ./cmd/$(BINARY)
	@GOFLAGS=$(GOFLAGS) GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILD_FOLDER)/$(BINARY)-darwin-amd64 ./cmd/$(BINARY)

## do-checksum: Generate checksums for all binaries
do-checksum:
	@cd build && sha256sum $(BINARY)-linux-amd64 $(BINARY)-linux-arm64 $(BINARY)-darwin-amd64 > $(BINARY)_checksum

## build-with-checksum: Build binaries for all platforms and generate checksums
build-with-checksum: build-all do-checksum

.PHONY: install build build-all do-checksum build-with-checksum

## mocks: Generate mocks
mocks:
	@echo Generating mocks
	@go install github.com/vektra/mockery/v2
	@go generate ./...

## test-unit: Run unit tests
test-unit:
	@echo Running unit tests...
	@GOFLAGS=$(GOFLAGS) go test -race -failfast -v ./...

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@rm -rf $(BUILD_FOLDER) 2> /dev/null
	@go clean ./...

.PHONY: mocks test-unit clean

## go-mod-cache: Retrieve the go modules and store them in the local cache
go-mod-cache: go.sum
	@echo "--> Retrieve the go modules and store them in the local cache."
	@go mod download

## go.sum: Ensure dependencies have not been modified
go.sum: go.mod
	@echo "--> Make sure that the dependencies haven't been altered."
	@go mod verify

# Add check to make sure we are using the proper Go version before proceeding with anything
check-version:
	@if ! go version | grep -q "go1.19"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.19 is required for compiling elysd. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.19 and retry. Thank you!"; \
		exit 1; \
	fi

.PHONY: go-mod-cache go.sum check-version

help: Makefile
	@echo
	@echo " Choose a command run in "$(BINARY)", or just run 'make' for install"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

.PHONY: help

.DEFAULT_GOAL := install

build-docker:
	@bash $(DOCKERNET_HOME)/build.sh -${build} ${BUILDDIR}

start-docker: build-docker
	@bash $(DOCKERNET_HOME)/start_network.sh

clean-docker:
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) stop
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) down
	rm -rf $(DOCKERNET_HOME)/state
	docker image prune -a

stop-docker:
	@bash $(DOCKERNET_HOME)/pkill.sh
	docker-compose -f $(DOCKERNET_COMPOSE_FILE) down --remove-orphans
