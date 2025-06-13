#!/bin/bash
# DoTask Health Check and Monitoring Script

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

# Configuration
FRONTEND_URL="http://localhost"
BACKEND_URL="http://localhost/api"
DEV_FRONTEND_URL="http://localhost:5173"
DEV_BACKEND_URL="http://localhost:8080"
DB_CONTAINER="dotask-postgres"

# Function to print colored output
print_status() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }

# Function to check HTTP endpoint
check_endpoint() {
    local url=$1
    local name=$2
    local timeout=${3:-5}
    
    if curl -s --max-time $timeout "$url" > /dev/null 2>&1; then
        print_success "$name: ✅ UP"
        return 0
    else
        print_error "$name: ❌ DOWN"
        return 1
    fi
}

# Function to check container health
check_container() {
    local container=$1
    local name=$2
    
    if docker ps --filter name="$container" --format "{{.Names}}" | grep -q "$container"; then
        local status=$(docker inspect --format='{{.State.Health.Status}}' "$container" 2>/dev/null)
        if [ "$status" = "healthy" ]; then
            print_success "$name Container: ✅ HEALTHY"
        elif [ "$status" = "unhealthy" ]; then
            print_error "$name Container: ❌ UNHEALTHY"
        else
            print_warning "$name Container: ⚠️ NO HEALTH CHECK"
        fi
        return 0
    else
        print_error "$name Container: ❌ NOT RUNNING"
        return 1
    fi
}

# Function to check database connectivity
check_database() {
    # Check if database container is running first
    local db_container_name
    if docker ps --filter name=dotask-postgres-dev -q | grep -q .; then
        db_container_name="dotask-postgres-dev"
    elif docker ps --filter name=dotask-postgres -q | grep -q .; then
        db_container_name="dotask-postgres"
    else
        print_error "Database: ❌ CONTAINER NOT RUNNING"
        return 1
    fi
    
    if docker exec "$db_container_name" pg_isready -U postgres -d dotask > /dev/null 2>&1; then
        print_success "Database: ✅ CONNECTED"
        
        # Get database stats
        local db_size=$(docker exec "$db_container_name" psql -U postgres dotask -t -c "SELECT pg_size_pretty(pg_database_size('dotask'));" 2>/dev/null | xargs)
        local connection_count=$(docker exec "$db_container_name" psql -U postgres dotask -t -c "SELECT count(*) FROM pg_stat_activity WHERE datname='dotask';" 2>/dev/null | xargs)
        
        echo "  📊 Database size: $db_size"
        echo "  🔗 Active connections: $connection_count"
        return 0
    else
        print_error "Database: ❌ NOT ACCESSIBLE"
        return 1
    fi
}

# Function to get resource usage
check_resources() {
    print_status "Resource Usage:"
    
    # Docker system info
    local containers=$(docker ps --filter name=dotask -q | wc -l)
    echo "  📦 Running containers: $containers"
    
    # Get memory usage for each container
    if [ $containers -gt 0 ]; then
        echo "  💾 Memory usage:"
        docker stats --no-stream --format "table   {{.Name}}\t{{.MemUsage}}\t{{.CPUPerc}}" $(docker ps --filter name=dotask --format "{{.Names}}")
    fi
    
    # Disk usage
    echo "  💿 Volume usage:"
    docker system df | grep -E "(TYPE|dotask|postgres)" || echo "    No DoTask volumes found"
}

# Function to test API endpoints
test_api() {
    print_status "Testing API endpoints..."
    
    # Determine which environment is running
    local backend_url
    if docker ps --filter name=dotask-backend-dev -q | grep -q .; then
        backend_url="$DEV_BACKEND_URL"
    else
        backend_url="$BACKEND_URL"
    fi
    
    # Test GraphQL endpoint
    local graphql_response=$(curl -s --max-time 5 -X POST \
        -H "Content-Type: application/json" \
        -d '{"query":"query { __schema { types { name } } }"}' \
        "$backend_url/query" 2>/dev/null)
    
    if echo "$graphql_response" | grep -q "types"; then
        print_success "GraphQL API: ✅ RESPONDING"
    else
        print_error "GraphQL API: ❌ NOT RESPONDING"
        echo "  Response: ${graphql_response:0:100}..."
    fi
}

