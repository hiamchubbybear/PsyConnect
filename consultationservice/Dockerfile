FROM golang:1.24-alpine
RUN apk add --no-cache git
WORKDIR /app
COPY consultationservice/consultationservice/go.mod consultationservice/consultationservice/go.sum ./
RUN go mod download
COPY ./consultationservice /app
RUN go build -o consultation-service ./cmd
EXPOSE 8080
CMD ["./consultation-service"]
