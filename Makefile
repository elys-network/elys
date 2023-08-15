#!/usr/bin/make -f

# Project variables
NAME?=elys
BINARY?=$(NAME)d
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags --match 'v*' --abbrev=8 | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')
GOFLAGS:=""
GOTAGS:=ledger
SHELL := /bin/bash # Use bash syntax

# currently installed Go version
GO_MAJOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f1)
GO_MINOR_VERSION = $(shell go version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)

# minimum supported Go version
GO_MINIMUM_MAJOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f1)
GO_MINIMUM_MINOR_VERSION = $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f2 | cut -d'.' -f2)

RED=\033[0;31m
GREEN=\033[0;32m
LGREEN=\033[1;32m
NOCOLOR=\033[0m
GO_CURR_VERSION=$(shell echo -e "Current Go version: $(LGREEN)$(GO_MAJOR_VERSION).$(GREEN)$(GO_MINOR_VERSION)$(NOCOLOR)")
GO_VERSION_ERR_MSG=$(shell echo -e '$(RED)❌ ERROR$(NOCOLOR): Go version $(LGREEN)$(GO_MINIMUM_MAJOR_VERSION).$(GREEN)$(GO_MINIMUM_MINOR_VERSION)$(NOCOLOR)+ is required')

GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)

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
	@echo "Generating mocks"
	@if ! command -v mockery &> /dev/null; then \
		echo "mockery binary not found, installing..."; \
		go install github.com/vektra/mockery/v2@latest; \
	fi
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
	@echo '$(GO_CURR_VERSION)'
	@if [[ $(GO_MAJOR_VERSION) -eq $(GO_MINIMUM_MAJOR_VERSION) && $(GO_MINOR_VERSION) -ge $(GO_MINIMUM_MINOR_VERSION) ]]; then \
		exit 0; \
	elif [[ $(GO_MAJOR_VERSION) -lt $(GO_MINIMUM_MAJOR_VERSION) ]]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
		exit 1; \
	elif [[ $(GO_MINOR_VERSION) -lt $(GO_MINIMUM_MINOR_VERSION) ]]; then \
		echo '$(GO_VERSION_ERR_MSG)'; \
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

.PHONY: build-docker start-docker clean-docker stop-docker

###############################################################################
###                                Release                                  ###
###############################################################################

GORELEASER_IMAGE := ghcr.io/goreleaser/goreleaser-cross:v$(GO_VERSION)
COSMWASM_VERSION := $(shell go list -m github.com/CosmWasm/wasmvm | sed 's/.* //')

ifdef GITHUB_TOKEN
release:
	docker run \
		--rm \
		-e GITHUB_TOKEN=$(GITHUB_TOKEN) \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/elysd \
		-w /go/src/elysd \
		$(GORELEASER_IMAGE) \
		release \
		--clean
else
release:
	@echo "Error: GITHUB_TOKEN is not defined. Please define it before running 'make release'."
endif

release-dry-run:
	docker run \
		--rm \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/elysd \
		-w /go/src/elysd \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--skip-publish

release-snapshot:
	docker run \
		--rm \
		-e COSMWASM_VERSION=$(COSMWASM_VERSION) \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/elysd \
		-w /go/src/elysd \
		$(GORELEASER_IMAGE) \
		release \
		--clean \
		--snapshot \
		--skip-validate \
		--skip-publish

.PHONY: release release-dry-run release-snapshot