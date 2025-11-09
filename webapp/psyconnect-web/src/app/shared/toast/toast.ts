import { animate, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { ToastType } from './toast-type';

@Component({
  selector: 'app-toast',
  standalone: true,
  imports: [CommonModule,TranslateModule],
  templateUrl: './toast.html',
  styleUrls: ['./toast.scss'],
  animations: [
    trigger('toastAnimation', [
      transition(':enter', [
        style({ opacity: 0, transform: 'translateY(10px)' }),
        animate(
          '200ms ease-out',
          style({ opacity: 1, transform: 'translateY(0)' })
        ),
      ]),
      transition(':leave', [
        animate(
          '200ms ease-in',
          style({ opacity: 0, transform: 'translateY(10px)' })
        ),
      ]),
    ]),
  ],
})
export class ToastComponent {
  @Input() title!: string;
  @Input() id!: string;
  @Input() message!: string;
  @Input() type: ToastType = ToastType.Info;
  @Input() onClose!: () => void;
  @Output() close: EventEmitter<string> = new EventEmitter<string>();
  get icon(): string {
    switch (this.type) {
      case ToastType.Success:
        return 'check_circle';
      case ToastType.Info:
        return 'info';
      case ToastType.Warning:
        return 'warning';
      case ToastType.Error:
        return 'error';
      default:
        return 'info';
    }
  }

  handleClose() {
    if (this.onClose) this.onClose();
    this.close.emit(this.id);
  }
}
