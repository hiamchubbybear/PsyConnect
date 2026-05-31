import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/utils/utils.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/models/profile_mood.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/post_provider.dart';
import 'package:PsyConnect/services/profile_service/mood.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:PsyConnect/ui/widgets/posts/create_mood.dart';
import 'package:PsyConnect/ui/widgets/posts/mood.dart';
import 'package:PsyConnect/ui/widgets/posts/post.dart';
import 'package:PsyConnect/validate/validate.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:pull_to_refresh_flutter3/pull_to_refresh_flutter3.dart';
import 'package:google_fonts/google_fonts.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

class _HomePageScrollViewState extends State<HomePageScrollView> {
  void handleMoodCreated() {
    _fetchMoods();
  }

  late Future<List<ProfileMoodModel>> moodFuture;
  late Future<UserProfile> userProfile;

  final RefreshController _refreshController =
      RefreshController(initialRefresh: false);

  void _onRefresh() async {
    _fetchMoods();
    await Provider.of<PostProvider>(context, listen: false).fetchPosts(refresh: true);
    if (mounted) setState(() {});
    _refreshController.refreshCompleted();
  }

  void _onLoading() async {
    await Provider.of<PostProvider>(context, listen: false).fetchPosts(refresh: false);
    _refreshController.loadComplete();
  }

  @override
  void initState() {
    super.initState();
    _fetchMoods();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<PostProvider>(context, listen: false).fetchPosts(refresh: true);
    });
  }

  void _fetchMoods() async {
    moodFuture = (() async {
      final moods = await MoodService().getProfileWithMood(context: context);
      final user = await ProfileService().getUserProfile();
      final hasSelfMood = moods.any((m) => m.profileId == user.getProfileId);
      if (!hasSelfMood) {
        moods.insert(
          0,
          ProfileMoodModel(
            profileId: user.getProfileId,
            fullName: user.getFirstName,
            avatarUri: user.getAvatarUri,
            mood: "",
            moodId: '',
            moodDescription: '',
            visibility: '',
            createdAt: 0,
            expiresAt: 0,
          ),
        );
      }
      return moods;
    })();
    userProfile = ProfileService().getUserProfile();
  }

  void _onMoodCreated() {
    setState(() {
      _fetchMoods();
    });
  }

  @override
  Widget build(BuildContext context) {
    final themeProvider = Provider.of<ThemeProvider>(context);
    final bool isDark = themeProvider.isDarkMode;
    final postProvider = Provider.of<PostProvider>(context);

    final double topPadding = MediaQuery.of(context).padding.top + 25;

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      body: SmartRefresher(
        controller: _refreshController,
        enablePullDown: true,
        enablePullUp: true,
        onRefresh: _onRefresh,
        onLoading: _onLoading,
        child: CustomScrollView(
          slivers: [
            SliverPadding(
              padding: EdgeInsets.only(top: topPadding, bottom: 10),
              sliver: SliverToBoxAdapter(
                child: FutureBuilder<List<ProfileMoodModel>>(
                  future: moodFuture,
                  builder: (context, snapshot) {
                    if (snapshot.connectionState == ConnectionState.waiting) {
                      return const Center(child: CircularProgressIndicator());
                    } else if (snapshot.hasError) {
                      return const Center(child: Text("Failed to load data"));
                    } else if (!snapshot.hasData || snapshot.data!.isEmpty) {
                      return FutureBuilder<UserProfile>(
                        future: userProfile,
                        builder: (context, userSnap) {
                          if (!userSnap.hasData) {
                            return const Center(
                                child: CircularProgressIndicator());
                          }
                          final user = userSnap.data!;
                          final selfMood = ProfileMoodModel(
                            profileId: user.getProfileId,
                            fullName: user.getFirstName,
                            avatarUri: user.getAvatarUri,
                            mood: "",
                            moodId: '',
                            moodDescription: '',
                            visibility: '',
                            createdAt: 0,
                            expiresAt: 0,
                          );
                          return StoriesWidget(
                            profilesMood: [selfMood],
                            onMoodCreated: _onMoodCreated,
                            isDark: isDark,
                          );
                        },
                      );
                    }
                    return StoriesWidget(
                      profilesMood: snapshot.data!,
                      onMoodCreated: _onMoodCreated,
                      isDark: isDark,
                    );
                  },
                ),
              ),
            ),
            
            // Live dynamic posts sliver list
            SliverList(
              delegate: SliverChildBuilderDelegate(
                (context, index) {
                  if (postProvider.posts.isEmpty) {
                    if (postProvider.isLoading) {
                      return Center(
                        child: Padding(
                          padding: const EdgeInsets.symmetric(vertical: 40),
                          child: CircularProgressIndicator(
                            color: isDark ? secondaryColor : primaryColor,
                          ),
                        ),
                      );
                    }
                    return Center(
                      child: Padding(
                        padding: const EdgeInsets.symmetric(vertical: 60, horizontal: 20),
                        child: Column(
                          children: [
                            Icon(Icons.article_outlined, size: 48, color: Colors.grey[400]),
                            const SizedBox(height: 12),
                            Text(
                              "No articles written yet.",
                              style: GoogleFonts.quicksand(
                                fontWeight: FontWeight.bold,
                                color: Colors.grey[500],
                                fontSize: 16,
                              ),
                            ),
                          ],
                        ),
                      ),
                    );
                  }
                  final post = postProvider.posts[index];
                  return PostWidget(
                    isDark: isDark,
                    post: post,
                  );
                },
                childCount: postProvider.posts.isEmpty ? 1 : postProvider.posts.length,
              ),
            ),
          ],
        ),
      ),
    );
  }
}

