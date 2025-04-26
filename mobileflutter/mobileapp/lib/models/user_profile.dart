class UserProfile {
  String? username;
  String? firstName;
  String? lastName;
  String? address;
  String? gender;
  String? avatarUri;
  String? description;
  UserProfile(
      {this.username,
      this.firstName,
      this.lastName,
      this.address,
      this.gender,
      this.avatarUri,
      this.description});
  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      username: json["username"] as String,
      firstName: json["firstName"] as String,
      lastName: json["lastName"] as String,
      address: json["address"] as String,
      gender: json["gender"] as String,
      avatarUri: json["avatarUri"] as String,
      description: json["description"] as String,
    );
  }
  Map<String, dynamic> toJson() {
    return {
      "username": username,
      "firstName": firstName,
      "lastName": lastName,
      "address": address,
      "gender": gender,
      "avatarUri": avatarUri,
      "description": description,
    };
  }

  String get getUsername => username ?? "";
  String get getFirstName => firstName ?? "";
  String get getLastName => lastName ?? "";
  String get getAddress => address ?? "";
  String get getGender => gender ?? "";
  String get getAvatarUri => avatarUri ?? "";
  String get getDescription => description ?? "";

  int countMissingFields() {
    int count = 0;
    if (username == null || username!.isEmpty) count++;
    if (firstName == null || firstName!.isEmpty) count++;
    if (lastName == null || lastName!.isEmpty) count++;
    if (address == null || address!.isEmpty) count++;
    if (gender == null || gender!.isEmpty) count++;
    if (avatarUri == null || avatarUri!.isEmpty) count++;
    if (description == null || description!.isEmpty) count++;
    return count;
  }
}
