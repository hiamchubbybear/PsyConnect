import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { AuthService } from '../../../services/auth/auth.service';
import { ConsultationProfileService } from '../../../services/consultation/consultation-profile.service';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-create-consultation-profile',
  standalone: true,
  imports: [CommonModule, FormsModule, ReactiveFormsModule, TranslateModule],
  template: `
    <div class="create-profile-layout fade-in">
      <div class="profile-wizard-container">
        
        <!-- Left Sidebar: Steps -->
        <aside class="wizard-sidebar">
          <div class="sidebar-header">
            <h3>{{ 'CONSULTATION.CreateProfile.Title' | translate }}</h3>
            <span class="role-badge">{{ role }}</span>
          </div>

          <div class="steps-progress">
            <div class="step-item" [class.active]="currentStep === 1" [class.completed]="currentStep > 1" (click)="goToStep(1)">
              <div class="step-indicator">
                <svg *ngIf="currentStep > 1" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span *ngIf="currentStep <= 1">1</span>
              </div>
              <span class="step-label">Configuration</span>
            </div>
            
            <div class="step-item" [class.active]="currentStep === 2" [class.completed]="currentStep > 2" (click)="goToStep(2)">
              <div class="step-indicator">
                <svg *ngIf="currentStep > 2" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span *ngIf="currentStep <= 2">2</span>
              </div>
              <span class="step-label">Attributes</span>
            </div>

            <div class="step-item" [class.active]="currentStep === 3" [class.completed]="currentStep > 3" (click)="goToStep(3)">
              <div class="step-indicator">
                <span *ngIf="currentStep <= 3">3</span>
              </div>
              <span class="step-label">Details</span>
            </div>

            <!-- Animated sliding highlight for active step -->
            <div class="step-highlight" [style.transform]="'translateY(' + ((currentStep - 1) * 60) + 'px)'"></div>
          </div>
        </aside>

        <!-- Right Content Area -->
        <main class="wizard-content">
          <form [formGroup]="profileForm" (ngSubmit)="createProfile()">
            <div class="content-header">
              <h2>{{ getStepTitle() }}</h2>
              <div class="header-actions">
                <button type="button" class="btn-cancel" (click)="goBack()">Cancel</button>
                <button *ngIf="currentStep < 3" type="button" class="btn-next" (click)="nextStep()">Continue &rsaquo;</button>
                <button *ngIf="currentStep === 3" type="submit" class="btn-next" [disabled]="profileForm.invalid || creating">
                  {{ creating ? 'Saving...' : 'Save & Finish' }}
                </button>
              </div>
            </div>

            <!-- Step 1: Configuration (Modes & Price) -->
            <div class="step-pane" [class.active-pane]="currentStep === 1">
              <div class="section-group" style="--anim-delay: 1">
                <h4>Consultation Modes</h4>
                <p class="section-desc">Choose how you would like to conduct consultations.</p>
                <div class="card-grid">
                  <div *ngFor="let mode of availableModes" 
                       class="selectable-card" 
                       [class.selected]="isModeSelected(mode.id)"
                       (click)="toggleMode(mode.id)">
                    <div class="selection-ring"></div>
                    <div class="card-icon" [innerHTML]="mode.icon"></div>
                    <span class="card-title">{{ mode.label }}</span>
                  </div>
                </div>
              </div>

              <div class="section-group" style="--anim-delay: 2">
                <h4>Budget / Rate (per hour)</h4>
                <p class="section-desc">Enter your expected {{ role === 'therapist' ? 'hourly rate' : 'budget' }} in USD.</p>
                <input type="number" formControlName="rage_price" class="form-input" placeholder="Amount (e.g. 50)" />
              </div>
            </div>

            <!-- Step 2: Attributes (Languages & Specializations) -->
            <div class="step-pane" [class.active-pane]="currentStep === 2">
              <div class="section-group" style="--anim-delay: 1">
                <h4>Languages</h4>
                <p class="section-desc">Select the languages you are comfortable with.</p>
                <div class="card-grid">
                  <div *ngFor="let lang of availableLanguages" 
                       class="selectable-card" 
                       [class.selected]="isLanguageSelected(lang.id)"
                       (click)="toggleLanguage(lang.id)">
                    <div class="selection-ring"></div>
                    <span class="card-title">{{ lang.label }}</span>
                  </div>
                </div>
              </div>

              <div class="section-group" style="--anim-delay: 2">
                <h4>{{ role === 'therapist' ? 'Specializations' : 'Issues' }}</h4>
                <p class="section-desc">Select relevant areas of focus.</p>
                <div class="card-grid columns-3">
                  <div *ngFor="let topic of availableTopics" 
                       class="selectable-card" 
                       [class.selected]="isTopicSelected(topic.id)"
                       (click)="toggleTopic(topic.id)">
                    <div class="selection-ring"></div>
                    <span class="card-title">{{ topic.label }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 3: Details (Address & Experience) -->
            <div class="step-pane" [class.active-pane]="currentStep === 3">
              <div class="section-group" style="--anim-delay: 1">
                <h4>Location Details</h4>
                <p class="section-desc">Enter your current city or regional address.</p>
                <input type="text" formControlName="address" class="form-input" placeholder="City, Country" />
              </div>

              <div class="section-group" style="--anim-delay: 2" *ngIf="role === 'therapist'">
                <h4>Experience</h4>
                <p class="section-desc">Years of professional experience.</p>
                <input type="number" formControlName="experience" class="form-input" placeholder="e.g. 5" />
              </div>

              <div class="summary-box" style="--anim-delay: 3">
                <p>You are almost done! Review your information and click <strong>Save & Finish</strong> to complete your profile.</p>
              </div>
            </div>

          </form>
        </main>
      </div>
    </div>
  `,
  styles: [`
    .create-profile-layout {
      min-height: calc(100vh - 80px); /* Adjust based on navbar height */
      background: var(--color-bg-secondary);
      display: flex;
      justify-content: center;
      padding: 2rem;
    }

    .profile-wizard-container {
      display: flex;
      width: 100%;
      max-width: 1100px;
      background: var(--color-surface);
      border-radius: 12px;
      box-shadow: 0 10px 40px -10px rgba(0,0,0,0.08);
      overflow: hidden;
      border: 1px solid var(--color-border);
    }

    /* Sidebar Styles */
    .wizard-sidebar {
      width: 280px;
      background: rgba(var(--color-bg-rgb), 0.5);
      border-right: 1px solid var(--color-border);
      padding: 2rem 0;
      display: flex;
      flex-direction: column;
    }

    .sidebar-header {
      padding: 0 2rem 2rem;
      border-bottom: 1px solid var(--color-border);
      margin-bottom: 1.5rem;

      h3 {
        font-size: 1.25rem;
        font-weight: 600;
        margin: 0 0 0.5rem 0;
        color: var(--color-text);
      }

      .role-badge {
        display: inline-block;
        padding: 0.25rem 0.75rem;
        background: var(--color-primary);
        color: #fff;
        border-radius: 999px;
        font-size: 0.75rem;
        font-weight: 600;
        text-transform: uppercase;
        letter-spacing: 0.5px;
      }
    }

    .steps-progress {
      position: relative;
      padding: 0 1rem;
    }

    .step-item {
      display: flex;
      align-items: center;
      gap: 1rem;
      padding: 1rem;
      border-radius: 8px;
      cursor: pointer;
      position: relative;
      z-index: 2;
      transition: all 0.3s ease;
      height: 60px;
      color: var(--color-text-muted);

      &:hover {
        background: rgba(var(--color-text-rgb), 0.04);
      }

      &.active {
        color: var(--color-primary);
        font-weight: 600;
      }

      &.completed {
        color: var(--color-text);
      }
    }

    .step-indicator {
      width: 28px;
      height: 28px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      background: var(--color-bg-secondary);
      border: 2px solid var(--color-border);
      font-size: 0.875rem;
      font-weight: 600;
      transition: all 0.3s ease;

      .step-item.active & {
        background: var(--color-primary);
        border-color: var(--color-primary);
        color: #fff;
      }
      .step-item.completed & {
        background: var(--color-success, #10b981);
        border-color: var(--color-success, #10b981);
        color: #fff;
      }
    }

    .step-highlight {
      position: absolute;
      top: 0;
      left: 1rem;
      right: 1rem;
      height: 60px;
      background: rgba(var(--color-primary-rgb, 102,126,234), 0.08); /* Fallback */
      border-radius: 8px;
      z-index: 1;
      transition: transform 0.4s cubic-bezier(0.25, 1, 0.5, 1);
      pointer-events: none;
    }

    /* Main Content Area */
    .wizard-content {
      flex: 1;
      padding: 2.5rem;
      display: flex;
      flex-direction: column;
      position: relative;
      overflow: hidden;
    }

    .content-header {
      display: flex;
      align-items: center;
      justify-content: space-between;
      margin-bottom: 2.5rem;
      padding-bottom: 1rem;
      border-bottom: 1px solid var(--color-border);

      h2 {
        font-size: 1.5rem;
        font-weight: 600;
        margin: 0;
      }

      .header-actions {
        display: flex;
        gap: 0.75rem;
      }
    }

    .btn-cancel {
      padding: 0.5rem 1rem;
      background: transparent;
      border: 1px solid var(--color-border);
      border-radius: 6px;
      font-weight: 500;
      color: var(--color-text);
      cursor: pointer;
      transition: all 0.2s;
      
      &:hover {
        background: var(--color-bg-secondary);
      }
    }

    .btn-next {
      padding: 0.5rem 1.25rem;
      background: var(--color-primary);
      color: #fff;
      border: none;
      border-radius: 6px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;

      &:hover:not(:disabled) {
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(var(--color-primary-rgb, 102,126,234), 0.3);
      }
      &:active:not(:disabled) {
        transform: translateY(0) scale(0.98);
      }
      &:disabled {
        opacity: 0.6;
        cursor: not-allowed;
      }
    }

    /* Steps Animation */
    .step-pane {
      display: none;
      flex-direction: column;
      gap: 2.5rem;
      animation: slideIn 0.4s ease forwards;

      &.active-pane {
        display: flex;
      }
    }

    @keyframes slideIn {
      from { opacity: 0; transform: translateX(20px); }
      to { opacity: 1; transform: translateX(0); }
    }

    /* Section Groups & Cards */
    .section-group {
      opacity: 0;
      animation: fadeUpIn 0.5s ease forwards;
      /* use inline style --anim-delay to stagger */
      animation-delay: calc(var(--anim-delay) * 0.1s);

      h4 {
        font-size: 1.125rem;
        font-weight: 600;
        margin: 0 0 0.25rem 0;
      }
      .section-desc {
        font-size: 0.875rem;
        color: var(--color-text-muted);
        margin: 0 0 1.25rem 0;
      }
    }

    @keyframes fadeUpIn {
      from { opacity: 0; transform: translateY(10px); }
      to { opacity: 1; transform: translateY(0); }
    }

    .card-grid {
      display: grid;
      grid-template-columns: repeat(4, 1fr);
      gap: 1rem;

      &.columns-3 {
        grid-template-columns: repeat(3, 1fr);
      }
    }

    .selectable-card {
      position: relative;
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: 8px;
      padding: 1.25rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      gap: 0.75rem;
      cursor: pointer;
      transition: all 0.2s cubic-bezier(0.2, 0.8, 0.2, 1);
      user-select: none;
      min-height: 100px;
      text-align: center;

      &:hover {
        transform: translateY(-2px);
        box-shadow: 0 6px 16px rgba(0,0,0,0.06);
        border-color: rgba(var(--color-text-rgb), 0.2);
      }

      &:active {
        transform: translateY(0) scale(0.97);
      }

      &.selected {
        border-color: var(--color-primary);
        background: rgba(var(--color-primary-rgb, 102,126,234), 0.04);
        box-shadow: 0 0 0 1px var(--color-primary);

        .selection-ring {
          border-color: var(--color-primary);
          &::after {
            content: '';
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            width: 8px;
            height: 8px;
            border-radius: 50%;
            background: var(--color-primary);
            animation: popIn 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275) forwards;
          }
        }
      }
    }

    @keyframes popIn {
      0% { transform: translate(-50%, -50%) scale(0); }
      100% { transform: translate(-50%, -50%) scale(1); }
    }

    .selection-ring {
      position: absolute;
      top: 0.75rem;
      left: 0.75rem;
      width: 16px;
      height: 16px;
      border-radius: 50%;
      border: 1.5px solid var(--color-border);
      transition: all 0.2s;
    }

    .card-icon {
      color: var(--color-text-muted);
      transition: color 0.2s;
      display: flex;
      
      ::ng-deep svg {
        width: 24px;
        height: 24px;
      }

      .selected & {
        color: var(--color-primary);
      }
    }

    .card-title {
      font-size: 0.875rem;
      font-weight: 500;
      color: var(--color-text);
    }

    .form-input {
      width: 100%;
      max-width: 400px;
      padding: 0.75rem 1rem;
      border: 1px solid var(--color-border);
      border-radius: 6px;
      background: var(--color-bg);
      color: var(--color-text);
      font-family: inherit;
      font-size: 0.9375rem;
      transition: border-color 0.2s, box-shadow 0.2s;

      &:focus {
        outline: none;
        border-color: var(--color-primary);
        box-shadow: 0 0 0 3px rgba(var(--color-primary-rgb, 102,126,234), 0.15);
      }
    }

    .summary-box {
      padding: 1.5rem;
      border-radius: 8px;
      background: rgba(var(--color-success-rgb, 16, 185, 129), 0.1);
      border: 1px solid rgba(var(--color-success-rgb, 16, 185, 129), 0.2);
      
      p {
        margin: 0;
        color: var(--color-text);
        font-size: 0.9375rem;
      }
    }

    /* Responsive adjustments */
    @media (max-width: 768px) {
      .profile-wizard-container {
        flex-direction: column;
      }
      .wizard-sidebar {
        width: 100%;
        border-right: none;
        border-bottom: 1px solid var(--color-border);
        padding: 1.5rem;
      }
      .step-highlight {
        display: none; /* Simplify on mobile for now */
      }
      .steps-progress {
        display: flex;
        justify-content: space-around;
        padding: 0;
      }
      .step-item {
        flex-direction: column;
        gap: 0.5rem;
        padding: 0.5rem;
        height: auto;
      }
      .step-label {
        font-size: 0.75rem;
      }
      .card-grid {
        grid-template-columns: repeat(2, 1fr) !important;
      }
      .wizard-content {
        padding: 1.5rem;
      }
    }
  `]
})
export class CreateConsultationProfileComponent implements OnInit {
  role: 'therapist' | 'client' | null = null;
  creating = false;
  profileForm: FormGroup;
  currentStep = 1;

