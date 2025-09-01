// therapist-profile-overlay.component.ts
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';

export interface TherapistProfile {
  profile_id: string;
  name?: string;
  address: string;
  languages: string[];
  specialization: string[];
  consultation_modes: string[];
  experience: number;
  rating: number;
  currency?: string;
  availability: {
    days: string[];
    time_slots: string[];
  };
  avatar_override?: string;
  professional_info: {
    title: {
      code: string;
      display: string;
    };
    degrees?: string[] | null;
    certifications?: string[] | null;
    experience_years: number;
  };
}

@Component({
  selector: 'app-therapist-profile-overlay',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div
      class="therapist-overlay"
      [class.therapist-overlay--visible]="isVisible"
      (click)="onBackdropClick($event)"
    >
      <div class="therapist-card" (click)="$event.stopPropagation()">
        <header class="therapist-card__header">
          <button
            class="therapist-card__close-btn"
            (click)="onClose()"
            type="button"
          >
            <svg
              width="24"
              height="24"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <line x1="18" y1="6" x2="6" y2="18"></line>
              <line x1="6" y1="6" x2="18" y2="18"></line>
            </svg>
          </button>

          <div class="therapist-card__profile">
            <div class="therapist-card__avatar">
              <svg
                width="32"
                height="32"
                viewBox="0 0 24 24"
                fill="currentColor"
              >
                <path
                  d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"
                />
              </svg>
            </div>
            <div class="therapist-card__info">
              <h2 class="therapist-card__name">
                {{ therapist?.name || 'Therapist Name' }}
              </h2>
              <p class="therapist-card__title">
                {{
                  therapist?.professional_info?.title?.display ||
                    'Licensed Therapist'
                }}
              </p>
              <div class="therapist-card__rating">
                <svg
                  class="therapist-card__star"
                  width="16"
                  height="16"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <polygon
                    points="12,2 15.09,8.26 22,9.27 17,14.14 18.18,21.02 12,17.77 5.82,21.02 7,14.14 2,9.27 8.91,8.26"
                  />
                </svg>
                <span class="therapist-card__rating-value">{{
                  therapist?.rating || 0
                }}</span>
                <span class="therapist-card__experience"
                  >{{ therapist?.experience || 0 }} years experience</span
                >
              </div>
            </div>
          </div>
        </header>

        <div class="therapist-card__content">
          <section class="therapist-info">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path>
                <circle cx="12" cy="10" r="3"></circle>
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Location</h3>
              <p class="therapist-info__text">{{ therapist?.address }}</p>
            </div>
          </section>

          <section class="therapist-info">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="m5 8 6 6"></path>
                <path d="m4 14 6-6 2-3"></path>
                <path d="M2 5h12"></path>
                <path d="M7 2h1"></path>
                <path d="m22 22-5-10-5 10"></path>
                <path d="M14 18h6"></path>
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Languages</h3>
              <div class="therapist-tags">
                <span
                  class="therapist-tag therapist-tag--primary"
                  *ngFor="let language of therapist?.languages"
                >
                  {{ language }}
                </span>
              </div>
            </div>
          </section>

          <section class="therapist-info">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  d="M9.5 2A2.5 2.5 0 0 1 12 4.5v15a2.5 2.5 0 0 1-4.96.44 2.5 2.5 0 0 1-2.96-3.08 3 3 0 0 1-.34-5.58 2.5 2.5 0 0 1 1.32-4.24 2.5 2.5 0 0 1 1.98-3A2.5 2.5 0 0 1 9.5 2Z"
                />
                <path
                  d="M14.5 2A2.5 2.5 0 0 0 12 4.5v15a2.5 2.5 0 0 0 4.96.44 2.5 2.5 0 0 0 2.96-3.08 3 3 0 0 0 .34-5.58 2.5 2.5 0 0 0-1.32-4.24 2.5 2.5 0 0 0-1.98-3A2.5 2.5 0 0 0 14.5 2Z"
                />
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Specializations</h3>
              <div class="therapist-tags">
                <span
                  class="therapist-tag therapist-tag--success"
                  *ngFor="let spec of therapist?.specialization"
                >
                  {{ spec }}
                </span>
              </div>
            </div>
          </section>

          <section class="therapist-info">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path
                  d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
                ></path>
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Consultation Options</h3>
              <div class="therapist-tags">
                <span
                  class="therapist-tag therapist-tag--purple"
                  *ngFor="let mode of therapist?.consultation_modes"
                >
                  {{ mode }}
                </span>
              </div>
            </div>
          </section>

          <section class="therapist-info" *ngIf="hasQualifications()">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <path d="M6 9H4.5a2.5 2.5 0 0 1 0-5H6"></path>
                <path d="M18 9h1.5a2.5 2.5 0 0 0 0-5H18"></path>
                <path d="M4 22h16"></path>
                <path
                  d="M10 14.66V17c0 .55.47.98.97 1.21C12.15 18.75 14 20 16 20s3.85-1.25 5.03-1.79c.5-.23.97-.66.97-1.21v-2.34"
                ></path>
                <path d="M18 10h-1a4 4 0 0 1-4-4V4a2 2 0 1 1 4 0v2"></path>
                <path d="M6 10h1a4 4 0 0 0 4-4V4a2 2 0 0 0-4 0v2"></path>
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Qualifications</h3>
              <div class="therapist-qualifications">
                <div
                  class="therapist-qualifications__section"
                  *ngIf="therapist?.professional_info?.degrees"
                >
                  <p class="therapist-qualifications__label">Education:</p>
                  <ul class="therapist-qualifications__list">
                    <li
                      class="therapist-qualifications__item"
                      *ngFor="
                        let degree of therapist?.professional_info?.degrees
                      "
                    >
                      <svg
                        width="14"
                        height="14"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                      >
                        <path
                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      {{ degree }}
                    </li>
                  </ul>
                </div>
                <div
                  class="therapist-qualifications__section"
                  *ngIf="therapist?.professional_info?.certifications"
                >
                  <p class="therapist-qualifications__label">Certifications:</p>
                  <ul class="therapist-qualifications__list">
                    <li
                      class="therapist-qualifications__item"
                      *ngFor="
                        let cert of therapist?.professional_info?.certifications
                      "
                    >
                      <svg
                        width="14"
                        height="14"
                        viewBox="0 0 24 24"
                        fill="currentColor"
                      >
                        <path
                          d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z"
                        />
                      </svg>
                      {{ cert }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </section>

          <section class="therapist-info">
            <div class="therapist-info__icon">
              <svg
                width="18"
                height="18"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
              >
                <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
                <line x1="16" y1="2" x2="16" y2="6"></line>
                <line x1="8" y1="2" x2="8" y2="6"></line>
                <line x1="3" y1="10" x2="21" y2="10"></line>
              </svg>
            </div>
            <div class="therapist-info__content">
              <h3 class="therapist-info__title">Availability</h3>
              <div class="therapist-availability">
                <div class="therapist-availability__days">
                  <p class="therapist-availability__label">Available Days:</p>
                  <div class="therapist-tags">
                    <span
                      class="therapist-tag therapist-tag--neutral"
                      *ngFor="let day of therapist?.availability?.days"
                    >
                      {{ day }}
                    </span>
                  </div>
                </div>
                <div class="therapist-availability__times">
                  <p class="therapist-availability__label">Time Slots:</p>
                  <ul class="therapist-availability__list">
                    <li
                      class="therapist-availability__item"
                      *ngFor="let slot of therapist?.availability?.time_slots"
                    >
                      <svg
                        width="14"
                        height="14"
                        viewBox="0 0 24 24"
                        fill="none"
                        stroke="currentColor"
                        stroke-width="2"
                      >
                        <circle cx="12" cy="12" r="10"></circle>
                        <polyline points="12,6 12,12 16,14"></polyline>
                      </svg>
                      {{ slot }}
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </section>
        </div>

        <footer class="therapist-card__actions">
          <button
            class="therapist-btn therapist-btn--primary"
            (click)="onBookAppointment()"
            type="button"
          >
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect>
              <line x1="16" y1="2" x2="16" y2="6"></line>
              <line x1="8" y1="2" x2="8" y2="6"></line>
              <line x1="3" y1="10" x2="21" y2="10"></line>
            </svg>
            Book Appointment
          </button>
          <button
            class="therapist-btn therapist-btn--secondary"
            (click)="onSendMessage()"
            type="button"
          >
            <svg
              width="16"
              height="16"
              viewBox="0 0 24 24"
              fill="none"
              stroke="currentColor"
              stroke-width="2"
            >
              <path
                d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"
              ></path>
            </svg>
            Send Message
          </button>
        </footer>
      </div>
    </div>
  `,
  styleUrls: ['./profile-overlay.scss'],
})
export class TherapistProfileOverlayComponent implements OnInit {
  @Input() therapist: TherapistProfile | null = null;
  @Input() isVisible: boolean = false;

  @Output() closeOverlay = new EventEmitter<void>();
  @Output() bookAppointment = new EventEmitter<TherapistProfile>();
  @Output() sendMessage = new EventEmitter<TherapistProfile>();

  ngOnInit(): void {
    if (!this.therapist) {
      this.therapist = {
        profile_id: 't001',
        name: 'Dr. Sarah Johnson',
        address: '123 Main Street, City, Country',
        languages: ['English', 'Spanish'],
        specialization: ['Cognitive Behavioral Therapy', 'Anxiety Management'],
        consultation_modes: ['Online', 'In-Person'],
        experience: 10,
        rating: 4.8,
        availability: {
          days: ['Monday', 'Wednesday', 'Friday'],
          time_slots: ['10:00 AM - 12:00 PM', '2:00 PM - 4:00 PM'],
        },
        professional_info: {
          title: {
            code: 'LCP',
            display: 'Licensed Clinical Psychologist',
          },
          degrees: [
            'Ph.D. in Clinical Psychology',
            'M.A. in Counseling Psychology',
          ],
          certifications: ['CBT Certified', 'Anxiety Disorders Specialist'],
          experience_years: 10,
        },
      };
    }
  }

  onClose(): void {
    console.log('Overlay close clicked');
    this.closeOverlay.emit();
  }

  onBackdropClick(event: Event): void {
    if (event.target === event.currentTarget) {
      this.onClose();
    }
  }

  onBookAppointment(): void {
    if (this.therapist) {
      this.bookAppointment.emit(this.therapist);
    }
  }

  onSendMessage(): void {
    if (this.therapist) {
      this.sendMessage.emit(this.therapist);
    }
  }

  hasQualifications(): boolean {
    return !!(
      this.therapist?.professional_info?.degrees?.length ||
      this.therapist?.professional_info?.certifications?.length
    );
  }
}
