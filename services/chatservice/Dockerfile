# Stage 1: Build
FROM golang:1.24-alpine AS builder
RUN apk add --no-cache git

# Context root
WORKDIR /app

# Copy go.mod and go.sum
COPY services/chatservice/go.mod services/chatservice/go.sum ./chatservice/
RUN cd chatservice && go mod download

# Copy source code
COPY services/chatservice ./chatservice

# Build single unified binary
WORKDIR /app/chatservice
RUN go build -o chat-service ./cmd

# Stage 2: Production image
FROM alpine:3.19
WORKDIR /app

# Copy binary
COPY --from=builder /app/chatservice/chat-service .

# Expose single port for both REST and WebSocket
EXPOSE 8083

# Start unified server
CMD ["./chat-service"]
