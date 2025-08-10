#!/bin/bash

# Dừng script nếu có lỗi
set -e

# Chạy API server
echo " Starting API server..."
go run cmd/main.go &
API_PID=$!

# Chạy WebSocket server
echo " Starting WebSocket server..."
go run internal/ws/cmd/main.go &
WS_PID=$!

# Chờ cả 2 process
wait $API_PID $WS_PID
