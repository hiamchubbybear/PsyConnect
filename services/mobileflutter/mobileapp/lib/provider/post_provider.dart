import 'package:flutter/material.dart';
import 'package:PsyConnect/models/post.dart';
import 'package:PsyConnect/models/comment.dart';
import 'package:PsyConnect/services/profile_service/post_service.dart';

class PostProvider with ChangeNotifier {
  final PostService _postService = PostService();

  List<Post> _posts = [];
  List<Post> get posts => _posts;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  // Active post comments
  List<Comment> _activeComments = [];
  List<Comment> get activeComments => _activeComments;
  
  bool _isLoadingComments = false;
  bool get isLoadingComments => _isLoadingComments;

  // Cache user votes: postId -> 'up' | 'down' | null
  final Map<String, String?> _userVotes = {};
  String? getUserVote(String postId) => _userVotes[postId];

  Future<void> fetchPosts({bool refresh = true}) async {
    _isLoading = true;
    _errorMessage = null;
    if (refresh) {
      _posts = [];
    }
    notifyListeners();

    try {
      final fetched = await _postService.getNewsfeed(
        limit: 20,
        skip: refresh ? 0 : _posts.length,
      );
      if (refresh) {
        _posts = fetched;
      } else {
        _posts.addAll(fetched);
      }
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> fetchComments(String postId) async {
    _isLoadingComments = true;
    _activeComments = [];
    notifyListeners();

    try {
      _activeComments = await _postService.getComments(postId);
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
    } finally {
      _isLoadingComments = false;
      notifyListeners();
    }
  }

  Future<bool> publishPost({
    required String title,
    required String content,
    required List<String> tags,
    required List<String> categories,
  }) async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final newPost = await _postService.createPost(
        title: title,
        content: content,
        tags: tags,
        categories: categories,
      );
      _posts.insert(0, newPost);
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> toggleVote(String postId, String voteType) async {
    final currentVote = _userVotes[postId];

    try {
      if (currentVote == voteType) {
        // Remove reaction
        await _postService.removeReaction(postId);
        _userVotes[postId] = null;
        
        // Update local counts
        final index = _posts.indexWhere((p) => p.id == postId);
        if (index != -1) {
          final p = _posts[index];
          _posts[index] = Post(
            id: p.id,
            title: p.title,
            content: p.content,
            authorId: p.authorId,
            authorName: p.authorName,
            authorAvatar: p.authorAvatar,
            tags: p.tags,
            categories: p.categories,
            upvoteCount: voteType == 'up' ? p.upvoteCount - 1 : p.upvoteCount,
            downvoteCount: voteType == 'down' ? p.downvoteCount - 1 : p.downvoteCount,
            commentCount: p.commentCount,
            createdAt: p.createdAt,
            updatedAt: p.updatedAt,
          );
        }
      } else {
        // Add or change reaction
        await _postService.addReaction(postId, voteType);
        
        final index = _posts.indexWhere((p) => p.id == postId);
        if (index != -1) {
          final p = _posts[index];
          int upDiff = 0;
          int downDiff = 0;

          if (currentVote == 'up') upDiff = -1;
          if (currentVote == 'down') downDiff = -1;

          if (voteType == 'up') upDiff += 1;
          if (voteType == 'down') downDiff += 1;

          _posts[index] = Post(
            id: p.id,
            title: p.title,
            content: p.content,
            authorId: p.authorId,
            authorName: p.authorName,
            authorAvatar: p.authorAvatar,
            tags: p.tags,
            categories: p.categories,
            upvoteCount: p.upvoteCount + upDiff,
            downvoteCount: p.downvoteCount + downDiff,
            commentCount: p.commentCount,
            createdAt: p.createdAt,
            updatedAt: p.updatedAt,
          );
        }
        _userVotes[postId] = voteType;
      }
      notifyListeners();
      return true;
    } catch (e) {
      print("Vote error: $e");
      return false;
    }
  }

  Future<bool> addComment({
    required String postId,
    required String content,
    String? parentCommentId,
  }) async {
    _errorMessage = null;

    try {
      final newComment = await _postService.createComment(
        postId: postId,
        content: content,
        parentCommentId: parentCommentId,
      );

      if (parentCommentId == null) {
        // Root comment
        _activeComments.add(newComment);
      } else {
        // Reply: find parent and insert
        _insertReply(_activeComments, parentCommentId, newComment);
      }
      
      // Update comment count on post list
      final index = _posts.indexWhere((p) => p.id == postId);
      if (index != -1) {
        final p = _posts[index];
        _posts[index] = Post(
          id: p.id,
          title: p.title,
          content: p.content,
          authorId: p.authorId,
          authorName: p.authorName,
          authorAvatar: p.authorAvatar,
          tags: p.tags,
          categories: p.categories,
          upvoteCount: p.upvoteCount,
          downvoteCount: p.downvoteCount,
          commentCount: p.commentCount + 1,
          createdAt: p.createdAt,
          updatedAt: p.updatedAt,
        );
      }
      
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      notifyListeners();
      return false;
    }
  }

  void _insertReply(List<Comment> list, String parentId, Comment newReply) {
    for (var i = 0; i < list.length; i++) {
      if (list[i].id == parentId) {
        list[i].replies.add(newReply);
        return;
      }
      if (list[i].replies.isNotEmpty) {
        _insertReply(list[i].replies, parentId, newReply);
      }
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
