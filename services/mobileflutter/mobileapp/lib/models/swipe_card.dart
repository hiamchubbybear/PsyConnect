class SwipeCard {
  final String profileId;
  final String name;
  final String address;
  final List<String> languages;
  final List<String> specialization;
  final List<String> consultationModes;
  final int experience;
  final double rating;
  final String currency;
  final double price;
  final String avatarUrl;
  final double points;
  final List<String> reasons;
  final String title;

  SwipeCard({
    required this.profileId,
    required this.name,
    required this.address,
    required this.languages,
    required this.specialization,
    required this.consultationModes,
    required this.experience,
    required this.rating,
    required this.currency,
    required this.price,
    required this.avatarUrl,
    required this.points,
    required this.reasons,
    required this.title,
  });

  factory SwipeCard.fromJson(Map<String, dynamic> json) {
    String parsedTitle = "Therapist";
    if (json['professional_info'] != null && json['professional_info']['title'] != null) {
      parsedTitle = json['professional_info']['title']['display'] ?? "Therapist";
    }

    return SwipeCard(
      profileId: json['profile_id'] ?? '',
      name: json['name'] ?? 'Mental Health Specialist',
      address: json['address'] ?? '',
      languages: List<String>.from(json['languages'] ?? []),
      specialization: List<String>.from(json['specialization'] ?? []),
      consultationModes: List<String>.from(json['consultation_modes'] ?? []),
      experience: json['experience'] ?? 0,
      rating: (json['rating'] as num?)?.toDouble() ?? 0.0,
      currency: json['currency'] ?? 'VND',
      price: (json['rage_price'] as num?)?.toDouble() ?? 0.0,
      avatarUrl: json['avatar_override'] ?? '',
      points: (json['points'] as num?)?.toDouble() ?? 0.0,
      reasons: List<String>.from(json['reasons'] ?? []),
      title: parsedTitle,
    );
  }
}
