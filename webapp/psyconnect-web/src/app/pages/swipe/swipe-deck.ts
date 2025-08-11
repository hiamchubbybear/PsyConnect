import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { SwipeCardComponent } from '../../components/swipe-card/swipe-card';
import { Therapist } from '../../models/swipe-card';

@Component({
  selector: 'app-swipe-deck',
  templateUrl: './swipe-deck.html',
  styleUrls: ['./swipe-deck.scss'],
  imports: [SwipeCardComponent, CommonModule],
})
export class SwipeDeckComponent {
  therapists: Therapist[] = [
    {
      id: '1',
      name: 'Dr. Sarah Johnson',
      title: 'Licensed Clinical Psychologist',
      rating: 4.9,
      reviews: 127,
      specialties: ['Anxiety', 'Depression', 'PTSD'],
      nextAvailable: 'Tomorrow at 2:00 PM',
      pricePerHour: '$120 per hour',
      avatarUrl: '',
    },
    {
      id: '2',
      name: 'Dr. Mark Lee',
      title: 'Counseling Psychologist',
      rating: 4.7,
      reviews: 86,
      specialties: ['Stress', 'Relationships'],
      nextAvailable: 'Today at 5:30 PM',
      pricePerHour: '$95 per hour',
      avatarUrl: '',
    },
  ];

  liked: Therapist[] = [];
  passed: Therapist[] = [];

  onCardSwiped(event: { direction: 'left' | 'right'; therapist: Therapist }) {
    const tIndex = this.therapists.findIndex(
      (t) => t.id === event.therapist.id
    );

    if (tIndex > -1) {
      this.therapists.splice(tIndex, 1);
    }

    if (event.direction === 'right') {
      this.liked.push(event.therapist);
      console.log('Matched', event.therapist);
    } else {
      this.passed.push(event.therapist);
      console.log('Passed', event.therapist);
    }
  }
}
