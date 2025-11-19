
FROM maven:3.9.6-eclipse-temurin-21 AS builder
WORKDIR /app

COPY services/apigateway/pom.xml ./apigateway/pom.xml
COPY services/apigateway/src ./apigateway/src


RUN mvn -f apigateway/pom.xml clean package -DskipTests


FROM eclipse-temurin:21-jre-alpine
WORKDIR /app


RUN addgroup -S spring && adduser -S spring -G spring
USER spring:spring


COPY --from=builder /app/apigateway/target/apigateway.jar app.jar

EXPOSE 8888

ENTRYPOINT ["java", "-jar", "app.jar"]
