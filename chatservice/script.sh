#!/bin/bash
set -e
echo " Starting API server..."
go run cmd/main.go &
API_PID=$!
echo " Starting WebSocket server..."
go run internal/ws/cmd/main.go &
WS_PID=$!
wait $API_PID $WS_PID
