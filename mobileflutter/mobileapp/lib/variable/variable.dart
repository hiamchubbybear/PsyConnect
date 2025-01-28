import 'package:flutter/material.dart';

final String androidBaseUrlProfileService = "http://10.0.2.2:8080";
final String iosBaseUrlProfileService = "http://localhost:8080";
final String androidBaseUrlIdentityService = "http://10.0.2.2:8080";
final String iosBaseUrlIdentityService = "http://localhost:8080";
final String loginUriAndroidString = "http://10.0.2.2:8080/auth/login";
final String loginUriIosString = "http://localhost:8080/auth/login";
final String profileService = "http://10.0.2.2:8081";
final Color warningError = Colors.red.shade200;
final Color error = Colors.red.shade300;
final Color successStatus = Colors.blue.shade200;
final Map<String, String> headers = {
  'Content-Type': 'application/json; charset=UTF-8',
};
