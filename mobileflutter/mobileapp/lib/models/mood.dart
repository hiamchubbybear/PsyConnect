class MoodModel {
  final String mood;
  final String moodDescription;
  final String visibility;
  MoodModel({
    required this.mood,
    required this.moodDescription,
    required this.visibility,
  });
  factory MoodModel.fromJson(Map<String, dynamic> json) {
    return MoodModel(
      mood: json['mood'] as String,
      moodDescription: json['moodDescription'] as String,
      visibility: json['visibility'] as String,
    );
  }
  Map<String, dynamic> toJson() {
    return {
      'mood': mood,
      'moodDescription': moodDescription,
      'visibility': visibility,
    };
  }
}
