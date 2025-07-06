# User Profile Service - PsyConnect

## Overview
The **User Profile Service** manages extended user profile information and additional user data within the PsyConnect platform. This service complements the main Profile Service by handling specialized profile features and extended user information.

## Features
- Extended user profile management
- Profile customization options
- User preference management
- Profile privacy settings
- Profile verification and validation
- User statistics and analytics

## Technology Stack
- **Backend**: [Technology to be implemented]
- **Database**: [Database to be implemented]
- **API Communication**: RESTful APIs
- **Security**: OAuth 2.0, JWT Authentication

## API Endpoints

### Extended Profile Management
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/userprofile/{id}` | Get extended user profile |
| POST | `/userprofile` | Create extended user profile |
| PUT | `/userprofile/{id}` | Update extended user profile |
| DELETE | `/userprofile/{id}` | Delete extended user profile |

### Profile Customization
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/userprofile/{id}/preferences` | Get user preferences |
| PUT | `/userprofile/{id}/preferences` | Update user preferences |
| GET | `/userprofile/{id}/themes` | Get available themes |
| PUT | `/userprofile/{id}/theme` | Set user theme |

### Privacy Settings
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/userprofile/{id}/privacy` | Get privacy settings |
| PUT | `/userprofile/{id}/privacy` | Update privacy settings |
| GET | `/userprofile/{id}/visibility` | Get profile visibility settings |
| PUT | `/userprofile/{id}/visibility` | Update profile visibility |

### Profile Statistics
| Method | Endpoint | Description |
|--------|----------|-------------|
| GET | `/userprofile/{id}/stats` | Get user profile statistics |
| GET | `/userprofile/{id}/activity` | Get user activity summary |
| GET | `/userprofile/{id}/achievements` | Get user achievements |

### Verification
| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/userprofile/{id}/verify` | Submit profile for verification |
| GET | `/userprofile/{id}/verification-status` | Get verification status |
| PUT | `/userprofile/{id}/verification` | Update verification documents |

*Note: This service is currently in development. API endpoints are subject to change.*

## Setup & Configuration

### Environment Variables
```env
PORT=8088
DATABASE_URL={your-database-url}
JWT_SECRET={your-jwt-secret}
VERIFICATION_SERVICE_URL={verification-service-url}
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/userprofile
   ```
2. Install dependencies and start the service (commands will be added when implemented)

## Integration
This service integrates with:
- **Profile Service**: Core profile data
- **Identity Service**: User authentication
- **Notification Service**: Profile update notifications

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
