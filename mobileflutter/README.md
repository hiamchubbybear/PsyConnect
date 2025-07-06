# Mobile Application - PsyConnect

## Overview
The **Mobile Application** is the cross-platform mobile interface for the PsyConnect platform, built with Flutter. It provides users with a seamless mobile experience for mental health services, consultations, and resources.

## Features
- Cross-platform compatibility (iOS & Android)
- Native mobile performance
- Offline functionality
- Push notifications
- Real-time chat and video calls
- Intuitive mobile UI/UX
- Biometric authentication
- Deep linking support

## Technology Stack
- **Framework**: Flutter
- **Language**: Dart
- **State Management**: Provider/Riverpod
- **HTTP Client**: Dio
- **Local Storage**: Hive/SQLite
- **Real-time**: WebSocket
- **Authentication**: OAuth2, JWT
- **Push Notifications**: Firebase Cloud Messaging

## Key Features

### User Experience
- Smooth animations and transitions
- Responsive design for all screen sizes
- Accessibility features
- Multi-language support
- Dark/light theme toggle

### Core Functionality
- User registration and authentication
- Profile management
- Therapist discovery with filters
- Appointment booking and management
- Real-time messaging
- Video consultation support
- Progress tracking
- Resource library access

### Mobile-Specific Features
- Push notifications for appointments
- Biometric login (fingerprint/face ID)
- Offline data synchronization
- Camera integration for profile photos
- GPS location services
- Calendar integration

## API Integration
The mobile app integrates with all PsyConnect microservices through the API Gateway:

### Service Integration
- **Identity Service**: Authentication and account management
- **Profile Service**: User profiles and social features
- **Consultation Service**: Therapist matching and sessions
- **Notification Service**: Push notifications
- **Chat Service**: Real-time messaging
- **Blog Service**: Content and resources
- **Search Service**: Platform search
- **Recommendation Service**: Personalized recommendations

## Setup & Development

### Prerequisites
```bash
Flutter SDK (v3.0+)
Dart SDK (v2.17+)
Android Studio / Xcode
```

### Environment Configuration
```dart
// lib/core/config/app_config.dart
const String baseUrl = 'http://localhost:8888';
const String websocketUrl = 'ws://localhost:8083';
const String googleClientId = 'your-google-client-id';
```

### Installation
1. Clone the repository:
   ```bash
   git clone https://github.com/hiamchubbybear/PsyConnect.git
   cd PsyConnect/mobileflutter/mobileapp
   ```

2. Install dependencies:
   ```bash
   flutter pub get
   ```

3. Run the app:
   ```bash
   flutter run
   ```

### Build Commands
- `flutter run` - Start development
- `flutter build apk` - Build Android APK
- `flutter build ios` - Build iOS app
- `flutter test` - Run unit tests
- `flutter analyze` - Code analysis

## Project Structure
```
lib/
├── core/
│   ├── config/             # App configuration
│   ├── constants/          # Constants and enums
│   ├── services/           # API and business services
│   ├── utils/              # Utility functions
│   └── theme/              # App theming
├── features/
│   ├── auth/               # Authentication feature
│   ├── profile/            # Profile management
│   ├── consultation/       # Therapist and sessions
│   ├── chat/               # Messaging feature
│   └── resources/          # Blog and resources
├── shared/
│   ├── widgets/            # Reusable widgets
│   ├── models/             # Data models
│   └── providers/          # State management
└── main.dart               # App entry point
```

## Key Screens
- **Onboarding**: Welcome and app introduction
- **Authentication**: Login/register screens
- **Dashboard**: Main user dashboard
- **Profile**: User profile management
- **Therapists**: Therapist discovery and details
- **Appointments**: Booking and schedule management
- **Chat**: Real-time messaging interface
- **Video Call**: Video consultation screen
- **Resources**: Blog and educational content
- **Settings**: App preferences and account settings

## Deployment

### Android
1. Generate signed APK:
   ```bash
   flutter build apk --release
   ```
2. Upload to Google Play Console

### iOS
1. Build iOS app:
   ```bash
   flutter build ios --release
   ```
2. Upload to App Store Connect

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
