import {
  animate,
  state,
  style,
  transition,
  trigger,
} from '@angular/animations';
import { CommonModule } from '@angular/common';
import {
  Component,
  EventEmitter,
  Input,
  OnDestroy,
  OnInit,
  Output,
} from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { ToastAction, ToastData } from './toast.model';

@Component({
  selector: 'app-toast',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './toast.html',
  styleUrls: ['./toast.scss'],
  animations: [
    trigger('toastAnimation', [
      state(
        'void',
        style({
          opacity: 0,
          transform: 'translateX(30px) scale(0.95)',
          filter: 'blur(4px)',
        }),
      ),
      transition(':enter', [animate('400ms cubic-bezier(0.16, 1, 0.3, 1)')]),
      transition(':leave', [
        animate(
          '300ms cubic-bezier(0.7, 0, 0.84, 0)',
          style({
            opacity: 0,
            transform: 'translateX(20px) scale(0.9)',
            filter: 'blur(4px)',
          }),
        ),
      ]),
    ]),
  ],
})
export class ToastComponent implements OnInit, OnDestroy {
  @Input() data!: ToastData;
  @Input() index = 0;
  @Output() close = new EventEmitter<string>();

  remainingSeconds = 0;
  private intervalId?: any;

  ngOnInit() {
    if (this.data.duration && this.data.duration > 0) {
      this.remainingSeconds = Math.ceil(this.data.duration / 1000);
      this.startCountdown();
    }
  }

  ngOnDestroy() {
    this.stopCountdown();
  }

  private startCountdown() {
    this.intervalId = setInterval(() => {
      if (this.remainingSeconds > 0) {
        this.remainingSeconds--;
      } else {
        this.stopCountdown();
      }
    }, 1000);
  }

  private stopCountdown() {
    if (this.intervalId) {
      clearInterval(this.intervalId);
    }
  }

  handleAction(action: ToastAction) {
    action.action();
    this.close.emit(this.data.id);
  }

  handleClose() {
    this.close.emit(this.data.id);
  }

  get stackingStyle() {
    return {
      'z-index': 100 - this.index,
      'pointer-events': 'auto',
      position: 'relative',
    };
  }
}
