import 'dart:convert';

import 'package:PsyConnect/models/user_profile.dart';
import 'package:shared_preferences/shared_preferences.dart';

class SharedPreferencesProvider {
  static final SharedPreferencesProvider _instance =
      SharedPreferencesProvider._internal();

  factory SharedPreferencesProvider() {
    return _instance;
  }
  SharedPreferencesProvider._internal();

  static const String _accessTokenKey = "accessToken";
  static const String _userIdKey = "userId";
  static const String _usernameKey = "username";
  static const String _emailKey = "email";
  static const String _roleKey = "role";
  static const String _themeModeKey = "isDarkMode";
  static const String _userProfile = "userProfile";

  Future<void> setJwt(String token) async =>
      (await SharedPreferences.getInstance()).setString(_accessTokenKey, token);
  Future<void> setUserProfile(UserProfile userProfile) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setString(_userProfile, jsonEncode(userProfile.toJson()));
  }

  Future<void> setUserId(String userId) async =>
      (await SharedPreferences.getInstance()).setString(_userIdKey, userId);

  Future<void> setUsername(String username) async =>
      (await SharedPreferences.getInstance()).setString(_usernameKey, username);

  Future<void> setEmail(String email) async =>
      (await SharedPreferences.getInstance()).setString(_emailKey, email);

  Future<void> setRole(String role) async =>
      (await SharedPreferences.getInstance()).setString(_roleKey, role);

  Future<void> setThemeMode(bool isDarkMode) async {
    final prefs = await SharedPreferences.getInstance();
    await prefs.setBool(_themeModeKey, isDarkMode);
  }

  Future<String?> getJwt() async =>
      (await SharedPreferences.getInstance()).getString(_accessTokenKey);

  Future<String?> getUserId() async =>
      (await SharedPreferences.getInstance()).getString(_userIdKey);

  Future<String?> getUsername() async =>
      (await SharedPreferences.getInstance()).getString(_usernameKey);

  Future<String?> getEmail() async =>
      (await SharedPreferences.getInstance()).getString(_emailKey);

  Future<String?> getRole() async =>
      (await SharedPreferences.getInstance()).getString(_roleKey);
  Future<bool> isDarkMode() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(_themeModeKey) ?? false;
  }

  Future<String?> getUserProfile() async =>
      (await SharedPreferences.getInstance()).getString(_userProfile);
  Future<void> clearAll() async =>
      (await SharedPreferences.getInstance()).clear();
}
