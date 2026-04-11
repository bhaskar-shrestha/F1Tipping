# Backend Dockerfile - Multi-stage build for Go API

# Stage 1: Build
FROM golang:1.21-alpine AS builder

WORKDIR /build

# Install dependencies for CGO if needed
RUN apk add --no-cache gcc musl-dev

# Copy go mod files
COPY go.mod ./
# Download dependencies (go.sum will be generated if it doesn't exist)
RUN go mod download || true

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./src/main.go

# Stage 2: Runtime
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS and curl for health checks
RUN apk add --no-cache ca-certificates curl

# Copy binary from builder
COPY --from=builder /build/main .

# Expose port
EXPOSE 8080

# Run with PORT environment variable (default to 8080)
ENV PORT=8080
ENV GOOS=linux
ENV GOARCH=amd64

CMD ["./main"]
