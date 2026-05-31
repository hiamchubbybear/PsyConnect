import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/group.dart';
import 'package:PsyConnect/services/api/api_service.dart';

class GroupService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<List<Group>> getGroups({String? category, String? q}) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    // Build endpoint with query parameters
    String query = "";
    final List<String> params = [];
    if (category != null && category.isNotEmpty) {
      params.add("category=${Uri.encodeComponent(category)}");
    }
    if (q != null && q.isNotEmpty) {
      params.add("q=${Uri.encodeComponent(q)}");
    }
    if (params.isNotEmpty) {
      query = "?${params.join("&")}";
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/groups$query",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      // Backend might return wrapped data structure like { code: 200, data: [...] }
      if (decoded['code'] == 200 && decoded['data'] != null) {
        final List<dynamic> list = decoded['data'];
        return list.map((json) => Group.fromJson(json)).toList();
      } else if (decoded is List) {
        return decoded.map((json) => Group.fromJson(json)).toList();
      }
      return [];
    } else {
      throw Exception("Failed to load community groups");
    }
  }

  Future<Group> createGroup(Map<String, dynamic> groupData) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/groups",
      token: token,
      body: groupData,
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final decoded = jsonDecode(response.body);
      if (decoded['code'] == 200 && decoded['data'] != null) {
        return Group.fromJson(decoded['data']);
      } else {
        return Group.fromJson(decoded);
      }
    } else {
      throw Exception("Failed to create community group");
    }
  }

  Future<void> joinGroup(String groupId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/groups/$groupId/join",
      token: token,
      body: {},
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to join community group");
    }
  }
}
