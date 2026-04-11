# Backend Dockerfile - Multi-stage build for Go API with PostgreSQL

# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install necessary packages
RUN apk add --no-cache git

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY src/ ./src/

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/main.go

# Stage 2: Runtime
FROM alpine:latest

WORKDIR /app

# Install necessary packages
RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    bash \
    postgresql-client

# Copy binary from builder
COPY --from=builder /build/main /app/main

# Only set non-sensitive environment variables
ENV GIN_MODE=release

# Database and sensitive configs are injected at runtime via docker-compose
EXPOSE 8080

# Run the application directly
CMD ["./main"]
