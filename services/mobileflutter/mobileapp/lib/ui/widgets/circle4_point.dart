import 'dart:math';
import 'package:flutter/material.dart';

class Circle4Point extends CustomPainter {
  final double thickness;
  final double gap;

  final Color bottomRightColor;
  final Color bottomLeftColor;
  final Color topRightColor;
  final Color topLeftColor;

  Circle4Point({
    required this.gap,
    required this.thickness,
    required this.bottomLeftColor,
    required this.bottomRightColor,
    required this.topLeftColor,
    required this.topRightColor,
  });

  double deg2rad(double deg) => deg * pi / 180;

  @override
  void paint(Canvas canvas, Size size) {
    final center = Offset(size.width / 2, size.height / 2);
    final radius = size.width / 2;

    Paint paintTr = Paint()
      ..color = topRightColor
      ..strokeCap = StrokeCap.round
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    Paint paintBr = Paint()
      ..color = bottomRightColor
      ..strokeCap = StrokeCap.round
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    Paint paintBl = Paint()
      ..color = bottomLeftColor
      ..strokeCap = StrokeCap.round
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    Paint paintTl = Paint()
      ..color = topLeftColor
      ..strokeCap = StrokeCap.round
      ..style = PaintingStyle.stroke
      ..strokeWidth = thickness;

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(320 + gap),
      deg2rad(50 - (gap * 2)),
      false,
      paintBr,
    );

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(24 + gap),
      deg2rad(54 - (gap * 2)),
      false,
      paintBl,
    );

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(94 + gap),
      deg2rad(144 - (gap * 2)),
      false,
      paintTl,
    );

    canvas.drawArc(
      Rect.fromCircle(center: center, radius: radius),
      deg2rad(256 + gap),
      deg2rad(50 - (gap * 2)),
      false,
      paintTr,
    );
  }

  @override
  bool shouldRepaint(covariant Circle4Point oldDelegate) {
    return oldDelegate.gap != gap ||
        oldDelegate.thickness != thickness ||
        oldDelegate.bottomLeftColor != bottomLeftColor ||
        oldDelegate.bottomRightColor != bottomRightColor ||
        oldDelegate.topLeftColor != topLeftColor ||
        oldDelegate.topRightColor != topRightColor;
  }
}
