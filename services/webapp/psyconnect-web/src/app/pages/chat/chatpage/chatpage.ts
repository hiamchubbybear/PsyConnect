import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, signal, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { combineLatest, finalize, Subscription } from 'rxjs';
import { environment } from '../../../../environments/environment';
import { CallType } from '../../../components/call-options-menu/call-options-menu.component';
import { ChatListComponent } from '../../../components/chat/chat-list/chat-list';
import { ChatMainComponent } from '../../../components/chat/chat-main/chat-main';
import { ChatSkeletonComponent } from '../../../components/chat/chat-skeleton/chat-skeleton';
import { IncomingCallComponent } from '../../../components/incoming-call/incoming-call.component';
import { VideoCallComponent } from '../../../components/video-call/video-call.component';
import { Friend, Message } from '../../../models/chat.models';
import { ChatUser, mapUserProfileToChatUser } from '../../../models/map.utils';
import { ChatService } from '../../../services/chat/chat.service';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { LoaderService } from '../../../services/loader/loader';
import {
  IncomingCall,
  WebRTCSignalingService,
} from '../../../services/webrtc/webrtc-signaling.service';
import { WebRTCService } from '../../../services/webrtc/webrtc.service';

@Component({
  selector: 'chat-page',
  standalone: true,
  templateUrl: './chatpage.html',
  styleUrls: ['./chatpage.scss'],
  imports: [
    ChatListComponent,
    ChatMainComponent,
    ChatSkeletonComponent,
    VideoCallComponent,
    IncomingCallComponent,
    CommonModule,
    FormsModule,
  ],
})
export class ChatComponent implements OnInit, OnDestroy {
  private sub?: Subscription;
  private wsSubscription?: Subscription;
  private incomingCallSubscription?: Subscription;
  messageText = '';
  conversationId = signal<string>('');
  currentUser: ChatUser = mapUserProfileToChatUser(null);
  selectedFriend = signal<Friend | null>(null);
  messages = signal<Message[]>([]);
  friends: Friend[] = [];
  isLoading = false;
  isLoadingFriends = true;
  isLoadingMessages = false;

  // WebRTC properties
  showVideoCall = false;
  incomingCall: IncomingCall | null = null;
  currentCallType: CallType = 'video'; // Track current call type
  activeCallerId: string | null = null; // Store caller ID for answer
  activeSessionId: string | null = null; // Store session ID for answer
  private isEndingCall = false; // Prevent infinite loop

  // Ringtone audio
  private ringtoneAudio: HTMLAudioElement | null = null;
  private outgoingAudio: HTMLAudioElement | null = null;
  private audioUnlocked = false; // Track if audio context is unlocked
  @ViewChild(VideoCallComponent) videoCallComponent?: VideoCallComponent;

  private readonly GROUPING_THRESHOLD = 5 * 60 * 1000; // 5 minutes

  constructor(
    private friendService: FriendService,
    private chatService: ChatService,
    private loader: LoaderService,
    private route: ActivatedRoute,
    private router: Router,
    private webrtcService: WebRTCService,
    private signalingService: WebRTCSignalingService,
  ) {}

