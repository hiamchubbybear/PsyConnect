import 'package:PsyConnect/ui/widgets/posts/post.dart';
import 'package:PsyConnect/variable/variable.dart';
import 'package:flutter/material.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

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
                      );
                    },
                    childCount: 3,
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
    'Nguyễn Thị Thu Meo',
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
                          storyImages[index],
                          fit: BoxFit.cover,
                        ),
                      ),
                    ),
                    if (isOwner)
                      Positioned(
                        bottom: 0,
                        right: 0,
                        child: GestureDetector(
                          onTap: () =>showPostDialog(context) ,
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
                        ),
                      ),
                  ],
                ),
                const SizedBox(height: 2),
                Text(
                  namesplite(name: username[index]),
                  style: quickSand12Font,
                ),
              ],
            ),
          );
        },
      ),
    );
  }
    String namesplite({required String name}) {
      if (name.length <= 14)
        return name;
      else {
        List<String> splited = name.split(" ");
        String res = "";
        for (int i = 0; i < splited.length; i++) {
          if (res.length + splited[i].length <= 14) {
            res += splited[i] + " ";
          } else
            return res;
        }
      }
      return name;
    }
  }
void showPostDialog(BuildContext context) {
  showGeneralDialog(
    context: context,
    barrierDismissible: true,
    barrierLabel: "Post Dialog",
    transitionDuration: Duration(milliseconds: 200),
    pageBuilder: (context, animation, secondaryAnimation) {
      return Center(
        child: Material(
          borderRadius: BorderRadius.circular(20),
          color: Colors.white,
          child: Container(
            width: MediaQuery.of(context).size.width * 1.3,
            height: 270,
            padding: EdgeInsets.symmetric(vertical: 20 , horizontal:  10),
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text("Share your mind", style: quickSand12Font),
                const SizedBox(height: 10),
                const  TextField(
                  maxLines: 3,
                  maxLength: 200,
                  decoration: InputDecoration(
                    hintText: "What are you thinking about?",
                    border: OutlineInputBorder(),
                  ),
                ),
                Spacer(),
                Row(
                  mainAxisAlignment: MainAxisAlignment.end,
                  children: [
                    TextButton(
                      child: Text("Cancel"),
                      onPressed: () => Navigator.of(context).pop(),
                    ),
                    SizedBox(width: 8),
                    ElevatedButton(
                      child: Text("Post"),
                      onPressed: () {
                        // Xử lý đăng bài
                        Navigator.of(context).pop();
                      },
                    ),
                  ],
                )
              ],
            ),
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
