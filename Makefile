NAME?=elys
BINARY?=$(NAME)d
COMMIT:=$(shell git log -1 --format='%H')
VERSION:=$(shell git describe --tags)

GOFLAGS:=""
GOTAGS:=ledger

ldflags = -X github.com/cosmos/cosmos-sdk/version.Name=$(NAME) \
		  -X github.com/cosmos/cosmos-sdk/version.ServerName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.ClientName=$(BINARY) \
		  -X github.com/cosmos/cosmos-sdk/version.Version=$(VERSION) \
		  -X github.com/cosmos/cosmos-sdk/version.Commit=$(COMMIT)

BUILD_FLAGS := -ldflags '$(ldflags)' -tags '$(GOTAGS)'

build:
	GOFLAGS=$(GOFLAGS) go build $(BUILD_FLAGS) -o ./build/$(BINARY) ./cmd/$(BINARY)/main.go

install:
	GOFLAGS=$(GOFLAGS) go install $(BUILD_FLAGS) ./cmd/$(BINARY)/main.go

build-all:
	GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-linux-amd64 ./cmd/$(BINARY)/main.go
	GOFLAGS=$(GOFLAGS) GOOS=linux GOARCH=arm64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-linux-arm64 ./cmd/$(BINARY)/main.go
	GOFLAGS=$(GOFLAGS) GOOS=darwin GOARCH=amd64 go build $(BUILD_FLAGS) -o ./build/$(BINARY)-darwin-amd64 ./cmd/$(BINARY)/main.go

do-checksum:
	cd build && sha256sum $(BINARY)-linux-amd64 $(BINARY)-linux-arm64 $(BINARY)-darwin-amd64 > $(BINARY)_checksum

build-with-checksum: build-all do-checksum

clean:
	@rm -rf build

test:
	GOFLAGS=$(GOFLAGS) go test -v ./...