# Consultative Service - PsyConnect

## Overview

The **Consultative Service** manages therapist and client information related to mental health consultations in the
PsyConnect ecosystem. It is designed to handle profile registration, updates, therapist status management, matching, and session management.

---

## Features

- Register, update, and retrieve **Therapist** and **Client** profiles
- Change therapist status (e.g., available/unavailable)
- Manage consultation sessions
- Handle therapist-client matching
- Recommendation system for therapist matching

---

## API Endpoints

### Therapist Management

| Method | Endpoint                                | Description                      | Role Required               | Headers Required |
|--------|-----------------------------------------|----------------------------------|-----------------------------|------------------|
| GET    | `/consultation/therapist`               | Get therapist info               | role.therapist:permission   | X-Roles          |
| POST   | `/consultation/therapist`               | Create therapist profile         | role.therapist:permission   | X-Roles          |
| PUT    | `/consultation/therapist`               | Update therapist profile         | role.therapist:permission   | X-Roles          |
| PUT    | `/consultation/therapist/status/{status}` | Change therapist status          | role.therapist:permission   | X-Roles          |

### Client Management

| Method | Endpoint                                | Description                      | Role Required               | Headers Required |
|--------|-----------------------------------------|----------------------------------|-----------------------------|------------------|
| GET    | `/consultation/client`                  | Get client info                  | role.client:permission      | X-Roles          |
| POST   | `/consultation/client`                  | Create client profile            | role.client:permission      | X-Roles          |
| PUT    | `/consultation/client`                  | Update client profile            | role.client:permission      | X-Roles          |
| POST   | `/consultation/client/recommend`        | Trigger recommendation update    | role.client:permission      | X-Roles          |
| GET    | `/consultation/client/recommend/top`    | Get top 5 recommended therapists | role.client:permission      | X-Roles          |
| POST   | `/consultation/client/match`            | Request therapist match          | role.client:permission      | X-Roles          |

### Matching Management

| Method | Endpoint                                | Description                      | Role Required               | Headers Required |
|--------|-----------------------------------------|----------------------------------|-----------------------------|------------------|
| GET    | `/consultation/therapist/match`         | Get all matching therapists      | None                        | None             |
| POST   | `/consultation/therapist/match`         | Request therapist match          | None                        | None             |

### Session Management

| Method | Endpoint                                | Description                      | Role Required               | Headers Required |
|--------|-----------------------------------------|----------------------------------|-----------------------------|------------------|
| GET    | `/consultation/admin/session`           | Get all sessions (admin)         | role.admin:permission       | X-Roles          |
| GET    | `/consultation/session/all`             | Get all sessions by user ID      | None (any authenticated)    | X-Roles          |
| POST   | `/consultation/session`                 | Create new session               | None (any authenticated)    | X-Roles          |
| DELETE | `/consultation/session`                 | Delete current session           | None (any authenticated)    | X-Roles          |
| GET    | `/consultation/session/{id}`            | Get session by ID                | None (any authenticated)    | X-Roles          |

---

## Technology Stack

- **Language**: Go (Gin Framework)
- **Communication**: RESTful API
- **Runtime**: Docker / Kubernetes-ready

---

## Setup & Run

1. Clone the repository:

```bash
git clone https://github.com/hiamchubbybear/psyconnect-dev.git
cd psyconnect-dev/consultationservice
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
- **GitHub**: [PsyConnect](https://github.com/hiamchubbybear/psyconnect-dev)

---

## Contributing

We welcome PRs and collaboration:

1. Fork the repo
2. Create a new branch (`feature/your-feature`)
3. Commit your changes
4. Push and create a Pull Request
