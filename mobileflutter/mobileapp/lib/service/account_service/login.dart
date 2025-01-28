import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/screens/home/my_home_page.dart';
import 'package:PsyConnect/variable/variable.dart';
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
  Future<void> loginHandle(String username, String password, String loginType,
      BuildContext context, String platform) async {
    loginType = "NORMAL";

    if (username.isEmpty || password.isEmpty) {
      if (ScaffoldMessenger.of(context).mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(
            content: Text(
              'Please enter both username and password before proceeding',
              style: TextStyle(fontSize: 20),
            ),
          ),
        );
      }
    } else {
      try {
        var response = await loginHandleIdentityService(
            username, password, "NORMAL", context);

        if (response.statusCode == 200) {
          if (ScaffoldMessenger.of(context).mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                showCloseIcon: true,
                content:
                    Text('Login successful!', style: TextStyle(fontSize: 20)),
              ),
            );
          }
          Navigator.pushReplacement(
            context,
            MaterialPageRoute(
                builder: (context) => const MyHomePage(title: 'Home Page')),
          );
        } else if (response.statusCode == 500) {
          if (ScaffoldMessenger.of(context).mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                showCloseIcon: true,
                content: Text('Incorrect username or password!',
                    style: TextStyle(fontSize: 20)),
              ),
            );
          }
        } else {
          if (ScaffoldMessenger.of(context).mounted) {
            ScaffoldMessenger.of(context).showSnackBar(
              const SnackBar(
                showCloseIcon: true,
                content: Text(
                  'There was an error during login. Please try again later.',
                  style: TextStyle(fontSize: 20),
                ),
              ),
            );
          }
        }
      } catch (e) {
        if (ScaffoldMessenger.of(context).mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            const SnackBar(
              showCloseIcon: true,
              content: Text(
                  'An unexpected error occurred. Please try again later.',
                  style: TextStyle(fontSize: 20)),
            ),
          );
        }
      }
    }
  }
}
