import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { SecureStorageService } from '../../encrypt/secure';

@Injectable({ providedIn: 'root' })
export class HeaderStateService {
  private key = 'headerIsMini';
  private miniSubject: BehaviorSubject<boolean>;

  isMini$;

  constructor(private secureStorageService: SecureStorageService) {
    const initial = this.loadMiniState();
    this.miniSubject = new BehaviorSubject<boolean>(initial);
    this.isMini$ = this.miniSubject.asObservable();
  }

  setMini(value: boolean) {
    this.miniSubject.next(value);
    this.secureStorageService.setItem(this.key, value);
  }

  getMini(): boolean {
    return this.miniSubject.getValue();
  }

  private loadMiniState(): boolean {
    return this.secureStorageService.getItem(this.key) ?? false;
  }
}
