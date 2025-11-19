# 📜 Scripts Directory

All utility scripts for managing PsyConnect services.

## Available Scripts

### `start-all.sh` - Background Mode ⭐
Start all dependencies and microservices in background.

```bash
cd scripts
./start-all.sh
```

**Features:**
- ✅ Auto-cleans old logs
- ✅ Starts all services in background
- ✅ Logs to `../logs/*.log`
- ✅ Fast startup

**Use when:** You want services running in background while working on other things.

---

### `start-foreground.sh` - Foreground Mode 🖥️
Start all services in tmux with logs visible in terminal.

```bash
cd scripts
./start-foreground.sh
```

**Features:**
- ✅ All logs visible in terminal
- ✅ Split screen with tmux
- ✅ Easy to monitor all services
- ✅ Ctrl+B then arrow keys to switch panes

**Use when:** You want to see real-time logs while developing.

**Requirements:** `brew install tmux`

---

### `view-logs.sh` - View All Logs 📊
Stream all service logs in real-time.

```bash
cd scripts
./view-logs.sh
```

Tails all log files with timestamps. Press Ctrl+C to stop.

---

### `stop-all.sh` - Stop Services 🛑
Stop all running services and dependencies.

```bash
cd scripts
./stop-all.sh
```

Gracefully stops all microservices and Docker containers.

---

### `dev-check.sh` - Profile Manager 🔧
Switch between dev (localhost) and cicd (Docker) profiles.

```bash
cd scripts
./dev-check.sh check  # View current profiles
./dev-check.sh dev    # Switch to localhost
./dev-check.sh cicd   # Switch to Docker hosts
```

---

## Quick Start

### Background Mode (Recommended)
```bash
cd scripts
./start-all.sh      # Start everything
./view-logs.sh      # View logs
./stop-all.sh       # Stop when done
```

### Foreground Mode (For Debugging)
```bash
cd scripts
./start-foreground.sh  # Start in tmux
# Ctrl+B then D to detach
# tmux attach -t psyconnect to reattach
```

---

## Logs Management

### Background Mode
- **Location:** `logs/[service-name].log`
- **Auto-cleanup:** Yes (on restart)
- **View:** `./view-logs.sh` or `tail -f logs/*.log`

### Foreground Mode
- **Location:** Terminal panes
- **Auto-cleanup:** On session close
- **View:** Built-in (live in tmux)

---

## Service URLs

| Service | URL | Port |
|---------|-----|------|
| **API Gateway** | http://localhost:8888 | 8888 |
| **Identity** | http://localhost:8080 | 8080 |
| **Profile** | http://localhost:8081 | 8081 |
| **Notification** | http://localhost:8082 | 8082 |
| **Chat** | http://localhost:8083 | 8083 |
| **Consultation** | http://localhost:8084 | 8084 |

## Troubleshooting

### Port already in use
```bash
lsof -i :8080  # Find process
kill -9 [PID]  # Kill it
```

### Service won't start (background)
```bash
cat logs/identityservice.log  # Check logs
```

### tmux not found (foreground)
```bash
brew install tmux
```

### Docker not running
```bash
colima start
```
