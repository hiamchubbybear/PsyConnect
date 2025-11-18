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

@Component({
  selector: 'app-contact-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__contact container">
      <h2 class="section-title">{{ "contact_us" | translate }}</h2>

      <div class="contact-grid">
        <div class="contact-info">
          <div class="contact-item">
            <i class="icon-phone"></i>
            <div>
              <strong>{{ "hotline" | translate }}</strong>
              <p>{{ contact.hotline }}</p>
            </div>
          </div>

          <div class="contact-item">
            <i class="icon-mail"></i>
            <div>
              <strong>{{ "email" | translate }}</strong>
              <p>
                <a href="mailto:{{ contact.email }}">{{ contact.email }}</a>
              </p>
            </div>
          </div>

          <div class="contact-item">
            <i class="icon-map-pin"></i>
            <div>
              <strong>{{ "address" | translate }}</strong>
              <p>{{ contact.address | translate }}</p>
            </div>
          </div>

          <div class="contact-item">
            <i class="icon-clock"></i>
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
              <i class="icon-phone"></i>
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
          <i class="icon-facebook"></i>
        </a>
        <a
          [href]="contact.socials.instagram"
          target="_blank"
          class="social-btn"
        >
          <i class="icon-instagram"></i>
        </a>
        <a
          [href]="contact.socials.linkedin"
          target="_blank"
          class="social-btn"
        >
          <i class="icon-linkedin"></i>
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
