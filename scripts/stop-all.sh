#!/bin/bash

set -e

echo "🛑 Stopping all PsyConnect services..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function: stop service by port
stop_by_port() {
    local service_name=$1
    local port=$2

    echo -e "${YELLOW}Checking $service_name on port $port...${NC}"

    # Tìm PID bằng port
    local pid=$(lsof -t -i :$port || true)

    if [ -n "$pid" ]; then
        echo -e "${YELLOW}Stopping $service_name (PID: $pid)...${NC}"
        kill $pid 2>/dev/null || true
        sleep 1

        # Force kill nếu còn chạy
        if ps -p $pid > /dev/null 2>&1; then
            kill -9 $pid 2>/dev/null || true
        fi

        echo -e "${GREEN}✅ $service_name stopped${NC}"
    else
        echo -e "${YELLOW}⚠️  $service_name not running on port $port${NC}"
    fi
}

echo "Stopping microservices..."

stop_by_port "API Gateway"         8888
stop_by_port "Identity Service"    8080
stop_by_port "Profile Service"     8081
stop_by_port "Notification Service" 8082
stop_by_port "Chat Service"         8083
stop_by_port "Consultation Service" 8084

echo ""
echo "Stopping dependencies (Docker)..."
cd ../dev
docker-compose down
cd ../scripts

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All services stopped${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
