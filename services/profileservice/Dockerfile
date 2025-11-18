# Stage 1: Install parent POM and build grpccommon
FROM maven:3.9.6-eclipse-temurin-21 AS grpccommon-builder
WORKDIR /build

# Copy and install parent POM
COPY services/pom.xml ./pom.xml
RUN mvn -N install -DskipTests && \
    echo "=== Parent POM installed ===" && \
    ls -la /root/.m2/repository/dev/psyconnect/

# Copy grpccommon files
COPY services/grpccommon/pom.xml ./grpccommon/pom.xml
COPY services/grpccommon/src ./grpccommon/src

# Build and install grpccommon
RUN mvn -f grpccommon/pom.xml clean install -DskipTests && \
    echo "=== grpccommon installed ===" && \
    ls -la /root/.m2/repository/dev/psyconnect/grpccommon/

# Stage 2: Build profileservice
FROM maven:3.9.6-eclipse-temurin-21 AS builder
WORKDIR /build

# Copy parent POM and install
COPY services/pom.xml ./pom.xml
RUN mvn -N install -DskipTests

# Copy Maven repository from stage 1
COPY --from=grpccommon-builder /root/.m2/repository /root/.m2/repository

# Verify grpccommon is available
RUN echo "=== Checking Maven repository ===" && \
    ls -la /root/.m2/repository/dev/psyconnect/ && \
    ls -la /root/.m2/repository/dev/psyconnect/grpccommon/ || echo "grpccommon not found!"

# Copy profileservice files
COPY services/profileservice/pom.xml ./profileservice/pom.xml
COPY services/profileservice/src ./profileservice/src

# Build profileservice with verbose output
RUN mvn -f profileservice/pom.xml clean package -DskipTests -B -X 2>&1 | tee /tmp/build.log || \
    (echo "=== BUILD FAILED - Last 100 lines ===" && tail -100 /tmp/build.log && exit 1)

# Stage 3: Runtime
FROM eclipse-temurin:21-jre-alpine
WORKDIR /app

RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring

COPY --from=builder /build/profileservice/target/*.jar app.jar

EXPOSE 8081

ENV JAVA_OPTS="-XX:+UseContainerSupport -XX:MaxRAMPercentage=75.0"

ENTRYPOINT ["sh", "-c", "java $JAVA_OPTS -jar app.jar"]
