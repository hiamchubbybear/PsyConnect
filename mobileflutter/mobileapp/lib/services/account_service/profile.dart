import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/services/api/api_service.dart';

class ProfileService {
  Future<UserProfile> getUserProfile() async {
    String? token = await SharedPreferencesProvider().getUserProfile();
    if (token == null || token.isEmpty) throw Exception("Invalid access token");
    var response =
        await ApiService.getWithAccessToken(endpoint: "profile", token: token);
    if (response.statusCode == 200) {
      final responseDecoded = jsonDecode(response.body);
      SharedPreferencesProvider().setUserProfile(responseDecoded["data"]);
    } else {
      final error = jsonDecode(response.body);
      throw Exception(
          "Error ${response.statusCode}: ${error['message'] ?? 'Unknown error'}");
    }
    final String? profileString =
        await SharedPreferencesProvider().getUserProfile();
    if (profileString == null || profileString.isEmpty) {
      throw new Exception("Failed to create get user profile");
    }
    final userProfile = UserProfile.fromJson(jsonDecode(profileString));
    return userProfile;
  }
}
