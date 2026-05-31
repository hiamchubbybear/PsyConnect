import 'dart:convert';
import 'package:firebase_core/firebase_core.dart';
import 'package:firebase_messaging/firebase_messaging.dart';
import 'package:flutter/foundation.dart';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/services/api/api_service.dart';

class NotificationService {
  static final NotificationService _instance = NotificationService._internal();
  factory NotificationService() => _instance;
  NotificationService._internal();

  final FirebaseMessaging _fcm = FirebaseMessaging.instance;
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<void> initializeNotificationHooks() async {
    try {
      // 1. Request notifications authorization permission from users
      NotificationSettings settings = await _fcm.requestPermission(
        alert: true,
        announcement: false,
        badge: true,
        carPlay: false,
        criticalAlert: false,
        provisional: false,
        sound: true,
      );

      if (settings.authorizationStatus == AuthorizationStatus.authorized) {
        if (kDebugMode) {
          print('User granted authorization permission for push notifications.');
        }

        // 2. Retrieve FCM registration device token
        String? token = await _fcm.getToken();
        if (token != null) {
          await _syncFcmTokenWithServer(token);
        }

        // 3. Monitor token rotations
        _fcm.onTokenRefresh.listen((newToken) async {
          await _syncFcmTokenWithServer(newToken);
        });

        // 4. Handle foreground notifications
        FirebaseMessaging.onMessage.listen((RemoteMessage message) {
          if (kDebugMode) {
            print('Received a foreground notification message: ${message.messageId}');
          }
          // We can parse or trigger a local custom overlay notification banner here
        });

        // 5. Handle user tapping notifications to open the app
        FirebaseMessaging.onMessageOpenedApp.listen((RemoteMessage message) {
          if (kDebugMode) {
            print('User tapped notification to open app: ${message.data}');
          }
          // Dynamic router triggers will go here
        });
      }
    } catch (e) {
      if (kDebugMode) {
        print("Error initializing FCM messaging push notifications: $e");
      }
    }
  }

  Future<void> _syncFcmTokenWithServer(String fcmToken) async {
    try {
      final token = await _prefs.getJwt();
      if (token == null || token.isEmpty) return;

      final body = {
        "token": fcmToken,
      };

      // gateways /notification/token maps to notificationservice
      final response = await ApiService.postWithAccessTokenAndBody(
        endpoint: "/notification/token",
        token: token,
        body: body,
      );

      if (response.statusCode == 200) {
        if (kDebugMode) {
          print("Successfully registered device FCM Token with central notification server");
        }
      }
    } catch (e) {
      if (kDebugMode) {
        print("Failed to sync FCM Token with backend gateway: $e");
      }
    }
  }
}
