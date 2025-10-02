import { Injectable } from '@angular/core';
import { catchError, map, Observable, of } from 'rxjs';
import { Chat } from '../../models/chat.models';
import { ProfileResponse } from '../../models/profile';
import { Profile } from '../profile/profile';
import { UserContextService, UserProfile } from '../profile/profile-service';

@Injectable({
  providedIn: 'root',
})
export class ChatService {
  constructor(
    private userContext: UserContextService,
    private profileService: Profile
  ) {}

  getCurrentUser(): Observable<UserProfile | null> {
    return this.profileService.getProfile().pipe(
      map((res) => res.data ?? null),
      catchError(() => of(null))
    );
  }

  getCurrentUserFromApi(): Observable<UserProfile> {
    return this.profileService
      .getProfile()
      .pipe(map((res: ProfileResponse) => res.data));
  }

  attachProfileToChat(chat: Chat): Observable<Chat & { user?: UserProfile }> {
    const localUser = this.userContext.getUser();
    if (localUser) {
      return of({ ...chat, user: localUser });
    }

    return this.getCurrentUserFromApi().pipe(
      map((profile) => ({ ...chat, user: profile }))
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
