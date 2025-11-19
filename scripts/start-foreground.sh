#!/bin/bash

# Script to start all services in FOREGROUND mode
# Logs will be visible in terminal (not background)

set -e

echo "🚀 Starting PsyConnect - FOREGROUND MODE"
echo ""

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

# Step 1: Start dependencies
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📦 Starting Dependencies${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
cd ../dev
./scripts/start-all.sh

# Step 2: Switch to dev profile
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🔧 Switching to Dev Profile${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
cd ..
./scripts/dev-check.sh dev > /dev/null 2>&1
echo -e "${GREEN}✅ Switched to dev profile${NC}"

# Step 3: Start services in foreground with tmux
echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🚀 Starting Microservices (Foreground)${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Check if tmux is installed
if ! command -v tmux &> /dev/null; then
    echo -e "${RED}❌ tmux is not installed!${NC}"
    echo "Install with: brew install tmux"
    exit 1
fi

# Kill existing session if exists
tmux kill-session -t psyconnect 2>/dev/null || true

# Create new tmux session
echo -e "${GREEN}Creating tmux session 'psyconnect'...${NC}"
echo ""

# Start tmux session with first service
tmux new-session -d -s psyconnect -n "API-Gateway" "cd services/apigateway && mvn spring-boot:run -DskipTests; read"

# Split and add more services
tmux split-window -h -t psyconnect "cd services/identityservice && mvn spring-boot:run -DskipTests; read"
tmux split-window -v -t psyconnect "cd services/profileservice && mvn spring-boot:run -DskipTests; read"
tmux select-pane -t 0
tmux split-window -v -t psyconnect "cd services/chatservice && go run cmd/main.go; read"

# Create new window for more services
tmux new-window -t psyconnect -n "Services-2"
tmux send-keys -t psyconnect:1 "cd services/consultationservice && go run cmd/main.go; read" C-m
tmux split-window -h -t psyconnect "cd services/notificationservice && node main.js; read"

# Select first window
tmux select-window -t psyconnect:0

echo ""
echo -e "${GREEN}✅ Services started in tmux!${NC}"
echo ""
echo "📋 tmux Commands:"
echo "  Attach:        tmux attach -t psyconnect"
echo "  Detach:        Ctrl+B then D"
echo "  Switch pane:   Ctrl+B then arrow keys"
echo "  Kill session:  tmux kill-session -t psyconnect"
echo ""
echo -e "${YELLOW}Attaching to tmux session...${NC}"

# Attach to session
tmux attach -t psyconnect
