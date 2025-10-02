import { CommonModule } from '@angular/common';
import { Component, input, output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Chat } from '../../../models/chat.models';

@Component({
  selector: 'app-chat-list',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: `./chat-list.html`,
  styleUrl: './chat-list.scss',
})
export class ChatListComponent {
  currentUser = input.required<{
    id: string;
    name: string;
    avatar: string;
    isOnline: boolean;
    status: string;
  }>();
  chats = input.required<Chat[]>();
  chatSelected = output<Chat>();
  searchQuery = '';
  visibleCount = 4;
  get limitedChats(): Chat[] {
    return this.chats().slice(0, this.visibleCount);
  }
  loadMoreChats() {
    if (this.visibleCount < this.chats().length) {
      this.visibleCount += 4;
    }
  }
  findChatById(chatId: string): Chat | undefined {
    return this.chats().find((chat) => chat.id === chatId);
  }
}
