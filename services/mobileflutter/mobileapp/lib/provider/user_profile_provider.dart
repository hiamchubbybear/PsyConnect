import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:flutter/material.dart';

class UserProfileProvider with ChangeNotifier {
  final SharedPreferencesProvider sharedPreferencesProvider = SharedPreferencesProvider();
  Map<String, dynamic>? _profile;
  Map<String, dynamic>? get profile => _profile;
  void setProfile(Map<String, dynamic> profileData) {
    _profile = profileData;
    notifyListeners();
  }

  void clearData() {
    print("Clear user data");
    _profile = null;
    notifyListeners();
  }
}
