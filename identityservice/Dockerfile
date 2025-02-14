FROM maven:3.9.6-eclipse-temurin-21 AS builder

WORKDIR /app

COPY . .

RUN mvn clean package -DskipTests

FROM openjdk:23-jdk
WORKDIR /app

COPY --from=builder /app/target/identityservice.jar identityservice.jar

EXPOSE 8080

CMD ["java", "-jar", "identityservice.jar"]
