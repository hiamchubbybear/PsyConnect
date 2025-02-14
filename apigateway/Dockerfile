FROM maven:3.9.6-eclipse-temurin-21 AS builder

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

FROM openjdk:23-jdk
WORKDIR /app
# Copy JAR đã build từ stage trước
COPY --from=builder /app/target/apigateway.jar apigateway.jar

# Expose port
EXPOSE 8888

# Run Spring Boot app
CMD ["java", "-jar", "apigateway.jar"]
