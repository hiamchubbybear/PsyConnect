import { CommonModule } from '@angular/common';
import { Component, input } from '@angular/core';
import { Message } from '../../../models/chat.models';

@Component({
  selector: 'app-message-bubble',
  standalone: true,
  imports: [CommonModule],
  templateUrl: `./chat-bubble.html`,
  styleUrl: `./chat-bubble.scss`,
})
export class MessageBubbleComponent {
  message = input.required<Message>();
  hasReaction = true;

  formatTime(date: Date): string {
    const hours = date.getHours();
    const minutes = date.getMinutes();
    const ampm = hours >= 12 ? 'PM' : 'AM';
    const displayHours = hours % 12 || 12;
    const displayMinutes = minutes.toString().padStart(2, '0');
    return `${displayHours}:${displayMinutes} ${ampm}`;
  }
}
