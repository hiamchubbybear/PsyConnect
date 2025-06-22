import 'dart:async';
import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/auth_token_provider.dart';
import 'package:PsyConnect/provider/user_profile_provider.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:PsyConnect/ui/screens/my_home_page.dart';
import 'package:flutter/material.dart';
import 'package:http/http.dart' as http;
import 'package:uni_links/uni_links.dart';
import 'package:url_launcher/url_launcher.dart';

class LoginService {
  ProfileService profileService = ProfileService();
  StreamSubscription<String?>? _linkSubscription;

  Future<void> oauth2LoginHandle(String provider) async {
    String oauthUrl = baseUrl + "/oauth2/authorization/" + provider;

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
      final uri = Uri.parse(baseUrl +
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
      String provider,
      BuildContext context,
      String platform,
      AuthTokenProvider tokenProvider,
      UserProfileProvider profileProvider) async {
    if (username.isEmpty || password.isEmpty) {
      ToastService.showToast(
          context: context,
          message: "Please enter both username and password before proceeding",
          title: "Success",
          type: ToastType.warning);
    } else {
      try {
        var response = await loginHandleIdentityService(
            username, password, provider, platform, context);

        if (response.statusCode == 200) {
          ToastService.showToast(
              context: context,
              message: "Login successful",
              title: "Success",
              type: ToastType.success);
          final token = jsonDecode(response.body)["data"]["token"].toString();
          print("${token}");
          tokenProvider.setToken(token);
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
        } else if (response.statusCode == 404) {
          ToastService.showToast(
              context: context,
              message: "Username not found",
              title: "Not found",
              type: ToastType.error);
        } else if (response.statusCode == 500) {
          print(response.body);
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
    final parameter = {"provider": "NORMAL", "platform": "MOBILE"};

    Uri requestUri =
        Uri.parse(baseUrl+"/auth/login")
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
