import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:flutter/material.dart';
import 'package:get/get.dart';

class ThemeProvider extends ChangeNotifier {
  SharedPreferencesProvider preferencesProvider = SharedPreferencesProvider();
  ThemeMode _themeMode = ThemeMode.system;

  ThemeMode get themeMode {
    print('[ThemeProvider] get themeMode: $_themeMode');
    return _themeMode;
  }

  bool get isDarkMode {
    print('[ThemeProvider] get isDarkMode: ${_themeMode == ThemeMode.dark}');
    return _themeMode == ThemeMode.dark;
  }

  ThemeProvider() {
    print('[ThemeProvider] Constructor called');
    _loadThemeMode();
  }

  Future<void> _loadThemeMode() async {
    print('[ThemeProvider] _loadThemeMode() called');
    bool isDark = await preferencesProvider.isDarkMode();
    print('[ThemeProvider] isDark from SharedPreferences: $isDark');
    _themeMode = isDark ? ThemeMode.dark : ThemeMode.light;
    print('[ThemeProvider] _themeMode set to: $_themeMode');
    Get.changeThemeMode(_themeMode);
    notifyListeners();
  }

  void toggleTheme(bool isOn) {
    print('[ThemeProvider] toggleTheme() called with isOn: $isOn');
    _themeMode = isOn ? ThemeMode.dark : ThemeMode.light;
    print('[ThemeProvider] _themeMode updated to: $_themeMode');
    preferencesProvider.setThemeMode(isOn);
    Get.changeThemeMode(_themeMode);
    notifyListeners();
  }
}
