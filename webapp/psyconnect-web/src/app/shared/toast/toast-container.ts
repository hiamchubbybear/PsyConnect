import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { ToastData, ToastType } from './toast.model';

@Component({
  selector: 'app-toast-container',
  imports: [CommonModule],
  standalone : true,
  templateUrl: './toast-container.html',
  styleUrls: ['./toast-container.scss'],
})
export class ToastContainerComponent {
  toasts: ToastData[] = [];

  addToast(toast: ToastData) {
    this.toasts.push(toast);
    setTimeout(() => this.removeToast(toast.id), toast.duration || 4000);
  }

  removeToast(id: string) {
    this.toasts = this.toasts.filter((t) => t.id !== id);
  }

  getIcon(type: ToastType): string {
    switch (type) {
      case ToastType.Success:
        return 'check_circle';
      case ToastType.Info:
        return 'info';
      case ToastType.Warning:
        return 'warning';
      case ToastType.Error:
        return 'error';
    }
  }
}
