import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { WebRTCService } from '../webrtc/webrtc.service';

export interface WebSocketMessage {
  type: 'chat' | 'offer' | 'answer' | 'ice' | 'leave';
  conversationId: string;
  senderId: string;
  receiverId: string;
  data: any;
}

export interface IncomingCall {
  conversationId: string;
  callerId: string;
  callerName: string;
  sessionId: string;
  offer: RTCSessionDescriptionInit;
}

@Injectable({
  providedIn: 'root',
})
export class WebRTCSignalingService {
  private ws: WebSocket | null = null;
  private currentUserId: string = '';
  private conversationId: string = '';

  private incomingCallSubject = new BehaviorSubject<IncomingCall | null>(null);
  public incomingCall$ = this.incomingCallSubject.asObservable();

  private connectedSubject = new BehaviorSubject<boolean>(false);
  public connected$ = this.connectedSubject.asObservable();

  constructor(
    private webrtcService: WebRTCService,
    private secureStorage: SecureStorageService
  ) {}

  // Connect to WebSocket
  connect(
    wsUrl: string,
    conversationId: string,
    receiverId: string,
    currentUserId: string
  ) {
    this.currentUserId = currentUserId;
    this.conversationId = conversationId;

    // Get JWT token from secure storage
    const token = this.secureStorage.getItem<string>(
      environment.accessTokenKey
    );
    if (!token) {
      console.error('❌ No JWT token found in secure storage');
      return;
    }

    // Add token to URL query parameters
    const url = `${wsUrl}?conversationId=${conversationId}&receiver=${receiverId}&token=${token}`;
    console.log('🔌 Connecting to WebSocket:', wsUrl);

    this.ws = new WebSocket(url);

    this.ws.onopen = () => {
      console.log('✅ WebSocket connected for WebRTC signaling');
      console.log('🔔 Setting connectedSubject to true');
      this.connectedSubject.next(true);
      console.log('🔔 connectedSubject value:', this.connectedSubject.value);
    };

    this.ws.onmessage = (event) => {
      this.handleMessage(JSON.parse(event.data));
    };

    this.ws.onerror = (error) => {
      console.error('❌ WebSocket error:', error);
    };

    this.ws.onclose = (event) => {
      console.log(
        '🔌 WebSocket disconnected. Code:',
        event.code,
        'Reason:',
        event.reason
      );
    };
  }

  // Handle incoming WebSocket messages
  private handleMessage(message: WebSocketMessage) {
    console.log('📨 Received message:', message.type, message);

    switch (message.type) {
      case 'offer':
        this.handleIncomingOffer(message);
        break;

      case 'answer':
        this.handleIncomingAnswer(message);
        break;

      case 'ice':
        this.handleIncomingICE(message);
        break;

      case 'leave':
        this.handleIncomingLeave(message);
        break;
    }
  }

  // Handle incoming call offer
  private handleIncomingOffer(message: WebSocketMessage) {
    console.log('📞 Incoming call from:', message.senderId);
    console.log('📞 Message data:', message.data);

    const incomingCall: IncomingCall = {
      conversationId: message.conversationId,
      callerId: message.senderId,
      callerName: message.data.callerName || 'Unknown',
      sessionId: message.data.sessionId,
      offer: {
        sdp: message.data.sdp,
        type: message.data.type,
      },
    };

    console.log('🔔 Emitting incomingCall:', incomingCall);
    this.incomingCallSubject.next(incomingCall);
    console.log(
      '✅ incomingCall emitted, current value:',
      this.incomingCallSubject.value
    );
  }

  // Handle incoming answer
  private async handleIncomingAnswer(message: WebSocketMessage) {
    console.log('📥 Received answer from:', message.senderId);

    const answer: RTCSessionDescriptionInit = {
      sdp: message.data.sdp,
      type: message.data.type,
    };

    await this.webrtcService.handleAnswer(answer);
  }

  // Handle incoming ICE candidate
  private async handleIncomingICE(message: WebSocketMessage) {
    console.log('🧊 Received ICE candidate from:', message.senderId);

    const candidate: RTCIceCandidateInit = {
      candidate: message.data.candidate,
      sdpMid: message.data.sdpMid,
      sdpMLineIndex: message.data.sdpMLineIndex,
    };

    await this.webrtcService.addIceCandidate(candidate);
  }

  // Handle call end
  private handleIncomingLeave(message: WebSocketMessage) {
    console.log('👋 Call ended by:', message.senderId);
    this.webrtcService.endCall();
    this.incomingCallSubject.next(null);
  }

  // Send offer
  sendOffer(receiverId: string, offer: RTCSessionDescriptionInit): string {
    // Generate session ID for this call
    const sessionId = crypto.randomUUID();

    const message: WebSocketMessage = {
      type: 'offer',
      conversationId: this.conversationId,
      senderId: this.currentUserId,
      receiverId: receiverId,
      data: {
        sdp: offer.sdp,
        type: offer.type,
        sessionId: sessionId, // Include session ID in offer
      },
    };

    console.log('📤 Sending offer message:', message);
    console.log(
      '📤 WebSocket state:',
      this.ws?.readyState,
      'OPEN=',
      WebSocket.OPEN
    );
    this.send(message);
    console.log('📤 Sent offer to:', receiverId, 'with sessionId:', sessionId);

    return sessionId; // Return session ID
  }

  // Send answer
  sendAnswer(
    receiverId: string,
    answer: RTCSessionDescriptionInit,
    sessionId: string
  ) {
    const message: WebSocketMessage = {
      type: 'answer',
      conversationId: this.conversationId,
      senderId: this.currentUserId,
      receiverId: receiverId,
      data: {
        sdp: answer.sdp,
        type: answer.type,
        sessionId: sessionId,
      },
    };

    this.send(message);
    console.log('📤 Sent answer to:', receiverId);
  }

  // Send ICE candidate
  sendICECandidate(receiverId: string, candidate: RTCIceCandidate) {
    const message: WebSocketMessage = {
      type: 'ice',
      conversationId: this.conversationId,
      senderId: this.currentUserId,
      receiverId: receiverId,
      data: {
        candidate: candidate.candidate,
        sdpMid: candidate.sdpMid,
        sdpMLineIndex: candidate.sdpMLineIndex,
      },
    };

    this.send(message);
    console.log('📤 Sent ICE candidate to:', receiverId);
  }

  // Send leave (end call)
  sendLeave(receiverId: string, sessionId?: string) {
    const message: WebSocketMessage = {
      type: 'leave',
      conversationId: this.conversationId,
      senderId: this.currentUserId,
      receiverId: receiverId,
      data: {
        sessionId: sessionId,
      },
    };

    this.send(message);
    console.log('📤 Sent leave to:', receiverId);
  }

  // Send message via WebSocket
  private send(message: WebSocketMessage) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      this.ws.send(JSON.stringify(message));
    } else {
      console.error('❌ WebSocket not connected');
    }
  }

  // Disconnect WebSocket
  disconnect() {
    if (this.ws) {
      this.ws.close();
      this.ws = null;
    }
  }

  // Clear incoming call
  clearIncomingCall() {
    this.incomingCallSubject.next(null);
  }
}
