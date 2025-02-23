import 'dart:convert';

import 'package:PsyConnect/page/verified_page.dart';
import 'package:PsyConnect/provider/user_provider.dart';
import 'package:PsyConnect/service/api/api_service.dart';
import 'package:PsyConnect/toasting&loading/toast.dart';
import 'package:flutter/material.dart';

class RegisterService {
  Future<Map<String, dynamic>?> registerHandle({
    required Map<String, String> requestBody,
    required UserProvider userProvider,
    required BuildContext context,
  }) async {
    try {
      final response = await ApiService.post(
        endpoint: "identity/create",
        body: requestBody,
      );

      if (response.statusCode == 200) {
        final responseDecoded = jsonDecode(response.body);
        final userData = responseDecoded["data"];
        if (responseDecoded is! Map<String, dynamic>) {
          throw Exception("Invalid response format");
        }
        userProvider.setUser(userData);
        debugPrint("Email from UserProvider: ${userProvider.user?["email"]}");
        ToastService.showToast(
            message: "Register successful!",
            context: context,
            title: 'Notification',
            type: ToastType.success);
        if (!context.mounted) return responseDecoded;
        Navigator.pushReplacement(
          context,
          MaterialPageRoute(builder: (context) => const VerifiedPage()),
        );

        return responseDecoded;
      }
      switch (response.statusCode) {
        case 409:
          ToastService.showToast(
              message: "username or email already used!",
              context: context,
              title: 'Exists field',
              type: ToastType.error);
          throw Exception("Username or email already used");
        case 500:
          throw Exception("Server error: ${response.statusCode}");
        default:
          ToastService.showToast(
              message: "Un error ocurred!",
              context: context,
              title: 'Error',
              type: ToastType.error);
          throw Exception("Unexpected error: ${response.statusCode}");
      }
    } catch (e) {
      debugPrint("RegisterService Error: $e");
      ToastService.showToast(
          message: "There are an expected error with server}!",
          context: context,
          title: 'Unknown error',
          type: ToastType.error);
      return null;
    }
  }
}
