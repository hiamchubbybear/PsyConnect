#!/bin/bash

# Script to view logs for a specific service

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

SERVICE=$1

if [ -z "$SERVICE" ]; then
    echo "Usage: ./view-service-log.sh <service-name>"
    echo ""
    echo "Available services:"
    if [ -d "../logs/latest" ]; then
        ls -1 ../logs/latest/ | while read dir; do
            echo "  - $dir"
        done
    else
        echo "  No services running yet"
    fi
    exit 1
fi

SERVICE_LOG="../logs/latest/$SERVICE/service.log"

if [ ! -f "$SERVICE_LOG" ]; then
    echo -e "${YELLOW}⚠️  Log file not found: $SERVICE_LOG${NC}"
    echo ""
    echo "Available services:"
    ls -1 ../logs/latest/ 2>/dev/null | while read dir; do
        echo "  - $dir"
    done
    exit 1
fi

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}📊 Viewing logs for: $SERVICE${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""

# Show log file path and PID
echo "Log file: $SERVICE_LOG"
if [ -f "../logs/latest/$SERVICE/pid" ]; then
    PID=$(cat "../logs/latest/$SERVICE/pid")
    if ps -p $PID > /dev/null 2>&1; then
        echo -e "Status: ${GREEN}Running${NC} (PID: $PID)"
    else
        echo -e "Status: ${YELLOW}Stopped${NC}"
    fi
fi
echo ""

# Tail the log
tail -f "$SERVICE_LOG"
