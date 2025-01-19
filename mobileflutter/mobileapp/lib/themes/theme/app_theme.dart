import 'package:PsyConnect/themes/colors/color.dart';
import 'package:flutter/material.dart';


class AppTheme {
  static ThemeData get lightTheme => ThemeData(
        colorScheme: ColorScheme.fromSeed(seedColor: primaryColor),
        useMaterial3: true,
      );

  static ThemeData get darkTheme => ThemeData.dark();
}
