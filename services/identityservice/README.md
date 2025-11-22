# Identity Service - PsyConnect

![Java](https://img.shields.io/badge/Java-17-orange?style=flat&logo=java&logoColor=white)
![Spring Boot](https://img.shields.io/badge/Spring_Boot-3.0-green?style=flat&logo=spring-boot&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0-blue?style=flat&logo=mysql&logoColor=white)

## 📖 Overview

The **Identity Service** is the backbone of user authentication and authorization in PsyConnect. It handles user registration, secure login (OAuth2/OIDC), account management, and role-based access control (RBAC).

## ✨ Features

- **User Management**: Registration, activation, profile updates, and deletion.
- **Authentication**: Secure login with JWT issuance and validation.
- **OAuth2 Integration**: Support for Google and Facebook login.
- **RBAC**: Role-based permission management (Admin, Therapist, Client).
- **Password Management**: Reset and change password functionality.

## 🛠 Technology Stack

- **Framework**: Java Spring Boot 3
- **Database**: MySQL 8.0
- **Security**: Spring Security, OAuth2 Resource Server
- **Build Tool**: Maven

## 🔌 API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/auth/login` | Authenticate user and get token |
| `POST` | `/auth/introspect` | Validate token validity |
| `POST` | `/auth/refresh` | Refresh access token |
| `POST` | `/auth/logout` | Invalidate token |

### User Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/identity/create` | Register new user |
| `POST` | `/identity/activate` | Activate account via code |
| `GET` | `/account/info` | Get current user info |
| `PUT` | `/account/update` | Update user details |
| `DELETE` | `/account/delete` | Delete account |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `MYSQL_SPRING_DATASOURCE_PASSWORD`: Your MySQL password
   - `SIGNER_KEY`: JWT signing key
   - `SERVER_PORT`: Service port (default: 8080)

## 🚀 Installation & Run

### Prerequisites
- Java JDK 17+
- Maven
- MySQL Server

### Local Run
```bash
# 1. Navigate to directory
cd services/identityservice

# 2. Install dependencies
mvn clean install

# 3. Run application
mvn spring-boot:run
```

### Docker Run
```bash
docker build -t psyconnect/identityservice .
docker run -p 8080:8080 --env-file .env psyconnect/identityservice
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
