#!/bin/bash
# DoTask Docker Master Control Script

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Function to get the correct docker compose command
get_compose_cmd() {
    if docker compose version &> /dev/null; then
        echo "docker compose"
    else
        echo "docker-compose"
    fi
}

# Function to check Docker installation
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    
    if ! docker info > /dev/null 2>&1; then
        print_error "Docker is not running. Please start Docker and try again."
        exit 1
    fi
    
    if ! docker compose version &> /dev/null && ! command -v docker-compose &> /dev/null; then
        print_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

# Function to show system status
show_status() {
    print_status "DoTask System Status"
    echo "===================="
    
    # Check Docker
    check_docker
    print_success "Docker is running"
    
    # Check containers
    if [ "$(docker ps -q --filter name=dotask)" ]; then
        print_success "DoTask containers are running"
        docker ps --filter name=dotask --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    else
        print_warning "No DoTask containers are running"
    fi
    
    echo ""
}

# Function to show help
show_help() {
    echo "DoTask Docker Management"
    echo "======================="
    echo ""
    echo "Usage: $0 [COMMAND]"
    echo ""
    echo "Commands:"
    echo "  dev               Start development environment"
    echo "  prod              Start production environment"
    echo "  stop              Stop all DoTask containers"
    echo "  status            Show system status"
    echo "  logs              Show recent logs"
    echo "  follow            Follow logs in real-time"
    echo "  clean             Clean up containers (keep data)"
    echo "  clean-all         Full cleanup (remove all data)"
    echo "  rebuild           Rebuild and restart containers"
    echo "  db                Open database management menu"
    echo "  health            Check container health"
    echo "  help              Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 dev            # Start development environment"
    echo "  $0 prod           # Start production environment"
    echo "  $0 logs           # Show recent logs"
    echo "  $0 clean          # Clean up containers"
    echo ""
}

# Function to start development environment
start_dev() {
    print_status "Starting DoTask Development Environment..."
    check_docker
    
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD -f docker-compose.dev.yml up --build -d
    
    if [ $? -eq 0 ]; then
        print_success "Development environment started successfully!"
        echo ""
        echo "🌐 Available endpoints:"
        echo "  • Frontend (Dev): http://localhost:5173"
        echo "  • Backend (Dev): http://localhost:8080"
        echo "  • GraphQL Playground: http://localhost:8080"
        echo "  • Database: localhost:5432"
        echo ""
        echo "💡 Use '$0 logs' to see application logs"
        echo "💡 Use '$0 stop' to stop the environment"
    else
        print_error "Failed to start development environment"
        exit 1
    fi
}

# Function to start production environment
start_prod() {
    print_status "Starting DoTask Production Environment..."
    check_docker
    
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD up --build -d
    
    if [ $? -eq 0 ]; then
        print_success "Production environment started successfully!"
        echo ""
        echo "🌐 Available endpoints:"
        echo "  • Application: http://localhost"
        echo "  • API: http://localhost/api"
        echo "  • Health Check: http://localhost/health"
        echo ""
        echo "💡 Use '$0 logs' to see application logs"
        echo "💡 Use '$0 stop' to stop the environment"
    else
        print_error "Failed to start production environment"
        exit 1
    fi
}

# Function to stop containers
stop_containers() {
    print_status "Stopping DoTask containers..."
    
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD down 2>/dev/null
    $COMPOSE_CMD -f docker-compose.dev.yml down 2>/dev/null
    
    print_success "DoTask containers stopped"
}

# Function to show logs
show_logs() {
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    if docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        print_status "Showing development environment logs..."
        $COMPOSE_CMD -f docker-compose.dev.yml logs --tail=50
    elif docker ps --filter name=dotask-frontend -q | grep -q .; then
        print_status "Showing production environment logs..."
        $COMPOSE_CMD logs --tail=50
    else
        print_warning "No DoTask containers are running"
    fi
}

# Function to follow logs
follow_logs() {
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    if docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        print_status "Following development environment logs (Ctrl+C to stop)..."
        $COMPOSE_CMD -f docker-compose.dev.yml logs -f
    elif docker ps --filter name=dotask-frontend -q | grep -q .; then
        print_status "Following production environment logs (Ctrl+C to stop)..."
        $COMPOSE_CMD logs -f
    else
        print_warning "No DoTask containers are running"
    fi
}

# Function to clean containers
clean_containers() {
    print_status "Cleaning up DoTask containers (preserving data)..."
    
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    $COMPOSE_CMD down --remove-orphans 2>/dev/null
    $COMPOSE_CMD -f docker-compose.dev.yml down --remove-orphans 2>/dev/null
    docker system prune -f
    
    print_success "Cleanup completed (data preserved)"
}

