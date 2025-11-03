import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription, combineLatest, finalize } from 'rxjs';
import { ChatListComponent } from '../../../components/chat/chat-list/chat-list';
import { ChatMainComponent } from '../../../components/chat/chat-main/chat-main';
import { Friend, Message } from '../../../models/chat.models';
import { ChatUser, mapUserProfileToChatUser } from '../../../models/map.utils';
import { ChatService } from '../../../services/chat/chat.service';
import { FriendService } from '../../../services/chat/profile.chat.service';
import { LoaderService } from '../../../services/loader/loader';

@Component({
  selector: 'chat-page',
  standalone: true,
  styleUrls: ['./chatpage.scss'],
  imports: [ChatListComponent, ChatMainComponent, CommonModule, FormsModule],
  template: `
    <div class="chat-layout">
      <div class="chat-sidebar">
        <app-chat-list
          [currentUser]="currentUser"
          [friends]="friends"
          [selectedFriendId]="selectedFriend()?.profileId || null"
          (friendSelected)="onFriendSelected($event, true)"
        ></app-chat-list>
      </div>
      <app-chat-main
        [selectedFriend]="selectedFriend()"
        [messages]="messages()"
        [conversationId]="conversationId()"
        [currentUserId]="currentUser.profileId"
        [currentUserName]="currentUser.name"
        [currentUserAvatar]="currentUser.avatar"
      ></app-chat-main>
    </div>
  `,
})
export class ChatComponent implements OnInit, OnDestroy {
  private sub?: Subscription;
  messageText = '';
  conversationId = signal<string>('');
  currentUser: ChatUser = mapUserProfileToChatUser(null);
  selectedFriend = signal<Friend | null>(null);
  messages = signal<Message[]>([]);
  friends: Friend[] = [];
  isLoading = false;

  constructor(
    private friendService: FriendService,
    private chatService: ChatService,
    private loader: LoaderService,
    private route: ActivatedRoute,
    private router: Router
  ) {}

  ngOnInit(): void {
    combineLatest([
      this.chatService.getCurrentUser(),
      this.friendService.getMyFriends(),
      this.route.paramMap,
    ]).subscribe(([profile, friends, params]) => {
      this.currentUser = mapUserProfileToChatUser(profile);
      this.friends = friends;
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
    });
  }
  ngOnDestroy(): void {
    this.sub?.unsubscribe();
    this.chatService.disconnect();
  }

  onFriendSelected(friend: Friend, updateUrl = false) {
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

        this.chatService.connect(friend.profileId, conv.id).subscribe({
          next: (msg) => {
            const enriched = this.enrichMessage(msg, friend);
            this.messages.update((prev) => [...prev, enriched]);
          },
          error: (err) => console.error('Connection error:', err),
        });

        this.loader.show();
        this.chatService
          .getChatsByConversation(conv.id, this.currentUser.profileId)
          .pipe(finalize(() => this.loader.hide()))
          .subscribe((msgs) => {
            const enrichedMsgs = msgs.map((m) => this.enrichMessage(m, friend));
            this.messages.set(enrichedMsgs);
          });
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
}
