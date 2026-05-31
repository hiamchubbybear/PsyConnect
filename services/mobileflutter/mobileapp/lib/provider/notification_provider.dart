import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/variable/variable.dart';

class SystemNotification {
  final String id;
  final String title;
  final String message;
  final String type;
  final bool isRead;
  final String createdAt;

  SystemNotification({
    required this.id,
    required this.title,
    required this.message,
    required this.type,
    required this.isRead,
    required this.createdAt,
  });

  factory SystemNotification.fromJson(Map<String, dynamic> json) {
    return SystemNotification(
      id: json['id'] ?? json['_id'] ?? '',
      title: json['title'] ?? '',
      message: json['message'] ?? json['content'] ?? '',
      type: json['type'] ?? 'info',
      isRead: json['isRead'] ?? json['is_read'] ?? false,
      createdAt: json['createdAt'] ?? json['created_at'] ?? '',
    );
  }
}

class NotificationProvider with ChangeNotifier {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  List<SystemNotification> _notifications = [];
  List<SystemNotification> get notifications => _notifications;

  int _unreadCount = 0;
  int get unreadCount => _unreadCount;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  Future<void> fetchNotifications(String userId) async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final token = await _prefs.getJwt();
      if (token == null) throw Exception("Unauthorized session");

      final url = Uri.parse("$baseUrl/notification/$userId");
      final response = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        final List<dynamic> data = decoded is List ? decoded : (decoded['data'] ?? []);
        
        _notifications = data.map((json) => SystemNotification.fromJson(json)).toList();
        _unreadCount = _notifications.where((n) => !n.isRead).length;
      } else {
        throw Exception("Failed to load notifications history");
      }
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> markAsRead(String notificationId) async {
    try {
      final token = await _prefs.getJwt();
      if (token == null) throw Exception("Unauthorized session");

      final url = Uri.parse("$baseUrl/notification/$notificationId/read");
      final response = await http.put(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200 || response.statusCode == 204) {
        // Toggle locally
        final index = _notifications.indexWhere((n) => n.id == notificationId);
        if (index != -1) {
          final n = _notifications[index];
          _notifications[index] = SystemNotification(
            id: n.id,
            title: n.title,
            message: n.message,
            type: n.type,
            isRead: true,
            createdAt: n.createdAt,
          );
          _unreadCount = _notifications.where((n) => !n.isRead).length;
          notifyListeners();
        }
        return true;
      }
      return false;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
