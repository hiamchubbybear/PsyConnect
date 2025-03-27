# Identity Service - PsyConnect - v0.1

## Overview
The **Identity Service** is responsible for user authentication, account management, and authorization within the PsyConnect platform. It provides secure and scalable identity management features, integrating OAuth2 and JWT authentication mechanisms.

## Features
- User registration and activation
- Authentication with email/password and OAuth2 providers
- Token-based authentication (JWT)
- Account management (update, delete, retrieve)
- OAuth2 authentication (Google, etc.)
- Role-based access control (RBAC)
- Account deletion confirmation

## Technology Stack
- **Backend**: Java Spring Boot
- **Database**: MySQL
- **Security**: OAuth 2.0, JWT Authentication
- **API Communication**: RESTful APIs

## API Endpoints

### Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/auth/login` | Authenticate user and generate JWT token |
| POST | `/auth/introspect` | Logout user and invalidate token |

### User Account Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/identity/create` | Register a new user account |
| POST | `/identity/activate` | Activate user account |
| POST | `/identity/req/activate` | Request activation link |
| GET | `/identity/hello` | Test API endpoint |

### Account Settings
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/account/delete` | Request account deletion |
| DELETE | `/account/delete` | Confirm account deletion |
| GET | `/account/info` | Retrieve account details |
| PUT | `/account/update/{uuid}` | Update account information |

### OAuth2 Authentication
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/oauth2/userInfo/google` | Authenticate user via Google OAuth2 |
| GET | `/oauth2/callback/google` | Google OAuth2 callback handler |

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
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/identity-service
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
- **GitHub Repository**: [PsyConnect](https://github.com/hiamchubbybear/PsyConnect)