import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

class ThemeProvider extends ChangeNotifier {
  SharedPreferencesProvider preferencesProvider = SharedPreferencesProvider();
  ThemeMode _themeMode = ThemeMode.system;

  ThemeMode get themeMode {
    return _themeMode;
  }

  bool get isDarkMode {
    return _themeMode == ThemeMode.dark;
  }

  ThemeProvider() {
    _loadThemeMode();
  }

  Future<void> _loadThemeMode() async {
    bool isDark = await preferencesProvider.isDarkMode();
    _themeMode = isDark ? ThemeMode.dark : ThemeMode.light;
    Get.changeThemeMode(_themeMode);
    notifyListeners();
  }

  void toggleTheme(bool isOn) {
    _themeMode = isOn ? ThemeMode.dark : ThemeMode.light;
    preferencesProvider.setThemeMode(isOn);
    Get.changeThemeMode(_themeMode);
    notifyListeners();
  }
}
