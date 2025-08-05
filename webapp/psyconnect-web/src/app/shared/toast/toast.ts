import { CommonModule } from '@angular/common';
import { Component, Input, OnInit } from '@angular/core';
import { ToastType } from './toast-type';

import { animate, style, transition, trigger } from '@angular/animations';

@Component({
  selector: 'app-toast',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './toast.html',
  styleUrls: ['./toast.scss'],
  animations: [
    trigger('fade', [
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
export class ToastComponent implements OnInit {
  @Input() title!: string;
  @Input() message!: string;
  @Input() type: ToastType = ToastType.Info;
  @Input() onClose!: () => void;

  icon = '';
  bgColor = '';
  accentColor = '';

  ngOnInit(): void {
    switch (this.type) {
      case ToastType.Success:
        this.bgColor = '#e6f4ea';
        this.icon = 'check_circle';
        this.accentColor = '#28a745';
        break;
      case ToastType.Info:
        this.bgColor = '#e8f0fe';
        this.icon = 'info';
        this.accentColor = '#1e88e5';
        break;
      case ToastType.Warning:
        this.bgColor = '#fff3cd';
        this.icon = 'warning';
        this.accentColor = '#ffc107';
        break;
      case ToastType.Error:
        this.bgColor = '#f8d7da';
        this.icon = 'error';
        this.accentColor = '#dc3545';
        break;
    }
  }
}
