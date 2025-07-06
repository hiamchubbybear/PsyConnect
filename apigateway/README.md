# API Gateway Service

## Overview
The **API Gateway Service** is responsible for routing requests to different microservices within the PsyConnect platform. It uses **Spring Cloud Gateway** to define routes and filters, ensuring efficient communication between clients and backend services.

## Features
- Routes incoming requests to appropriate microservices
- Adds user-related headers for downstream services
- Logs events for debugging and monitoring
- Request routing to identity, profile, notification, and consultation services

## Configuration

### Environment Variables
| Variable       | Description                         |
|---------------|-------------------------------------|
| `SIGNER_KEY`  | Secret key used for token validation |
| `server.port` | Port on which the API Gateway runs |

### Application Properties
The routing configuration is defined in `application.properties`:
```properties
# Service name
spring.application.name=api_service

# Base URLs
base.url=http://localhost
identity.port=8080
profile.port=8081

# Port
server.port=8888

# Identity Service Config
spring.cloud.gateway.routes[0].id=identity_service
spring.cloud.gateway.routes[0].uri=http://localhost:8080
spring.cloud.gateway.routes[0].predicates[0]=Path=/identity/**,/auth/**,/oauth2/**,/account/**

# Profile Service Config
spring.cloud.gateway.routes[1].id=profile_service
spring.cloud.gateway.routes[1].uri=http://localhost:8081
spring.cloud.gateway.routes[1].predicates[0]=Path=/profile/**,/profile/friends/**,/mood/**,/user-setting/**

# Notification Service Config
spring.cloud.gateway.routes[2].id=notification_service
spring.cloud.gateway.routes[2].uri=http://localhost:8082
spring.cloud.gateway.routes[2].predicates[0]=Path=/noti/**

# Consultation Service Config
spring.cloud.gateway.routes[3].id=consultation_service
spring.cloud.gateway.routes[3].uri=http://localhost:8084
spring.cloud.gateway.routes[3].predicates[0]=Path=/consultation/**

# Logging
logging.level.org.springframework.web=DEBUG
logging.level.org.springframework.cloud.gateway=TRACE
```

## Request Processing
The gateway processes requests and forwards them to the appropriate microservice based on the configured routes.

### Headers Added by Gateway
| Header Name    | Description                                    |
|----------------|------------------------------------------------|
| `X-User-Id`    | User ID extracted from request context         |
| `X-Profile-Id` | Profile ID extracted from request context      |
| `X-Roles`      | Roles and permissions for authorization         |

### Example Request Flow
#### Request:
```http
GET /profile HTTP/1.1
Host: localhost:8888
```

#### Forwarded Request (after processing):
```http
GET /profile HTTP/1.1
Host: localhost:8081
X-User-Id: 12345
X-Profile-Id: 67890
X-Roles: role.client:permission
```

## Running the Service
To start the API Gateway Service, run:
```sh
mvn spring-boot:run
```

## Logging & Debugging
- Set logging level in `application.properties`:
  ```properties
  logging.level.org.springframework.web=DEBUG
  logging.level.org.springframework.cloud.gateway=TRACE
  ```
- Logs will show routing details and filter actions

## Contributing
Feel free to contribute by submitting issues or pull requests!

