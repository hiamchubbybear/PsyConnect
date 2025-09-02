// therapist-profile-overlay.component.ts
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

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
  imports: [CommonModule, TranslateModule],
  templateUrl: 'profile-overlay.html',
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
    this.closeOverlay.emit();
  }

  onBackdropClick(event: Event): void {
    if (event.target === event.currentTarget) {
      this.onClose();
    }
  }

  onBookAppointment() {
    if (this.therapist) this.bookAppointment.emit(this.therapist);
    this.onClose();
  }

  onSendMessage() {
    if (this.therapist) this.sendMessage.emit(this.therapist);
    this.onClose();
  }

  hasQualifications(): boolean {
    return !!(
      this.therapist?.professional_info?.degrees?.length ||
      this.therapist?.professional_info?.certifications?.length
    );
  }
}