  // Data for cards
  availableModes = [
    { id: 'online', label: 'Online', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>' },
    { id: 'in-person', label: 'In-person', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>' },
    { id: 'chat', label: 'Chat', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>' },
    { id: 'phone', label: 'Phone', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"></rect><line x1="12" y1="18" x2="12.01" y2="18"></line></svg>' }
  ];

  availableLanguages = [
    { id: 'en', label: 'English' },
    { id: 'vi', label: 'Vietnamese' },
    { id: 'zh', label: 'Chinese' },
    { id: 'ja', label: 'Japanese' }
  ];

  /* Example topics for specialization/issues */
  availableTopics = [
    { id: 'anxiety', label: 'Anxiety' },
    { id: 'depression', label: 'Depression' },
    { id: 'stress', label: 'Stress' },
    { id: 'relationship', label: 'Relationships' },
    { id: 'trauma', label: 'Trauma' },
    { id: 'career', label: 'Career' },
  ];

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private authService: AuthService,
    private consultationProfileService: ConsultationProfileService,
    private toastService: ToastService,
    private fb: FormBuilder
  ) {
    this.profileForm = this.fb.group({
      address: ['', Validators.required],
      rage_price: [null, [Validators.required, Validators.min(0)]],
      languages: [[]],
      consultation_modes: [[]],
      issue_detail: [[]], // Client
      specialization: [[]], // Therapist
      experience: [null]
    });
  }

  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      this.role = params['role'] || this.authService.getRole();
      if (!this.role) {
        this.toastService.error('Error', 'No role found. Please login again.');
        this.router.navigate(['/login']);
      }
    });
  }

