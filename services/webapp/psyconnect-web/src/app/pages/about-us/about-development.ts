import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { Dev } from './about-us-data-service';

interface Release {
  version: string;
  date: string;
  notes: string[];
  upcomingFeatures: string[];
}

@Component({
  selector: 'app-development-section',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <section class="about-us__development container">
      <div class="release-card">
        <div class="release-header">
          <h3>{{ "current_version" | translate }}: {{ release.version }}</h3>
          <div class="release-date">{{ release.date | date }}</div>
        </div>

        <div class="release-content">
          <div class="release-notes">
            <h4>{{ "whats_new" | translate }}</h4>
            <ul>
              <li *ngFor="let note of release.notes">
                {{ note | translate }}
              </li>
            </ul>
          </div>

          <div class="upcoming-features">
            <h4>{{ "coming_soon" | translate }}</h4>
            <ul>
              <li *ngFor="let feature of release.upcomingFeatures">
                {{ feature | translate }}
              </li>
            </ul>
          </div>
        </div>

        <div class="dev-team">
          <h4>{{ "development_team" | translate }}</h4>
          <div class="dev-grid">
            <div *ngFor="let d of devs" class="dev-card">
              <img class="dev-avatar" [src]="d.avatar" [alt]="d.name" />
              <div class="dev-info">
                <strong>{{ d.name }}</strong>
                <p class="dev-role">{{ d.role | translate }}</p>
                <p class="dev-bio">{{ d.bio | translate }}</p>
                <div class="dev-links">
                  <a
                    *ngIf="d.github"
                    [href]="'https://' + d.github"
                    target="_blank"
                    class="link"
                  >
                    <i class="icon-github"></i>
                  </a>
                  <a
                    *ngIf="d.linkedin"
                    [href]="'https://' + d.linkedin"
                    target="_blank"
                    class="link"
                  >
                    <i class="icon-linkedin"></i>
                  </a>
                </div>
              </div>
            </div>
          </div>
        </div>

        <div class="release-cta">
          <button class="btn btn--ghost" (click)="viewRepo.emit()">
            <i class="icon-github"></i>
            {{ "view_repository" | translate }}
          </button>
        </div>
      </div>
    </section>
  `,
  styleUrls: ['./about-development.scss'],
})
export class DevelopmentSectionComponent {
  @Input() devs: Dev[] = [];
  @Input() release!: Release;
  @Output() viewRepo = new EventEmitter<void>();
}
