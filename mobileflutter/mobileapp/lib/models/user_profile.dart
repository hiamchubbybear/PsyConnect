
class UserProfile {
  String? accountId;
  String? profileId;
  String? username;
  String? firstName;
  String? lastName;
  String? address;
  String? gender;
  String? avatarUri;
  String? description;

  UserProfile({
    this.accountId,
    this.profileId,
    this.username,
    this.firstName,
    this.lastName,
    this.address,
    this.gender,
    this.avatarUri,
    this.description,
  });

  factory UserProfile.fromJson(Map<String, dynamic> json) {
    return UserProfile(
      accountId: json["accountId"]?.toString(),
      profileId: json["profileId"]?.toString(),
      username: json["username"]?.toString(),
      firstName: json["firstName"]?.toString(),
      lastName: json["lastName"]?.toString(),
      address: json["address"]?.toString(),
      gender: json["gender"]?.toString(),
      avatarUri: json["avatarUri"]?.toString(),
      description: json["description"]?.toString(),
    );
  }

  Map<String, dynamic> toJson() {
    return {
      "accountId": accountId,
      "profileId": profileId,
      "username": username,
      "firstName": firstName,
      "lastName": lastName,
      "address": address,
      "gender": gender,
      "avatarUri": avatarUri,
      "description": description,
    };
  }

  String get getaccountId => accountId ?? "";
  String get getProfileId => profileId ?? "";
  String get getUsername => username ?? "";
  String get getFirstName => firstName ?? "";
  String get getLastName => lastName ?? "";
  String get getAddress => address ?? "";
  String get getGender => gender ?? "";
  String get getAvatarUri => avatarUri ?? "";
  String get getDescription => description ?? "";

  int countMissingFields() {
    int count = 0;
    if (accountId == null || accountId!.isEmpty) count++;
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
