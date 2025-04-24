import 'package:flutter/material.dart';
import 'package:provider/src/provider.dart';

class AnimatedApp extends StatefulWidget {
  final Widget child;

  const AnimatedApp(MultiProvider multiProvider, {super.key, required this.child});

  @override
  State<AnimatedApp> createState() => _AnimatedAppState();
}

class _AnimatedAppState extends State<AnimatedApp>
    with SingleTickerProviderStateMixin {
  late AnimationController _controller;
  late Animation<double> _animation;

  @override
  void initState() {
    super.initState();
    _controller =
        AnimationController(vsync: this, duration: const Duration(seconds: 1))
          ..forward();

    _animation = CurvedAnimation(parent: _controller, curve: Curves.easeIn);
  }

  @override
  Widget build(BuildContext context) {
    return FadeTransition(
      opacity: _animation,
      child: widget.child,
    );
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }
}
