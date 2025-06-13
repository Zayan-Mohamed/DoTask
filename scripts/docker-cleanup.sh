#!/bin/bash
# DoTask Docker Cleanup Script

echo "🧹 Cleaning up DoTask Docker resources..."

# Stop all containers
echo "⏹️ Stopping containers..."
if docker compose version &> /dev/null; then
    COMPOSE_CMD="docker compose"
else
    COMPOSE_CMD="docker-compose"
fi

$COMPOSE_CMD down --remove-orphans 2>/dev/null || true
$COMPOSE_CMD -f docker-compose.dev.yml down --remove-orphans 2>/dev/null || true

# Remove containers, networks, and volumes (optional)
if [ "$1" = "--full" ]; then
    echo "🗑️ Performing full cleanup (including volumes)..."
    $COMPOSE_CMD down --volumes --remove-orphans 2>/dev/null || true
    $COMPOSE_CMD -f docker-compose.dev.yml down --volumes --remove-orphans 2>/dev/null || true
    
    # Remove unused Docker resources
    docker system prune -f --volumes
    
    echo "✅ Full cleanup completed!"
    echo "⚠️ All data has been removed. You'll need to recreate the database."
else
    echo "🔄 Performing standard cleanup (preserving volumes)..."
    
    # Remove unused Docker resources but keep volumes
    docker system prune -f
    
    echo "✅ Standard cleanup completed!"
    echo "💾 Database volumes preserved."
    echo "💡 Use '--full' flag to remove all data."
fi

echo ""
echo "📊 Current Docker status:"
docker ps -a --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"
