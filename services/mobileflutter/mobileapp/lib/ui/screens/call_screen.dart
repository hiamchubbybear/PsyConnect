import 'dart:async';
import 'package:flutter/material.dart';
import 'package:flutter_webrtc/flutter_webrtc.dart';
import 'package:google_fonts/google_fonts.dart';
import 'package:PsyConnect/services/api/websocket_service.dart';
import 'package:PsyConnect/core/toasting&loading/toast.dart';

class CallScreen extends StatefulWidget {
  final String conversationId;
  final String receiverId;
  final String receiverName;
  final bool isCaller;
  final String? initialOfferSdp; // Passed to callee
  final String? initialSessionId; // Passed to callee

  const CallScreen({
    super.key,
    required this.conversationId,
    required this.receiverId,
    required this.receiverName,
    required this.isCaller,
    this.initialOfferSdp,
    this.initialSessionId,
  });

  @override
  State<CallScreen> createState() => _CallScreenState();
}

class _CallScreenState extends State<CallScreen> {
  final _wsService = WebSocketService();
  StreamSubscription<Map<String, dynamic>>? _wsSubscription;

  final _localRenderer = RTCVideoRenderer();
  final _remoteRenderer = RTCVideoRenderer();
  
  RTCPeerConnection? _peerConnection;
  MediaStream? _localStream;
  
  bool _isAudioMuted = false;
  bool _isVideoOff = false;
  bool _isConnected = false;
  String? _sessionId;

  @override
  void initState() {
    super.initState();
    _sessionId = widget.initialSessionId;
    _initRenderers();
  }

  Future<void> _initRenderers() async {
    await _localRenderer.initialize();
    await _remoteRenderer.initialize();
    await _startConsultationCall();
  }

  Future<void> _startConsultationCall() async {
    // 1. Request hardware capture permissions & get media stream
    try {
      final Map<String, dynamic> mediaConstraints = {
        'audio': true,
        'video': {
          'facingMode': 'user',
        }
      };

      _localStream = await navigator.mediaDevices.getUserMedia(mediaConstraints);
      if (mounted) {
        setState(() {
          _localRenderer.srcObject = _localStream;
        });
      }

      // 2. Create peer connection with Google public STUN server
      final Map<String, dynamic> configuration = {
        'iceServers': [
          {'urls': 'stun:stun.l.google.com:19302'},
        ]
      };

      _peerConnection = await createPeerConnection(configuration);

      // Add local stream tracks to connection
      _localStream!.getTracks().forEach((track) {
        _peerConnection!.addTrack(track, _localStream!);
      });

      // 3. Register Signaling Callbacks
      _peerConnection!.onIceCandidate = (RTCIceCandidate candidate) {
        _wsService.sendMessage({
          "type": "ice",
          "conversationId": widget.conversationId,
          "receiverId": widget.receiverId,
          "data": {
            "candidate": candidate.candidate,
            "sdpMid": candidate.sdpMid,
            "sdpMLineIndex": candidate.sdpMLineIndex
          }
        });
      };

      _peerConnection!.onTrack = (RTCTrackEvent event) {
        if (event.streams.isNotEmpty && mounted) {
          setState(() {
            _remoteRenderer.srcObject = event.streams[0];
            _isConnected = true;
          });
        }
      };

      _peerConnection!.onConnectionState = (RTCPeerConnectionState state) {
        print("WebRTC connection state mutated: $state");
        if (state == RTCPeerConnectionState.RTCPeerConnectionStateDisconnected ||
            state == RTCPeerConnectionState.RTCPeerConnectionStateFailed) {
          _handleHangUp(remoteEnded: true);
        }
      };

      // 4. Listen to signaling events from WebSocket
      _wsSubscription = _wsService.messageStream.listen(_handleSignalingEvent);

      // 5. Setup SDP Handshake
      if (widget.isCaller) {
        // Create SDP Offer
        RTCSessionDescription offer = await _peerConnection!.createOffer();
        await _peerConnection!.setLocalDescription(offer);

        _wsService.sendMessage({
          "type": "offer",
          "conversationId": widget.conversationId,
          "receiverId": widget.receiverId,
          "data": {
            "sdp": offer.sdp,
            "type": "offer"
          }
        });
      } else {
        // Callee: Consume received SDP Offer and reply with SDP Answer
        if (widget.initialOfferSdp != null) {
          await _peerConnection!.setRemoteDescription(
            RTCSessionDescription(widget.initialOfferSdp!, 'offer'),
          );
          
          RTCSessionDescription answer = await _peerConnection!.createAnswer();
          await _peerConnection!.setLocalDescription(answer);

          _wsService.sendMessage({
            "type": "answer",
            "conversationId": widget.conversationId,
            "receiverId": widget.receiverId,
            "data": {
              "sdp": answer.sdp,
              "type": "answer",
              "sessionId": _sessionId ?? ""
            }
          });
        }
      }
    } catch (e) {
      print("Error initializing WebRTC Call: $e");
      if (mounted) {
        ToastService.showToast(
          context: context,
          message: "Failed to initialize video call: $e",
          title: "Hardware Permission Error",
          type: ToastType.error,
        );
        Navigator.pop(context);
      }
    }
  }

