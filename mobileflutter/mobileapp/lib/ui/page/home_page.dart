import 'package:flutter/material.dart';

class HomePageScrollView extends StatefulWidget {
  const HomePageScrollView({super.key});

  @override
  State<HomePageScrollView> createState() => _HomePageScrollViewState();
}

class _HomePageScrollViewState extends State<HomePageScrollView> {
  @override
  Widget build(BuildContext context) {
    return Stack(
      children: [
        CustomScrollView(
          slivers: [
            const SliverAppBar(
              pinned: true,
              expandedHeight: 200.0,
              flexibleSpace: FlexibleSpaceBar(
                centerTitle: true,
                title: Text('SliverAppBar Example'),
                background: FlutterLogo(),
              ),
            ),
            SliverList(
              delegate: SliverChildBuilderDelegate(
                (BuildContext context, int index) {
                  return ListTile(
                    leading: const Icon(Icons.star),
                    title: Text('Item #$index'),
                    subtitle: const Text('Subtitle here'),
                  );
                },
                childCount: 50,
              ),
            ),
          ],
        ),
      ],
    );
  }
}
