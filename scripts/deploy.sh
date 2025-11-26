#!/bin/bash

# Deploy script for PsyConnect services with logging infrastructure

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🚀 PsyConnect Deployment Script${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Change to deploy directory
cd "$(dirname "$0")/../deploy"

echo -e "${YELLOW}Cleaning old Docker resources...${NC}"

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

echo -e "${GREEN}✅ Cleanup completed${NC}"
echo ""

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🚀 Starting PsyConnect Services${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# 1) Infrastructure layer (Kafka, Zookeeper)
echo -e "${YELLOW}📨 Starting Kafka infrastructure...${NC}"
docker-compose up -d zookeeper kafka
echo -e "${GREEN}✅ Kafka infrastructure started${NC}"

echo "⏳ Waiting 30 seconds for Kafka to initialize..."
sleep 30

# 2) Database layer
echo ""
echo -e "${YELLOW}🗄️  Starting database layer...${NC}"
docker-compose up -d db mongodb redis neo4j
echo -e "${GREEN}✅ Databases started${NC}"

echo "⏳ Waiting 30 seconds for databases to initialize..."
sleep 30

# 3) Observability stack (Loki, Prometheus, Grafana)
echo ""
echo -e "${YELLOW}📊 Starting observability stack...${NC}"
docker-compose up -d loki prometheus grafana
echo -e "${GREEN}✅ Observability stack started${NC}"

echo "⏳ Waiting 20 seconds for observability stack..."
sleep 20

# 4) Logging Service
echo ""
echo -e "${YELLOW}📝 Starting Logging Service...${NC}"
docker-compose up -d logging-service
echo -e "${GREEN}✅ Logging Service started${NC}"

echo "⏳ Waiting 10 seconds..."
sleep 10

# 5) Core Microservices (Profile + Notification)
echo ""
echo -e "${YELLOW}🔧 Starting core microservices...${NC}"
docker-compose up -d profileservice notificationservice
echo -e "${GREEN}✅ Core microservices started${NC}"

echo "⏳ Waiting 15 seconds..."
sleep 15

# 6) IdentityService
echo ""
echo -e "${YELLOW}🔐 Starting Identity Service...${NC}"
docker-compose up -d identityservice
echo -e "${GREEN}✅ Identity Service started${NC}"

echo "⏳ Waiting 15 seconds..."
sleep 15

# 7) API Gateway
echo ""
echo -e "${YELLOW}🌐 Starting API Gateway...${NC}"
docker-compose up -d apigateway
echo -e "${GREEN}✅ API Gateway started${NC}"

echo "⏳ Waiting 10 seconds..."
sleep 10

# 8) Remaining services (Consultation, Chat)
echo ""
echo -e "${YELLOW}💬 Starting Consultation and Chat services...${NC}"
docker-compose up -d consultationservice chatservice
echo -e "${GREEN}✅ Consultation and Chat services started${NC}"

echo "⏳ Waiting 10 seconds..."
sleep 10

# 9) Monitoring tools
echo ""
echo -e "${YELLOW}🔍 Starting monitoring tools...${NC}"
docker-compose up -d kafdrop dozzle
echo -e "${GREEN}✅ Monitoring tools started${NC}"

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All Services Started Successfully!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

echo "📋 Service URLs:"
echo "  🌐 API Gateway:        http://localhost:8888"
echo "  🔐 Identity Service:   http://localhost:8080"
echo "  👤 Profile Service:    http://localhost:8081"
echo "  📧 Notification:       http://localhost:8082"
echo "  💬 Chat Service:       http://localhost:8083"
echo "  🏥 Consultation:       http://localhost:8084"
echo "  📊 Logging Service:    http://localhost:8085"
echo ""
echo "📊 Observability Stack:"
echo "  📈 Grafana:            http://localhost:3000 (admin/16032004)"
echo "  📉 Prometheus:         http://localhost:9090"
echo "  🔍 Loki:               http://localhost:3100"
echo ""
echo "🔍 Monitoring Tools:"
echo "  📨 Kafdrop:            http://localhost:9000"
echo "  📋 Dozzle:             http://localhost:9999"
echo ""
echo "📊 Databases:"
echo "  🗄️  MySQL:     localhost:3306"
echo "  🔵 Neo4j:     http://localhost:7474"
echo "  🍃 MongoDB:   localhost:27017"
echo "  🔴 Redis:     localhost:6379"
echo "  📨 Kafka:     localhost:9092"
echo ""
echo -e "${YELLOW}⏳ Services are starting... It may take 2-3 minutes for all services to be fully ready.${NC}"
echo -e "${YELLOW}   Check logs with: docker-compose logs -f [service-name]${NC}"
echo ""
