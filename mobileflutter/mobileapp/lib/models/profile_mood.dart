class ProfileMoodModel {
  final String profileId;
  final String fullName;
  final String avatarUri;
  final String moodId;
  final String mood;
  final String moodDescription;
  final String visibility;
  final int createdAt;
  final int expiresAt;

  ProfileMoodModel({
    required this.profileId,
    required this.fullName,
    required this.avatarUri,
    required this.moodId,
    required this.mood,
    required this.moodDescription,
    required this.visibility,
    required this.createdAt,
    required this.expiresAt,
  });

  factory ProfileMoodModel.fromJson(Map<String, dynamic> json) {
    return ProfileMoodModel(
      profileId: json['profileId'] as String,
      fullName: json['fullName'] as String,
      avatarUri: json['avatarUrl'] as String,
      moodId: json['moodId'] as String,
      mood: json['mood'] as String,
      moodDescription: json['description'] as String,
      visibility: json['visibility'] as String,
      createdAt: json['createdAt'] as int,
      expiresAt: json['expiresAt'] as int,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'profileId': profileId,
      'fullName': fullName,
      'avatarUrl': avatarUri,
      'moodId': moodId,
      'mood': mood,
      'description': moodDescription,
      'visibility': visibility,
      'createdAt': createdAt,
      'expiresAt': expiresAt,
    };
  }
}
