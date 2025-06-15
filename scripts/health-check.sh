#!/bin/bash
# DoTask AWS Monitoring and Health Check Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

print_step() { echo -e "${BLUE}>>> $1${NC}"; }
print_success() { echo -e "${GREEN}✅ $1${NC}"; }
print_error() { echo -e "${RED}❌ $1${NC}"; }
print_warning() { echo -e "${YELLOW}⚠️  $1${NC}"; }

print_step "DoTask AWS Health Check"
echo "========================================="

# Check Docker
print_step "Checking Docker status..."
if systemctl is-active --quiet docker; then
    print_success "Docker is running"
else
    print_error "Docker is not running"
    exit 1
fi

# Check containers
print_step "Checking application containers..."
if docker-compose -f docker-compose.aws.yml ps | grep -q "Up"; then
    print_success "Application containers are running"
    docker-compose -f docker-compose.aws.yml ps
else
    print_error "Application containers are not running properly"
    docker-compose -f docker-compose.aws.yml ps
fi

# Check backend health
print_step "Checking backend API health..."
if curl -s http://localhost:8080/health > /dev/null; then
    print_success "Backend API is responding"
else
    print_error "Backend API is not responding"
fi

# Check frontend
print_step "Checking frontend..."
if curl -s http://localhost:80 > /dev/null; then
    print_success "Frontend is responding"
else
    print_error "Frontend is not responding"
fi

# Check database connectivity
print_step "Checking database connectivity..."
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
    if docker run --rm postgres:15-alpine pg_isready -h "$RDS_ENDPOINT" -p 5432 -U postgres; then
        print_success "Database is accessible"
    else
        print_error "Cannot connect to database"
    fi
else
    print_warning "No .env file found, skipping database check"
fi

# System resources
print_step "System resource usage..."
echo "Memory usage:"
free -h
echo ""
echo "Disk usage:"
df -h
echo ""
echo "CPU usage:"
top -bn1 | grep "Cpu(s)" | awk '{print $2 + $4"%"}'

print_step "Health check completed!"
