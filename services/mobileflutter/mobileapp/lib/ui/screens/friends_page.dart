import 'package:flutter/cupertino.dart';
import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:PsyConnect/models/friend.dart';
import 'package:PsyConnect/provider/friend_provider.dart';
import 'package:PsyConnect/provider/chat_provider.dart';
import 'package:PsyConnect/ui/screens/chat_detail_screen.dart';
import 'package:PsyConnect/ui/widgets/common/custom_button.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';

class FriendsPage extends StatefulWidget {
  const FriendsPage({super.key});

  @override
  State<FriendsPage> createState() => _FriendsPageState();
}

class _FriendsPageState extends State<FriendsPage> with SingleTickerProviderStateMixin {
  late TabController _tabController;
  final TextEditingController _searchController = TextEditingController();
  String _searchQuery = "";

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 3, vsync: this);
    
    // Fetch all friends data on load
    WidgetsBinding.instance.addPostFrameCallback((_) {
      Provider.of<FriendProvider>(context, listen: false).fetchAllSocialData();
    });

    _searchController.addListener(() {
      setState(() {
        _searchQuery = _searchController.text.toLowerCase();
      });
    });
  }

  @override
  void dispose() {
    _tabController.dispose();
    _searchController.dispose();
    super.dispose();
  }

  Future<void> _refreshData() async {
    await Provider.of<FriendProvider>(context, listen: false).fetchAllSocialData();
  }

  void _startChatWithFriend(Friend friend) async {
    final chatProvider = Provider.of<ChatProvider>(context, listen: false);
    
    // Show loading
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (context) => const Center(
        child: CircularProgressIndicator(),
      ),
    );

    try {
      final conversationId = await chatProvider.getOrCreateConversation(friend.profileId);
      Navigator.pop(context); // Dismiss loading dialog

      if (conversationId != null && conversationId.isNotEmpty) {
        Navigator.push(
          context,
          MaterialPageRoute(
            builder: (context) => ChatDetailScreen(
              conversationId: conversationId,
              companionId: friend.profileId,
              companionName: friend.fullName,
              companionAvatar: friend.avatarUri.isNotEmpty
                  ? friend.avatarUri
                  : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg',
            ),
          ),
        );
      } else {
        ToastService.showToast(
          context: context,
          message: "Failed to open conversation channel. Please try again.",
          title: "Error",
          type: ToastType.error,
        );
      }
    } catch (e) {
      Navigator.pop(context); // Dismiss loading dialog
      ToastService.showToast(
        context: context,
        message: "Failed to open chat: $e",
        title: "Error",
        type: ToastType.error,
      );
    }
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final isDark = theme.brightness == Brightness.dark;
    final primaryColor = isDark ? Colors.blue[300]! : Colors.blue;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        title: Text(
          "Social Connections",
          style: GoogleFonts.quicksand(fontWeight: FontWeight.bold),
        ),
        centerTitle: true,
        bottom: TabBar(
          controller: _tabController,
          labelColor: primaryColor,
          unselectedLabelColor: isDark ? Colors.grey[400] : Colors.grey[600],
          indicatorColor: primaryColor,
          labelStyle: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 15),
          unselectedLabelStyle: GoogleFonts.quicksand(fontWeight: FontWeight.w600, fontSize: 15),
          tabs: const [
            Tab(text: "Friends"),
            Tab(text: "Pending"),
            Tab(text: "Suggestions"),
          ],
        ),
      ),
      body: Consumer<FriendProvider>(
        builder: (context, friendProvider, child) {
          if (friendProvider.isLoading && friendProvider.friends.isEmpty && friendProvider.requests.isEmpty) {
            return const Center(child: CircularProgressIndicator());
          }

          if (friendProvider.errorMessage != null) {
            return Center(
              child: Padding(
                padding: const EdgeInsets.all(24.0),
                child: Column(
                  mainAxisAlignment: MainAxisAlignment.center,
                  children: [
                    Text(
                      "Error: ${friendProvider.errorMessage}",
                      textAlign: TextAlign.center,
                      style: GoogleFonts.quicksand(color: Colors.redAccent, fontSize: 16),
                    ),
                    const SizedBox(height: 16),
                    CustomButton(
                      onPressed: _refreshData,
                      text: "Retry",
                      width: 150,
                      height: 44,
                    )
                  ],
                ),
              ),
            );
          }

          return TabBarView(
            controller: _tabController,
            children: [
              _buildFriendsTab(friendProvider, theme, isDark, primaryColor),
              _buildPendingTab(friendProvider, theme, isDark, primaryColor),
              _buildSuggestionsTab(friendProvider, theme, isDark, primaryColor),
            ],
          );
        },
      ),
    );
  }

  Widget _buildFriendsTab(
    FriendProvider friendProvider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    final filteredFriends = friendProvider.friends.where((f) {
      return f.fullName.toLowerCase().contains(_searchQuery);
    }).toList();

    return RefreshIndicator(
      onRefresh: _refreshData,
      child: Column(
        children: [
          _buildSearchBar(theme, isDark),
          Expanded(
            child: filteredFriends.isEmpty
                ? _buildEmptyState(
                    icon: CupertinoIcons.person_3,
                    title: "No Friends Found",
                    subtitle: _searchQuery.isEmpty
                        ? "Check out suggestions to find and connect with friends!"
                        : "No matching friends found for '$_searchQuery'",
                  )
                : ListView.builder(
                    padding: const EdgeInsets.all(16),
                    itemCount: filteredFriends.length,
                    itemBuilder: (context, index) {
                      final friend = filteredFriends[index];
                      return _buildFriendCard(friend, friendProvider, theme, isDark, primaryColor);
                    },
                  ),
          ),
        ],
      ),
    );
  }

  Widget _buildPendingTab(
    FriendProvider friendProvider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    return RefreshIndicator(
      onRefresh: _refreshData,
      child: friendProvider.requests.isEmpty
          ? _buildEmptyState(
              icon: CupertinoIcons.envelope,
              title: "No Pending Invites",
              subtitle: "You don't have any incoming friend request invitations at the moment.",
            )
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: friendProvider.requests.length,
              itemBuilder: (context, index) {
                final request = friendProvider.requests[index];
                return _buildPendingCard(request, friendProvider, theme, isDark, primaryColor);
              },
            ),
    );
  }

  Widget _buildSuggestionsTab(
    FriendProvider friendProvider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    return RefreshIndicator(
      onRefresh: _refreshData,
      child: friendProvider.suggestions.isEmpty
          ? _buildEmptyState(
              icon: CupertinoIcons.sparkles,
              title: "No Suggestions Yet",
              subtitle: "We couldn't find any recommendations right now. Try updating your interests!",
            )
          : ListView.builder(
              padding: const EdgeInsets.all(16),
              itemCount: friendProvider.suggestions.length,
              itemBuilder: (context, index) {
                final suggestion = friendProvider.suggestions[index];
                return _buildSuggestionCard(suggestion, friendProvider, theme, isDark, primaryColor);
              },
            ),
    );
  }

  Widget _buildSearchBar(ThemeData theme, bool isDark) {
    return Padding(
      padding: const EdgeInsets.fromLTRB(16, 16, 16, 8),
      child: TextField(
        controller: _searchController,
        style: GoogleFonts.quicksand(fontSize: 16),
        decoration: InputDecoration(
          hintText: "Search friends...",
          hintStyle: GoogleFonts.quicksand(color: Colors.grey[500]),
          prefixIcon: Icon(CupertinoIcons.search, color: Colors.grey[500]),
          filled: true,
          fillColor: isDark ? Colors.grey[850] : Colors.grey[100],
          contentPadding: const EdgeInsets.symmetric(vertical: 0, horizontal: 16),
          border: OutlineInputBorder(
            borderRadius: BorderRadius.circular(30),
            borderSide: BorderSide.none,
          ),
        ),
      ),
    );
  }

  Widget _buildEmptyState({
    required IconData icon,
    required String title,
    required String subtitle,
  }) {
    return Center(
      child: SingleChildScrollView(
        physics: const AlwaysScrollableScrollPhysics(),
        padding: const EdgeInsets.all(32),
        child: Column(
          mainAxisAlignment: MainAxisAlignment.center,
          children: [
            Icon(icon, size: 64, color: Colors.grey[400]),
            const SizedBox(height: 16),
            Text(
              title,
              style: GoogleFonts.quicksand(
                fontSize: 20,
                fontWeight: FontWeight.bold,
              ),
            ),
            const SizedBox(height: 8),
            Text(
              subtitle,
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

  Widget _buildFriendCard(
    Friend friend,
    FriendProvider provider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    final avatarUrl = friend.avatarUri.isNotEmpty
        ? friend.avatarUri
        : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

    return Card(
      elevation: 2,
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Row(
          children: [
            CircleAvatar(
              radius: 28,
              backgroundImage: NetworkImage(avatarUrl),
              backgroundColor: isDark ? Colors.grey[800] : Colors.grey[200],
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    friend.fullName,
                    style: GoogleFonts.quicksand(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (friend.role != null && friend.role!.isNotEmpty) ...[
                    const SizedBox(height: 2),
                    Container(
                      padding: const EdgeInsets.symmetric(horizontal: 6, vertical: 2),
                      decoration: BoxDecoration(
                        color: primaryColor.withOpacity(0.1),
                        borderRadius: BorderRadius.circular(4),
                      ),
                      child: Text(
                        friend.role!.toUpperCase(),
                        style: GoogleFonts.quicksand(
                          fontSize: 10,
                          fontWeight: FontWeight.bold,
                          color: primaryColor,
                        ),
                      ),
                    ),
                  ],
                ],
              ),
            ),
            IconButton(
              icon: Icon(CupertinoIcons.chat_bubble, color: primaryColor),
              onPressed: () => _startChatWithFriend(friend),
            ),
            IconButton(
              icon: const Icon(CupertinoIcons.person_badge_minus, color: Colors.redAccent),
              onPressed: () => _confirmUnfriend(context, friend, provider),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildPendingCard(
    Friend request,
    FriendProvider provider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    final avatarUrl = request.avatarUri.isNotEmpty
        ? request.avatarUri
        : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

    return Card(
      elevation: 2,
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Column(
          children: [
            Row(
              children: [
                CircleAvatar(
                  radius: 28,
                  backgroundImage: NetworkImage(avatarUrl),
                  backgroundColor: isDark ? Colors.grey[800] : Colors.grey[200],
                ),
                const SizedBox(width: 12),
                Expanded(
                  child: Column(
                    crossAxisAlignment: CrossAxisAlignment.start,
                    children: [
                      Text(
                        request.fullName,
                        style: GoogleFonts.quicksand(
                          fontSize: 16,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                      Text(
                        "Wants to connect with you",
                        style: GoogleFonts.quicksand(
                          fontSize: 12,
                          color: Colors.grey[500],
                        ),
                      ),
                    ],
                  ),
                ),
              ],
            ),
            const SizedBox(height: 12),
            Row(
              mainAxisAlignment: MainAxisAlignment.end,
              children: [
                OutlinedButton(
                  onPressed: () async {
                    final success = await provider.cancelRequest(request.profileId);
                    if (success && mounted) {
                      ToastService.showToast(
                        context: context,
                        message: "Declined request from ${request.fullName}",
                        title: "Declined",
                        type: ToastType.info,
                      );
                    }
                  },
                  style: OutlinedButton.styleFrom(
                    side: const BorderSide(color: Colors.grey),
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                  ),
                  child: Text(
                    "Decline",
                    style: GoogleFonts.quicksand(color: isDark ? Colors.white : Colors.black87),
                  ),
                ),
                const SizedBox(width: 12),
                ElevatedButton(
                  onPressed: () async {
                    final success = await provider.acceptRequest(request.profileId);
                    if (success && mounted) {
                      ToastService.showToast(
                        context: context,
                        message: "You are now friends with ${request.fullName}!",
                        title: "Connected",
                        type: ToastType.success,
                      );
                    } else if (mounted) {
                      ToastService.showToast(
                        context: context,
                        message: provider.errorMessage ?? "Failed to accept connection invitation",
                        title: "Error",
                        type: ToastType.error,
                      );
                    }
                  },
                  style: ElevatedButton.styleFrom(
                    backgroundColor: primaryColor,
                    shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
                  ),
                  child: Text(
                    "Accept",
                    style: GoogleFonts.quicksand(color: Colors.white, fontWeight: FontWeight.bold),
                  ),
                ),
              ],
            )
          ],
        ),
      ),
    );
  }

  Widget _buildSuggestionCard(
    Friend suggestion,
    FriendProvider provider,
    ThemeData theme,
    bool isDark,
    Color primaryColor,
  ) {
    final avatarUrl = suggestion.avatarUri.isNotEmpty
        ? suggestion.avatarUri
        : 'https://i.pinimg.com/736x/83/21/ec/8321ec3e2ed58da8e46f1926f10373dc.jpg';

    return Card(
      elevation: 2,
      margin: const EdgeInsets.only(bottom: 12),
      shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(16)),
      child: Padding(
        padding: const EdgeInsets.all(12.0),
        child: Row(
          children: [
            CircleAvatar(
              radius: 28,
              backgroundImage: NetworkImage(avatarUrl),
              backgroundColor: isDark ? Colors.grey[800] : Colors.grey[200],
            ),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    suggestion.fullName,
                    style: GoogleFonts.quicksand(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  if (suggestion.role != null && suggestion.role!.isNotEmpty) ...[
                    const SizedBox(height: 2),
                    Text(
                      suggestion.role!.toUpperCase(),
                      style: GoogleFonts.quicksand(
                        fontSize: 11,
                        color: primaryColor,
                        fontWeight: FontWeight.w600,
                      ),
                    ),
                  ],
                ],
              ),
            ),
            ElevatedButton(
              onPressed: () async {
                final success = await provider.sendFriendRequest(suggestion.profileId);
                if (success && mounted) {
                  ToastService.showToast(
                    context: context,
                    message: "Friend request sent to ${suggestion.fullName}!",
                    title: "Request Sent",
                    type: ToastType.success,
                  );
                } else if (mounted) {
                  ToastService.showToast(
                    context: context,
                    message: provider.errorMessage ?? "Failed to send friend request",
                    title: "Error",
                    type: ToastType.error,
                  );
                }
              },
              style: ElevatedButton.styleFrom(
                backgroundColor: primaryColor,
                shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(8)),
              ),
              child: Text(
                "Connect",
                style: GoogleFonts.quicksand(color: Colors.white, fontWeight: FontWeight.bold),
              ),
            ),
          ],
        ),
      ),
    );
  }

  void _confirmUnfriend(BuildContext context, Friend friend, FriendProvider provider) {
    showCupertinoDialog(
      context: context,
      builder: (context) => CupertinoAlertDialog(
        title: Text(
          "Unfriend?",
          style: GoogleFonts.quicksand(fontWeight: FontWeight.bold),
        ),
        content: Text(
          "Are you sure you want to unfriend ${friend.fullName}?",
          style: GoogleFonts.quicksand(),
        ),
        actions: [
          CupertinoDialogAction(
            child: Text("Cancel", style: GoogleFonts.quicksand()),
            onPressed: () => Navigator.pop(context),
          ),
          CupertinoDialogAction(
            isDestructiveAction: true,
            child: Text("Unfriend", style: GoogleFonts.quicksand(fontWeight: FontWeight.bold)),
            onPressed: () async {
              Navigator.pop(context);
              final success = await provider.removeFriend(friend.profileId);
              if (success && mounted) {
                ToastService.showToast(
                  context: context,
                  message: "Successfully unfriended ${friend.fullName}",
                  title: "Removed",
                  type: ToastType.info,
                );
              } else if (mounted) {
                ToastService.showToast(
                  context: context,
                  message: provider.errorMessage ?? "Failed to unfriend",
                  title: "Error",
                  type: ToastType.error,
                );
              }
            },
          ),
        ],
      ),
    );
  }
}