  ngOnInit(): void {
    // Setup incoming call subscription IMMEDIATELY
    // This ensures we can receive calls even if we haven't initiated a call
    console.log('👂 Setting up incoming call subscription in ngOnInit...');
    this.incomingCallSubscription =
      this.signalingService.incomingCall$.subscribe((call) => {
        console.log('🔔 incomingCall$ emitted:', call);
        this.incomingCall = call;
        console.log('✅ this.incomingCall set to:', this.incomingCall);

        // Play ringtone when receiving call
        if (call) {
          this.playRingtone();
        } else {
          this.stopRingtone();
        }
      });

    // Subscribe to call ended event to close UI when remote user hangs up
    this.webrtcService.callEnded$.subscribe(() => {
      console.log('📞 Call ended by remote user');
      // Only cleanup UI, don't send leave again (already sent by remote)
      this.cleanupCallUI();
    });

    // Load current user immediately
    this.chatService.getCurrentUser().subscribe((profile) => {
      if (profile) {
        this.currentUser = mapUserProfileToChatUser(profile);
        console.log('✅ Current user loaded:', this.currentUser.name);
      }
    });

    // Load friends separately
    combineLatest([
      this.friendService.getMyFriends(),
      this.route.paramMap,
    ]).subscribe({
      next: ([friends, params]) => {
        console.log('✅ Friends loaded:', friends.length);
        this.friends = friends;
        console.log('📋 this.friends:', this.friends);

        // Set loading false AFTER friends are set
        this.isLoadingFriends = false;
        console.log('🔍 isLoadingFriends NOW:', this.isLoadingFriends);

        const idFromParam = params.get('id');
        if (idFromParam) {
          const friend = friends.find((f) => f.profileId === idFromParam);
          if (friend) {
            this.onFriendSelected(friend, false);
            return;
          }
        }
        if (friends.length > 0) {
          this.onFriendSelected(friends[0], true);
        }
      },
      error: (err) => {
        console.error('❌ Error loading friends:', err);
        // Even if error, stop loading so user can see their profile (empty list)
        this.isLoadingFriends = false;
        this.friends = [];
      },
    });

    // Unlock audio on first user interaction
    this.unlockAudio();
  }

  // Unlock audio context for autoplay
  private unlockAudio() {
    const unlock = () => {
      if (this.audioUnlocked) return;

      // Play silent audio to unlock audio context
      const silentAudio = new Audio();
      silentAudio.src =
        'data:audio/mp3;base64,SUQzBAAAAAABEVRYWFgAAAAtAAADY29tbWVudABCaWdTb3VuZEJhbmsuY29tIC8gTGFTb25vdGhlcXVlLm9yZwBURU5DAAAAHQAAA1N3aXRjaCBQbHVzIMKpIE5DSCBTb2Z0d2FyZQBUSVQyAAAABgAAAzIyMzUAVFNTRQAAAA8AAANMYXZmNTcuODMuMTAwAAAAAAAAAAAAAAD/80DEAAAAA0gAAAAATEFNRTMuMTAwVVVVVVVVVVVVVUxBTUUzLjEwMFVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVf/zQsRbAAADSAAAAABVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVf/zQMSkAAADSAAAAABVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVVV';
      silentAudio.volume = 0.01;
      silentAudio
        .play()
        .then(() => {
          console.log('🔓 Audio context unlocked');
          this.audioUnlocked = true;
          document.removeEventListener('click', unlock);
          document.removeEventListener('touchstart', unlock);
        })
        .catch(() => {
          // Ignore errors
        });
    };

    // Listen for first click/touch
    document.addEventListener('click', unlock, { once: true });
    document.addEventListener('touchstart', unlock, { once: true });
  }
  ngOnDestroy(): void {
    this.sub?.unsubscribe();
    this.wsSubscription?.unsubscribe();
    this.incomingCallSubscription?.unsubscribe();
    this.chatService.disconnect();
    this.stopRingtone();
    this.stopOutgoingTone();
  }

