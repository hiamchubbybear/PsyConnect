import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { switchMap } from 'rxjs';
import { TherapistProfileOverlayComponent } from '../../components/profile-overlay/profile-overlay';
import { SwipeCardComponent } from '../../components/swipe-card/swipe-card';
import { SecureStorageService } from '../../encrypt/secure';
import { mapTherapistResponse, Therapist } from '../../models/swipe-card';
import { SwipeService } from '../../services/swipe/swipe.service';
import { ToastType } from '../../shared/toast/toast.model';
import { ToastService } from '../../shared/toast/toast.service';

@Component({
  selector: 'app-swipe-deck',
  templateUrl: './swipe-deck.html',
  standalone: true,
  styleUrls: ['./swipe-deck.scss'],
  imports: [SwipeCardComponent, CommonModule, TherapistProfileOverlayComponent],
})
export class SwipeDeckComponent implements OnInit {
  therapists: Therapist[] = [];
  currentIndex = 0;
  isOverlayOpen = false;
  selectedTherapist: Therapist | null = null;

  constructor(
    private swipeService: SwipeService,
    private secureStorage: SecureStorageService,
    private toastService : ToastService
  ) {}

  ngOnInit() {
    this.loadTherapists();
  }

  private loadTherapists() {
    this.swipeService
      .triggerUpdate()
      .pipe(switchMap(() => this.swipeService.getSwipeData()))
      .subscribe({
        next: (res) => {
          if (res?.data?.length) {
            this.therapists = mapTherapistResponse(res);
          } else {
           this.onHandleUpdateTherapist();
          }
        },
        error: (err) => {
          this.therapists = this.getFallbackData();
        },
      });
  }
  onHandleUpdateTherapist() {
    this.swipeService.getUpdateTherapistHandler().subscribe({

      next: (res) => {
        if(res) this.toastService.show("Updated" , "Therapist updated" , ToastType.Success)
          else this.toastService.show("Updated" , "Therapist updated failed"  , ToastType.Success)
      },
      error : err => this.toastService.show("Updated" , err  , ToastType.Success)
    }
    )
  }

  private getFallbackData(): Therapist[] {
    return [
      {
        profileId: 't001',
        name: 'Dr. Sarah Johnson',
        address: '123 Main Street, City, Country',
        languages: ['English', 'Spanish'],
        specialization: ['Cognitive Behavioral Therapy', 'Anxiety Management'],
        consultationModes: ['Online', 'In-Person'],
        experience: 10,
        rating: 4.8,
        currency: 'USD',
        ragePrice: 120,
        isAvailable: true,
        availability: {
          days: ['Monday', 'Wednesday', 'Friday'],
          timeSlots: ['10:00 AM - 12:00 PM', '2:00 PM - 4:00 PM'],
        },
        currentSession: [],
        matchedClients: [],
        avatarOverride: '',
        professionalInfo: {
          title: {
            code: 'LCP',
            display: 'Licensed Clinical Psychologist',
          },
          degrees: [
            {
              type: 'PhD',
              field: 'Clinical Psychology',
              institution: 'Harvard',
              year: 2012,
            },
            {
              type: 'MA',
              field: 'Counseling Psychology',
              institution: 'Stanford',
              year: 2009,
            },
          ],
          certifications: [
            { name: 'CBT Certified', issuer: 'APA', year: 2015 },
            { name: 'Anxiety Disorders Specialist', issuer: 'ABP', year: 2016 },
          ],
          experienceYears: 10,
        },
      },
      {
        profileId: 't002',
        name: 'Dr. Michael Chen',
        address: '456 Oak Avenue, Downtown, Country',
        languages: ['English', 'Mandarin'],
        specialization: ['Family Therapy', 'Couples Counseling'],
        consultationModes: ['Online', 'In-Person'],
        experience: 8,
        rating: 4.9,
        currency: 'USD',
        ragePrice: 100,
        isAvailable: false,
        availability: {
          days: ['Tuesday', 'Thursday', 'Saturday'],
          timeSlots: ['9:00 AM - 11:00 AM', '1:00 PM - 3:00 PM'],
        },
        currentSession: [],
        matchedClients: [],
        avatarOverride: '',
        professionalInfo: {
          title: {
            code: 'LMFT',
            display: 'Licensed Marriage and Family Therapist',
          },
          degrees: [
            {
              type: 'MS',
              field: 'Marriage and Family Therapy',
              institution: 'UCLA',
              year: 2014,
            },
            {
              type: 'BA',
              field: 'Psychology',
              institution: 'UC Berkeley',
              year: 2010,
            },
          ],
          certifications: [
            {
              name: 'Gottman Method Couples Therapy',
              issuer: 'Gottman Institute',
              year: 2016,
            },
            { name: 'EFT Certified', issuer: 'ICEEFT', year: 2017 },
          ],
          experienceYears: 8,
        },
      },
    ];
  }

  openOverlay(therapist: Therapist) {
    console.log('[SwipeDeck] Opening overlay for:', therapist);
    this.selectedTherapist = { ...therapist };
    this.isOverlayOpen = true;
    document.body.classList.add('overlay-open');
  }

  onCloseOverlay() {
    console.log('[SwipeDeck] Closing overlay');
    this.isOverlayOpen = false;
    this.selectedTherapist = null;
    document.body.classList.remove('overlay-open');
  }

  onCardSwiped(event: {
    direction: 'left' | 'right' | 'up';
    therapist: Therapist;
  }) {
    console.log('[SwipeDeck] Card swiped:', event.direction, event.therapist);

    if (event.direction === 'up') {
      this.openOverlay(event.therapist);
      return;
    }

    this.therapists = this.therapists.filter(
      (t) => t.profileId !== event.therapist.profileId
    );

    this.secureStorage.setItem('therapists', this.therapists);

    if (event.direction === 'right') {
      console.log('[SwipeDeck] Matched:', event.therapist.name);
    } else if (event.direction === 'left') {
      console.log('[SwipeDeck] Passed:', event.therapist.name);
    }
  }

  onBookAppointment(therapist: Therapist) {
    console.log('[SwipeDeck] Book appointment with:', therapist.name);
    this.onCloseOverlay();
  }

  onSendMessage(therapist: Therapist) {
    console.log('[SwipeDeck] Send message to:', therapist.name);
    this.onCloseOverlay();
  }

  get currentTherapist(): Therapist | null {
    return this.therapists.length > 0 ? this.therapists[0] : null;
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
