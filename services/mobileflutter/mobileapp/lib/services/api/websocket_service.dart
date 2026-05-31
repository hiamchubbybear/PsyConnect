import 'dart:async';
import 'dart:convert';
import 'dart:io';

class WebSocketService {
  static final WebSocketService _instance = WebSocketService._internal();
  factory WebSocketService() => _instance;
  WebSocketService._internal();

  WebSocket? _socket;
  bool _isConnected = false;
  bool get isConnected => _isConnected;

  final _messageController = StreamController<Map<String, dynamic>>.broadcast();
  Stream<Map<String, dynamic>> get messageStream => _messageController.stream;

  final _connectionController = StreamController<bool>.broadcast();
  Stream<bool> get connectionStream => _connectionController.stream;

  Timer? _reconnectTimer;
  String? _currentUrl;

  Future<void> connect(String url) async {
    if (_isConnected && _currentUrl == url) return;
    
    _currentUrl = url;
    _reconnectTimer?.cancel();
    
    try {
      print("Connecting to WebSocket: $url");
      _socket = await WebSocket.connect(url).timeout(const Duration(seconds: 10));
      _isConnected = true;
      _connectionController.add(true);
      print("WebSocket connected successfully");

      _socket!.listen(
        (data) {
          try {
            final Map<String, dynamic> decoded = jsonDecode(data as String);
            _messageController.add(decoded);
          } catch (e) {
            print("Error parsing WebSocket message: $e");
          }
        },
        onError: (err) {
          print("WebSocket error: $err");
          _handleDisconnect();
        },
        onDone: () {
          print("WebSocket connection closed by server");
          _handleDisconnect();
        },
        cancelOnError: true,
      );
    } catch (e) {
      print("Failed to connect to WebSocket: $e");
      _handleDisconnect();
    }
  }

  void _handleDisconnect() {
    _isConnected = false;
    _connectionController.add(false);
    _socket = null;
    
    // Auto-reconnect after 5 seconds if a URL was provided
    if (_currentUrl != null) {
      _reconnectTimer?.cancel();
      _reconnectTimer = Timer(const Duration(seconds: 5), () {
        print("Attempting to reconnect to WebSocket...");
        if (_currentUrl != null) {
          connect(_currentUrl!);
        }
      });
    }
  }

  void sendMessage(Map<String, dynamic> message) {
    if (_socket != null && _isConnected) {
      final String payload = jsonEncode(message);
      _socket!.add(payload);
      print("Sent WebSocket message: $payload");
    } else {
      print("Cannot send message. WebSocket is not connected.");
    }
  }

  void disconnect() {
    _currentUrl = null;
    _reconnectTimer?.cancel();
    _socket?.close();
    _isConnected = false;
    _connectionController.add(false);
    _socket = null;
    print("WebSocket disconnected manually");
  }
}
