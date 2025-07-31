.PHONY: build run test clean docker-build docker-up docker-down

# Variables
BINARY_NAME=elys-indexer
DOCKER_IMAGE=elys-indexer:latest

# Build the binary
build:
	go build -o $(BINARY_NAME) ./cmd/indexer

# Run the indexer
run: build
	./$(BINARY_NAME) -c config.yaml

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f $(BINARY_NAME)
	go clean

# Docker commands
docker-build:
	docker build -t $(DOCKER_IMAGE) .

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f indexer

# Database commands
db-setup:
	docker-compose up -d postgres
	sleep 5
	docker-compose exec postgres psql -U indexer -d elys_indexer -f /docker-entrypoint-initdb.d/01-schema.sql

db-migrate:
	# Add migration commands here when needed

# Development commands
dev-deps:
	docker-compose up -d postgres redis

dev: dev-deps
	go run ./cmd/indexer -c config.yaml -l debug

# Format code
fmt:
	go fmt ./...
	gofmt -s -w .

# Lint code
lint:
	golangci-lint run

# Generate mocks (if needed)
mocks:
	go generate ./...