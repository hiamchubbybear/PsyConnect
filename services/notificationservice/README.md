# Notification Service - PsyConnect

## Overview
The **Notification Service** handles email notifications and messaging features within the PsyConnect platform. It provides endpoints for sending activation emails, account update notifications, and other email-based communications to users.

## Features
- Email notifications for account activation
- Account update notifications
- Kafka consumer for processing notification events
- Template-based email system
- Integration with external email providers

## Technology Stack
- **Backend**: Node.js with Express
- **Message Queue**: Apache Kafka
- **Email Service**: SMTP/External email provider
- **API Communication**: RESTful APIs

## API Endpoints

### Email Notifications
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/noti/activate` | Send account activation email |
| POST | `/noti/account-update` | Send account update notification email |

### Request Examples

#### Send Activation Email
```bash
POST /noti/activate
Content-Type: application/json

{
  "username": "john_doe",
  "code": "12345",
  "email": "john@example.com",
  "fullname": "John Doe"
}
```

#### Send Account Update Email
```bash
POST /noti/account-update
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com"
}
```

## Environment Variables
Ensure the following environment variables are set before running the service:
```env
PORT=8082
SMTP_HOST={your-smtp-host}
SMTP_PORT={your-smtp-port}
MAIL_USERNAME={your-email-username}
MAIL_PASSWORD={your-email-password}
KAFKA_BROKER={your-kafka-broker}
```

## Setup & Configuration

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/notificationservice
   ```
2. Install dependencies:
   ```bash
   npm install
   ```
3. Start the service:
   ```bash
   npm start
   ```

### Docker
```bash
docker build -t psyconnect/notification-service .
docker run -p 8082:8082 psyconnect/notification-service
```

## Kafka Integration
The service includes a Kafka consumer that automatically processes notification events from other microservices, enabling real-time email notifications.

## Contributing
We welcome contributions! Please follow the standard Git workflow:
1. Fork the repository
2. Create your feature branch (`git checkout -b feature/YourFeature`)
3. Commit your changes (`git commit -m 'Add YourFeature'`)
4. Push to your branch (`git push origin feature/YourFeature`)
5. Open a Pull Request

## Contact
For inquiries, reach out via:
- **Project Lead**: Chessy
- **Email**: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- **GitHub Repository**: [PsyConnect](https://github.com/hiamchubbybear/PsyConnect)