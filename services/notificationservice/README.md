# Notification Service - PsyConnect

![Node.js](https://img.shields.io/badge/Node.js-18-green?style=flat&logo=nodedotjs&logoColor=white)
![Express](https://img.shields.io/badge/Express-4.18-black?style=flat&logo=express&logoColor=white)
![Kafka](https://img.shields.io/badge/Kafka-Events-black?style=flat&logo=apachekafka&logoColor=white)

## 📖 Overview

The **Notification Service** handles all outbound communications from the platform. It listens for events (like "User Registered" or "Appointment Booked") and delivers notifications via Email or Push Notifications.

## ✨ Features

- **Email Notifications**: SMTP integration for transactional emails.
- **Event-Driven**: Consumes Kafka topics to trigger notifications asynchronously.
- **Templates**: HTML email templates for consistent branding.
- **Push Notifications**: (Planned) Mobile push support.

## 🛠 Technology Stack

- **Runtime**: Node.js 18
- **Framework**: Express.js
- **Message Broker**: Apache Kafka (Consumer)
- **Email**: Nodemailer / SMTP

## 🔌 API Endpoints

While primarily event-driven, it exposes endpoints for testing or direct triggers.

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/noti/activate` | Send activation email |
| `POST` | `/noti/account-update` | Send security alert |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `SMPT_HOST`, `SMPT_USER`, `SMPT_PASS`: Email credentials
   - `KAFKA_BROKER`: Kafka address
   - `PORT`: Service port (default: 8082)

## 🚀 Installation & Run

### Prerequisites
- Node.js 18+
- Kafka (for event consumption)

### Local Run
```bash
# 1. Navigate to directory
cd services/notificationservice

# 2. Install dependencies
npm install

# 3. Run application
npm start
```

### Docker Run
```bash
docker build -t psyconnect/notificationservice .
docker run -p 8082:8082 --env-file .env psyconnect/notificationservice
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