# Function to full clean
full_clean() {
    print_warning "This will remove ALL DoTask data including the database!"
    echo -n "Are you sure? (y/N): "
    read -r response
    
    if [[ "$response" =~ ^[Yy]$ ]]; then
        print_status "Performing full cleanup..."
        
        cd "$PROJECT_ROOT"
        COMPOSE_CMD=$(get_compose_cmd)
        $COMPOSE_CMD down --volumes --remove-orphans 2>/dev/null
        $COMPOSE_CMD -f docker-compose.dev.yml down --volumes --remove-orphans 2>/dev/null
        docker system prune -f --volumes
        
        print_success "Full cleanup completed (all data removed)"
    else
        print_status "Full cleanup cancelled"
    fi
}

# Function to rebuild containers
rebuild_containers() {
    print_status "Rebuilding DoTask containers..."
    
    cd "$PROJECT_ROOT"
    COMPOSE_CMD=$(get_compose_cmd)
    
    # Determine which environment is running
    if docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        print_status "Rebuilding development environment..."
        $COMPOSE_CMD -f docker-compose.dev.yml down
        $COMPOSE_CMD -f docker-compose.dev.yml build --no-cache
        $COMPOSE_CMD -f docker-compose.dev.yml up -d
    elif docker ps --filter name=dotask-frontend -q | grep -q .; then
        print_status "Rebuilding production environment..."
        $COMPOSE_CMD down
        $COMPOSE_CMD build --no-cache
        $COMPOSE_CMD up -d
    else
        print_warning "No DoTask environment is currently running"
        echo "Use '$0 dev' or '$0 prod' to start an environment first"
        exit 1
    fi
    
    print_success "Rebuild completed"
}

# Function to check container health
check_health() {
    print_status "Checking DoTask container health..."
    
    if [ "$(docker ps -q --filter name=dotask)" ]; then
        echo ""
        docker ps --filter name=dotask --format "table {{.Names}}\t{{.Status}}"
        echo ""
        
        # Test endpoints
        print_status "Testing endpoints..."
        
        if docker ps --filter name=dotask-frontend --format "{{.Names}}" | grep -q dotask-frontend; then
            if curl -s http://localhost/health > /dev/null; then
                print_success "Frontend health check: ✅ OK"
            else
                print_error "Frontend health check: ❌ FAILED"
            fi
        fi
        
        if docker ps --filter name=dotask-backend --format "{{.Names}}" | grep -q dotask-backend; then
            if curl -s http://localhost/api > /dev/null; then
                print_success "Backend health check: ✅ OK"
            else
                print_error "Backend health check: ❌ FAILED"
            fi
        fi
        
        if docker ps --filter name=dotask-frontend-dev --format "{{.Names}}" | grep -q dotask-frontend-dev; then
            if curl -s http://localhost:5173 > /dev/null; then
                print_success "Dev frontend health check: ✅ OK"
            else
                print_error "Dev frontend health check: ❌ FAILED"
            fi
        fi
        
        if docker ps --filter name=dotask-backend-dev --format "{{.Names}}" | grep -q dotask-backend-dev; then
            if curl -s http://localhost:8080 > /dev/null; then
                print_success "Dev backend health check: ✅ OK"
            else
                print_error "Dev backend health check: ❌ FAILED"
            fi
        fi
    else
        print_warning "No DoTask containers are running"
    fi
}

# Function to open database menu
database_menu() {
    print_status "Opening database management..."
    exec "$SCRIPT_DIR/docker-db.sh" "$@"
}

# Main script logic
case "$1" in
    "dev"|"development")
        start_dev
        ;;
    "prod"|"production")
        start_prod
        ;;
    "stop"|"down")
        stop_containers
        ;;
    "status")
        show_status
        ;;
    "logs")
        show_logs
        ;;
    "follow"|"logs-follow")
        follow_logs
        ;;
    "clean"|"cleanup")
        clean_containers
        ;;
    "clean-all"|"cleanup-all")
        full_clean
        ;;
    "rebuild")
        rebuild_containers
        ;;
    "health"|"healthcheck")
        check_health
        ;;
    "db"|"database")
        shift
        database_menu "$@"
        ;;
    "help"|"--help"|"-h")
        show_help
        ;;
    "")
        show_status
        echo ""
        echo "💡 Use '$0 help' to see available commands"
        ;;
    *)
        print_error "Unknown command: $1"
        echo ""
        show_help
        exit 1
        ;;
esac
