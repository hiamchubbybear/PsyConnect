import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/core/utils/utils.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/mood.dart';
import 'package:PsyConnect/models/profile_mood.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';
import 'package:PsyConnect/services/profile_service/mood.dart';
import 'package:PsyConnect/ui/widgets/posts/mood.dart';
import 'package:PsyConnect/ui/widgets/posts/post.dart';
import 'package:flutter/material.dart';
import 'package:pull_to_refresh_flutter3/pull_to_refresh_flutter3.dart';

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
    if (mounted) setState(() {});
    _refreshController.refreshCompleted();
  }

  void _onLoading() async {
    await Future.delayed(Duration(milliseconds: 1000));
    if (mounted)
      setState(() {
        _onMoodCreated();
      });
    _refreshController.loadComplete();
  }

  @override
  void initState() {
    super.initState();
    _fetchMoods();
  }

  void _fetchMoods() async {
    moodFuture = (() async {
      final moods = await moodService.getProfileWithMood(context: context);
      final user = await ProfileService().getUserProfile();

      final hasSelfMood = moods.any((m) => m.profileId == user.getProfileId);
      for (var toElement in moods) {
        print("Tên người dùng: ${toElement.fullName}");
      }

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
    return Scaffold(
      backgroundColor: Theme.of(context).scaffoldBackgroundColor,
      body: SafeArea(
        child: SmartRefresher(
          controller: _refreshController,
          enablePullDown: true,
          enablePullUp: true,
          onRefresh: _onRefresh,
          onLoading: _onLoading,
          child: CustomScrollView(
            slivers: [
              SliverToBoxAdapter(
                child: Padding(
                  padding: const EdgeInsets.symmetric(vertical: 10.0),
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
                            );
                          },
                        );
                      }
                      return StoriesWidget(
                        profilesMood: snapshot.data!,
                        onMoodCreated: _onMoodCreated,
                      );
                    },
                  ),
                ),
              ),
              SliverList(
                delegate: SliverChildBuilderDelegate(
                  childCount: 4,
                  (context, index) => const PostWidget(
                    profileId: "1",
                    avatarUri:
                        "https://upload.wikimedia.org/wikipedia/commons/9/9b/Photo_of_a_kitten.jpg",
                    username: "chessy1603",
                    name: "Phong Khê",
                    postedTime: 1740478871,
                    privacy: "PUBLIC",
                    postImageUri:
                        "https://i.pinimg.com/236x/7c/89/df/7c89dfc7f3be5c1df083b01864cfb3a3.jpg",
                    liked: ["huytran", "congdanhhihi", "thuhaaa", "hphunggg"],
                    comment: ["Dễ thương vậy", "Haha"],
                    nol: 37,
                    noc: 30,
                    content: 'Xin chào thế giới',
                    postId: '',
                  ),
                ),
              ),
            ],
          ),
        ),
      ),
    );
  }
}

class StoriesWidget extends StatelessWidget {
  final List<ProfileMoodModel> profilesMood;
  final VoidCallback? onMoodCreated;