# Function to check logs for errors
check_logs() {
    print_status "Checking recent logs for errors..."
    
    local error_count=0
    
    # Check backend logs (exclude common warnings)
    if docker ps --filter name=dotask-backend -q | grep -q .; then
        local backend_errors=$(docker logs dotask-backend --since=5m 2>&1 | grep -i "error\|fatal\|panic" | grep -v "trusted all proxies" | wc -l)
        if [ $backend_errors -gt 0 ]; then
            print_warning "Backend: $backend_errors error(s) in last 5 minutes"
            error_count=$((error_count + backend_errors))
        fi
    elif docker ps --filter name=dotask-backend-dev -q | grep -q .; then
        local backend_errors=$(docker logs dotask-backend-dev --since=5m 2>&1 | grep -i "error\|fatal\|panic" | grep -v "trusted all proxies" | wc -l)
        if [ $backend_errors -gt 0 ]; then
            print_warning "Backend: $backend_errors error(s) in last 5 minutes"
            error_count=$((error_count + backend_errors))
        fi
    fi
    
    # Check frontend logs (exclude Vite optimization messages)
    if docker ps --filter name=dotask-frontend -q | grep -q .; then
        local frontend_errors=$(docker logs dotask-frontend --since=5m 2>&1 | grep -i "error\|fatal" | grep -v "optimized\|dependencies" | wc -l)
        if [ $frontend_errors -gt 0 ]; then
            print_warning "Frontend: $frontend_errors error(s) in last 5 minutes"
            error_count=$((error_count + frontend_errors))
        fi
    elif docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        local frontend_errors=$(docker logs dotask-frontend-dev --since=5m 2>&1 | grep -i "error\|fatal" | grep -v "optimized\|dependencies\|new dependencies\|reloading" | wc -l)
        if [ $frontend_errors -gt 0 ]; then
            print_warning "Frontend: $frontend_errors error(s) in last 5 minutes"
            error_count=$((error_count + frontend_errors))
        fi
    fi
    
    # Check database logs (exclude normal recovery messages)
    if docker ps --filter name=dotask-postgres -q | grep -q .; then
        local db_errors=$(docker logs dotask-postgres --since=5m 2>&1 | grep -i "error\|fatal" | grep -v "recovery\|checkpoint\|interrupted" | wc -l)
        if [ $db_errors -gt 0 ]; then
            print_warning "Database: $db_errors error(s) in last 5 minutes"
            error_count=$((error_count + db_errors))
        fi
    elif docker ps --filter name=dotask-postgres-dev -q | grep -q .; then
        local db_errors=$(docker logs dotask-postgres-dev --since=5m 2>&1 | grep -i "error\|fatal" | grep -v "recovery\|checkpoint\|interrupted" | wc -l)
        if [ $db_errors -gt 0 ]; then
            print_warning "Database: $db_errors error(s) in last 5 minutes"
            error_count=$((error_count + db_errors))
        fi
    fi
    
    if [ $error_count -eq 0 ]; then
        print_success "No errors found in recent logs"
    else
        print_warning "Total errors in last 5 minutes: $error_count"
    fi
}

# Function to run comprehensive health check
run_health_check() {
    echo "🏥 DoTask Health Check Report"
    echo "============================="
    echo "Timestamp: $(date)"
    echo ""
    
    local overall_status=0
    
    # Check containers
    print_status "Container Health:"
    check_container "dotask-postgres" "Database" || overall_status=1
    
    if docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        check_container "dotask-frontend-dev" "Frontend (Dev)" || overall_status=1
        check_container "dotask-backend-dev" "Backend (Dev)" || overall_status=1
    else
        check_container "dotask-frontend" "Frontend" || overall_status=1
        check_container "dotask-backend" "Backend" || overall_status=1
    fi
    echo ""
    
    # Check endpoints
    print_status "Endpoint Health:"
    if docker ps --filter name=dotask-frontend-dev -q | grep -q .; then
        check_endpoint "$DEV_FRONTEND_URL" "Frontend (Dev)" || overall_status=1
        check_endpoint "$DEV_BACKEND_URL" "Backend (Dev)" || overall_status=1
    else
        check_endpoint "$FRONTEND_URL" "Frontend" || overall_status=1
        check_endpoint "$FRONTEND_URL/health" "Health Check" || overall_status=1
    fi
    echo ""
    
    # Check database
    print_status "Database Health:"
    check_database || overall_status=1
    echo ""
    
    # Test API
    test_api
    echo ""
    
    # Check resources
    check_resources
    echo ""
    
    # Check logs
    check_logs
    echo ""
    
    # Overall status
    if [ $overall_status -eq 0 ]; then
        print_success "🎉 Overall Status: HEALTHY"
    else
        print_error "⚠️ Overall Status: ISSUES DETECTED"
    fi
    
    return $overall_status
}

# Function to generate monitoring report
generate_report() {
    local report_file="health_report_$(date +%Y%m%d_%H%M%S).txt"
    run_health_check | tee "$report_file"
    echo ""
    print_success "Report saved to: $report_file"
}

# Function to watch system continuously
watch_system() {
    print_status "Starting continuous monitoring (Ctrl+C to stop)..."
    echo ""
    
    while true; do
        clear
        run_health_check
        echo ""
        echo "Next check in 30 seconds..."
        sleep 30
    done
}

# Main script logic
case "$1" in
    "check"|"")
        run_health_check
        ;;
    "report")
        generate_report
        ;;
    "watch")
        watch_system
        ;;
    "api")
        test_api
        ;;
    "resources")
        check_resources
        ;;
    "logs")
        check_logs
        ;;
    *)
        echo "DoTask Health Monitor"
        echo "===================="
        echo ""
        echo "Usage: $0 [COMMAND]"
        echo ""
        echo "Commands:"
        echo "  check     - Run full health check (default)"
        echo "  report    - Generate and save health report"
        echo "  watch     - Continuous monitoring"
        echo "  api       - Test API endpoints only"
        echo "  resources - Check resource usage only"
        echo "  logs      - Check logs for errors only"
        echo ""
        ;;
esac
