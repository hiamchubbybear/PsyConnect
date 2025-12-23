#!/bin/sh

# Entrypoint script to wait for Kafka before starting service

echo "⏳ Waiting for Kafka to be ready..."

# Wait for Kafka
KAFKA_HOST=$(echo $KAFKA_BROKERS | cut -d: -f1)
KAFKA_PORT=$(echo $KAFKA_BROKERS | cut -d: -f2)

echo "📡 Checking Kafka at $KAFKA_HOST:$KAFKA_PORT..."

# Wait up to 60 seconds
for i in $(seq 1 30); do
    if nc -z $KAFKA_HOST $KAFKA_PORT 2>/dev/null; then
        echo "✅ Kafka is ready!"
        break
    fi
    echo "⏳ Waiting for Kafka... ($i/30)"
    sleep 2
done

# Additional wait to ensure Kafka is fully ready
echo "⏳ Waiting additional 5 seconds for Kafka to be fully ready..."
sleep 5

echo "🚀 Starting notification service..."
exec ./notification-service
