import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/post.dart';
import 'package:PsyConnect/provider/post_provider.dart';
import 'package:PsyConnect/ui/screens/post_detail_screen.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:timeago/timeago.dart' as timeago;

class PostWidget extends StatelessWidget {
  final bool isDark;
  final Post post;

  const PostWidget({
    super.key,
    required this.isDark,
    required this.post,
  });

  @override
  Widget build(BuildContext context) {
    final postProvider = Provider.of<PostProvider>(context);
    final theme = Theme.of(context);
    
    // Check if user has liked (upvoted) this post
    final userVote = postProvider.getUserVote(post.id);
    final isUpvoted = userVote == 'up';
    final isDownvoted = userVote == 'down';

    // Parse date
    DateTime postTime;
    try {
      postTime = DateTime.parse(post.createdAt).toLocal();
    } catch (_) {
      postTime = DateTime.now();
    }

    final avatarUrl = post.authorAvatar ?? 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';
    final authorName = post.authorName ?? 'Anonymous';

    return GestureDetector(
      onTap: () {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => PostDetailScreen(post: post),
          ),
        );
      },
      child: Card(
        color: isDark ? theme.cardColor : whiteColor,
        margin: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
        elevation: 0.5,
        shape: RoundedRectangleBorder(
          borderRadius: BorderRadius.circular(12),
          side: BorderSide(
            color: isDark ? Colors.grey[850]! : Colors.grey[200]!,
            width: 1,
          ),
        ),
        child: Padding(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              // Author Header
              Row(
                children: [
                  CircleAvatar(
                    radius: 20,
                    backgroundImage: NetworkImage(avatarUrl),
                    backgroundColor: isDark ? Colors.grey[800] : Colors.grey[300],
                  ),
                  const SizedBox(width: 12),
                  Expanded(
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
                        Text(
                          authorName,
                          style: GoogleFonts.quicksand(
                            fontWeight: FontWeight.bold,
                            fontSize: 14,
                            color: isDark ? Colors.white : Colors.black87,
                          ),
                        ),
                        const SizedBox(height: 2),
                        Row(
                          children: [
                            Text(
                              timeago.format(postTime),
                              style: GoogleFonts.quicksand(
                                fontSize: 11,
                                fontWeight: FontWeight.w400,
                                color: isDark ? Colors.grey[400] : Colors.grey[600],
                              ),
                            ),
                            const SizedBox(width: 8),
                            Icon(
                              Icons.public_rounded,
                              size: 13,
                              color: isDark ? Colors.grey[600] : Colors.grey[500],
                            ),
                          ],
                        ),
                      ],
                    ),
                  ),
                ],
              ),
              
              const SizedBox(height: 12),
              
              // Post Title
              if (post.title.isNotEmpty)
                Padding(
                  padding: const EdgeInsets.only(bottom: 6),
                  child: Text(
                    post.title,
                    style: GoogleFonts.quicksand(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: isDark ? Colors.white : Colors.black87,
                    ),
                  ),
                ),

              // Post Content
              Text(
                post.content,
                style: GoogleFonts.quicksand(
                  fontSize: 13,
                  height: 1.4,
                  fontWeight: FontWeight.w400,
                  color: isDark ? Colors.white70 : Colors.black87,
                ),
              ),

              // Post Categories / Tags chips
              if (post.categories.isNotEmpty || post.tags.isNotEmpty) ...[
                const SizedBox(height: 12),
                Wrap(
                  spacing: 6,
                  runSpacing: 4,
                  children: [
                    ...post.categories.map((c) => _buildBadge(c, Colors.teal, isDark)),
                    ...post.tags.map((t) => _buildBadge("#$t", Colors.blue, isDark)),
                  ],
                ),
              ],

              const SizedBox(height: 8),
              const Divider(height: 16),

              // Footer interaction buttons (Upvote, Downvote, Comments)
              Row(
                mainAxisAlignment: MainAxisAlignment.spaceBetween,
                children: [
                  Row(
                    children: [
                      // Upvote
                      IconButton(
                        icon: Icon(
                          isUpvoted ? Icons.thumb_up_alt_rounded : Icons.thumb_up_off_alt_rounded,
                          color: isUpvoted ? Colors.blue : Colors.grey,
                          size: 20,
                        ),
                        onPressed: () {
                          postProvider.toggleVote(post.id, 'up');
                        },
                      ),
                      Text(
                        "${post.upvoteCount}",
                        style: GoogleFonts.quicksand(
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                          color: isDark ? Colors.white70 : Colors.black87,
                        ),
                      ),
                      
                      const SizedBox(width: 8),

                      // Downvote
                      IconButton(
                        icon: Icon(
                          isDownvoted ? Icons.thumb_down_alt_rounded : Icons.thumb_down_off_alt_rounded,
                          color: isDownvoted ? Colors.red : Colors.grey,
                          size: 20,
                        ),
                        onPressed: () {
                          postProvider.toggleVote(post.id, 'down');
                        },
                      ),
                      Text(
                        "${post.downvoteCount}",
                        style: GoogleFonts.quicksand(
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                          color: isDark ? Colors.white70 : Colors.black87,
                        ),
                      ),

                      const SizedBox(width: 16),

                      // Comments
                      IconButton(
                        icon: Icon(
                          Icons.mode_comment_outlined,
                          color: isDark ? Colors.grey : Colors.grey[600],
                          size: 18,
                        ),
                        onPressed: () {
                          Navigator.push(
                            context,
                            MaterialPageRoute(
                              builder: (context) => PostDetailScreen(post: post),
                            ),
                          );
                        },
                      ),
                      Text(
                        "${post.commentCount}",
                        style: GoogleFonts.quicksand(
                          fontSize: 12,
                          fontWeight: FontWeight.w600,
                          color: isDark ? Colors.white70 : Colors.black87,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildBadge(String text, Color color, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 3),
      decoration: BoxDecoration(
        color: color.withOpacity(0.08),
        borderRadius: BorderRadius.circular(6),
        border: Border.all(color: color.withOpacity(0.2), width: 0.5),
      ),
      child: Text(
        text,
        style: GoogleFonts.quicksand(
          fontSize: 10,
          fontWeight: FontWeight.bold,
          color: isDark ? color.withOpacity(0.8) : color,
        ),
      ),
    );
  }
}
