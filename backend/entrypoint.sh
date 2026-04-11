#!/bin/bash

# F1 Tipping Backend Entrypoint Script
# This script handles database initialization before starting the application

set -e

echo "========================================"
echo "F1 Tipping Backend Entrypoint"
echo "========================================"
echo "Database Host: ${DB_HOST:-localhost}"
echo "Database Port: ${DB_PORT:-5432}"
echo "Database Name: ${DB_NAME:-f1_prediction}"
echo "Database User: ${DB_USER:-postgres}"
echo "========================================"

# Wait for database to be ready
echo "Waiting for PostgreSQL to be ready..."
MAX_ATTEMPTS=30
ATTEMPT=0

while [ $ATTEMPT -lt $MAX_ATTEMPTS ]; do
    if pg_isready -h "${DB_HOST:-localhost}" -p "${DB_PORT:-5432}" -U "${DB_USER:-postgres}" 2>/dev/null; then
        echo "PostgreSQL is ready!"
        break
    else
        ATTEMPT=$((ATTEMPT + 1))
        echo "Attempt $ATTEMPT/$MAX_ATTEMPTS: PostgreSQL not ready, waiting 2 seconds..."
        sleep 2
    fi
done

if [ $ATTEMPT -eq $MAX_ATTEMPTS ]; then
    echo "ERROR: PostgreSQL not ready after $MAX_ATTEMPTS attempts"
    exit 1
fi

echo "========================================"
echo "Starting F1 Tipping Backend..."
echo "========================================"

# Run the main application
exec ./main
