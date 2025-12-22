# Runtime only
FROM alpine:3.19
WORKDIR /app
RUN apk add --no-cache ca-certificates
COPY chat-service .
EXPOSE 8083
CMD ["./chat-service"]
