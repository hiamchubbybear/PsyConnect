FROM maven:3.9.6-eclipse-temurin-17 AS builder
WORKDIR /build
COPY . /build
RUN mvn clean install -pl profileservice -am -DskipTests -X
FROM eclipse-temurin:21
WORKDIR /app
COPY --from=builder /build/profileservice/target/*.jar app.jar
ENTRYPOINT ["java", "-jar", "app.jar"]
