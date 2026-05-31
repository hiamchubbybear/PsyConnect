import 'dart:async';
import 'dart:convert';
import 'package:flutter/foundation.dart';
import 'package:http/http.dart' as http;
import 'package:PsyConnect/core/preferences/sharepreference_provider.dart';
import 'package:PsyConnect/core/variable/variable.dart';
import 'package:PsyConnect/models/user_profile.dart';
import 'package:PsyConnect/services/api/api_service.dart';
import 'package:PsyConnect/services/api/websocket_service.dart';
import 'package:PsyConnect/services/profile_service/profile.dart';

class ChatProvider with ChangeNotifier {
  final ProfileService _profileService = ProfileService();
  final WebSocketService _wsService = WebSocketService();
  StreamSubscription<Map<String, dynamic>>? _wsMessageSubscription;

  List<dynamic> _conversations = [];
  List<dynamic> get conversations => _conversations;

  bool _isLoadingConversations = false;
  bool get isLoadingConversations => _isLoadingConversations;

  List<dynamic> _activeMessages = [];
  List<dynamic> get activeMessages => _activeMessages;

  bool _isLoadingMessages = false;
  bool get isLoadingMessages => _isLoadingMessages;

  final Map<String, UserProfile> _cachedProfiles = {};
  Map<String, UserProfile> get cachedProfiles => _cachedProfiles;

  String? _activeConversationId;
  String? get activeConversationId => _activeConversationId;

  String? _myProfileId;
  String? get myProfileId => _myProfileId;

  UserProfile? _myProfile;
  UserProfile? get myProfile => _myProfile;

  Map<String, dynamic>? _activeIncomingCall;
  Map<String, dynamic>? get activeIncomingCall => _activeIncomingCall;

  ChatProvider() {
    // Listen to incoming WebSocket messages globally
    _wsMessageSubscription = _wsService.messageStream.listen(_handleIncomingWSMessage);
  }

  @override
  void dispose() {
    _wsMessageSubscription?.cancel();
    super.dispose();
  }

  Future<void> initMyProfile() async {
    if (_myProfileId != null) return;
    try {
      _myProfile = await _profileService.getUserProfile();
      _myProfileId = _myProfile?.profileId;
      notifyListeners();
    } catch (e) {
      print("Error loading own profile: $e");
    }
  }

