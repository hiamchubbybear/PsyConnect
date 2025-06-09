import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/setting.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:http/http.dart' as http;

class SettingService {
  Future<Setting> getSetting() async {
    final provider = SharedPreferencesProvider();
    final jwtTokenData = await provider.getJwt();
    if (jwtTokenData == null || jwtTokenData.isEmpty) {
      throw Exception("Access token is missing.");
    }

    final http.Response res = await ApiService.getWithAccessToken(
      endpoint: "user-setting",
      token: jwtTokenData,
    );

    if (res.statusCode == 200) {
      final jsonResponse = jsonDecode(res.body);
      print("Loaded setting JSON: ${jsonResponse["data"].toString()}");
      return Setting.fromJson(jsonResponse["data"]);
    } else if (res.statusCode >= 400 && res.statusCode < 500) {
      throw Exception("Client error: ${res.statusCode} ${res.body}");
    } else if (res.statusCode >= 500) {
      throw Exception("Server error: ${res.statusCode}");
    } else {
      throw Exception("Unexpected error: ${res.statusCode}");
    }
  }
}
