class Post {
  final String id;
  final String title;
  final String content;
  final String authorId;
  final String? authorName;
  final String? authorAvatar;
  final List<String> tags;
  final List<String> categories;
  final int upvoteCount;
  final int downvoteCount;
  final int commentCount;
  final String createdAt;
  final String updatedAt;

  Post({
    required this.id,
    required this.title,
    required this.content,
    required this.authorId,
    this.authorName,
    this.authorAvatar,
    required this.tags,
    required this.categories,
    required this.upvoteCount,
    required this.downvoteCount,
    required this.commentCount,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Post.fromJson(Map<String, dynamic> json) {
    return Post(
      id: json['id'] ?? '',
      title: json['title'] ?? '',
      content: json['content'] ?? '',
      authorId: json['author_id'] ?? '',
      authorName: json['author_name'],
      authorAvatar: json['author_avatar'],
      tags: List<String>.from(json['tags'] ?? []),
      categories: List<String>.from(json['categories'] ?? []),
      upvoteCount: json['upvote_count'] ?? 0,
      downvoteCount: json['downvote_count'] ?? 0,
      commentCount: json['comment_count'] ?? 0,
      createdAt: json['created_at'] ?? '',
      updatedAt: json['updated_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'title': title,
      'content': content,
      'author_id': authorId,
      'author_name': authorName,
      'author_avatar': authorAvatar,
      'tags': tags,
      'categories': categories,
      'upvote_count': upvoteCount,
      'downvote_count': downvoteCount,
      'comment_count': commentCount,
      'created_at': createdAt,
      'updated_at': updatedAt,
    };
  }
}
