import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-reviews',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <div class="reviews-container">
      <header class="page-header">
        <h1>{{ 'POST.Reviews.Title' | translate }}</h1>
        <p class="subtitle">{{ 'POST.Reviews.Subtitle' | translate }}</p>
      </header>
      <div class="coming-soon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polygon points="12 2 15.09 8.26 22 9.27 17 14.14 18.18 21.02 12 17.77 5.82 21.02 7 14.14 2 9.27 8.91 8.26 12 2"></polygon></svg>
        <h2>{{ 'COMMON.ComingSoon' | translate }}</h2>
        <p>{{ 'POST.Reviews.ComingSoonMessage' | translate }}</p>
      </div>
    </div>
  `,
  styles: [`
    .reviews-container {
      padding: var(--spacing-2xl);
      max-width: 1400px;
      margin: 0 auto;
    }
    .page-header h1 {
      font-size: var(--font-size-3xl);
      font-weight: 700;
      color: var(--color-text);
      margin: 0 0 var(--spacing-sm) 0;
    }
    .subtitle {
      color: var(--color-text-muted);
    }
    .coming-soon {
      text-align: center;
      padding: var(--spacing-3xl);
      i {
        font-size: 4rem;
        color: var(--color-text-muted);
        opacity: 0.3;
      }
      h2 {
        color: var(--color-text);
        margin: var(--spacing-lg) 0 var(--spacing-sm) 0;
      }
      p {
        color: var(--color-text-muted);
      }
    }
  `]
})
export class Reviews {}
