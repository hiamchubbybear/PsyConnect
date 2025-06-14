import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';

final String androidBaseUrl = "http://localhost:8888";
final String iosBaseUrl = "http://localhost:8888";
final String androidBaseUrlIdentityService =
    "http://localhost:8888/identity/create";
final String iosBaseUrlIdentityService =
    "http://localhost:8888/identity/create";
final String loginUriAndroidString = "http://localhost:8888/auth/login";
final String loginUriIosString = "http://localhost:8888/auth/login";
final String profileService = "http://localhost:8888";
final Color warningError = Colors.red.shade200;
const Color blackColor = Colors.black;
const Color whiteColor = Colors.white;
final Color error = Colors.red.shade300;
final Color acceptColor=Colors.green.shade300;
final Color declineColor = Colors.green.shade300;
final Color successStatus = Colors.blue.shade200;
final Color secondaryColor = Colors.grey.shade400;
final Map<String, String> headers = {
  'Content-Type': 'application/json; charset=UTF-8',
};
ThemeProvider themeProvider = ThemeProvider();
TextStyle textStyle = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.bold,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);
TextStyle quickSand15Font = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.w600,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);

TextStyle quickSand12Font = GoogleFonts.quicksand(
  fontSize: 11,
  fontWeight: FontWeight.w400,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);

TextStyle textFieldStyle = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.w600,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);
TextStyle headingStyle = GoogleFonts.quicksand(
  fontSize: 30,
  fontWeight: FontWeight.bold,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);
TextStyle subHeadingStyle = GoogleFonts.quicksand(
  fontSize: 25,
  fontWeight: FontWeight.bold,
  color: themeProvider.isDarkMode ? Colors.white : Colors.black,
);
const double kDefault = 16.0;
const double kCircle = 20;

const TextStyle kHeadingStyle = TextStyle(
  fontSize: 24,
  fontWeight: FontWeight.bold,
  color: Colors.black,
);

const TextStyle kSubHeadingStyle = TextStyle(
  fontSize: 16,
  fontWeight: FontWeight.w500,
  color: Colors.grey,
);
const TextStyle kSubHintStyle = TextStyle(
  fontSize: 13,
  fontWeight: FontWeight.w300,
  color: Colors.grey,
);
