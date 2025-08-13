import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class AuthStateService {
  private sidebarVisible = new BehaviorSubject<boolean>(false);
  sidebarVisible$ = this.sidebarVisible.asObservable();

  showSidebar() {
    this.sidebarVisible.next(true);
  }

  hideSidebar() {
    this.sidebarVisible.next(false);
  }
}