  // --- Step Navigation ---
  getStepTitle(): string {
    switch (this.currentStep) {
      case 1: return 'Configuration';
      case 2: return 'Attributes';
      case 3: return 'Details';
      default: return '';
    }
  }

  nextStep() {
    // Optional: add step validation here
    if (this.currentStep < 3) {
      this.currentStep++;
    }
  }

  goToStep(step: number) {
    // Optional: only allow jumping if previous steps are valid
    this.currentStep = step;
  }

  goBack() {
    this.router.navigate(['/']);
  }

  // --- Array Control Helpers (for Cards) ---
  toggleArrayValue(controlName: string, value: string) {
    const control = this.profileForm.get(controlName);
    if (!control) return;
    const currentList: string[] = control.value || [];
    const index = currentList.indexOf(value);
    
    if (index > -1) {
      currentList.splice(index, 1);
    } else {
      currentList.push(value);
    }
    control.setValue([...currentList]);
  }

  isArraySelected(controlName: string, value: string): boolean {
    const control = this.profileForm.get(controlName);
    if (!control) return false;
    return (control.value || []).includes(value);
  }

  toggleMode(id: string) { this.toggleArrayValue('consultation_modes', id); }
  isModeSelected(id: string) { return this.isArraySelected('consultation_modes', id); }

