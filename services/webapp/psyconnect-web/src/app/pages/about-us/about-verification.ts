import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-verification-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__verification container">
      <h2 class="section-title">{{ "expert_verification" | translate }}</h2>
      <p class="section-desc">{{ "verification_description" | translate }}</p>

      <div class="verification-grid">
        <div
          class="verification-step"
          *ngFor="let step of verificationSteps; let i = index"
        >
          <div class="step-number">{{ i + 1 }}</div>
          <div class="step-content">
            <h4>{{ step | translate }}</h4>
          </div>
        </div>
      </div>

      <div class="quality-assurance">
        <h3 class="subhead">{{ "quality_standards" | translate }}</h3>
        <ul class="quality-list">
          <li *ngFor="let standard of qualityStandards">
            <i class="icon-check"></i>
            {{ standard | translate }}
          </li>
        </ul>
      </div>
    </section>
  `,styleUrls: ['./about-verification.scss'],
})
export class VerificationSectionComponent {
  @Input() verificationSteps: string[] = [];
  @Input() qualityStandards: string[] = [];
}
