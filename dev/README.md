# Development Dependencies

Infrastructure dependencies for running PsyConnect services locally.

## Quick Start

```bash
# Start all dependencies
cd dev
./scripts/start-all.sh

# Or manually
docker-compose up -d
```

## Services

| Service | Port | Purpose | Credentials |
|---------|------|---------|-------------|
| **MySQL** | 3306 | Identity Service database | root / (see .env) |
| **Neo4j** | 7474, 7687 | Profile Service graph DB | neo4j / (see .env) |
| **MongoDB** | 27017 | Consultation/Chat Services | root / (see .env) |
| **Redis** | 6379 | Caching layer | (see .env) |
| **Kafka** | 9092 | Message broker | - |
| **Zookeeper** | 2181 | Kafka coordination | - |

## Usage

### Start All Dependencies
```bash
./scripts/start-all.sh
```

### Start Specific Services
```bash
# Databases only
docker-compose up -d mysql neo4j mongodb

# Messaging only
docker-compose up -d zookeeper kafka

# Cache only
docker-compose up -d redis
```

### Stop All
```bash
./scripts/stop-all.sh
# or
docker-compose down
```

### Clean All Data
```bash
./scripts/clean.sh
# or
docker-compose down -v
```

### View Logs
```bash
# All services
docker-compose logs -f

# Specific service
docker-compose logs -f mysql
docker-compose logs -f kafka
```

## Configuration

1. Copy environment template:
   ```bash
   cp .env.example .env
   ```

2. Edit `.env` with your credentials:
   ```env
   MYSQL_PASSWORD=your-password
   NEO4J_PASSWORD=your-password
   MONGO_PASSWORD=your-password
   REDIS_PASSWORD=your-password
   ```

## Running Services Locally

After starting dependencies, run services in dev mode:

```bash
# Switch all services to dev profile (localhost)
./dev-check.sh dev

# Run Java services
cd services/apigateway && mvn spring-boot:run
cd services/identityservice && mvn spring-boot:run
cd services/profileservice && mvn spring-boot:run

# Run Go services
cd services/consultationservice && go run cmd/main.go
cd services/chatservice && go run main.go

# Run Node.js services
cd services/notificationservice && npm start
```

## Accessing Services

### MySQL
```bash
mysql -h localhost -P 3306 -u root -p
# Password from .env
```

### Neo4j Browser
Open: http://localhost:7474
- Username: `neo4j`
- Password: (from .env)

### MongoDB
```bash
mongosh mongodb://root:password@localhost:27017
```

### Redis
```bash
redis-cli -h localhost -p 6379 -a password
```

## Troubleshooting

### Port Already in Use
```bash
# Kill process using port
./scripts/../script/KILLPORT.sh 3306
./scripts/../script/KILLPORT.sh 9092
```

### Service Won't Start
```bash
# Check logs
docker-compose logs [service-name]

# Restart specific service
docker-compose restart [service-name]
```

### Clean Start
```bash
# Stop and remove everything
docker-compose down -v

# Start fresh
./scripts/start-all.sh
```

## Health Checks

All services have health checks configured. Check status:
```bash
docker-compose ps
```

Healthy services show `healthy` status.
