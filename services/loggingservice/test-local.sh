#!/bin/bash

# Colors
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}   🧪 PsyConnect Logging Service - Local Test Suite${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test counter
PASSED=0
FAILED=0

# Function to check service
check_service() {
    local name=$1
    local url=$2
    local port=$3

    echo -e "${YELLOW}Testing $name...${NC}"

    if curl -s --max-time 5 "$url" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ $name is running on port $port${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ $name is NOT accessible on port $port${NC}"
        ((FAILED++))
        return 1
    fi
}

# Function to check port
check_port() {
    local name=$1
    local port=$2

    echo -e "${YELLOW}Checking port $port ($name)...${NC}"

    if nc -z localhost $port 2>/dev/null; then
        echo -e "${GREEN}✅ Port $port is open ($name)${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ Port $port is closed ($name)${NC}"
        ((FAILED++))
        return 1
    fi
}

echo -e "${BLUE}━━━ 1. Infrastructure Services ━━━${NC}"
echo ""

# Check Kafka
check_port "Kafka" 9092

# Check Zookeeper
check_port "Zookeeper" 2181

# Check MySQL
check_port "MySQL" 3306

# Check Neo4j
check_port "Neo4j" 7687

# Check MongoDB
check_port "MongoDB" 27017

# Check Redis
check_port "Redis" 6379

echo ""
echo -e "${BLUE}━━━ 2. Observability Stack ━━━${NC}"
echo ""

# Check Loki
check_service "Loki" "http://localhost:3100/ready" 3100

# Check Grafana
check_service "Grafana" "http://localhost:3000/api/health" 3000

# Check Prometheus
check_service "Prometheus" "http://localhost:9090/-/ready" 9090

# Check Tempo
check_service "Tempo" "http://localhost:3200/ready" 3200

echo ""
echo -e "${BLUE}━━━ 3. Logging Service ━━━${NC}"
echo ""

# Check if logging service is running
if pgrep -f "logging-service" > /dev/null || pgrep -f "go run main.go" > /dev/null; then
    echo -e "${GREEN}✅ Logging service process is running${NC}"
    ((PASSED++))

    # Check health endpoint
    check_service "Logging Service Health" "http://localhost:8080/health" 8080

    # Check metrics endpoint
    if curl -s http://localhost:8080/metrics > /dev/null 2>&1; then
        echo -e "${GREEN}✅ Metrics endpoint is accessible${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ Metrics endpoint is NOT accessible${NC}"
        ((FAILED++))
    fi
else
    echo -e "${RED}❌ Logging service is NOT running${NC}"
    echo -e "${YELLOW}   Start it with: go run main.go${NC}"
    ((FAILED++))
fi

echo ""
echo -e "${BLUE}━━━ 4. Grafana Configuration ━━━${NC}"
echo ""

# Check Grafana datasources
if curl -s -u admin:16032004 http://localhost:3000/api/datasources 2>/dev/null | grep -q "Loki"; then
    echo -e "${GREEN}✅ Loki datasource is configured${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Loki datasource is NOT configured${NC}"
    ((FAILED++))
fi

if curl -s -u admin:16032004 http://localhost:3000/api/datasources 2>/dev/null | grep -q "Prometheus"; then
    echo -e "${GREEN}✅ Prometheus datasource is configured${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Prometheus datasource is NOT configured${NC}"
    ((FAILED++))
fi

# Check Grafana dashboards
DASHBOARD_COUNT=$(curl -s -u admin:16032004 http://localhost:3000/api/search?type=dash-db 2>/dev/null | grep -o "\"title\"" | wc -l | tr -d ' ')
if [ "$DASHBOARD_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✅ Found $DASHBOARD_COUNT dashboards in Grafana${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  No dashboards found (may need time to provision)${NC}"
fi

# Check alert rules
ALERT_COUNT=$(curl -s -u admin:16032004 http://localhost:3000/api/ruler/grafana/api/v1/rules 2>/dev/null | grep -o "\"title\"" | wc -l | tr -d ' ')
if [ "$ALERT_COUNT" -gt 0 ]; then
    echo -e "${GREEN}✅ Found $ALERT_COUNT alert rules configured${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  No alert rules found (may need time to provision)${NC}"
fi

echo ""
echo -e "${BLUE}━━━ 5. Send Test Logs ━━━${NC}"
echo ""

# Check if kafkacat is installed
if command -v kafkacat &> /dev/null; then
    echo -e "${YELLOW}Sending test log events...${NC}"

    # Send INFO log
    echo '{
        "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
        "level": "INFO",
        "service": "test-service",
        "message": "Test INFO log from local test",
        "userId": "test-user-local",
        "action": "TEST",
        "traceId": "test-trace-'$(date +%s)'"
    }' | kafkacat -P -b localhost:9092 -t logging-service 2>/dev/null

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ INFO log sent successfully${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ Failed to send INFO log${NC}"
        ((FAILED++))
    fi

    # Send ERROR log
    echo '{
        "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
        "level": "ERROR",
        "service": "test-service",
        "message": "Test ERROR log from local test",
        "error": "This is a test error",
        "userId": "test-user-local"
    }' | kafkacat -P -b localhost:9092 -t logging-service 2>/dev/null

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✅ ERROR log sent successfully${NC}"
        ((PASSED++))
    else
        echo -e "${RED}❌ Failed to send ERROR log${NC}"
        ((FAILED++))
    fi

    echo -e "${YELLOW}   Wait a few seconds and check Grafana for the logs${NC}"
else
    echo -e "${YELLOW}⚠️  kafkacat not installed, skipping log sending${NC}"
    echo -e "${YELLOW}   Install: brew install kafkacat (macOS)${NC}"
fi

echo ""
echo -e "${BLUE}━━━ 6. File System Checks ━━━${NC}"
echo ""

# Check log files
if [ -f "./logs/app.log" ]; then
    LOG_SIZE=$(du -h ./logs/app.log | cut -f1)
    echo -e "${GREEN}✅ app.log exists (size: $LOG_SIZE)${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  app.log not found (will be created when logs are received)${NC}"
fi

if [ -f "./logs/error.log" ]; then
    ERR_SIZE=$(du -h ./logs/error.log | cut -f1)
    echo -e "${GREEN}✅ error.log exists (size: $ERR_SIZE)${NC}"
    ((PASSED++))
else
    echo -e "${YELLOW}⚠️  error.log not found (will be created when errors are logged)${NC}"
fi

# Check configuration files
CONFIG_FILES=(
    "grafana/dashboards/system-overview.json"
    "grafana/dashboards/identity-service.json"
    "grafana/dashboards/profile-service.json"
    "grafana/dashboards/error-analysis.json"
    "grafana/dashboards/user-activity.json"
    "grafana/provisioning/alerting/alerts.yml"
    "grafana/provisioning/alerting/notifications.yml"
)

MISSING_FILES=0
for file in "${CONFIG_FILES[@]}"; do
    if [ ! -f "$file" ]; then
        echo -e "${RED}❌ Missing: $file${NC}"
        ((MISSING_FILES++))
    fi
done

if [ $MISSING_FILES -eq 0 ]; then
    echo -e "${GREEN}✅ All configuration files present${NC}"
    ((PASSED++))
else
    echo -e "${RED}❌ Missing $MISSING_FILES configuration files${NC}"
    ((FAILED++))
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}   📊 Test Results${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "   ${GREEN}Passed: $PASSED${NC}"
echo -e "   ${RED}Failed: $FAILED${NC}"
echo ""

TOTAL=$((PASSED + FAILED))
if [ $TOTAL -gt 0 ]; then
    SUCCESS_RATE=$(awk "BEGIN {printf \"%.1f\", ($PASSED/$TOTAL)*100}")
    echo -e "   Success Rate: ${GREEN}$SUCCESS_RATE%${NC}"
fi

echo ""
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}   🔗 Quick Links${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "   📊 Grafana:    ${BLUE}http://localhost:3000${NC} (admin/16032004)"
echo -e "   📈 Prometheus: ${BLUE}http://localhost:9090${NC}"
echo -e "   🔍 Loki:       ${BLUE}http://localhost:3100${NC}"
echo -e "   🏥 Health:     ${BLUE}http://localhost:8080/health${NC}"
echo -e "   📊 Metrics:    ${BLUE}http://localhost:8080/metrics${NC}"
echo ""

if [ $FAILED -gt 0 ]; then
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}   ⚠️  Some tests failed. Troubleshooting:${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "   1. Start observability stack: ${BLUE}docker-compose up -d${NC}"
    echo -e "   2. Check logs: ${BLUE}docker-compose logs -f${NC}"
    echo -e "   3. Start logging service: ${BLUE}go run main.go${NC}"
    echo -e "   4. Wait 30s for services to initialize"
    echo ""
else
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${GREEN}   ✅ All tests passed! System is healthy!${NC}"
    echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
fi

exit $FAILED
