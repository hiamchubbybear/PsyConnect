import 'package:PsyConnect/core/colors/color.dart';
import 'package:flutter/material.dart';
import 'package:loading_animation_widget/loading_animation_widget.dart';

loadingBar() {
  return Center(
    child: LoadingAnimationWidget.fallingDot(color: primaryColor, size: 200),
  );
}
