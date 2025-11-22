# 🧠 PsyConnect - Mental Health Consultation Platform

![Project Status](https://img.shields.io/badge/status-active-success.svg)
![License](https://img.shields.io/badge/license-MIT-blue.svg)
![Docker](https://img.shields.io/badge/docker-%230db7ed.svg?style=flat&logo=docker&logoColor=white)
![Angular](https://img.shields.io/badge/angular-%23DD0031.svg?style=flat&logo=angular&logoColor=white)
![Flutter](https://img.shields.io/badge/Flutter-%2302569B.svg?style=flat&logo=Flutter&logoColor=white)
![Spring Boot](https://img.shields.io/badge/Spring_Boot-%236DB33F.svg?style=flat&logo=spring-boot&logoColor=white)
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=flat&logo=go&logoColor=white)

**PsyConnect** is a comprehensive, modern mental health platform designed to bridge the gap between clients and mental health professionals. Built on a robust microservices architecture, it ensures scalability, security, and a seamless user experience across web and mobile platforms.

---

## 🚀 Build Status

| Service | Status | Tech Stack |
|---------|--------|------------|
| **API Gateway** | ![Build](https://img.shields.io/badge/build-passing-success) | Java / Spring Cloud Gateway |
| **Identity Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Java / Spring Boot / MySQL |
| **Profile Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Java / Spring Boot / Neo4j |
| **Consultation Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Go / Gin / MongoDB |
| **Chat Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Go / MongoDB / WebSocket |
| **Notification Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Node.js / Express |
| **Recommendation Service** | ![Build](https://img.shields.io/badge/build-passing-success) | Python / Flask / Scikit-learn |

---

## ✨ Key Features

### 👤 For Clients
- **Smart Therapist Matching**: AI-driven recommendation system based on specialization, language, and availability.
- **Secure Consultations**: Private, secure channels for booking and managing therapy sessions.
- **Real-time Chat**: Instant messaging with therapists using WebSocket technology.
- **Mood Tracking**: Daily mood logging and analysis to track mental well-being.
- **Social Community**: Safe space to share experiences and connect with peers (Feed, Friends).
- **Cross-Platform**: Seamless experience on both Web (Angular) and Mobile (Flutter).

### 👨‍⚕️ For Therapists
- **Professional Profile**: Customizable profiles highlighting expertise and experience.
- **Schedule Management**: Flexible availability settings and session tracking.
- **Patient Management**: Tools to manage client interactions and history.

### ⚙️ System Capabilities
- **Centralized Authentication**: OAuth2/OIDC compliant identity management.
- **Real-time Notifications**: Email and push notifications for appointments and messages.
- **High Performance**: Caching with Redis and asynchronous processing with Kafka.
- **Observability**: Centralized logging (ELK) and monitoring (Prometheus/Grafana).

---

## 🏗 Architecture

PsyConnect employs a **Microservices Architecture** to ensure loose coupling and independent scalability.

### 🔌 Communication
- **Synchronous**: REST API (Client-facing), gRPC (Inter-service).
- **Asynchronous**: Apache Kafka (Event-driven architecture for notifications, logging, and data sync).

### 💾 Data Persistence
- **Relational**: MySQL (Identity/Auth data).
- **Graph**: Neo4j (Social connections, Profile relationships).
- **Document**: MongoDB (Chat history, Consultation records, Feeds).
- **In-Memory**: Redis (Caching, Session management).

---

## 🛠 Technology Stack

### Backend Services
- **Identity Service**: Spring Boot 3, Spring Security, OAuth2.
- **Profile Service**: Spring Boot 3, Spring Data Neo4j.
- **Consultation Service**: Golang, Gin Framework.
- **Chat Service**: Golang, Gorilla WebSocket.
- **Notification Service**: Node.js.
- **Recommendation Service**: Python, Pandas, Scikit-learn.
- **API Gateway**: Spring Cloud Gateway.

### Frontend
- **Web**: Angular 17+, RxJS, TailwindCSS/Material.
- **Mobile**: Flutter (Dart).

### DevOps & Infrastructure
- **Containerization**: Docker, Docker Compose.
- **Message Broker**: Apache Kafka, Zookeeper.
- **Monitoring**: Prometheus, Grafana, ELK Stack (Elasticsearch, Logstash, Kibana).
- **CI/CD**: GitHub Actions.

---

## 🚀 Getting Started

### Prerequisites
- Docker & Docker Compose
- Java JDK 17+
- Go 1.21+
- Node.js 18+
- Python 3.9+

### Installation & Setup

1.  **Clone the repository**
    ```bash
    git clone https://github.com/hiamchubbybear/PsyConnect.git
    cd PsyConnect
    ```

2.  **Configure Environment Variables**

    The project uses a hierarchical configuration system. You need to set up the root `.env` for infrastructure and individual `.env` files for each service.

    **Step 1: Root Configuration**
    ```bash
    cp .env.example .env
    # Edit .env to set your local passwords for MySQL, Neo4j, etc.
    ```

    **Step 2: Service Configuration**
    Copy the example config for each service:
    ```bash
    cp services/identityservice/.env.example services/identityservice/.env
    cp services/profileservice/.env.example services/profileservice/.env
    cp services/notificationservice/.env.example services/notificationservice/.env
    cp services/consultationservice/.env.example services/consultationservice/.env
    cp services/chatservice/.env.example services/chatservice/.env
    cp services/apigateway/.env.example services/apigateway/.env
    cp services/recommendationservice/.env.example services/recommendationservice/.env
    ```

3.  **Run with Docker Compose**
    ```bash
    # Start infrastructure (DBs, Kafka, Redis) and Services
    docker-compose up -d --build
    ```

4.  **Access the Application**
    - **Web App**: `http://localhost:4200` (if running locally) or via Gateway.
    - **API Gateway**: `http://localhost:8888`
    - **Eureka/Consul** (if applicable): `http://localhost:8761`

---

## 📂 Directory Structure

```
PsyConnect/
├── .github/            # CI/CD Workflows
├── dev/                # Development infrastructure (docker-compose, init scripts)
├── services/           # Microservices source code
│   ├── apigateway/
│   ├── identityservice/
│   ├── profileservice/
│   ├── consultationservice/
│   ├── chatservice/
│   ├── notificationservice/
│   ├── recommendationservice/
│   ├── mobileflutter/  # Mobile App
│   └── webapp/         # Web App
└── README.md           # Project Documentation
```

---

## 🤝 Contributing

We welcome contributions! Please follow these steps:

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 📞 Contact

- **Project Lead**: Chessy
- **Email**: tranvanhuy16032004@gmail.com
- **Repository**: [github.com/hiamchubbybear/PsyConnect](https://github.com/hiamchubbybear/PsyConnect)
