# 🔄 MIGRATION GUIDE: Node.js → Go

Complete guide for migrating from Node.js notification service to Go version.

---

## 📊 OVERVIEW

| Aspect           | Node.js     | Go            |
| ---------------- | ----------- | ------------- |
| **Language**     | JavaScript  | Go            |
| **Runtime**      | Node.js 18+ | Native binary |
| **Memory**       | ~150MB      | ~20-50MB      |
| **Startup**      | 1-2s        | <1s           |
| **Image Size**   | ~200MB      | ~50MB         |
| **Dependencies** | npm         | go mod        |
| **Features**     | 100%        | 100%          |

---

## 🎯 MIGRATION STRATEGY

### Phase 1: Preparation (1 hour)

1. **Review current Node.js service**
2. **Setup Go environment**
3. **Configure database**
4. **Setup Firebase**

### Phase 2: Testing (2 hours)

1. **Run Go service alongside Node.js**
2. **Test all features**
3. **Compare outputs**
4. **Performance testing**

### Phase 3: Switchover (1 hour)

1. **Stop Node.js service**
2. **Start Go service**
3. **Monitor logs**
4. **Verify functionality**

### Phase 4: Cleanup (30 minutes)

1. **Remove Node.js service**
2. **Update documentation**
3. **Archive old code**

---

## 📋 PRE-MIGRATION CHECKLIST

- [ ] Go service built successfully
- [ ] Database accessible
- [ ] Kafka cluster running
- [ ] Firebase credentials ready
- [ ] SMTP credentials valid
- [ ] All environment variables documented
- [ ] Backup of current service
- [ ] Rollback plan ready

---

## 🚀 STEP-BY-STEP MIGRATION

### Step 1: Setup Go Service (Parallel)

Run Go service on different port for testing:

```bash
cd services/notificationservice-go

# Edit .env
PORT=8092  # Different port

# Run
go run cmd/main.go
```

Both services will now run:

- Node.js: `http://localhost:8082`
- Go: `http://localhost:8092`

### Step 2: Database Migration

**Option A: Share same database**

Both services can use the same MySQL database:

```env
# Node.js (Sequelize)
DB_HOST=localhost
DB_NAME=notification_db

# Go (GORM)
DB_HOST=localhost
DB_NAME=notification_db
```

GORM will auto-migrate tables if they don't exist.

**Option B: Separate databases**

```sql
CREATE DATABASE notification_db_go;
```

Then migrate data:

```bash
mysqldump notification_db | mysql notification_db_go
```

### Step 3: Test Parity

**Test 1: Email Sending**

```bash
# Node.js
curl -X POST http://localhost:8082/mail/activate \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","code":"12345","fullname":"Test"}'

# Go
curl -X POST http://localhost:8092/mail/activate \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","code":"12345","fullname":"Test"}'
```

**Test 2: Kafka Consumer**

Both services will consume same Kafka topics (same consumer group).

Monitor logs:

```bash
# Node.js
docker logs -f psyconnect-notificationservice-dev

# Go
docker logs -f psyconnect-notification-go
```

**Test 3: Push Notifications**

```bash
# Save FCM token
curl -X POST http://localhost:8092/notification/token \
  -H "Content-Type: application/json" \
  -d '{"userId":"test123","token":"fcm-token"}'

# Trigger notification via Kafka
# Check both services handle it
```

### Step 4: Performance Comparison

**Memory Usage**:

```bash
# Node.js
docker stats psyconnect-notificationservice-dev

# Go
docker stats psyconnect-notification-go
```

**Response Time**:

```bash
# Node.js
time curl http://localhost:8082/health

# Go
time curl http://localhost:8092/health
```

### Step 5: Update Docker Compose

Edit `dev/docker-compose.yml`:

```yaml
# BEFORE (Node.js)
notification-service:
  build: ./services/notificationservice
  image: dev-notificationservice
  container_name: psyconnect-notificationservice-dev
  ports:
    - "8082:8082"
  # ... rest of config

# AFTER (Go)
notification-service:
  build:
    context: ./services/notificationservice-go
  container_name: psyconnect-notification-go-dev
  ports:
    - "8082:8082"
  environment:
    - KAFKA_BROKERS=kafka:9094
    - DATABASE_DSN=root:${DB_PASSWORD}@tcp(mysql:3306)/notification_db?charset=utf8mb4&parseTime=True
    - SMPT_MAIL=${SMPT_MAIL}
    - SMPT_APP_PASS=${SMPT_APP_PASS}
    - FIREBASE_CREDENTIALS=/app/firebase-credentials.json
  volumes:
    - ./firebase-credentials.json:/app/firebase-credentials.json:ro
  depends_on:
    - kafka
    - mysql
```

