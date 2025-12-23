# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-w -s" -o notification-service ./cmd/main.go

# Runtime stage
FROM alpine:latest

# Install ca-certificates for HTTPS, tzdata for timezones, wget for health check, and netcat for Kafka check
RUN apk --no-cache add ca-certificates tzdata wget netcat-openbsd

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/notification-service .

# Copy templates
COPY --from=builder /app/templates ./templates

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Copy .env file (optional, can use env vars from docker-compose)
COPY .env .env

# Expose HTTP port
EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application with entrypoint
ENTRYPOINT ["/entrypoint.sh"]
