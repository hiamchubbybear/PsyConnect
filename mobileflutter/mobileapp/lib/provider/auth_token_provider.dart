import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:flutter/foundation.dart';
import 'package:flutter/material.dart';

class AuthTokenProvider with ChangeNotifier {
  final SharedPreferencesProvider sharePref = SharedPreferencesProvider();
  String? _token;
  String? get token => _token;
  void setToken(String tokenData) {
    _token = tokenData;
    sharePref.setJwt(tokenData);
    print("Saved token to shared pref $tokenData");
    notifyListeners();
  }

  String? get getToken => _token;

  void clearToken() {
    sharePref.setJwt("");
    print("Cleared access token");
    _token = null;
    notifyListeners();
  }
}
