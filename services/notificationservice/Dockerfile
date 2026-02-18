# Build stage
FROM golang:1.23-alpine AS builder

# Force enable Go modules
ENV GO111MODULE=on

# Use official Go proxy to avoid network issues
ENV GOPROXY=https://proxy.golang.org,direct

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Verify module setup and build
RUN ls -la && \
    cat go.mod && \
    go list -m && \
    CGO_ENABLED=0 GOOS=linux go build -v -o notification-service ./cmd/...

# Runtime stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS, tzdata for timezones, wget for health check
RUN apk add --no-cache ca-certificates tzdata wget netcat-openbsd

# Copy binary from builder
COPY --from=builder /app/notification-service .

# Copy templates
COPY --from=builder /app/templates ./templates

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 8082

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8082/health || exit 1

# Run the application with entrypoint
ENTRYPOINT ["/entrypoint.sh"]
