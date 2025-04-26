import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:flutter/material.dart';

class MoodService {
  Future<void> createMoodPost(
      {required MoodModel mood, required BuildContext context}) async {
    SharedPreferencesProvider sharedPreferencesProvider =
        SharedPreferencesProvider();
    String? accessToken = await sharedPreferencesProvider.getJwt();
    print("Access token is : $accessToken");
    if (accessToken == null || accessToken.isEmpty)
      ToastService.showToast(
          context: context,
          message: "Your current session expired . Login again please.",
          title: "Session expired",
          type: ToastType.error);
    final response = await ApiService.postWithAccessTokenAndBody(
        endpoint: "/mood/add", token: accessToken as String, body: mood.toJson());
    if (response.statusCode == 404) {
      ToastService.showToast(
          context: context,
          message: response.body,
          title: "Failed",
          type: ToastType.error);
          return;
    } else if (response.statusCode == 500) {
      ToastService.showToast(
          context: context,
          message: "Failed to connect with server",
          title: "Server error",
          type: ToastType.error);
          return;
    }
    ToastService.showToast(
        context: context,
        message: "Create mood success",
        title: "Success",
        type: ToastType.success);
        Navigator.of(context).pop();
  }
}
