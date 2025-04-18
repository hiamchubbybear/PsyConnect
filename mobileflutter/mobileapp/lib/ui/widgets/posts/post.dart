import 'package:PsyConnect/ui/themes/colors/color.dart';
import 'package:PsyConnect/variable/variable.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:timeago/timeago.dart' as timeago;

class PostWidget extends StatefulWidget {
  final String avatarUri;
  final String username;
  final String name;
  final String content;
  final String? postImageUri;
  final int postedTime;
  final String privacy;
  final List<String> liked;
  final List<String> comment;
  final int nol;
  final int noc;

  const PostWidget({
    super.key,
    required this.avatarUri,
    required this.username,
    required this.name,
    required this.content,
    required this.postedTime,
    this.postImageUri,
    required this.privacy,
    required this.liked,
    required this.comment,
    required this.nol,
    required this.noc,
  });

  @override
  State<PostWidget> createState() => _PostWidgetState();
}

class _PostWidgetState extends State<PostWidget> {
  bool isLiked = false;

  @override
  Widget build(BuildContext context) {
    return Card(
      margin: const EdgeInsets.symmetric(vertical: 10, horizontal: 5),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Row(children: [
                  CircleAvatar(
                    radius: 25,
                    backgroundImage: NetworkImage(widget.avatarUri),
                  ),
                ]),
                const SizedBox(width: 10),
                Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Row(children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.end,
                        children: [
                          Text(
                            widget.name,
                            style: GoogleFonts.quicksand(
                                fontWeight: FontWeight.bold, fontSize: 16),
                            overflow: TextOverflow.ellipsis,
                          ),
                          IconButton(
                            icon: const Icon(Icons.abc),
                            onPressed: () => _handleOnPostOptions(),
                          ),
                        ],
                      )
                    ]),
                    Row(
                      children: [
                        Text(
                            timeago.format(DateTime.fromMillisecondsSinceEpoch(
                                widget.postedTime)),
                            style: quickSand12Font),
                        const SizedBox(width: 5),
                        Icon(
                            widget.privacy == "public"
                                ? Icons.public
                                : Icons.lock,
                            size: 15,
                            color: Colors.grey[600]),
                      ],
                    ),
                  ],
                ),
              ],
            ),
            const SizedBox(height: 10),
            Text("${widget.content}"),
            const SizedBox(height: 5),
            (widget.postImageUri) != null
                ? ClipRRect(
                    borderRadius: BorderRadius.all(Radius.circular(10)),
                    child: Image.network(
                      '${widget.postImageUri!}',
                      width: double.infinity,
                      height: 300,
                      fit: BoxFit.fitHeight,
                    ))
                : const SizedBox(height: 10),
            Row(
              mainAxisAlignment: MainAxisAlignment.spaceBetween,
              children: [
                Row(
                  children: [
                    IconButton(
                      icon: Icon(
                          isLiked ? Icons.favorite : Icons.favorite_border,
                          color: isLiked ? primaryColor : Colors.grey),
                      onPressed: () {
                        setState(() {
                          isLiked = !isLiked;
                        });
                      },
                    ),
                    Text("${widget.nol + (isLiked ? 1 : 0)} Likes",
                        style: quickSand15Font),
                    SizedBox(
                      width: 10,
                    ),
                    Icon(
                      Icons.comment_sharp,
                      color: Colors.grey,
                      size: 20,
                    ),
                    const SizedBox(width: 5),
                    Text("${widget.noc}"),
                  ],
                ),
                Row(
                  children: [],
                ),
              ],
            ),
          ],
        ),
      ),
    );
  }

  _handleOnPostOptions() {}
}
