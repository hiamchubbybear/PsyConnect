import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { SwipeCardComponent } from '../../components/swipe-card/swipe-card';
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

  constructor(private swipeService: SwipeService) {}

  ngOnInit() {
    this.swipeService.getSwipeData().subscribe((data) => {
      this.therapists = [...data];
    });
  }

  onCardSwiped(event: { direction: 'left' | 'right'; therapist: Therapist }) {
    this.therapists = this.therapists.filter(
      (t) => t.id !== event.therapist.id
    );

    if (event.direction === 'right') {
      console.log('Matched', event.therapist);
    } else {
      console.log('Passed', event.therapist);
    }
  }
}
