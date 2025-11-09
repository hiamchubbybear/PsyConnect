import { CommonModule } from '@angular/common';
import {
    AfterViewInit,
    Component,
    ElementRef,
    EventEmitter,
    Input,
    OnChanges,
    Output,
    SimpleChanges,
    ViewChild,
} from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { finalize } from 'rxjs';
import { Friend, Message } from '../../../models/chat.models';
import { ChatService } from '../../../services/chat/chat.service';
import { MessageBubbleComponent } from '../chat-bubble/chat-bubble';

@Component({
  selector: 'app-chat-main',
  standalone: true,
  imports: [CommonModule, FormsModule, MessageBubbleComponent, TranslateModule],
  templateUrl: './chat-main.html',
  styleUrls: ['./chat-main.scss'],
})
export class ChatMainComponent implements AfterViewInit, OnChanges {
  @Input() selectedFriend: Friend | null = null;
  @Input() messages: Message[] = [];
  @Input() conversationId: string | null = null;
  @Input() currentUserId: string = '';
  @Input() currentUserName: string = '';
  @Input() currentUserAvatar: string = '';
  @Output() messagesChange = new EventEmitter<Message[]>();
  @ViewChild('messagesContainer')
  messagesContainer!: ElementRef<HTMLDivElement>;
  showLoadOlderButton = false;

  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;

  isLoadingOld = false;
  messageText = '';
  isTyping = false;
  private autoScrollPending = false;

  constructor(private chatService: ChatService) {}

  ngAfterViewInit(): void {
    this.scrollToBottom();
  }

  ngOnChanges(changes: SimpleChanges) {
    if (changes['messages'] && !this.isLoadingOld) {
      this.autoScrollPending = true;
      queueMicrotask(() => {
        if (this.autoScrollPending) this.scrollToBottom();
      });
    }
  }
  onScroll() {
    const container = this.messagesContainer.nativeElement;
    const threshold = 50;
    this.showLoadOlderButton = container.scrollTop <= threshold;
  }
  onSendMessage() {
    const friend = this.selectedFriend;
    const text = this.messageText.trim();
    if (!text || !friend) return;

    if (!this.conversationId || !this.currentUserId) {
      console.warn('Missing conversationId or currentUserId');
      return;
    }

    this.chatService.sendMessage({
      conversationId: this.conversationId,
      content: text,
      senderId: this.currentUserId,
    });

    this.messageText = '';
  }

  loadOlderMessages() {
    if (!this.conversationId || !this.currentUserId || !this.selectedFriend)
      return;

    this.isLoadingOld = true;
    const friend = this.selectedFriend;
    const oldest = this.messages[0];
    const oldestTime = oldest ? oldest.timestamp : undefined;

    this.chatService
      .getChatsByConversation(
        this.conversationId,
        this.currentUserId,
        10,
        oldestTime
      )
      .pipe(finalize(() => (this.isLoadingOld = false)))
      .subscribe((olderMsgs) => {
        if (!olderMsgs.length) return;

        const enriched = olderMsgs.map((m) => this.enrichMessage(m, friend!));

        const merged = [...enriched, ...this.messages];
        this.messages = merged;
        this.messagesChange.emit(merged);
      });
  }

  private enrichMessage(msg: Message, friend: Friend): Message {
    const isMine = msg.senderId === this.currentUserId;
    return {
      ...msg,
      isMine,
      userName: isMine
        ? this.currentUserName
        : `${friend.firstName} ${friend.lastName}`,
      userAvatar: isMine
        ? this.currentUserAvatar
        : friend.avatarUri || 'assets/avatars/default.png',
    };
  }

  private scrollToBottom() {
    if (!this.scrollContainer) return;
    try {
      const el = this.scrollContainer.nativeElement;
      el.scrollTop = el.scrollHeight;
      this.autoScrollPending = false;
    } catch (err) {
      console.warn('scrollToBottom error:', err);
    }
  }
  trackByMessageId(index: number, message: any): string {
    return message.id || index;
  }
}
