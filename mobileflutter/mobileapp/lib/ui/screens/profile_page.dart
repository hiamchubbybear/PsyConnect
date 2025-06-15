import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/route/route_animation.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:PsyConnect/ui/screens/login_page.dart';
import 'package:PsyConnect/ui/screens/consultation_profile_page.dart';
import 'package:PsyConnect/ui/screens/setting_page.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({Key? key}) : super(key: key);

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  ProfileService profileService = ProfileService();
  UserProfile userProfile = UserProfile();

  // Tạo các handler functions riêng biệt cho từng button
  void handleSetProfileDetails(BuildContext context) {
    print("Handle Set Profile Details");
    // TODO: Navigate to profile details page
    handleOnProfile(context); // Tạm thời dùng handleOnProfile
  }

  void handleUploadResume(BuildContext context) {
    print("Handle Upload Resume");
    // TODO: Navigate to upload resume page
    handleOnProfile(context); // Tạm thời dùng handleOnProfile
  }

  void handleAddSkills(BuildContext context) {
    print("Handle Add Skills");
    // TODO: Navigate to add skills page
    handleOnProfile(context); // Tạm thời dùng handleOnProfile
  }

  // profileCompletionCards với các handler riêng biệt
  List<ProfileCompletionCard> get profileCompletionCards => [
    ProfileCompletionCard(
      title: "Set Your Profile Details",
      icon: CupertinoIcons.person_circle,
      buttonText: "Continue",
      onTap: handleSetProfileDetails,
    ),
    ProfileCompletionCard(
      title: "Upload your resume",
      icon: CupertinoIcons.doc,
      buttonText: "Upload",
      onTap: handleUploadResume,
    ),
    ProfileCompletionCard(
      title: "Add your skills",
      icon: CupertinoIcons.square_list,
      buttonText: "Add",
      onTap: handleAddSkills,
    ),
  ];
  @override
  void initState() {
    super.initState();
    loadUserData();
  }

  void loadUserData() async {
    try {
      final data = await SharedPreferencesProvider().getUserProfile();

      if (data == null) {
        await _fetchAndSetUserProfileFromApi();
        return;
      }

      final Map<String, dynamic> jsonMap = jsonDecode(data);

      final user = UserProfile.fromJson(jsonMap);

      if (user.accountId == null || user.username == null) {
        await _fetchAndSetUserProfileFromApi();
        return;
      }

      setState(() {
        userProfile = user;
      });
    } catch (e, stacktrace) {
      ToastService.showToast(
        context: context,
        message: "Your current session is expired. Please login again! $e",
        title: "Failed",
        type: ToastType.error,
      );
    }
  }

  Future<void> _fetchAndSetUserProfileFromApi() async {
    try {
      final userFromApi = await profileService.getUserProfile();
      setState(() {
        userProfile = userFromApi;
        print(jsonEncode(userProfile.toJson()));
      });
    } catch (e) {
      ToastService.showToast(
        context: context,
        message: "Failed to load profile. Please try again later.",
        title: "Error",
        type: ToastType.error,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        elevation: 0,
        backgroundColor: Colors.transparent,
        foregroundColor: Colors.black,
        title: Text(
          "Profile",
          style: quickSand15Font,
        ),
        centerTitle: true,
        actions: [
          IconButton(
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(builder: (context) => const SettingsPage()),
              );
            },
            icon: const Icon(Icons.settings_rounded),
          )
        ],
      ),
      body: ListView(
        padding: const EdgeInsets.fromLTRB(10, 10, 10, 50),
        children: [
          Column(
            children: [
              GestureDetector(
                child: CircleAvatar(
                  radius: 50,
                  backgroundImage: NetworkImage(userProfile
                          .getAvatarUri.isNotEmpty
                      ? userProfile.getAvatarUri
                      : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg'),
                ),
                onTap: () => handleOnProfile(context),
              ),
              const SizedBox(height: 10),
              Text(
                (userProfile.getFirstName.isNotEmpty)
                    ? userProfile.getFirstName
                    : "Meo",
                style: subHeadingStyle,
              ),
              Text((userProfile.getDescription.isNotEmpty)
                  ? userProfile.getDescription
                  : "Meorapist")
            ],
          ),
          const SizedBox(height: 25),
          FutureBuilder<int>(
            future: checkComplete(userProfile),
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return const SizedBox();
              } else if (snapshot.hasError) {
                return Text('Error: ${snapshot.error}');
              } else if (snapshot.hasData) {
                int value = snapshot.data!;
                if (value > 0) {
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          const Padding(
                            padding: EdgeInsets.only(right: 5),
                            child: Text(
                              "Complete your profile",
                              style: TextStyle(
                                fontWeight: FontWeight.bold,
                              ),
                            ),
                          ),
                          Text(
                            "($value/8)",
                            style: TextStyle(
                              color: successStatus,
                            ),
                          )
                        ],
                      ),
                      const SizedBox(height: 10),
                      Row(
                        children: List.generate(5, (index) {
                          return Expanded(
                            child: Container(
                              height: 7,
                              margin:
                                  EdgeInsets.only(right: index == 4 ? 0 : 6),
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(10),
                                color: index < (7 - value)
                                    ? successStatus
                                    : Colors.black12,
                              ),
                            ),
                          );
                        }),
                      ),
                    ],
                  );
                } else {
                  return const SizedBox();
                }
              } else {
                return const SizedBox();
              }
            },
          ),
          const SizedBox(height: 10),
          SizedBox(
            height: 180,
            child: ListView.separated(
              physics: const BouncingScrollPhysics(),
              scrollDirection: Axis.horizontal,
              itemBuilder: (context, index) {
                final card = profileCompletionCards[index];
                return SizedBox(
                  width: 160,
                  child: Card(
                    shadowColor: Colors.black12,
                    child: Padding(
                      padding: const EdgeInsets.all(15),
                      child: Column(
                        children: [
                          Icon(
                            card.icon,
                            size: 30,
                          ),
                          const SizedBox(height: 10),
                          Text(
                            card.title,
                            textAlign: TextAlign.center,
                          ),
                          const Spacer(),
                          ElevatedButton(
                            onPressed: () {
                              print("Button pressed for: ${card.title}");
                              if (card.onTap != null) {
                                print("Executing onTap function");
                                card.onTap!(context);
                              } else {
                                print("onTap is null for: ${card.title}");
                              }
                            },
                            style: ElevatedButton.styleFrom(
                              elevation: 0,
                              shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(10)),
                            ),
                            child: Text(card.buttonText),
                          )
                        ],
                      ),
                    ),
                  ),
                );
              },
              separatorBuilder: (context, index) =>
                  const Padding(padding: EdgeInsets.only(right: 5)),
              itemCount: profileCompletionCards.length,
            ),
          ),
          const SizedBox(height: 35),
          ...List.generate(
            customListTiles.length,
            (index) {
              final tile = customListTiles[index];
              return Padding(
                padding: const EdgeInsets.only(bottom: 5),
                child: Card(
                  elevation: 4,
                  shadowColor: Colors.black12,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: InkWell(
                    borderRadius: BorderRadius.circular(12),
                    onTap: () => tile.onTap?.call(context),
                    child: ListTile(
                      leading: Icon(tile.icon),
                      title: Text(tile.title),
                      trailing: const Icon(Icons.chevron_right),
                    ),
                  ),
                ),
              );
            },
          )
        ],
      ),
      bottomNavigationBar: BottomNavigationBar(
        currentIndex: 3,
        type: BottomNavigationBarType.fixed,
        items: const [
          BottomNavigationBarItem(
            icon: Icon(CupertinoIcons.home),
            label: "Home",
          ),
          BottomNavigationBarItem(
            icon: Icon(CupertinoIcons.chat_bubble_2),
            label: "Messages",
          ),
          BottomNavigationBarItem(
            icon: Icon(CupertinoIcons.book),
            label: "Discover",
          ),
          BottomNavigationBarItem(
            icon: Icon(CupertinoIcons.person),
            label: "Profile",
          ),
        ],
      ),
    );
  }
}

