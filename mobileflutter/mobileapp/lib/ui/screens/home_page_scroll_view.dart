import 'package:PsyConnect/core/utils/utils.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/ui/widgets/posts/mood.dart';
import 'package:PsyConnect/ui/widgets/posts/post.dart';
import 'package:flutter/material.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

Utils utils = Utils();

class _HomePageScrollViewState extends State<HomePageScrollView> {
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
                    padding: EdgeInsets.symmetric(vertical: 10.0),
                    child: StoriesWidget(),
                  ),
                ),
                SliverList(
                  delegate: SliverChildBuilderDelegate(
                    (BuildContext context, int index) {
                      return PostWidget(
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
                        userId: '',
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
  final List<UserProfile> storyImages = [
    UserProfile(
      userId: "1",
      username: "chessy1603",
      avatarUri:
          "https://upload.wikimedia.org/wikipedia/commons/9/9b/Photo_of_a_kitten.jpg",
    ),
    UserProfile(
      userId: "2",
      username: "huytran",
      avatarUri:
          "https://upload.wikimedia.org/wikipedia/commons/4/47/PNG_transparency_demonstration_1.png",
    ),
    UserProfile(
      userId: "3",
      username: "congdanhhihi",
      avatarUri:
          "https://upload.wikimedia.org/wikipedia/commons/thumb/a/a0/Flag_of_Vietnam.svg/1920px-Flag_of_Vietnam.svg.png",
    ),
    UserProfile(
      userId: "4",
      username: "thuhaaa",
      avatarUri:
          "https://upload.wikimedia.org/wikipedia/commons/thumb/f/fc/Flag_of_the_People%27s_Republic_of_China.svg/1920px-Flag_of_the_People%27s_Republic_of_China.svg.png",
    ),
  ];
  StoriesWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 80,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: storyImages.length,
        itemBuilder: (context, index) {
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
                          storyImages[index].avatarUri ?? '',
                          errorBuilder: (context, error, stackTrace) {
                            return const Icon(Icons.error, size: 60);
                          },
                          loadingBuilder: (context, child, loadingProgress) {
                            if (loadingProgress == null) return child;
                            return const Center(
                              child: CircularProgressIndicator(),
                            );
                          },
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                    if (isOwner)
                      Positioned(bottom: 0, right: 0, child: MoodWidget()),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  utils.namesplite(name: storyImages[index].username ?? ''),
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
