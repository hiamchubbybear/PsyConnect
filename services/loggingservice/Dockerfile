# Build stage
FROM golang:1.24-alpine AS builder

# Install build dependencies including librdkafka
RUN apk add --no-cache \
    git \
    gcc \
    g++ \
    musl-dev \
    librdkafka-dev \
    pkgconf \
    make

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application with dynamic linking to system librdkafka
RUN CGO_ENABLED=1 GOOS=linux go build -tags dynamic -o logging-service .

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ca-certificates \
    tzdata \
    librdkafka

WORKDIR /root/

# Copy the binary from builder
COPY --from=builder /app/logging-service .

# Create logs directory
RUN mkdir -p /root/logs

# Expose HTTP port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the application
CMD ["./logging-service"]
