#!/bin/bash
# DoTask Docker Status and Logs Script

echo "📊 DoTask Docker Status"
echo "======================="

# Check if any DoTask containers are running
if [ "$(docker ps -q --filter name=dotask)" ]; then
    echo "✅ DoTask containers are running"
    echo ""
    
    echo "📋 Container Status:"
    docker ps --filter name=dotask --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
    echo ""
    
    echo "🏥 Health Status:"
    docker ps --filter name=dotask --format "table {{.Names}}\t{{.Status}}" | grep -E "(healthy|unhealthy)" || echo "No health checks configured"
    echo ""
    
    echo "💾 Volume Usage:"
    docker system df -v | grep dotask || echo "No DoTask volumes found"
    echo ""
    
    echo "🌐 Available Endpoints:"
    if docker ps --filter name=dotask-frontend --format "{{.Names}}" | grep -q dotask-frontend; then
        echo "  • Frontend: http://localhost"
        echo "  • Health Check: http://localhost/health"
    fi
    
    if docker ps --filter name=dotask-backend --format "{{.Names}}" | grep -q dotask-backend; then
        echo "  • Backend API: http://localhost/api"
        echo "  • GraphQL: http://localhost/api/query"
    fi
    
    if docker ps --filter name=dotask-frontend-dev --format "{{.Names}}" | grep -q dotask-frontend-dev; then
        echo "  • Dev Frontend: http://localhost:5173"
    fi
    
    if docker ps --filter name=dotask-backend-dev --format "{{.Names}}" | grep -q dotask-backend-dev; then
        echo "  • Dev Backend: http://localhost:8080"
        echo "  • GraphQL Playground: http://localhost:8080"
    fi
    
    if docker ps --filter name=dotask-postgres --format "{{.Names}}" | grep -q dotask-postgres; then
        echo "  • Database: localhost:5432"
    fi
    
    echo ""
    
    # Show recent logs if requested
    if [ "$1" = "--logs" ] || [ "$1" = "-l" ]; then
        echo "📋 Recent Logs (last 20 lines):"
        echo "================================"
        if docker compose version &> /dev/null; then
            docker compose logs --tail=20 2>/dev/null || docker compose -f docker-compose.dev.yml logs --tail=20
        else
            docker-compose logs --tail=20 2>/dev/null || docker-compose -f docker-compose.dev.yml logs --tail=20
        fi
    fi
    
    # Show real-time logs if requested
    if [ "$1" = "--follow" ] || [ "$1" = "-f" ]; then
        echo "📋 Following logs (Ctrl+C to stop):"
        echo "==================================="
        if docker compose version &> /dev/null; then
            docker compose logs -f 2>/dev/null || docker compose -f docker-compose.dev.yml logs -f
        else
            docker-compose logs -f 2>/dev/null || docker-compose -f docker-compose.dev.yml logs -f
        fi
    fi
    
else
    echo "❌ No DoTask containers are running"
    echo ""
    echo "🚀 To start DoTask:"
    echo "  Development: ./scripts/docker-dev.sh"
    echo "  Production:  ./scripts/docker-prod.sh"
fi

echo ""
echo "💡 Available commands:"
echo "  ./scripts/docker-status.sh --logs    # Show recent logs"
echo "  ./scripts/docker-status.sh --follow  # Follow logs in real-time"
echo "  ./scripts/docker-cleanup.sh          # Clean up containers"
echo "  ./scripts/docker-cleanup.sh --full   # Full cleanup (removes data)"
