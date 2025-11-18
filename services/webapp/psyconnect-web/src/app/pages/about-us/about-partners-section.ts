import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Certificate, Partner } from './about-us-data-service';

@Component({
  selector: 'app-partners-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__partners container">
      <h2 class="section-title">{{ "featured_experts" | translate }}</h2>
      <p class="section-desc">{{ "featured_experts_desc" | translate }}</p>

      <div class="partner-grid">
        <div *ngFor="let p of partners" class="partner-card">
          <div class="partner-card__header">
            <img
              class="partner-card__avatar"
              [src]="p.avatar"
              [alt]="p.name"
            />
            <div class="partner-card__verified" *ngIf="p.verified">
              <i class="icon-shield-check"></i>
              {{ "verified" | translate }}
            </div>
          </div>

          <div class="partner-card__meta">
            <h4 class="partner-card__name">{{ p.name }}</h4>
            <p class="partner-card__role">{{ p.role | translate }}</p>

            <div class="partner-card__stats">
              <div class="stat">
                <i class="icon-star"></i>
                {{ p.rating }} ({{ p.reviewCount }}
                {{ "reviews" | translate }})
              </div>
              <div class="stat">
                <i class="icon-clock"></i>
                {{ p.experience }} {{ "years_experience" | translate }}
              </div>
            </div>

            <div class="partner-card__specialization">
              <span class="chip" *ngFor="let spec of p.specialization">
                {{ spec | translate }}
              </span>
            </div>

            <div class="partner-card__education">
              <h5>{{ "education" | translate }}:</h5>
              <ul>
                <li *ngFor="let edu of p.education">{{ edu }}</li>
              </ul>
            </div>

            <div class="partner-card__certificates">
              <h5>{{ "certificates" | translate }}:</h5>
              <div class="certificate-list">
                <div
                  *ngFor="let cert of p.certificates"
                  class="certificate-item"
                >
                  <div class="cert-info">
                    <strong>{{ cert.name }}</strong>
                    <small
                      >{{ cert.issuer }} - {{ cert.issueDate | date }}</small
                    >
                  </div>
                  <button
                    *ngIf="cert.verificationUrl"
                    class="btn btn--sm btn--ghost"
                    (click)="verifyCertificate.emit(cert)"
                  >
                    {{ "verify" | translate }}
                  </button>
                </div>
              </div>
            </div>

            <div class="partner-card__languages">
              <span *ngFor="let lang of p.languages" class="language-tag">{{
                lang
              }}</span>
            </div>
          </div>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-partners-section.scss'],
})
export class PartnersSectionComponent {
  @Input() partners: Partner[] = [];
  @Output() verifyCertificate = new EventEmitter<Certificate>();
}
