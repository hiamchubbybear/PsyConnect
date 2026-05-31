import 'dart:convert';

import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/route/route_animation.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:PsyConnect/ui/screens/consultation_profile_page.dart';
import 'package:PsyConnect/ui/screens/login_page.dart';
import 'package:PsyConnect/ui/screens/setting_page.dart';
import 'package:PsyConnect/ui/widgets/common/custom_button.dart';
import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/validate/validate.dart';

class ProfilePage extends StatefulWidget {
  const ProfilePage({super.key});

  @override
  State<ProfilePage> createState() => _ProfilePageState();
}

class _ProfilePageState extends State<ProfilePage> {
  final ProfileService profileService = ProfileService();
  UserProfile userProfile = UserProfile();

  void handleSetProfileDetails(BuildContext context) {
    handleOnProfile(context);
  }

  void handleUploadResume(BuildContext context) {
    handleOnProfile(context);
  }

  void handleAddSkills(BuildContext context) {
    handleOnProfile(context);
  }

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
    } catch (e) {
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
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: const Text("Profile"),
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
        padding: const EdgeInsets.fromLTRB(16, 16, 16, 50),
        children: [
          Column(
            children: [
              GestureDetector(
                onTap: () => handleOnProfile(context),
                child: Container(
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    border: Border.all(
                      color: isDark ? Colors.blue[300]! : Colors.blue,
                      width: 2,
                    ),
                  ),
                  child: CircleAvatar(
                    radius: 50,
                    backgroundColor: isDark ? Colors.grey[800] : Colors.grey[200],
                    child: FutureBuilder<bool>(
                      future: checkImageExists(userProfile.getAvatarUri),
                      builder: (context, snapshot) {
                        final String imageUrl = snapshot.hasData && snapshot.data == true
                            ? userProfile.getAvatarUri
                            : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

                        return CircleAvatar(
                          radius: 48,
                          backgroundImage: NetworkImage(imageUrl),
                        );
                      },
                    ),
                  ),
                ),
              ),
              const SizedBox(height: 12),
              Text(
                userProfile.getFirstName.isNotEmpty ? userProfile.getFirstName : "Guest User",
                style: theme.textTheme.titleLarge?.copyWith(fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 4),
              Text(
                userProfile.getDescription.isNotEmpty ? userProfile.getDescription : "Mental Health Enthusiast",
                style: theme.textTheme.bodyMedium?.copyWith(color: isDark ? Colors.grey[400] : Colors.grey[600]),
              )
            ],
          ),
          const SizedBox(height: 32),
          FutureBuilder<int>(
            future: checkComplete(userProfile),
            builder: (context, snapshot) {
              if (snapshot.connectionState == ConnectionState.waiting) {
                return const SizedBox();
              } else if (snapshot.hasError) {
                return Text('Error: ${snapshot.error}');
              } else if (snapshot.hasData) {
                final int value = snapshot.data!;
                if (value > 0) {
                  return Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        children: [
                          Text(
                            "Complete your profile",
                            style: theme.textTheme.titleMedium?.copyWith(fontWeight: FontWeight.bold),
                          ),
                          const SizedBox(width: 6),
                          Text(
                            "($value/8)",
                            style: TextStyle(
                              color: isDark ? Colors.blue[300] : Colors.blue,
                              fontWeight: FontWeight.bold,
                            ),
                          )
                        ],
                      ),
                      const SizedBox(height: 12),
                      Row(
                        children: List.generate(5, (index) {
                          return Expanded(
                            child: Container(
                              height: 7,
                              margin: EdgeInsets.only(right: index == 4 ? 0 : 6),
                              decoration: BoxDecoration(
                                borderRadius: BorderRadius.circular(10),
                                color: index < (8 - value)
                                    ? (isDark ? Colors.blue[300] : Colors.blue)
                                    : (isDark ? Colors.grey[800] : Colors.grey[200]),
                              ),
                            ),
                          );
                        }),
                      ),
                    ],
                  );
                }
              }
              return const SizedBox();
            },
          ),
          const SizedBox(height: 20),
          SizedBox(
            height: 160,
            child: ListView.separated(
              physics: const BouncingScrollPhysics(),
              scrollDirection: Axis.horizontal,
              itemBuilder: (context, index) {
                final card = profileCompletionCards[index];
                return SizedBox(
                  width: 160,
                  child: Card(
                    elevation: 1,
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(12)),
                    child: Padding(
                      padding: const EdgeInsets.all(12),
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          Icon(
                            card.icon,
                            size: 28,
                            color: isDark ? Colors.blue[300] : Colors.blue,
                          ),
                          const SizedBox(height: 8),
                          Text(
                            card.title,
                            textAlign: TextAlign.center,
                            maxLines: 2,
                            overflow: TextOverflow.ellipsis,
                            style: theme.textTheme.bodyMedium?.copyWith(fontWeight: FontWeight.bold),
                          ),
                          const Spacer(),
                          CustomButton(
                            onPressed: () {
                              if (card.onTap != null) {
                                card.onTap!(context);
                              }
                            },
                            text: card.buttonText,
                            height: 36,
                            borderRadius: 8,
                          )
                        ],
                      ),
                    ),
                  ),
                );
              },
              separatorBuilder: (context, index) => const SizedBox(width: 8),
              itemCount: profileCompletionCards.length,
            ),
          ),
          const SizedBox(height: 32),
          ...List.generate(
            customListTiles.length,
            (index) {
              final tile = customListTiles[index];
              return Padding(
                padding: const EdgeInsets.only(bottom: 8),
                child: Card(
                  elevation: 1,
                  shape: RoundedRectangleBorder(
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: InkWell(
                    borderRadius: BorderRadius.circular(12),
                    onTap: () => tile.onTap?.call(context),
                    child: ListTile(
                      leading: Icon(tile.icon, color: isDark ? Colors.grey[400] : Colors.grey[600]),
                      title: Text(
                        tile.title,
                        style: theme.textTheme.bodyLarge?.copyWith(fontWeight: FontWeight.w600),
                      ),
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
  Navigator.push(
      context, createSlideFromBottomRoute(const ConsultationProfilePage()));
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
