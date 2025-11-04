import { HttpClient } from '@angular/common/http';
import { Injectable, inject } from '@angular/core';
import { Messaging, getToken, onMessage } from '@angular/fire/messaging';
import { environment } from '../../../environments/environment';

@Injectable({ providedIn: 'root' })
export class NotificationService {
  private messaging = inject(Messaging);

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
      });
    } catch (error) {
      console.error('Lỗi khi khởi tạo notification:', error);
    }
  }
}
