import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Subscription } from 'rxjs';
import { fakeMessages } from '../../../../../src/app/models/chat.models';
import { ChatListComponent } from '../../../components/chat/chat-list/chat-list';
import { ChatMainComponent } from '../../../components/chat/chat-main/chat-main';
import { Chat, Message } from '../../../models/chat.models';
import { ChatUser, mapUserProfileToChatUser } from '../../../models/map.utils';
import { ChatService } from '../../../services/chat/chat.service';

@Component({
  selector: 'chat-page',
  styleUrls: ['./chatpage.scss'],
  imports: [ChatListComponent, ChatMainComponent, CommonModule, FormsModule],
  standalone: true,
  template: `
    <div class="chat-layout">
      <div class="chat-sidebar">
        <app-chat-list
          [currentUser]="currentUser"
          [chats]="chats"
          (chatSelected)="onChatSelected($event)"
        ></app-chat-list>
      </div>
      <app-chat-main [selectedChat]="selectedChat()" [messages]="messages()">
      </app-chat-main>
    </div>
  `,
})
export class ChatComponent implements OnInit, OnDestroy {
  private sub?: Subscription;

  currentUser: ChatUser = mapUserProfileToChatUser(null);

  constructor(private chatService: ChatService) {}
  fakeMessages = fakeMessages;

  ngOnInit(): void {
    this.sub = this.chatService.getCurrentUser().subscribe((profile) => {
      this.currentUser = mapUserProfileToChatUser(profile);
    });
  }

  ngOnDestroy(): void {
    this.sub?.unsubscribe();
  }

  chats: Chat[] = [
    {
      id: 'c-001',
      title: 'Alice',
      avatar:
        'https://i.pinimg.com/1200x/8d/15/fb/8d15fbf1a92360fd84e30e65e8af9136.jpg',
      lastMessage: 'Hey, how are you?',
      timestamp: '10:30 AM',
      isTyping: false,
    },
    {
      id: 'c-002',
      title: 'Bob',
      avatar:
        'https://i.pinimg.com/736x/ef/e7/20/efe7207830144f9a582e70d1b9643dcb.jpg',
      lastMessage: 'Did you check the docs?',
      timestamp: '09:45 AM',
      isTyping: true,
    },
  ];

  selectedChat = signal<Chat | null>(this.chats[0] ?? null);

  messages = signal<Message[]>([
    {
      id: 'm-001',
      userId: 'u-001',
      userName: 'Alice',
      userAvatar: 'https://i.pinimg.com/1200x/8d/15/fb/8d15fbf1a92360fd84e30e65e8af9136.jpg',
      content: 'Hi there 👋',
      timestamp: new Date('2025-09-22T10:15:00'),
      isMine: false,
    },
    {
      id: 'm-002',
      userId: 'me',
      userName: 'You',
      userAvatar: 'assets/avatars/me.png',
      content: 'Hey Alice! I’m good, how about you?',
      timestamp: new Date('2025-09-22T10:16:00'),
      isMine: true,
    },
    {
      id: 'm-003',
      userId: 'u-001',
      userName: 'Alice',
      userAvatar: 'https://i.pinimg.com/1200x/8d/15/fb/8d15fbf1a92360fd84e30e65e8af9136.jpg',
      content: 'I’m fine too. Wanna grab coffee later? ☕',
      timestamp: new Date('2025-09-22T10:18:00'),
      isMine: false,
    },
  ]);
  onChatSelected(chat: Chat) {
    this.selectedChat.set(chat);
  }
}
