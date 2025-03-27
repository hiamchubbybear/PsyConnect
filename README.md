# PsyConnect

## Mental Health Platform Documentation

PsyConnect is a comprehensive mental health platform built with modern microservices architecture, focusing on security, scalability, and user experience.

## 📋 Table of Contents

- [Architecture Overview](#architecture-overview)
- [Technology Stack](#technology-stack)
- [Core Services](#core-services)
    - [Identity Service](identityservice/README.md)
    - [Profile Service](profileservice/README.md)
    - [Payment Service](#payment-service)
    - [Matching Service](#matching-service)
    - [Messaging Service](#messaging-service)
    - [Review & Feedback](#review-feedback)
    - [Booking & Scheduling](#booking-scheduling)
    - [Content Services](#content-services)
    - [Support Services](#support-services)
    - [Api Gateway](apigateway/README.md)
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

### Identity Service

- User authentication and registration
- Multifactorial authentication
- Authorization management

### Profile Service

- User profile management
- Customization options
- Privacy controls

### Payment Service

- Multiple payment method support
- Subscription management
- Transaction history

### Matching Service

- Advanced therapist-client matching
- Preference-based recommendations
- Availability matching

### Messaging Service

- Real-time communication
- Secure file sharing
- Message history

### Review & Feedback

- Therapist/client reviews
- Session feedback system
- Rating analytics

### Booking & Scheduling

- Appointment management
- Calendar integration
- Automated reminders

### Content Services

- Blog management
- Support group forums
- Community interactions

### Support Services

- Push notifications
- Email/SMS alerts
- System logging
- Performance monitoring

## 🌍 Setup Environment

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

## 🚀 Installation

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

## 🤝 Contributing

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

## 📜 License

This project is licensed under the MIT License. See [LICENSE](LICENSE) for details.

## 📞 Contact

- Project Lead: Chessy
- Email: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- Repository: [github.com/hiamchubbybear/PsyConnect](https://github.com/hiamchubbybear/PsyConnect)

