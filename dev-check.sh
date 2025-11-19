#!/bin/bash

# Script để kiểm tra và switch config giữa dev và cicd

set -e

SERVICES_DIR="services"
MODE="${1:-check}"  # check, dev, cicd

echo "🔍 Checking application configs..."
echo ""

# Function to check active profile
check_service_config() {
    local service=$1
    local config_path=$2

    echo "📦 $service:"

    if [ -f "$config_path" ]; then
        # Check for active profile
        if grep -q "spring.profiles.active" "$config_path" 2>/dev/null; then
            active_profile=$(grep "spring.profiles.active" "$config_path" | cut -d'=' -f2 | tr -d ' ')
            echo "   ✅ Active profile: $active_profile"
        elif grep -q "active:" "$config_path" 2>/dev/null; then
            active_profile=$(grep -A1 "profiles:" "$config_path" | grep "active:" | cut -d':' -f2 | tr -d ' ')
            echo "   ✅ Active profile: $active_profile"
        else
            echo "   ⚠️  No active profile set"
        fi

        # Show key configs
        echo "   📄 Config file: $config_path"

        # Check for common dev indicators
        if grep -q "localhost" "$config_path" 2>/dev/null; then
            echo "   🏠 Uses localhost"
        fi

        if grep -q "docker" "$config_path" 2>/dev/null || grep -q "host.docker.internal" "$config_path" 2>/dev/null; then
            echo "   🐳 Uses Docker hosts"
        fi
    else
        echo "   ❌ Config not found: $config_path"
    fi
    echo ""
}

# Function to set profile
set_profile() {
    local config_path=$1
    local profile=$2

    if [ -f "$config_path" ]; then
        # For .properties files
        if [[ "$config_path" == *.properties ]]; then
            if grep -q "spring.profiles.active" "$config_path"; then
                sed -i.bak "s/spring.profiles.active=.*/spring.profiles.active=$profile/" "$config_path"
                echo "   ✅ Updated to profile: $profile"
            else
                echo "spring.profiles.active=$profile" >> "$config_path"
                echo "   ✅ Added profile: $profile"
            fi
        # For .yml files
        elif [[ "$config_path" == *.yml ]] || [[ "$config_path" == *.yaml ]]; then
            if grep -q "active:" "$config_path"; then
                sed -i.bak "s/active:.*/active: $profile/" "$config_path"
                echo "   ✅ Updated to profile: $profile"
            else
                echo "spring:" >> "$config_path"
                echo "  profiles:" >> "$config_path"
                echo "    active: $profile" >> "$config_path"
                echo "   ✅ Added profile: $profile"
            fi
        fi
    fi
}

# Check Java services
echo "═══════════════════════════════════════"
echo "🔧 Java/Spring Boot Services"
echo "═══════════════════════════════════════"

# API Gateway
if [ -d "$SERVICES_DIR/apigateway" ]; then
    check_service_config "API Gateway" "$SERVICES_DIR/apigateway/src/main/resources/application.yml"
    if [ "$MODE" != "check" ]; then
        set_profile "$SERVICES_DIR/apigateway/src/main/resources/application.yml" "$MODE"
    fi
fi

# Identity Service
if [ -d "$SERVICES_DIR/identityservice" ]; then
    check_service_config "Identity Service" "$SERVICES_DIR/identityservice/src/main/resources/application.properties"
    if [ "$MODE" != "check" ]; then
        set_profile "$SERVICES_DIR/identityservice/src/main/resources/application.properties" "$MODE"
    fi
fi

# Profile Service
if [ -d "$SERVICES_DIR/profileservice" ]; then
    check_service_config "Profile Service" "$SERVICES_DIR/profileservice/src/main/application.properties"
    if [ "$MODE" != "check" ]; then
        set_profile "$SERVICES_DIR/profileservice/src/main/application.properties" "$MODE"
    fi
fi

# Check Go services
echo "═══════════════════════════════════════"
echo "🐹 Go Services"
echo "═══════════════════════════════════════"

for service in consultationservice chatservice loggingservice; do
    if [ -d "$SERVICES_DIR/$service" ]; then
        echo "📦 $service:"

        # Check for .env files
        if [ -f "$SERVICES_DIR/$service/.env" ]; then
            echo "   ✅ Has .env file"
            echo "   📄 Config: $SERVICES_DIR/$service/.env"
        elif [ -f "$SERVICES_DIR/$service/.env.example" ]; then
            echo "   ⚠️  Only .env.example found"
            if [ "$MODE" == "dev" ]; then
                cp "$SERVICES_DIR/$service/.env.example" "$SERVICES_DIR/$service/.env"
                echo "   ✅ Created .env from .env.example"
            fi
        else
            echo "   ❌ No .env file"
        fi
        echo ""
    fi
done

# Check Node.js services
echo "═══════════════════════════════════════"
echo "📦 Node.js Services"
echo "═══════════════════════════════════════"

if [ -d "$SERVICES_DIR/notificationservice" ]; then
    echo "📦 Notification Service:"
    if [ -f "$SERVICES_DIR/notificationservice/.env" ]; then
        echo "   ✅ Has .env file"
    else
        echo "   ⚠️  No .env file"
    fi
    echo ""
fi

echo "═══════════════════════════════════════"
echo "📋 Summary"
echo "═══════════════════════════════════════"
echo ""
echo "Usage:"
echo "  ./dev-check.sh          # Check current configs"
echo "  ./dev-check.sh dev      # Switch all to dev profile"
echo "  ./dev-check.sh cicd     # Switch all to cicd profile"
echo ""
echo "💡 Tips:"
echo "  - Dev profile: Uses localhost, local databases"
echo "  - CICD profile: Uses Docker hostnames, containerized services"
echo ""