  onFriendSelected(friend: Friend, updateUrl = false) {
    // Cleanup previous WebSocket subscription
    this.wsSubscription?.unsubscribe();

    this.selectedFriend.set(friend);
    if (updateUrl) {
      this.router.navigate(['/chat', friend.profileId]);
    }
    this.isLoading = true;
    this.chatService
      .getOrCreateConversation(this.currentUser.profileId, friend.profileId)
      .pipe(finalize(() => (this.isLoading = false)))
      .subscribe((conv) => {
        this.conversationId.set(conv.id);

        // Load chat history FIRST
        this.isLoadingMessages = true;
        this.chatService
          .getChatsByConversation(conv.id, this.currentUser.profileId)
          .pipe(finalize(() => (this.isLoadingMessages = false)))
          .subscribe((msgs) => {
            const enrichedMsgs = msgs.map((m) => this.enrichMessage(m, friend));
            const groupedMsgs = this.groupMessages(enrichedMsgs);
            this.messages.set(groupedMsgs);

            // THEN connect WebSocket to receive new messages
            this.wsSubscription = this.chatService
              .connect(friend.profileId, conv.id)
              .subscribe({
                next: (wsMsg: any) => {
                  console.log('📨 Received WebSocket message:', wsMsg);

                  // Check if this is a WebRTC signaling message
                  const webrtcTypes = ['offer', 'answer', 'ice', 'leave'];
                  if (webrtcTypes.includes(wsMsg.type)) {
                    console.log(
                      '🔀 Routing WebRTC message to signaling service:',
                      wsMsg.type,
                    );

                    // Stop outgoing tone immediately when answer received
                    if (wsMsg.type === 'answer') {
                      console.log(
                        '📞 Answer received - stopping outgoing tone',
                      );
                      this.stopOutgoingTone();
                      this.stopRingtone();
                      this.playCallConnected();
                    }

                    if (wsMsg.type === 'leave') {
                      console.log('👋 LEAVE message received:', wsMsg);
                    }
                    // Forward to signaling service's message handler
                    (this.signalingService as any).handleMessage(wsMsg);
                    return; // Don't process as chat message
                  }

                  // Extract actual message data from nested structure
                  const msg: Message = {
                    id: wsMsg.data?.id || '',
                    conversationId:
                      wsMsg.conversationId || wsMsg.data?.conversationId,
                    senderId: wsMsg.senderId || wsMsg.data?.senderId,
                    content: wsMsg.data?.content || '',
                    timestamp: new Date(wsMsg.data?.timestamp || Date.now()),
                    userName: '', // Will be enriched
                    userAvatar: '', // Will be enriched
                    isMine: false, // Will be enriched
                  };

                  // Enrich with friend info
                  const enrichedMsg = this.enrichMessage(msg, friend);

                  // Re-group all messages after adding new one
                  this.messages.update((prev) => {
                    const updated = [...prev, enrichedMsg];
                    return this.groupMessages(updated);
                  });
                },
                error: (err) => console.error('❌ WebSocket error:', err),
              });
          });
      });
  }

  private groupMessages(messages: Message[]): Message[] {
    return messages.map((msg, index) => {
      const prev = messages[index - 1];
      const next = messages[index + 1];

      const isFirstInGroup =
        !prev ||
        prev.senderId !== msg.senderId ||
        msg.timestamp.getTime() - prev.timestamp.getTime() >
          this.GROUPING_THRESHOLD;

      const isLastInGroup =
        !next ||
        next.senderId !== msg.senderId ||
        next.timestamp.getTime() - msg.timestamp.getTime() >
          this.GROUPING_THRESHOLD;

      return {
        ...msg,
        isFirstInGroup,
        isLastInGroup,
        showTimestamp: isFirstInGroup || isLastInGroup,
      };
    });
  }

  private enrichMessage(msg: Message, friend: Friend): Message {
    const isMine = msg.senderId === this.currentUser.profileId;

    return {
      ...msg,
      isMine,
      userName: isMine
        ? this.currentUser.name
        : `${friend.firstName} ${friend.lastName}`,
      userAvatar: isMine
        ? this.currentUser.avatar
        : friend.avatarUri || 'assets/avatars/default.png',
    };
  }

