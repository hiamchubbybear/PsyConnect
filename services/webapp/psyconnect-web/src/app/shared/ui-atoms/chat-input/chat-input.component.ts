import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-chat-input',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './chat-input.component.html',
  styleUrls: ['./chat-input.component.scss'],
})
export class ChatInputComponent {
  @Input() placeholder = 'CHAT.EXAMPLE.placeholder';
  @Input() disabled = false;

  @Output() sendMessage = new EventEmitter<string>();
  @Output() attachFile = new EventEmitter<void>();
  @Output() addEmoji = new EventEmitter<void>();

  messageText = '';

  onSend() {
    const text = this.messageText.trim();
    if (!text || this.disabled) return;

    this.sendMessage.emit(text);
    this.messageText = '';
  }

  onKeyUpEnter(event: Event) {
    event.preventDefault();
    this.onSend();
  }
}
