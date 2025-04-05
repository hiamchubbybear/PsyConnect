FROM maven:3.9.6-eclipse-temurin-17 AS builder

RUN apt-get update && apt-get install -y unzip curl

# Tự động phát hiện kiến trúc
RUN if [ "$(uname -m)" = "aarch64" ]; then \
      PROTOC_ARCH="aarch_64"; \
    else \
      PROTOC_ARCH="x86_64"; \
    fi && \
    curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.19.6/protoc-3.19.6-linux-${PROTOC_ARCH}.zip && \
    unzip protoc-3.19.6-linux-${PROTOC_ARCH}.zip -d /usr/local && \
    chmod +x /usr/local/bin/protoc

WORKDIR /app
COPY . .

# Build với tham số OS
RUN mvn clean install -DskipTests \
    -Dos.detected.classifier=linux-$(if [ "$(uname -m)" = "aarch64" ]; then echo "aarch_64"; else echo "x86_64"; fi)