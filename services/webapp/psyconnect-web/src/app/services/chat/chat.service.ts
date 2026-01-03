import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, map, Observable, of, Subject, tap } from 'rxjs';
import { webSocket, WebSocketSubject } from 'rxjs/webSocket';
import { environment } from '../../../environments/environment';
import { SecureStorageService } from '../../encrypt/secure';
import { Chat, ChatFromApi, Message } from '../../models/chat.models';
import { ProfileResponse } from '../../models/profile';
import { Profile } from '../profile/profile';
import { UserContextService, UserProfile } from '../profile/profile-service';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  private readonly apiUrl = environment.apiUrl;
  private readonly wsUrl = environment.wsUrl;
  private socket$: WebSocketSubject<any> | null = null;
  private reconnectAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 3000;
  private heartbeatInterval?: any;

  private connectionState$ = new Subject<boolean>();

  constructor(
    private http: HttpClient,
    private userContext: UserContextService,
    private profileService: Profile,
    private secureStorage: SecureStorageService
  ) {}
  private readonly AccessTokenKey = environment.accessTokenKey;
  getCurrentUser(): Observable<UserProfile | null> {
    return this.profileService.getProfile().pipe(
      map((res) => res.data ?? null),
      catchError(() => of(null))
    );
  }

  getChatsByConversation(
    conversationId: string,
    currentUserId: string,
    limit = 10,
    before?: Date
  ): Observable<Message[]> {
    let params: any = { limit };
    if (before) params.before = before.toISOString();

    return this.http
      .get<{ data: ChatFromApi[] }>(
        `${this.apiUrl}/chats/conversation/${conversationId}`,
        { params }
      )
      .pipe(
        map((res) => {
          const rawChats = res.data ?? [];
          return rawChats.map((c, i) => {
            let parsedContent: string;
            try {
              parsedContent =
                typeof c.text === 'string'
                  ? JSON.parse(c.text)
                  : String(c.text);
            } catch {
              parsedContent = String(c.text);
            }

            const senderId = c.senderId || '';
            const isMine = senderId === currentUserId;

            const msg: Message = {
              id: c.id,
              conversationId: c.conversationId,
              senderId,
              userName: '',
              userAvatar: '',
              content: parsedContent,
              timestamp: new Date(c.createdAt),
              isMine,
            };
            return msg;
          });
        }),
        catchError((err) => {
          return of([] as Message[]);
        })
      );
  }

  createConversation(userIds: string[]): Observable<{ id: string }> {
    return this.http
      .post<{ data: { id: string } }>(`${this.apiUrl}/chats/conversations`, {
        userIds,
      })
      .pipe(map((res) => res.data));
  }
  createChat(chat: Partial<Chat>): Observable<Chat | null> {
    return this.http.post<{ data: Chat }>(this.apiUrl, chat).pipe(
      map((res) => res.data),
      catchError(() => of(null))
    );
  }

  deleteChat(id: string): Observable<boolean> {
    return this.http.delete(`${this.apiUrl}/${id}`).pipe(
      map(() => true),
      catchError(() => of(false))
    );
  }
  private isConnected = false;

  connect(receiverId: string, conversationId: string): Observable<Message> {
    const token = this.secureStorage.getItem(this.AccessTokenKey);
    const url = `${this.wsUrl}?token=${token}&receiver=${receiverId}&conversationId=${conversationId}`;

    console.log('%connecting to WS:', 'color: cyan', url);

    this.socket$ = webSocket<Message>({
      url,
      openObserver: {
        next: () => {
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.connectionState$.next(true);
          console.log('%c WebSocket connected', 'color: green');

          this.startHeartbeat();
        },
      },
      closeObserver: {
        next: () => {
          this.isConnected = false;
          this.connectionState$.next(false);
          console.warn('%c WebSocket closed', 'color: orange');
          this.stopHeartbeat();
          this.tryReconnect(receiverId, conversationId);
        },
      },
      deserializer: (e) => JSON.parse(e.data),
      serializer: (msg) => JSON.stringify(msg),
    });

    return this.socket$.asObservable().pipe(
      tap({
        error: (err) => console.error('ebSocket error:', err),
      })
    );
  }
  private tryReconnect(receiverId: string, conversationId: string) {
    if (this.reconnectAttempts >= this.maxReconnectAttempts) {
      console.error(' Max reconnect attempts reached. Stop retrying.');
      return;
    }

    this.reconnectAttempts++;
    const delayMs = this.reconnectDelay * this.reconnectAttempts;

    console.log(
      ` Attempting reconnect #${this.reconnectAttempts} after ${
        delayMs / 1000
      }s`
    );

    setTimeout(() => {
      this.connect(receiverId, conversationId).subscribe({
        next: (msg) => {
          console.log(' Reconnected and received:', msg);
        },
        error: (err) => {
          console.error('econnect failed:', err);
        },
      });
    }, delayMs);
  }
  private startHeartbeat() {}

  private stopHeartbeat() {}
  sendMessage(msg: {
    conversationId: string;
    content: string;
    senderId: string;
  }) {
    if (!this.socket$ || !this.isConnected) {
      console.warn('WebSocket not connected, message dropped');
      return;
    }

    const payload = {
      type: 'chat', // Message type for backend routing
      conversationId: msg.conversationId,
      senderId: msg.senderId,
      data: {
        text: msg.content,
      },
    };

    console.log('Sending chat message:', payload);
    this.socket$.next(payload);
  }

  disconnect(): void {
    this.socket$?.complete();
  }

  attachProfileToChat(chat: Chat): Observable<Chat & { user?: UserProfile }> {
    const localUser = this.userContext.getUser();
    if (localUser) {
      return of({ ...chat, user: localUser });
    }
    return this.profileService
      .getProfile()
      .pipe(map((res: ProfileResponse) => ({ ...chat, user: res.data })));
  }
  getOrCreateConversation(
    user1Id: string,
    user2Id: string
  ): Observable<{ id: string }> {
    console.log('getOrCreateConversation response~!!!:');
    return this.http
      .get<{ code: number; message: string; data: { id: string } }>(
        `${this.apiUrl}/conversations/by-users`,
        { params: { user1: user1Id, user2: user2Id } }
      )
      .pipe(
        tap((res) => console.log('getOrCreateConversation response:', res)),
        map((res) => {
          const conv = res.data;
          const id = conv.id || (conv as any)._id || '';
          if (!id) console.warn('Conversation ID missing in response', conv);
          return { id };
        }),
        catchError((err) => {
          console.error('getOrCreateConversation error:', err);
          return of({ id: '' });
        })
      );
  }

  attachProfilesToChatList(
    chats: Chat[]
  ): Observable<(Chat & { user?: UserProfile })[]> {
    return new Observable((observer) => {
      const enriched: (Chat & { user?: UserProfile })[] = [];
      let completed = 0;

      chats.forEach((chat, index) => {
        this.attachProfileToChat(chat).subscribe({
          next: (result) => {
            enriched[index] = result;
            completed++;
            if (completed === chats.length) {
              observer.next(enriched);
              observer.complete();
            }
          },
          error: (err) => observer.error(err),
        });
      });
    });
  }
}
