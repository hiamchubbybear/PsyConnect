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
      .loading-bar {
        position: fixed;
        top: 100%;
        left: 0;
        width: auto;
        height: 3px;
        background: linear-gradient(90deg, #000, #555);
        animation: loading 1s infinite linear;
      }

      @keyframes loading {
        0% {
          transform: translateX(-100%);
        }
        100% {
          transform: translateX(100%);
        }
      }
    `,
  ],
})
export class LoaderComponent {
  constructor(public loaderService: LoaderService) {}
}
