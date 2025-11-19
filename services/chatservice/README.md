# Chat Service - Unified Server Architecture

Real-time messaging service with REST API and WebSocket on a single port.

## 🚀 Quick Start

### Local Development
```bash
# Start dependencies
cd dev && ./scripts/start-all.sh

# Switch to dev profile
./dev-check.sh dev

# Run service
cd services/chatservice
go run cmd/main.go
```

**Access:**
- REST API: `http://localhost:8083`
- WebSocket: `ws://localhost:8083/ws`

### Docker
```bash
# Build
docker build -f ./services/chatservice/Dockerfile -t chatservice .

# Run
docker run -p 8083:8083 --env-file .env.dev chatservice
```

---

## 📡 API Endpoints

### REST API (Port 8083)

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/chats` | Create chat message |
| `GET` | `/chats/:id` | Get chat by ID |
| `GET` | `/chats/conversation/:conversationId` | Get conversation messages |
| `PUT` | `/chats/:id` | Update chat message |
| `DELETE` | `/chats/:id` | Delete chat message |
| `POST` | `/conversations` | Create conversation |
| `GET` | `/conversations/by-users` | Get conversation by users |
| `GET` | `/health` | Health check |

### WebSocket (Port 8083)

**Endpoint:** `GET /ws?conversationId=xxx&receiver=yyy`

**Connection:**
```javascript
const ws = new WebSocket('ws://localhost:8083/ws?conversationId=abc&receiver=user2');

ws.onmessage = (event) => {
  console.log('Message:', event.data);
};

ws.send(JSON.stringify({ text: 'Hello!' }));
```

---

## 🏗️ Architecture

### Unified Server Design
```
┌─────────────────────────────────┐
│   Chat Service (Port 8083)      │
│                                 │
│  ┌──────────────────────────┐  │
│  │  Gin Router              │  │
│  │  - REST endpoints        │  │
│  │  - WebSocket endpoint    │  │
│  │  - Shared middleware     │  │
│  └──────────────────────────┘  │
│                                 │
│  ┌──────────────────────────┐  │
│  │  WebSocket Hub           │  │
│  │  - Client management     │  │
│  │  - Message broadcasting  │  │
│  └──────────────────────────┘  │
│                                 │
│  ┌──────────────────────────┐  │
│  │  MongoDB Repository      │  │
│  │  - Message persistence   │  │
│  └──────────────────────────┘  │
└─────────────────────────────────┘
```

**Benefits:**
- ✅ Single port (8083)
- ✅ Single process
- ✅ Shared authentication
- ✅ Unified logging
- ✅ Simple deployment

---

## 🧰 Tech Stack

| Component | Technology |
|-----------|------------|
| Language | Go 1.24 |
| Framework | Gin |
| WebSocket | Gorilla WebSocket |
| Database | MongoDB |
| Message Queue | Kafka |
| Containerization | Docker |

---

## 🔧 Configuration

### Environment Variables

```env
# Service
SERVICE_PORT=8083
SERVICE_HOST=localhost

# Database
DB_HOST=localhost
DB_PORT=27017
DB_USER=root
DB_PASS=password
DB_NAME=chatservice

# Kafka
KAFKA_ADDRESS=localhost:9092
KAFKA_TOPIC=chat-message

# gRPC
GRPC_ADDRESS=localhost:50051

# JWT
JWT_SIGNER_KEY=your-secret-key
```

### Profiles
- **Dev**: `.env.dev` (localhost)
- **CICD**: `.env.cicd` (Docker service names)

Switch profiles: `./dev-check.sh dev|cicd`

---

## 🧪 Testing

### REST API
```bash
# Create conversation
curl -X GET "http://localhost:8083/conversations/by-users?user1=123&user2=456"

# Send message
curl -X POST http://localhost:8083/chats \
  -H "Content-Type: application/json" \
  -d '{"senderId":"123","conversationId":"abc","text":"Hello"}'

# Health check
curl http://localhost:8083/health
```

### WebSocket
```bash
# Using wscat
npm install -g wscat
wscat -c "ws://localhost:8083/ws?conversationId=abc&receiver=456"

# Send message
> {"text": "Hello from WebSocket!"}
```

### HTML Test Client
Open `home.html` in browser:
- Update URL to `ws://localhost:8083/ws`
- Connect and test real-time messaging

---

## 📦 Dependencies

```go
require (
    github.com/gin-gonic/gin
    github.com/gin-contrib/cors
    github.com/gorilla/websocket
    go.mongodb.org/mongo-driver
    github.com/segmentio/kafka-go
)
```

---

## 🚀 Deployment

### Docker Compose
```yaml
chatservice:
  image: hiamchubbybear/chatservice:latest
  ports:
    - "8083:8083"
  environment:
    - SERVICE_PORT=8083
    - DB_HOST=mongodb
    - KAFKA_ADDRESS=kafka:9092
  depends_on:
    - mongodb
    - kafka
```

### Kubernetes
```yaml
apiVersion: v1
kind: Service
metadata:
  name: chatservice
spec:
  ports:
    - port: 8083
      targetPort: 8083
  selector:
    app: chatservice
```

---

## 📊 Monitoring

### Health Check
```bash
curl http://localhost:8083/health
# Response: {"status":"healthy"}
```

### Logs
```bash
# Docker
docker logs -f chatservice

# Local
go run cmd/main.go
```

---

## 🔒 Security

- JWT authentication on WebSocket connections
- CORS enabled for web clients
- Environment-based configuration
- No hardcoded credentials

---

## 📝 Development

### Run Locally
```bash
go run cmd/main.go
```

### Build
```bash
go build -o chat-service ./cmd
```

### Test
```bash
go test ./...
```

---

## 🎯 Port Summary

| Service | Port | Protocols |
|---------|------|-----------|
| Chat Service | 8083 | HTTP + WebSocket |

**Note:** Both REST API and WebSocket run on the same port (8083) for simplified deployment.
