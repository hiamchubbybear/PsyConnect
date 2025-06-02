import 'package:PsyConnect/core/utils/utils.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/profile_mood.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/services/account_service/profile.dart';
import 'package:PsyConnect/ui/widgets/posts/mood.dart';
import 'package:PsyConnect/ui/widgets/posts/post.dart';
import 'package:flutter/material.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

class _HomePageScrollViewState extends State<HomePageScrollView> {
  Future<List<ProfileMoodModel>>? moodFuture;
  late Future<UserProfile> userProfile;
  ProfileService profileService = ProfileService();
  @override
  void initState() {
    super.initState();
    moodFuture = moodService.getProfileWithMood(context: context);
    userProfile = profileService.getUserProfile();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      body: SafeArea(
        child: Stack(
          children: [
            CustomScrollView(
              slivers: [
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 10.0),
                    child: FutureBuilder<List<ProfileMoodModel>>(
                      future: moodFuture,
                      builder: (context, snapshot) {
                        if (snapshot.connectionState ==
                            ConnectionState.waiting) {
                          return const Center(
                              child: CircularProgressIndicator());
                        } else if (snapshot.hasError) {
                          return const Center(
                              child: Text("Lỗi khi tải dữ liệu mood"));
                        } else if (!snapshot.hasData ||
                            snapshot.data!.isEmpty) {
                          return FutureBuilder<UserProfile>(
                            future: userProfile,
                            builder: (context, userSnapshot) {
                              if (userSnapshot.connectionState ==
                                  ConnectionState.waiting) {
                                return const Center(
                                    child: CircularProgressIndicator());
                              } else if (userSnapshot.hasError ||
                                  !userSnapshot.hasData) {
                                return const Center(
                                    child: Text(
                                        "Không thể tải thông tin người dùng"));
                              }

                              final user = userSnapshot.data!;
                              final selfMood = ProfileMoodModel(
                                profileId: user.getProfileId,
                                fullName: user.getFirstName,
                                avatarUri: user.getAvatarUri,
                                mood: mood,
                                moodId: '',
                                moodDescription: '',
                                visibility: '',
                                createdAt: 0,
                                expiresAt: 0,
                              );

                              return StoriesWidget(profilesMood: [selfMood]);
                            },
                          );
                        }

                        return StoriesWidget(profilesMood: snapshot.data!);
                      },
                    ),
                  ),
                ),
                SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (BuildContext context, int index) {
                      return PostWidget(
                        profileId: "1",
                        avatarUri:
                            "https://upload.wikimedia.org/wikipedia/commons/9/9b/Photo_of_a_kitten.jpg",
                        username: "chessy1603",
                        name: "Phong Khê",
                        postedTime: 1740478871,
                        privacy: "PUBLIC",
                        postImageUri:
                            "https://i.pinimg.com/236x/7c/89/df/7c89dfc7f3be5c1df083b01864cfb3a3.jpg",
                        liked: [
                          "huytran",
                          "congdanhhihi",
                          "thuhaaa",
                          "hphunggg"
                        ],
                        comment: ["Dễ thương vậy", "Haha"],
                        nol: 100,
                        noc: 30,
                        content: 'Xin chào thế giới',
                        postId: '',
                      );
                    },
                    childCount: 4,
                  ),
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }
}

class StoriesWidget extends StatelessWidget {
  final List<ProfileMoodModel> profilesMood;

  const StoriesWidget({super.key, required this.profilesMood});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 80,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: profilesMood.length,
        itemBuilder: (context, index) {
          final mood = profilesMood[index];
          final isOwner = index == 0;

          return Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Column(
              children: [
                Stack(
                  children: [
                    Container(
                      width: 60,
                      height: 60,
                      decoration: BoxDecoration(
                        shape: BoxShape.circle,
                        border: Border.all(color: secondaryColor, width: 3),
                      ),
                      child: ClipOval(
                        child: Image.network(
                          mood.avatarUri ?? '',
                          errorBuilder: (context, error, stackTrace) =>
                              const Icon(Icons.error, size: 60),
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return const Center(
                                child: CircularProgressIndicator());
                          },
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                    if (isOwner)
                      Positioned(
                        bottom: 0,
                        right: 0,
                        child: MoodWidget(),
                      ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  Utils().namesplite(name: mood.fullName ?? ''),
                  style: quickSand12Font,
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
