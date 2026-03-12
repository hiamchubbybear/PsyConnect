import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, ViewChild, ElementRef, AfterViewChecked } from '@angular/core';
import { RouterModule, Router } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { Friend, Message } from '../../../models/chat.models';
import { ChatService } from '../../../services/chat/chat.service';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { ImgFallbackDirective } from '../../../shared/directives/img-fallback.directive';
import { UserContextService, UserProfile } from '../../../services/profile/profile-service';
import { NotificationService } from '../../../services/notification/notification.service';
import { SessionService } from '../../../services/consultation/session.service';
import { ConsultationSession } from '../../../models/consultation.model';
import { SecureStorageService } from '../../../encrypt/secure';
import { environment } from '../../../../environments/environment';

@Component({
  selector: 'app-floating-chat',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, AvatarFallbackPipe, ImgFallbackDirective, FormsModule],
  templateUrl: './floating-chat.html',
  styleUrls: ['./floating-chat.scss']
})
export class FloatingChatComponent implements OnInit, OnDestroy, AfterViewChecked {
  @ViewChild('scrollViewport') private scrollViewport?: ElementRef;
  isOpen = false;
  recentFriends: Friend[] = [];
  unreadTotal = 0;
  unreadMap: { [profileId: string]: number } = {};
  
  // Compact chat state
  activeChatFriend: Friend | null = null;
  compactMessages: Message[] = [];
  allCompactMessages: Message[] = [];
  viewOffset = 0; // Steps from the end
  viewWindowSize = 12; // Show more messages for better UX
  private shouldScrollToBottom = false;
  
  newMessage = '';
  conversationId: string | null = null;
  upcomingSession: ConsultationSession | null = null;
  chatBoxBottom = 80;
  
  private sub?: Subscription;
  private notifSub?: Subscription;
  private chatSub?: Subscription;

  constructor(
    private chatService: ChatService,
    private friendService: FriendService,
    private userContext: UserContextService,
    private notificationService: NotificationService,
    private sessionService: SessionService,
    private router: Router
  ) {}

  ngOnInit() {
    this.loadRecentChats();

    this.notifSub = this.notificationService.notifications$.subscribe((notifs) => {
      const chatNotifs = notifs.filter((n) => !n.isRead && n.title === 'New message');
      this.unreadTotal = chatNotifs.length;

      const map: { [profileId: string]: number } = {};
      
      chatNotifs.forEach((n) => {
        let fromProfileId = '';
        if (n.metadata && n.metadata.fromId) {
          fromProfileId = n.metadata.fromId;
        } else if (typeof n.metadata === 'string') {
          try {
            const parsed = JSON.parse(n.metadata);
            fromProfileId = parsed.fromId || parsed.profileId;
          } catch (e) {}
        }
        
        if (fromProfileId) {
          map[fromProfileId] = (map[fromProfileId] || 0) + 1;
          
          // Only reload if the message is from someone ELSE or if we don't have an active connection
          if (this.activeChatFriend && fromProfileId === this.activeChatFriend.profileId && !this.chatSub) {
             this.loadCompactMessages(this.activeChatFriend.profileId);
          }
        }
      });
      
      this.unreadMap = map;
    });
  }

  ngOnDestroy() {
    this.sub?.unsubscribe();
    this.notifSub?.unsubscribe();
    this.chatSub?.unsubscribe();
    this.chatService.disconnect();
  }

  toggleMenu() {
    this.isOpen = !this.isOpen;
    if (!this.isOpen) {
      this.activeChatFriend = null;
    }
    if (this.isOpen && this.recentFriends.length === 0) {
      this.loadRecentChats();
    }
    if (this.isOpen) {
      this.shouldScrollToBottom = true;
    }
  }

  openChat(friend: Friend) {
    if (this.activeChatFriend?.profileId === friend.profileId) {
      this.activeChatFriend = null;
      return;
    }

    this.activeChatFriend = friend;
    this.chatBoxBottom = 80; // Fixed position at bottom
    this.fetchUpcomingSession(friend.profileId);
    this.loadCompactMessages(friend.profileId);
  }

  fetchUpcomingSession(friendId: string) {
    const userId = this.userContext.getUser()?.profileId;
    if (!userId) return;

    this.sessionService.getAllSessions().subscribe((sessions: ConsultationSession[]) => {
      const now = new Date();
      this.upcomingSession = sessions.find((s: ConsultationSession) => {
        const involvesBoth = (s.client_id === userId && s.therapist_id === friendId) || 
                            (s.client_id === friendId && s.therapist_id === userId);
        const isFuture = new Date(s.start_time) > now;
        const isActive = s.status === 'CONFIRMED' || s.status === 'PENDING_PAYMENT';
        return involvesBoth && isFuture && isActive;
      }) || null;
    });
  }

  closeCompactChat() {
    this.activeChatFriend = null;
    this.chatSub?.unsubscribe();
    this.chatService.disconnect();
  }

