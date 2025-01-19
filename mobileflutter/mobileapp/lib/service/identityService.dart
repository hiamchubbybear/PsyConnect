import 'dart:async';
import 'dart:convert';

import 'package:PsyConnect/variable/variable.dart';
import 'package:http/http.dart' as http;
import 'package:http/http.dart';

Future<Response> loginHandleIdentityService(
    String username, String password, String loginType) async {
  try {
    Map<String, String> requestBody = {
      "username": username,
      "password": password,
    };
    final parameter = {"loginType": "NORMAL"};
    Uri requestUri =
        Uri.parse(loginUriString).replace(queryParameters: parameter);

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
