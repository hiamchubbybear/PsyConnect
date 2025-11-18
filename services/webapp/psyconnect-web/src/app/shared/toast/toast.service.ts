import { Injectable } from '@angular/core';
import { Subject } from 'rxjs';
import { ToastData, ToastType } from './toast.model';

@Injectable({
  providedIn: 'root',
})
export class ToastService {
  private toastSubject = new Subject<ToastData>();
  toasts$ = this.toastSubject.asObservable();

  show(
    title: string,
    message: string,
    type: ToastType = ToastType.Info,
    duration = 4000
  ) {
    const toast: ToastData = {
      id: crypto.randomUUID(),
      title,
      message,
      type,
      duration,
    };
    this.toastSubject.next(toast);
  }

  success(title: string, message: string, duration = 40000) {
    this.show(title, message, ToastType.Success, duration);
  }

  info(title: string, message: string, duration = 4000) {
    this.show(title, message, ToastType.Info, duration);
  }

  warning(title: string, message: string, duration = 4000) {
    this.show(title, message, ToastType.Warning, duration);
  }

  error(title: string, message: string, duration = 4000) {
    this.show(title, message, ToastType.Error, duration);
  }
}
function generateUUID(): string {
  if ('randomUUID' in crypto) {
    return crypto.randomUUID();
  }
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = crypto.getRandomValues(new Uint8Array(1))[0] & 15;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}