  // WebRTC Handlers
  handleCallRequest(callType: CallType) {
    console.log('📞', callType, 'call requested');
    this.currentCallType = callType; // Store call type

    const friend = this.selectedFriend();
    if (!friend || !this.conversationId()) {
      console.error('❌ No friend or conversation selected');
      return;
    }

    // Start call notification API (Fire and forget, or wait?)
    // This triggers the push notification to receiving user
    const sessionId = crypto.randomUUID(); // Generate session ID locally
    this.chatService
      .startCall({
        conversationId: this.conversationId(),
        callerId: this.currentUser.profileId,
        callerName: this.currentUser.name,
        receiverId: friend.profileId,
        sessionId: sessionId,
      })
      .subscribe();

    // Store call info for ending call later
    this.activeCallerId = friend.profileId;
    // Session ID will be set when offer is created. Wait, I generated it above?
    // Actually, sendOffer generates *another* session ID?
    // WebRTCSignalingService.sendOffer generates sessionId.
    // Ideally I should synchronise them, but for notification purposes, just sending A session ID is fine.
    // The actual WebRTC session ID is what matters for connection.
    // Let's rely on WebRTC service to manage the actual call flow.

    // ... rest of existing logic ...

    // Store call info for ending call later
    this.activeCallerId = friend.profileId;
    // Session ID will be set when offer is created
    console.log('✅ Set activeCallerId:', this.activeCallerId);

    // Connect WebSocket for signaling if not already connected
    // const wsUrl = 'ws://localhost:8083/ws';
    const wsUrl = environment.wsUrl;
    this.signalingService.connect(
      wsUrl,
      this.conversationId(),
      friend.profileId,
      this.currentUser.profileId,
    );

    // Show video call component
    this.showVideoCall = true;

    // Wait for WebSocket to connect before starting call
    const connectionSub = this.signalingService.connected$.subscribe(
      (connected) => {
        if (connected) {
          console.log('✅ WebSocket ready, starting call...');

          // Play outgoing tone
          this.playOutgoingTone();

          setTimeout(() => {
            if (this.videoCallComponent) {
              this.videoCallComponent.startCall();
            }
            connectionSub.unsubscribe(); // Cleanup
          }, 100);
        }
      },
    );
  }

  onCallEnded() {
    // Prevent infinite loop - if already ending, don't process again
    if (this.isEndingCall) {
      console.log('⚠️  Already ending call, skipping...');
      return;
    }

    this.isEndingCall = true;
    console.log('📞 Call ended - cleaning up...');

    // Send leave message to remote user FIRST (only if we initiated the end)
    if (this.activeCallerId && this.activeSessionId) {
      console.log('📤 Sending leave message to:', this.activeCallerId);
      this.signalingService.sendLeave(
        this.activeCallerId,
        this.activeSessionId,
      );
    }

    // End call in WebRTC service (stops all streams)
    // This will emit callEnded$, but our flag prevents re-entry
    this.webrtcService.endCall();

    // Cleanup UI (without calling endCall again)
    this.cleanupCallUI();

    console.log('✅ Call cleanup complete');

    // Reset flag after a delay
    setTimeout(() => {
      this.isEndingCall = false;
    }, 1000);
  }

  // Cleanup call UI without sending leave (for remote disconnect)
  private cleanupCallUI() {
    console.log('🧹 Cleaning up call UI...');

    // Stop all tones
    this.stopRingtone();
    this.stopOutgoingTone();

    // Play call end sound (only once)
    if (!this.isEndingCall) {
      this.playCallEnd();
    }

    // DON'T call webrtcService.endCall() here - it's already called by the service
    // when it emits callEnded$. Calling it again creates infinite loop!

    // Clear call session info
    this.activeCallerId = null;
    this.activeSessionId = null;

    // Close video call UI - CRITICAL!
    this.showVideoCall = false;
    console.log('✅ showVideoCall set to false');

    console.log('✅ UI cleanup complete');
  }

  onOffer(offer: RTCSessionDescriptionInit) {
    console.log('📥 onOffer called with:', offer);
    const friend = this.selectedFriend();
    if (!friend) {
      console.error('❌ No friend selected');
      return;
    }

    console.log('📤 Calling sendOffer for friend:', friend.profileId);
    const sessionId = this.signalingService.sendOffer(friend.profileId, offer);

    // Store session ID for ending call later
    this.activeSessionId = sessionId;
    console.log('✅ Set activeSessionId:', this.activeSessionId);
  }

  onAnswer(answer: RTCSessionDescriptionInit) {
    console.log('📥 onAnswer called with:', answer);

    // Stop outgoing tone when answer received (call connected)
    this.stopOutgoingTone();

    // Use stored activeCallerId and activeSessionId instead of incomingCall
    if (!this.activeCallerId || !this.activeSessionId) {
      console.error('❌ No active call session info');
      return;
    }

    this.signalingService.sendAnswer(
      this.activeCallerId,
      answer,
      this.activeSessionId,
    );
  }

  onICECandidate(candidate: RTCIceCandidate) {
    const friend = this.selectedFriend();
    if (!friend) return;

    this.signalingService.sendICECandidate(friend.profileId, candidate);
  }

