import 'package:flutter/material.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:provider/provider.dart';
import 'package:PsyConnect/provider/chat_provider.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';
import 'package:PsyConnect/ui/screens/call_screen.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';
import 'package:intl/intl.dart';

class ChatDetailScreen extends StatefulWidget {
  final String conversationId;
  final String companionId;
  final String companionName;
  final String companionAvatar;

  const ChatDetailScreen({
    super.key,
    required this.conversationId,
    required this.companionId,
    required this.companionName,
    required this.companionAvatar,
  });

  @override
  State<ChatDetailScreen> createState() => _ChatDetailScreenState();
}

class _ChatDetailScreenState extends State<ChatDetailScreen> {
  final TextEditingController _messageController = TextEditingController();
  final ScrollController _scrollController = ScrollController();

  @override
  void initState() {
    super.initState();
    WidgetsBinding.instance.addPostFrameCallback((_) {
      final chatProvider = Provider.of<ChatProvider>(context, listen: false);
      chatProvider.addListener(_onProviderChange);
      chatProvider.openConversation(
        conversationId: widget.conversationId,
        companionId: widget.companionId,
      );
    });
  }

  void _onProviderChange() {
    if (mounted) {
      Future.delayed(const Duration(milliseconds: 150), _scrollToBottom);
    }
  }

  void _scrollToBottom() {
    if (_scrollController.hasClients) {
      _scrollController.animateTo(
        _scrollController.position.maxScrollExtent,
        duration: const Duration(milliseconds: 300),
        curve: Curves.easeOut,
      );
    }
  }