### Step 6: Switchover

**Stop Node.js service**:

```bash
docker stop psyconnect-notificationservice-dev
docker rm psyconnect-notificationservice-dev
```

**Start Go service**:

```bash
cd dev
docker-compose up -d notification-service
```

**Verify**:

```bash
# Check logs
docker logs -f psyconnect-notification-go-dev

# Test health
curl http://localhost:8082/health

# Test email
curl -X POST http://localhost:8082/mail/activate \
  -H "Content-Type: application/json" \
  -d '{"username":"test","email":"test@example.com","code":"12345","fullname":"Test"}'
```

### Step 7: Monitor

Monitor for 24-48 hours:

```bash
# Logs
docker logs -f psyconnect-notification-go-dev

# Errors
docker logs psyconnect-notification-go-dev | grep "❌"

# Email activity
docker logs psyconnect-notification-go-dev | grep "📧"

# Kafka activity
docker logs psyconnect-notification-go-dev | grep "📨"
```

### Step 8: Cleanup

After confirming Go service works:

```bash
# Archive Node.js service
mv services/notificationservice services/notificationservice-old

# Or delete
rm -rf services/notificationservice

# Update .gitignore if needed
echo "services/notificationservice-old/" >> .gitignore
```

---

## 🔄 ROLLBACK PLAN

If issues occur, rollback to Node.js:

```bash
# Stop Go service
docker stop psyconnect-notification-go-dev

# Start Node.js service
docker-compose up -d notification-service

# Verify
curl http://localhost:8082/health
```

---

## 📊 VERIFICATION CHECKLIST

After migration:

- [ ] Service starts successfully
- [ ] Health checks passing
- [ ] Database connected
- [ ] Kafka consumer connected
- [ ] Emails sending correctly
- [ ] Push notifications working
- [ ] All 22 topics handled
- [ ] No errors in logs
- [ ] Memory usage < 50MB
- [ ] Response time < 100ms
- [ ] All API endpoints working

---

## 🐛 COMMON ISSUES

### Issue 1: Database Connection Failed

**Cause**: DSN format different between Sequelize and GORM

**Node.js (Sequelize)**:

```env
DB_HOST=localhost
DB_USER=root
DB_PASSWORD=pass
DB_NAME=notification_db
```

**Go (GORM)**:

```env
DATABASE_DSN=root:pass@tcp(localhost:3306)/notification_db?charset=utf8mb4&parseTime=True
```

### Issue 2: Kafka Topics Not Consuming

**Cause**: Consumer group conflict

**Solution**: Use different consumer group during testing:

```env
# Node.js
KAFKA_GROUP_ID=notification-group

# Go (testing)
KAFKA_GROUP_ID=notification-group-go

# Go (production)
KAFKA_GROUP_ID=notification-group
```

### Issue 3: Firebase Credentials Not Found

**Cause**: Path mismatch

**Solution**: Use absolute path:

```env
FIREBASE_CREDENTIALS=/app/firebase-credentials.json
```

And mount volume:

```yaml
volumes:
  - ./firebase-credentials.json:/app/firebase-credentials.json:ro
```

---

## 📈 EXPECTED IMPROVEMENTS

After migration:

**Performance**:

- ✅ 3x less memory usage
- ✅ 2x faster startup
- ✅ 4x smaller Docker image
- ✅ Lower CPU usage

**Reliability**:

- ✅ Type-safe code
- ✅ Compile-time error checking
- ✅ Better error handling
- ✅ No runtime surprises

**Deployment**:

- ✅ Single binary
- ✅ No node_modules
- ✅ Faster builds
- ✅ Easier debugging

---

## 🎯 POST-MIGRATION TASKS

1. **Update documentation**
2. **Train team on Go service**
3. **Setup monitoring**
4. **Archive Node.js code**
5. **Celebrate!** 🎉

---

## 📞 SUPPORT

If you encounter issues during migration:

1. Check logs of both services
2. Compare database states
3. Verify Kafka consumer groups
4. Test each feature individually
5. Use rollback plan if needed

---

**Migration Guide Complete!** 🚀

**Estimated migration time**: 4-5 hours
**Recommended migration window**: Low-traffic period
**Rollback time**: <5 minutes
