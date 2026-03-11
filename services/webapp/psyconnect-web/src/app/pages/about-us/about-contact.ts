import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

interface ContactInfo {
  hotline: string;
  email: string;
  emergency: string;
  socials: {
    facebook: string;
    instagram: string;
    linkedin: string;
  };
  address: string;
  businessHours: {
    weekdays: string;
    weekends: string;
  };
}

import { IconComponent } from '../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-contact-section',
  standalone: true,
  imports: [CommonModule, TranslateModule, IconComponent],
  template: `
    <section class="about-us__contact container">
      <h2 class="section-title">{{ "contact_us" | translate }}</h2>

      <div class="contact-grid">
        <div class="contact-info">
          <div class="contact-item">
            <app-icon name="phone" [size]="20"></app-icon>
            <div>
              <strong>{{ "hotline" | translate }}</strong>
              <p>{{ contact.hotline }}</p>
            </div>
          </div>

          <div class="contact-item">
            <app-icon name="mail" [size]="20"></app-icon>
            <div>
              <strong>{{ "email" | translate }}</strong>
              <p>
                <a href="mailto:{{ contact.email }}">{{ contact.email }}</a>
              </p>
            </div>
          </div>

          <div class="contact-item">
            <app-icon name="map-pin" [size]="20"></app-icon>
            <div>
              <strong>{{ "address" | translate }}</strong>
              <p>{{ contact.address | translate }}</p>
            </div>
          </div>

          <div class="contact-item">
            <app-icon name="clock" [size]="20"></app-icon>
            <div>
              <strong>{{ "business_hours" | translate }}</strong>
              <p>
                {{ "weekdays" | translate }}:
                {{ contact.businessHours.weekdays }}
              </p>
              <p>
                {{ "weekends" | translate }}:
                {{ contact.businessHours.weekends }}
              </p>
            </div>
          </div>
        </div>

        <div class="emergency-contact">
          <div class="emergency-card">
            <h3>{{ "crisis_support" | translate }}</h3>
            <p>{{ "crisis_support_desc" | translate }}</p>
            <button class="btn btn--danger" (click)="emergencyHelp.emit()">
              <app-icon name="phone" [size]="18"></app-icon>
              {{ "emergency_call" | translate }}: {{ contact.emergency }}
            </button>
          </div>
        </div>
      </div>

      <div class="social-links">
        <a
          [href]="contact.socials.facebook"
          target="_blank"
          class="social-btn"
        >
          <app-icon name="facebook" [size]="20"></app-icon>
        </a>
        <a
          [href]="contact.socials.instagram"
          target="_blank"
          class="social-btn"
        >
          <app-icon name="instagram" [size]="20"></app-icon>
        </a>
        <a
          [href]="contact.socials.linkedin"
          target="_blank"
          class="social-btn"
        >
          <app-icon name="linkedin" [size]="20"></app-icon>
        </a>
      </div>
    </section>
  `,
  styleUrls: ['./about-contact.scss'],
})
export class ContactSectionComponent {
  @Input() contact!: ContactInfo;
  @Output() emergencyHelp = new EventEmitter<void>();
}
