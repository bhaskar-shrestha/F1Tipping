#!/bin/bash

# F1 Tipping App - Quick Start Script

echo "========================================"
echo "  F1 Tipping App - Docker Quick Start"
echo "========================================"
echo ""

# Check if docker is installed
if ! command -v docker &> /dev/null; then
    echo "Error: Docker is not installed or not in PATH"
    exit 1
fi

# Check if docker-compose is installed
if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
    echo "Error: Docker Compose is not installed or not in PATH"
    exit 1
fi

echo "Starting F1 Tipping App with Docker..."
echo ""

# Start services
if docker-compose version &> /dev/null; then
    docker-compose up -d
else
    docker compose up -d
fi

# Wait for services to start
echo "Waiting for services to start..."
sleep 5

# Check status
echo ""
echo "Service Status:"
docker-compose ps 2>/dev/null || docker compose ps

echo ""
echo "========================================"
echo "  Services are running!"
echo "========================================"
echo ""
echo "Web UI: http://localhost"
echo "API:    http://localhost:8080"
echo ""
echo "View logs: docker-compose logs -f"
echo "Stop:      docker-compose down"
echo ""
