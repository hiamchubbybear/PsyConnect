import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-audience-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__audience container">
      <h2 class="section-title">{{ 'who_should_use' | translate }}</h2>
      <div class="audience-grid">
        <div *ngFor="let a of targetAudience" class="audience-card">
          <i class="icon-users"></i>
          <h4>{{ a | translate }}</h4>
        </div>
      </div>
      <p class="lead">{{ 'getting_started_message' | translate }}</p>
    </section>
  `,
  styleUrls: ['./about-audience.scss'],
})
export class AudienceSectionComponent {
  @Input() targetAudience: string[] = [];
}
