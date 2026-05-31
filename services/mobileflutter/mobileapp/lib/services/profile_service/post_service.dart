import 'dart:convert';
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/models/post.dart';
import 'package:PsyConnect/models/comment.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:http/http.dart' as http;

class PostService {
  final SharedPreferencesProvider _prefs = SharedPreferencesProvider();

  Future<List<Post>> getNewsfeed({int limit = 20, int skip = 0}) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    // Direct fetch because ApiService doesn't have query params support for basic GET
    final url = Uri.parse("$baseUrl/v1/consultation/posts?limit=$limit&skip=$skip");
    final response = await http.get(
      url,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      final List<dynamic> decoded = jsonDecode(response.body);
      return decoded.map((json) => Post.fromJson(json)).toList();
    } else {
      throw Exception("Failed to fetch newsfeed posts");
    }
  }

  Future<Post> getPostDetails(String postId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.getWithAccessToken(
      endpoint: "v1/consultation/posts/$postId",
      token: token,
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      return Post.fromJson(decoded);
    } else {
      throw Exception("Failed to fetch post details");
    }
  }

  Future<Post> createPost({
    required String title,
    required String content,
    required List<String> tags,
    required List<String> categories,
  }) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "title": title,
      "content": content,
      "tags": tags,
      "categories": categories,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/posts",
      token: token,
      body: body,
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final decoded = jsonDecode(response.body);
      return Post.fromJson(decoded);
    } else {
      throw Exception("Failed to create post");
    }
  }

  Future<void> deletePost(String postId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.deleteWithAccessToken(
      endpoint: "v1/consultation/posts/$postId",
      token: token,
    );

    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception("Failed to delete post");
    }
  }

  Future<void> addReaction(String postId, String reactionType) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "reaction_type": reactionType, // 'up' or 'down'
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/posts/$postId/react",
      token: token,
      body: body,
    );

    if (response.statusCode != 200 && response.statusCode != 201) {
      throw Exception("Failed to persist vote reaction");
    }
  }

  Future<void> removeReaction(String postId) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final response = await ApiService.deleteWithAccessToken(
      endpoint: "v1/consultation/posts/$postId/react",
      token: token,
    );

    if (response.statusCode != 200 && response.statusCode != 204) {
      throw Exception("Failed to remove reaction");
    }
  }

  Future<List<Comment>> getComments(String postId, {int limit = 20, int skip = 0}) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final url = Uri.parse("$baseUrl/v1/consultation/posts/$postId/comments?limit=$limit&skip=$skip&depth=3");
    final response = await http.get(
      url,
      headers: {
        'Content-Type': 'application/json',
        'Authorization': 'Bearer $token',
      },
    );

    if (response.statusCode == 200) {
      final decoded = jsonDecode(response.body);
      final List<dynamic> data = decoded['comments'] ?? [];
      return data.map((json) => Comment.fromJson(json)).toList();
    } else {
      throw Exception("Failed to fetch comments");
    }
  }

  Future<Comment> createComment({
    required String postId,
    required String content,
    String? parentCommentId,
  }) async {
    final token = await _prefs.getJwt();
    if (token == null || token.isEmpty) {
      throw Exception("Invalid access token");
    }

    final body = {
      "content": content,
      if (parentCommentId != null) "parent_comment_id": parentCommentId,
    };

    final response = await ApiService.postWithAccessTokenAndBody(
      endpoint: "/v1/consultation/posts/$postId/comments",
      token: token,
      body: body,
    );

    if (response.statusCode == 200 || response.statusCode == 201) {
      final decoded = jsonDecode(response.body);
      return Comment.fromJson(decoded);
    } else {
      throw Exception("Failed to write comment");
    }
  }
}