void handleOnProfile(BuildContext context) {
  print("handleOnProfile called - navigating to ConsultationProfilePage");
  Navigator.push(
      context, createSlideFromBottomRoute(ConsultationProfilePage()));
}

class ProfileCompletionCard {
  final String title;
  final String buttonText;
  final IconData icon;
  final void Function(BuildContext context)? onTap;
  ProfileCompletionCard({
    required this.title,
    required this.buttonText,
    required this.icon,
    required this.onTap,
  });
}

// profileCompletionCards đã được di chuyển vào trong class _ProfilePageState

class CustomListTile {
  final IconData icon;
  final String title;
  final void Function(BuildContext context)? onTap;
  CustomListTile({required this.icon, required this.title, this.onTap});
}

Future<int> checkComplete(UserProfile userProfile) async {
  return userProfile.countMissingFields();
}

List<CustomListTile> customListTiles = [
  CustomListTile(
    icon: Icons.insights,
    title: "Activity",
    onTap: (context) {},
  ),
  CustomListTile(
    icon: Icons.history,
    title: "History",
    onTap: (context) {},
  ),
  CustomListTile(
    title: "Notifications",
    icon: CupertinoIcons.bell,
    onTap: (context) {},
  ),
  CustomListTile(
    title: "Logout",
    icon: CupertinoIcons.arrow_right_arrow_left,
    onTap: (context) {
      SharedPreferencesProvider().clearAll();
      Navigator.of(context).pushAndRemoveUntil(
        MaterialPageRoute(builder: (context) => const LoginPage()),
        (route) => false,
      );
    },
  ),
];
