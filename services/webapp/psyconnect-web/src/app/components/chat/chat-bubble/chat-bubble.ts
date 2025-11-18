import { CommonModule } from '@angular/common';
import { Component, Input, SimpleChanges } from '@angular/core';
import { Message } from '../../../models/chat.models';

@Component({
  selector: 'app-message-bubble',
  standalone: true,
  imports: [CommonModule],
  templateUrl: `./chat-bubble.html`,
  styleUrl: `./chat-bubble.scss`,
})
export class MessageBubbleComponent {
  @Input() message!: Message;
  hasReaction = true;
  ngOnChanges(changes: SimpleChanges) {
    if (changes['message']) {
      console.log('[MessageBubble] New message input:', this.message);

      if (this.message) {
        console.log('id:', this.message.id, typeof this.message.id);
        console.log(
          'conversationId:',
          this.message.conversationId,
          typeof this.message.conversationId
        );
        console.log('userId:', this.message.id, typeof this.message.id);
        console.log(
          'userName:',
          this.message.userName,
          typeof this.message.userName
        );
        console.log(
          'userAvatar:',
          this.message.userAvatar,
          typeof this.message.userAvatar
        );
        console.log(
          'content:',
          this.message.content,
          typeof this.message.content
        );
        if (this.message.timestamp instanceof Date) {
          console.log(
            'timestamp:',
            this.message.timestamp.toISOString(),
            ' Date object'
          );
    } else {
          console.warn(
            'timestamp is not a Date object, converting...',
            this.message.timestamp
          );
          this.message.timestamp = new Date(this.message.timestamp);
        }

        console.log('isMine:', this.message.isMine, typeof this.message.isMine);
      } else {
        console.warn('[MessageBubble] message is null or undefined');
      }
    }
  }
  formatTime(date: any): string {
    if (!date) return '';

    const d = date instanceof Date ? date : new Date(date);

    if (isNaN(d.getTime())) return '';

    const hours = d.getHours();
    const minutes = d.getMinutes();
    const ampm = hours >= 12 ? 'PM' : 'AM';
    const displayHours = hours % 12 || 12;
    const displayMinutes = minutes.toString().padStart(2, '0');
    return `${displayHours}:${displayMinutes} ${ampm}`;
  }
}
