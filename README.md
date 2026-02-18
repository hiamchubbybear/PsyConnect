# PsyConnect - Mental Health Consultation Platform

PsyConnect is a comprehensive mental health platform designed to bridge the gap between clients and mental health professionals. Built on a robust microservices architecture, it ensures scalability, security, and a seamless user experience across web and mobile platforms.

## Build Status

| Service                | Tech Stack                    |
| ---------------------- | ----------------------------- |
| API Gateway            | Java / Spring Cloud Gateway   |
| Identity Service       | Java / Spring Boot 3 / MySQL  |
| Profile Service        | Java / Spring Boot 3 / Neo4j  |
| Consultation Service   | Go / Gin / MongoDB            |
| Chat Service           | Go / MongoDB / WebSocket      |
| Notification Service   | Go / Kafka Consumer           |
| Recommendation Service | Python / Flask / Scikit-learn |
| **Mobile Application** | **Flutter / Dart**            |
| **Web Application**    | Angular 17+ / TypeScript      |

## Technology Stack

### Backend Services

- **Communication**: gRPC (Inter-service), REST API (Client-facing).
- **API Gateway**: Java, Spring Cloud Gateway.
- **Identity Service**: Java, Spring Boot 3, Spring Security, OAuth2.
- **Profile Service**: Java, Spring Boot 3, Spring Data Neo4j.
- **Consultation Service**: Go (Golang), Gin Framework.
- **Chat Service**: Go (Golang), Gorilla WebSocket.
- **Notification Service**: Go (Golang).
- **Recommendation Service**: Python, Flask, Scikit-learn, Pandas.

### Frontend Applications

- **Mobile App**: **Flutter (Dart)** for iOS and Android.
- **Web Portal**: Angular 17+, RxJS, TailwindCSS, Angular Material.

### Data & Storage

- **Relational**: MySQL (Identity/Auth).
- **Graph**: Neo4j (Social Connections).
- **Document**: MongoDB (Chat History, Consultations).
- **Caching**: Redis.

### DevOps & Observability

- **Containerization**: Docker, Docker Compose.
- **Orchestration**: Kubernetes (K3s/MiniKube/Cloud).
- **CI/CD**: GitHub Actions (Self-hosted runners).
- **Message Broker**: Apache Kafka.
- **Monitoring**: Prometheus, Grafana.
- **Logging**: ELK Stack (Elasticsearch, Logstash, Kibana).
- **Tunneling**: Cloudflare Tunnel.

## System Architecture

PsyConnect employs a Microservices Architecture governed by Kubernetes.

### Core Components

- **API Gateway (8888)**: Entry point for client requests, handling routing and initial authorization.
- **Identity Service (8080)**: Manages authentication, user accounts, and OAuth2 integration.
- **Profile Service (8081)**: Handles user profiles and social graph connections (Neo4j).
- **Consultation Service (8083)**: Manages therapy sessions, bookings, and payments.
- **Chat Service (8084)**: Real-time messaging using WebSockets.
- **Notification Service (8082)**: Handles email and push notifications via Kafka events.

## Getting Started

### Prerequisites

- Docker & Docker Compose (for local infra)
- Kubernetes Cluster (K3s recommended) or Minikube
- Java JDK 17+
- Go 1.21+
- Node.js 18+ (for Webapp build)
- Python 3.9+

### Configuration

The project uses **Environment Variables** for configuration, injected via Kubernetes ConfigMaps and Secrets.

**Key Configuration Files:**

- `k8s/base/configmaps.yaml`: Service URLs, Hostnames, Ports.
- `k8s/base/secrets.yaml`: Sensitive credentials (DB passwords, API keys). _Note: This file is not git-tracked. Use `secrets-dev.yaml` as a reference._

### Installation & Deployment

1.  **Clone the repository**

    ```bash
    git clone https://github.com/hiamchubbybear/PsyConnect.git
    cd PsyConnect
    ```

2.  **Infrastructure Setup**
    Ensure your Kubernetes cluster is running.

    ```bash
    # Apply Infrastructure (Namespace, ConfigMaps, Secrets)
    kubectl apply -f k8s/base/namespace.yaml
    kubectl apply -f k8s/base/configmaps.yaml
    # See "Secrets Management" below for applying secrets
    ```

3.  **Secrets Management**
    Copy the development secrets template and populate it with your credentials:

    ```bash
    cp k8s/base/secrets-dev.yaml k8s/base/secrets.yaml
    # Edit k8s/base/secrets.yaml with real values
    kubectl apply -f k8s/base/secrets.yaml
    ```

4.  **Deploy Services**
    You can deploy services individually or all at once using the provided K8s manifests in `k8s/services/`.
    ```bash
    kubectl apply -f k8s/services/identity/
    kubectl apply -f k8s/services/apigateway/
    # ... apply other services
    ```

## Directory Structure

```
PsyConnect/
├── .github/            # CI/CD Workflows
├── k8s/                # Kubernetes Manifests
│   ├── base/           # Common Configs (ConfigMap, Secret, Namespace)
│   ├── infrastructure/ # Db/Kafka setup
│   └── services/       # Per-service Deployments
├── services/           # Microservices Source Code
│   ├── apigateway/
│   ├── identityservice/
│   ├── profileservice/
│   ├── consultationservice/
│   ├── chatservice/
│   ├── notificationservice/
│   └── ...
└── README.md           # Project Documentation
```

## Contributing

1.  Fork the repository.
2.  Create your feature branch (`git checkout -b feature/AmazingFeature`).
3.  Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4.  Push to the branch (`git push origin feature/AmazingFeature`).
5.  Open a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Contact

- **Project Lead**: Chessy
- **Repository**: https://github.com/hiamchubbybear/PsyConnect
