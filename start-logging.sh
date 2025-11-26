#!/bin/bash

# Quick Start Script for PsyConnect Logging System

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}   🚀 PsyConnect Logging System - Quick Start${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check if running from correct directory
if [ ! -f "ALL_SERVICES_COMPLETE.md" ]; then
    echo -e "${RED}❌ Please run this script from the PsyConnect root directory${NC}"
    exit 1
fi

# Function to check if command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Check prerequisites
echo -e "${YELLOW}Checking prerequisites...${NC}"

if ! command_exists docker; then
    echo -e "${RED}❌ Docker is not installed${NC}"
    exit 1
fi

if ! command_exists docker-compose; then
    echo -e "${RED}❌ Docker Compose is not installed${NC}"
    exit 1
fi

echo -e "${GREEN}✅ All prerequisites met${NC}"
echo ""

# Ask for environment
echo -e "${YELLOW}Select environment:${NC}"
echo "  1) Development (localhost)"
echo "  2) CI/CD (Docker network)"
read -p "Enter choice [1-2]: " env_choice

case $env_choice in
    1)
        ENVIRONMENT="dev"
        echo -e "${GREEN}Selected: Development${NC}"
        ;;
    2)
        ENVIRONMENT="cicd"
        echo -e "${GREEN}Selected: CI/CD${NC}"
        ;;
    *)
        echo -e "${YELLOW}Invalid choice, using Development${NC}"
        ENVIRONMENT="dev"
        ;;
esac

echo ""

# Start observability stack
echo -e "${BLUE}━━━ Starting Observability Stack ━━━${NC}"
cd services/loggingservice

echo -e "${YELLOW}Starting Loki, Grafana, Prometheus, Tempo...${NC}"
ENVIRONMENT=$ENVIRONMENT docker-compose up -d

echo -e "${GREEN}✅ Observability stack started${NC}"
echo ""

# Wait for services to be ready
echo -e "${YELLOW}Waiting for services to be ready...${NC}"
sleep 10

# Check service health
echo -e "${BLUE}━━━ Checking Service Health ━━━${NC}"

check_service() {
    local name=$1
    local url=$2

    if curl -s "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $name is ready${NC}"
        return 0
    else
        echo -e "${RED}❌ $name is not ready${NC}"
        return 1
    fi
}

check_service "Loki" "http://localhost:3100/ready"
check_service "Grafana" "http://localhost:3000/api/health"
check_service "Prometheus" "http://localhost:9090/-/ready"

echo ""

# Display access information
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}   ✅ Logging System Started Successfully!${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Access URLs:${NC}"
echo -e "  📊 Grafana:    ${GREEN}http://localhost:3000${NC}"
echo -e "     Username:   ${BLUE}admin${NC}"
echo -e "     Password:   ${BLUE}16032004${NC}"
echo ""
echo -e "  📈 Prometheus: ${GREEN}http://localhost:9090${NC}"
echo -e "  🔍 Loki:       ${GREEN}http://localhost:3100${NC}"
echo -e "  🔗 Tempo:      ${GREEN}http://localhost:3200${NC}"
echo ""
echo -e "${YELLOW}Logging Service:${NC}"
echo -e "  🏥 Health:     ${GREEN}http://localhost:8080/health${NC}"
echo -e "  📊 Metrics:    ${GREEN}http://localhost:8080/metrics${NC}"
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${YELLOW}Next Steps:${NC}"
echo -e "  1. Open Grafana: ${GREEN}http://localhost:3000${NC}"
echo -e "  2. Go to Dashboards → PsyConnect"
echo -e "  3. Start your services to see logs!"
echo ""
echo -e "${YELLOW}To stop:${NC}"
echo -e "  cd services/loggingservice"
echo -e "  docker-compose down"
echo ""
echo -e "${GREEN}Happy Logging! 🎉${NC}"
echo ""
