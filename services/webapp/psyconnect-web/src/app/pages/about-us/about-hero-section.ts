import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Statistic } from './about-us-data-service';
import { IconComponent } from '../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-hero-section',
  standalone: true,
  imports: [CommonModule, TranslateModule, IconComponent],
  template: `
    <section class="about-us__hero">
      <div class="about-us__inner container">
        <div class="hero-content">
          <img
            class="about-us__logo"
            [src]="brand.logo"
            [alt]="brand.name + ' logo'"
          />
          <div class="about-us__meta">
            <h1 class="about-us__title">{{ brand.name }}</h1>
            <p class="about-us__tagline">{{ "hero_tagline" | translate }}</p>
            <p class="about-us__mission">{{ mission }}</p>

            <div class="about-us__values">
              <span class="pill" *ngFor="let v of coreValues">{{
                v | translate
              }}</span>
            </div>

            <div class="about-us__cta">
              <button class="btn btn--primary" (click)="bookNow.emit()">
                {{ "find_expert" | translate }}
              </button>
              <button class="btn btn--ghost" (click)="listAsPartner.emit()">
                {{ "become_partner" | translate }}
              </button>
            </div>
          </div>
        </div>

        <div class="statistics-grid">
          <div class="stat-card" *ngFor="let stat of statistics">
            <div class="stat-card__icon">
              <app-icon [name]="stat.icon"></app-icon>
            </div>
            <div class="stat-card__value">{{ stat.value }}</div>
            <div class="stat-card__label">{{ stat.label | translate }}</div>
          </div>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-hero-section.scss'],
})
export class HeroSectionComponent {
  @Input() brand: any;
  @Input() statistics: Statistic[] = [];
  @Input() mission: string = '';
  @Input() coreValues: string[] = [];

  @Output() bookNow = new EventEmitter<void>();
  @Output() listAsPartner = new EventEmitter<void>();
}
