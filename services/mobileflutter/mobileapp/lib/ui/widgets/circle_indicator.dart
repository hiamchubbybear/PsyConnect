import 'dart:math';
import 'package:flutter/material.dart';

class CircleIndicator extends CustomPainter {
  final Color mColor;
  final double thickness;
  final double progress;

  const CircleIndicator({
    this.progress = 150.0,
    required this.mColor,
    required this.thickness,
  });

  double deg2rad(double deg) => deg * pi / 180;

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width / 2;

    final mPaint = Paint()
      ..color = Colors.grey.withOpacity(.23)
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    final mPaintColor = Paint()
      ..color = mColor
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(1),
      deg2rad(358),
      false,
      mPaint,
    );

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(270),
      deg2rad(progress - 2),
      false,
      mPaintColor,
    );
  }

  @override
  bool shouldRepaint(covariant CircleIndicator oldDelegate) {
    return oldDelegate.progress != progress ||
        oldDelegate.mColor != mColor ||
        oldDelegate.thickness != thickness;
  }
}