  Future<void> fetchRecentConversations() async {
    await initMyProfile();
    if (_myProfileId == null) return;

    _isLoadingConversations = true;
    notifyListeners();

    try {
      final token = await SharedPreferencesProvider().getJwt();
      if (token == null) throw Exception("No authorization token found");

      // chatservice maps /conversations/me using Gateway routing
      final url = Uri.parse("$baseUrl/conversations/me?userId=$_myProfileId");
      final response = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        _conversations = decoded['data'] ?? [];
        
        // Asynchronously load participant profiles to display details nicely
        for (var conv in _conversations) {
          final participants = conv['participants'] as List<dynamic>? ?? [];
          for (var pId in participants) {
            if (pId != _myProfileId) {
              fetchUserProfileSilently(pId.toString());
            }
          }
        }
      } else {
        print("Failed to fetch conversations: ${response.statusCode} - ${response.body}");
      }
    } catch (e) {
      print("Exception in fetchRecentConversations: $e");
    } finally {
      _isLoadingConversations = false;
      notifyListeners();
    }
  }

  Future<void> fetchConversationMessages(String conversationId) async {
    _isLoadingMessages = true;
    _activeMessages = [];
    notifyListeners();

    try {
      final token = await SharedPreferencesProvider().getJwt();
      if (token == null) throw Exception("No authorization token");

      // chatservice /chats/conversation/:id
      final url = Uri.parse("$baseUrl/chats/conversation/$conversationId?limit=50");
      final response = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        final List<dynamic> fetchedChats = decoded['data'] ?? [];
        
        // Reverse so that newer messages render at the bottom in the chat UI
        _activeMessages = fetchedChats.reversed.toList();
      } else {
        print("Failed to fetch chats: ${response.statusCode} - ${response.body}");
      }
    } catch (e) {
      print("Exception fetching chats: $e");
    } finally {
      _isLoadingMessages = false;
      notifyListeners();
    }
  }

  Future<UserProfile?> fetchUserProfileSilently(String profileId) async {
    if (_cachedProfiles.containsKey(profileId)) {
      return _cachedProfiles[profileId];
    }

    try {
      final token = await SharedPreferencesProvider().getJwt();
      if (token == null) return null;

      final url = Uri.parse("$baseUrl/profile/$profileId");
      final response = await http.get(
        url,
        headers: {
          'Content-Type': 'application/json',
          'Authorization': 'Bearer $token',
        },
      );

      if (response.statusCode == 200) {
        final decoded = jsonDecode(response.body);
        final data = decoded['data'];
        if (data != null) {
          final profile = UserProfile.fromJson(data);
          _cachedProfiles[profileId] = profile;
          notifyListeners();
          return profile;
        }
      }
    } catch (e) {
      print("Error fetching profile $profileId: $e");
    }
    return null;
  }

  Future<void> openConversation({
    required String conversationId,
    required String companionId,
  }) async {
    await initMyProfile();
    _activeConversationId = conversationId;

    final token = await SharedPreferencesProvider().getJwt();
    if (token == null) return;

    // Convert Gateway baseUrl http://... to ws://...
    String wsBaseUrl = baseUrl.replaceFirst("http://", "ws://").replaceFirst("https://", "wss://");
    final wsUrl = "$wsBaseUrl/ws?token=$token&conversationId=$conversationId&receiver=$companionId";

    // Connect WebSocket
    await _wsService.connect(wsUrl);

    // Load messages
    await fetchConversationMessages(conversationId);
  }

  void closeConversation() {
    _activeConversationId = null;
    _activeMessages = [];
    _wsService.disconnect();
  }

  void sendChatMessage({
    required String text,
    required String companionId,
  }) {
    if (_activeConversationId == null || _myProfileId == null) return;

    final senderName = _myProfile != null 
        ? "${_myProfile!.firstName ?? ''} ${_myProfile!.lastName ?? ''}".trim() 
        : "User";

    final payload = {
      "type": "chat",
      "conversationId": _activeConversationId,
      "receiverId": companionId,
      "data": {
        "text": text,
        "receiverId": companionId,
        "senderName": senderName.isEmpty ? "User" : senderName
      }
    };

    _wsService.sendMessage(payload);
  }

  void _handleIncomingWSMessage(Map<String, dynamic> message) {
    final type = message['type']?.toString();
    final conversationId = message['conversationId']?.toString();

    if (type == 'chat' && conversationId == _activeConversationId) {
      final rawData = message['data'];
      
      // Parse broadcast response data
      // Broadcast data contains: id, senderId, content, timestamp, conversationId
      if (rawData != null) {
        final parsedChat = {
          "id": rawData["id"],
          "senderId": rawData["senderId"],
          "conversationId": rawData["conversationId"],
          "text": rawData["content"],
          "createdAt": rawData["timestamp"],
          "isSystem": false
        };

        // Append to list if not already present
        final exists = _activeMessages.any((msg) => msg['id'] == parsedChat['id']);
        if (!exists) {
          _activeMessages.add(parsedChat);
          notifyListeners();
        }
      }
    } else if (type == 'offer') {
      final rawData = message['data'];
      if (rawData != null) {
        final sdp = rawData['sdp']?.toString();
        final sessionId = rawData['sessionId']?.toString();
        final callerId = message['senderId']?.toString();
        if (sdp != null && callerId != null) {
          _handleIncomingCall(
            callerId: callerId,
            sdp: sdp,
            sessionId: sessionId,
            conversationId: conversationId ?? "",
          );
        }
      }
    } else if (type == 'leave') {
      clearIncomingCall();
    }

    // Refresh recent conversation list to update last messages & positions
    fetchRecentConversations();
  }

  void _handleIncomingCall({
    required String callerId,
    required String sdp,
    required String? sessionId,
    required String conversationId,
  }) async {
    final callerProfile = await fetchUserProfileSilently(callerId);
    final callerName = callerProfile != null 
        ? "${callerProfile.firstName ?? ''} ${callerProfile.lastName ?? ''}".trim()
        : "Specialist";
        
    _activeIncomingCall = {
      "callerId": callerId,
      "callerName": callerName.isEmpty ? "Specialist" : callerName,
      "sdp": sdp,
      "sessionId": sessionId,
      "conversationId": conversationId,
    };
    notifyListeners();
  }

  void clearIncomingCall() {
    _activeIncomingCall = null;
    notifyListeners();
  }
}
