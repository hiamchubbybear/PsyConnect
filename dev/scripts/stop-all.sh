#!/bin/bash
set -e

echo "🛑 Stopping all development dependencies..."
docker-compose down

echo "✅ All services stopped!"
