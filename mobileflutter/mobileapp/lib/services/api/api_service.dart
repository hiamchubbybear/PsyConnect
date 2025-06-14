import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/core/variable/variable.dart';
import 'package:http/http.dart' as http;

class ApiService {
  static final String _baseUrl =
      Platform.isAndroid ? androidBaseUrl : iosBaseUrl;

  static Future<http.Response> post({
    required String endpoint,
    required Map<String, dynamic> body,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl/$endpoint");
    var encodedJson = jsonEncode(body);

    try {
      return await http.post(
        uri,
        headers: {'Content-Type': 'application/json'},
        body: encodedJson,
      );
    } catch (e) {
      throw Exception("Failed to connect to backend: $e");
    }
  }

  static Future<http.Response> postWithParameters({
    required String endpoint,
    required String param,
    required String paramValue,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl/$endpoint?$param=$paramValue");
    try {
      return await http.post(uri);
    } catch (e) {
      throw new Exception("There are some error $e");
    }
  }

  static Future<http.Response> getWithAccessToken({
    required String endpoint,
    required String token,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl/$endpoint");
    try {
      return await http.get(
        uri,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );
    } catch (e) {
      throw Exception("Failed to connect to backend: $e");
    }
  }

  static Future<http.Response> postWithAccessTokenAndBody({
    required String endpoint,
    required String token,
    required Map<String, dynamic> body,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl$endpoint");
    try {
      return await http.post(
        uri,
        body: jsonEncode(body),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );
    } catch (e) {
      throw Exception("Failed to connect to backend: $e");
    }
  }

  static Future<http.Response> putWithTokenAndBody({
    required String? token,
    required String? endpoint,
    required Map<String, dynamic> body,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl/$endpoint");
    try {
      return await http.put(
        uri,
        body: jsonEncode(body),
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );
    } catch (e) {
      throw Exception("Failed to connect to backend: $e");
    }
  }
}
