#!/bin/bash
# Database migration script

set -e

echo "Running database migrations..."

cd "$(dirname "$0")/.."

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '^#' | xargs)
fi

# Run migrations using Go
go run backend/shared/cmd/migrate/main.go

echo "Migrations completed successfully!"

