import 'package:flutter/material.dart';

class HomeBar extends StatelessWidget {
  const HomeBar({super.key});

  @override
  Widget build(BuildContext context) {
    return const SliverAppBar(
      pinned: true,
      toolbarHeight: 50,
      title: Column(
        children: [ SizedBox(height: 10)],
      ),
    );
  }
}
