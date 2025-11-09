# PsyConnect System Architecture Overview

## 1. Tổng quan toàn hệ thống

### Các service chính
- **API Gateway** (Java Spring Boot): Điểm vào duy nhất cho tất cả requests
- **Identity Service** (Java Spring Boot): Xác thực người dùng và quản lý tài khoản
- **Profile Service** (Java Spring Boot): Hồ sơ người dùng và tính năng xã hội
- **Consultation Service** (Go Gin): Ghép đôi therapist/client và quản lý phiên tư vấn
- **Chat Service** (Go Gin): Nhắn tin thời gian thực (đang phát triển)
- **Notification Service** (Node.js): Email và push notifications
- **Recommendation Service** (Python): AI matching therapist-client
- **Logging Service** (Go): Centralized logging với OpenTelemetry
- **Web App** (Angular): Progressive Web App
- **Mobile App** (Flutter): Ứng dụng di động đa nền tảng

### Technology stack
- **Backend**: Java Spring Boot, Go Gin, Node.js Express, Python FastAPI
- **Frontend**: Angular (PWA), Flutter (mobile)
- **Databases**: MySQL, Neo4j, MongoDB, Redis
- **Infrastructure**: Docker, docker-compose
- **Communication**: REST APIs, gRPC, Kafka, WebSocket
- **Monitoring**: OpenTelemetry, Loki, Promtail

### Kiểu kiến trúc
Microservices architecture với API Gateway pattern, event-driven communication qua Kafka, và polyglot technology stack.

## 2. Sơ đồ giao tiếp (text only)

```
┌─────────────────┐    ┌─────────────────┐
│   Web App       │    │   Mobile App    │
│   (Angular)     │    │   (Flutter)     │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │
          ┌─────────────────────┐
          │    API Gateway      │
          │   (Spring Boot)     │
          │     Port: 8888      │
          └─────────┬───────────┘
                    │
          ┌─────────┼─────────┐
          │         │         │
┌─────────▼──┐ ┌────▼───┐ ┌───▼─────────┐
│Identity    │ │Profile │ │Consultation │
│Service     │ │Service │ │Service      │
│Port: 8080  │ │8081    │ │Configurable │
│MySQL       │ │Neo4j   │ │MongoDB      │
│JWT/OAuth2  │ │Redis   │ │Redis        │
└────────────┘ └────────┘ └─────────────┘
          │         │         │
          │         │         │
┌─────────▼──┐ ┌────▼───┐ ┌───▼─────────┐
│Notification│ │Chat    │ │Recommendation│
│Service     │ │Service │ │Service      │
│Port: 8082  │ │Dev     │ │Port: 8087   │
│Node.js     │ │Go      │ │Python       │
│Kafka       │ │WebSock │ │AI/ML        │
└────────────┘ └────────┘ └─────────────┘
          │         │         │
          └─────────┼─────────┘
                    │
          ┌─────────────────────┐
          │   Logging Service   │
          │       (Go)         │
          │  OpenTelemetry     │
          └─────────────────────┘
```

### Service Communication Details:
- **API Gateway** → **Identity Service**: REST (user auth)
- **API Gateway** → **Profile Service**: REST (user profiles)
- **API Gateway** → **Consultation Service**: REST (matching/sessions)
- **API Gateway** → **Notification Service**: REST (emails)
- **API Gateway** → **Chat Service**: WebSocket (real-time)
- **API Gateway** → **Recommendation Service**: REST (AI matching)

- **Identity Service** → **Profile Service**: gRPC (user creation)
- **Consultation Service** → **Profile Service**: gRPC (profile checks)
- **All Services** → **Logging Service**: Events/logs
- **Services** → **Kafka**: Event publishing/consuming

## 3. Database overview

### Service Database Mapping:
- **Identity Service**: MySQL (accounts, tokens, permissions)
- **Profile Service**: Neo4j (social graph, relationships)
- **Consultation Service**: MongoDB (therapists, clients, sessions)
- **Chat Service**: MongoDB (planned for messages)
- **Notification Service**: MySQL (FCM tokens)
- **Recommendation Service**: JSON files (therapist/client data)
- **Logging Service**: Log files (app.log, error.log)

### Collections/Tables chính:
- **MySQL (Identity)**: accounts, blacklist_tokens, permissions, roles, tokens
- **Neo4j (Profile)**: Profile nodes, FriendRelationship edges, Mood nodes
- **MongoDB (Consultation)**: therapists, clients, sessions, matched, swipes
- **Redis**: Session caching, real-time data

### Relationships:
- User accounts (MySQL) → User profiles (Neo4j)
- Therapist profiles (MongoDB) → User profiles (Neo4j)
- Consultation sessions (MongoDB) → Client/Therapist profiles

