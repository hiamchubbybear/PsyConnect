# Chat Service - PsyConnect

## Overview
The **Chat Service** handles real-time messaging and communication features within the PsyConnect platform. This service enables secure messaging between users, therapists, and clients for consultation sessions and support.

## Features
- Real-time messaging between users
- Secure communication channels
- Message history and persistence
- Group chat functionality
- File sharing capabilities
- Message encryption and security

## Technology Stack
- **Backend**: [Technology to be implemented]
- **Real-time Communication**: WebSocket/Socket.io
- **Database**: [Database to be implemented]
- **API Communication**: RESTful APIs
- **Security**: End-to-end encryption

## API Endpoints

### Chat Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chat/create` | Create new chat session |
| GET | `/chat/{id}` | Get chat session details |
| DELETE | `/chat/{id}` | Delete chat session |

### Message Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chat/{id}/message` | Send message to chat |
| GET | `/chat/{id}/messages` | Get chat message history |
| PUT | `/chat/{id}/message/{messageId}` | Update message |
| DELETE | `/chat/{id}/message/{messageId}` | Delete message |

### User Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/chat/{id}/users` | Add user to chat |
| DELETE | `/chat/{id}/users/{userId}` | Remove user from chat |
| GET | `/chat/{id}/users` | Get chat participants |

*Note: This service is currently in development. API endpoints are subject to change.*

## Setup & Configuration

### Environment Variables
```env
PORT=8083
DATABASE_URL={your-database-url}
REDIS_URL={your-redis-url}
JWT_SECRET={your-jwt-secret}
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/chatservice
   ```
2. Install dependencies and start the service (commands will be added when implemented)

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
