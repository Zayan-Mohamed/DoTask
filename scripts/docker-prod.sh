#!/bin/bash# Build and s# Show conta# Show logs for a few seconds
echo ""
echo "📋 Recent logs:"
sleep 2
if docker compose version &> /dev/null; then
    docker compose logs --tail=10
else
    docker-compose logs --tail=10
fistatus
echo ""
echo "📊 Container Status:"
if docker compose version &> /dev/null; then
    docker compose ps
else
    docker-compose ps
fi production containers
echo "📦 Building and starting production containers..."
if docker compose version &> /dev/null; then
    docker compose up --build -d
else
    docker-compose up --build -d
fiDoTask Production Environment Starter

echo "🚀 Starting DoTask Production Environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Build and start production containers
echo "📦 Building and starting production containers..."
docker compose up --build -d

echo "✅ Production environment is running!"
echo "🌐 Application: http://localhost"
echo "📊 Health checks enabled"

# Show container status
echo ""
echo "📊 Container Status:"
docker compose ps

# Show logs for a few seconds
echo ""
echo "📋 Recent logs:"
sleep 2
docker compose logs --tail=10
