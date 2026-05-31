import 'package:flutter/material.dart';
import 'package:PsyConnect/models/consultation_session.dart';
import 'package:PsyConnect/services/profile_service/session_service.dart';

class SessionProvider with ChangeNotifier {
  final SessionService _sessionService = SessionService();

  List<ConsultationSession> _sessions = [];
  List<ConsultationSession> get sessions => _sessions;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  Future<void> fetchSessions() async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      _sessions = await _sessionService.getMySessions();
      _sessions.sort((a, b) => b.createdAt.compareTo(a.createdAt)); // Sort by newest
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      _sessions = [];
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<String?> getPaymentUrl(String sessionId) async {
    try {
      return await _sessionService.getSessionPaymentUrl(sessionId);
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      notifyListeners();
      return null;
    }
  }

  Future<bool> bookSession({
    required String therapistId,
    required String startTime,
    required String endTime,
    required String scheduledDate,
    required String mode,
    required double price,
  }) async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      await _sessionService.createSession(
        therapistId: therapistId,
        startTime: startTime,
        endTime: endTime,
        scheduledDate: scheduledDate,
        mode: mode,
        price: price,
      );
      await fetchSessions(); // Reload sessions
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
