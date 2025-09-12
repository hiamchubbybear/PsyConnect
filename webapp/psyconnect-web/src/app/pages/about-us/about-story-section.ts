import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-story-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__story container">
      <div class="story-content">
        <h2 class="section-title">{{ storyTitle }}</h2>
        <p class="lead">{{ story }}</p>

        <div class="mission-vision">
          <div class="mission-card">
            <h3 class="subhead">{{ "our_mission" | translate }}</h3>
            <p>{{ mission }}</p>
          </div>
          <div class="vision-card">
            <h3 class="subhead">{{ "our_vision" | translate }}</h3>
            <p>{{ vision }}</p>
          </div>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-story-section.scss'],
})
export class StorySectionComponent {
  @Input() storyTitle: string = '';
  @Input() story: string = '';
  @Input() mission: string = '';
  @Input() vision: string = '';
}
