# Build stage
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git make

# Create the expected directory structure for the replace directive
WORKDIR /workspace

# Copy the entire elys project to maintain the correct relative paths
COPY . ./elys

# Change to indexer directory
WORKDIR /workspace/elys/indexer

# Download dependencies
RUN go mod download

# Build the binary
RUN go build -o elys-indexer ./cmd/indexer

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /workspace/elys/indexer/elys-indexer /app/elys-indexer

# Create non-root user
RUN addgroup -g 1000 -S indexer && \
    adduser -u 1000 -S indexer -G indexer

USER indexer

EXPOSE 8080 9090

ENTRYPOINT ["/app/elys-indexer"]