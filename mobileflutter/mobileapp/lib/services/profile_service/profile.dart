import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/services/api/api_service.dart';

class ProfileService {
  Future<UserProfile> getUserProfile() async {
    String? userProfileRaw = await SharedPreferencesProvider().getUserProfile();
    if (userProfileRaw != null && userProfileRaw.isNotEmpty) {
      return UserProfile.fromJson(jsonDecode(userProfileRaw));
    }
    String? token = await SharedPreferencesProvider().getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }
    final response = await ApiService.getWithAccessToken(
      endpoint: "profile",
      token: token,
    );
    if (response.statusCode == 200) {
      final responseDecoded = jsonDecode(response.body)["data"];
      UserProfile userProfile = UserProfile.fromJson(responseDecoded);
      await SharedPreferencesProvider().setUserProfile(userProfile);
      return userProfile;
    } else {
      final error = jsonDecode(response.body);
      throw Exception(
          "Error ${response.statusCode}: ${error['message'] ?? 'Unknown error'}");
    }
  }
}