  void _handleSignalingEvent(Map<String, dynamic> message) {
    final type = message['type']?.toString();
    final senderId = message['senderId']?.toString();
    if (senderId != widget.receiverId) return; // Ignore other peers

    final data = message['data'];
    if (data == null) return;

    if (type == 'answer' && widget.isCaller) {
      final sdp = data['sdp']?.toString();
      _sessionId = data['sessionId']?.toString();
      if (sdp != null) {
        _peerConnection?.setRemoteDescription(RTCSessionDescription(sdp, 'answer'));
      }
    } else if (type == 'ice') {
      final candidate = data['candidate']?.toString();
      final sdpMid = data['sdpMid']?.toString();
      final sdpMLineIndex = data['sdpMLineIndex'] as int?;
      if (candidate != null && sdpMid != null && sdpMLineIndex != null) {
        _peerConnection?.addCandidate(
          RTCIceCandidate(candidate, sdpMid, sdpMLineIndex),
        );
      }
    } else if (type == 'leave') {
      _handleHangUp(remoteEnded: true);
    }
  }

  void _handleHangUp({bool remoteEnded = false}) {
    // 1. Release streams and peer connection
    _localStream?.getTracks().forEach((track) => track.stop());
    _peerConnection?.close();
    _wsSubscription?.cancel();
    
    _localRenderer.srcObject = null;
    _remoteRenderer.srcObject = null;

    if (!remoteEnded) {
      // Send leave message to remote peer
      _wsService.sendMessage({
        "type": "leave",
        "conversationId": widget.conversationId,
        "receiverId": widget.receiverId,
        "data": {
          "sessionId": _sessionId ?? ""
        }
      });
    }

    if (mounted) {
      ToastService.showToast(
        context: context,
        message: remoteEnded ? "Call hung up by consultant." : "Consultation room call disconnected.",
        title: "Call Ended",
        type: ToastType.info,
      );
      Navigator.pop(context);
    }
  }

  void _toggleMute() {
    if (_localStream != null) {
      final audioTrack = _localStream!.getAudioTracks().firstOrNull;
      if (audioTrack != null) {
        audioTrack.enabled = !audioTrack.enabled;
        setState(() {
          _isAudioMuted = !audioTrack.enabled;
        });
      }
    }
  }

  void _toggleVideo() {
    if (_localStream != null) {
      final videoTrack = _localStream!.getVideoTracks().firstOrNull;
      if (videoTrack != null) {
        videoTrack.enabled = !videoTrack.enabled;
        setState(() {
          _isVideoOff = !videoTrack.enabled;
        });
      }
    }
  }

