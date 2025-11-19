# Profile Service

## Overview

The **Profile Service** is a core microservice in the PsyConnect platform, responsible for managing user profiles,
settings, activity logs, and moods. It provides APIs for handling user data securely and efficiently.

## Features

- User profile creation and management
- User settings configuration
- Activity log tracking
- Mood status updates
- Friend relationships and profile interactions

## Technology Stack

- **Backend**: Java Spring Boot
- **Database**: Neo4j, MySQL
- **API Communication**: RESTful APIs

## 📌 Version History

| Version    | Date       | Changes                                             |
|------------|------------|-----------------------------------------------------|
| **v1.1.0** | 2025-08-19 | - Refactor new endpoint "/me" into myself endpoints |
| **v1.0.0** | 2025-08-18 | - Initial release                                   |

## API Endpoints - v0.1

### User Profile Management

| Method | Endpoint                     | Description                           | Role Required         | Headers Required |
|--------|------------------------------|---------------------------------------|-----------------------|------------------|
| POST   | ~~`/profile/internal/user`~~ | Create a new user profile (Internal)  | None                  | None             |
| PUT    | ~~`/profile`~~               | Update an existing user profile       | None                  | X-Profile-Id     |
| GET    | ~~`/profile`~~               | Retrieve user profile details         | None                  | X-Profile-Id     |
| GET    | `/profile/all`               | Get paginated list of user profiles   | role.admin:permission | X-Roles          |
| GET    | ~~`/profile/friends`~~       | Retrieve user friends and their moods | None                  | X-Profile-Id     |

### User Settings

| Method | Endpoint                    | Description                    | Role Required | Headers Required |
|--------|-----------------------------|--------------------------------|---------------|------------------|
| GET    | ~~`/user-setting`~~         | Retrieve user settings         | None          | X-Profile-Id     |
| PUT    | ~~`/user-setting`~~         | Update user settings           | None          | X-Profile-Id     |
| POST   | ~~`/user-setting/default`~~ | Reset user settings to default | None          | X-Profile-Id     |

### Mood Management

| Method | Endpoint            | Description                | Role Required | Headers Required |
|--------|---------------------|----------------------------|---------------|------------------|
| POST   | ~~`/mood/me/add`~~  | Add mood status            | None          | X-Profile-Id     |
| GET    | ~~`/mood`~~         | Get current mood status    | None          | X-Profile-Id     |
| PUT    | ~~`/mood/me`~~      | Update mood status         | None          | X-Profile-Id     |
| DELETE | ~~`/mood/me`~~      | Delete mood status         | None          | X-Profile-Id     |
| GET    | ~~`/mood/friends`~~ | Get friends' mood statuses | None          | X-Profile-Id     |

---

## API Endpoints - v1

## User Profile Management

| Method | Endpoint               | Description                         | Role Required         | Headers Required |
|--------|------------------------|-------------------------------------|-----------------------|------------------|
| POST   | `/v1/profile/internal/user` | Create a new user profile (Internal)  | None                  | None             |
| PUT    | `/v1/profile/me`       | Update an existing user profile     | None                  | X-Profile-Id     |
| GET    | `/v1/profile/me`       | Retrieve user profile details       | None                  | X-Profile-Id     |
| GET    | `/v1/profile/all`      | Get paginated list of user profiles | role.admin:permission | X-Roles          |
| GET    | `/v1/profile/friends`  | Retrieve user friends and their moods | None                | X-Profile-Id     |

---

## User Settings

| Method | Endpoint                | Description                    | Role Required | Headers Required |
|--------|-------------------------|--------------------------------|---------------|------------------|
| GET    | `/v1/user-setting/me`   | Retrieve user settings         | None          | X-Profile-Id     |
| PUT    | `/v1/user-setting/me`   | Update user settings           | None          | X-Profile-Id     |
| POST   | `/v1/user-setting/default` | Reset user settings to default | None        | X-Profile-Id     |

---

## Mood Management

| Method | Endpoint               | Description                | Role Required | Headers Required |
|--------|------------------------|----------------------------|---------------|------------------|
| POST   | `/v1/mood/me/add`      | Add mood status            | None          | X-Profile-Id     |
| GET    | `/v1/mood/me`          | Get current mood status    | None          | X-Profile-Id     |
| PUT    | `/v1/mood/me`          | Update mood status         | None          | X-Profile-Id     |
| DELETE | `/v1/mood/me`          | Delete mood status         | None          | X-Profile-Id     |
| GET    | `/v1/mood/me/friends`  | Get friends' mood statuses | None          | X-Profile-Id     |

---

## Friend Management (Social)

| Method | Endpoint                   | Description                          | Role Required | Headers Required |
|--------|----------------------------|--------------------------------------|---------------|------------------|
| POST   | `/v1/profile/friend/request` | Send a friend request                | None          | X-Profile-Id     |
| POST   | `/v1/profile/friend/accept`  | Accept a friend request              | None          | X-Profile-Id     |
| POST   | `/v1/profile/friend/unfriend`| Remove a friend                      | None          | X-Profile-Id     |
| POST   | `/v1/profile/friend/undo`    | Cancel (undo) a sent friend request  | None          | X-Profile-Id     |
| DELETE | `/v1/profile/friend/request` | Decline a friend request             | None          | X-Profile-Id     |
****
---

## Setup & Configuration

### Environment Variables

Ensure the following environment variables are set before running the service:

```env
NEO4J_PASSWORD_PROFILE_SERVICE={your-variable}
MYSQL_SPRING_DATASOURCE_PASSWORD={your-variable}
MYSQL_SPRING_DATASOURCE_USERNAME={your-variable}
SPRING_DATASOURCE_URL={your-variable}
SERVER_PORT={your-variable}
```

### Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/psyconnect-dev.git
   cd psyconnect-dev/profileservice
   ```
2. Build and run the service:
   ```bash
   mvn clean install
   mvn spring-boot:run
   ```

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
- **GitHub Repository**: [psyconnect-dev](https://github.com/hiamchubbybear/psyconnect-dev)
