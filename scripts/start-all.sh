#!/bin/bash

# Script to start all dependencies and microservices

set -e

echo "🚀 Starting PsyConnect Development Environment"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Step 1: Start dependencies
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📦 Step 1/3: Starting Dependencies${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
cd ../dev
./scripts/start-all.sh

# Step 2: Switch to dev profile
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🔧 Step 2/3: Switching to Dev Profile${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
cd ..
./scripts/dev-check.sh dev > /dev/null 2>&1
echo -e "${GREEN}✅ Switched to dev profile${NC}"

# Step 3: Start all microservices
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🚀 Step 3/3: Starting Microservices${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Create timestamped logs directory for this run
TIMESTAMP=$(date '+%Y-%m-%d_%H-%M-%S')
LOGS_DIR="logs/${TIMESTAMP}"
mkdir -p "$LOGS_DIR"

# Create symlink to latest logs
rm -f logs/latest
ln -s "${TIMESTAMP}" logs/latest

echo "📝 Logs will be saved to: $LOGS_DIR"
echo "   Shortcut: logs/latest/"
echo -e "${GREEN}✅ Logs directory ready${NC}"
echo ""

# Function to start a service
start_service() {
    local service_name=$1
    local service_dir=$2
    local start_command=$3
    local port=$4
    local log_name=$(echo "$service_name" | tr '[:upper:] ' '[:lower:]_')

    echo -e "${YELLOW}Starting $service_name on port $port...${NC}"

    # Create service-specific log directory
    mkdir -p "$LOGS_DIR/$log_name"

    # Start from service directory to ensure .env and config files are found
    (
        cd "$service_dir"
        nohup bash -c "$start_command" > "../../$LOGS_DIR/$log_name/service.log" 2>&1 &
        echo $! > "../../$LOGS_DIR/$log_name/pid"
    )

    # Get PID from file
    local pid=$(cat "$LOGS_DIR/$log_name/pid")
    echo -e "${GREEN}✅ $service_name started (PID: $pid)${NC}"
    echo "   Log: $LOGS_DIR/$log_name/service.log"
}

# Start Java services
echo -e "${BLUE}☕ Starting Java Services...${NC}"
start_service "API Gateway" "services/apigateway" "mvn spring-boot:run -DskipTests -Dspring-boot.run.profiles=dev" "8888"
start_service "Identity Service" "services/identityservice" "mvn spring-boot:run -DskipTests -Dspring-boot.run.profiles=dev" "8080"
start_service "Profile Service" "services/profileservice" "mvn spring-boot:run -DskipTests -Dspring-boot.run.profiles=dev" "8081"

echo ""
echo -e "${BLUE}🐹 Starting Go Services...${NC}"
start_service "Consultation Service" "services/consultationservice" "go run cmd/main.go" "8084"
start_service "Chat Service" "services/chatservice" "go run cmd/main.go" "8083"
# start_service "Logging Service" "services/loggingservice" "go run main.go" "8086"

echo ""
echo -e "${BLUE}📦 Starting Node.js Services...${NC}"
start_service "Notification Service" "services/notificationservice" "node main.js" "8082"

# Summary
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
echo ""
echo "📊 Database Services:"
echo "  🗄️  MySQL:     localhost:3306"
echo "  🔵 Neo4j:     http://localhost:7474"
echo "  🍃 MongoDB:   localhost:27017"
echo "  🔴 Redis:     localhost:6379"
echo "  📨 Kafka:     localhost:9092"
echo ""
echo "📝 Logs: ./logs/*.log"
echo "🛑 Stop all: ./stop-all.sh"
echo ""
echo -e "${YELLOW}⏳ Services are starting... Check logs for details.${NC}"
echo -e "${YELLOW}   It may take 1-2 minutes for all services to be fully ready.${NC}"
