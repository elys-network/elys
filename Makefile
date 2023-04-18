NAME?=elys
BINARY?=$(NAME)d
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags --match 'v*' --abbrev=8 | sed 's/-g/-/' | sed 's/-[0-9]*-/-/')

GOFLAGS:=""
GOTAGS:=ledger

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.AppName=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(GOTAGS)'

build: check-version go.sum
	GOFLAGS=$(GOFLAGS) go build $(BUILD_FLAGS) -o ./build/$(BINARY) ./cmd/$(BINARY)

install: check-version go.sum
	GOFLAGS=$(GOFLAGS) go install $(BUILD_FLAGS) ./cmd/$(BINARY)

build-all:
	GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-linux-amd64 ./cmd/$(BINARY)/main.go
	GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-linux-arm64 ./cmd/$(BINARY)/main.go
	GOFLAGS=$(GOFLAGS) GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-darwin-amd64 ./cmd/$(BINARY)/main.go

do-checksum:
	cd build && sha256sum $(BINARY)-linux-amd64 $(BINARY)-linux-arm64 $(BINARY)-darwin-amd64 > $(BINARY)_checksum

build-with-checksum: build-all do-checksum

go-mod-cache: go.sum
	@echo "--> Retrieve the go modules and store them in the local cache."
	@go mod download

go.sum: go.mod
	@echo "--> Make sure that the dependencies haven't been altered."
	@go mod verify

clean:
	@rm -rf build

test:
	GOFLAGS=$(GOFLAGS) go test -v ./...

# Add check to make sure we are using the proper Go version before proceeding with anything
check-version:
	@if ! go version | grep -q "go1.19"; then \
		echo "\033[0;31mERROR:\033[0m Go version 1.19 is required for compiling elysd. It looks like you are using" "$(shell go version) \nThere are potential consensus-breaking changes that can occur when running binaries compiled with different versions of Go. Please download Go version 1.19 and retry. Thank you!"; \
		exit 1; \
	fi