#!/bin/bash

echo "Cleaning old Docker resources..."

# Stop all containers
docker-compose down --volumes --remove-orphans

# Prune unused containers, networks
docker system prune -f

# Prune unused images
docker image prune -a -f

# Remove buildx cache
docker buildx prune --all --force

# Optional: clean old build cache directories
rm -rf ~/.docker/buildx

echo "Cleanup completed."

echo "Starting PsyConnect services..."

# 1) Database layer first
echo "Starting MySQL, MongoDB, Redis, Neo4j..."
docker-compose up -d db mongodb redis neo4j

echo "Waiting 30 seconds for DB layer to initialize..."
sleep 30

# 2) Start Core Microservices (Profile + Notification)
echo "Starting ProfileService and NotificationService..."
docker-compose up -d profileservice notificationservice

echo "Waiting 10 seconds..."
sleep 10

# 3) IdentityService
echo "Starting IdentityService..."
docker-compose up -d identityservice

echo "Waiting 10 seconds..."
sleep 10

# 4) API Gateway
echo "Starting APIGateway..."
docker-compose up -d apigateway

echo "Waiting 5 seconds..."
sleep 5

# 5) Remaining services
echo "Starting ConsultationService and Dozzle..."
docker-compose up -d consultationservice dozzle

echo "All services started successfully."
