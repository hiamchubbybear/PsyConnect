import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/swipe_card.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:http/http.dart' as http;

class TherapistService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<SwipeCard> getTherapistById(String therapistId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/therapists/$therapistId",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      final data = decoded['data'] ?? decoded;
      return SwipeCard.fromJson(data);
    } else {
      throw Exception("Failed to load therapist details");
    }
  }
}
