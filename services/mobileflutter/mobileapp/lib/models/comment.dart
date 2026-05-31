class Comment {
  final String id;
  final String postId;
  final String authorId;
  final String? authorName;
  final String? authorAvatar;
  final String content;
  final String? parentCommentId;
  final String createdAt;
  final List<Comment> replies;

  Comment({
    required this.id,
    required this.postId,
    required this.authorId,
    this.authorName,
    this.authorAvatar,
    required this.content,
    this.parentCommentId,
    required this.createdAt,
    required this.replies,
  });

  factory Comment.fromJson(Map<String, dynamic> json) {
    var rawReplies = json['replies'] as List<dynamic>? ?? [];
    List<Comment> parsedReplies = rawReplies.map((r) => Comment.fromJson(r)).toList();

    return Comment(
      id: json['id'] ?? '',
      postId: json['post_id'] ?? '',
      authorId: json['author_id'] ?? '',
      authorName: json['author_name'],
      authorAvatar: json['author_avatar'],
      content: json['content'] ?? '',
      parentCommentId: json['parent_comment_id'],
      createdAt: json['created_at'] ?? '',
      replies: parsedReplies,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'post_id': postId,
      'author_id': authorId,
      'author_name': authorName,
      'author_avatar': authorAvatar,
      'content': content,
      'parent_comment_id': parentCommentId,
      'created_at': createdAt,
      'replies': replies.map((r) => r.toJson()).toList(),
    };
  }
}