  acceptIncomingCall() {
    if (!this.incomingCall) {
      console.error('❌ No incoming call to accept');
      return;
    }

    console.log('📞 Accepting incoming call from:', this.incomingCall.callerId);

    // Stop ringtone when accepting
    this.stopRingtone();

    // Store call info for onAnswer handler
    this.activeCallerId = this.incomingCall.callerId;
    this.activeSessionId = this.incomingCall.sessionId;
    const offer = this.incomingCall.offer;

    // Clear incoming call notification
    this.incomingCall = null;

    // Show video call component
    this.showVideoCall = true;

    // CRITICAL: Connect WebRTC signaling WebSocket BEFORE accepting call
    const friend = this.selectedFriend();
    if (!friend) {
      console.error('❌ No friend selected');
      return;
    }

    const wsUrl = environment.wsUrl;
    console.log('🔌 Connecting WebRTC signaling WebSocket for receiver...');
    this.signalingService.connect(
      wsUrl,
      this.conversationId(),
      friend.profileId,
      this.currentUser.profileId,
    );

    // Wait for WebSocket to connect, then accept call
    const connectionSub = this.signalingService.connected$.subscribe(
      (connected) => {
        if (connected) {
          console.log('✅ WebRTC WebSocket connected, accepting call...');

          setTimeout(() => {
            if (this.videoCallComponent) {
              console.log('📞 Calling videoCallComponent.acceptCall()');
              this.videoCallComponent.acceptCall(offer);
            } else {
              console.error('❌ videoCallComponent not available');
            }
            connectionSub.unsubscribe();
          }, 100);
        }
      },
    );
  }

  rejectIncomingCall() {
    if (!this.incomingCall) return;

    console.log('📞 Rejecting incoming call');
    this.signalingService.sendLeave(
      this.incomingCall.callerId,
      this.incomingCall.sessionId,
    );
    this.signalingService.clearIncomingCall();
    this.stopRingtone(); // Stop ringtone when rejecting
  }

  // Ringtone methods
  private playRingtone() {
    try {
      if (!this.ringtoneAudio) {
        this.ringtoneAudio = new Audio('/sounds/ring-tone.mp3');
        this.ringtoneAudio.loop = true;
        this.ringtoneAudio.volume = 0.3; // Reduced from 0.5
      }
      this.ringtoneAudio.play().catch((err) => {
        console.error('❌ Error playing ringtone:', err);
      });
    } catch (error) {
      console.error('❌ Error initializing ringtone:', error);
    }
  }

  private stopRingtone() {
    if (this.ringtoneAudio) {
      this.ringtoneAudio.pause();
      this.ringtoneAudio.currentTime = 0;
    }
  }

  private playOutgoingTone() {
    try {
      if (!this.outgoingAudio) {
        this.outgoingAudio = new Audio('/sounds/outgoing-call.mp3');
        this.outgoingAudio.loop = true;
        this.outgoingAudio.volume = 0.3; // Reduced from 0.5
      }
      this.outgoingAudio.play().catch((err) => {
        console.error('❌ Error playing outgoing tone:', err);
      });
    } catch (error) {
      console.error('❌ Error initializing outgoing tone:', error);
    }
  }

  private stopOutgoingTone() {
    if (this.outgoingAudio) {
      this.outgoingAudio.pause();
      this.outgoingAudio.currentTime = 0;
    }
  }

  private playCallConnected() {
    try {
      const audio = new Audio('/sounds/call-connected.mp3');
      audio.volume = 0.4; // Reduced from 0.7
      audio.play().catch((err) => {
        console.error('❌ Error playing call connected sound:', err);
      });
    } catch (error) {
      console.error('❌ Error playing call connected sound:', error);
    }
  }

  private playCallEnd() {
    try {
      const audio = new Audio('/sounds/end-call.mp3');
      audio.volume = 0.4; // Reduced from 0.7
      audio.play().catch((err) => {
        console.error('❌ Error playing call end sound:', err);
      });
    } catch (error) {
      console.error('❌ Error playing call end sound:', error);
    }
  }
}
