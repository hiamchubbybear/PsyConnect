#!/bin/bash

# Script to view all service logs in real-time

set -e

# Colors
GREEN='\033[0;32m'
BLUE='\033[0;34m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}📊 Viewing All Service Logs${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

cd ..

# Check if logs directory exists
if [ ! -d "logs/latest" ]; then
    echo -e "${YELLOW}⚠️  No logs directory found. Services not started yet?${NC}"
    exit 1
fi

# Count log files
log_count=$(find logs/latest -name "service.log" 2>/dev/null | wc -l)
if [ "$log_count" -eq 0 ]; then
    echo -e "${YELLOW}⚠️  No log files found. Services not started yet?${NC}"
    exit 1
fi

echo -e "${GREEN}Streaming logs from all services...${NC}"
echo -e "${YELLOW}Press Ctrl+C to stop${NC}"
echo ""
echo "Log directory: logs/latest/"
echo ""

# Tail all log files with service name prefix
tail -f logs/latest/*/service.log 2>/dev/null | while read line; do
    echo -e "${BLUE}[$(date '+%H:%M:%S')]${NC} $line"
done
