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
