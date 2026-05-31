import 'package:flutter/material.dart';
import 'package:PsyConnect/models/swipe_card.dart';
import 'package:PsyConnect/services/profile_service/swipe_service.dart';

class SwipeProvider with ChangeNotifier {
  final SwipeService _swipeService = SwipeService();

  List<SwipeCard> _cards = [];
  List<SwipeCard> get cards => _cards;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  bool _missingProfile = false;
  bool get missingProfile => _missingProfile;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  Future<void> fetchRecommendations() async {
    _isLoading = true;
    _missingProfile = false;
    _errorMessage = null;
    notifyListeners();

    try {
      // 1. Initial recommendation fetch
      List<SwipeCard> fetched = await _swipeService.getSwipeRecommendations();
      
      // 2. Fallback Generation Trigger if recommendation deck is empty
      if (fetched.isEmpty) {
        await _swipeService.triggerRecommendationGeneration();
        // Retry fetch
        fetched = await _swipeService.getSwipeRecommendations();
      }
      
      _cards = fetched;
    } catch (e) {
      final msg = e.toString().toLowerCase();
      // Detect if user has no consultation client profile set up yet
      if (msg.contains("503") || msg.contains("404") || msg.contains("500") || msg.contains("profile")) {
        _missingProfile = true;
      }
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      _cards = [];
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<void> swipeLeft(String therapistId) async {
    // Passed
    try {
      await _swipeService.recordSwipe(therapistId, 'passed');
    } catch (e) {
      print("Swipe pass persist error: $e");
    }
    
    // Remove locally
    _cards.removeWhere((c) => c.profileId == therapistId);
    notifyListeners();
  }

  Future<void> swipeRight(String therapistId) async {
    // Swiped (Interested)
    try {
      await _swipeService.recordSwipe(therapistId, 'swiped');
    } catch (e) {
      print("Swipe like persist error: $e");
    }

    // Remove locally
    _cards.removeWhere((c) => c.profileId == therapistId);
    notifyListeners();
  }
}
