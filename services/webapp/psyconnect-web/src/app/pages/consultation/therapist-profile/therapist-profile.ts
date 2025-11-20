import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormArray, FormBuilder, FormGroup, ReactiveFormsModule, Validators } from '@angular/forms';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import {
    CONSULTATION_MODES,
    LANGUAGES,
    PROFESSIONAL_TITLES,
    SPECIALIZATIONS,
    TherapistV1,
    TIME_SLOTS,
    WEEKDAYS,
} from '../../../models/consultation.model';
import { TherapistService } from '../../../services/consultation/therapist.service';
import { ToastService } from '../../../shared/toast/toast.service';

@Component({
  selector: 'app-therapist-profile',
  standalone: true,
  imports: [CommonModule, ReactiveFormsModule, TranslateModule],
  templateUrl: './therapist-profile.html',
  styleUrls: ['./therapist-profile.scss'],
})
export class TherapistProfileComponent implements OnInit {
  profileForm!: FormGroup;
  loading = false;
  existingProfile: TherapistV1 | null = null;
  isEditMode = false;

  // Reference data
  weekdays = WEEKDAYS;
  timeSlots = TIME_SLOTS;
  consultationModes = CONSULTATION_MODES;
  languages = LANGUAGES;
  specializations = SPECIALIZATIONS;
  professionalTitles = PROFESSIONAL_TITLES;

  constructor(
    private fb: FormBuilder,
    private therapistService: TherapistService,
    private toastService: ToastService,
    private router: Router
  ) {}

  ngOnInit() {
    this.initForm();
    this.loadExistingProfile();
  }

  initForm() {
    this.profileForm = this.fb.group({
      name: ['', [Validators.required, Validators.minLength(2)]],
      address: ['', Validators.required],
      languages: [[], Validators.required],
      specialization: [[], Validators.required],
      consultation_modes: [[], Validators.required],
      experience: [0, [Validators.required, Validators.min(0)]],
      currency: ['USD', Validators.required],
      rage_price: [0, [Validators.required, Validators.min(0)]],
      availability: this.fb.group({
        days: [[], Validators.required],
        time_slots: [[], Validators.required],
      }),
      professional_info: this.fb.group({
        title: this.fb.group({
          code: ['', Validators.required],
          display: ['', Validators.required],
        }),
        degrees: this.fb.array([]),
        certifications: this.fb.array([]),
        experience_years: [0, [Validators.required, Validators.min(0)]],
      }),
    });
  }

  loadExistingProfile() {
    this.loading = true;
    this.therapistService.getTherapistProfile().subscribe({
      next: (profile) => {
        if (profile && profile.profile_id) {
          this.existingProfile = profile;
          this.isEditMode = true;
          this.patchFormValues(profile);
        }
        this.loading = false;
      },
      error: (err) => {
        console.error('Error loading profile:', err);
        this.loading = false;
      },
    });
  }

  patchFormValues(profile: TherapistV1) {
    this.profileForm.patchValue({
      name: profile.name,
      address: profile.address,
      languages: profile.languages,
      specialization: profile.specialization,
      consultation_modes: profile.consultation_modes,
      experience: profile.experience,
      currency: profile.currency,
      rage_price: profile.rage_price,
      availability: {
        days: profile.availability?.days || [],
        time_slots: profile.availability?.time_slots || [],
      },
    });

    if (profile.professional_info) {
      this.profileForm.patchValue({
        professional_info: {
          title: profile.professional_info.title,
          experience_years: profile.professional_info.experience_years,
        },
      });

      // Add degrees
      const degreesArray = this.degrees;
      profile.professional_info.degrees?.forEach((degree) => {
        degreesArray.push(this.createDegreeGroup(degree));
      });

      // Add certifications
      const certsArray = this.certifications;
      profile.professional_info.certifications?.forEach((cert) => {
        certsArray.push(this.createCertificationGroup(cert));
      });
    }
  }

  get degrees(): FormArray {
    return this.profileForm.get('professional_info.degrees') as FormArray;
  }

  get certifications(): FormArray {
    return this.profileForm.get('professional_info.certifications') as FormArray;
  }

