import 'package:flutter/material.dart';
import 'package:PsyConnect/models/friend.dart';
import 'package:PsyConnect/services/profile_service/friend_service.dart';

class FriendProvider with ChangeNotifier {
  final FriendService _friendService = FriendService();

  List<Friend> _friends = [];
  List<Friend> get friends => _friends;

  List<Friend> _requests = [];
  List<Friend> get requests => _requests;

  List<Friend> _suggestions = [];
  List<Friend> get suggestions => _suggestions;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  Future<void> fetchAllSocialData() async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final myFriends = await _friendService.getMyFriends();
      final myRequests = await _friendService.getReceivedRequests();
      final mySuggestions = await _friendService.getFriendSuggestions();

      _friends = myFriends;
      _requests = myRequests;
      _suggestions = mySuggestions;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> sendFriendRequest(String targetId) async {
    try {
      await _friendService.sendFriendRequest(targetId);
      // Remove from suggestions locally
      _suggestions.removeWhere((f) => f.profileId == targetId);
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    }
  }

  Future<bool> acceptRequest(String targetId) async {
    try {
      await _friendService.acceptFriendRequest(targetId);
      
      // Move from requests to friends locally
      final reqIndex = _requests.indexWhere((f) => f.profileId == targetId);
      if (reqIndex != -1) {
        final accepted = _requests.removeAt(reqIndex);
        _friends.add(accepted);
      }
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    }
  }

  Future<bool> cancelRequest(String targetId) async {
    // Decline just removes locally for mock/simplified flow
    _requests.removeWhere((f) => f.profileId == targetId);
    notifyListeners();
    return true;
  }

  Future<bool> removeFriend(String targetId) async {
    try {
      await _friendService.unfriend(targetId);
      _friends.removeWhere((f) => f.profileId == targetId);
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
