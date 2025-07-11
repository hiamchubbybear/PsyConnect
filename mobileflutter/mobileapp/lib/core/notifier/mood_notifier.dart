import 'package:flutter/material.dart';

class MoodNotifier extends ChangeNotifier {
  void notifyMoodChanged() {
    notifyListeners();
  }
}
