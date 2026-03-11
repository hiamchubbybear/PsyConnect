import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

import { IconComponent } from '../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-privacy-section',
  standalone: true,
  imports: [CommonModule, TranslateModule, IconComponent],
  template: `
    <section class="about-us__privacy container">
      <h2 class="section-title">{{ "privacy_security" | translate }}</h2>

      <div class="privacy-grid">
        <div class="privacy-card">
          <h3 class="subhead">{{ "privacy_protection" | translate }}</h3>
          <ul>
            <li *ngFor="let feature of privacyFeatures">
              <app-icon name="shield" [size]="18"></app-icon>
              {{ feature | translate }}
            </li>
          </ul>
        </div>

        <div class="security-card">
          <h3 class="subhead">{{ "security_measures" | translate }}</h3>
          <ul>
            <li *ngFor="let measure of securityMeasures">
              <app-icon name="lock" [size]="18"></app-icon>
              {{ measure | translate }}
            </li>
          </ul>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-privacy-section.scss'],
})
export class PrivacySectionComponent {
  @Input() privacyFeatures: string[] = [];
  @Input() securityMeasures: string[] = [];
}
