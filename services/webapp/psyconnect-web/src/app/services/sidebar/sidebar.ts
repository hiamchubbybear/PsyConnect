import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';

@Injectable({ providedIn: 'root' })
export class SidebarService {
  private collapsedSubject = new BehaviorSubject<boolean>(false);
  collapsed$ = this.collapsedSubject.asObservable();

  toggleSidebar(collapsed?: boolean) {
    if (collapsed !== undefined) {
      this.collapsedSubject.next(collapsed);
    } else {
      this.collapsedSubject.next(!this.collapsedSubject.value);
    }
  }
}
