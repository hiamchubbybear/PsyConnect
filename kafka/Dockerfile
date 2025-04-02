FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o app ./cmd
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/app .
COPY ./config/config.yaml .
ENV CONFIG_PATH=config.yaml
CMD ["./app"]