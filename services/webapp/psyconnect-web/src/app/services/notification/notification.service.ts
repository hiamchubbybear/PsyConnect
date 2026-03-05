import { HttpClient, HttpHeaders } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Messaging, getToken, onMessage } from '@angular/fire/messaging';
import { BehaviorSubject } from 'rxjs';
import { environment } from '../../../environments/environment';

export interface InAppNotification {
  id: number;
  userId: string;
  title: string;
  body: string;
  type: string;
  metadata: any;
  isRead: boolean;
  createdAt: string;
}

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private messaging = inject(Messaging);
  private notificationsSubject = new BehaviorSubject<InAppNotification[]>([]);
  notifications$ = this.notificationsSubject.asObservable();
  unreadCountSubject = new BehaviorSubject<number>(0);
  unreadCount$ = this.unreadCountSubject.asObservable();

  constructor(private http: HttpClient) {}

  async init(userId: string) {
    try {
      const permission = await Notification.requestPermission();
      if (permission !== 'granted') {
        console.warn('User từ chối cấp quyền notification');
        return;
      }
      const token = await getToken(this.messaging, {
        vapidKey: environment.vapidKey,
      });

      if (token) {
        console.log('FCM Token:', token);
        await this.http
          .post(`${environment.apiUrl}/notification/token`, { userId, token })
          .toPromise();
      } else {
        console.warn('Không lấy được token — user chưa cho phép thông báo');
      }

      onMessage(this.messaging, (payload) => {
        console.log('Notification nhận được:', payload);
        const { title, body } = payload.notification || {};
        new Notification(title || 'Thông báo', { body });
        
        this.fetchNotifications(userId);
      });

      
      this.fetchNotifications(userId);
    } catch (error) {
      console.error('Lỗi khi khởi tạo notification:', error);
    }
  }

  fetchNotifications(userId: string) {
    const headers = new HttpHeaders().set('X-Profile-Id', userId);
    this.http
      .get<InAppNotification[]>(`${environment.apiUrl}/notification/me`, {
        headers,
      })
      .subscribe((notifications) => {
        this.notificationsSubject.next(notifications);
        const unread = notifications.filter((n) => !n.isRead).length;
        this.unreadCountSubject.next(unread);
      });
  }

  markAsRead(userId: string, notificationId: number) {
    const headers = new HttpHeaders().set('X-Profile-Id', userId);
    return this.http
      .patch(
        `${environment.apiUrl}/notification/${notificationId}/read`,
        {},
        { headers }
      )
      .subscribe(() => {
        
        const current = this.notificationsSubject.value;
        const updated = current.map((n) =>
          n.id === notificationId ? { ...n, isRead: true } : n
        );
        this.notificationsSubject.next(updated);
        this.unreadCountSubject.next(updated.filter((n) => !n.isRead).length);
      });
  }
}
