import 'package:flutter/material.dart';

class UserProvider with ChangeNotifier {
  Map<String, dynamic>? _user;
  Map<String, dynamic>? get user => _user;
  void setUser(Map<String, dynamic> userData) {
    _user = userData;
    notifyListeners();
  }

  Map<String, dynamic>? get getData => _user?["data"];

  void clearUser() {
    print("Clear user data");
    _user = null;
    notifyListeners();
  }
}
