#!/bin/bash

# Script to stop all running services

set -e

echo "🛑 Stopping all PsyConnect services..."
echo ""

# Colors
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

# Function to stop a service by PID file
stop_service() {
    local service_name=$1
    local log_name=$(echo "$service_name" | tr '[:upper:] ' '[:lower:]_')
    local pid_file="../logs/latest/$log_name/pid"

    if [ -f "$pid_file" ]; then
        local pid=$(cat "$pid_file")
        if ps -p $pid > /dev/null 2>&1; then
            echo -e "${YELLOW}Stopping $service_name (PID: $pid)...${NC}"
            kill $pid 2>/dev/null || true
            sleep 1
            # Force kill if still running
            if ps -p $pid > /dev/null 2>&1; then
                kill -9 $pid 2>/dev/null || true
            fi
            echo -e "${GREEN}✅ $service_name stopped${NC}"
        else
            echo -e "${YELLOW}⚠️  $service_name not running${NC}"
        fi
        rm -f "$pid_file"
    else
        echo -e "${YELLOW}⚠️  No PID file for $service_name${NC}"
    fi
}

# Stop all services
echo "Stopping microservices..."
stop_service "API Gateway"
stop_service "Identity Service"
stop_service "Profile Service"
stop_service "Consultation Service"
stop_service "Chat Service"
stop_service "Notification Service"
# stop_service "Logging Service"

echo ""
echo "Stopping dependencies..."
cd ../dev
docker-compose down
cd ../scripts

echo ""
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All services stopped${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
