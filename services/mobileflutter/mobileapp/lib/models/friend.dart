class Friend {
  final String profileId;
  final String firstName;
  final String lastName;
  final String avatarUri;
  final String? role;

  Friend({
    required this.profileId,
    required this.firstName,
    required this.lastName,
    required this.avatarUri,
    this.role,
  });

  factory Friend.fromJson(Map<String, dynamic> json) {
    return Friend(
      profileId: json['profileId'] ?? json['profile_id'] ?? '',
      firstName: json['firstName'] ?? json['first_name'] ?? '',
      lastName: json['lastName'] ?? json['last_name'] ?? '',
      avatarUri: json['avatarUri'] ?? json['avatar_uri'] ?? '',
      role: json['role'],
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'profileId': profileId,
      'firstName': firstName,
      'lastName': lastName,
      'avatarUri': avatarUri,
      'role': role,
    };
  }

  String get fullName => "$firstName $lastName".trim();
}