class StoriesWidget extends StatelessWidget {
  final List<ProfileMoodModel> profilesMood;
  final VoidCallback? onMoodCreated;
  final bool isDark;
  const StoriesWidget({
    super.key,
    required this.profilesMood,
    this.onMoodCreated,
    required this.isDark,
  });

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
                  clipBehavior: Clip.none,
                  children: [
                    Container(
                        width: 60,
                        height: 60,
                        decoration: BoxDecoration(
                          shape: BoxShape.circle,
                          border: Border.all(color: secondaryColor, width: 3),
                        ),
                        child: CircleAvatar(
                          radius: 50,
                          backgroundColor: Colors.grey,
                          child: FutureBuilder<bool>(
                              future: checkImageExists(mood.avatarUri),
                              builder: (context, snapshot) {
                                String imageUrl = snapshot.hasData &&
                                        snapshot.data == true
                                    ? mood.avatarUri
                                    : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';
                                return CircleAvatar(
                                  radius: 50,
                                  backgroundImage: NetworkImage(imageUrl),
                                );
                              }),
                        )),
                    if (mood.mood.trim().isNotEmpty)
                      Positioned(
                        bottom: 45,
                        right: -35,
                        child: GestureDetector(
                          behavior: HitTestBehavior.translucent,
                          onTap: () => _showMoodOptions(
                              context, mood, onMoodCreated, isDark),
                          child: Container(
                            padding: const EdgeInsets.all(8),
                            constraints: const BoxConstraints(
                              minWidth: 80,
                              minHeight: 40,
                            ),
                            child: MoodNoteBubbleWithSmoke(
                              text: mood.mood,
                              isDark: isDark,
                            ),
                          ),
                        ),
                      ),
                    if (isOwner && mood.mood.trim().isEmpty)
                      Positioned(
                        bottom: 45,
                        right: -10,
                        child: MoodWidget(),
                      ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  mood.fullName.isNotEmpty
                      ? Utils().namesplite(name: mood.fullName)
                      : "Unknown",
                  style: quickSand12FontMoodCreate,
                ),
              ],
            ),
          );
        },
      ),
    );
  }

  void _showMoodOptions(BuildContext context, ProfileMoodModel mood,
      VoidCallback? onMoodCreated, bool isDarkMode) {
    final isDark = isDarkMode;
    showModalBottomSheet(
      context: context,
      backgroundColor: isDark ? Colors.grey.shade900 : Colors.white,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (BuildContext context) {
        return Container(
          padding: const EdgeInsets.all(20),
          child: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                width: 40,
                height: 4,
                decoration: BoxDecoration(
                  color: isDark ? Colors.grey.shade600 : Colors.grey.shade300,
                  borderRadius: BorderRadius.circular(2),
                ),
              ),
              const SizedBox(height: 20),
              Text(
                "Mood Options",
                style: kSubHeadingStyle.copyWith(
                  fontSize: 18,
                  fontWeight: FontWeight.w600,
                ),
              ),
              const SizedBox(height: 20),
              ListTile(
                leading: Icon(
                  Icons.edit,
                  color: isDark ? Colors.blue.shade300 : Colors.blue,
                ),
                title: Text(
                  "Edit Mood",
                  style: kSubHeadingStyle.copyWith(fontWeight: FontWeight.w500),
                ),
                onTap: () {
                  Navigator.pop(context);
                  _showEditMoodDialog(context, mood, onMoodCreated, isDark);
                },
              ),
              ListTile(
                leading: Icon(
                  Icons.delete,
                  color: isDark ? Colors.red.shade300 : Colors.red,
                ),
                title: Text(
                  "Delete Mood",
                  style: kSubHeadingStyle.copyWith(fontWeight: FontWeight.w500),
                ),
                onTap: () {
                  Navigator.pop(context);
                  _showDeleteConfirmation(context, mood, onMoodCreated, isDark);
                },
              ),
              const SizedBox(height: 10),
            ],
          ),
        );
      },
    );
  }

  void _showEditMoodDialog(BuildContext context, ProfileMoodModel mood,
      VoidCallback? onMoodCreated, bool isDark) {
    String editedMoodDescription = mood.moodDescription;
    String editedVisibility = mood.visibility;

    showGeneralDialog(
      context: context,
      barrierDismissible: true,
      barrierLabel: "Edit Mood Dialog",
      transitionDuration: const Duration(milliseconds: 200),
      pageBuilder: (context, animation, secondaryAnimation) {
        return Center(
          child: Material(
            borderRadius: BorderRadius.circular(20),
            color: isDark ? Colors.grey.shade900 : Colors.white,
            elevation: 8,
            child: StatefulBuilder(
              builder: (context, setState) {
                return Container(
                  width: MediaQuery.of(context).size.width * 0.9,
                  height: 320,
                  padding: const EdgeInsets.all(24),
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Text(
                            "Edit your mood",
                            style: kSubHeadingStyle.copyWith(
                              fontWeight: FontWeight.w600,
                              fontSize: 18,
                            ),
                          ),
                          Container(
                            padding: const EdgeInsets.symmetric(
                                horizontal: 12, vertical: 4),
                            decoration: BoxDecoration(
                              color: isDark
                                  ? Colors.grey.shade700
                                  : Colors.grey.shade100,
                              borderRadius: BorderRadius.circular(8),
                            ),
                            child: DropdownButton<String>(
                              value: editedVisibility.isEmpty
                                  ? "Private"
                                  : editedVisibility,
                              underline: const SizedBox(),
                              icon: Icon(
                                Icons.keyboard_arrow_down,
                                color: isDark
                                    ? Colors.white70
                                    : Colors.grey.shade600,
                              ),
                              items: ["Private", "Public", "Friends only"]
                                  .map((item) => DropdownMenuItem<String>(
                                      value: item,
                                      child: Text(
                                        item,
                                        style: kSubHeadingStyle.copyWith(
                                          fontSize: 13,
                                          fontWeight: FontWeight.w500,
                                        ),
                                      )))
                                  .toList(),
                              onChanged: (String? value) {
                                setState(() {
                                  editedVisibility = value!;
                                });
                              },
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 20),
                      Expanded(
                        child: TextField(
                          controller: TextEditingController(
                              text: editedMoodDescription),
                          maxLines: null,
                          expands: true,
                          maxLength: 200,
                          style: kSubHeadingStyle.copyWith(
                            fontSize: 15,
                            height: 1.4,
                          ),
                          decoration: InputDecoration(
                            hintText: "What are you thinking about?",
                            hintStyle: kSubHintStyle,
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(
                                color: isDark
                                    ? Colors.grey.shade600
                                    : Colors.grey.shade300,
                              ),
                            ),
                            enabledBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(
                                color: isDark
                                    ? Colors.grey.shade600
                                    : Colors.grey.shade300,
                              ),
                            ),
                            focusedBorder: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(12),
                              borderSide: BorderSide(
                                color: acceptColor,
                                width: 2,
                              ),
                            ),
                            fillColor: isDark
                                ? Colors.grey.shade800
                                : Colors.grey.shade50,
                            filled: true,
                            contentPadding: const EdgeInsets.all(16),
                          ),
                          onChanged: (value) {
                            editedMoodDescription = value;
                          },
                        ),
                      ),
                      const SizedBox(height: 20),
                      Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          TextButton(
                            style: TextButton.styleFrom(
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 20, vertical: 12),
                              shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(10)),
                            ),
                            child: Text(
                              "Cancel",
                              style: kSubHeadingStyle.copyWith(
                                color: isDark
                                    ? Colors.white70
                                    : Colors.grey.shade600,
                                fontWeight: FontWeight.w500,
                              ),
                            ),
                            onPressed: () => Navigator.of(context).pop(),
                          ),
                          const SizedBox(width: 12),
                          ElevatedButton(
                            style: ElevatedButton.styleFrom(
                              backgroundColor: acceptColor,
                              padding: const EdgeInsets.symmetric(
                                  horizontal: 24, vertical: 12),
                              shape: RoundedRectangleBorder(
                                  borderRadius: BorderRadius.circular(10)),
                              elevation: 2,
                            ),
                            child: Text(
                              "Update",
                              style: kSubHeadingStyle.copyWith(
                                color: Colors.white,
                                fontWeight: FontWeight.w600,
                              ),
                            ),
                            onPressed: () async {
                              if (editedMoodDescription.trim().isEmpty) {
                                ToastService.showToast(
                                    context: context,
                                    message: "Please enter your mood",
                                    title: "Warning",
                                    type: ToastType.warning);
                                return;
                              }
                              MoodModel updatedMood = MoodModel(
                                mood: editedMoodDescription,
                                moodDescription: editedMoodDescription,
                                visibility: editedVisibility,
                              );
                              await MoodService().updateMoodPost(
                                mood: updatedMood,
                                context: context,
                              );

                              if (onMoodCreated != null) {
                                onMoodCreated!();
                              }
                              Navigator.pop(context);
                            },
                          ),
                        ],
                      )
                    ],
                  ),
                );
              },
            ),
          ),
        );
      },
      transitionBuilder: (context, anim1, anim2, child) {
        return FadeTransition(
          opacity: anim1,
          child: ScaleTransition(
            scale: Tween<double>(begin: 0.8, end: 1.0).animate(
              CurvedAnimation(parent: anim1, curve: Curves.easeOutBack),
            ),
            child: child,
          ),
        );
      },
    );
  }
}

