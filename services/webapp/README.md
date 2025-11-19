# Web Application - PsyConnect

## Overview
The **Web Application** is the main frontend interface for the PsyConnect platform, providing users with a comprehensive web-based experience for mental health services, consultations, and resources.

## Features
- Responsive web interface
- User dashboard and profile management
- Therapist search and booking
- Real-time chat and video consultations
- Resource library and blog access
- Progressive Web App (PWA) capabilities
- Multi-language support

## Technology Stack
- **Frontend Framework**: Angular
- **UI Library**: Angular Material / Bootstrap
- **State Management**: NgRx
- **HTTP Client**: Angular HTTP Client
- **Real-time**: WebSocket / Socket.io
- **Build Tool**: Angular CLI
- **Testing**: Jasmine, Karma

## Key Features

### User Interface
- Modern, intuitive design
- Mobile-responsive layout
- Accessibility compliance (WCAG 2.1)
- Dark/light theme support
- Customizable user interface

### Core Functionality
- User registration and authentication
- Profile management and customization
- Therapist discovery and filtering
- Appointment scheduling and management
- Real-time messaging and video calls
- Payment processing integration
- Progress tracking and analytics

### Progressive Web App
- Offline functionality
- Push notifications
- App-like experience
- Installation prompts
- Service worker integration

## API Integration
The web application integrates with all PsyConnect microservices:

### Service Endpoints
- **Identity Service**: User authentication and account management
- **Profile Service**: User profiles and social features
- **Consultation Service**: Therapist and session management
- **Notification Service**: Email and push notifications
- **Chat Service**: Real-time messaging
- **Blog Service**: Content and resources
- **Search Service**: Platform-wide search
- **Recommendation Service**: AI-powered recommendations

## Setup & Development

### Prerequisites
```bash
Node.js (v16+)
npm (v8+)
Angular CLI (v15+)
```

### Environment Variables
```env
API_BASE_URL=http://localhost:8888
WEBSOCKET_URL=ws://localhost:8083
GOOGLE_CLIENT_ID={your-google-client-id}
ENVIRONMENT=development
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/webapp
   ```

2. Install dependencies:
   ```bash
   npm install
   ```

3. Start development server:
   ```bash
   ng serve
   ```

4. Build for production:
   ```bash
   ng build --prod
   ```

### Available Scripts
- `ng serve` - Start development server
- `ng build` - Build the application
- `ng test` - Run unit tests
- `ng e2e` - Run end-to-end tests
- `ng lint` - Run linting

## Project Structure
```
src/
├── app/
│   ├── components/          # Reusable components
│   ├── pages/              # Page components
│   ├── services/           # API and business logic services
│   ├── guards/             # Route guards
│   ├── interceptors/       # HTTP interceptors
│   ├── models/             # Data models and interfaces
│   └── shared/             # Shared modules and utilities
├── assets/                 # Static assets
├── environments/           # Environment configurations
└── styles/                 # Global styles
```

## Key Pages
- **Home**: Landing page and platform overview
- **Dashboard**: User dashboard and quick actions
- **Profile**: User profile management
- **Therapists**: Therapist discovery and filtering
- **Appointments**: Booking and schedule management
- **Chat**: Real-time messaging interface
- **Resources**: Blog posts and mental health resources
- **Settings**: User preferences and account settings

## Deployment
The application can be deployed to:
- Static hosting (Netlify, Vercel)
- CDN (AWS CloudFront, Azure CDN)
- Container platforms (Docker, Kubernetes)
- Traditional web servers (Apache, Nginx)

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
