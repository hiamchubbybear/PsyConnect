import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/models/profile_mood.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:flutter/material.dart';

class MoodService {
  Future<void> createMoodPost(
      {required MoodModel mood, required BuildContext context}) async {
    SharedPreferencesProvider sharedPreferencesProvider =
        SharedPreferencesProvider();
        print("Mood data "+ mood.mood);
    if (mood.mood == "" || mood.mood.isEmpty) {
      ToastService.showToast(
          context: context,
          message: "Mood can not be empty",
          title: "Fill the fill",
          type: ToastType.warning);
      Navigator.of(context).pop();
      return;
    }
    String? accessToken = await sharedPreferencesProvider.getJwt();
    print("Access token is : $accessToken");
    if (accessToken == null || accessToken.isEmpty) {
      ToastService.showToast(
          context: context,
          message: "Your current session expired . Login again please.",
          title: "Session expired",
          type: ToastType.error);
    }
    final response = await ApiService.postWithAccessTokenAndBody(
        endpoint: "/mood/add",
        token: accessToken as String,
        body: mood.toJson());
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
    } else if (response.statusCode == 409) {
      ToastService.showToast(
          context: context,
          message: jsonDecode(response.body)["message"],
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

  Future<List<ProfileMoodModel>> getProfileWithMood(
      {required BuildContext context}) async {
    SharedPreferencesProvider prefs = SharedPreferencesProvider();
    String? accessToken = await prefs.getJwt();

    if (accessToken == null || accessToken.isEmpty) {
      ToastService.showToast(
        context: context,
        message: "Your current session expired. Please log in again.",
        title: "Session expired",
        type: ToastType.error,
      );
      return [];
    }

    final ownMoodResponse = await ApiService.getWithAccessToken(
        endpoint: "mood", token: accessToken);

    final friendsResponse = await ApiService.getWithAccessToken(
        endpoint: "mood/friends", token: accessToken);

    if (ownMoodResponse.statusCode != 200 &&
        friendsResponse.statusCode != 200) {
      print("Failed to fetch from server " +
          jsonDecode(ownMoodResponse.body)["message"]);
      ToastService.showToast(
        context: context,
        message: "Unable to fetch mood data",
        title: "Error",
        type: ToastType.error,
      );
      return [];
    }
    List<ProfileMoodModel> result = [];
    if (ownMoodResponse.statusCode == 200) {
      final data = jsonDecode(ownMoodResponse.body)['data'];
      if (data != null) {
        result.add(ProfileMoodModel.fromJson(data));
      }
    }
    if (friendsResponse.statusCode == 200) {
      final data = jsonDecode(friendsResponse.body)['data'];
      if (data != null) {
        result.addAll((data as List).map((e) => ProfileMoodModel.fromJson(e)));
      }
    }
    await prefs.setProfileMood(result);
    return result;
  }
}
