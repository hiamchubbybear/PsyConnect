# PsyConnect

## Mental Health Platform Documentation

PsyConnect is a comprehensive mental health platform built with modern microservices architecture, focusing on security, scalability, and user experience.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Technology Stack](#technology-stack)
- [Core Services](#core-services)
- [Setup Environment](#setup-environment)
- [Installation](#installation)
- [Development & Deployment](#development-deployment)
- [Contributing](#contributing)

## 🏗 Architecture Overview 

### Microservices Architecture

Our platform utilizes a distributed system design with the following key features:

- Modular, independently scalable microservices
- Cross-platform compatibility (web & mobile)
- Container-based deployment with Docker
- Automated CI/CD pipeline through GitHub
- Secure gRPC communication between microservices

## 🛠 Technology Stack

### Backend Services

| Category      | Technologies          |
| ------------- | --------------------- |
| Languages     | Java, Node.js, Golang |
| Communication | gRPC, RESTful APIs    |
| Caching       | Redis                 |

### Database Infrastructure

- **Multi-Database Strategy**
    - Neo4j: Identity and authentication services
    - MySQL, PostgreSQL: Relational data storage
    - Redis: Real-time operations and caching
- **Scalability**: Optimized indexing and partitioning

### Frontend Technologies

- **Web Platform**: Angular (Progressive Web App)
- **Mobile Platform**: Flutter (Cross-platform development)

### Infrastructure & DevOps

- **Containerization**: Docker
- **CI/CD**: GitHub Actions
- **Security**: OAuth 2.0, JWT Authentication

## 🔧 Core Services

### Backend Services
- **[Identity Service](identityservice/README.md)** - User authentication, account management, OAuth2 integration
- **[Profile Service](profileservice/README.md)** - User profiles, social features, mood tracking, settings
- **[Consultation Service](consultationservice/README.md)** - Therapist-client matching, session management
- **[Notification Service](notificationservice/README.md)** - Email notifications, messaging system
- **[Chat Service](chatservice/README.md)** - Real-time messaging and communication
- **[Blog Service](blogservice/README.md)** - Content management, articles, resources
- **[Search Service](searchservice/README.md)** - Platform-wide search, filtering, recommendations
- **[Recommendation Service](recommendationservice/README.md)** - AI-powered therapist-client matching
- **[User Profile Service](userprofile/README.md)** - Extended profile management and customization

### Frontend Applications
- **[Web Application](webapp/README.md)** - Angular-based Progressive Web App
- **[Mobile Application](mobileflutter/README.md)** - Flutter cross-platform mobile app

### Infrastructure Services
- **[API Gateway](apigateway/README.md)** - Request routing, authentication, security
- **[Logging Service](loggingservice/)** - Centralized logging and monitoring
- **[Kafka](kafka/)** - Message queue and event streaming

## 📚 API Documentation Overview

### Core API Endpoints
All API requests go through the API Gateway at `http://localhost:8888`

#### Authentication & Identity
- `POST /auth/login` - User authentication
- `POST /identity/create` - User registration
- `GET /account/info` - Account information
- `GET /oauth2/userInfo` - OAuth2 authentication

#### User Profiles & Social
- `GET /profile` - User profile details
- `POST /profile/friend/request` - Send friend request
- `GET /mood` - Current mood status
- `PUT /user-setting` - Update user settings

#### Consultation & Matching
- `GET /consultation/therapist` - Therapist information
- `POST /consultation/client` - Create client profile
- `POST /consultation/therapist/match` - Request therapist matching
- `GET /consultation/session/all` - User sessions

#### Notifications & Communication
- `POST /noti/activate` - Send activation email
- `POST /chat/create` - Create chat session
- `GET /chat/{id}/messages` - Chat message history

#### Content & Search
- `GET /blog` - Blog posts and articles
- `GET /search/therapists` - Search therapists
- `POST /recommend/therapist` - Get AI recommendations

### Authentication
All protected endpoints require JWT token in Authorization header:
```
Authorization: Bearer <your_jwt_token>
```

For detailed API documentation, see individual service README files.

![Image](https://private-user-images.githubusercontent.com/51482452/403526354-a5116da7-6296-4387-9599-62641acd6993.png?jwt=eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnaXRodWIuY29tIiwiYXVkIjoicmF3LmdpdGh1YnVzZXJjb250ZW50LmNvbSIsImtleSI6ImtleTUiLCJleHAiOjE3NDUwNzYzNzksIm5iZiI6MTc0NTA3NjA3OSwicGF0aCI6Ii81MTQ4MjQ1Mi80MDM1MjYzNTQtYTUxMTZkYTctNjI5Ni00Mzg3LTk1OTktNjI2NDFhY2Q2OTkzLnBuZz9YLUFtei1BbGdvcml0aG09QVdTNC1ITUFDLVNIQTI1NiZYLUFtei1DcmVkZW50aWFsPUFLSUFWQ09EWUxTQTUzUFFLNFpBJTJGMjAyNTA0MTklMkZ1cy1lYXN0LTElMkZzMyUyRmF3czRfcmVxdWVzdCZYLUFtei1EYXRlPTIwMjUwNDE5VDE1MjExOVomWC1BbXotRXhwaXJlcz0zMDAmWC1BbXotU2lnbmF0dXJlPTA4Yzk2ZjkwYmQ2NWIxZWUyMmMwN2ZlNmJhY2M4OWVhYTYwYWNjZjY2Yjc0NGQ4MTI2ZjMzNWIwOTQ5YzJkNWMmWC1BbXotU2lnbmVkSGVhZGVycz1ob3N0In0.XaZux9tZfHm5XNoYGEhj32HYGg6zPkeVB2P7o8ruNgQ)

## Setup Environment

Ensure you have the following environment variables set up before running the project:

```env
GOOGLE_CLIENT_ID={your-variable}
GOOGLE_CLIENT_SECRET={your-variable}
MAIL_USERNAME={your-variable}
MAIL_PASSWORD={your-variable}
JWT_SIGNER_KEY={your-variable}
NEO4J_PASSWORD_PROFILE_SERVICE={your-variable}
MYSQL_SPRING_DATASOURCE_PASSWORD={your-variable}
MYSQL_SPRING_DATASOURCE_USERNAME={your-variable}
SPRING_DATASOURCE_URL={your-variable}
SERVER_PORT={your-variable}
```

**Note**: All environment variables must be configured correctly before running the application.

## Installation

### Prerequisites

```
- Docker
- Java SDK
- Node.js
- Golang
- Flutter SDK
```

### Local Development Setup

```bash
# Clone repository
git clone https://github.com/hiamchubbybear/PsyConnect.git
cd PsyConnect

# Build containers
docker-compose up --build
```

### Monitoring Stack

- Logging: ELK Stack
- Metrics: Prometheus & Grafana
- Performance: New Relic

##  Contributing

We welcome contributions! Please follow these steps:

1. Fork the repository
2. Create your feature branch
   ```bash
   git checkout -b feature/YourFeature
   ```
3. Commit your changes
   ```bash
   git commit -m 'Add YourFeature'
   ```
4. Push to your branch
   ```bash
   git push origin feature/YourFeature
   ```
5. Open a Pull Request

##  License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

##  Contact

- Project Lead: Chessy
- Email: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- Repository: [github.com/hiamchubbybear/PsyConnect](https://github.com/hiamchubbybear/PsyConnect)

