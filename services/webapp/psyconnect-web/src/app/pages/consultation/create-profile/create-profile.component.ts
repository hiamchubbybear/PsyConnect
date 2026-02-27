import { CommonModule } from '@angular/common';
import { HttpClient } from '@angular/common/http';
import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, FormsModule, ReactiveFormsModule, Validators } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Subject, of } from 'rxjs';
import { catchError, debounceTime, distinctUntilChanged, switchMap } from 'rxjs/operators';
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
              <span class="step-label">{{ 'CONSULTATION.CreateProfile.StepConfig' | translate }}</span>
            </div>
            
            <div class="step-item" [class.active]="currentStep === 2" [class.completed]="currentStep > 2" (click)="goToStep(2)">
              <div class="step-indicator">
                <svg *ngIf="currentStep > 2" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span *ngIf="currentStep <= 2">2</span>
              </div>
              <span class="step-label">{{ 'CONSULTATION.CreateProfile.StepAttributes' | translate }}</span>
            </div>

            <div class="step-item" [class.active]="currentStep === 3" [class.completed]="currentStep > 3" (click)="goToStep(3)">
              <div class="step-indicator">
                <svg *ngIf="currentStep > 3" viewBox="0 0 24 24" width="16" height="16" stroke="currentColor" stroke-width="3" fill="none" stroke-linecap="round" stroke-linejoin="round"><polyline points="20 6 9 17 4 12"></polyline></svg>
                <span *ngIf="currentStep <= 3">3</span>
              </div>
              <span class="step-label">{{ 'CONSULTATION.CreateProfile.StepDetails' | translate }}</span>
            </div>

            <div class="step-item" *ngIf="role === 'therapist'" [class.active]="currentStep === 4" [class.completed]="currentStep > 4" (click)="goToStep(4)">
              <div class="step-indicator">
                <span *ngIf="currentStep <= 4">4</span>
              </div>
              <span class="step-label">{{ 'CONSULTATION.THERAPIST.ProfessionalInfo' | translate }}</span>
            </div>

            <!-- Animated sliding highlight for active step -->
            <div class="step-highlight" [style.transform]="'translateY(' + ((currentStep - 1) * 60) + 'px)'"></div>
          </div>
        </aside>

        <!-- Right Content Area -->
        <main class="wizard-content">
          <form [formGroup]="profileForm" (ngSubmit)="createProfile()">
            <div class="content-header">
              <h2>{{ getStepTitle() | translate }}</h2>
              <div class="header-actions">
                <button type="button" class="btn-cancel" (click)="goBack()">{{ 'COMMON.Cancel' | translate }}</button>
                <button *ngIf="currentStep < (role === 'therapist' ? 4 : 3)" type="button" class="btn-next" (click)="nextStep()">{{ 'COMMON.Next' | translate }} &rsaquo;</button>
                <button *ngIf="currentStep === (role === 'therapist' ? 4 : 3)" type="submit" class="btn-next" [disabled]="profileForm.invalid || creating">
                  {{ creating ? ('COMMON.Saving' | translate) : (isEditing ? ('COMMON.Update' | translate) : ('COMMON.SaveProfile' | translate)) }}
                </button>
              </div>
            </div>

            <!-- Step 1: Configuration (Modes & Price) -->
            <div class="step-pane" [class.active-pane]="currentStep === 1">
              <div class="section-group" style="--anim-delay: 1">
                <h4>{{ 'CONSULTATION.THERAPIST.ConsultationModes' | translate }}</h4>
                <p class="section-desc">{{ 'CONSULTATION.CreateProfile.ConsultationModesDesc' | translate }}</p>
                <div class="card-grid">
                  <div *ngFor="let mode of availableModes" 
                       class="selectable-card" 
                       [class.selected]="isModeSelected(mode.id)"
                       (click)="toggleMode(mode.id)">
                    <div class="selection-ring"></div>
                    <div class="card-icon" [innerHTML]="mode.icon"></div>
                    <span class="card-title">{{ mode.label | translate }}</span>
                  </div>
                </div>
              </div>

              <div class="section-group" style="--anim-delay: 2">
                <h4>{{ 'CONSULTATION.THERAPIST.Price' | translate }}</h4>
                <p class="section-desc">{{ role === 'therapist' ? ('CONSULTATION.THERAPIST.PriceDesc' | translate) : ('CONSULTATION.CLIENT.BudgetDesc' | translate) }}</p>
                <div style="display: flex; gap: 1rem; max-width: 400px;">
                  <input type="number" formControlName="rage_price" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.AmountPlaceholder' | translate" style="flex: 2;"/>
                  <select *ngIf="role === 'therapist'" formControlName="currency" class="form-input" style="flex: 1;">
                    <option *ngFor="let cur of availableCurrencies" [value]="cur.id">{{ cur.label }}</option>
                  </select>
                </div>
              </div>
            </div>

            <!-- Step 2: Attributes (Languages & Specializations) -->
            <div class="step-pane" [class.active-pane]="currentStep === 2">
              <div class="section-group" style="--anim-delay: 1">
                <h4>{{ 'CONSULTATION.THERAPIST.Languages' | translate }}</h4>
                <p class="section-desc">{{ 'CONSULTATION.CreateProfile.LanguagesDesc' | translate }}</p>
                <div class="card-grid">
                  <div *ngFor="let lang of availableLanguages" 
                       class="selectable-card" 
                       [class.selected]="isLanguageSelected(lang.id)"
                       (click)="toggleLanguage(lang.id)">
                    <div class="selection-ring"></div>
                    <span class="card-title">{{ lang.label | translate }}</span>
                  </div>
                </div>
              </div>

              <div class="section-group" style="--anim-delay: 2">
                <h4>{{ role === 'therapist' ? ('CONSULTATION.THERAPIST.Specialization' | translate) : ('CONSULTATION.CLIENT.IssueDetail' | translate) }}</h4>
                <p class="section-desc">{{ 'CONSULTATION.CreateProfile.SpecializationDesc' | translate }}</p>
                <div class="card-grid columns-3">
                  <div *ngFor="let topic of availableTopics" 
                       class="selectable-card" 
                       [class.selected]="isTopicSelected(topic.id)"
                       (click)="toggleTopic(topic.id)">
                    <div class="selection-ring"></div>
                    <span class="card-title">{{ topic.label | translate }}</span>
                  </div>
                </div>
              </div>
            </div>

            <!-- Step 3: Details (Address & Experience) -->
            <div class="step-pane" [class.active-pane]="currentStep === 3">
              <div class="section-group" style="--anim-delay: 1">
                <h4>{{ role === 'therapist' ? ('CONSULTATION.THERAPIST.BasicDetails' | translate) : ('CONSULTATION.CreateProfile.LocationDetails' | translate) }}</h4>
                <p class="section-desc">{{ role === 'therapist' ? ('CONSULTATION.THERAPIST.BasicDetailsDesc' | translate) : ('CONSULTATION.CreateProfile.LocationDesc' | translate) }}</p>
                <div style="display: flex; flex-direction: column; gap: 1rem;">
                  <input *ngIf="role === 'therapist'" type="text" formControlName="name" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.Name' | translate" />
                  <input type="text" formControlName="address" class="form-input" [placeholder]="'CONSULTATION.CreateProfile.CityCountry' | translate" />
                </div>
              </div>

              <div class="section-group" style="--anim-delay: 2" *ngIf="role === 'therapist'">
                <h4>{{ 'CONSULTATION.THERAPIST.ProfessionalTitleAndExperience' | translate }}</h4>
                <div style="display: flex; gap: 1rem; margin-bottom: 1rem;">
                  <select formControlName="professional_title_code" class="form-input" style="flex: 2;">
                    <option value="">{{ 'CONSULTATION.THERAPIST.SelectTitle' | translate }}</option>
                    <option *ngFor="let title of availableTitles" [value]="title.code">{{ title.display }}</option>
                  </select>
                  <input type="number" formControlName="experience" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.ExperienceYears' | translate" style="flex: 1;" />
                </div>
                <input type="text" formControlName="professional_title_display" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.CustomTitle' | translate" />
              </div>

              <div class="summary-box" style="--anim-delay: 3" *ngIf="role === 'client'">
                <p [innerHTML]="'CONSULTATION.CreateProfile.AlmostDoneClient' | translate"></p>
              </div>
            </div>

            <!-- Step 4: Professional Info (Degrees & Certifications for Therapist) -->
            <div class="step-pane" [class.active-pane]="currentStep === 4" *ngIf="role === 'therapist'">
              <div class="section-group" style="--anim-delay: 1">
                <h4>{{ 'CONSULTATION.THERAPIST.Degrees' | translate }}</h4>
                <p class="section-desc">{{ 'CONSULTATION.THERAPIST.DegreesDesc' | translate }}</p>
                <div formArrayName="degrees" class="array-list">
                  <div *ngFor="let deg of degrees.controls; let i = index" [formGroupName]="i" class="array-item">
                    <button type="button" class="btn-remove" (click)="removeDegree(i)">&times;</button>
                    <select formControlName="type" class="form-input">
                      <option value="">{{ 'CONSULTATION.THERAPIST.DegreeType' | translate }}</option>
                      <option *ngFor="let type of availableDegreeTypes" [value]="type">{{ type }}</option>
                    </select>
                    <input type="text" formControlName="field" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.FieldOfStudy' | translate" />
                    
                    <div style="position: relative; flex: 1; min-width: 200px;">
                      <input type="text" formControlName="institution" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.Institution' | translate" (input)="onInstitutionSearch($event, i)" (focus)="activeDegreeIndex = i" (blur)="onInstitutionBlur()" autocomplete="off" style="width: 100%;" />
                      <div class="autocomplete-dropdown" *ngIf="activeDegreeIndex === i && universityResults.length > 0">
                        <div class="autocomplete-item" *ngFor="let uni of universityResults" (mousedown)="selectUniversity(uni.name, i)">
                          {{ uni.name }} <small *ngIf="uni.country">({{ uni.country }})</small>
                        </div>
                      </div>
                    </div>
                    
                    <input type="number" formControlName="year" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.Year' | translate" style="max-width: 100px;" />
                  </div>
                </div>
                <button type="button" class="btn-add" (click)="addDegree()">{{ 'CONSULTATION.THERAPIST.AddDegree' | translate }}</button>
              </div>

              <div class="section-group" style="--anim-delay: 2">
                <h4>{{ 'CONSULTATION.THERAPIST.Certifications' | translate }}</h4>
                <p class="section-desc">{{ 'CONSULTATION.THERAPIST.CertificationsDesc' | translate }}</p>
                <div formArrayName="certifications" class="array-list">
                  <div *ngFor="let cert of certifications.controls; let i = index" [formGroupName]="i" class="array-item">
                    <button type="button" class="btn-remove" (click)="removeCertification(i)">&times;</button>
                    <input type="text" formControlName="name" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.CertificationName' | translate" />
                    <input type="text" formControlName="issuer" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.Issuer' | translate" />
                    <input type="number" formControlName="year" class="form-input" [placeholder]="'CONSULTATION.THERAPIST.Year' | translate" style="max-width: 100px;" />
                  </div>
                </div>
                <button type="button" class="btn-add" (click)="addCertification()">{{ 'CONSULTATION.THERAPIST.AddCertification' | translate }}</button>
              </div>

              <div class="summary-box" style="--anim-delay: 3">
                <p [innerHTML]="'CONSULTATION.CreateProfile.AlmostDoneTherapist' | translate"></p>
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

    /* Dynamic Form Array Styles */
    .array-list {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      margin-top: 1rem;
    }
    
    .array-item {
      display: flex;
      flex-wrap: wrap;
      gap: 0.75rem;
      padding: 1rem;
      background: var(--color-bg);
      border: 1px dashed var(--color-border);
      border-radius: 6px;
      position: relative;
    }

    .array-item .form-input {
      flex: 1 1 200px;
    }

    .btn-remove {
      position: absolute;
      top: -10px;
      right: -10px;
      width: 24px;
      height: 24px;
      border-radius: 50%;
      background: var(--color-error, #ef4444);
      color: white;
      border: none;
      display: flex;
      align-items: center;
      justify-content: center;
      cursor: pointer;
      font-size: 14px;
      line-height: 1;
      padding: 0;
    }

    .btn-add {
      margin-top: 0.5rem;
      padding: 0.5rem 1rem;
      background: transparent;
      border: 1px solid var(--color-primary);
      color: var(--color-primary);
      border-radius: 6px;
      font-weight: 500;
      cursor: pointer;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
    }

    .btn-add:hover {
      background: rgba(var(--color-primary-rgb, 102,126,234), 0.08);
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

  // --- Search Subscriptions ---
  private institutionSearch$ = new Subject<{term: string, index: number}>();
  universityResults: any[] = [];
  activeDegreeIndex: number | null = null;

  // Data for cards
  availableModes = [
    { id: 'online', label: 'CONSULTATION.Modes.Online', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M22 16.92v3a2 2 0 0 1-2.18 2 19.79 19.79 0 0 1-8.63-3.07 19.5 19.5 0 0 1-6-6 19.79 19.79 0 0 1-3.07-8.67A2 2 0 0 1 4.11 2h3a2 2 0 0 1 2 1.72 12.84 12.84 0 0 0 .7 2.81 2 2 0 0 1-.45 2.11L8.09 9.91a16 16 0 0 0 6 6l1.27-1.27a2 2 0 0 1 2.11-.45 12.84 12.84 0 0 0 2.81.7A2 2 0 0 1 22 16.92z"></path></svg>' },
    { id: 'in-person', label: 'CONSULTATION.Modes.InPerson', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z"></path><circle cx="12" cy="10" r="3"></circle></svg>' },
    { id: 'chat', label: 'CONSULTATION.Modes.Chat', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M21 15a2 2 0 0 1-2 2H7l-4 4V5a2 2 0 0 1 2-2h14a2 2 0 0 1 2 2z"></path></svg>' },
    { id: 'phone', label: 'CONSULTATION.Modes.Phone', icon: '<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="5" y="2" width="14" height="20" rx="2" ry="2"></rect><line x1="12" y1="18" x2="12.01" y2="18"></line></svg>' }
  ];

  availableLanguages = [
    { id: 'en', label: 'CONSULTATION.Languages.English' },
    { id: 'vi', label: 'CONSULTATION.Languages.Vietnamese' },
    { id: 'zh', label: 'CONSULTATION.Languages.Chinese' },
    { id: 'ja', label: 'CONSULTATION.Languages.Japanese' }
  ];

  availableTopics = [
    { id: 'anxiety', label: 'CONSULTATION.Topics.Anxiety' },
    { id: 'depression', label: 'CONSULTATION.Topics.Depression' },
    { id: 'stress', label: 'CONSULTATION.Topics.Stress' },
    { id: 'relationship', label: 'CONSULTATION.Topics.Relationships' },
    { id: 'trauma', label: 'CONSULTATION.Topics.Trauma' },
    { id: 'career', label: 'CONSULTATION.Topics.Career' },
  ];

  availableTitles = [
    { code: 'phd', display: 'Doctor of Philosophy (Ph.D)' },
    { code: 'psyd', display: 'Doctor of Psychology (Psy.D)' },
    { code: 'ms', display: 'Master of Science (MS)' },
    { code: 'ma', display: 'Master of Arts (MA)' },
    { code: 'lcsw', display: 'Licensed Clinical Social Worker (LCSW)' },
    { code: 'lpc', display: 'Licensed Professional Counselor (LPC)' }
  ];

  availableDegreeTypes = ['Bachelors', 'Masters', 'Doctorate', 'MD', 'Other'];
  availableCurrencies = [{ id: 'VND', label: 'VND' }, { id: 'USD', label: 'USD' }];

  constructor(
    private router: Router,
    private route: ActivatedRoute,
    private authService: AuthService,
    private consultationProfileService: ConsultationProfileService,
    private toastService: ToastService,
    private fb: FormBuilder,
    private http: HttpClient
  ) {
    this.profileForm = this.fb.group({
      address: ['', Validators.required],
      rage_price: [null, [Validators.required, Validators.min(0)]],
      languages: [[]],
      consultation_modes: [[]],
      issue_detail: [[]], // Client
      specialization: [[]], // Therapist
      experience: [null], // Therapist
      name: [''], // Therapist
      currency: ['USD'], // Therapist
      professional_title_code: [''], // Therapist
      professional_title_display: [''], // Therapist
      degrees: this.fb.array([]), // Therapist
      certifications: this.fb.array([]) // Therapist
    });
  }

  isEditing = false;
  
  ngOnInit() {
    this.route.queryParams.subscribe(params => {
      const rawRole = params['role'] || this.authService.getRole();
      if (!rawRole) {
        this.toastService.error('Error', 'No role found. Please login again.');
        this.router.navigate(['/login']);
      } else {
        this.role = (rawRole as string).toLowerCase() as 'therapist' | 'client';
        // Add one initial empty degree/cert if therapist
        if (this.role === 'therapist') {
          // Initialize empty state is fine, user can click Add
        }
        this.loadExistingProfile();
      }
    });

    // Subscribe to institution searching
    this.institutionSearch$.pipe(
      debounceTime(300),
      distinctUntilChanged((prev, curr) => prev.term === curr.term && prev.index === curr.index),
      switchMap(request => {
        if (!request.term || request.term.length < 2) {
          return of([]);
        }
        return this.http.get<any[]>(`http://universities.hipolabs.com/search?name=${encodeURIComponent(request.term)}`).pipe(
          catchError(() => of([]))
        );
      })
    ).subscribe(results => {
      this.universityResults = results.slice(0, 15);
    });
  }

  get degrees() {
    return this.profileForm.get('degrees') as FormArray;
  }

  get certifications() {
    return this.profileForm.get('certifications') as FormArray;
  }

  addDegree() {
    this.degrees.push(this.fb.group({
      type: [''],
      field: [''],
      institution: [''],
      year: [new Date().getFullYear()]
    }));
  }

  removeDegree(index: number) {
    this.degrees.removeAt(index);
  }

  addCertification() {
    this.certifications.push(this.fb.group({
      name: [''],
      issuer: [''],
      year: [new Date().getFullYear()]
    }));
  }

  removeCertification(index: number) {
    this.certifications.removeAt(index);
  }

  loadExistingProfile() {
    this.consultationProfileService.getCurrentProfile().subscribe({
      next: (res: any) => {
        if (res && res.data) {
          const profile = res.data;
          this.isEditing = true;
          this.profileForm.patchValue({
            address: profile.address || '',
            rage_price: profile.rage_price || null,
            languages: profile.languages || [],
            consultation_modes: profile.consultation_modes || [],
            experience: profile.experience || null,
            name: profile.name || '',
            currency: profile.currency || 'USD'
          });
          
          if (this.role === 'therapist') {
            if (profile.specialization) {
              this.profileForm.patchValue({ specialization: profile.specialization });
            }
            if (profile.professional_info) {
              if (profile.professional_info.title) {
                 this.profileForm.patchValue({
                   professional_title_code: profile.professional_info.title.code || '',
                   professional_title_display: profile.professional_info.title.display || '',
                 });
              }
              if (profile.professional_info.degrees && Array.isArray(profile.professional_info.degrees)) {
                 profile.professional_info.degrees.forEach((deg: any) => {
                    this.degrees.push(this.fb.group({
                      type: [deg.type || ''],
                      field: [deg.field || ''],
                      institution: [deg.institution || ''],
                      year: [deg.year || new Date().getFullYear()]
                    }));
                 });
              }
              if (profile.professional_info.certifications && Array.isArray(profile.professional_info.certifications)) {
                 profile.professional_info.certifications.forEach((cert: any) => {
                    this.certifications.push(this.fb.group({
                      name: [cert.name || ''],
                      issuer: [cert.issuer || ''],
                      year: [cert.year || new Date().getFullYear()]
                    }));
                 });
              }
            }
          } else if (this.role === 'client' && profile.issue_detail) {
             this.profileForm.patchValue({ issue_detail: profile.issue_detail });
          }
        }
      },
      error: (err) => {
         console.warn('Could not load existing profile, assuming creation flow.', err);
      }
    });
  }

  // --- Step Navigation ---
  getStepTitle(): string {
    switch (this.currentStep) {
      case 1: return 'CONSULTATION.CreateProfile.StepConfig';
      case 2: return 'CONSULTATION.CreateProfile.StepAttributes';
      case 3: return 'CONSULTATION.CreateProfile.StepDetails';
      case 4: return 'CONSULTATION.THERAPIST.ProfessionalInfo';
      default: return '';
    }
  }

  nextStep() {
    // Optional: add step validation here
    const maxSteps = this.role === 'therapist' ? 4 : 3;
    if (this.currentStep < maxSteps) {
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
      profileData.name = formVal.name;
      profileData.currency = formVal.currency;
      profileData.experience = Number(formVal.experience) || 0;
      profileData.specialization = formVal.specialization;
      profileData.professional_info = {
        title: {
          code: formVal.professional_title_code,
          display: formVal.professional_title_display || formVal.professional_title_code // Fallback name
        },
        experience_years: Number(formVal.experience) || 0,
        degrees: formVal.degrees || [],
        certifications: formVal.certifications || []
      };
    } else {
      profileData.issue_detail = formVal.issue_detail;
    }

    let updateMethod;
    if (this.isEditing) {
       updateMethod = this.role === 'therapist'
         ? this.consultationProfileService.updateTherapistProfile(profileData)
         : this.consultationProfileService.updateClientProfile(profileData);
    } else {
       updateMethod = this.role === 'therapist'
         ? this.consultationProfileService.createTherapistProfile(profileData)
         : this.consultationProfileService.createClientProfile(profileData);
    }

    updateMethod.subscribe({
      next: (profile: any) => {
        this.toastService.success('Success', this.isEditing ? 'Profile updated successfully!' : 'Profile created successfully!');
        this.router.navigate(['/']);
      },
      error: (error: any) => {
        console.error('Failed to save profile:', error);
        this.toastService.error('Error', 'Failed to save profile. Please try again.');
        this.creating = false;
      }
    });
  }

  // --- University Autocomplete Methods ---
  onInstitutionSearch(event: any, index: number) {
    this.activeDegreeIndex = index;
    this.institutionSearch$.next({ term: event.target.value, index });
  }

  selectUniversity(uniName: string, index: number) {
    const degrees = this.profileForm.get('degrees') as FormArray;
    const group = degrees.at(index) as FormGroup;
    group.patchValue({ institution: uniName });
    this.universityResults = [];
    this.activeDegreeIndex = null;
  }

  onInstitutionBlur() {
    // Delay hiding to allow mousedown to fire on the list item first
    setTimeout(() => {
      this.activeDegreeIndex = null;
      this.universityResults = [];
    }, 200);
  }
}
