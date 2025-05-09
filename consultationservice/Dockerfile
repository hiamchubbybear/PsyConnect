FROM golang:1.24-alpine

RUN apk add --no-cache git
WORKDIR /app
COPY consultationservice/consultationservice/go.mod consultationservice/consultationservice/go.sum ./
RUN go mod download
COPY consultationservice/consultationservice ./
RUN go build -o consultation-service ./consultationservice/cmd
EXPOSE 8084
CMD ["./consultation-service"]
