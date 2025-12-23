#!/bin/bash

# Test script for Notification Service (Go)

set -e

BLUE='\033[0;34m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8082"

echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${BLUE}🧪 Testing Notification Service (Go)${NC}"
echo -e "${BLUE}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""

# Test 1: Health Check
echo -e "${YELLOW}Test 1: Health Check${NC}"
response=$(curl -s ${BASE_URL}/health)
if echo "$response" | grep -q "healthy"; then
    echo -e "${GREEN}✅ Health check passed${NC}"
else
    echo -e "${RED}❌ Health check failed${NC}"
    exit 1
fi
echo ""

# Test 2: Ready Check
echo -e "${YELLOW}Test 2: Ready Check${NC}"
response=$(curl -s ${BASE_URL}/ready)
if echo "$response" | grep -q "ready"; then
    echo -e "${GREEN}✅ Ready check passed${NC}"
else
    echo -e "${RED}❌ Ready check failed${NC}"
    exit 1
fi
echo ""

# Test 3: Send Activation Email
echo -e "${YELLOW}Test 3: Send Activation Email${NC}"
response=$(curl -s -X POST ${BASE_URL}/mail/activate \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "code": "12345",
    "fullname": "Test User"
  }')

if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✅ Activation email test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Activation email test: ${response}${NC}"
fi
echo ""

# Test 4: Send Account Update Email
echo -e "${YELLOW}Test 4: Send Account Update Email${NC}"
response=$(curl -s -X POST ${BASE_URL}/mail/account-update \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com"
  }')

if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✅ Account update email test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Account update email test: ${response}${NC}"
fi
echo ""

# Test 5: Save FCM Token (requires database)
echo -e "${YELLOW}Test 5: Save FCM Token${NC}"
response=$(curl -s -X POST ${BASE_URL}/notification/token \
  -H "Content-Type: application/json" \
  -d '{
    "userId": "test-user-123",
    "token": "fcm-test-token-xyz"
  }')

if echo "$response" | grep -q "success"; then
    echo -e "${GREEN}✅ FCM token save test passed${NC}"
else
    echo -e "${YELLOW}⚠️  FCM token save test (requires database): ${response}${NC}"
fi
echo ""

# Test 6: Get Notifications (requires database)
echo -e "${YELLOW}Test 6: Get Notifications${NC}"
response=$(curl -s ${BASE_URL}/notification/test-user-123?limit=10)

if echo "$response" | grep -q "notifications"; then
    echo -e "${GREEN}✅ Get notifications test passed${NC}"
else
    echo -e "${YELLOW}⚠️  Get notifications test (requires database): ${response}${NC}"
fi
echo ""

echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo -e "${GREEN}✅ All tests completed!${NC}"
echo -e "${GREEN}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
echo ""
echo -e "${BLUE}Note:${NC} Some tests require database and Firebase to be configured."
echo -e "${BLUE}Email tests require valid SMTP credentials in .env${NC}"
echo ""
