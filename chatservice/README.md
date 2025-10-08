# 🧠 Chat Service - PsyConnect

## Overview
The **Chat Service** manages real-time messaging between users and therapists within the **PsyConnect** platform.
It is designed to support secure, persistent, and scalable chat communication.

---

## 🚧 Status
**Under Development** – Core API endpoints are defined, implementation in progress.

---

## 🧩 Features (Planned)
- 💬 Real-time one-to-one messaging between clients and therapists
- 👥 Group chat functionality
- 🕓 Message history and persistence
- 📎 File & media sharing
- 🔔 Push notifications for new messages
- 🧠 Integration with the **Therapist Recommendation** and **Profile** services

---

## ⚙️ API Endpoints

| Method | Endpoint | Description |
|--------|-----------|-------------|
| `POST` | `/chats` | Create a new chat message |
| `GET` | `/chats/:id` | Get chat by its ID |
| `GET` | `/chats/conversation/:conversationId` | Get all chats in a conversation |
| `PUT` | `/chats/:id` | Update an existing chat message |
| `DELETE` | `/chats/:id` | Delete a chat message |

---

## 🧰 Tech Stack
| Component | Technology |
|------------|-------------|
| Language | Go (Golang) |
| Framework | Gin |
| Database | MongoDB (planned) |
| Realtime | WebSocket (planned) |
| Queue | Redis / Kafka (planned) |
| Containerization | Docker |
| Deployment | Kubernetes-ready |

---

## 🚀 Running Locally

### Prerequisites
- Go 1.22+
- Docker (optional)
- MongoDB (if persistence is implemented)


## Contact
For inquiries, reach out via:
- **Project Lead**: Chessy
    - **Email**: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- **GitHub Repository**: [psyconnect-dev](https://github.com/hiamchubbybear/psyconnect-dev)
