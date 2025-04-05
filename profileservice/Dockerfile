FROM maven:3.9.6-eclipse-temurin-17 AS builder

RUN apt-get update && apt-get install -y unzip curl

RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v3.19.6/protoc-3.19.6-linux-aarch_64.zip \
    && unzip protoc-3.19.6-linux-aarch_64.zip -d /usr/local \
    && chmod +x /usr/local/bin/protoc

WORKDIR /app
COPY . .

RUN mvn clean install
