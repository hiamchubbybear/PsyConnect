# 📊 PsyConnect Logging Service

A centralized logging service for the PsyConnect platform that collects, processes, and visualizes logs from all microservices.

## 🎯 Features

- ✅ **Kafka Integration**: Consumes log events from Kafka topics
- ✅ **Loki Integration**: Pushes logs to Grafana Loki for querying and visualization
- ✅ **File Logging**: Writes logs to rotating files as backup
- ✅ **Structured Logging**: JSON-formatted logs with rich metadata
- ✅ **Health Checks**: HTTP endpoints for monitoring service health
- ✅ **Metrics**: Prometheus-compatible metrics endpoint
- ✅ **Distributed Tracing**: Tempo integration for request tracing
- ✅ **Grafana Dashboards**: Pre-built dashboards for log visualization

## 🏗 Architecture

```
┌─────────────────┐
│  Microservices  │
│  (Java/Go/Node) │
└────────┬────────┘
         │ Kafka Topics
         ▼
┌─────────────────┐
│ Logging Service │
│   (Go + Kafka)  │
└────┬───────┬────┘
     │       │
     │       └──────────────┐
     ▼                      ▼
┌─────────┐         ┌──────────────┐
│  Loki   │◄────────│ File Logger  │
└────┬────┘         └──────────────┘
     │
     ▼
┌─────────┐
│ Grafana │
└─────────┘
```

## 🚀 Quick Start

### Prerequisites

- Go 1.24+
- Docker & Docker Compose
- Kafka running on `localhost:9092`

### 1. Clone and Setup

```bash
cd services/loggingservice
cp .env.example .env
```

### 2. Start Observability Stack

```bash
docker-compose up -d
```

This starts:
- **Loki** (port 3100) - Log aggregation
- **Grafana** (port 3000) - Visualization
- **Prometheus** (port 9090) - Metrics
- **Tempo** (port 3200) - Tracing

### 3. Run Logging Service

```bash
go run main.go
```

Or build and run:

```bash
go build -o logging-service
./logging-service
```

### 4. Access Grafana

Open http://localhost:3000

- **Username**: `admin`
- **Password**: `16032004`

Navigate to **Dashboards** → **PsyConnect** → **Logs Overview**

## 📁 Project Structure

```
loggingservice/
├── main.go                          # Entry point
├── pkg/
│   ├── settings/                    # Configuration management
│   │   └── settings.go
│   ├── models/                      # Data models
│   │   └── log_event.go
│   ├── kafka/                       # Kafka consumer
│   │   └── consumer.go
│   ├── loki/                        # Loki client
│   │   └── client.go
│   ├── filelogger/                  # File logging
│   │   └── filelogger.go
│   └── server/                      # HTTP server
│       └── server.go
├── grafana/
│   ├── provisioning/                # Auto-provisioning configs
│   │   ├── datasources/
│   │   └── dashboards/
│   └── dashboards/                  # Dashboard JSON files
│       └── logs-overview.json
├── prometheus/
│   └── prometheus.yml               # Prometheus config
├── tempo/
│   └── tempo.yaml                   # Tempo config
├── logs/                            # Log files (gitignored)
├── Docker-compose.yml
├── .env.example
└── README.md
```

## ⚙️ Configuration

All configuration is done via environment variables. See `.env.example` for all options.

### Key Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `KAFKA_BOOTSTRAP_SERVERS` | `127.0.0.1:9092` | Kafka broker address |
| `LOKI_ENABLED` | `true` | Enable Loki integration |
| `LOKI_URL` | `http://localhost:3100` | Loki endpoint |
| `FILE_LOGGER_ENABLED` | `true` | Enable file logging |
| `SERVER_PORT` | `8080` | HTTP server port |

## 📊 Log Event Schema

```json
{
  "timestamp": "2025-11-25T18:30:00Z",
  "level": "INFO",
  "service": "identity-service",
  "message": "User logged in successfully",
  "traceId": "abc123",
  "userId": "user-456",
  "action": "LOGIN",
  "metadata": {
    "ip": "192.168.1.1",
    "userAgent": "Mozilla/5.0..."
  }
}
```

### Supported Log Levels

- `DEBUG` - Detailed debugging information
- `INFO` - General informational messages
- `LOG` - Application-specific logs
- `AUDIT` - Audit trail events
- `WARN` - Warning messages
- `ERROR` - Error messages
- `FATAL` - Critical errors

## 🔍 Querying Logs in Grafana

### Example LogQL Queries

**All logs from a specific service:**
```logql
{service="identity-service"}
```

**Error logs only:**
```logql
{level="ERROR"}
```

**Logs for a specific user:**
```logql
{userId="user-123"}
```

**Logs with specific action:**
```logql
{action="LOGIN"} |= "failed"
```

**Logs by trace ID:**
```logql
{traceId="abc123"}
```

## 🏥 Health Checks

The service exposes HTTP endpoints for monitoring:

### Health Check
```bash
curl http://localhost:8080/health
```

Response:
```json
{
  "status": "healthy",
  "uptime": "2h30m15s",
  "messages_count": 15234,
  "errors_count": 12,
  "timestamp": "2025-11-25T18:30:00Z",
  "version": "1.0.0"
}
```

### Metrics
```bash
curl http://localhost:8080/metrics
```

### Readiness Probe
```bash
curl http://localhost:8080/ready
```

## 📈 Metrics (Prometheus)

The service exposes metrics that can be scraped by Prometheus:

- `messages_processed_total` - Total messages processed
- `errors_total` - Total errors encountered
- `error_rate` - Percentage of errors

## 🔧 Development

### Run Tests
```bash
go test ./...
```

### Build
```bash
go build -o logging-service
```

### Run with Debug Logging
```bash
LOG_LEVEL=debug go run main.go
```

## 🐳 Docker Deployment

Build Docker image:
```bash
docker build -t psyconnect/logging-service:latest .
```

Run with Docker:
```bash
docker run -d \
  --name logging-service \
  -p 8080:8080 \
  -e KAFKA_BOOTSTRAP_SERVERS=kafka:9092 \
  -e LOKI_URL=http://loki:3100 \
  psyconnect/logging-service:latest
```

## 🛠 Troubleshooting

### Logs not appearing in Grafana

1. Check Loki is running: `docker ps | grep loki`
2. Check logging service logs: `docker logs logging-service`
3. Verify Kafka connectivity: Check `KAFKA_BOOTSTRAP_SERVERS`
4. Check Grafana datasource configuration

### High memory usage

- Reduce `LOKI_BATCH_SIZE` (default: 100)
- Increase flush interval in Loki client
- Enable log compression: `LOG_COMPRESS=true`

### Kafka consumer lag

- Increase Kafka consumer instances
- Optimize log processing handlers
- Check Kafka broker health

## 📚 Related Documentation

- [Grafana Loki Documentation](https://grafana.com/docs/loki/latest/)
- [Prometheus Documentation](https://prometheus.io/docs/)
- [Kafka Documentation](https://kafka.apache.org/documentation/)
- [PsyConnect Main README](../../README.md)

## 🤝 Contributing

See main project [CONTRIBUTING.md](../../CONTRIBUTING.md)

## 📄 License

MIT License - see [LICENSE](../../LICENSE)

---

**Built with ❤️ for PsyConnect**