  loadCompactMessages(friendId: string) {
    const userId = this.userContext.getUser()?.profileId;
    if (!userId) return;

    this.chatService.getOrCreateConversation(userId, friendId).subscribe(({ id }) => {
      this.conversationId = id;
      if (!id) return;

      this.chatSub?.unsubscribe();
      this.chatService.disconnect();

      this.chatService.getChatsByConversation(id, userId, 25).subscribe(msgs => {
        // Correct chronological order: latest messages at the bottom
        this.allCompactMessages = msgs;
        this.viewOffset = 0;
        this.updateViewWindow();
        this.shouldScrollToBottom = true;
      });

      this.chatSub = this.chatService.connect(friendId, id).subscribe((msg: Message) => {
        if (msg.conversationId === this.conversationId) {
          // 1. Handle "Echo" of own message (isMine)
          if (msg.isMine) {
            const tempIndex = this.allCompactMessages.findIndex(m => 
              m.senderId === msg.senderId && 
              m.id.startsWith('temp-') && 
              m.content === msg.content
            );

            if (tempIndex !== -1) {
              // Update optimistic message with real server data
              this.allCompactMessages[tempIndex].id = msg.id;
              if (msg.timestamp) this.allCompactMessages[tempIndex].timestamp = msg.timestamp;
              return; // Done, no need to push
            }
          }

          // 2. Standard Duplication Check (for incoming or already echoed messages)
          const isDuplicate = this.allCompactMessages.some(m => 
            m.id === msg.id || 
            (m.senderId === msg.senderId && 
             m.content === msg.content && 
             Math.abs(new Date(m.timestamp).getTime() - new Date(msg.timestamp).getTime()) < 2000)
          );

          if (!isDuplicate) {
            this.allCompactMessages.push(msg);
            if (this.viewOffset === 0) {
              this.updateViewWindow();
              this.shouldScrollToBottom = true;
            }
          }
        }
      });
    });
  }

  updateViewWindow() {
    const end = this.allCompactMessages.length - this.viewOffset;
    const start = Math.max(0, end - this.viewWindowSize);
    this.compactMessages = this.allCompactMessages.slice(start, end);
  }

  scrollUp() {
    if (this.viewOffset + 1 <= this.allCompactMessages.length - this.viewWindowSize) {
      this.viewOffset++;
      this.updateViewWindow();
    }
  }

  scrollDown() {
    if (this.viewOffset > 0) {
      this.viewOffset--;
      this.updateViewWindow();
    }
  }

  sendQuickMessage() {
    if (!this.newMessage.trim() || !this.conversationId) return;
    
    const userId = this.userContext.getUser()?.profileId;
    if (!userId) return;

    const userProfile = this.userContext.getUser();
    if (!userProfile) return;

    const msgContent = this.newMessage.trim();
    const tempMsg: Message = {
      id: 'temp-' + Date.now(),
      conversationId: this.conversationId!,
      content: msgContent,
      senderId: userProfile.profileId,
      userName: userProfile.firstName + ' ' + userProfile.lastName,
      userAvatar: userProfile.avatarUri || '',
      timestamp: new Date(),
      isMine: true
    };

    // Optimistic update
    this.allCompactMessages.push(tempMsg);
    if (this.viewOffset === 0) {
      this.updateViewWindow();
      this.shouldScrollToBottom = true;
    }

    this.chatService.sendMessage({
      conversationId: this.conversationId,
      content: msgContent,
      senderId: userId
    });

    this.newMessage = '';
  }

  onKeydown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      this.sendQuickMessage();
    }
  }

  openSearch() {
    this.isOpen = false;
    this.activeChatFriend = null;
    this.router.navigate(['/feature/chat']);
  }

  private loadRecentChats() {
    const userId = this.userContext.getUser()?.profileId;
    if (!userId) return;

    this.sub = this.chatService.getRecentConversations(userId).subscribe((conversations: any[]) => {
      const unknownIds = new Set<string>();
      for (const conv of conversations) {
        for (const pId of conv.participants || []) {
          if (pId !== userId) {
            unknownIds.add(pId);
          }
        }
      }

      if (unknownIds.size > 0) {
        const sortedIds = Array.from(unknownIds);
        this.friendService.getProfilesBatch(sortedIds).subscribe((friends: Friend[]) => {
          // Maintaining the order from conversations/me
          this.recentFriends = sortedIds
            .map(id => friends.find(f => f.profileId === id))
            .filter((f): f is Friend => !!f)
            .slice(0, 4);
          this.recalculateUnreadMap();
        });
      }
    });
  }

  private recalculateUnreadMap() {
    // Manually trigger the observable logic if needed by re-evaluating the current notifs
    this.notificationService.fetchNotifications(this.userContext.getUser()?.profileId || '');
  }

  ngAfterViewChecked() {
    if (this.shouldScrollToBottom) {
      this.scrollToBottom();
      this.shouldScrollToBottom = false;
    }
  }

  private scrollToBottom(): void {
    if (this.scrollViewport) {
      try {
        const element = this.scrollViewport.nativeElement;
        element.scrollTop = element.scrollHeight;
      } catch (err) {}
    }
  }
}
