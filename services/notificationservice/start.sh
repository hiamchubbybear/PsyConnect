#!/bin/bash

# Quick start script for notification service

cd "$(dirname "$0")"

echo "🚀 Starting Notification Service (Go)..."
echo ""

# Check if Kafka is running
if ! nc -z 127.0.0.1 19092 2>/dev/null; then
    echo "❌ Kafka is not running on 127.0.0.1:19092"
    echo "Please start Kafka first:"
    echo "  cd dev && docker-compose up -d kafka"
    exit 1
fi

echo "✅ Kafka is running"
echo ""

# Run service
go run cmd/main.go
