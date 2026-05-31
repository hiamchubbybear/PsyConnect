class Group {
  final String id;
  final String name;
  final String description;
  final String category;
  final int memberCount;
  final bool isPublic;
  final String createdAt;

  Group({
    required this.id,
    required this.name,
    required this.description,
    required this.category,
    required this.memberCount,
    required this.isPublic,
    required this.createdAt,
  });

  factory Group.fromJson(Map<String, dynamic> json) {
    return Group(
      id: json['id'] ?? json['_id'] ?? '',
      name: json['name'] ?? '',
      description: json['description'] ?? '',
      category: json['category'] ?? '',
      memberCount: json['memberCount'] ?? json['member_count'] ?? 0,
      isPublic: json['isPublic'] ?? json['is_public'] ?? true,
      createdAt: json['createdAt'] ?? json['created_at'] ?? '',
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'name': name,
      'description': description,
      'category': category,
      'memberCount': memberCount,
      'isPublic': isPublic,
      'createdAt': createdAt,
    };
  }
}
