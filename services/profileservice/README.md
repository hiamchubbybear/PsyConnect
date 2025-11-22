# Profile Service - PsyConnect

![Java](https://img.shields.io/badge/Java-17-orange?style=flat&logo=java&logoColor=white)
![Spring Boot](https://img.shields.io/badge/Spring_Boot-3.0-green?style=flat&logo=spring-boot&logoColor=white)
![Neo4j](https://img.shields.io/badge/Neo4j-5-blue?style=flat&logo=neo4j&logoColor=white)

## 📖 Overview

The **Profile Service** manages rich user profiles and social connections. Utilizing a Graph Database (Neo4j), it efficiently handles complex relationships like friendships, followers, and professional networks between clients and therapists.

## ✨ Features

- **Profile Management**: Extended user details (Bio, Avatar, Specializations).
- **Social Graph**: Manage friends, followers, and connection requests.
- **Mood Tracking**: Log and track daily mood status.
- **Settings**: Manage user preferences and privacy settings.

## 🛠 Technology Stack

- **Framework**: Java Spring Boot 3
- **Database**: Neo4j (Graph DB)
- **Communication**: REST API, Kafka (Event Consumer)
- **Build Tool**: Maven

## 🔌 API Endpoints

### Profile
| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/v1/profile/me` | Get my profile |
| `PUT` | `/v1/profile/me` | Update my profile |
| `GET` | `/v1/profile/{id}` | Get public profile of user |

### Social
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/v1/profile/friend/request` | Send friend request |
| `POST` | `/v1/profile/friend/accept` | Accept friend request |
| `GET` | `/v1/profile/friends` | List friends |

### Mood
| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/v1/mood/me/add` | Log current mood |
| `GET` | `/v1/mood/me` | Get mood history |

## ⚙️ Configuration

### Environment Variables
This service uses a `.env` file for configuration.
1. Copy the example file:
   ```bash
   cp .env.example .env
   ```
2. Update the variables in `.env`:
   - `SPRING_NEO4J_AUTHENTICATION_PASSWORD`: Neo4j password
   - `SERVER_PORT`: Service port (default: 8081)

## 🚀 Installation & Run

### Prerequisites
- Java JDK 17+
- Maven
- Neo4j Database

### Local Run
```bash
# 1. Navigate to directory
cd services/profileservice

# 2. Install dependencies
mvn clean install

# 3. Run application
mvn spring-boot:run
```

### Docker Run
```bash
docker build -t psyconnect/profileservice .
docker run -p 8081:8081 --env-file .env psyconnect/profileservice
```

## 🤝 Contributing
Please refer to the root [README](../../README.md) for contributing guidelines.
