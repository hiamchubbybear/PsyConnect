import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { TeamMember } from './about-us-data-service';

import { IconComponent } from '../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-team-section',
  standalone: true,
  imports: [CommonModule, TranslateModule, IconComponent],
  template: `
    <section class="about-us__team container">
      <h2 class="section-title">{{ "our_team" | translate }}</h2>
      <p class="section-desc">{{ "team_description" | translate }}</p>

      <div class="team-grid">
        <div *ngFor="let member of team" class="team-card">
          <img
            class="team-card__avatar"
            [src]="member.avatar"
            [alt]="member.name"
          />
          <div class="team-card__info">
            <h4 class="team-card__name">{{ member.name }}</h4>
            <p class="team-card__role">{{ member.role | translate }}</p>
            <p class="team-card__bio">{{ member.bio | translate }}</p>
            <div class="team-card__social" *ngIf="member.social">
              <a
                *ngIf="member.social.linkedin"
                [href]="member.social.linkedin"
                target="_blank"
              >
                <app-icon name="linkedin" [size]="18"></app-icon>
              </a>
              <a
                *ngIf="member.social.email"
                [href]="'mailto:' + member.social.email"
              >
                <app-icon name="mail" [size]="18"></app-icon>
              </a>
            </div>
          </div>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-team-section.scss'],
})
export class TeamSectionComponent {
  @Input() team: TeamMember[] = [];
}
