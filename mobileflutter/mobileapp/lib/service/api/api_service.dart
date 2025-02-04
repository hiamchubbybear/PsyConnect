import 'dart:convert';
import 'dart:io';

import 'package:PsyConnect/variable/variable.dart';
import 'package:http/http.dart' as http;

class ApiService {
  static final String _baseUrl = Platform.isAndroid
      ? androidBaseUrlProfileService
      : iosBaseUrlProfileService;

  static Future<Map<String, dynamic>> post({
    required String endpoint,
    required Map<String, String> body,
  }) async {
    final Uri uri = Uri.parse("$_baseUrl/$endpoint");
    try {
      var encodedJson = jsonEncode(body);
      final response = await http.post(
        uri,
        headers: {'Content-Type': 'application/json'},
        body: encodedJson,
      );

      print("${encodedJson}");
      if (response.statusCode == 200) {

        return jsonDecode(response.body);
      } else {
        throw Exception("Error : ${response.statusCode}");
      }
    } catch (e) {
      throw Exception("Failed to make request to backend: $e");
    }
  }
}
