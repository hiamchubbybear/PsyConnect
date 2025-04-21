import 'package:flutter/material.dart';
import 'package:get/get.dart';
import 'package:google_fonts/google_fonts.dart';

final String androidBaseUrl = "http://10.0.2.2:8888";
final String iosBaseUrl = "http://localhost:8888";
final String androidBaseUrlIdentityService =
    "http://10.0.2.2:8888/identity/create";
final String iosBaseUrlIdentityService =
    "http://localhost:8888/identity/create";
final String loginUriAndroidString = "http://10.0.2.2:8888/auth/login";
final String loginUriIosString = "http://localhost:8888/auth/login";
final String profileService = "http://10.0.2.2:8888";
final Color warningError = Colors.red.shade200;
const Color blackColor = Colors.black;
const Color whiteColor = Colors.white;
final Color error = Colors.red.shade300;
final Color successStatus = Colors.blue.shade200;
final Color secondaryColor= Colors.grey.shade400;
final Map<String, String> headers = {
  'Content-Type': 'application/json; charset=UTF-8',
};
TextStyle textStyle = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.bold,
    color: Get.isDarkMode ? Colors.white : Colors.black,
);
TextStyle quickSand15Font = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.w600,
    color: Get.isDarkMode ? Colors.white : Colors.black,
);

TextStyle quickSand12Font = GoogleFonts.quicksand(
  fontSize: 12,
  fontWeight: FontWeight.w600,
  color: Get.isDarkMode ? Colors.white : Colors.black,
);

TextStyle textFieldStyle = GoogleFonts.quicksand(
  fontSize: 17,
  fontWeight: FontWeight.w600,
    color: Get.isDarkMode ? Colors.white : Colors.black,
);