## 4. Event / Queue / Kafka Map

### Topics:
- **user.registered**: User account creation
- **account.activated**: Email verification
- **profile.updated**: Profile changes
- **consultation.booked**: Session booking
- **chat.message**: Real-time messages
- **notification.sent**: Email/push delivery
- **gateway.requests**: API gateway logs
- **app.logs**: Application logs

### Producers/Consumers:
- **Producers**:
  - Identity Service: user.registered, account.activated
  - Profile Service: profile.updated
  - Consultation Service: consultation.booked
  - Chat Service: chat.message
  - API Gateway: gateway.requests
  - All Services: app.logs

- **Consumers**:
  - Notification Service: user.registered, account.activated, consultation.booked
  - Logging Service: app.logs, gateway.requests
  - Profile Service: consultation.booked (activity logs)

### Schema rút gọn:
```json
{
  "topic": "user.registered",
  "payload": {
    "user_id": "uuid",
    "email": "string",
    "timestamp": "datetime"
  }
}
```

## 5. Backend flow

### Auth Flow:
1. Client → API Gateway (/auth/login)
2. API Gateway → Identity Service
3. Identity Service validates credentials
4. JWT token generated and returned
5. Subsequent requests include JWT in headers

### Payment Flow:
(Not implemented yet)
1. Client initiates payment
2. Payment Service processes transaction
3. Webhook updates consultation status
4. Notification sent to both parties

### Recommendation Flow:
1. Client requests therapist recommendations
2. Consultation Service → Recommendation Service
3. AI matching using cosine similarity
4. Top matches returned with compatibility scores
5. Client can swipe/accept matches

### Message Flow:
1. Client opens chat
2. WebSocket connection established
3. Messages sent via WebSocket
4. Chat Service stores in database
5. Real-time delivery to recipient
6. Push notification if offline

## 6. Environment variables overview

### API Gateway:
- JWT_SIGNER_KEY: JWT signing secret
- KAFKA_BROKER: Kafka broker URL

### Identity Service:
- MYSQL_SPRING_DATASOURCE_*: MySQL connection details
- GOOGLE_CLIENT_ID/SECRET: OAuth2 credentials
- MAIL_USERNAME/PASSWORD: SMTP settings
- JWT_SIGNER_KEY: Token signing

### Profile Service:
- NEO4J_*: Neo4j connection
- KAFKA_BROKER: Event streaming

### Consultation Service:
- DB_HOST/PORT/USER/PASS: MongoDB/Redis config
- SERVICE_PORT: Service port

### Notification Service:
- SMTP_HOST/PORT: Email server
- KAFKA_BROKER: Event consumption

### Recommendation Service:
- MODEL_PATH: AI model location
- DATA_PATH: Training data location

### Common:
- BOOTSTRAP_ADDRESS: Kafka brokers
- SERVICE_ENV: Environment (dev/prod)

**Note**: Most services lack .env.example files, making onboarding difficult.

## 7. List tất cả WARNINGS

### Configuration Issues:
- **Missing .env.example files** for all services
- **Hardcoded URLs** in docker-compose.yml
- **No centralized config management**
- **Security**: Sensitive data in plain environment variables

### Architecture Issues:
- **No service discovery** mechanism
- **Tight coupling** between some services
- **Inconsistent API versioning** across services
- **Missing health checks** for some services
- **No circuit breakers** implemented

### Development Issues:
- **Inconsistent logging** formats
- **Missing API documentation** (Swagger/OpenAPI)
- **No integration tests** setup
- **Mixed technology stack** without standards
- **Incomplete services** (Chat, Logging partially implemented)

### Production Readiness:
- **No monitoring/alerting** system
- **Missing backup strategies**
- **No rate limiting** implemented
- **Security headers** not configured
- **No proper error handling** standards

### Database Issues:
- **Multiple database types** without clear data ownership
- **No database migrations** strategy
- **Missing indexes** and optimization
- **Data consistency** concerns across services

### Deployment Issues:
- **No CI/CD pipeline** configured
- **Manual deployment** process
- **No container orchestration** (Kubernetes)
- **Missing resource limits** in Docker configs

### Recommendations:
1. Implement service mesh (Istio/Linkerd)
2. Add comprehensive monitoring (Prometheus/Grafana)
3. Implement proper secret management
4. Add API gateway features (rate limiting, caching)
5. Standardize error handling and logging
6. Implement health checks and metrics
7. Add integration and e2e testing
8. Implement proper API versioning
9. Add database migration strategies
10. Set up CI/CD with automated testing
