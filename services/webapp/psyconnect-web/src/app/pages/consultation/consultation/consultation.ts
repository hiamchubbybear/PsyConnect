
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import {
    FormBuilder,
    FormGroup,
    ReactiveFormsModule,
    Validators,
} from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import { finalize } from 'rxjs';
import {
    Client,
    CONSULTATION_MODES,
    EXPERIENCE_LEVELS,
    GENDER_PREFERENCES,
    LANGUAGES,
    SPECIALIZATIONS,
    TIME_SLOTS,
    URGENCY_LEVELS,
    WEEKDAYS,
} from '../../../models/consultation.model';
import { ConsultationClientService } from '../../../services/consultation/consultation.service';

@Component({
  selector: 'app-consultation-profile',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './consultation.html',
  styleUrls: ['./consultation.scss'],
})
export class ProfileFormClientComponent implements OnInit {
  @Input() initialData?: Partial<Client>;
  @Input() isSaving = false;
  @Output() save = new EventEmitter<Client>();
  @Output() cancel = new EventEmitter<void>();

  form!: FormGroup;
  currentStep = 1;

  weekdays = WEEKDAYS;
  timeSlots = TIME_SLOTS;
  consultationModes = CONSULTATION_MODES;
  languages = LANGUAGES;
  specializations = SPECIALIZATIONS;
  urgencyLevels = URGENCY_LEVELS;
  experienceLevels = EXPERIENCE_LEVELS;
  genderPreferences = GENDER_PREFERENCES;
  loading = false;

  constructor(
    private fb: FormBuilder,
    private consultationService: ConsultationClientService
  ) {}

  ngOnInit(): void {
    this.initForm();
    this.loadConsultationProfile();
  }
  loadConsultationProfile(): void {
    this.loading = true;
    this.consultationService.getClientConsultationProfile().subscribe({
      next: (res: any) => {
        const data = res.data;
        if (!data) return;

        this.form.patchValue({
          address: data.address ?? '',
          languages: data.languages ?? [],
          consultation_modes: data.consultation_modes ?? [],
          availability: {
            days: data.availability?.days ?? [],
            time_slots: data.availability?.time_slots ?? [],
          },
          therapist_specialization: data.specialization ?? [],
        });
      },
      error: (err) => console.error('Error loading consultation profile:', err),
      complete: () => (this.loading = false),
    });
  }
  initForm(): void {
    this.form = this.fb.group({
      address: ['', Validators.required],
      languages: [[], Validators.required],
      issue_detail: [[]],
      consultation_modes: [[], Validators.required],
      rage_price: [0, [Validators.required, Validators.min(0)]],
      availability: this.fb.group({
        days: [[], Validators.required],
        time_slots: [[], Validators.required],
      }),
      preferred_therapist_gender: ['any'],
      experience_level: ['any'],
      therapist_specialization: [[]],
      urgency_level: ['flexible'],
      session_duration: [60],
      preferred_therapist_language: [[]],
      is_flexible_with_schedule: [false],
    });
  }
  nextStep(): void {
    if (this.validateCurrentStep()) {
      if (this.currentStep < 4) {
        this.currentStep++;
        this.scrollToTop();
      }
    } else {
      this.markCurrentStepTouched();
    }
  }

  previousStep(): void {
    if (this.currentStep > 1) {
      this.currentStep--;
      this.scrollToTop();
    }
  }

  validateCurrentStep(): boolean {
    switch (this.currentStep) {
      case 1:
        return (
          (this.form.get('address')?.valid &&
            this.form.get('languages')?.valid) ||
          false
        );
      case 2:
        return (
          (this.form.get('consultation_modes')?.valid &&
            this.form.get('rage_price')?.valid) ||
          false
        );
      case 3:
        return (
          (this.form.get('availability.days')?.valid &&
            this.form.get('availability.time_slots')?.valid) ||
          false
        );
      case 4:
        return true;
      default:
        return false;
    }
  }

  markCurrentStepTouched(): void {
    switch (this.currentStep) {
      case 1:
        this.form.get('address')?.markAsTouched();
        this.form.get('languages')?.markAsTouched();
        break;
      case 2:
        this.form.get('consultation_modes')?.markAsTouched();
        this.form.get('rage_price')?.markAsTouched();
        break;
      case 3:
        this.form.get('availability.days')?.markAsTouched();
        this.form.get('availability.time_slots')?.markAsTouched();
        break;
    }
  }

  scrollToTop(): void {
    window.scrollTo({ top: 0, behavior: 'smooth' });
  }

  toggleArrayValue(controlName: string, value: string): void {
    const control = this.form.get(controlName);
    if (!control) return;

    const currentValues: string[] = control.value || [];
    const index = currentValues.indexOf(value);

    if (index > -1) {
      currentValues.splice(index, 1);
    } else {
      currentValues.push(value);
    }

    control.setValue([...currentValues]);
    control.markAsTouched();
  }

  toggleNestedArrayValue(
    parentName: string,
    childName: string,
    value: string
  ): void {
    const control = this.form.get(`${parentName}.${childName}`);
    if (!control) return;

    const currentValues: string[] = control.value || [];
    const index = currentValues.indexOf(value);

    if (index > -1) {
      currentValues.splice(index, 1);
    } else {
      currentValues.push(value);
    }

    control.setValue([...currentValues]);
    control.markAsTouched();
  }

  isSelected(controlName: string, value: string): boolean {
    const control = this.form.get(controlName);
    return control?.value?.includes(value) || false;
  }

  isNestedSelected(
    parentName: string,
    childName: string,
    value: string
  ): boolean {
    const control = this.form.get(`${parentName}.${childName}`);
    return control?.value?.includes(value) || false;
  }

  onSubmit(): void {
    if (this.form.invalid) {
      this.markFormGroupTouched(this.form);
      return;
    }
    this.isSaving = true;
    const payload = this.form.value as Client;
    this.consultationService
      .saveClientProfile(payload)
      .pipe(finalize(() => (this.isSaving = false)))
      .subscribe({
        next: (response) => {
          console.log('Profile saved successfully', response);
          this.save.emit(response);
        },
        error: (err) => {
          console.error('Error saving profile:', err);
        },
      });
  }

  onCancel(): void {
    this.cancel.emit();
  }

  private markFormGroupTouched(formGroup: FormGroup): void {
    Object.keys(formGroup.controls).forEach((key) => {
      const control = formGroup.get(key);
      control?.markAsTouched();
      if (control instanceof FormGroup) {
        this.markFormGroupTouched(control);
      }
    });
  }

  getFieldError(fieldName: string): string | null {
    const control = this.form.get(fieldName);
    if (control?.touched && control?.errors) {
      if (control.errors['required']) return 'This field is required';
      if (control.errors['min']) return 'Value must be greater than 0';
    }
    return null;
  }
}
