import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { catchError, filter, map, Observable, of, Subject, tap } from 'rxjs';
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
    private secureStorage: SecureStorageService,
  ) {}
  private readonly AccessTokenKey = environment.accessTokenKey;
  getCurrentUser(): Observable<UserProfile | null> {
    return this.profileService.getProfile().pipe(
      map((res) => res.data ?? null),
      catchError(() => of(null)),
    );
  }

  getChatsByConversation(
    conversationId: string,
    currentUserId: string,
    limit = 10,
    before?: Date,
  ): Observable<Message[]> {
    let params: any = { limit };
    if (before) params.before = before.toISOString();

    return this.http
      .get<{
        data: ChatFromApi[];
      }>(`${this.apiUrl}/chats/conversation/${conversationId}`, { params })
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
        }),
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
      catchError(() => of(null)),
    );
  }

  deleteChat(id: string): Observable<boolean> {
    return this.http.delete(`${this.apiUrl}/${id}`).pipe(
      map(() => true),
      catchError(() => of(false)),
    );
  }

  startCall(payload: {
    conversationId: string;
    callerId: string;
    callerName: string;
    receiverId: string;
    sessionId?: string;
  }): Observable<any> {
    return this.http.post(`${this.apiUrl}/chats/call/start`, payload).pipe(
      tap((res) => console.log('✅ Call started:', res)),
      catchError((err) => {
        console.error('❌ Failed to start call:', err);
        return of(null);
      }),
    );
  }

  private isConnected = false;

  connect(receiverId: string, conversationId: string): Observable<Message> {
    const token = this.secureStorage.getItem(this.AccessTokenKey);
    const url = `${this.wsUrl}?token=${token}&receiver=${receiverId}&conversationId=${conversationId}`;
    console.log('[ChatService] Connection Request:', { receiverId, conversationId, url });

    const currentUserId = this.userContext.getUser()?.profileId || '';

    this.socket$ = webSocket<any>({
      url,
      openObserver: {
        next: () => {
          this.isConnected = true;
          this.reconnectAttempts = 0;
          this.connectionState$.next(true);
          console.log('%c [ChatService] WebSocket connected', 'color: green');

          this.startHeartbeat();
        },
      },
      closeObserver: {
        next: (closeEvent) => {
          this.isConnected = false;
          this.connectionState$.next(false);
          console.warn('%c [ChatService] WebSocket closed', 'color: orange', closeEvent);
          this.stopHeartbeat();
          this.tryReconnect(receiverId, conversationId);
        },
      },
      deserializer: (e) => {
        try {
          const parsed = JSON.parse(e.data);
          console.debug('[ChatService] Incoming WS Data:', parsed);
          return parsed;
        } catch (err) {
          console.error('[ChatService] WS Parse Error:', err, e.data);
          return null;
        }
      },
      serializer: (msg) => JSON.stringify(msg),
    });

    return this.socket$.asObservable().pipe(
      // 1. Filter out control messages (pings, connection confirmations, etc.)
      filter((data: any) => {
        if (!data) return false;
        // Only allow 'chat' type or messages that clearly have text content
        const messageData = data.data || data;
        const valid = data.type === 'chat' || !!data.text || !!data.content || !!messageData.text || !!messageData.content;
        if (!valid) console.debug('[ChatService] Filtering out message:', data);
        return valid;
      }),
      map((data: any) => {
        // Handle different possible Backend message structures
        // Some backends put data inside a 'data' property, others are flat
        const messageData = data.data || data;
        const senderId = data.senderId || messageData.senderId || data.userId || '';
        const isMine = senderId === currentUserId;

        // Try to find content in various possible locations
        let content = '';
        if (messageData.text) {
          content = messageData.text;
        } else if (messageData.content) {
          content = messageData.content;
        } else if (data.text) {
          content = data.text;
        } else if (data.content) {
          content = data.content;
        }

        // Handle JSON parsing of content if needed
        try {
          if (content && typeof content === 'string' && (content.startsWith('{') || content.startsWith('['))) {
            const parsed = JSON.parse(content);
            if (typeof parsed === 'string') content = parsed;
          }
        } catch (e) {}

        const msg: Message = {
          id: data.id || messageData.id || ('ws-' + Date.now()),
          conversationId: data.conversationId || messageData.conversationId || conversationId,
          senderId,
          userName: data.userName || messageData.userName || '',
          userAvatar: data.userAvatar || messageData.userAvatar || '',
          content: String(content),
          timestamp: data.createdAt ? new Date(data.createdAt) : 
                    (messageData.createdAt ? new Date(messageData.createdAt) : new Date()),
          isMine,
          sessionData: data.sessionData || messageData.sessionData
        };
        return msg;
      }),
      tap({
        error: (err) => console.error('WebSocket error:', err),
      }),
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
      }s`,
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
  private startHeartbeat() {
    this.stopHeartbeat();
    this.heartbeatInterval = setInterval(() => {
      if (this.socket$ && this.isConnected) {
        this.socket$.next({ type: 'ping' });
      }
    }, 30000); // 30 seconds
  }

  private stopHeartbeat() {
    if (this.heartbeatInterval) {
      clearInterval(this.heartbeatInterval);
      this.heartbeatInterval = null;
    }
  }
  sendMessage(msg: {
    conversationId: string;
    content: string;
    senderId: string;
    sessionData?: any;
  }) {
    if (!this.socket$ || !this.isConnected) {
      console.warn('WebSocket not connected, message dropped');
      return;
    }

    const payload = {
      type: 'chat', 
      conversationId: msg.conversationId,
      senderId: msg.senderId,
      data: {
        text: msg.content,
        sessionData: msg.sessionData,
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
    user2Id: string,
  ): Observable<{ id: string }> {
    console.log('getOrCreateConversation response~!!!:');
    return this.http
      .get<{
        code: number;
        message: string;
        data: { id: string };
      }>(`${this.apiUrl}/conversations/by-users`, {
        params: { user1: user1Id, user2: user2Id },
      })
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
        }),
      );
  }

  getRecentConversations(userId: string): Observable<any[]> {
    return this.http
      .get<{
        code: number;
        message: string;
        data: any[];
      }>(`${this.apiUrl}/conversations/me`, { params: { userId } })
      .pipe(
        map((res) => res.data || []),
        catchError((err) => {
          console.error('getRecentConversations error:', err);
          return of([]);
        }),
      );
  }

  attachProfilesToChatList(
    chats: Chat[],
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
