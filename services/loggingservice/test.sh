#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

echo -e "${GREEN}🧪 Testing PsyConnect Logging Service${NC}"
echo ""

# Check if logging service is running
echo -e "${YELLOW}1. Checking if logging service is running...${NC}"
if curl -s http://localhost:8080/health > /dev/null; then
    echo -e "${GREEN}✅ Logging service is running${NC}"
    curl -s http://localhost:8080/health | jq .
else
    echo -e "${RED}❌ Logging service is not running${NC}"
    echo "Please start it with: go run main.go"
    exit 1
fi

echo ""

# Check if Kafka is accessible
echo -e "${YELLOW}2. Checking Kafka connectivity...${NC}"
if nc -z localhost 9092 2>/dev/null; then
    echo -e "${GREEN}✅ Kafka is accessible${NC}"
else
    echo -e "${RED}❌ Kafka is not accessible${NC}"
    echo "Please ensure Kafka is running on localhost:9092"
fi

echo ""

# Check if Loki is accessible
echo -e "${YELLOW}3. Checking Loki connectivity...${NC}"
if curl -s http://localhost:3100/ready > /dev/null; then
    echo -e "${GREEN}✅ Loki is accessible${NC}"
else
    echo -e "${RED}❌ Loki is not accessible${NC}"
    echo "Please start the observability stack: docker-compose up -d"
fi

echo ""

# Check if Grafana is accessible
echo -e "${YELLOW}4. Checking Grafana connectivity...${NC}"
if curl -s http://localhost:3000/api/health > /dev/null; then
    echo -e "${GREEN}✅ Grafana is accessible${NC}"
    echo "   URL: http://localhost:3000"
    echo "   Username: admin"
    echo "   Password: 16032004"
else
    echo -e "${RED}❌ Grafana is not accessible${NC}"
    echo "Please start the observability stack: docker-compose up -d"
fi

echo ""

# Send test log event
echo -e "${YELLOW}5. Sending test log event to Kafka...${NC}"

# Check if kafkacat is installed
if command -v kafkacat &> /dev/null; then
    TEST_LOG='{
        "timestamp": "'$(date -u +"%Y-%m-%dT%H:%M:%SZ")'",
        "level": "INFO",
        "service": "test-service",
        "message": "Test log from test script",
        "userId": "test-user-123",
        "action": "TEST",
        "traceId": "test-trace-'$(date +%s)'",
        "metadata": {
            "test": true,
            "script": "test.sh"
        }
    }'

    echo "$TEST_LOG" | kafkacat -P -b localhost:9092 -t logging-service
    echo -e "${GREEN}✅ Test log sent to Kafka${NC}"
    echo "Check Grafana to see if it appears!"
else
    echo -e "${YELLOW}⚠️  kafkacat not installed${NC}"
    echo "Install with: brew install kafkacat (macOS) or apt-get install kafkacat (Linux)"
    echo ""
    echo "Alternatively, you can manually send a log event using Kafka console producer"
fi

echo ""

# Check metrics
echo -e "${YELLOW}6. Checking service metrics...${NC}"
curl -s http://localhost:8080/metrics | jq .

echo ""
echo -e "${GREEN}✅ Test complete!${NC}"
echo ""
echo "Next steps:"
echo "  1. Open Grafana: http://localhost:3000"
echo "  2. Navigate to Dashboards → PsyConnect → Logs Overview"
echo "  3. You should see the test log event"
echo ""
