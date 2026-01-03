import { CommonModule } from '@angular/common';
import { Component, HostBinding, Input } from '@angular/core';
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

  @HostBinding('class.first-in-group')
  get isFirstInGroup() {
    return this.message?.isFirstInGroup;
  }

  @HostBinding('class.last-in-group')
  get isLastInGroup() {
    return this.message?.isLastInGroup;
  }

  hasReaction = true;

  formatTime(date: any): string {
    if (!date) return '';

    const d = date instanceof Date ? date : new Date(date);
    if (isNaN(d.getTime())) return '';

    const now = new Date();
    const isToday =
      d.getDate() === now.getDate() &&
      d.getMonth() === now.getMonth() &&
      d.getFullYear() === now.getFullYear();

    const hours = d.getHours();
    const minutes = d.getMinutes();
    const ampm = hours >= 12 ? 'PM' : 'AM';
    const displayHours = hours % 12 || 12;
    const displayMinutes = minutes.toString().padStart(2, '0');
    const timeString = `${displayHours}:${displayMinutes} ${ampm}`;

    if (isToday) return timeString;

    const dateString = `${d.getDate().toString().padStart(2, '0')}/${(
      d.getMonth() + 1
    )
      .toString()
      .padStart(2, '0')}/${d.getFullYear()}`;

    return `${dateString}, ${timeString}`;
  }
}
