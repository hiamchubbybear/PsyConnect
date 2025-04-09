
FROM golang:1.24-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY bootstrap .
RUN go build -o consultation-service
EXPOSE 8080
CMD ["./consultation-service"]