  void _switchCamera() {
    if (_localStream != null) {
      final videoTrack = _localStream!.getVideoTracks().firstOrNull;
      if (videoTrack != null) {
        Helper.switchCamera(videoTrack);
      }
    }
  }

  @override
  void dispose() {
    _localRenderer.dispose();
    _remoteRenderer.dispose();
    _wsSubscription?.cancel();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    final theme = Theme.of(context);

    return Scaffold(
      backgroundColor: Colors.black87,
      body: Stack(
        children: [
          // 1. Remote Full-screen Video
          Positioned.fill(
            child: _remoteRenderer.srcObject != null
                ? RTCVideoView(
                    _remoteRenderer,
                    objectFit: RTCVideoViewObjectFit.RTCVideoViewObjectFitCover,
                  )
                : Container(
                    color: Colors.black87,
                    child: Center(
                      child: Column(
                        mainAxisAlignment: MainAxisAlignment.center,
                        children: [
                          const CircularProgressIndicator(color: Colors.white),
                          const SizedBox(height: 16),
                          Text(
                            "Calling ${widget.receiverName}...",
                            style: GoogleFonts.quicksand(
                              color: Colors.white,
                              fontSize: 18,
                              fontWeight: FontWeight.w600,
                            ),
                          ),
                        ],
                      ),
                    ),
                  ),
          ),

          // 2. Small Local Camera Overlay (Pip style)
          Positioned(
            top: 50,
            right: 20,
            child: ClipRRect(
              borderRadius: BorderRadius.circular(16.0),
              child: Container(
                height: 160,
                width: 120,
                color: Colors.grey[900],
                child: _localRenderer.srcObject != null && !_isVideoOff
                    ? RTCVideoView(
                        _localRenderer,
                        mirror: true,
                        objectFit: RTCVideoViewObjectFit.RTCVideoViewObjectFitCover,
                      )
                    : const Center(
                        child: Icon(
                          Icons.videocam_off_rounded,
                          color: Colors.white,
                          size: 32,
                        ),
                      ),
              ),
            ),
          ),

          // 3. Floating Custom Control HUD
          Positioned(
            bottom: 40,
            left: 20,
            right: 20,
            child: Container(
              padding: const EdgeInsets.symmetric(vertical: 16, horizontal: 24),
              decoration: BoxDecoration(
                // ignore: deprecated_member_use
                color: Colors.black.withOpacity(0.6),
                borderRadius: BorderRadius.circular(30),
              ),
              child: Row(
                mainAxisAlignment: MainAxisAlignment.spaceEvenly,
                children: [
                  // Mute Mic Button
                  _buildControlBtn(
                    icon: _isAudioMuted ? Icons.mic_off_rounded : Icons.mic_rounded,
                    color: _isAudioMuted ? Colors.red : Colors.grey[800]!,
                    onPressed: _toggleMute,
                  ),
                  // Turn Off Camera Button
                  _buildControlBtn(
                    icon: _isVideoOff ? Icons.videocam_off_rounded : Icons.videocam_rounded,
                    color: _isVideoOff ? Colors.red : Colors.grey[800]!,
                    onPressed: _toggleVideo,
                  ),
                  // Switch Camera Button
                  _buildControlBtn(
                    icon: Icons.cameraswitch_rounded,
                    color: Colors.grey[800]!,
                    onPressed: _switchCamera,
                  ),
                  // Red Hangup Button
                  _buildControlBtn(
                    icon: Icons.call_end_rounded,
                    color: Colors.red,
                    onPressed: () => _handleHangUp(remoteEnded: false),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildControlBtn({
    required IconData icon,
    required Color color,
    required VoidCallback onPressed,
  }) {
    return Container(
      height: 52,
      width: 52,
      decoration: BoxDecoration(
        color: color,
        shape: BoxShape.circle,
      ),
      child: IconButton(
        icon: Icon(icon, color: Colors.white, size: 24),
        onPressed: onPressed,
      ),
    );
  }
}
