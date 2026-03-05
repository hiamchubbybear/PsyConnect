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
    
    
    
    
    
    

    
    

    console.log('TODO: Save this token to backend:', userId, token);
    
  }
}
