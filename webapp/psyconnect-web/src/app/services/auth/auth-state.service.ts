import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { SecureStorageService } from '../../encrypt/secure';

@Injectable({ providedIn: 'root' })
export class AuthStateService {
  private loggedIn = new BehaviorSubject<boolean>(false);
  loggedIn$ = this.loggedIn.asObservable();

  private sidebarVisible = new BehaviorSubject<boolean>(false);
  sidebarVisible$ = this.sidebarVisible.asObservable();

  constructor(private secureStorage: SecureStorageService) {
    const token = this.secureStorage.getItem('access_token');
    this.loggedIn.next(!!token);
    this.sidebarVisible.next(!!token);
  }

  setLoggedIn(value: boolean) {
    this.loggedIn.next(value);
    if (value) {
      this.showSidebar();
    } else {
      this.hideSidebar();
    }
  }

  showSidebar() {
    this.sidebarVisible.next(true);
  }

  hideSidebar() {
    this.sidebarVisible.next(false);
  }
}
