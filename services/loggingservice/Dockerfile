# Runtime stage
FROM alpine:latest
WORKDIR /root/
COPY logging-service .
RUN mkdir -p /root/logs
EXPOSE 8889
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8889/health || exit 1
CMD ["./logging-service"]
