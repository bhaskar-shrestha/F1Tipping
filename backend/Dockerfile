# Backend Dockerfile - Multi-stage build for Go API with PostgreSQL

# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install necessary build packages
RUN apk add --no-cache git ca-certificates

# Copy only go.mod and go.sum for dependency caching
COPY go.mod go.sum ./

# Cache Go module dependencies
RUN go mod download && go mod graph

# Copy all source-code
COPY . .

# Build the binary with all Go build flags
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags="-w -s" -o /app/main ./src/main.go

# Stage 2: Runtime  
FROM alpine:latest

WORKDIR /app

# Install runtime dependencies (ca-certificates for HTTPS, tzdata for timezone support)
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main ./main

# Copy database schema to expected location
COPY src/db ./src/db

# Set environment
ENV GIN_MODE=release

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD [ "test", "-f", "/app/main" ]

CMD ["./main"]
