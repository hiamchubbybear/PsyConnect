class Setting {
  final bool showLastSeen;
  final bool showProfilePicture;
  final bool showMood;
  final bool notificationsEnabled;
  final bool emailNotifications;
  final bool pushNotifications;
  final bool smsNotifications;
  final bool twoFactorAuth;
  final bool allowLoginAlerts;
  final String trustedDevices;
  final bool autoDeleteOldMoods;
  String language;
  String theme;

  Setting({
    required this.showLastSeen,
    required this.showProfilePicture,
    required this.showMood,
    required this.notificationsEnabled,
    required this.emailNotifications,
    required this.pushNotifications,
    required this.smsNotifications,
    required this.twoFactorAuth,
    required this.allowLoginAlerts,
    required this.trustedDevices,
    required this.autoDeleteOldMoods,
    this.language = 'en',
    this.theme = 'light',
  });

  factory Setting.fromJson(Map<String, dynamic> json) => Setting(
        showLastSeen: json['showLastSeen'] ?? true,
        showProfilePicture: json['showProfilePicture'] ?? true,
        showMood: json['showMood'] ?? false,
        notificationsEnabled: json['notificationsEnabled'] ?? false,
        emailNotifications: json['emailNotifications'] ?? false,
        pushNotifications: json['pushNotifications'] ?? false,
        smsNotifications: json['smsNotifications'] ?? false,
        twoFactorAuth: json['twoFactorAuth'] ?? false,
        allowLoginAlerts: json['allowLoginAlerts'] ?? false,
        trustedDevices: json['trustedDevices'] ?? '',
        autoDeleteOldMoods: json['autoDeleteOldMoods'] ?? true,
        language: json['language'] ?? 'en',
        theme: json['theme'] ?? 'light',
      );

  Map<String, dynamic> toJson() => {
        'showLastSeen': showLastSeen,
        'showProfilePicture': showProfilePicture,
        'showMood': showMood,
        'notificationsEnabled': notificationsEnabled,
        'emailNotifications': emailNotifications,
        'pushNotifications': pushNotifications,
        'smsNotifications': smsNotifications,
        'twoFactorAuth': twoFactorAuth,
        'allowLoginAlerts': allowLoginAlerts,
        'trustedDevices': trustedDevices,
        'language': language,
        'theme': theme,
        'autoDeleteOldMoods': autoDeleteOldMoods,
      };

  Setting copyWith({
    bool? showLastSeen,
    bool? showProfilePicture,
    bool? showMood,
    bool? notificationsEnabled,
    bool? emailNotifications,
    bool? pushNotifications,
    bool? smsNotifications,
    bool? twoFactorAuth,
    bool? allowLoginAlerts,
    String? trustedDevices,
    String? language,
    String? theme,
    bool? autoDeleteOldMoods,
  }) {
    return Setting(
      showLastSeen: showLastSeen ?? this.showLastSeen,
      showProfilePicture: showProfilePicture ?? this.showProfilePicture,
      showMood: showMood ?? this.showMood,
      notificationsEnabled: notificationsEnabled ?? this.notificationsEnabled,
      emailNotifications: emailNotifications ?? this.emailNotifications,
      pushNotifications: pushNotifications ?? this.pushNotifications,
      smsNotifications: smsNotifications ?? this.smsNotifications,
      twoFactorAuth: twoFactorAuth ?? this.twoFactorAuth,
      allowLoginAlerts: allowLoginAlerts ?? this.allowLoginAlerts,
      trustedDevices: trustedDevices ?? this.trustedDevices,
      language: language ?? this.language,
      theme: theme ?? this.theme,
      autoDeleteOldMoods: autoDeleteOldMoods ?? this.autoDeleteOldMoods,
    );
  }
}
