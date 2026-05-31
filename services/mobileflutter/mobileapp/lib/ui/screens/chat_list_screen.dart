import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:timeago/timeago.dart' as timeago;
import 'package:PsyConnect/provider/chat_provider.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/ui/screens/chat_detail_screen.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

class ChatListScreen extends StatefulWidget {
  const ChatListScreen({super.key});

  @override
  State<ChatListScreen> createState() => _ChatListScreenState();
}

class _ChatListScreenState extends State<ChatListScreen> {
  @override
  void initState() {
    super.initState();
    // Fetch conversations as soon as screen is initialized
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<ChatProvider>(context, listen: false).fetchRecentConversations();
    });
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chatProvider = Provider.of<ChatProvider>(context);
    final primaryColor = theme.primaryColor;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: Text(
          'Chats',
          style: GoogleFonts.quicksand(
            fontWeight: FontWeight.bold,
            fontSize: 24,
            color: theme.appBarTheme.titleTextStyle?.color ?? theme.textTheme.titleLarge?.color,
          ),
        ),
        centerTitle: false,
        backgroundColor: Colors.transparent,
        elevation: 0,
        actions: [
          IconButton(
            icon: Icon(Icons.refresh_rounded, color: primaryColor),
            onPressed: () => chatProvider.fetchRecentConversations(),
          ),
        ],
      ),
      body: Column(
        children: [
          // Elegant search bar
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
            child: Container(
              decoration: BoxDecoration(
                color: theme.cardColor,
                borderRadius: BorderRadius.circular(16.0),
                boxShadow: [
                  BoxShadow(
                    color: Colors.black.withOpacity(0.03),
                    blurRadius: 10,
                    offset: const Offset(0, 4),
                  )
                ],
              ),
              child: TextField(
                style: GoogleFonts.quicksand(fontWeight: FontWeight.w500),
                decoration: InputDecoration(
                  hintText: 'Search chats...',
                  hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontWeight: FontWeight.w400),
                  prefixIcon: const Icon(Icons.search_rounded, color: Colors.grey),
                  border: InputBorder.none,
                  contentPadding: const EdgeInsets.symmetric(vertical: 14.0),
                ),
              ),
            ),
          ),
          
          Expanded(
            child: chatProvider.isLoadingConversations && chatProvider.conversations.isEmpty
                ? Center(
                    child: LoadingAnimationWidget.threeArchedCircle(
                      color: primaryColor,
                      size: 40,
                    ),
                  )
                : RefreshIndicator(
                    onRefresh: () => chatProvider.fetchRecentConversations(),
                    color: primaryColor,
                    child: chatProvider.conversations.isEmpty
                        ? _buildEmptyState(theme, primaryColor)
                        : ListView.separated(
                            padding: const EdgeInsets.all(16.0),
                            itemCount: chatProvider.conversations.length,
                            separatorBuilder: (_, __) => const SizedBox(height: 12.0),
                            itemBuilder: (context, index) {
                              final conv = chatProvider.conversations[index];
                              return _buildConversationItem(context, conv, chatProvider, theme);
                            },
                          ),
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildEmptyState(ThemeData theme, Color primaryColor) {
    return Center(
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Container(
              height: 100,
              width: 100,
              decoration: BoxDecoration(
                color: primaryColor.withOpacity(0.1),
                shape: BoxShape.circle,
              ),
              child: Icon(
                Icons.chat_bubble_outline_rounded,
                size: 50,
                color: primaryColor,
              ),
            ),
            const SizedBox(height: 24),
            Text(
              'No active chats yet',
              style: GoogleFonts.quicksand(
                fontSize: 20,
                fontWeight: FontWeight.bold,
                color: theme.textTheme.bodyLarge?.color,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              'Your messages will appear here once you start\na conversation with a psychologist or doctor.',
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                fontSize: 14,
                color: Colors.grey[500],
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildConversationItem(
    BuildContext context,
    dynamic conv,
    ChatProvider chatProvider,
    ThemeData theme,
  ) {
    final conversationId = conv['id']?.toString() ?? '';
    final participants = conv['participants'] as List<dynamic>? ?? [];
    
    // Companion is the participant ID that isn't the logged-in user
    final companionId = participants.firstWhere(
      (id) => id.toString() != chatProvider.myProfileId,
      orElse: () => '',
    ).toString();

    final lastMessageData = conv['lastMessage'];
    String lastMessageText = '';
    String timeString = '';
    
    if (lastMessageData != null) {
      lastMessageText = lastMessageData['text']?.toString() ?? '';
      final createdAtRaw = lastMessageData['createdAt']?.toString();
      if (createdAtRaw != null) {
        try {
          final dt = DateTime.parse(createdAtRaw).toLocal();
          timeString = timeago.format(dt);
        } catch (_) {}
      }
    } else {
      lastMessageText = 'Start chatting now...';
      final createdAtRaw = conv['createdAt']?.toString();
      if (createdAtRaw != null) {
        try {
          final dt = DateTime.parse(createdAtRaw).toLocal();
          timeString = timeago.format(dt);
        } catch (_) {}
      }
    }

    final UserProfile? companionProfile = chatProvider.cachedProfiles[companionId];
    
    if (companionProfile == null && companionId.isNotEmpty) {
      // Fetch details asynchronously
      chatProvider.fetchUserProfileSilently(companionId);
    }

    final displayName = companionProfile != null
        ? "${companionProfile.firstName ?? ''} ${companionProfile.lastName ?? ''}".trim()
        : 'Loading user...';
    
    final finalDisplayName = displayName.isNotEmpty ? displayName : (companionProfile?.username ?? 'Anonymous');
    
    final avatarUrl = companionProfile?.avatarUri;

    return Container(
      decoration: BoxDecoration(
        color: theme.cardColor,
        borderRadius: BorderRadius.circular(20.0),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.02),
            blurRadius: 8,
            offset: const Offset(0, 4),
          )
        ],
      ),
      child: Material(
        color: Colors.transparent,
        borderRadius: BorderRadius.circular(20.0),
        child: InkWell(
          borderRadius: BorderRadius.circular(20.0),
          onTap: () {
            if (companionId.isEmpty) return;
            Navigator.push(
              context,
              MaterialPageRoute(
                builder: (context) => ChatDetailScreen(
                  conversationId: conversationId,
                  companionId: companionId,
                  companionName: finalDisplayName,
                  companionAvatar: avatarUrl ?? '',
                ),
              ),
            );
          },
          child: Padding(
            padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 14.0),
            child: Row(
              children: [
                // Avatar representation
                _buildAvatar(avatarUrl, finalDisplayName, theme),
                const SizedBox(width: 16),
                
                // Content section
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Row(
                        mainAxisAlignment: MainAxisAlignment.spaceBetween,
                        children: [
                          Flexible(
                            child: Text(
                              finalDisplayName,
                              maxLines: 1,
                              overflow: TextOverflow.ellipsis,
                              style: GoogleFonts.quicksand(
                                fontSize: 16,
                                fontWeight: FontWeight.bold,
                                color: theme.textTheme.bodyLarge?.color,
                              ),
                            ),
                          ),
                          Text(
                            timeString,
                            style: GoogleFonts.quicksand(
                              fontSize: 11,
                              color: Colors.grey,
                              fontWeight: FontWeight.w500,
                            ),
                          ),
                        ],
                      ),
                      const SizedBox(height: 6),
                      Text(
                        lastMessageText,
                        maxLines: 1,
                        overflow: TextOverflow.ellipsis,
                        style: GoogleFonts.quicksand(
                          fontSize: 13.5,
                          fontWeight: lastMessageData != null ? FontWeight.w500 : FontWeight.w400,
                          color: lastMessageData != null ? theme.textTheme.bodyMedium?.color : Colors.grey,
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildAvatar(String? avatarUrl, String displayName, ThemeData theme) {
    if (avatarUrl != null && avatarUrl.isNotEmpty) {
      return CircleAvatar(
        radius: 26,
        backgroundImage: NetworkImage(avatarUrl),
      );
    }
    
    // Generate lovely initials background gradient
    final initial = displayName.isNotEmpty ? displayName[0].toUpperCase() : 'P';
    
    return Container(
      height: 52,
      width: 52,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          colors: [
            theme.primaryColor,
            theme.primaryColor.withBlue(220).withGreen(200),
          ],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
      ),
      child: Center(
        child: Text(
          initial,
          style: GoogleFonts.quicksand(
            fontSize: 20,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      ),
    );
  }
}
