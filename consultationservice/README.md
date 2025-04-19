# Consultative Service - PsyConnect

## Overview

The **Consultative Service** manages therapist and client information related to mental health consultations in the
PsyConnect ecosystem. It is designed to handle profile registration, updates, and therapist status management in a
structured and secure way.

---

## Features

- Register, update, and retrieve **Therapist** and **Client** profiles
- Change therapist status (e.g., available/unavailable)
- JWT-based authorization
- Works with API Gateway and user identity headers

---

## API Endpoints

### Therapist Management

| Method | Endpoint                         | Description              |
|--------|----------------------------------|--------------------------|
| GET    | `/consultation/therapist`        | Get therapist info       |
| POST   | `/consultation/therapist`        | Create therapist profile |
| PUT    | `/consultation/therapist`        | Update therapist profile |
| PUT    | `/consultation/therapist:status` | Change therapist status  |

### Client Management

| Method | Endpoint               | Description           |
|--------|------------------------|-----------------------|
| GET    | `/consultation/client` | Get client info       |
| POST   | `/consultation/client` | Create client profile |
| PUT    | `/consultation/client` | Update client profile |

---

## Authentication

All endpoints require a valid **JWT token** in the request header.

### Request Format

**Header:**

```http
Authorization: Bearer <your_jwt_token>
```

### Example with `curl`

```bash
curl -X GET http://localhost:8080/consultation/therapist \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  -H "Content-Type: application/json"
```

### Example with Postman

1. Go to the **Authorization** tab.
2. Select **Bearer Token**.
3. Paste your token in the token field.

### Headers Added by API Gateway

| Header         | Description                   |
|----------------|-------------------------------|
| `X-User-Id`    | User account ID               |
| `X-Profile-Id` | Profile ID linked to the user |

These headers are automatically extracted and injected from the JWT by the API Gateway.

---

## Technology Stack

- **Language**: Go (Gin Framework)
- **Auth**: JWT via API Gateway
- **Communication**: RESTful API
- **Runtime**: Docker / Kubernetes-ready

---

## Setup & Run

1. Clone the repository:

```bash
git clone https://github.com/hiamchubbybear/PsyConnect.git
cd PsyConnect/consultative-service
```

2. Build & Run the service:

```bash
go build -o consultative-service .
./consultative-service
```

Or with Docker:

```bash
docker build -t psyconnect/consultative-service .
docker run -p 8084:8084 psyconnect/consultative-service
```

---

## Contact

For questions or contributions, contact:

- **Project Lead**: Chessy
- **Email**: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
- **GitHub**: [PsyConnect](https://github.com/hiamchubbybear/PsyConnect)

---

## Contributing

We welcome PRs and collaboration:

1. Fork the repo
2. Create a new branch (`feature/your-feature`)
3. Commit your changes
4. Push and create a Pull Request
