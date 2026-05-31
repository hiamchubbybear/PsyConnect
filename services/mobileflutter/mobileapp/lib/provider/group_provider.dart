import 'package:flutter/material.dart';
import 'package:PsyConnect/models/group.dart';
import 'package:PsyConnect/services/profile_service/group_service.dart';

class GroupProvider with ChangeNotifier {
  final GroupService _groupService = GroupService();

  List<Group> _groups = [];
  List<Group> get groups => _groups;

  bool _isLoading = false;
  bool get isLoading => _isLoading;

  String? _errorMessage;
  String? get errorMessage => _errorMessage;

  String _selectedCategory = "";
  String get selectedCategory => _selectedCategory;

  String _searchQuery = "";
  String get searchQuery => _searchQuery;

  // Track which groups the user has joined during this session locally
  final Set<String> _joinedGroupIds = {};
  Set<String> get joinedGroupIds => _joinedGroupIds;

  void setCategory(String category) {
    _selectedCategory = category;
    fetchGroups();
  }

  void setSearchQuery(String query) {
    _searchQuery = query;
    fetchGroups();
  }

  Future<void> fetchGroups() async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final categoryFilter = _selectedCategory.isEmpty ? null : _selectedCategory;
      final qFilter = _searchQuery.isEmpty ? null : _searchQuery;

      final fetched = await _groupService.getGroups(
        category: categoryFilter,
        q: qFilter,
      );
      _groups = fetched;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  Future<bool> joinGroup(String groupId) async {
    try {
      await _groupService.joinGroup(groupId);
      _joinedGroupIds.add(groupId);
      
      // Update local member count for visual completeness
      final index = _groups.indexWhere((g) => g.id == groupId);
      if (index != -1) {
        final currentGroup = _groups[index];
        _groups[index] = Group(
          id: currentGroup.id,
          name: currentGroup.name,
          description: currentGroup.description,
          category: currentGroup.category,
          memberCount: currentGroup.memberCount + 1,
          isPublic: currentGroup.isPublic,
          createdAt: currentGroup.createdAt,
        );
      }
      
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    }
  }

  Future<bool> createGroup({
    required String name,
    required String description,
    required String category,
    bool isPublic = true,
  }) async {
    _isLoading = true;
    _errorMessage = null;
    notifyListeners();

    try {
      final payload = {
        "name": name,
        "description": description,
        "category": category,
        "isPublic": isPublic,
      };
      final newGroup = await _groupService.createGroup(payload);
      _groups.insert(0, newGroup);
      _joinedGroupIds.add(newGroup.id);
      
      notifyListeners();
      return true;
    } catch (e) {
      _errorMessage = e.toString().replaceAll("Exception: ", "");
      return false;
    } finally {
      _isLoading = false;
      notifyListeners();
    }
  }

  void clearError() {
    _errorMessage = null;
    notifyListeners();
  }
}
