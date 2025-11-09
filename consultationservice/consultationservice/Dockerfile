FROM golang:1.24-alpine

# Cài git để go mod hoạt động
RUN apk add --no-cache git

WORKDIR /app

# Copy go.mod/go.sum trước để tận dụng cache
COPY go.mod go.sum ./

# Tải dependency, đồng bộ vendor
RUN go mod download
COPY vendor ./vendor

# Copy toàn bộ source code
COPY . .

# Build với vendor
RUN go build -mod=vendor -o consultation-service ./cmd

EXPOSE 8084

CMD ["./consultation-service"]
