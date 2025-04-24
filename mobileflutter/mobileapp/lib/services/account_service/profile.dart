import 'dart:convert';

import 'package:PsyConnect/services/api/api_service.dart';

class ProfileService {
  Future<Map<String, dynamic>> getUserProfile({required String token}) async {
    try {
      var response =
          await ApiService.getUserAccount(endpoint: "profile", token: token);
      if (response.statusCode == 200) {
        final responseDecoded = jsonDecode(response.body);
        return responseDecoded["data"];
      } else {
        final error = jsonDecode(response.body);
        throw Exception("Error ${response.statusCode}: ${error['message'] ?? 'Unknown error'}");
      }
    } catch (e) {
      throw Exception("Failed to fetch user profile: $e");
    }
  }
}