void _showDeleteConfirmation(BuildContext context, ProfileMoodModel mood,
    VoidCallback? onMoodCreated, bool isDark) {
  showDialog(
    context: context,
    builder: (BuildContext context) {
      return AlertDialog(
        backgroundColor: isDark ? Colors.grey.shade900 : Colors.white,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(15),
        ),
        title: Text(
          "Delete Mood",
          style: kSubHeadingStyle.copyWith(
            fontWeight: FontWeight.w600,
            fontSize: 18,
          ),
        ),
        content: Text(
          "Are you sure you want to delete this mood? This action cannot be undone.",
          style: kSubHeadingStyle.copyWith(
            fontSize: 14,
            color: isDark ? Colors.white70 : Colors.grey.shade600,
          ),
        ),
        actions: [
          TextButton(
            child: Text(
              "Cancel",
              style: kSubHeadingStyle.copyWith(
                color: isDark ? Colors.white70 : Colors.grey.shade600,
                fontWeight: FontWeight.w500,
              ),
            ),
            onPressed: () => Navigator.of(context).pop(),
          ),
          ElevatedButton(
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.red,
              shape: RoundedRectangleBorder(
                borderRadius: BorderRadius.circular(8),
              ),
            ),
            child: Text(
              "Delete",
              style: kSubHeadingStyle.copyWith(
                color: Colors.white,
                fontWeight: FontWeight.w600,
              ),
            ),
            onPressed: () async {
              try {
                await MoodService().deleteMoodPost(
                  moodId: mood.moodId,
                  context: context,
                );

                if (onMoodCreated != null) {
                  onMoodCreated!();
                }
                Navigator.pop(context);
                ToastService.showToast(
                    context: context,
                    message: "Mood deleted successfully",
                    title: "Success",
                    type: ToastType.success);
              } catch (e) {
                print(e.toString());
                ToastService.showToast(
                    context: context,
                    message: "Failed to delete mood",
                    title: "Error",
                    type: ToastType.error);
              }
            },
          ),
        ],
      );
    },
  );
}
