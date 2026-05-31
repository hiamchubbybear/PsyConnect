import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/post.dart';
import 'package:PsyConnect/models/comment.dart';
import 'package:PsyConnect/provider/post_provider.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:timeago/timeago.dart' as timeago;

class PostDetailScreen extends StatefulWidget {
  final Post post;

  const PostDetailScreen({
    super.key,
    required this.post,
  });

  @override
  State<PostDetailScreen> createState() => _PostDetailScreenState();
}

class _PostDetailScreenState extends State<PostDetailScreen> {
  final TextEditingController _commentController = TextEditingController();
  final FocusNode _focusNode = FocusNode();
  
  Comment? _replyingToComment;
  bool _isSubmitting = false;

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<PostProvider>(context, listen: false).fetchComments(widget.post.id);
    });
  }

  @override
  void dispose() {
    _commentController.dispose();
    _focusNode.dispose();
    super.dispose();
  }

  void _setReplyTarget(Comment comment) {
    setState(() {
      _replyingToComment = comment;
    });
    _focusNode.requestFocus();
  }

  void _clearReplyTarget() {
    setState(() {
      _replyingToComment = null;
    });
  }

  Future<void> _submitComment(PostProvider postProvider) async {
    final text = _commentController.text.trim();
    if (text.isEmpty) return;

    setState(() {
      _isSubmitting = true;
    });

    final success = await postProvider.addComment(
      postId: widget.post.id,
      content: text,
      parentCommentId: _replyingToComment?.id,
    );

    if (success && mounted) {
      _commentController.clear();
      _focusNode.unfocus();
      _clearReplyTarget();
    }

    setState(() {
      _isSubmitting = false;
    });
  }

  @override
  Widget build(BuildContext context) {
    final postProvider = Provider.of<PostProvider>(context);
    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    final theme = Theme.of(context);

    // Parse date
    DateTime postTime;
    try {
      postTime = DateTime.parse(widget.post.createdAt).toLocal();
    } catch (_) {
      postTime = DateTime.now();
    }

    final avatarUrl = widget.post.authorAvatar ?? 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';
    final authorName = widget.post.authorName ?? 'Anonymous';

    return Scaffold(
      backgroundColor: isDark ? Colors.black : const Color(0xFFF7F9FC),
      appBar: AppBar(
        title: Text(
          "Post Discussion",
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold,
            fontSize: 18,
            color: theme.textTheme.titleLarge?.color,
          ),
        ),
        elevation: 0.5,
        backgroundColor: isDark ? theme.cardColor : Colors.white,
        leading: IconButton(
          icon: Icon(Icons.arrow_back_ios_new_rounded, size: 20, color: theme.iconTheme.color),
          onPressed: () => Navigator.pop(context),
        ),
      ),
      body: Column(
        children: [
          Expanded(
            child: CustomScrollView(
              slivers: [
                // 1. Post Header Box
                SliverToBoxAdapter(
                  child: Container(
                    margin: const EdgeInsets.only(bottom: 12),
                    padding: const EdgeInsets.all(16),
                    color: isDark ? theme.cardColor : Colors.white,
                    child: Column(
                      crossAxisAlignment: CrossAxisAlignment.start,
                      children: [
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
                        const SizedBox(height: 16),
                        if (widget.post.title.isNotEmpty)
                          Padding(
                            padding: const EdgeInsets.only(bottom: 8),
                            child: Text(
                              widget.post.title,
                              style: GoogleFonts.quicksand(
                                fontSize: 18,
                                fontWeight: FontWeight.bold,
                                color: isDark ? Colors.white : Colors.black87,
                              ),
                            ),
                          ),
                        Text(
                          widget.post.content,
                          style: GoogleFonts.quicksand(
                            fontSize: 14,
                            height: 1.5,
                            color: isDark ? Colors.white70 : Colors.black87,
                          ),
                        ),
                        if (widget.post.categories.isNotEmpty || widget.post.tags.isNotEmpty) ...[
                          const SizedBox(height: 16),
                          Wrap(
                            spacing: 6,
                            runSpacing: 4,
                            children: [
                              ...widget.post.categories.map((c) => _buildBadge(c, Colors.teal, isDark)),
                              ...widget.post.tags.map((t) => _buildBadge("#$t", Colors.blue, isDark)),
                            ],
                          ),
                        ],
                      ],
                    ),
                  ),
                ),

                // 2. Comments Count Header
                SliverToBoxAdapter(
                  child: Padding(
                    padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 8),
                    child: Text(
                      "Comments (${postProvider.activeComments.length})",
                      style: GoogleFonts.quicksand(
                        fontSize: 15,
                        fontWeight: FontWeight.bold,
                        color: isDark ? Colors.white70 : Colors.black87,
                      ),
                    ),
                  ),
                ),

                // 3. Nested Comments Tree list
                postProvider.isLoadingComments
                    ? const SliverFillRemaining(
                        child: Center(child: CircularProgressIndicator()),
                      )
                    : postProvider.activeComments.isEmpty
                        ? SliverFillRemaining(
                            hasScrollBody: false,
                            child: Center(
                              child: Column(
                                mainAxisAlignment: MainAxisAlignment.center,
                                children: [
                                  Icon(Icons.forum_outlined, size: 40, color: Colors.grey[400]),
                                  const SizedBox(height: 12),
                                  Text(
                                    "No comments yet. Write the first response!",
                                    style: GoogleFonts.quicksand(
                                      color: Colors.grey[500],
                                      fontSize: 13,
                                      fontWeight: FontWeight.w500,
                                    ),
                                  ),
                                ],
                              ),
                            ),
                          )
                        : SliverList(
                            delegate: SliverChildBuilderDelegate(
                              (context, index) {
                                final rootComment = postProvider.activeComments[index];
                                return _buildCommentTree(rootComment, 0, isDark);
                              },
                              childCount: postProvider.activeComments.length,
                            ),
                          ),
              ],
            ),
          ),

          // 4. Persistent Custom Bottom Bar input dock
          SafeArea(
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12, vertical: 8),
              decoration: BoxDecoration(
                color: isDark ? theme.cardColor : Colors.white,
                border: Border(
                  top: BorderSide(
                    color: isDark ? Colors.grey[850]! : Colors.grey[200]!,
                    width: 1,
                  ),
                ),
              ),
              child: Column(
                mainAxisSize: MainAxisSize.min,
                children: [
                  // Reply Target indicator pill
                  if (_replyingToComment != null)
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 10, vertical: 6),
                      margin: const EdgeInsets.only(bottom: 8),
                      decoration: BoxDecoration(
                        color: isDark ? Colors.blueGrey[900] : primaryColor.withOpacity(0.08),
                        borderRadius: BorderRadius.circular(20),
                      ),
                      child: Row(
                        children: [
                          Icon(Icons.reply, size: 14, color: isDark ? Colors.blue[300] : primaryColor),
                          const SizedBox(width: 6),
                          Expanded(
                            child: Text(
                              "Replying to @${_replyingToComment!.authorName ?? 'User'}",
                              style: GoogleFonts.quicksand(
                                fontSize: 11,
                                fontWeight: FontWeight.bold,
                                color: isDark ? Colors.blue[300] : primaryColor,
                              ),
                            ),
                          ),
                          GestureDetector(
                            onTap: _clearReplyTarget,
                            child: Icon(Icons.close, size: 14, color: isDark ? Colors.white70 : Colors.black54),
                          )
                        ],
                      ),
                    ),

                  Row(
                    children: [
                      Expanded(
                        child: TextField(
                          controller: _commentController,
                          focusNode: _focusNode,
                          maxLines: null,
                          style: GoogleFonts.quicksand(
                            fontSize: 14,
                            color: isDark ? Colors.white : Colors.black87,
                          ),
                          decoration: InputDecoration(
                            hintText: _replyingToComment != null ? "Write a reply..." : "Write a comment...",
                            hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontSize: 13),
                            border: OutlineInputBorder(
                              borderRadius: BorderRadius.circular(24),
                              borderSide: BorderSide.none,
                            ),
                            fillColor: isDark ? Colors.grey[900] : Colors.grey[100],
                            filled: true,
                            contentPadding: const EdgeInsets.symmetric(horizontal: 16, vertical: 10),
                          ),
                        ),
                      ),
                      const SizedBox(width: 8),
                      _isSubmitting
                          ? const SizedBox(
                              height: 24,
                              width: 24,
                              child: CircularProgressIndicator(strokeWidth: 2),
                            )
                          : Container(
                              decoration: BoxDecoration(
                                color: isDark ? Colors.blueGrey[800] : primaryColor,
                                shape: BoxShape.circle,
                              ),
                              child: IconButton(
                                icon: const Icon(Icons.send, color: Colors.white, size: 18),
                                onPressed: () => _submitComment(postProvider),
                              ),
                            ),
                    ],
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildCommentTree(Comment comment, int depth, bool isDark) {
    return Column(
      crossAxisAlignment: CrossAxisAlignment.stretch,
      children: [
        _buildCommentCard(comment, depth, isDark),
        if (comment.replies.isNotEmpty)
          ...comment.replies.map((reply) => _buildCommentTree(reply, depth + 1, isDark)),
      ],
    );
  }

  Widget _buildCommentCard(Comment comment, int depth, bool isDark) {
    // Indent based on depth level
    final double leftPadding = 16.0 + (depth * 20.0);
    final avatarUrl = comment.authorAvatar ?? 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';
    
    DateTime commentTime;
    try {
      commentTime = DateTime.parse(comment.createdAt).toLocal();
    } catch (_) {
      commentTime = DateTime.now();
    }

    return Padding(
      padding: EdgeInsets.only(left: leftPadding, right: 16, top: 10, bottom: 6),
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          // Visual connector lines for nested depths
          if (depth > 0)
            Container(
              width: 2,
              height: 36,
              margin: const EdgeInsets.only(right: 12),
              color: isDark ? Colors.grey[800] : Colors.grey[300],
            ),

          CircleAvatar(
            radius: 14,
            backgroundImage: NetworkImage(avatarUrl),
            backgroundColor: isDark ? Colors.grey[800] : Colors.grey[300],
          ),
          const SizedBox(width: 10),
          
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                // Author & Time header
                Row(
                  children: [
                    Text(
                      comment.authorName ?? 'User',
                      style: GoogleFonts.quicksand(
                        fontWeight: FontWeight.bold,
                        fontSize: 12,
                        color: isDark ? Colors.white : Colors.black87,
                      ),
                    ),
                    const SizedBox(width: 8),
                    Text(
                      timeago.format(commentTime),
                      style: GoogleFonts.quicksand(
                        fontSize: 10,
                        color: Colors.grey,
                      ),
                    ),
                  ],
                ),
                const SizedBox(height: 4),

                // Comment Content
                Text(
                  comment.content,
                  style: GoogleFonts.quicksand(
                    fontSize: 12,
                    height: 1.4,
                    color: isDark ? Colors.white70 : Colors.black87,
                  ),
                ),
                
                // Reply Action Button
                GestureDetector(
                  onTap: () => _setReplyTarget(comment),
                  child: Padding(
                    padding: const EdgeInsets.symmetric(vertical: 4),
                    child: Text(
                      "Reply",
                      style: GoogleFonts.quicksand(
                        fontSize: 11,
                        fontWeight: FontWeight.bold,
                        color: isDark ? Colors.blue[300] : primaryColor,
                      ),
                    ),
                  ),
                ),
              ],
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildBadge(String text, Color color, bool isDark) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 2),
      decoration: BoxDecoration(
        color: color.withOpacity(0.08),
        borderRadius: BorderRadius.circular(4),
        border: Border.all(color: color.withOpacity(0.15), width: 0.5),
      ),
      child: Text(
        text,
        style: GoogleFonts.quicksand(
          fontSize: 9,
          fontWeight: FontWeight.bold,
          color: isDark ? color.withOpacity(0.8) : color,
        ),
      ),
    );
  }
}
