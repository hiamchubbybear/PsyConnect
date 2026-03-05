import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, signal, ViewChild } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { combineLatest, finalize, of, Subscription } from 'rxjs';
import { catchError } from 'rxjs/operators';
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
    RouterModule,
    TranslateModule,
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
  suggestedUsers: Friend[] = [];
  isLoading = false;
  isLoadingFriends = true;
  isLoadingMessages = false;

  
  showConsentDialog = false;
  private readonly CONSENT_KEY = 'psy_chat_consent_accepted';

  
  
  private readonly STRANGERS_KEY = 'psy_chat_strangers';
  private strangers: Friend[] = [];

  
  showVideoCall = false;
  incomingCall: IncomingCall | null = null;
  currentCallType: CallType = 'video'; 
  activeCallerId: string | null = null; 
  activeSessionId: string | null = null; 
  private isEndingCall = false; 

  
  private ringtoneAudio: HTMLAudioElement | null = null;
  private outgoingAudio: HTMLAudioElement | null = null;
  private audioUnlocked = false; 
  @ViewChild(VideoCallComponent) videoCallComponent?: VideoCallComponent;

  private readonly GROUPING_THRESHOLD = 5 * 60 * 1000; 

  constructor(
    private friendService: FriendService,
    private chatService: ChatService,
    private loader: LoaderService,
    private route: ActivatedRoute,
    private router: Router,
    private webrtcService: WebRTCService,
    private signalingService: WebRTCSignalingService,
    private translate: TranslateService,
  ) {}

  ngOnInit(): void {
    
    if (!sessionStorage.getItem(this.CONSENT_KEY)) {
      this.showConsentDialog = true;
    }

    
    try {
      const stored = sessionStorage.getItem(this.STRANGERS_KEY);
      this.strangers = stored ? JSON.parse(stored) : [];
    } catch {
      this.strangers = [];
    }

    
    
    console.log('👂 Setting up incoming call subscription in ngOnInit...');
    this.incomingCallSubscription =
      this.signalingService.incomingCall$.subscribe((call) => {
        console.log('🔔 incomingCall$ emitted:', call);
        this.incomingCall = call;
        console.log('✅ this.incomingCall set to:', this.incomingCall);

        
        if (call) {
          this.playRingtone();
        } else {
          this.stopRingtone();
        }
      });

    
    this.webrtcService.callEnded$.subscribe(() => {
      console.log('📞 Call ended by remote user');
      
      this.cleanupCallUI();
    });

    
    this.chatService.getCurrentUser().subscribe((profile) => {
      if (profile) {
        this.currentUser = mapUserProfileToChatUser(profile);
        console.log('✅ Current user loaded:', this.currentUser.name);
        this.loadFriendsAndConversations();
      }
    });

    
    this.sub = this.route.paramMap.subscribe((params) => {
      const id = params.get('id');
      if (!id || this.isLoadingFriends) return;
      const alreadySelected = this.selectedFriend()?.profileId === id;
      if (alreadySelected) return;
      const friend = this.friends.find((f) => f.profileId === id);
      if (friend) {
        this.onFriendSelected(friend, false);
      } else {
        this.openChatByProfileId(id);
      }
    });

    
    this.unlockAudio();

    
    this.friendService.getFriendSuggestions().subscribe({
      next: (suggestions) => {
        this.suggestedUsers = suggestions.slice(0, 4);
      },
      error: () => {
        this.suggestedUsers = [];
      },
    });
  }

  private loadFriendsAndConversations() {
    const userId = this.currentUser.profileId;
    if (!userId || userId === 'fallback' || userId === 'unknown') {
      console.warn(
        '⚠️ loadFriendsAndConversations: profileId not ready yet:',
        userId,
      );
      
      this.friendService.getMyFriends().subscribe({
        next: (friends) => {
          this.friends = friends;
          this.isLoadingFriends = false;
          const currentId = this.route.snapshot.paramMap.get('id');
          if (currentId) {
            const friend = friends.find((f) => f.profileId === currentId);
            if (friend) this.onFriendSelected(friend, false);
            else this.openChatByProfileId(currentId);
          } else if (friends.length > 0) {
            this.onFriendSelected(friends[0], true);
          }
        },
        error: () => {
          this.isLoadingFriends = false;
        },
      });
      return;
    }

    console.log('📡 Loading conversations for profileId:', userId);
    combineLatest([
      this.friendService.getMyFriends(),
      this.chatService.getRecentConversations(userId).pipe(
        catchError((err) => {
          console.error('Error fetching recent conversations:', err);
          return of([]);
        }),
      ),
    ]).subscribe({
      next: ([friends, conversations]) => {
        console.log('✅ Friends and Conversations loaded');

        
        const enrichedFriends: Friend[] = friends.map((f): Friend => {
          
          const conv = conversations.find((c) =>
            c.participants?.includes(f.profileId),
          );
          if (conv && conv.lastMessage) {
            let parsedText = conv.lastMessage.text;
            try {
              parsedText =
                typeof parsedText === 'string'
                  ? JSON.parse(parsedText)
                  : String(parsedText);
            } catch (e) {
              
            }
            return {
              ...f,
              lastMessage: parsedText,
              lastMessageTime: new Date(conv.lastMessage.createdAt),
            };
          } else if (conv) {
            return {
              ...f,
              lastMessageTime: new Date(conv.createdAt),
            };
          }
          return f;
        });

        
        enrichedFriends.sort((a, b) => {
          const timeA = a.lastMessageTime?.getTime() || 0;
          const timeB = b.lastMessageTime?.getTime() || 0;
          return timeB - timeA;
        });

        this.friends = enrichedFriends;
        this.isLoadingFriends = false;

        
        const currentId = this.route.snapshot.paramMap.get('id');
        if (currentId) {
          const friend = this.friends.find((f) => f.profileId === currentId);
          if (friend) {
            this.onFriendSelected(friend, false);
          } else {
            this.openChatByProfileId(currentId);
          }
        } else if (this.friends.length > 0) {
          this.onFriendSelected(this.friends[0], true);
        }
      },
      error: (err) => {
        console.error('❌ Error loading friends/conversations:', err);
        this.isLoadingFriends = false;
        this.friends = [];
      },
    });
  }

  
  private unlockAudio() {
    const unlock = () => {
      if (this.audioUnlocked) return;

      
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
          
        });
    };

    
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

  onBack() {
    this.selectedFriend.set(null);
    this.router.navigate(['/feature/chat']);
  }

  
  openChatByProfileId(profileId: string) {
    
    const existingFriend = this.friends.find((f) => f.profileId === profileId);
    if (existingFriend) {
      this.onFriendSelected(existingFriend, false);
      return;
    }

    
    const existingStranger = this.strangers.find(
      (f) => f.profileId === profileId,
    );
    if (existingStranger) {
      this.onFriendSelected(existingStranger, false);
      return;
    }

    
    const minimalFriend: Friend = {
      profileId,
      firstName: 'User',
      lastName: '',
      avatarUri: '',
    };

    this.strangers = [minimalFriend, ...this.strangers];
    this.saveStrangers();
    this.onFriendSelected(minimalFriend, false);
  }

  
  get allContacts(): Friend[] {
    const friendIds = new Set(this.friends.map((f) => f.profileId));
    const uniqueStrangers = this.strangers.filter(
      (s) => !friendIds.has(s.profileId),
    );
    return [...this.friends, ...uniqueStrangers];
  }

  private saveStrangers() {
    try {
      sessionStorage.setItem(
        this.STRANGERS_KEY,
        JSON.stringify(this.strangers),
      );
    } catch {}
  }

  
  acceptConsent() {
    sessionStorage.setItem(this.CONSENT_KEY, 'true');
    this.showConsentDialog = false;
    this.addSystemMessage(
      this.translate.instant('CHAT.SYSTEM.ConsentAccepted'),
    );
  }

  private addSystemMessage(content: string) {
    const sysMsg: Message = {
      id: 'system-' + Date.now(),
      conversationId: this.conversationId(),
      senderId: 'system',
      content: content,
      timestamp: new Date(),
      userName: 'System',
      userAvatar: '',
      isMine: false,
      isSystem: true,
    };

    this.messages.update((prev) => [...prev, sysMsg]);
  }

  
  declineConsent() {
    this.router.navigate(['/feature/feed']);
  }

  onFriendSelected(friend: Friend, updateUrl = false) {
    
    this.wsSubscription?.unsubscribe();

    this.selectedFriend.set(friend);
    if (updateUrl) {
      this.router.navigate(['/feature/chat', friend.profileId]);
    }
    this.isLoading = true;
    this.chatService
      .getOrCreateConversation(this.currentUser.profileId, friend.profileId)
      .pipe(finalize(() => (this.isLoading = false)))
      .subscribe((conv) => {
        this.conversationId.set(conv.id);

        
        this.isLoadingMessages = true;
        this.chatService
          .getChatsByConversation(conv.id, this.currentUser.profileId)
          .pipe(finalize(() => (this.isLoadingMessages = false)))
          .subscribe((msgs) => {
            const enrichedMsgs = msgs.map((m) => this.enrichMessage(m, friend));
            const groupedMsgs = this.groupMessages(enrichedMsgs);
            this.messages.set(groupedMsgs);

            
            this.wsSubscription = this.chatService
              .connect(friend.profileId, conv.id)
              .subscribe({
                next: (wsMsg: any) => {
                  console.log('📨 Received WebSocket message:', wsMsg);

                  
                  const webrtcTypes = ['offer', 'answer', 'ice', 'leave'];
                  if (webrtcTypes.includes(wsMsg.type)) {
                    console.log(
                      '🔀 Routing WebRTC message to signaling service:',
                      wsMsg.type,
                    );

                    
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
                    
                    (this.signalingService as any).handleMessage(wsMsg);
                    return; 
                  }

                  
                  const msg: Message = {
                    id: wsMsg.data?.id || '',
                    conversationId:
                      wsMsg.conversationId || wsMsg.data?.conversationId,
                    senderId: wsMsg.senderId || wsMsg.data?.senderId,
                    content: wsMsg.data?.content || '',
                    timestamp: new Date(wsMsg.data?.timestamp || Date.now()),
                    userName: '', 
                    userAvatar: '', 
                    isMine: false, 
                  };

                  
                  const enrichedMsg = this.enrichMessage(msg, friend);

                  
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

  
  handleCallRequest(callType: CallType) {
    console.log('📞', callType, 'call requested');
    this.currentCallType = callType; 

    const friend = this.selectedFriend();
    if (!friend || !this.conversationId()) {
      console.error('❌ No friend or conversation selected');
      return;
    }

    
    
    const sessionId = crypto.randomUUID(); 
    this.chatService
      .startCall({
        conversationId: this.conversationId(),
        callerId: this.currentUser.profileId,
        callerName: this.currentUser.name,
        receiverId: friend.profileId,
        sessionId: sessionId,
      })
      .subscribe();

    
    this.activeCallerId = friend.profileId;
    
    
    
    
    
    

    

    
    this.activeCallerId = friend.profileId;
    
    console.log('✅ Set activeCallerId:', this.activeCallerId);

    
    
    const wsUrl = environment.wsUrl;
    this.signalingService.connect(
      wsUrl,
      this.conversationId(),
      friend.profileId,
      this.currentUser.profileId,
    );

    
    this.showVideoCall = true;

    
    const connectionSub = this.signalingService.connected$.subscribe(
      (connected) => {
        if (connected) {
          console.log('✅ WebSocket ready, starting call...');

          
          this.playOutgoingTone();

          setTimeout(() => {
            if (this.videoCallComponent) {
              this.videoCallComponent.startCall();
            }
            connectionSub.unsubscribe(); 
          }, 100);
        }
      },
    );
  }

  onCallEnded() {
    
    if (this.isEndingCall) {
      console.log('⚠️  Already ending call, skipping...');
      return;
    }

    this.isEndingCall = true;
    console.log('📞 Call ended - cleaning up...');

    
    if (this.activeCallerId && this.activeSessionId) {
      console.log('📤 Sending leave message to:', this.activeCallerId);
      this.signalingService.sendLeave(
        this.activeCallerId,
        this.activeSessionId,
      );
    }

    
    
    this.webrtcService.endCall();

    
    this.cleanupCallUI();

    console.log('✅ Call cleanup complete');

    
    setTimeout(() => {
      this.isEndingCall = false;
    }, 1000);
  }

  
  private cleanupCallUI() {
    console.log('🧹 Cleaning up call UI...');

    
    this.stopRingtone();
    this.stopOutgoingTone();

    
    if (!this.isEndingCall) {
      this.playCallEnd();
    }

    
    

    
    this.activeCallerId = null;
    this.activeSessionId = null;

    
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

    
    this.activeSessionId = sessionId;
    console.log('✅ Set activeSessionId:', this.activeSessionId);
  }

  onAnswer(answer: RTCSessionDescriptionInit) {
    console.log('📥 onAnswer called with:', answer);

    
    this.stopOutgoingTone();

    
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

    
    this.stopRingtone();

    
    this.activeCallerId = this.incomingCall.callerId;
    this.activeSessionId = this.incomingCall.sessionId;
    const offer = this.incomingCall.offer;

    
    this.incomingCall = null;

    
    this.showVideoCall = true;

    
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
    this.stopRingtone(); 
  }

  
  private playRingtone() {
    try {
      if (!this.ringtoneAudio) {
        this.ringtoneAudio = new Audio('/sounds/ring-tone.mp3');
        this.ringtoneAudio.loop = true;
        this.ringtoneAudio.volume = 0.3; 
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
        this.outgoingAudio.volume = 0.3; 
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
      audio.volume = 0.4; 
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
      audio.volume = 0.4; 
      audio.play().catch((err) => {
        console.error('❌ Error playing call end sound:', err);
      });
    } catch (error) {
      console.error('❌ Error playing call end sound:', error);
    }
  }
}
