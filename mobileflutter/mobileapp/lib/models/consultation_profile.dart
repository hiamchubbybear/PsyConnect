class ConsultationProfile {
  final String address;
  final List<String> languages;
  final List<String> specialization;
  final List<String> consultationModes;
  final int experience;
  final double rating;
  final double pricePerSession;
  final Availability availability;
  ConsultationProfile({
    required this.address,
    required this.languages,
    required this.specialization,
    required this.consultationModes,
    required this.experience,
    required this.rating,
    required this.pricePerSession,
    required this.availability,
  });


  factory ConsultationProfile.fromJson(Map<String, dynamic> json) {
    return ConsultationProfile(
      address: json['address'],
      languages: List<String>.from(json['languages']),
      specialization: List<String>.from(json['specialization']),
      consultationModes: List<String>.from(json['consultation_modes']),
      experience: json['experience'],
      rating: (json['rating'] as num).toDouble(),
      pricePerSession: (json['price_per_session'] as num).toDouble(),
      availability: Availability.fromJson(json['availability']),
    );
  }

  Map<String, dynamic> toJson() => {
        'address': address,
        'languages': languages,
        'specialization': specialization,
        'consultation_modes': consultationModes,
        'experience': experience,
        'rating': rating,
        'price_per_session': pricePerSession,
        'availability': availability.toJson(),
      };
}

class Availability {
  final List<String> days;
  final List<String> timeSlots;

  Availability({
    required this.days,
    required this.timeSlots,
  });

  factory Availability.fromJson(Map<String, dynamic> json) {
    return Availability(
      days: List<String>.from(json['days']),
      timeSlots: List<String>.from(json['time_slots']),
    );
  }

  Map<String, dynamic> toJson() => {
        'days': days,
        'time_slots': timeSlots,
      };
}
