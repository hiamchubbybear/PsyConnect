FROM golang:1.24-alpine

# Cài thêm git nếu dùng go get
RUN apk add --no-cache git

# Tạo thư mục làm việc
WORKDIR /app

# Copy go.mod và go.sum trước (tận dụng cache)
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ source code vào container
COPY . .

# Build từ file main.go trong thư mục cmd/
RUN go build -o consultation-service ./cmd

EXPOSE 8080

# Chạy binary
CMD ["./consultation-service"]
