#!/bin/bash
# Database Migration Script for AWS RDS

set -e

# Load environment variables
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m'

print_step() { echo -e "${BLUE}>>> $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }

print_step "Starting database migration..."

# Check if RDS is accessible
print_step "Testing RDS connection..."
if ! docker run --rm postgres:15-alpine pg_isready -h "$RDS_ENDPOINT" -p 5432 -U postgres; then
    print_error "Cannot connect to RDS instance. Check your security groups and RDS endpoint."
    exit 1
fi

print_success "RDS connection successful!"

# Run migrations
print_step "Running database migrations..."
cd do-task-backend

# Use docker to run migrations since we might not have Go installed locally
docker run --rm \
    -v "$(pwd)/migrations:/migrations" \
    -e "DATABASE_URL=$DATABASE_URL" \
    migrate/migrate:latest \
    -path=/migrations \
    -database="$DATABASE_URL" \
    up

print_success "Database migrations completed!"

# Test database with a simple query
print_step "Testing database with sample query..."
docker run --rm postgres:15-alpine \
    psql "$DATABASE_URL" \
    -c "SELECT version();"

print_success "Database setup completed successfully!"
