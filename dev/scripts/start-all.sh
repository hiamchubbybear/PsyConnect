#!/bin/bash
set -e

echo "🚀 Starting all development dependencies..."
echo ""

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  No .env file found. Creating from .env.example..."
    cp .env.example .env
    echo "✅ Created .env file. Please review and update if needed."
    echo ""
fi

# Start services
docker-compose up -d

echo ""
echo "⏳ Waiting for services to be healthy..."
sleep 5

# Check status
docker-compose ps

echo ""
echo "✅ All services started!"
echo ""
echo "📋 Service URLs:"
echo "  MySQL:    localhost:3306 (user: root, password: see .env)"
echo "  Neo4j:    http://localhost:7474 (bolt://localhost:7687)"
echo "  MongoDB:  localhost:27017 (user: root, password: see .env)"
echo "  Redis:    localhost:6379 (password: see .env)"
echo "  Kafka:    localhost:9092"
echo ""
echo "🔍 Useful commands:"
echo "  Check status:  docker-compose ps"
echo "  View logs:     docker-compose logs -f [service]"
echo "  Stop all:      docker-compose down"
echo "  Clean volumes: docker-compose down -v"
