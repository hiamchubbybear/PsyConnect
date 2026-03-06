import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-settings',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <div class="settings-container">
      <header class="page-header">
        <h1>{{ 'POST.Settings.Title' | translate }}</h1>
        <p class="subtitle">{{ 'POST.Settings.Subtitle' | translate }}</p>
      </header>
      <div class="coming-soon">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="3"></circle><path d="M19.4 15a1.65 1.65 0 0 0 .33 1.82l.06.06a2 2 0 0 1 0 2.83 2 2 0 0 1-2.83 0l-.06-.06a1.65 1.65 0 0 0-1.82-.33 1.65 1.65 0 0 0-1 1.51V21a2 2 0 0 1-2 2 2 2 0 0 1-2-2v-.09A1.65 1.65 0 0 0 9 19.4a1.65 1.65 0 0 0-1.82.33l-.06.06a2 2 0 0 1-2.83 0 2 2 0 0 1 0-2.83l.06-.06a1.65 1.65 0 0 0 .33-1.82 1.65 1.65 0 0 0-1.51-1H3a2 2 0 0 1-2-2 2 2 0 0 1 2-2h.09A1.65 1.65 0 0 0 4.6 9a1.65 1.65 0 0 0-.33-1.82l-.06-.06a2 2 0 0 1 0-2.83 2 2 0 0 1 2.83 0l.06.06a1.65 1.65 0 0 0 1.82.33H9a1.65 1.65 0 0 0 1-1.51V3a2 2 0 0 1 2-2 2 2 0 0 1 2 2v.09a1.65 1.65 0 0 0 1 1.51 1.65 1.65 0 0 0 1.82-.33l.06-.06a2 2 0 0 1 2.83 0 2 2 0 0 1 0 2.83l-.06.06a1.65 1.65 0 0 0-.33 1.82V9a1.65 1.65 0 0 0 1.51 1H21a2 2 0 0 1 2 2 2 2 0 0 1-2 2h-.09a1.65 1.65 0 0 0-1.51 1z"></path></svg>
        <h2>{{ 'COMMON.ComingSoon' | translate }}</h2>
        <p>{{ 'POST.Settings.ComingSoonMessage' | translate }}</p>
      </div>
    </div>
  `,
  styles: [`
    .settings-container {
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
export class Settings {}
