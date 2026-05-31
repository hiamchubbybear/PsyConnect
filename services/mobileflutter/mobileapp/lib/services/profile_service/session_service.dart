import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/consultation_session.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:http/http.dart' as http;

class SessionService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<List<ConsultationSession>> getMySessions() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/sessions/me",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      final List<dynamic> data = decoded['data'] ?? [];
      return data.map((json) => ConsultationSession.fromJson(json)).toList();
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['message'] ?? "Failed to fetch consultation sessions");
    }
  }

  Future<String> getSessionPaymentUrl(String sessionId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/sessions/$sessionId/payment-url",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      final data = decoded['data'];
      if (data != null && data['payment_url'] != null) {
        return data['payment_url'].toString();
      }
      throw Exception("Payment URL not returned by server");
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['message'] ?? "Failed to get session payment URL");
    }
  }

  Future<ConsultationSession> createSession({
    required String therapistId,
    required String startTime,
    required String endTime,
    required String scheduledDate,
    required String mode,
    required double price,
  }) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "therapist_id": therapistId,
      "start_time": startTime,
      "end_time": endTime,
      "scheduled_date": scheduledDate,
      "time_zone": DateTime.now().timeZoneName,
      "mode": mode,
      "price": price,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/sessions/me",
      token: token,
      body: body,
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final decoded = jsonDecode(response.body);
      return ConsultationSession.fromJson(decoded['data']);
    } else {
      final error = jsonDecode(response.body);
      throw Exception(error['message'] ?? "Failed to book session");
    }
  }
}
