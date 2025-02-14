import 'package:flutter/material.dart';
import 'package:fluttertoast/fluttertoast.dart';

class ToastService {
  static void showSuccessToast({required String message, int duration = 2}) {
    _showToast(
      message: message,
      backgroundColor: Colors.green,
      toastLength: duration,
      textColor: Colors.white,
      fontSize: 18.0,
      fontWeight: FontWeight.bold,
    );
  }

  static void showErrorToast({required String message, int duration = 3}) {
    _showToast(
      message: message,
      backgroundColor: Colors.red,
      toastLength: duration,
      textColor: Colors.white,
      fontSize: 18.0,
      fontWeight: FontWeight.bold,
    );
  }

  static void _showToast({
    required String message,
    required Color backgroundColor,
    required int toastLength,
    required Color textColor,
    required double fontSize,
    required FontWeight fontWeight,
  }) {
    Fluttertoast.showToast(
      msg: message,
      toastLength: toastLength == 2 ? Toast.LENGTH_SHORT : Toast.LENGTH_LONG,
      gravity: ToastGravity.CENTER,
      backgroundColor: backgroundColor,
      textColor: textColor,
      fontSize: fontSize,
      timeInSecForIosWeb: 1,
    );



    final TextStyle textStyle = TextStyle(
      fontSize: fontSize,
      fontWeight: fontWeight,
      color: textColor,
    );
  }
}
