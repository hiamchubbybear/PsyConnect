import 'package:PsyConnect/core/colors/color.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/provider/theme_provider.dart';
import 'package:PsyConnect/provider/chat_provider.dart';
import 'package:PsyConnect/ui/screens/home_page_scroll_view.dart';
import 'package:PsyConnect/ui/screens/profile_page.dart';
import 'package:PsyConnect/ui/screens/schedule_home_page.dart';
import 'package:PsyConnect/ui/screens/chat_list_screen.dart';
import 'package:PsyConnect/ui/screens/call_screen.dart';
import 'package:PsyConnect/ui/screens/create_post_screen.dart';
import 'package:PsyConnect/services/api/websocket_service.dart';
import 'package:flutter/material.dart';
import 'package:provider/provider.dart';
import 'package:google_fonts/google_fonts.dart';

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;
  int _currentIndex = 0;
  bool _isShowingCallDialog = false;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 5, vsync: this);
    _tabController.addListener(_handleTabSelection);
  }

  void _handleTabSelection() {
    if (_tabController.indexIsChanging) {
      setState(() {
        _currentIndex = _tabController.index;
      });
    }
  }

  @override
  void dispose() {
    _tabController.removeListener(_handleTabSelection);
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final chatProvider = Provider.of<ChatProvider>(context);
    if (chatProvider.activeIncomingCall != null && !_isShowingCallDialog) {
      _isShowingCallDialog = true;
      final callInfo = chatProvider.activeIncomingCall!;
      WidgetsBinding.instance.addPostFrameCallback((_) {
        _showIncomingCallDialog(context, callInfo, chatProvider);
      });
    }

    final isDark = Provider.of<ThemeProvider>(context).isDarkMode;
    Size size = MediaQuery.of(context).size;
    return Scaffold(
      extendBody: true,
      body: TabBarView(
        controller: _tabController,
        physics: const NeverScrollableScrollPhysics(),
        children: [
          const HomePageScrollView(),
          ScheduleHomePage(size: size),
          const ProfilePage(),
          const ChatListScreen(),
          const ProfilePage()
        ],
      ),
      floatingActionButton: FloatingActionButton(
        onPressed: () {},
        backgroundColor: isDark ? secondaryColor : primaryColor,
        elevation: 2,
        child: Icon(Icons.add,
            color: isDark ? blackColor : Colors.white, size: 28),
      ),
      floatingActionButtonLocation: FloatingActionButtonLocation.centerDocked,
      bottomNavigationBar: Container(
        height: 100,
        decoration: BoxDecoration(
          borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
          boxShadow: [
            BoxShadow(
              // ignore: deprecated_member_use
              color: Colors.black.withOpacity(0.1),
              blurRadius: 10,
              spreadRadius: 2,
            )
          ],
        ),
        child: ClipRRect(
          borderRadius: const BorderRadius.vertical(top: Radius.circular(20)),
          child: BottomAppBar(
            shape: const CircularNotchedRectangle(),
            notchMargin: 8,
            child: Container(
              padding: const EdgeInsets.symmetric(horizontal: 12),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceAround,
                children: [
                  _buildNavItem(
                    icon: Icons.home_outlined,
                    activeIcon: Icons.home,
                    label: 'Home',
                    index: 0,
                    isDark: isDark
                  ),
                  _buildNavItem(
                    icon: Icons.calendar_today_outlined,
                    activeIcon: Icons.calendar_today,
                    label: 'Schedule',
                    index: 1,
                    isDark: isDark
                  ),
                  const SizedBox(width: 40),
                  _buildNavItem(
                    icon: Icons.notifications_outlined,
                    activeIcon: Icons.chat,
                    label: 'Chat',
                    index: 3,
                    isDark: isDark
                  ),
                  _buildNavItem(
                    icon: Icons.person_outline,
                    activeIcon: Icons.person,
                    label: 'Profile',
                    index: 4,
                    isDark: isDark
                  ),
                ],
              ),
            ),
          ),
        ),
      ),
    );
  }

  Widget _buildNavItem({
    required IconData icon,
    required IconData activeIcon,
    required String label,
    required int index,
    required bool isDark,
  }) {
    final isActive = _currentIndex == index;
    final color = isActive ? Colors.grey : Colors.grey;

    return InkWell(
      onTap: () {
        setState(() {
          _currentIndex = index;
          _tabController.animateTo(index);
        });
      },
      splashColor: Colors.transparent,
      highlightColor: Colors.transparent,
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          Icon(
            isActive ? activeIcon : icon,
            color: color,
            size: 23,
          ),
          const SizedBox(height: 2),
          Text(label,
              style: quickSand12Font.copyWith(
                color: isDark ? Colors.white : Colors.black,
              )),
        ],
      ),
    );
  }

  void _showIncomingCallDialog(BuildContext context, Map<String, dynamic> callInfo, ChatProvider chatProvider) {
    showDialog(
      context: context,
      barrierDismissible: false,
      builder: (dialogCtx) {
        final theme = Theme.of(context);
        return AlertDialog(
          shape: RoundedRectangleBorder(borderRadius: BorderRadius.circular(20)),
          backgroundColor: theme.cardColor,
          title: Center(
            child: Text(
              "Incoming Call",
              style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 20),
            ),
          ),
          content: Column(
            mainAxisSize: MainAxisSize.min,
            children: [
              Container(
                height: 80,
                width: 80,
                decoration: BoxDecoration(
                  color: theme.primaryColor.withOpacity(0.1),
                  shape: BoxShape.circle,
                ),
                child: Icon(Icons.person, size: 40, color: theme.primaryColor),
              ),
              const SizedBox(height: 16),
              Text(
                callInfo["callerName"] ?? "Specialist",
                style: GoogleFonts.quicksand(fontWeight: FontWeight.bold, fontSize: 18),
              ),
              const SizedBox(height: 8),
              Text(
                "Inbound Consultation Video Call Room",
                style: GoogleFonts.quicksand(color: Colors.grey, fontSize: 13),
              ),
            ],
          ),
          actionsAlignment: MainAxisAlignment.spaceEvenly,
          actions: [
            // Decline Button
            IconButton(
              icon: const Icon(Icons.call_end_rounded, color: Colors.red, size: 36),
              onPressed: () {
                Navigator.pop(dialogCtx);
                setState(() {
                  _isShowingCallDialog = false;
                });
                WebSocketService().sendMessage({
                  "type": "leave",
                  "conversationId": callInfo["conversationId"],
                  "receiverId": callInfo["callerId"],
                  "data": {
                    "sessionId": callInfo["sessionId"] ?? ""
                  }
                });
                chatProvider.clearIncomingCall();
              },
            ),
            // Accept Button
            IconButton(
              icon: const Icon(Icons.call_rounded, color: Colors.green, size: 36),
              onPressed: () {
                Navigator.pop(dialogCtx);
                setState(() {
                  _isShowingCallDialog = false;
                });
                chatProvider.clearIncomingCall();
                
                Navigator.push(
                  context,
                  MaterialPageRoute(
                    builder: (context) => CallScreen(
                      conversationId: callInfo["conversationId"],
                      receiverId: callInfo["callerId"],
                      receiverName: callInfo["callerName"],
                      isCaller: false,
                      initialOfferSdp: callInfo["sdp"],
                      initialSessionId: callInfo["sessionId"],
                    ),
                  ),
                );
              },
            ),
          ],
        );
      },
    );
  }
}
