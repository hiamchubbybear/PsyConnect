import 'package:flutter/material.dart';


class ThemesApp {
  static final light= ThemeData(
    primaryColor: Colors.blue[200] ,
    brightness: Brightness.light
  );
  static final dark= ThemeData(
    primaryColor: Colors.yellow[100] ,
    brightness: Brightness.dark
  );
}
