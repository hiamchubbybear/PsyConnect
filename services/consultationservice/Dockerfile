# Stage 1: Build
FROM golang:1.23-alpine AS builder
RUN apk add --no-cache git

# Context root
WORKDIR /app

# Copy go.mod và go.sum của service
COPY consultationservice/go.mod consultationservice/go.sum ./consultationservice/
RUN cd consultationservice && go mod download

# Copy toàn bộ source code service
COPY consultationservice ./consultationservice

# Build từ folder service
WORKDIR /app/consultationservice
RUN go build -o consultation-service ./cmd

# Stage 2: Production image
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/consultationservice/consultation-service .
EXPOSE 8083
CMD ["./consultation-service"]
