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
import { ChatInputComponent } from '../../../shared/ui-atoms/chat-input/chat-input.component';
import { MessageBubbleComponent } from '../../../shared/ui-atoms/message-bubble/message-bubble.component';
import {
  CallOptionsMenuComponent,
  CallType,
} from '../../call-options-menu/call-options-menu.component';

@Component({
  selector: 'app-chat-main',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MessageBubbleComponent,
    ChatInputComponent,
    CallOptionsMenuComponent,
    TranslateModule,
  ],
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
  @Input() isLoadingMessages = false;
  @Output() messagesChange = new EventEmitter<Message[]>();
  @Output() callRequested = new EventEmitter<CallType>();
  @Output() back = new EventEmitter<void>();
  @ViewChild('messagesContainer')
  messagesContainer!: ElementRef<HTMLDivElement>;
  showLoadOlderButton = false;

  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;

  isLoadingOld = false;

  isTyping = false;
  isSending = false; // Guard against duplicate sends
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
  onSendMessage(text: string) {
    const friend = this.selectedFriend;
    if (!text || !friend) return;

    if (!this.conversationId || !this.currentUserId) {
      console.warn('Missing conversationId or currentUserId');
      return;
    }

    // Guard against multiple sends
    if (this.isSending) {
      console.warn('⚠️ Already sending, ignoring duplicate');
      return;
    }

    this.isSending = true;

    this.chatService.sendMessage({
      conversationId: this.conversationId,
      content: text,
      senderId: this.currentUserId,
    });

    // Reset guard after delay
    setTimeout(() => {
      this.isSending = false;
    }, 500);
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
        oldestTime,
      )
      .pipe(finalize(() => (this.isLoadingOld = false)))
      .subscribe((olderMsgs) => {
        if (!olderMsgs.length) return;

        const enriched = olderMsgs.map((m) => this.enrichMessage(m, friend!));

        // Merge and re-group all messages
        const merged = [...enriched, ...this.messages];
        const grouped = this.groupMessages(merged);

        this.messages = grouped;
        this.messagesChange.emit(grouped);
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

  private groupMessages(messages: Message[]): Message[] {
    const GROUPING_THRESHOLD = 5 * 60 * 1000; // 5 minutes

    return messages.map((msg, index) => {
      const prev = messages[index - 1];
      const next = messages[index + 1];

      const isFirstInGroup =
        !prev ||
        prev.senderId !== msg.senderId ||
        msg.timestamp.getTime() - prev.timestamp.getTime() > GROUPING_THRESHOLD;

      const isLastInGroup =
        !next ||
        next.senderId !== msg.senderId ||
        next.timestamp.getTime() - msg.timestamp.getTime() > GROUPING_THRESHOLD;

      return {
        ...msg,
        isFirstInGroup,
        isLastInGroup,
        showTimestamp: isFirstInGroup || isLastInGroup,
      };
    });
  }

  startCall(callType: CallType) {
    console.log('📞 Starting', callType, 'call...');
    this.callRequested.emit(callType);
  }
}
