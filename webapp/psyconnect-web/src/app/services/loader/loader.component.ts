import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { LoaderService } from './loader';

@Component({
  selector: 'app-loader',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="loader-backdrop" *ngIf="loader.loading$ | async" [@fadeIn]>
      <div class="loader-container">
        <div class="spinner"></div>
        <div class="loading-text">Loading...</div>
      </div>
    </div>
  `,
  styles: [
    `
      .loader-backdrop {
        position: fixed;
        inset: 0;
        background: rgba(255, 255, 255, 0.95);
        backdrop-filter: blur(2px);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
        animation: fadeIn 0.2s ease-out;
      }

      .loader-container {
        display: flex;
        flex-direction: column;
        align-items: center;
        gap: 1rem;
      }

      .spinner {
        width: 40px;
        height: 40px;
        border: 2px solid #f0f0f0;
        border-top: 2px solid #2196f3;
        border-radius: 50%;
        animation: spin 0.8s cubic-bezier(0.4, 0, 0.2, 1) infinite;
        position: relative;
      }

      .spinner::before {
        content: '';
        position: absolute;
        inset: -2px;
        border: 2px solid transparent;
        border-top: 2px solid rgba(33, 150, 243, 0.3);
        border-radius: 50%;
        animation: spin 1.2s cubic-bezier(0.4, 0, 0.2, 1) infinite reverse;
      }

      .loading-text {
        font-size: 14px;
        font-weight: 500;
        color: #666;
        letter-spacing: 0.5px;
        animation: pulse 1.5s ease-in-out infinite;
      }

      @keyframes spin {
        to {
          transform: rotate(360deg);
        }
      }

      @keyframes fadeIn {
        from {
          opacity: 0;
        }
        to {
          opacity: 1;
        }
      }

      @keyframes pulse {
        0%,
        100% {
          opacity: 0.6;
        }
        50% {
          opacity: 1;
        }
      }

      /* Dark theme support */
      @media (prefers-color-scheme: dark) {
        .loader-backdrop {
          background: rgba(17, 17, 17, 0.95);
        }

        .spinner {
          border-color: #333;
          border-top-color: #2196f3;
        }

        .spinner::before {
          border-top-color: rgba(33, 150, 243, 0.3);
        }

        .loading-text {
          color: #ccc;
        }
      }

      /* Reduced motion preference */
      @media (prefers-reduced-motion: reduce) {
        .spinner,
        .spinner::before {
          animation-duration: 2s;
        }

        .loading-text {
          animation: none;
          opacity: 0.8;
        }

        .loader-backdrop {
          animation: none;
        }
      }

      /* Mobile optimization */
      @media (max-width: 768px) {
        .spinner {
          width: 36px;
          height: 36px;
        }

        .loading-text {
          font-size: 13px;
        }
      }
    `,
  ],
})
export class LoaderComponent {
  constructor(public loader: LoaderService) {}
}
