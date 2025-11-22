# Consultation Service - PsyConnect

![Go](https://img.shields.io/badge/Go-1.21-blue?style=flat&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat&logo=go&logoColor=white)
![MongoDB](https://img.shields.io/badge/MongoDB-6.0-green?style=flat&logo=mongodb&logoColor=white)

## 📖 Overview

The **Consultation Service** is the core business logic handler for therapist-client interactions. It manages the matching process, consultation sessions, and the professional newsfeed.

## ✨ Features

- **Therapist Matching**: Algorithms to match clients with suitable therapists.
- **Session Management**: Booking, rescheduling, and tracking therapy sessions.
- **Professional Feed**: A social feed for therapists to share articles and updates.
- **Reviews & Ratings**: Feedback system for consultations.

## 🛠 Technology Stack

- **Language**: Golang 1.21
- **Framework**: Gin Web Framework
- **Database**: MongoDB
- **Communication**: gRPC (Internal), REST API (External)

## 🔌 API Endpoints

### Consultation
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/consultation/client/match` | Request therapist match |
| `POST` | `/consultation/session` | Create new session |
| `GET` | `/consultation/session/all` | Get user sessions |

### Feed
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/consultation/posts` | Get newsfeed |
| `POST` | `/v1/consultation/posts` | Create post |
| `POST` | `/v1/consultation/posts/:id/react` | React to post |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `DB_URI`: MongoDB connection string
   - `KAFKA_ADDRESS`: Kafka broker address
   - `SERVICE_PORT`: Service port (default: 8084)

## 🚀 Installation & Run

### Prerequisites
- Go 1.21+
- MongoDB

### Local Run
```bash
# 1. Navigate to directory
cd services/consultationservice

# 2. Download dependencies
go mod download

# 3. Run application
go run main.go
```

### Docker Run
```bash
docker build -t psyconnect/consultationservice .
docker run -p 8084:8084 --env-file .env psyconnect/consultationservice
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
