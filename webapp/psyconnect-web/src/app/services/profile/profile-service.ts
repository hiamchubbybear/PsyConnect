import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { SecureStorageService } from '../../encrypt/secure';

export interface UserProfile {
  accountId: string;
  profileId: string;
  firstName: string;
  lastName: string;
  dob: string;
  address: string;
  gender: string;
  avatarUri: string;
}

@Injectable({
  providedIn: 'root',
})
export class UserContextService {
  private userSubject = new BehaviorSubject<UserProfile | null>(null);
  user$ = this.userSubject.asObservable();
  private isInitialized = false;

  constructor(private secureStorage: SecureStorageService) {
    this.initializeUser();
  }

  private initializeUser() {
    if (this.isInitialized) {
      return;
    }

    try {
      const userString = this.secureStorage.getItem<string>('profile');
      if (userString) {
        const user = JSON.parse(userString);

        this.userSubject.next(user);
        this.isInitialized = true;
      } else {
      }
    } catch (error) {
      this.secureStorage.removeItem('profile');
    }
  }

  setUser(user: UserProfile) {
    try {
      const currentUser = this.userSubject.value;

      if (
        !currentUser ||
        JSON.stringify(currentUser) !== JSON.stringify(user)
      ) {
        this.userSubject.next(user);
        this.secureStorage.setItem('profile', JSON.stringify(user));
        this.isInitialized = true;
      } else {
      }
    } catch (error) {}
  }

  getUser(): UserProfile | null {
    return this.userSubject.value;
  }

  clear() {
    this.userSubject.next(null);
    this.secureStorage.removeItem('profile');
    this.isInitialized = false;
  }

  reloadFromStorage() {
    if (!this.isInitialized) {
      this.initializeUser();
    } else {
    }
  }

  hasUser(): boolean {
    return this.userSubject.value !== null;
  }
}
