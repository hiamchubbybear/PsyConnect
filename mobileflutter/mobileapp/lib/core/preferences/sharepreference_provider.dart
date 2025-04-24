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

  Future<void> setJwt(String token) async =>
      (await SharedPreferences.getInstance()).setString(_accessTokenKey, token);

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

  // Đọc trạng thái theme mode
  Future<bool> isDarkMode() async {
    final prefs = await SharedPreferences.getInstance();
    return prefs.getBool(_themeModeKey) ??
        false; // mặc định là false nếu không có dữ liệu
  }

  Future<void> clearAll() async =>
      (await SharedPreferences.getInstance()).clear();
}
