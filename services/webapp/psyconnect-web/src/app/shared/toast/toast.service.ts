import { Injectable } from '@angular/core';
import { BehaviorSubject } from 'rxjs';
import { ToastData, ToastType } from './toast.model';

@Injectable({
  providedIn: 'root',
})
export class ToastService {
  private toastsSubject = new BehaviorSubject<ToastData[]>([]);
  toasts$ = this.toastsSubject.asObservable();

  show(data: Partial<ToastData>) {
    const toast: ToastData = {
      id: data.id || generateUUID(),
      title: data.title,
      message: data.message || '',
      description: data.description,
      type: data.type || ToastType.Info,
      duration: data.duration ?? 5000,
      showCountdown: data.showCountdown ?? true,
      actions: data.actions || [],
    };

    const currentToasts = [...this.toastsSubject.value];
    const existingIndex = currentToasts.findIndex((t) => t.id === toast.id);

    if (existingIndex > -1) {
      currentToasts[existingIndex] = toast;
    } else {
      currentToasts.push(toast);
    }

    this.toastsSubject.next(currentToasts);

    if (toast.duration && toast.duration > 0) {
      const staggerDelay =
        (existingIndex > -1 ? existingIndex : currentToasts.length - 1) * 1000;
      setTimeout(() => {
        this.remove(toast.id);
      }, toast.duration + staggerDelay);
    }
  }

  remove(id: string) {
    const remaining = this.toastsSubject.value.filter((t) => t.id !== id);
    this.toastsSubject.next(remaining);
  }

  clearAll() {
    this.toastsSubject.next([]);
  }

  success(title: string, message: string, duration = 4000) {
    this.show({ title, message, type: ToastType.Success, duration });
  }

  error(title: string, message: string, duration = 5000) {
    this.show({ title, message, type: ToastType.Error, duration });
  }

  warning(title: string, message: string, duration = 4000) {
    this.show({ title, message, type: ToastType.Warning, duration });
  }

  info(title: string, message: string, duration = 4000) {
    this.show({ title, message, type: ToastType.Info, duration });
  }

  session(data: Partial<ToastData>) {
    // Persistent by default (duration=0)
    this.show({ ...data, type: ToastType.Session, duration: 0 });
  }
}

function generateUUID(): string {
  if (typeof crypto !== 'undefined' && crypto.randomUUID) {
    return crypto.randomUUID();
  }
  return 'xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx'.replace(/[xy]/g, (c) => {
    const r = (Math.random() * 16) | 0;
    const v = c === 'x' ? r : (r & 0x3) | 0x8;
    return v.toString(16);
  });
}
