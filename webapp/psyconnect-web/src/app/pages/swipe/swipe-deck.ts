import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import {
    TherapistProfile,
    TherapistProfileOverlayComponent,
} from '../../components/profile-overlay/profile-overlay';
import { SwipeCardComponent } from '../../components/swipe-card/swipe-card';
import { SecureStorageService } from '../../encrypt/secure';
import { SwipeService } from '../../services/swipe/swipe.service';

@Component({
  selector: 'app-swipe-deck',
  templateUrl: './swipe-deck.html',
  standalone: true,
  styleUrls: ['./swipe-deck.scss'],
  imports: [SwipeCardComponent, CommonModule, TherapistProfileOverlayComponent],
})
export class SwipeDeckComponent implements OnInit {
  therapists: TherapistProfile[] = [];
  currentIndex = 0;
  isOverlayOpen = true;
  selectedTherapist: TherapistProfile | null = null;

  constructor(
    private swipeService: SwipeService,
    private secureStorage: SecureStorageService
  ) {}

  ngOnInit() {
    if (!this.therapists.length) {
      this.therapists = [
        {
          profile_id: 't001',
          name: 'Dr. Sarah Johnson',
          address: '123 Main Street, City, Country',
          languages: ['English', 'Spanish'],
          specialization: [
            'Cognitive Behavioral Therapy',
            'Anxiety Management',
          ],
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
        },
        {
          profile_id: 't002',
          name: 'Dr. Michael Chen',
          address: '456 Oak Avenue, Downtown, Country',
          languages: ['English', 'Mandarin'],
          specialization: ['Family Therapy', 'Couples Counseling'],
          consultation_modes: ['Online', 'In-Person'],
          experience: 8,
          rating: 4.9,
          availability: {
            days: ['Tuesday', 'Thursday', 'Saturday'],
            time_slots: ['9:00 AM - 11:00 AM', '1:00 PM - 3:00 PM'],
          },
          professional_info: {
            title: {
              code: 'LMFT',
              display: 'Licensed Marriage and Family Therapist',
            },
            degrees: [
              'M.S. in Marriage and Family Therapy',
              'B.A. in Psychology',
            ],
            certifications: ['Gottman Method Couples Therapy', 'EFT Certified'],
            experience_years: 8,
          },
        },
      ];
    }
  }

  openOverlay(therapist: TherapistProfile) {
    console.log('Opening overlay for:', therapist.name);
    this.selectedTherapist = { ...therapist };
    this.isOverlayOpen = true;
    document.body.classList.add('overlay-open');
  }

  onCloseOverlay() {
    this.isOverlayOpen = false;
    this.selectedTherapist = null;
    document.body.classList.remove('overlay-open');
  }

  onCardSwiped(event: {
    direction: 'left' | 'right' | 'up';
    therapist: TherapistProfile;
  }) {
    console.log('Card swiped:', event.direction, event.therapist.name);

    if (event.direction === 'up') {
      this.openOverlay(event.therapist);
      return;
    }

    this.therapists = this.therapists.filter(
      (t) => t.profile_id !== event.therapist.profile_id
    );

    localStorage.setItem('therapists', JSON.stringify(this.therapists));

    if (event.direction === 'right') {
      console.log('Matched:', event.therapist.name);
    } else if (event.direction === 'left') {
      console.log('Passed:', event.therapist.name);
    }
  }

  onBookAppointment(therapist: TherapistProfile) {
    console.log('Book appointment with:', therapist.name);
    // Add booking logic here
    this.onCloseOverlay();
  }

  onSendMessage(therapist: TherapistProfile) {
    console.log('Send message to:', therapist.name);
    // Add messaging logic here
    this.onCloseOverlay();
  }

  get currentTherapist(): TherapistProfile | null {
    return this.therapists.length > 0 ? this.therapists[0] : null;
  }

  loadCurrentTherapist() {
    this.therapists =
      this.secureStorage.getItem<TherapistProfile[]>('therapists') ?? [];
  }

  formatNextAvailable(availability: any): string {
    if (!availability) return '';
    const days = Array.isArray(availability.days)
      ? availability.days.join(', ')
      : '';
    const times = Array.isArray(availability.time_slots)
      ? availability.time_slots.join(', ')
      : '';
    return `${days}${days && times ? ' - ' : ''}${times}`;
  }
}
