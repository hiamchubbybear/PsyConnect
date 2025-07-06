# Identity Service - PsyConnect - v0.1

## Overview
The **Identity Service** is responsible for user authentication, account management, and authorization within the PsyConnect platform. It provides secure and scalable identity management features.

## Features
- User registration and activation
- User authentication and logout
- Account management (update, delete, retrieve)
- Role-based access control (RBAC)
- Account deletion confirmation

## Technology Stack
- **Backend**: Java Spring Boot
- **Database**: MySQL
- **API Communication**: RESTful APIs

## API Endpoints

### Authentication

| Method | Endpoint                | Description                      | Role Required               | Headers Required |
|--------|-------------------------|----------------------------------|-----------------------------|------------------|
| POST   | `/auth/login`           | Authenticate user and get token  | None                        | None             |
| POST   | `/auth/introspect`      | Logout user and invalidate token | None                        | None             |
| POST   | `/auth/internal/valid`  | Internal token validation        | None                        | None             |

### User Account Management

| Method | Endpoint                | Description                      | Role Required               | Headers Required |
|--------|-------------------------|----------------------------------|-----------------------------|------------------|
| POST   | `/identity/create`      | Register a new user account      | None                        | None             |
| POST   | `/identity/activate`    | Activate user account            | None                        | None             |
| POST   | `/identity/req/activate`| Request activation link          | None                        | None             |
| GET    | `/identity/hello`       | Test API endpoint                | None                        | None             |

### Account Settings

| Method | Endpoint                | Description                      | Role Required               | Headers Required |
|--------|-------------------------|----------------------------------|-----------------------------|------------------|
| POST   | `/account/delete`       | Request account deletion         | None                        | X-User-Id        |
| DELETE | `/account/delete`       | Confirm account deletion         | None                        | X-User-Id        |
| GET    | `/account/info`         | Retrieve account details         | None                        | X-User-Id        |
| PUT    | `/account/update`       | Update account information       | None                        | X-User-Id        |
| GET    | `/account/all/{page}`   | Get paginated list of accounts   | role.admin:permission       | X-Roles          |

---

## Setup & Configuration
### Environment Variables
Ensure the following environment variables are set before running the service:
```env
MYSQL_SPRING_DATASOURCE_PASSWORD={your-variable}
MYSQL_SPRING_DATASOURCE_USERNAME={your-variable}
SPRING_DATASOURCE_URL={your-variable}
SERVER_PORT={your-variable}
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/psyconnect-dev.git
   cd psyconnect-dev/identityservice
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