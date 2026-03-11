import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { RouterModule, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Subscription } from 'rxjs';
import { Friend } from '../../../models/chat.models';
import { ChatService } from '../../../services/chat/chat.service';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { ImgFallbackDirective } from '../../../shared/directives/img-fallback.directive';
import { UserContextService } from '../../../services/profile/profile-service';
import { NotificationService } from '../../../services/notification/notification.service';

@Component({
  selector: 'app-floating-chat',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule, AvatarFallbackPipe, ImgFallbackDirective],
  templateUrl: './floating-chat.html',
  styleUrls: ['./floating-chat.scss']
})
export class FloatingChatComponent implements OnInit, OnDestroy {
  isOpen = false;
  recentFriends: Friend[] = [];
  unreadTotal = 0;
  unreadMap: { [profileId: string]: number } = {};
  
  private sub?: Subscription;
  private notifSub?: Subscription;

  constructor(
    private chatService: ChatService,
    private friendService: FriendService,
    private userContext: UserContextService,
    private notificationService: NotificationService,
    private router: Router
  ) {}

  ngOnInit() {
    this.loadRecentChats();

    this.notifSub = this.notificationService.notifications$.subscribe((notifs) => {
      // "New message" is the title sent by the backend for chat notifications
      const chatNotifs = notifs.filter((n) => !n.isRead && n.title === 'New message');
      this.unreadTotal = chatNotifs.length;

      const map: { [profileId: string]: number } = {};
      
      chatNotifs.forEach((n) => {
        let fromName = '';
        if (n.metadata && n.metadata.from) {
          fromName = n.metadata.from;
        } else if (typeof n.metadata === 'string') {
          try {
            const parsed = JSON.parse(n.metadata);
            fromName = parsed.from;
          } catch (e) {}
        }
        
        const friend = this.recentFriends.find(f => {
          const fullName = `${f.firstName} ${f.lastName}`.trim();
          return fullName === fromName;
        });
        
        if (friend && friend.profileId) {
          map[friend.profileId] = (map[friend.profileId] || 0) + 1;
        }
      });
      
      this.unreadMap = map;
    });
  }

  ngOnDestroy() {
    this.sub?.unsubscribe();
    this.notifSub?.unsubscribe();
  }

  toggleMenu() {
    this.isOpen = !this.isOpen;
    if (this.isOpen && this.recentFriends.length === 0) {
      this.loadRecentChats();
    }
  }

  openChat(profileId: string) {
    this.isOpen = false;
    this.router.navigate(['/feature/chat', profileId]);
  }

  openSearch() {
    this.isOpen = false;
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
        // Fetch profiles for these recent connections
        this.friendService.getProfilesBatch(Array.from(unknownIds)).subscribe((friends: Friend[]) => {
          // Sort by conversation time if possible, or just take first 3-4
          this.recentFriends = friends.slice(0, 4);
          // Re-calculate unread map in case new friends loaded
          this.recalculateUnreadMap();
        });
      }
    });
  }

  private recalculateUnreadMap() {
    // Manually trigger the observable logic if needed by re-evaluating the current notifs
    this.notificationService.fetchNotifications(this.userContext.getUser()?.profileId || '');
  }
}
