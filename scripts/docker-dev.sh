#!/bin/bash
# DoTask Development Environment Starter

echo "🚀 Starting DoTask Development Environment..."

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker is not running. Please start Docker and try again."
    exit 1
fi

# Build and start development containers
echo "📦 Building and starting development containers..."
if docker compose version &> /dev/null; then
    docker compose -f docker-compose.dev.yml up --build
else
    docker-compose -f docker-compose.dev.yml up --build
fi

echo "✅ Development environment is running!"
echo "🌐 Frontend: http://localhost:5173"
echo "⚡ Backend: http://localhost:8080"
echo "🗄️ Database: localhost:5432"
echo "🎮 GraphQL Playground: http://localhost:8080"
