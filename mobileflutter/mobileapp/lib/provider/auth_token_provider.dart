import 'package:flutter/material.dart';

class AuthTokenProvider with ChangeNotifier {
  String ? _token;
  String ? get token => _token;
  void setToken(String  tokenData) {
    _token = tokenData;
    notifyListeners();
  }

  String ? get getToken => _token;

  void clearToken() {
    print("Clear Token data");
    _token = null;
    notifyListeners();
  }
}
