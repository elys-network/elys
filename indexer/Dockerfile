# Build stage
FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN go build -o elys-indexer ./cmd/indexer

# Final stage
FROM alpine:latest

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/elys-indexer /app/elys-indexer

# Create non-root user
RUN addgroup -g 1000 -S indexer && \
    adduser -u 1000 -S indexer -G indexer

USER indexer

EXPOSE 8080 9090

ENTRYPOINT ["/app/elys-indexer"]