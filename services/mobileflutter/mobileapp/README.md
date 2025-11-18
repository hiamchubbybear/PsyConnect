# PsyConnect Mobile App

## Overview
A Flutter-based cross-platform mobile application for the PsyConnect mental health platform, providing users with seamless access to mental health services, consultations, and resources.

## Features
- Cross-platform compatibility (iOS & Android)
- User authentication and profile management
- Therapist discovery and booking
- Real-time chat and video consultations
- Push notifications
- Offline data synchronization
- Biometric authentication

## Technology Stack
- **Framework**: Flutter
- **Language**: Dart  
- **State Management**: Provider
- **HTTP Client**: Custom API Service
- **Authentication**: JWT, OAuth2

## Getting Started

### Prerequisites
- Flutter SDK (v3.0+)
- Dart SDK (v2.17+)
- Android Studio / Xcode

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

## API Integration
The app integrates with the PsyConnect backend through:
- API Gateway: `http://localhost:8888`
- WebSocket: Real-time communication
- JWT Authentication: Secure API access

## Project Structure
```
lib/
├── services/api/           # API communication
├── screens/               # App screens
├── widgets/               # Reusable widgets
├── models/                # Data models
├── providers/             # State management
└── core/                  # Core utilities
```

## Build Commands
- `flutter run` - Development mode
- `flutter build apk` - Android release
- `flutter build ios` - iOS release
- `flutter test` - Run tests

## Contributing
We welcome contributions! Please follow the standard Git workflow and ensure code quality with proper testing.

## Resources
- [Flutter Documentation](https://docs.flutter.dev/)
- [API Documentation](../README.md)
- [PsyConnect Platform](https://github.com/hiamchubbybear/PsyConnect)

## Contact
- **Project Lead**: Chessy
- **Email**: [tranvanhuy16032004@gmail.com](mailto:tranvanhuy16032004@gmail.com)