  @override
  void dispose() {
    final chatProvider = Provider.of<ChatProvider>(context, listen: false);
    chatProvider.removeListener(_onProviderChange);
    // Explicitly disconnect WebSocket and clear state on exit
    chatProvider.closeConversation();
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  void _handleSend() {
    final text = _messageController.text.trim();
    if (text.isEmpty) return;

    final chatProvider = Provider.of<ChatProvider>(context, listen: false);
    chatProvider.sendChatMessage(
      text: text,
      companionId: widget.companionId,
    );
    _messageController.clear();
    
    // Smooth immediate scroll
    Future.delayed(const Duration(milliseconds: 50), _scrollToBottom);
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);
    final chatProvider = Provider.of<ChatProvider>(context);
    final primaryColor = theme.primaryColor;

    return Scaffold(
      backgroundColor: theme.scaffoldBackgroundColor,
      appBar: AppBar(
        backgroundColor: theme.cardColor,
        elevation: 0.5,
        leading: IconButton(
          icon: const Icon(Icons.arrow_back_ios_new_rounded, size: 20),
          onPressed: () => Navigator.pop(context),
        ),
        titleSpacing: 0,
        title: Row(
          children: [
            _buildAppBarAvatar(widget.companionAvatar, widget.companionName, theme),
            const SizedBox(width: 12),
            Expanded(
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    widget.companionName,
                    style: GoogleFonts.quicksand(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                      color: theme.textTheme.bodyLarge?.color,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Row(
                    children: [
                      Container(
                        height: 7,
                        width: 7,
                        decoration: const BoxDecoration(
                          color: Colors.green,
                          shape: BoxShape.circle,
                        ),
                      ),
                      const SizedBox(width: 5),
                      Text(
                        'Online',
                        style: GoogleFonts.quicksand(
                          fontSize: 11,
                          color: Colors.green,
                          fontWeight: FontWeight.bold,
                        ),
                      ),
                    ],
                  ),
                ],
              ),
            ),
          ],
        ),
        actions: [
          // Audio Call Button
          IconButton(
            icon: Icon(Icons.phone_rounded, color: primaryColor),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => CallScreen(
                    conversationId: widget.conversationId,
                    receiverId: widget.companionId,
                    receiverName: widget.companionName,
                    isCaller: true,
                  ),
                ),
              );
            },
          ),
          // Video Call Button
          IconButton(
            icon: Icon(Icons.videocam_rounded, color: primaryColor, size: 26),
            onPressed: () {
              Navigator.push(
                context,
                MaterialPageRoute(
                  builder: (context) => CallScreen(
                    conversationId: widget.conversationId,
                    receiverId: widget.companionId,
                    receiverName: widget.companionName,
                    isCaller: true,
                  ),
                ),
              );
            },
          ),
          const SizedBox(width: 8),
        ],
      ),
      body: SafeArea(
        child: Column(
          children: [
            // Messages Pane
            Expanded(
              child: chatProvider.isLoadingMessages
                  ? Center(
                      child: LoadingAnimationWidget.threeArchedCircle(
                        color: primaryColor,
                        size: 40,
                      ),
                    )
                  : chatProvider.activeMessages.isEmpty
                      ? _buildEmptyChatState(theme, primaryColor)
                      : ListView.builder(
                          controller: _scrollController,
                          padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 20.0),
                          itemCount: chatProvider.activeMessages.length,
                          itemBuilder: (context, index) {
                            final chat = chatProvider.activeMessages[index];
                            return _buildMessageBubble(chat, chatProvider.myProfileId, theme);
                          },
                        ),
            ),
            
            // Text Input Pane
            _buildInputBar(theme, primaryColor),
          ],
        ),
      ),
    );
  }

  Widget _buildAppBarAvatar(String avatarUrl, String displayName, ThemeData theme) {
    if (avatarUrl.isNotEmpty) {
      return CircleAvatar(
        radius: 20,
        backgroundImage: NetworkImage(avatarUrl),
      );
    }
    final initial = displayName.isNotEmpty ? displayName[0].toUpperCase() : 'P';
    return Container(
      height: 40,
      width: 40,
      decoration: BoxDecoration(
        shape: BoxShape.circle,
        gradient: LinearGradient(
          colors: [
            theme.primaryColor,
            theme.primaryColor.withBlue(220),
          ],
        ),
      ),
      child: Center(
        child: Text(
          initial,
          style: GoogleFonts.quicksand(
            fontSize: 16,
            fontWeight: FontWeight.bold,
            color: Colors.white,
          ),
        ),
      ),
    );
  }

  Widget _buildEmptyChatState(ThemeData theme, Color primaryColor) {
    return Center(
      child: SingleChildScrollView(
        child: Column(
          children: [
            Icon(Icons.handshake_outlined, size: 60, color: primaryColor.withOpacity(0.6)),
            const SizedBox(height: 16),
            Text(
              "Start Conversation",
              style: GoogleFonts.quicksand(
                fontSize: 18,
                fontWeight: FontWeight.bold,
                color: theme.textTheme.bodyLarge?.color,
              ),
            ),
            const SizedBox(height: 8),
            Padding(
              padding: const EdgeInsets.symmetric(horizontal: 40),
              child: Text(
                "Send a message to initialize your session consultation safely.",
                textAlign: TextAlign.center,
                style: GoogleFonts.quicksand(
                  fontSize: 13,
                  color: Colors.grey[500],
                ),
              ),
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildMessageBubble(dynamic chat, String? myProfileId, ThemeData theme) {
    final senderId = chat['senderId']?.toString();
    final text = chat['text']?.toString() ?? '';
    final isSystem = chat['isSystem'] as bool? ?? false;
    final isMe = senderId == myProfileId;

    if (isSystem || senderId == 'SYSTEM') {
      return Padding(
        padding: const EdgeInsets.symmetric(vertical: 12.0),
        child: Center(
          child: Container(
            padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 8.0),
            decoration: BoxDecoration(
              // ignore: deprecated_member_use
              color: Colors.yellow.shade100.withOpacity(0.3),
              borderRadius: BorderRadius.circular(12.0),
              border: Border.all(color: Colors.yellow.shade400.withOpacity(0.3)),
            ),
            child: Text(
              text,
              textAlign: TextAlign.center,
              style: GoogleFonts.quicksand(
                fontSize: 12,
                color: theme.brightness == Brightness.dark ? Colors.yellow[300] : Colors.yellow[900],
                fontWeight: FontWeight.w600,
              ),
            ),
          ),
        ),
      );
    }

    String timeText = '';
    final createdAtRaw = chat['createdAt']?.toString();
    if (createdAtRaw != null) {
      try {
        final dt = DateTime.parse(createdAtRaw).toLocal();
        timeText = DateFormat('HH:mm').format(dt);
      } catch (_) {}
    }

    return Padding(
      padding: const EdgeInsets.only(bottom: 12.0),
      child: Column(
        crossAxisAlignment: isMe ? CrossAxisAlignment.end : CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: isMe ? MainAxisAlignment.end : MainAxisAlignment.start,
            crossAxisAlignment: CrossAxisAlignment.end,
            children: [
              if (!isMe) ...[
                const SizedBox(width: 4),
              ],
              Flexible(
                child: Container(
                  padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 12.0),
                  decoration: BoxDecoration(
                    color: isMe 
                        ? theme.primaryColor 
                        : theme.cardColor,
                    borderRadius: BorderRadius.only(
                      topLeft: const Radius.circular(16.0),
                      topRight: const Radius.circular(16.0),
                      bottomLeft: isMe ? const Radius.circular(16.0) : const Radius.circular(2.0),
                      bottomRight: isMe ? const Radius.circular(2.0) : const Radius.circular(16.0),
                    ),
                    boxShadow: [
                      BoxShadow(
                        color: Colors.black.withOpacity(0.02),
                        blurRadius: 4,
                        offset: const Offset(0, 2),
                      )
                    ],
                  ),
                  child: Text(
                    text,
                    style: GoogleFonts.quicksand(
                      fontSize: 14.5,
                      fontWeight: FontWeight.w500,
                      color: isMe ? Colors.white : theme.textTheme.bodyLarge?.color,
                    ),
                  ),
                ),
              ),
            ],
          ),
          const SizedBox(height: 4),
          Padding(
            padding: const EdgeInsets.symmetric(horizontal: 8.0),
            child: Text(
              timeText,
              style: GoogleFonts.quicksand(
                fontSize: 10,
                color: Colors.grey,
                fontWeight: FontWeight.w500,
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildInputBar(ThemeData theme, Color primaryColor) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 10.0),
      decoration: BoxDecoration(
        color: theme.cardColor,
        border: Border(
          top: BorderSide(
            color: theme.brightness == Brightness.dark ? Colors.grey[800]! : Colors.grey[200]!,
            width: 0.5,
          ),
        ),
      ),
      child: Row(
        children: [
          // Emoji icon placeholder
          IconButton(
            icon: Icon(Icons.sentiment_satisfied_alt_rounded, color: Colors.grey[600]),
            onPressed: () {},
          ),
          // Attachment icon placeholder
          IconButton(
            icon: Icon(Icons.attach_file_rounded, color: Colors.grey[600]),
            onPressed: () {},
          ),
          
          Expanded(
            child: Container(
              decoration: BoxDecoration(
                color: theme.brightness == Brightness.dark ? Colors.grey[900] : Colors.grey[100],
                borderRadius: BorderRadius.circular(24.0),
              ),
              child: TextField(
                controller: _messageController,
                style: GoogleFonts.quicksand(fontWeight: FontWeight.w500, fontSize: 14.5),
                maxLines: null,
                keyboardType: TextInputType.multiline,
                decoration: InputDecoration(
                  hintText: 'Type your message...',
                  hintStyle: GoogleFonts.quicksand(color: Colors.grey, fontWeight: FontWeight.w400),
                  border: InputBorder.none,
                  contentPadding: const EdgeInsets.symmetric(horizontal: 16.0, vertical: 10.0),
                ),
              ),
            ),
          ),
          
          const SizedBox(width: 8),
          
          // Send Button
          InkWell(
            onTap: _handleSend,
            borderRadius: BorderRadius.circular(24.0),
            child: CircleAvatar(
              radius: 22,
              backgroundColor: primaryColor,
              child: const Icon(
                Icons.send_rounded,
                color: Colors.white,
                size: 20,
              ),
            ),
          ),
        ],
      ),
    );
  }
}
