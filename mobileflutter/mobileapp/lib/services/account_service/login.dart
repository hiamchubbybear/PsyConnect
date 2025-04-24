import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/ui/screens/my_home_page.dart';
import 'package:PsyConnect/services/account_service/profile.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;

Future<http.Response> loginHandleIdentityService(String username,
    String password, String loginType, BuildContext context) async {
  try {
    Map<String, String> requestBody = {
      "username": username,
      "password": password,
    };
    final parameter = {"loginType": "NORMAL"};

    Uri requestUri =
        Uri.parse(Platform.isIOS ? loginUriIosString : loginUriAndroidString)
            .replace(queryParameters: parameter);
    final response = await http
        .post(
      requestUri,
      headers: {
        'Content-Type': 'application/json',
      },
      body: jsonEncode(requestBody),
    )
        .timeout(
      Duration(seconds: 10),
      onTimeout: () {
        throw Exception("Request timed out.");
      },
    );
    return response;
  } catch (e) {
    print("There was an error: $e");
    throw Exception("An error occurred during login.");
  }
}

class LoginService {
  ProfileService profileService = ProfileService();
  Future<void> loginHandle(
      String username,
      String password,
      String loginType,
      BuildContext context,
      String platform,
      AuthTokenProvider tokenProvider,
      UserProfileProvider profileProvider) async {
    loginType = "NORMAL";
    if (username.isEmpty || password.isEmpty) {
      ToastService.showToast(
          context: context,
          message: "Please enter both username and password before proceeding",
          title: "Success",
          type: ToastType.warning);
    } else {
      try {
        var response = await loginHandleIdentityService(
            username, password, "NORMAL", context);

        if (response.statusCode == 200) {
          ToastService.showToast(
              context: context,
              message: "Login successful",
              title: "Success",
              type: ToastType.success);
          final token = jsonDecode(response.body)["data"]["token"].toString();
          print("${token}");
          tokenProvider.setToken(token);
          profileService.getUserProfile(token: token);
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(
                builder: (context) => const MyHomePage(title: 'Home Page')),
          );
        } else if (response.statusCode == 400) {
          ToastService.showToast(
              context: context,
              message: "Wrong username or password.",
              title: "Failed",
              type: ToastType.error);
        } else if (response.statusCode == 404){

          ToastService.showToast(
              context: context,
              message:
                  "Username not found",
              title: "Not found",
              type: ToastType.error);
        } else if (response.statusCode == 500){

          ToastService.showToast(
              context: context,
              message:
                  "There was an error during login. Please try again later.",
              title: "Failed",
              type: ToastType.error);
        }
      } catch (e) {
        print("${e.toString()}");
        ToastService.showToast(
            context: context,
            message: "Occur error while login",
            title: "Failed",
            type: ToastType.error);
      }
    }
  }
}
