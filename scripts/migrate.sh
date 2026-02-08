#!/bin/bash

set -e

DATABASE_URL=${DATABASE_URL:-"postgres://urlshortener:password@localhost:5432/urlshortener?sslmode=disable"}

echo "Running migrations..."

for migration in migrations/*.sql; do
    echo "Applying $migration..."
    psql "$DATABASE_URL" -f "$migration"
done

echo "Migrations completed successfully!"