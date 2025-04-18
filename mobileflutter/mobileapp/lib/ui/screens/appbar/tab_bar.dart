import 'package:flutter/material.dart';

class CustomeTabBar extends StatelessWidget implements PreferredSizeWidget {
  const CustomeTabBar({super.key});

  @override
  Widget build(BuildContext context) {
    return AppBar(
      automaticallyImplyLeading: false,
      bottom: const TabBar(
        tabs: [
          Tab(
            icon: Icon(
              Icons.home,
              size: 36,
            ),
          ),
          Tab(
              icon: ImageIcon(
            AssetImage("assets/images/dark_mode_schedule.png"),
            size: 28,
          )),
          Tab(
              icon: ImageIcon(
            AssetImage("assets/images/dark_mode_chat.png"),
            size: 28,
          )),
          Tab(
              icon: ImageIcon(
            AssetImage("assets/images/dark_mode_notification.png"),
            size: 28,
          )),
          Tab(icon: Icon(Icons.settings, size: 36)),
        ],
      ),
    );
  }

  @override
  Size get preferredSize => const Size.fromHeight(kToolbarHeight + 48);
}
