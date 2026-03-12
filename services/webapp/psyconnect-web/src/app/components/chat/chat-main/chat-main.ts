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
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { finalize } from 'rxjs';
import { Friend, Message } from '../../../models/chat.models';
import { ConsultationSession } from '../../../models/consultation.model';
import { ChatService } from '../../../services/chat/chat.service';
import { SessionService } from '../../../services/consultation/session.service';
import { ImgFallbackDirective } from '../../../shared/directives/img-fallback.directive';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';
import { ChatInputComponent } from '../../../shared/ui-atoms/chat-input/chat-input.component';
import { MessageBubbleComponent } from '../../../shared/ui-atoms/message-bubble/message-bubble.component';
import { CallType } from '../../call-options-menu/call-options-menu.component';

@Component({
  selector: 'app-chat-main',
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    MessageBubbleComponent,
    ChatInputComponent,
    TranslateModule,
    AvatarFallbackPipe,
    ImgFallbackDirective,
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
  @Input() suggestedUsers: Friend[] = [];
  @Output() messagesChange = new EventEmitter<Message[]>();
  @Output() callRequested = new EventEmitter<CallType>();
  @Output() suggestionSelected = new EventEmitter<Friend>();
  @Output() back = new EventEmitter<void>();
  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;

  isLoadingOld = false;
  showLoadOlderButton = false;
  showScrollBottomButton = false;

  get hasUserMessages(): boolean {
    return this.messages.some((m) => !m.isSystem);
  }

  @Input() isTyping = false;
  isSending = false;
  private autoScrollPending = false;

  upcomingSession: ConsultationSession | null = null;

  constructor(
    private chatService: ChatService,
    private sessionService: SessionService,
    private dialog: MatDialog,
    private router: Router,
  ) {}

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

    if (changes['selectedFriend'] || changes['conversationId']) {
      this.fetchUpcomingSession();
    }
  }

  get isSessionLive(): boolean {
    if (!this.upcomingSession) return false;
    const now = new Date();
    const start = new Date(this.upcomingSession.start_time);
    const end = new Date(this.upcomingSession.end_time);
    // Active if started and not ended (with 15min early join buffer)
    const earlyJoinLimit = new Date(start.getTime() - 15 * 60 * 1000);
    return now >= earlyJoinLimit && now <= end;
  }

  private fetchUpcomingSession() {
    if (!this.selectedFriend) {
      this.upcomingSession = null;
      return;
    }

    this.sessionService.getAllSessions().subscribe((sessions) => {
      const friendId = this.selectedFriend?.profileId;
      const myId = this.currentUserId;

      const now = new Date();
      // Find the closest upcoming or current session
      const relevantSessions = sessions
        .filter((s) => {
          const involveMe = s.client_id === myId || s.therapist_id === myId;
          const involveFriend =
            s.client_id === friendId || s.therapist_id === friendId;
          const isParticipant = involveMe && involveFriend;

          const isValidStatus =
            s.status === 'PENDING_PAYMENT' || s.status === 'CONFIRMED';
          const endTime = new Date(s.end_time);

          return isParticipant && isValidStatus && endTime > now;
        })
        .sort(
          (a, b) =>
            new Date(a.start_time).getTime() - new Date(b.start_time).getTime(),
        );

      this.upcomingSession = relevantSessions[0] || null;
    });
  }

  joinSession() {
    if (!this.upcomingSession || !this.isSessionLive) return;
    if (this.upcomingSession.mode?.toUpperCase() === 'ONLINE') {
      this.callRequested.emit('video');
    }
  }
  onScroll() {
    if (!this.scrollContainer) return;
    const container = this.scrollContainer.nativeElement;

    const topThreshold = 50;
    this.showLoadOlderButton = container.scrollTop <= topThreshold;

    const bottomDistance =
      container.scrollHeight - container.scrollTop - container.clientHeight;
    this.showScrollBottomButton = bottomDistance > 200;
  }
  onSendMessage(text: string) {
    const friend = this.selectedFriend;
    if (!text || !text.trim() || !friend) return;
    text = text.trim();

    if (!this.conversationId || !this.currentUserId) {
      console.warn('Missing conversationId or currentUserId');
      return;
    }

    if (this.isSending) {
      console.warn('⚠️ Already sending, ignoring duplicate');
      return;
    }

    this.isSending = true;

    const tempId = `temp-${Date.now()}`;
    const optimisticMsg: Message = {
      id: tempId,
      conversationId: this.conversationId,
      senderId: this.currentUserId,
      userName: this.currentUserName,
      userAvatar: this.currentUserAvatar,
      content: text,
      timestamp: new Date(),
      isMine: true,
      status: 'sending',
    };

    this.messages = this.groupMessages([...this.messages, optimisticMsg]);
    this.messagesChange.emit(this.messages);

    try {
      this.chatService.sendMessage({
        conversationId: this.conversationId,
        content: text,
        senderId: this.currentUserId,
      });

      setTimeout(() => {
        this.messages = this.messages.map((m) =>
          m.id === tempId ? { ...m, status: 'sent' as const } : m,
        );
        this.messagesChange.emit(this.messages);

        setTimeout(() => {
          this.messages = this.messages.map((m) =>
            m.id === tempId ? { ...m, status: 'delivered' as const } : m,
          );
          this.messagesChange.emit(this.messages);
        }, 1200);
      }, 300);
    } catch (err) {
      console.error('❌ Failed to send message:', err);

      this.messages = this.messages.map((m) =>
        m.id === tempId ? { ...m, status: 'failed' as const } : m,
      );
      this.messagesChange.emit(this.messages);
    }

    setTimeout(() => {
      this.isSending = false;
    }, 500);
  }

  retryMessage(failedMsg: Message) {
    if (!this.conversationId || !this.currentUserId) return;

    this.messages = this.messages.map((m) =>
      m.id === failedMsg.id ? { ...m, status: 'sending' as const } : m,
    );
    this.messagesChange.emit(this.messages);

    try {
      this.chatService.sendMessage({
        conversationId: this.conversationId,
        content: failedMsg.content,
        senderId: this.currentUserId,
      });

      setTimeout(() => {
        this.messages = this.messages.map((m) =>
          m.id === failedMsg.id ? { ...m, status: 'sent' as const } : m,
        );
        this.messagesChange.emit(this.messages);
      }, 300);
    } catch (err) {
      this.messages = this.messages.map((m) =>
        m.id === failedMsg.id ? { ...m, status: 'failed' as const } : m,
      );
      this.messagesChange.emit(this.messages);
    }
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
        ? this.currentUserAvatar || ''
        : friend.avatarUri || '',
    };
  }

  scrollToBottom() {
    if (!this.scrollContainer) return;
    try {
      const el = this.scrollContainer.nativeElement;
      el.scrollTo({ top: el.scrollHeight, behavior: 'smooth' });
      this.autoScrollPending = false;
    } catch (err) {
      console.warn('scrollToBottom error:', err);
    }
  }
  trackByMessageId(index: number, message: any): string {
    return message.id || index;
  }

  private groupMessages(messages: Message[]): Message[] {
    const GROUPING_THRESHOLD = 5 * 60 * 1000;

    return messages.map((msg, index) => {
      const prev = messages[index - 1];
      const next = messages[index + 1];

      let showDateSeparator = false;
      let dateSeparatorText = '';
      if (!prev) {
        showDateSeparator = true;
      } else {
        const prevDate = new Date(prev.timestamp);
        const currDate = new Date(msg.timestamp);
        if (
          prevDate.getDate() !== currDate.getDate() ||
          prevDate.getMonth() !== currDate.getMonth() ||
          prevDate.getFullYear() !== currDate.getFullYear()
        ) {
          showDateSeparator = true;
        }
      }

      if (showDateSeparator) {
        const date = new Date(msg.timestamp);
        const today = new Date();
        const yesterday = new Date(today);
        yesterday.setDate(yesterday.getDate() - 1);

        if (date.toDateString() === today.toDateString()) {
          dateSeparatorText = 'Today';
        } else if (date.toDateString() === yesterday.toDateString()) {
          dateSeparatorText = 'Yesterday';
        } else {
          dateSeparatorText = date.toLocaleDateString('en-GB', {
            day: 'numeric',
            month: 'short',
            year:
              date.getFullYear() !== today.getFullYear()
                ? 'numeric'
                : undefined,
          });
        }
      }

      const isFirstInGroup =
        !prev ||
        prev.senderId !== msg.senderId ||
        showDateSeparator ||
        msg.timestamp.getTime() - prev.timestamp.getTime() > GROUPING_THRESHOLD;

      const isLastInGroup =
        !next ||
        next.senderId !== msg.senderId ||
        (() => {
          const currDate = new Date(msg.timestamp);
          const nextDate = new Date(next.timestamp);
          return (
            currDate.getDate() !== nextDate.getDate() ||
            currDate.getMonth() !== nextDate.getMonth() ||
            currDate.getFullYear() !== nextDate.getFullYear()
          );
        })() ||
        next.timestamp.getTime() - msg.timestamp.getTime() > GROUPING_THRESHOLD;

      return {
        ...msg,
        isFirstInGroup,
        isLastInGroup,
        showTimestamp: isFirstInGroup || isLastInGroup,
        showDateSeparator,
        dateSeparatorText,
      };
    });
  }

  startCall(callType: CallType) {
    console.log('📞 Starting', callType, 'call...');
    this.callRequested.emit(callType);
  }

  bookSession() {
    if (!this.selectedFriend) return;
    this.router.navigate([
      '/feature/consultation/book',
      this.selectedFriend.profileId,
    ]);
  }
}
