import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Service } from './about-us-data-service';

@Component({
  selector: 'app-services-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__services container">
      <h2 class="section-title">{{ "our_services" | translate }}</h2>
      <p class="section-desc">{{ "services_description" | translate }}</p>

      <div class="service-grid">
        <div *ngFor="let s of services" class="service-card">
          <div class="service-card__icon">
            <i class="icon-{{ s.icon }}"></i>
          </div>
          <div class="service-card__content">
            <h4 class="service-card__name">{{ s.name | translate }}</h4>
            <p class="service-card__desc">{{ s.desc | translate }}</p>
            <div class="service-card__details">
              <div class="detail" *ngIf="s.duration">
                <i class="icon-clock"></i>
                {{ s.duration }}
              </div>
              <div class="detail" *ngIf="s.priceRange">
                <i class="icon-dollar-sign"></i>
                {{ s.priceRange }}
              </div>
            </div>
          </div>
        </div>
      </div>

      <div class="methods-section">
        <h3 class="subhead">{{ "therapeutic_methods" | translate }}</h3>
        <div class="methods-grid">
          <span class="method-chip" *ngFor="let m of methods">{{
            m | translate
          }}</span>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-service-section.scss'],
})
export class ServicesSectionComponent {
  @Input() services: Service[] = [];
  @Input() methods: string[] = [];
}
