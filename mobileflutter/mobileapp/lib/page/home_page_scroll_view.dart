import 'package:PsyConnect/themes/colors/color.dart';
import 'package:PsyConnect/widgets/posts/post.dart';
import 'package:flutter/material.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

class _HomePageScrollViewState extends State<HomePageScrollView> {
  @override
  Widget build(BuildContext context) {
    return Stack(
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
                    liked: ["huytran", "congdanhhihi", "thuhaaa", "hphunggg"],
                    comment: ["Dễ thương vậy", "Haha"],
                    nol: 100,
                    noc: 30,
                    content: 'Xin chào thế giới',
                  );
                },
                childCount: 50,
              ),
            ),
          ],
        ),
      ],
    );
  }
}

class StoriesWidget extends StatelessWidget {
  final List<String> storyImages = [
    'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg',
    'https://th.bing.com/th/id/OIP.TvwRR9fnAiYIBm4zkGdxAwHaIM?w=926&h=1024&rs=1&pid=ImgDetMain',
    'https://th.bing.com/th/id/OIP.VhMZO57QYlGxH9AQM9-FxAHaHa?pid=ImgDet&w=196&h=196&c=7&dpr=1.8',
    'https://th.bing.com/th/id/OIP.EpUD2PGp7HmIA3a0UnKl6AAAAA?pid=ImgDet&w=196&h=196&c=7&dpr=1.8',
    'https://pbs.twimg.com/profile_images/1306417770356113408/1GLl3v0u_400x400.jpg',
  ];
  final List<String> username = [
    'Trần Meo',
    'Nguyễn Thị Thu Meo',
    'Phan Công Meo',
    'Phùng Đình Quang Meo',
    'Nguyễn Thị Như Meo',
  ];
  StoriesWidget({super.key});

  @override
  Widget build(BuildContext context) {
    return SizedBox(
      height: 110,
      child: ListView.builder(
        scrollDirection: Axis.horizontal,
        itemCount: storyImages.length,
        itemBuilder: (context, index) {
          return Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Column(
              children: [
                Container(
                  width: 70,
                  height: 70,
                  decoration: BoxDecoration(
                    shape: BoxShape.circle,
                    border: Border.all(color: secondaryColor, width: 3),
                  ),
                  child: Container(
                    child: ClipOval(
                      child: Image.network(
                        storyImages[index],
                        fit: BoxFit.cover,
                      ),
                    ),
                  ),
                ),
                const SizedBox(height: 5),
                Text(
                  "${username[index]}",
                  style: const TextStyle(fontSize: 12),
                ),
              ],
            ),
          );
        },
      ),
    );
  }
}
