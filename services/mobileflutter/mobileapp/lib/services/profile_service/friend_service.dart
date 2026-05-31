import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/friend.dart';
import 'package:PsyConnect/services/api/api_service.dart';

class FriendService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<List<Friend>> getMyFriends() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/profile/friends",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      if (decoded['code'] == 200 && decoded['data'] != null) {
        final List<dynamic> list = decoded['data'];
        return list.map((json) => Friend.fromJson(json)).toList();
      }
      return [];
    } else {
      throw Exception("Failed to load friends list");
    }
  }

  Future<List<Friend>> getReceivedRequests() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/profile/friends/received",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      if (decoded['code'] == 200 && decoded['data'] != null) {
        final List<dynamic> list = decoded['data'];
        return list.map((json) => Friend.fromJson(json)).toList();
      }
      return [];
    } else {
      throw Exception("Failed to load received friend requests");
    }
  }

  Future<List<Friend>> getFriendSuggestions() async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/profile/friends/suggestions",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      if (decoded['code'] == 200 && decoded['data'] != null) {
        final List<dynamic> list = decoded['data'];
        return list.map((json) => Friend.fromJson(json)).toList();
      }
      return [];
    } else {
      throw Exception("Failed to load friend suggestions");
    }
  }

  Future<void> sendFriendRequest(String targetId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "target": targetId,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/profile/friend/request",
      token: token,
      body: body,
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to send friend request invitation");
    }
  }

  Future<void> acceptFriendRequest(String targetId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "target": targetId,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/profile/friend/accept",
      token: token,
      body: body,
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to accept friend request");
    }
  }

  Future<void> unfriend(String targetId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "target": targetId,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/profile/friend/unfriend",
      token: token,
      body: body,
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to unfriend target profile");
    }
  }
}
