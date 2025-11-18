import { CommonModule } from '@angular/common';
import { Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs';
import { ToastComponent } from './toast';
import { ToastData } from './toast.model';
import { ToastService } from './toast.service';

@Component({
  selector: 'app-toast-container',
  standalone: true,
  imports: [CommonModule, ToastComponent],
  templateUrl: './toast-container.html',
  styleUrls: ['./toast-container.scss'],
})
export class ToastContainerComponent implements OnDestroy {
  toasts: ToastData[] = [];
  private sub: Subscription;

  constructor(private toastService: ToastService) {
    this.sub = this.toastService.toasts$.subscribe((toast) => {
      this.addToast(toast);
    });
  }

  addToast(toast: ToastData) {
    this.toasts.push(toast);
    setTimeout(() => this.removeToast(toast.id), toast.duration || 4000);
  }

  removeToast(id: string) {
    this.toasts = this.toasts.filter((t) => t.id !== id);
  }

  ngOnDestroy() {
    this.sub.unsubscribe();
  }
}
