# 🚀 DEPLOYMENT GUIDE - Notification Service (Go)

## 📋 Prerequisites

- Docker & Docker Compose
- MySQL 8.0+
- Kafka cluster
- Firebase project (optional, for push notifications)
- Gmail account with App Password

---

## 🔧 SETUP

### 1. Environment Configuration

Create `.env` file:

```bash
cp .env.example .env
```

Edit `.env` with your values:

```env
# SMTP
SMPT_MAIL=your-email@gmail.com
SMPT_APP_PASS=your-app-password

# Database
DB_PASSWORD=your-mysql-password

# Firebase (optional)
FIREBASE_CREDENTIALS=/path/to/firebase-credentials.json
```

### 2. Firebase Setup (Optional)

If using push notifications:

1. Go to [Firebase Console](https://console.firebase.google.com/)
2. Create/select project
3. Go to Project Settings → Service Accounts
4. Generate new private key
5. Save as `firebase-credentials.json`
6. Place in service directory

### 3. Database Setup

**Option A: Use Docker Compose** (Recommended)

```bash
docker-compose up -d mysql
```

**Option B: Existing MySQL**

```sql
CREATE DATABASE notification_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

Tables will be auto-created on first run.

---

## 🚀 DEPLOYMENT OPTIONS

### Option 1: Docker Compose (Recommended)

**Full stack** (includes MySQL, Kafka, Zookeeper):

```bash
docker-compose up -d
```

**Only notification service** (if you have existing MySQL/Kafka):

```bash
docker-compose up -d notification-service
```

**View logs**:

```bash
docker-compose logs -f notification-service
```

### Option 2: Standalone Docker

**Build image**:

```bash
docker build -t psyconnect-notification-go:latest .
```

**Run container**:

```bash
docker run -d \
  --name notification-service \
  -p 8082:8082 \
  -e KAFKA_BROKERS=kafka:9094 \
  -e DATABASE_DSN="user:pass@tcp(mysql:3306)/notification_db?charset=utf8mb4&parseTime=True" \
  -e SMPT_MAIL=your-email@gmail.com \
  -e SMPT_APP_PASS=your-app-password \
  -v $(pwd)/firebase-credentials.json:/app/firebase-credentials.json:ro \
  --network psyconnect-network \
  psyconnect-notification-go:latest
```

### Option 3: Binary (Local Development)

**Build**:

```bash
go build -o notification-service cmd/main.go
```

**Run**:

```bash
./notification-service
```

Or:

```bash
go run cmd/main.go
```

---

## ✅ VERIFICATION

### 1. Check Service Health

```bash
curl http://localhost:8082/health
# Expected: {"status":"healthy","service":"notification-service"}
```

### 2. Check Logs

**Docker**:

```bash
docker logs -f psyconnect-notification-go
```

**Expected output**:

```
🚀 Starting PsyConnect Notification Service (Go)
✅ Configuration loaded
✅ Database connected
✅ Firebase/FCM initialized
✅ Email service initialized
✅ Notification service initialized
✅ Kafka consumer initialized
🌐 HTTP API server listening on port 8082
✅ Notification service is running...
```

### 3. Run Tests

```bash
./test.sh
```

### 4. Test Email Sending

```bash
curl -X POST http://localhost:8082/mail/activate \
  -H "Content-Type: application/json" \
  -d '{
    "username": "test",
    "email": "test@example.com",
    "code": "12345",
    "fullname": "Test User"
  }'
```

---

## 🔄 INTEGRATION WITH MAIN DOCKER COMPOSE

Add to your main `dev/docker-compose.yml`:

```yaml
notification-service:
  build:
    context: ../services/notificationservice-go
    dockerfile: Dockerfile
  container_name: psyconnect-notification-go-dev
  ports:
    - "8082:8082"
  environment:
    # SMTP
    - SMPT_HOST=smtp.gmail.com
    - SMPT_PORT=587
    - SMPT_MAIL=${SMPT_MAIL}
    - SMPT_APP_PASS=${SMPT_APP_PASS}

    # Kafka
    - KAFKA_BROKERS=kafka:9094
    - KAFKA_GROUP_ID=notification-group

    # Database
    - DATABASE_DSN=root:${DB_PASSWORD}@tcp(mysql:3306)/notification_db?charset=utf8mb4&parseTime=True

    # Firebase (optional)
    - FIREBASE_CREDENTIALS=/app/firebase-credentials.json

    # Server
    - PORT=8082
    - ENVIRONMENT=development

  volumes:
    - ../services/notificationservice-go/firebase-credentials.json:/app/firebase-credentials.json:ro

  depends_on:
    - kafka
    - mysql

  restart: unless-stopped

  networks:
    - psyconnect-network
```

---

## 🔧 TROUBLESHOOTING

### Database Connection Failed

**Check MySQL is running**:

```bash
docker ps | grep mysql
```

**Test connection**:

```bash
mysql -h localhost -u root -p
```

**Check DSN format**:

```
user:password@tcp(host:port)/database?charset=utf8mb4&parseTime=True
```

### Kafka Connection Failed

**Check Kafka is running**:

```bash
docker ps | grep kafka
```

**Check broker address**:

- Inside Docker: `kafka:9094`
- Outside Docker: `localhost:19092`

### Firebase Initialization Failed

**Check credentials file exists**:

```bash
ls -la firebase-credentials.json
```

**Check file is valid JSON**:

```bash
cat firebase-credentials.json | jq .
```

**Check environment variable**:

```bash
echo $FIREBASE_CREDENTIALS
```

### Email Sending Failed

**Check SMTP credentials**:

```bash
# Test with telnet
telnet smtp.gmail.com 587
```

**Check App Password** (not regular password):

- Go to Google Account → Security → 2-Step Verification → App Passwords
- Generate new app password
- Use that in `.env`

---

## 📊 MONITORING

### Health Checks

```bash
# Health
curl http://localhost:8082/health

# Ready
curl http://localhost:8082/ready
```

### Logs

```bash
# Docker
docker logs -f psyconnect-notification-go

# Follow specific events
docker logs -f psyconnect-notification-go | grep "📨"  # Kafka messages
docker logs -f psyconnect-notification-go | grep "📧"  # Emails
docker logs -f psyconnect-notification-go | grep "❌"  # Errors
```

### Metrics

Access metrics endpoint:

```bash
curl http://localhost:8082/metrics
```

---

## 🔄 UPDATES & MAINTENANCE

### Update Service

```bash
# Pull latest code
git pull

# Rebuild
docker-compose build notification-service

# Restart
docker-compose up -d notification-service
```

### Database Migrations

GORM auto-migrates on startup. No manual migration needed.

### Backup Database

```bash
docker exec psyconnect-mysql mysqldump -u root -p notification_db > backup.sql
```

---

## 🎯 PRODUCTION CHECKLIST

- [ ] Environment variables configured
- [ ] Database created and accessible
- [ ] Kafka cluster running
- [ ] Firebase credentials (if using push notifications)
- [ ] SMTP credentials valid
- [ ] Health checks passing
- [ ] Logs showing no errors
- [ ] Test email sent successfully
- [ ] Kafka consumer connected
- [ ] Database tables created
- [ ] Monitoring setup
- [ ] Backup strategy in place

---

## 📞 SUPPORT

If you encounter issues:

1. Check logs: `docker logs -f psyconnect-notification-go`
2. Verify environment variables
3. Test database connection
4. Test Kafka connection
5. Check Firebase credentials
6. Verify SMTP settings

---

**Deployment Guide Complete!** 🚀
