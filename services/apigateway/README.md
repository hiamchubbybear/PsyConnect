# API Gateway - PsyConnect

![Java](https://img.shields.io/badge/Java-17-orange?style=flat&logo=java&logoColor=white)
![Spring Cloud Gateway](https://img.shields.io/badge/Spring_Cloud-Gateway-green?style=flat&logo=spring&logoColor=white)

## 📖 Overview

The **API Gateway** is the single entry point for all client requests. It handles routing, load balancing, and initial security checks (like JWT validation) before forwarding requests to the appropriate microservices.

## ✨ Features

- **Centralized Routing**: Maps external paths to internal microservice addresses.
- **Authentication Filter**: Validates JWT tokens and extracts user context.
- **Header Injection**: Injects `X-User-Id`, `X-Roles` into downstream requests.
- **CORS Management**: Centralized Cross-Origin Resource Sharing configuration.

## 🛠 Technology Stack

- **Framework**: Java Spring Boot 3
- **Library**: Spring Cloud Gateway
- **Security**: Spring Security (Resource Server)

## 🔌 Routes

| Path Prefix | Target Service | Description |
|-------------|----------------|-------------|
| `/identity/**` | Identity Service | Auth & User accounts |
| `/profile/**` | Profile Service | Profiles & Social |
| `/consultation/**` | Consultation Service | Matching & Sessions |
| `/noti/**` | Notification Service | Notifications |
| `/chats/**` | Chat Service | Messaging |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `IDENTITY_SERVICE_URI`: URL for Identity Service
   - `PROFILE_SERVICE_URI`: URL for Profile Service
   - `SERVER_PORT`: Gateway port (default: 8888)

## 🚀 Installation & Run

### Prerequisites
- Java JDK 17+
- Maven

### Local Run
```bash
# 1. Navigate to directory
cd services/apigateway

# 2. Install dependencies
mvn clean install

# 3. Run application
mvn spring-boot:run
```

### Docker Run
```bash
docker build -t psyconnect/apigateway .
docker run -p 8888:8888 --env-file .env psyconnect/apigateway
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
