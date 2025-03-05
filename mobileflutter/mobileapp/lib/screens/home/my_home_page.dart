import 'package:PsyConnect/page/home_page_scroll_view.dart';
import 'package:PsyConnect/screens/component/setting.dart';
import 'package:flutter/material.dart';

class MyHomePage extends StatefulWidget {
  const MyHomePage({super.key, required this.title});

  final String title;

  @override
  State<MyHomePage> createState() => _MyHomePageState();
}

class _MyHomePageState extends State<MyHomePage>
    with SingleTickerProviderStateMixin {
  late TabController _tabController;

  @override
  void initState() {
    super.initState();
    _tabController = TabController(length: 5, vsync: this);
  }

  @override
  void dispose() {
    _tabController.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        automaticallyImplyLeading: false,
      ),
      body: Padding(
        padding: const EdgeInsets.all(16.0),
        child: TabBarView(
          controller: _tabController,
          children: const [
            HomePageScrollView(),
            Center(child: Text("Schedule Page")),
            Center(child: Text("Chat Page")),
            Center(child: Text("Notifications Page")),
            Setting(),
          ],
        ),
      ),
      bottomNavigationBar: Material(
        child: Padding(
          padding: const EdgeInsets.symmetric(vertical: 16.0),
          child: TabBar(
            controller: _tabController,
            tabs: const [
              Tab(icon: Icon(Icons.home, size: 37)),
              Tab(
                  icon: ImageIcon(
                      AssetImage("assets/images/dark_mode_schedule.png"),
                      size: 28)),
              Tab(
                  icon: ImageIcon(AssetImage("assets/images/dark_mode_chat.png"),
                      size: 28)),
              Tab(
                  icon: ImageIcon(
                      AssetImage("assets/images/dark_mode_notification.png"),
                      size: 28)),
              Tab(icon: Icon(Icons.settings, size: 36)),
            ],
          ),
        ),
      ),
    );
  }
}
