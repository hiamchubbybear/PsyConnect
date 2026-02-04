import { HttpClient } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Messaging, getToken, onMessage } from '@angular/fire/messaging';
import { BehaviorSubject } from 'rxjs';
import { environment } from '../../environments/environment';

@Injectable({
  providedIn: 'root',
})
export class FcmService {
  currentMessage = new BehaviorSubject<any>(null);

  constructor(
    private messaging: Messaging,
    private http: HttpClient,
  ) {}

  async requestPermission(userId: string) {
    console.log('Requesting permission...');
    const permission = await Notification.requestPermission();
    if (permission === 'granted') {
      console.log('Notification permission granted.');
      const token = await getToken(this.messaging, {
        vapidKey: environment.vapidKey,
      });
      console.log('FCM Token:', token);
      if (token) {
        // Send token to backend
        this.saveToken(userId, token);
      }
    } else {
      console.log('Unable to get permission to notify.');
    }
  }

  listen() {
    console.log('Listening for messages...');
    onMessage(this.messaging, (payload) => {
      console.log('Message received. ', payload);
      this.currentMessage.next(payload);
    });
  }

  private saveToken(userId: string, token: string) {
    // Assuming Identity Service or Notification Service endpoint
    // Adjust URL based on your backend
    // Using mock for now since exact endpoint isn't clarified in context,
    // but based on Notification Service handler: SaveFCMToken
    // It likely expects a calling service (Profile/Identity) or direct call?
    // Notification Service consumer listens to events, but we need to SAVE token.

    // Let's assume we post to Profile Service or directly to Notification if exposed.
    // Or simply log for manual testing if backend saving endpoint is not ready.

    console.log('TODO: Save this token to backend:', userId, token);
    // this.http.post(`${environment.apiUrl}/${environment.apiVersion}/notifications/token`, { token }).subscribe();
  }
}
