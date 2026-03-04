import { CommonModule } from '@angular/common';
import { Component, OnDestroy, OnInit } from '@angular/core';
import { Router, RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { forkJoin, of, Subscription } from 'rxjs';
import { catchError } from 'rxjs/operators';
import {
  CONSULTATION_MODES,
  ConsultationSession,
  LANGUAGES,
  SPECIALIZATIONS,
  TIME_SLOTS,
  WEEKDAYS,
} from '../../../models/consultation.model';
import { AuthService } from '../../../services/auth/auth.service';
import { ConsultationProfileService } from '../../../services/consultation/consultation-profile.service';
import { Profile } from '../../../services/profile/profile';

interface DisplayProfile {
  name: string;
  avatarUri: string;
  address: string;
  languages: string[];
  consultationModes: string[];
  specializations: string[];
  availabilityDays: string[];
  timeSlots: string[];
  urgencyLevel: string;
  sessionDuration: number;
  genderPreference: string;
  experienceLevel: string;
  isFlexible: boolean;
  sessionPrice: number;
}

@Component({
  selector: 'app-my-consultation-profile',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule],
  templateUrl: './my-profile.html',
  styleUrls: ['./my-profile.scss'],
})
export class MyConsultationProfileComponent implements OnInit, OnDestroy {
  profile: DisplayProfile | null = null;
  sessions: ConsultationSession[] = [];
  loading = true;
  error = false;
  activeTab: 'info' | 'sessions' | 'reviews' = 'info';
  hasProfile = false;
  role: 'therapist' | 'client' | null = null;

  private langSub?: Subscription;

  // Lookup maps for display labels
  readonly MODES_MAP = Object.fromEntries(
    CONSULTATION_MODES.map((m) => [m.value, m.labelVi]),
  );
  readonly LANG_MAP = Object.fromEntries(
    LANGUAGES.map((l) => [l.value, l.labelVi]),
  );
  readonly SPEC_MAP = Object.fromEntries(
    SPECIALIZATIONS.map((s) => [s.value, s.labelVi]),
  );
  readonly DAY_MAP = Object.fromEntries(
    WEEKDAYS.map((d) => [d.value, d.labelVi]),
  );
  readonly SLOT_MAP = Object.fromEntries(
    TIME_SLOTS.map((t) => [t.value, t.labelVi]),
  );

  constructor(
    private consultProfileService: ConsultationProfileService,
    private profileService: Profile,
    private router: Router,
    private authService: AuthService,
    private translate: TranslateService,
  ) {}

  ngOnInit(): void {
    this.loadAll();
    // Re-check labels on language change
    this.langSub = this.translate.onLangChange.subscribe(() => {
      // Logic for re-calculating labels if they were pre-calculated
    });
  }

  ngOnDestroy(): void {
    this.langSub?.unsubscribe();
  }

  loadAll(): void {
    this.loading = true;
    this.error = false;

    forkJoin({
      userProfile: this.profileService
        .getProfile()
        .pipe(catchError(() => of(null))),
      consultProfile: this.consultProfileService
        .getCurrentProfile()
        .pipe(catchError(() => of(null))),
    }).subscribe({
      next: ({ userProfile, consultProfile }) => {
        const userRes = (userProfile as any)?.data ?? userProfile;
        const consultRes = (consultProfile as any)?.data ?? consultProfile;

        if (!consultRes) {
          this.hasProfile = false;
          this.loading = false;
          return;
        }

        this.hasProfile = true;
        this.role = this.authService.getRole();

        this.profile = {
          name:
            `${userRes?.firstName ?? ''} ${userRes?.lastName ?? ''}`.trim() ||
            this.translate.instant('ROLE.USER'),
          avatarUri: userRes?.avatarUri ?? '',
          address: consultRes.address ?? '',
          languages: consultRes.languages ?? consultRes.languages_spoken ?? [],
          consultationModes:
            consultRes.consultation_modes ??
            consultRes.consultation_options ??
            [],
          specializations:
            consultRes.therapist_specialization ??
            consultRes.specializations ??
            [],
          availabilityDays: consultRes.availability?.days ?? [],
          timeSlots: consultRes.availability?.time_slots ?? [],
          urgencyLevel: consultRes.urgency_level ?? 'flexible',
          sessionDuration: consultRes.session_duration ?? 60,
          genderPreference: consultRes.preferred_therapist_gender ?? 'any',
          experienceLevel:
            consultRes.experience_level ??
            (consultRes.experience_years
              ? `${consultRes.experience_years}`
              : 'any'),
          isFlexible: consultRes.is_flexible_with_schedule ?? false,
          sessionPrice: consultRes.rage_price ?? consultRes.price_per_hour ?? 0,
        };
        this.loading = false;
      },
      error: () => {
        this.error = true;
        this.loading = false;
      },
    });
  }

  setTab(tab: 'info' | 'sessions' | 'reviews'): void {
    this.activeTab = tab;
  }

  editProfile(): void {
    this.router.navigate(['/feature/consultation/settings']);
  }

  createProfile(): void {
    this.router.navigate(['/feature/consultation/create-profile']);
  }

  label(map: Record<string, string>, val: string): string {
    return map[val] ?? val;
  }

  urgencyLabel(): string {
    if (!this.profile) return '';
    const level =
      this.profile.urgencyLevel.charAt(0).toUpperCase() +
      this.profile.urgencyLevel.slice(1);
    const key = `CONSULTATION.MyProfile.Labels.Urgency.${level}`;
    const trans = this.translate.instant(key);
    return trans !== key ? trans : this.profile.urgencyLevel;
  }

  genderLabel(): string {
    if (!this.profile) return '';
    const gender = this.profile.genderPreference;
    let keyPart = 'Any';
    if (gender === 'male') keyPart = 'Male';
    if (gender === 'female') keyPart = 'Female';
    if (gender === 'non-binary' || gender === 'non_binary')
      keyPart = 'NonBinary';

    const key = `CONSULTATION.MyProfile.Labels.Gender.${keyPart}`;
    const trans = this.translate.instant(key);
    return trans !== key ? trans : gender;
  }

  experienceLabel(): string {
    if (!this.profile) return '';
    const level = this.profile.experienceLevel;

    // Handle numeric years
    if (/^\d+$/.test(level)) {
      return `${level} ${this.translate.instant('SOCIAL.sidebar.items.wall') === 'Bảng tin' ? 'năm' : 'years'}`; // Simple fallback
    }

    const levelKey = level.charAt(0).toUpperCase() + level.slice(1);
    const key = `CONSULTATION.MyProfile.Labels.Experience.${levelKey}`;
    const trans = this.translate.instant(key);
    return trans !== key ? trans : level;
  }

  get avatarInitials(): string {
    if (!this.profile?.name) return 'U';
    const parts = this.profile.name.trim().split(' ');
    if (parts.length === 1) return parts[0]?.[0]?.toUpperCase() ?? 'U';
    return (
      (parts[0]?.[0] ?? '') + (parts[parts.length - 1]?.[0] ?? '')
    ).toUpperCase();
  }
}
