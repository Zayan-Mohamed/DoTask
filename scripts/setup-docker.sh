#!/bin/bash
# DoTask Docker Complete Setup and Test Script

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='    cd "$PROJECT_ROOT"
    
    print_info "Starting production cont    # Cleanup    cd "$PROJECT_ROOT"
    
    print_info "Starting development containers..."
    # Cleanup
     print_info "Starting database for testin    # Cleanup
    print_info "Stopping database..."
    $COMPOSE_CMD down."
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD up -d postgres print_info "Stopping development environment..."
    $COMPOSE_CMD -f docker-compose.dev.yml down COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD -f docker-compose.dev.yml up --build -d print_info "Stopping production environment..."
    $COMPOSE_CMD downers..."
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD up --build -d[1;33m'
BLUE='\033[0;34m'
BOLD='\033[1m'
NC='\033[0m'

# Configuration
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

print_header() {
    echo -e "${BOLD}${BLUE}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                    DoTask Docker Setup                         ║"
    echo "║                   Context 7 Implementation                     ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_step() {
    echo -e "${BOLD}${YELLOW}>>> $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

# Function to get the correct docker compose command
get_compose_cmd() {
    if docker compose version &> /dev/null; then
        echo "docker compose"
    else
        echo "docker-compose"
    fi
}

# Function to check prerequisites
check_prerequisites() {
    print_step "Checking Prerequisites"
    
    local errors=0
    
    # Check Docker
    if command -v docker &> /dev/null; then
        print_success "Docker is installed"
        if docker info > /dev/null 2>&1; then
            print_success "Docker is running"
        else
            print_error "Docker is not running"
            errors=$((errors + 1))
        fi
    else
        print_error "Docker is not installed"
        errors=$((errors + 1))
    fi
    
    # Check Docker Compose
    if docker compose version &> /dev/null; then
        print_success "Docker Compose is installed (v2)"
    elif command -v docker-compose &> /dev/null; then
        print_success "Docker Compose (legacy) is installed"
    else
        print_error "Docker Compose is not installed"
        errors=$((errors + 1))
    fi
    
    # Check curl
    if command -v curl &> /dev/null; then
        print_success "curl is available"
    else
        print_error "curl is not installed (needed for health checks)"
        errors=$((errors + 1))
    fi
    
    # Check required files
    local required_files=(
        "docker-compose.yml"
        "docker-compose.dev.yml"
        "do-task/Dockerfile"
        "do-task-backend/Dockerfile"
        "do-task/nginx.conf"
        "docker/init-db.sql"
    )
    
    for file in "${required_files[@]}"; do
        if [ -f "$PROJECT_ROOT/$file" ]; then
            print_success "Found $file"
        else
            print_error "Missing $file"
            errors=$((errors + 1))
        fi
    done
    
    # Check script permissions
    for script in docker-dev.sh docker-prod.sh docker-cleanup.sh docker-db.sh docker-status.sh docker-monitor.sh dotask.sh; do
        if [ -x "$SCRIPT_DIR/$script" ]; then
            print_success "Script $script is executable"
        else
            print_error "Script $script is not executable"
            chmod +x "$SCRIPT_DIR/$script" 2>/dev/null && print_success "Fixed permissions for $script" || errors=$((errors + 1))
        fi
    done
    
    if [ $errors -eq 0 ]; then
        print_success "All prerequisites met!"
    else
        print_error "Found $errors issue(s). Please fix them before continuing."
        exit 1
    fi
    
    echo ""
}

# Function to test production environment
test_production() {
    print_step "Testing Production Environment"
    
    cd "$PROJECT_ROOT"
    
    print_info "Starting production containers..."
    docker compose up --build -d
    
    print_info "Waiting for services to start..."
    sleep 30
    
    # Check container status
    local containers=(dotask-postgres dotask-backend dotask-frontend)
    for container in "${containers[@]}"; do
        if docker ps --filter name="$container" --format "{{.Names}}" | grep -q "$container"; then
            print_success "Container $container is running"
        else
            print_error "Container $container is not running"
            docker logs "$container" --tail=10
            return 1
        fi
    done
    
    # Test endpoints
    print_info "Testing endpoints..."
    
    # Wait for frontend to be ready
    local max_attempts=30
    local attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost/health > /dev/null 2>&1; then
            break
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if curl -s http://localhost/health > /dev/null 2>&1; then
        print_success "Frontend health check passed"
    else
        print_error "Frontend health check failed"
        return 1
    fi
    
    if curl -s http://localhost > /dev/null 2>&1; then
        print_success "Frontend is accessible"
    else
        print_error "Frontend is not accessible"
        return 1
    fi
    
    # Test GraphQL API
    local graphql_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{"query":"query { __schema { types { name } } }"}' \
        http://localhost/api/query 2>/dev/null)
    
    if echo "$graphql_response" | grep -q "types"; then
        print_success "GraphQL API is responding"
    else
        print_error "GraphQL API test failed"
        echo "Response: $graphql_response"
        return 1
    fi
    
    print_success "Production environment test completed!"
    
    # Cleanup
    print_info "Stopping production environment..."
    docker compose down
    
    echo ""
}

# Function to test development environment
test_development() {
    print_step "Testing Development Environment"
    
    cd "$PROJECT_ROOT"
    
    print_info "Starting development containers..."
    docker compose -f docker-compose.dev.yml up --build -d
    
    print_info "Waiting for services to start..."
    sleep 45  # Dev environment needs more time for compilation
    
    # Check container status
    local containers=(dotask-postgres-dev dotask-backend-dev dotask-frontend-dev)
    for container in "${containers[@]}"; do
        if docker ps --filter name="$container" --format "{{.Names}}" | grep -q "$container"; then
            print_success "Container $container is running"
        else
            print_error "Container $container is not running"
            docker logs "$container" --tail=10
            return 1
        fi
    done
    
    # Test endpoints
    print_info "Testing development endpoints..."
    
    # Wait for services to be ready
    local max_attempts=30
    local attempt=0
    
    # Test backend
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:8080 > /dev/null 2>&1; then
            break
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if curl -s http://localhost:8080 > /dev/null 2>&1; then
        print_success "Development backend is accessible"
    else
        print_error "Development backend is not accessible"
        return 1
    fi
    
    # Test frontend (might take longer to compile)
    attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if curl -s http://localhost:5173 > /dev/null 2>&1; then
            break
        fi
        attempt=$((attempt + 1))
        sleep 3
    done
    
    if curl -s http://localhost:5173 > /dev/null 2>&1; then
        print_success "Development frontend is accessible"
    else
        print_error "Development frontend is not accessible"
        docker logs dotask-frontend-dev --tail=20
        return 1
    fi
    
    # Test GraphQL API
    local graphql_response=$(curl -s -X POST \
        -H "Content-Type: application/json" \
        -d '{"query":"query { __schema { types { name } } }"}' \
        http://localhost:8080/query 2>/dev/null)
    
    if echo "$graphql_response" | grep -q "types"; then
        print_success "Development GraphQL API is responding"
    else
        print_error "Development GraphQL API test failed"
        echo "Response: $graphql_response"
        return 1
    fi
    
    print_success "Development environment test completed!"
    
    # Cleanup
    print_info "Stopping development environment..."
    docker compose -f docker-compose.dev.yml down
    
    echo ""
}

# Function to test database operations
test_database() {
    print_step "Testing Database Operations"
    
    cd "$PROJECT_ROOT"
    
    print_info "Starting database for testing..."
    docker compose up -d postgres
    
    # Wait for database
    local max_attempts=30
    local attempt=0
    while [ $attempt -lt $max_attempts ]; do
        if docker exec dotask-postgres pg_isready -U postgres -d dotask > /dev/null 2>&1; then
            break
        fi
        attempt=$((attempt + 1))
        sleep 2
    done
    
    if docker exec dotask-postgres pg_isready -U postgres -d dotask > /dev/null 2>&1; then
        print_success "Database connection established"
    else
        print_error "Database connection failed"
        return 1
    fi
    
    # Test database operations
    print_info "Testing database operations..."
    
    # Test backup
    if "$SCRIPT_DIR/docker-db.sh" backup > /dev/null 2>&1; then
        print_success "Database backup test passed"
        # Clean up backup file
        rm -f backup_*.sql 2>/dev/null
    else
        print_error "Database backup test failed"
    fi
    
    # Test status
    if "$SCRIPT_DIR/docker-db.sh" status > /dev/null 2>&1; then
        print_success "Database status check passed"
    else
        print_error "Database status check failed"
    fi
    
    print_success "Database tests completed!"
    
    # Cleanup
    print_info "Stopping database..."
    docker compose down
    
    echo ""
}

# Function to test management scripts
test_scripts() {
    print_step "Testing Management Scripts"
    
    # Test dotask master script
    if "$SCRIPT_DIR/dotask.sh" help > /dev/null 2>&1; then
        print_success "Master script (dotask.sh) works"
    else
        print_error "Master script (dotask.sh) failed"
    fi
    
    # Test individual scripts
    local scripts=(docker-dev.sh docker-prod.sh docker-cleanup.sh docker-status.sh docker-monitor.sh)
    for script in "${scripts[@]}"; do
        if [ -x "$SCRIPT_DIR/$script" ]; then
            print_success "Script $script is executable and available"
        else
            print_error "Script $script has issues"
        fi
    done
    
    print_success "Script tests completed!"
    echo ""
}

# Function to display usage information
show_usage() {
    print_step "Usage Information"
    
    echo -e "${BOLD}After setup completion, you can use:${NC}"
    echo ""
    echo -e "${GREEN}Development:${NC}"
    echo "  ./scripts/dotask.sh dev          # Start development environment"
    echo "  ./scripts/docker-dev.sh          # Alternative way"
    echo ""
    echo -e "${GREEN}Production:${NC}"
    echo "  ./scripts/dotask.sh prod         # Start production environment"
    echo "  ./scripts/docker-prod.sh         # Alternative way"
    echo ""
    echo -e "${GREEN}Management:${NC}"
    echo "  ./scripts/dotask.sh status       # Check system status"
    echo "  ./scripts/dotask.sh logs         # View logs"
    echo "  ./scripts/dotask.sh clean        # Cleanup containers"
    echo "  ./scripts/docker-monitor.sh      # Health monitoring"
    echo ""
    echo -e "${GREEN}Database:${NC}"
    echo "  ./scripts/dotask.sh db status    # Database status"
    echo "  ./scripts/dotask.sh db backup    # Create backup"
    echo "  ./scripts/dotask.sh db shell     # Open database shell"
    echo ""
    echo -e "${BOLD}Access URLs:${NC}"
    echo -e "${GREEN}Production:${NC}"
    echo "  • Application: http://localhost"
    echo "  • API: http://localhost/api"
    echo "  • Health: http://localhost/health"
    echo ""
    echo -e "${GREEN}Development:${NC}"
    echo "  • Frontend: http://localhost:5173"
    echo "  • Backend: http://localhost:8080"
    echo "  • GraphQL Playground: http://localhost:8080"
    echo ""
}

# Function to run complete setup
run_setup() {
    print_header
    
    echo -e "${BOLD}This script will test your complete DoTask Docker setup.${NC}"
    echo -e "${BOLD}It will start and stop containers during testing.${NC}"
    echo ""
    
    case "$1" in
        "--quick"|"-q")
            print_info "Running quick setup (prerequisites and scripts only)..."
            check_prerequisites
            test_scripts
            ;;
        "--full"|"-f")
            print_info "Running full setup test (includes starting containers)..."
            check_prerequisites
            test_scripts
            test_database
            test_production
            test_development
            ;;
        "--prod-only")
            print_info "Testing production environment only..."
            check_prerequisites
            test_production
            ;;
        "--dev-only")
            print_info "Testing development environment only..."
            check_prerequisites
            test_development
            ;;
        *)
            print_info "Running standard setup (prerequisites, scripts, and database)..."
            check_prerequisites
            test_scripts
            test_database
            ;;
    esac
    
    echo -e "${BOLD}${GREEN}"
    echo "╔════════════════════════════════════════════════════════════════╗"
    echo "║                   🎉 Setup Complete! 🎉                       ║"
    echo "║              Your DoTask Docker environment is ready!         ║"
    echo "╚════════════════════════════════════════════════════════════════╝"
    echo -e "${NC}"
    
    show_usage
}

# Main execution
case "$1" in
    "--help"|"-h"|"help")
        echo "DoTask Docker Setup Script"
        echo "=========================="
        echo ""
        echo "Usage: $0 [OPTIONS]"
        echo ""
        echo "Options:"
        echo "  --quick, -q      Quick setup (prerequisites and scripts only)"
        echo "  --full, -f       Full setup test (includes container testing)"
        echo "  --prod-only      Test production environment only"
        echo "  --dev-only       Test development environment only"
        echo "  --help, -h       Show this help"
        echo ""
        echo "Default: Standard setup (prerequisites, scripts, database)"
        ;;
    *)
        run_setup "$@"
        ;;
esac