  createDegreeGroup(degree?: any): FormGroup {
    return this.fb.group({
      type: [degree?.type || '', Validators.required],
      field: [degree?.field || '', Validators.required],
      institution: [degree?.institution || '', Validators.required],
      year: [degree?.year || new Date().getFullYear(), Validators.required],
    });
  }

  createCertificationGroup(cert?: any): FormGroup {
    return this.fb.group({
      name: [cert?.name || '', Validators.required],
      issuer: [cert?.issuer || '', Validators.required],
      year: [cert?.year || new Date().getFullYear(), Validators.required],
    });
  }

  addDegree() {
    this.degrees.push(this.createDegreeGroup());
  }

  removeDegree(index: number) {
    this.degrees.removeAt(index);
  }

  addCertification() {
    this.certifications.push(this.createCertificationGroup());
  }

  removeCertification(index: number) {
    this.certifications.removeAt(index);
  }

  onTitleChange(event: Event) {
    const select = event.target as HTMLSelectElement;
    const selectedTitle = this.professionalTitles.find((t) => t.code === select.value);
    if (selectedTitle) {
      this.profileForm.patchValue({
        professional_info: {
          title: {
            code: selectedTitle.code,
            display: selectedTitle.display,
          },
        },
      });
    }
  }

  onSubmit() {
    if (this.profileForm.invalid) {
      Object.keys(this.profileForm.controls).forEach((key) => {
        this.profileForm.get(key)?.markAsTouched();
      });
      return;
    }

    this.loading = true;
    const formValue = this.profileForm.value;

    const therapistData: TherapistV1 = {
      ...formValue,
      profile_id: this.existingProfile?.profile_id,
    };

    const request = this.isEditMode
      ? this.therapistService.updateTherapistProfile(therapistData)
      : this.therapistService.createTherapistProfile(therapistData);

    request.subscribe({
      next: () => {
        this.toastService.success(
          'Success',
          this.isEditMode ? 'Profile updated successfully!' : 'Profile created successfully!'
        );
        this.router.navigate(['/feature/schedule']);
      },
      error: (err) => {
        console.error('Failed to save profile:', err);
        this.toastService.error('Error', 'Failed to save profile. Please try again.');
        this.loading = false;
      },
    });
  }

  cancel() {
    this.router.navigate(['/feature/schedule']);
  }

  getErrorMessage(fieldName: string): string {
    const control = this.profileForm.get(fieldName);
    if (!control || !control.errors || !control.touched) return '';

    if (control.errors['required']) return `${fieldName} is required`;
    if (control.errors['minlength']) {
      return `${fieldName} must be at least ${control.errors['minlength'].requiredLength} characters`;
    }
    if (control.errors['min']) {
      return `${fieldName} must be at least ${control.errors['min'].min}`;
    }
    return '';
  }

  // Helper methods for checkbox arrays
  toggleLanguage(language: string, checked: boolean) {
    const current = this.profileForm.value.languages || [];
    const updated = checked
      ? [...current, language]
      : current.filter((l: string) => l !== language);
    this.profileForm.patchValue({ languages: updated });
  }

  toggleSpecialization(spec: string, checked: boolean) {
    const current = this.profileForm.value.specialization || [];
    const updated = checked
      ? [...current, spec]
      : current.filter((s: string) => s !== spec);
    this.profileForm.patchValue({ specialization: updated });
  }

  toggleConsultationMode(mode: string, checked: boolean) {
    const current = this.profileForm.value.consultation_modes || [];
    const updated = checked
      ? [...current, mode]
      : current.filter((m: string) => m !== mode);
    this.profileForm.patchValue({ consultation_modes: updated });
  }

  toggleDay(day: string, checked: boolean) {
    const current = this.profileForm.value.availability?.days || [];
    const updated = checked
      ? [...current, day]
      : current.filter((d: string) => d !== day);
    this.profileForm.patchValue({
      availability: {
        ...this.profileForm.value.availability,
        days: updated,
      },
    });
  }

  toggleTimeSlot(slot: string, checked: boolean) {
    const current = this.profileForm.value.availability?.time_slots || [];
    const updated = checked
      ? [...current, slot]
      : current.filter((t: string) => t !== slot);
    this.profileForm.patchValue({
      availability: {
        ...this.profileForm.value.availability,
        time_slots: updated,
      },
    });
  }
}
