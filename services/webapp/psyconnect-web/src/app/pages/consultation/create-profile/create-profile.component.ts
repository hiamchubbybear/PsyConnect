import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { AuthService } from '../../../services/auth/auth.service';
import { ConsultationProfileService } from '../../../services/consultation/consultation-profile.service';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-create-consultation-profile',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  template: `
    <div class="create-profile-container">
      <div class="create-profile-card">
        <h1>{{ 'CONSULTATION.CreateProfile.Title' | translate }}</h1>
        <p class="subtitle">{{ 'CONSULTATION.CreateProfile.Subtitle' | translate }}</p>

        <div class="role-info">
          <i class="fas fa-user-circle"></i>
          <span>{{ 'CONSULTATION.CreateProfile.Role' | translate }}: <strong>{{ role }}</strong></span>
        </div>

        <div class="info-message">
          <i class="fas fa-info-circle"></i>
          <p>{{ 'CONSULTATION.CreateProfile.Message' | translate }}</p>
        </div>

        <div class="actions">
          <button class="btn-primary" (click)="createProfile()" [disabled]="creating">
            <i class="fas fa-plus"></i>
            {{ creating ? ('COMMON.Creating' | translate) : ('CONSULTATION.CreateProfile.Create' | translate) }}
          </button>
          <button class="btn-secondary" (click)="goBack()">
            <i class="fas fa-arrow-left"></i>
            {{ 'COMMON.Back' | translate }}
          </button>
        </div>
      </div>
    </div>
  `,
  styles: [`
    .create-profile-container {
      display: flex;
      align-items: center;
      justify-content: center;
      min-height: 100vh;
      padding: 2rem;
      background: var(--color-bg-secondary);
    }

    .create-profile-card {
      max-width: 600px;
      width: 100%;
      background: var(--color-surface);
      border-radius: 8px;
      padding: 2rem;
      box-shadow: var(--shadow-card);
    }

    h1 {
      font-size: 1.75rem;
      font-weight: 600;
      color: var(--color-text);
      margin: 0 0 0.5rem 0;
    }

    .subtitle {
      color: var(--color-text-muted);
      margin: 0 0 2rem 0;
    }

    .role-info {
      display: flex;
      align-items: center;
      gap: 0.75rem;
      padding: 1rem;
      background: var(--color-accent-bg);
      border-radius: 6px;
      margin-bottom: 1.5rem;

      i {
        font-size: 1.5rem;
        color: var(--color-accent);
      }

      span {
        color: var(--color-text);
        font-size: 0.9375rem;

        strong {
          text-transform: capitalize;
          color: var(--color-accent);
        }
      }
    }

    .info-message {
      display: flex;
      gap: 0.75rem;
      padding: 1rem;
      background: rgba(59, 130, 246, 0.1);
      border-left: 4px solid #3b82f6;
      border-radius: 6px;
      margin-bottom: 2rem;

      i {
        color: #3b82f6;
        font-size: 1.25rem;
        flex-shrink: 0;
      }

      p {
        margin: 0;
        color: var(--color-text);
        font-size: 0.875rem;
        line-height: 1.5;
      }
    }

    .actions {
      display: flex;
      gap: 1rem;

      button {
        flex: 1;
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.5rem;
        padding: 0.75rem 1.5rem;
        border: none;
        border-radius: 6px;
        font-size: 0.9375rem;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;

        &:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      }

      .btn-primary {
        background: var(--color-primary);
        color: var(--color-on-accent);

        &:hover:not(:disabled) {
          background: var(--color-primary-hover);
        }
      }

      .btn-secondary {
        background: transparent;
        border: 1px solid var(--color-border);
        color: var(--color-text);

        &:hover {
          background: var(--color-bg-secondary);
        }
      }
    }
  `]
})
export class CreateConsultationProfileComponent implements OnInit {
  role: 'therapist' | 'client' | null = null;
  creating = false;

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private authService: AuthService,
    private consultationProfileService: ConsultationProfileService,
    private toastService: ToastService
  ) {}

  ngOnInit() {
    // Get role from query params or auth service
    this.route.queryParams.subscribe(params => {
      this.role = params['role'] || this.authService.getRole();

      if (!this.role) {
        this.toastService.error('Error', 'No role found. Please login again.');
        this.router.navigate(['/login']);
      }
    });
  }

  createProfile() {
    if (!this.role) {
      return;
    }

    this.creating = true;

    const profileData = {
      // Add minimal required data
      // Backend should create profile with user info from JWT
    };

    const createMethod = this.role === 'therapist'
      ? this.consultationProfileService.createTherapistProfile(profileData)
      : this.consultationProfileService.createClientProfile(profileData);

    createMethod.subscribe({
      next: (profile: any) => {
        this.toastService.success('Success', 'Profile created successfully!');
        this.router.navigate(['/']);
      },
      error: (error: any) => {
        console.error('Failed to create profile:', error);
        this.toastService.error('Error', 'Failed to create profile. Please try again.');
        this.creating = false;
      }
    });
  }

  goBack() {
    this.router.navigate(['/']);
  }
}
