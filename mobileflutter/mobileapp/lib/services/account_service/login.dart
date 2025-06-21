import 'dart:async';
import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:uni_links/uni_links.dart';
import 'package:url_launcher/url_launcher.dart';

class LoginService {
  ProfileService profileService = ProfileService();
  StreamSubscription<String?>? _linkSubscription;

  Future<void> oauth2LoginHandle(String provider) async {
    String oauthUrl = androidBaseUrl + "/oauth2/authorization/" + provider;

    if (await canLaunchUrl(Uri.parse(oauthUrl))) {
      await launchUrl(Uri.parse(oauthUrl),
          mode: LaunchMode.externalApplication);
    } else {
      throw 'Could not launch $oauthUrl';
    }
  }

  void initDeepLinkListener(
      BuildContext context, void Function() onLoginSuccess) {
    _linkSubscription?.cancel();

    _linkSubscription = linkStream.listen((String? link) {
      if (link != null) {
        final uri = Uri.parse(link);
        handleDeepLink(uri, context, onLoginSuccess);
      }
    });

    getInitialLink().then((String? link) {
      if (link != null) {
        final uri = Uri.parse(link);
        handleDeepLink(uri, context, onLoginSuccess);
      }
    });
  }

  void handleDeepLink(
      Uri uri, BuildContext context, void Function() onLoginSuccess) {
    if (uri.scheme == "psyconnect" && uri.path.endsWith("/callback/code")) {
      final String? code = uri.queryParameters['code'];
      final String? email = uri.queryParameters['email'];
      final String? provider = uri.queryParameters['provider'];

      if (code != null && email != null && provider != null) {
        loginHandleWithSessionCode(
          sessionCode: code,
          email: email,
          provider: provider,
          context: context,
          onLoginSuccess: () {
            onLoginSuccess();
          },
        );
      } else {
        if (context.mounted) {
          ToastService.showToast(
            context: context,
            message: "OAuth callback missing required parameters",
            title: "Error",
            type: ToastType.error,
          );
        }
      }
    } else {}
  }

  void dispose() {
    _linkSubscription?.cancel();
  }

  Future<void> loginHandleWithSessionCode({
    required String sessionCode,
    required String email,
    required String provider,
    required BuildContext context,
    required void Function() onLoginSuccess,
  }) async {
    SharedPreferencesProvider sharedPreferencesProvider =
        SharedPreferencesProvider();

    try {
      String platform = "MOBILE";
      final uri = Uri.parse(androidBaseUrl +
          "/auth/oauth2/callback/exchange-code"
              "?code=$sessionCode&email=$email&provider=$provider&platform=$platform");
      final response = await http.get(uri);
      if (response.statusCode == 200) {
        final responseBody = jsonDecode(response.body);
        final token = responseBody["data"]["token"];
        await sharedPreferencesProvider.setJwt(token);
        if (context.mounted) {
          ToastService.showToast(
            context: context,
            message: "Login successful via OAuth2",
            title: "Success",
            type: ToastType.success,
          );
        } else {}

        await Future.delayed(Duration(milliseconds: 500));

        onLoginSuccess();
      } else {
        if (context.mounted) {
          ToastService.showToast(
            context: context,
            message: "OAuth2 login failed with status ${response.statusCode}",
            title: "Login Failed",
            type: ToastType.error,
          );
        }
      }
    } catch (e) {
      if (context.mounted) {
        ToastService.showToast(
          context: context,
          message: "Exception occurred while logging in: $e",
          title: "Error",
          type: ToastType.error,
        );
      }
    }
  }

  Future<void> loginHandle(
      String username,
      String password,
      String loginType,
      BuildContext context,
      String platform,
      AuthTokenProvider tokenProvider,
      UserProfileProvider profileProvider) async {}
}

Future<http.Response> loginHandleIdentityService(
    String username,
    String password,
    String loginType,
    String platForm,
    BuildContext context) async {
  try {
    Map<String, String> requestBody = {
      "username": username,
      "password": password,
    };
    final parameter = {"loginType": "NORMAL", "platform ": "MOBILE"};

    Uri requestUri =
        Uri.parse(Platform.isIOS ? iosBaseUrl : loginUriAndroidString)
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
    throw Exception("An error occurred during login.");
  }
}
