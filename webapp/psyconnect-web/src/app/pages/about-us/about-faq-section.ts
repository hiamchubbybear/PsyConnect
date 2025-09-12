import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-faq-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__faq container">
      <h2 class="section-title">{{ "frequently_asked" | translate }}</h2>

      <div class="faq-list">
        <details class="faq-item" *ngFor="let item of faq">
          <summary class="faq-question">
            {{ item.question | translate }}
          </summary>
          <div class="faq-answer">{{ item.answer | translate }}</div>
        </details>
      </div>
    </section>
  `,
  styleUrls: ['./about-faq-section.scss'],
})
export class FaqSectionComponent {
  @Input() faq: { question: string; answer: string }[] = [];
}
