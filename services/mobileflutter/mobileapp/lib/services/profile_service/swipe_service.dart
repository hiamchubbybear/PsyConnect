import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/swipe_card.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:http/http.dart' as http;

class SwipeService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<List<SwipeCard>> getSwipeRecommendations() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/clients/me/recommend/top",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      final List<dynamic> data = decoded['data'] ?? [];
      return data.map((json) => SwipeCard.fromJson(json)).toList();
    } else {
      throw Exception("Failed to fetch swipe recommendations");
    }
  }

  Future<void> triggerRecommendationGeneration() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/clients/me/recommend",
      token: token,
      body: {},
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to generate swipe matches");
    }
  }

  Future<void> recordSwipe(String therapistId, String status) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "therapist_id": therapistId,
      "status": status, // 'passed' or 'swiped'
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/clients/me/swipe",
      token: token,
      body: body,
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to persist swipe selection");
    }
  }
}
