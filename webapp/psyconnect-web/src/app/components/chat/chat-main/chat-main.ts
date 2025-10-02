// 6. components/chat-main/chat-main.component.ts
import { CommonModule } from '@angular/common';
import { Component, effect, input, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { Chat, Message } from '../../../models/chat.models';
import { MessageBubbleComponent } from '../chat-bubble/chat-bubble';

@Component({
  selector: 'app-chat-main',
  standalone: true,
  imports: [CommonModule, FormsModule, MessageBubbleComponent, TranslateModule],
  templateUrl: `./chat-main.html`,
  styleUrl: `./chat-main.scss`,
})
export class ChatMainComponent  {
  selectedChat = input<Chat | null>();
  messages = input.required<Message[]>();
  messageText = '';
  isTyping = signal(true);

  constructor() {
    effect(() => {
      console.log('Messages changed:', this.messages());
    });
  }
  sendMessage() {
    if (this.messageText.trim()) {
      console.log('Sending message:', this.messageText);
      this.messageText = '';
    }
  }
}