  toggleLanguage(id: string) { this.toggleArrayValue('languages', id); }
  isLanguageSelected(id: string) { return this.isArraySelected('languages', id); }

  toggleTopic(id: string) {
    if (this.role === 'therapist') {
      this.toggleArrayValue('specialization', id);
    } else {
      this.toggleArrayValue('issue_detail', id);
    }
  }
  isTopicSelected(id: string) {
    return this.role === 'therapist' 
      ? this.isArraySelected('specialization', id)
      : this.isArraySelected('issue_detail', id);
  }

  // --- Submit ---
  createProfile() {
    // Ensure lists aren't empty
    if (!this.profileForm.value.languages?.length) {
      this.toastService.error('Validation', 'Please select at least one language.');
      this.currentStep = 2; return;
    }
    if (!this.profileForm.value.consultation_modes?.length) {
      this.toastService.error('Validation', 'Please select at least one consultation mode.');
      this.currentStep = 1; return;
    }

    if (!this.role || this.profileForm.invalid) return;

    this.creating = true;
    const formVal = this.profileForm.value;

    const baseData = {
      address: formVal.address,
      rage_price: Number(formVal.rage_price),
      languages: formVal.languages,
      consultation_modes: formVal.consultation_modes,
      availability: {
        days: ['monday', 'tuesday', 'wednesday', 'thursday', 'friday'],
        time_slots: ['morning', 'afternoon']
      }
    };

    let profileData: any = { ...baseData };

    if (this.role === 'therapist') {
      profileData.experience = Number(formVal.experience) || 0;
      profileData.specialization = formVal.specialization;
    } else {
      profileData.issue_detail = formVal.issue_detail;
    }

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
}
