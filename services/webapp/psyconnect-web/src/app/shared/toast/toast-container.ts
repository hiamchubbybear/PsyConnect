import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { Subscription } from 'rxjs';
import { HeaderStateService } from '../../components/header/header-state';
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
export class ToastContainerComponent implements OnInit, OnDestroy {
  toasts: ToastData[] = [];
  isHeaderMini = false;
  private sub = new Subscription();

  constructor(
    private toastService: ToastService,
    private headerState: HeaderStateService,
  ) {}

  ngOnInit() {
    this.sub.add(
      this.toastService.toasts$.subscribe((toasts) => {
        // Prioritize Session type to be always at the front (reverse makes last element first)
        // We sort sessions to the end of the original list so they end up at the start after reversal
        this.toasts = [...toasts]
          .sort((a, b) => {
            if (a.type === 'session' && b.type !== 'session') return 1;
            if (a.type !== 'session' && b.type === 'session') return -1;
            return 0;
          })
          .reverse();
      }),
    );

    this.sub.add(
      this.headerState.isMini$.subscribe((mini: boolean) => {
        this.isHeaderMini = mini;
      }),
    );
  }

  removeToast(id: string) {
    this.toastService.remove(id);
  }

  trackByFn(index: number, item: ToastData) {
    return item.id;
  }

  ngOnDestroy() {
    this.sub?.unsubscribe();
  }
}
