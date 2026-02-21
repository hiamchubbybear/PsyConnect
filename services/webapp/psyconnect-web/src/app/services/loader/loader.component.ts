import { animate, style, transition, trigger } from '@angular/animations';
import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { LoaderService } from './loader';
@Component({
  selector: 'app-loader',
  standalone: true,
  imports: [CommonModule],
  animations: [
    trigger('fadeIn', [
      transition(':enter', [
        style({ opacity: 0 }),
        animate('200ms ease-out', style({ opacity: 1 })),
      ]),
      transition(':leave', [animate('200ms ease-in', style({ opacity: 0 }))]),
    ]),
  ],
  template: `
    <div class="loader-bar" *ngIf="loaderService.loading$ | async"></div>
  `,
  styles: [
    `
      .loader-bar {
        position: fixed;
        top: 0;
        left: 0;
        width: 100vw;
        height: 3px;
        z-index: 999999;
        background: transparent;
        overflow: hidden;
      }
      .loader-bar::after {
        content: '';
        position: absolute;
        top: 0;
        bottom: 0;
        left: 0;
        width: 100%;
        background: linear-gradient(
          90deg,
          var(--color-primary, #667eea),
          var(--color-accent, #ed8936)
        );
        transform-origin: left;
        animation: indeterminate-loader 2.5s infinite ease-in-out;
      }
      @keyframes indeterminate-loader {
        0% {
          transform: translateX(-100%) scaleX(0.2);
        }
        50% {
          transform: translateX(0) scaleX(1);
        }
        100% {
          transform: translateX(100%) scaleX(0.2);
        }
      }
    `,
  ],
})
export class LoaderComponent {
  constructor(public loaderService: LoaderService) {}
}
