import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { SwipeCardComponent } from '../../components/swipe-card/swipe-card';
import { SecureStorageService } from '../../encrypt/secure';
import { Therapist } from '../../models/swipe-card';
import { SwipeService } from '../../services/swipe/swipe.service';

@Component({
  selector: 'app-swipe-deck',
  templateUrl: './swipe-deck.html',
  styleUrls: ['./swipe-deck.scss'],
  imports: [SwipeCardComponent, CommonModule],
})
export class SwipeDeckComponent implements OnInit {
  therapists: Therapist[] = [];
  currentIndex = 0;

  constructor(
    private swipeService: SwipeService,
    private secureStorage: SecureStorageService
  ) {}

  ngOnInit() {
    this.swipeService.getSwipeData().subscribe((res) => {
      if (res.status === 200 && Array.isArray(res.data)) {
        this.therapists = res.data.map((item: any) => ({
          id: item.profile_id,
          name: 'Unknown Name',
          title: 'Therapist',
          rating: item.rating,
          reviews: 0,
          specialties: item.specialization || [],
          nextAvailable: this.formatNextAvailable(item.availability),
          pricePerHour: `${item.rage_price} ${item.currency}`,
          avatarUrl:
            item.avatarUrl ||
            'https://psyconnect.chessy.dev/assets/images/avatar.jpeg',
          languages: item.languages,
          modes: item.consultation_modes,
        }));
        this.secureStorage.setItem('therapists', this.therapists);
        this.loadCurrentTherapist();
      }
    });
  }

  loadCurrentTherapist() {
    this.therapists =
      this.secureStorage.getItem<Therapist[]>('therapists') ?? [];
  }

  get currentTherapist(): Therapist | null {
    return this.therapists.length > 0 ? this.therapists[0] : null;
  }

  onCardSwiped(event: { direction: 'left' | 'right'; therapist: Therapist }) {
    this.therapists = this.therapists.filter(
      (t) => t.id !== event.therapist.id
    );
    console.log(this.therapists);

    localStorage.setItem('therapists', JSON.stringify(this.therapists));

    if (event.direction === 'right') {
      console.log('Matched', event.therapist);
    } else {
      console.log('Passed', event.therapist);
    }
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
