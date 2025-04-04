#!/usr/bin/make -f

# Project variables
NAME?=elys
BINARY?=$(NAME)d
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags --match 'v*' --abbrev=8 | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')
GOFLAGS:=""

# if rocksdb env variable is set, add the tag
ifdef ROCKSDB
	DBENGINE=rocksdb
else
	DBENGINE=pebbledb
endif

GOTAGS:=ledger,$(DBENGINE)

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

#GO_VERSION := $(shell cat go.mod | grep -E 'go [0-9].[0-9]+' | cut -d ' ' -f 2)
GO_VERSION := 1.23

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
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT) \
		  -X github.com/cosmos/cosmos-sdk/types.DBBackend=$(DBENGINE) \
		  -X github.com/cosmos/cosmos-sdk/version.BuildTags=netgo,ledger,muslc,osusergo,$(DBENGINE)
BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(GOTAGS)'


PROTO_VERSION=0.14.0
protoImageName=ghcr.io/cosmos/proto-builder:$(PROTO_VERSION)
protoImage=$(DOCKER) run --rm -v $(CURDIR):/workspace --workdir /workspace $(protoImageName)

.PHONY: proto

proto-all: proto-format proto-gen

proto:
	@echo "Generating Protobuf files"
	@sh ./scripts/protocgen.sh

proto-gen:
	@echo "Generating Protobuf files"
	@$(protoImage) sh ./scripts/protocgen.sh

proto-format:
	@$(protoImage) find ./ -name "*.proto" -exec clang-format -i {} \;

## install: Install elysd binary in $GOBIN
install: check-version go.sum
	@echo Installing Elysd binary...
	@GOFLAGS=$(GOFLAGS) go install $(BUILD_FLAGS) ./cmd/$(BINARY)
	@elysd version

$(BUILDDIR)/:
	@-mkdir -p $(BUILDDIR) 2> /dev/null

## build: Build the binary
build: check-version go.sum $(BUILDDIR)/
	@echo Building Elysd binary...
	@GOFLAGS=$(GOFLAGS) go build $(BUILD_FLAGS) -o $(BUILDDIR) ./cmd/$(BINARY)

## build-all: Build binaries for all platforms
build-all: $(BUILDDIR)/
	@echo Building Elysd binaries for all platforms...
	@GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILDDIR)/$(BINARY)-linux-amd64 ./cmd/$(BINARY)
	@GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o $(BUILDDIR)/$(BINARY)-linux-arm64 ./cmd/$(BINARY)
	@GOFLAGS=$(GOFLAGS) GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o $(BUILDDIR)/$(BINARY)-darwin-amd64 ./cmd/$(BINARY)

## do-checksum: Generate checksums for all binaries
do-checksum:
	@cd build && sha256sum $(BINARY)-linux-amd64 $(BINARY)-linux-arm64 $(BINARY)-darwin-amd64 > $(BINARY)_checksum

## build-with-checksum: Build binaries for all platforms and generate checksums
build-with-checksum: build-all do-checksum

.PHONY: install build build-all do-checksum build-with-checksum $(BUILDDIR)/

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

## ci-test-unit: Run unit tests
ci-test-unit:
	@echo Running unit tests via CI...
	@GOFLAGS=$(GOFLAGS) go test ./...

## clean: Clean build files. Runs `go clean` internally.
clean:
	@echo Cleaning build cache...
	@rm -rf $(BUILDDIR) 2> /dev/null
	@go clean ./...

.PHONY: mocks test-unit ci-test-unit clean

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

## build-docker: Build docker image
build-docker:
	@bash $(DOCKERNET_HOME)/build.sh -${build} ${BUILDDIR}

## start-docker: Start docker network
start-docker: build-docker
	@bash $(DOCKERNET_HOME)/start_network.sh

## clean-docker: Clean docker network
clean-docker:
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) stop
	@docker-compose -f $(DOCKERNET_COMPOSE_FILE) down
	rm -rf $(DOCKERNET_HOME)/state
	docker image prune -a

## stop-docker: Stop docker network
stop-docker:
	@bash $(DOCKERNET_HOME)/pkill.sh
	docker-compose -f $(DOCKERNET_COMPOSE_FILE) down --remove-orphans

.PHONY: build-docker start-docker clean-docker stop-docker

GORELEASER_IMAGE := ghcr.io/goreleaser/goreleaser-cross:v$(GO_VERSION)
COSMWASM_VERSION := $(shell go list -m github.com/CosmWasm/wasmvm/v2 | sed 's/.* //')

## release: Build binaries for all platforms and generate checksums
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

## release-dry-run: Dry-run build process for all platforms and generate checksums
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
		--skip=publish

## release-snapshot: Build snapshots for all platforms and generate checksums
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
		--skip=publish

.PHONY: release release-dry-run release-snapshot