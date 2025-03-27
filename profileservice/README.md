# Profile Service - PsyConnect - v0.1

## Overview
The **Profile Service** is a core microservice in the PsyConnect platform, responsible for managing user profiles, settings, activity logs, and moods. It provides APIs for handling user data securely and efficiently. 
  
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
- **Security**: OAuth 2.0, JWT Authentication

## API Endpoints

### User Profile Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/profile/internal/user` | Create a new user profile (Internal use) |
| POST | `/profile` | Update an existing user profile |
| GET | `/profile` | Retrieve user profile details |
| GET | `/profile/all` | Get paginated list of user profiles |
| GET | `/profile/friends` | Retrieve user friends and their moods |

### User Settings
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/user-setting/add` | Add new user settings |
| GET | `/user-setting` | Retrieve user settings |
| PUT | `/user-setting` | Update user settings |
| DELETE | `/user-setting` | Delete user settings |
| POST | `/user-setting/default` | Reset user settings to default |

### Activity Logs
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/activity-log/{profileId}` | Add an activity log |
| GET | `/activity-log/{profileId}` | Retrieve activity logs for a profile |
| DELETE | `/activity-log/{logId}` | Delete an activity log |

### Mood Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/mood/add` | Add mood status |
| GET | `/mood` | Get current mood status |
| PUT | `/mood` | Update mood status |
| DELETE | `/mood` | Delete mood status |

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
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/profile-service
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