  const StoriesWidget({
    super.key,
    required this.profilesMood,
    this.onMoodCreated,
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
                      child: ClipOval(
                        child: Image.network(
                          mood.avatarUri,
                          fit: BoxFit.cover,
                          errorBuilder: (context, error, stackTrace) =>
                              const Icon(Icons.error),
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return const Center(
                                child: CircularProgressIndicator());
                          },
                        ),
                      ),
                    ),
                    if (mood.moodDescription.trim().isNotEmpty)
                      Positioned(
                        bottom: 45,
                        right: -10,
                        child:
                            MoodNoteBubbleWithSmoke(text: mood.moodDescription),
                      ),
                    if (isOwner && mood.moodDescription.trim().isEmpty)
                      Positioned(
                        bottom: 45,
                        right: -10,
                        child: CreateMoodWidget(onMoodCreated: onMoodCreated),
                      ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  mood.fullName.isNotEmpty
                      ? Utils().namesplite(name: mood.fullName)
                      : "Unknown name",
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

class CreateMoodWidget extends StatefulWidget {
  final VoidCallback? onMoodCreated;
  const CreateMoodWidget({super.key, required this.onMoodCreated});
  @override
  State<CreateMoodWidget> createState() => _CreateMoodWidgetState();
}

String mood = "testing";
String moodDescription = "";
String visibility = "";
String selectedItem = items[0];
List<String> items = [
  "Private",
  "Public",
  "Friends only",
];
MoodService moodService = MoodService();

class _CreateMoodWidgetState extends State<CreateMoodWidget> {
  @override
  Widget build(BuildContext context) {
    return GestureDetector(
      onTap: () => showPostDialog(context, widget.onMoodCreated),
      child: Container(
        width: 20,
        height: 20,
        decoration: BoxDecoration(
          shape: BoxShape.circle,
          color: Colors.green[300],
          border: Border.all(color: Colors.white, width: 2),
        ),
        child: const Icon(
          Icons.add,
          size: 14,
          color: Colors.white,
        ),
      ),
    );
  }
}

void showPostDialog(BuildContext context, VoidCallback? onMoodCreated) {
  showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: "Post Dialog",
    transitionDuration: const Duration(milliseconds: 100),
    pageBuilder: (context, animation, secondaryAnimation) {
      return Center(
        child: Material(
          borderRadius: BorderRadius.circular(20),
          color: Colors.white,
          child: StatefulBuilder(
            builder: (context, setState) {
              return Container(
                width: MediaQuery.of(context).size.width * 1.3,
                height: 270,
                padding:
                    const EdgeInsets.symmetric(vertical: 20, horizontal: 10),
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(
                      mainAxisAlignment: MainAxisAlignment.spaceBetween,
                      children: [
                        Text("Share your mind", style: quickSand12Font),
                        DropdownButton<String>(
                          value: selectedItem,
                          items: items
                              .map((item) => DropdownMenuItem<String>(
                                  value: item,
                                  child: Text(item, style: quickSand12Font)))
                              .toList(),
                          onChanged: (String? value) {
                            setState(() {
                              selectedItem = value!;
                              visibility = value!;
                            });
                          },
                        ),
                      ],
                    ),
                    const SizedBox(height: 10),
                    TextField(
                      maxLines: 3,
                      maxLength: 15,
                      decoration: InputDecoration(
                        hintText: "Drop a thought",
                        hintStyle: kSubHintStyle,
                        border: OutlineInputBorder(
                            borderRadius: BorderRadius.circular(10)),
                      ),
                      onChanged: (value) => {
                        moodDescription = value,
                      },
                    ),
                    Row(
                      mainAxisAlignment: MainAxisAlignment.end,
                      children: [
                        TextButton(
                          child: const Text("Cancel"),
                          onPressed: () => Navigator.of(context).pop(),
                        ),
                        const SizedBox(width: 8),
                        ElevatedButton(
                          child: const Text("Post"),
                          onPressed: () {
                            if (mood.isEmpty) {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please type what you are thinking",
                                  title: "Warning",
                                  type: ToastType.warning);
                            } else if (visibility == "") {
                              ToastService.showToast(
                                  context: context,
                                  message: "Please choose post privacy",
                                  title: "Warning",
                                  type: ToastType.warning);
                            } else {
                              MoodModel moodModel = MoodModel(
                                  mood: mood,
                                  moodDescription: moodDescription,
                                  visibility: visibility);
                              moodService.createMoodPost(
                                  mood: moodModel, context: context);
                              if (onMoodCreated != null) {
                                onMoodCreated();
                              } else {
                                ToastService.showToast(
                                    context: context,
                                    message: "Error while render",
                                    title: "Uncatogerize exception",
                                    type: ToastType.warning);
                              }
                              Navigator.of(context).pop();
                            }
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
        child: ScaleTransition(scale: anim1, child: child),
      );
    },
  );
}
