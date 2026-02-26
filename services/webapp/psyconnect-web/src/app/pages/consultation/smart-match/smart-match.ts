import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { MatchRequest } from '../../../models/consultation.model';
import { Therapist } from '../../../models/swipe-card';
import { MatchingService } from '../../../services/consultation/matching.service';
import { SwipeService } from '../../../services/swipe/swipe.service';

@Component({
  selector: 'consultation-smart-match',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './smart-match.html',
  styleUrls: ['./smart-match.scss'],
})
export class SmartMatchComponent implements OnInit {
  recommendedTherapists: Therapist[] = [];
  isLoading = true;
  missingProfile = false;

  constructor(
    private swipeService: SwipeService,
    private matchingService: MatchingService,
    private router: Router,
  ) {}

  ngOnInit(): void {
    this.fetchRecommendations();
  }

  fetchRecommendations() {
    this.isLoading = true;
    this.missingProfile = false;
    this.swipeService.getSwipeData().subscribe({
      next: (res) => {
        this.recommendedTherapists = res.data || [];
        this.isLoading = false;
      },
      error: (err) => {
        console.error('Failed to load recommendations', err);
        if (err.status === 503) {
          this.missingProfile = true;
        }
        this.isLoading = false;
      },
    });
  }

  goToCreateProfile() {
    this.router.navigate(['/feature/consultation/create-profile']);
  }

  onSwipeMatch(therapistId: string) {
    console.log('Sending SwipeAndMatch API request for:', therapistId);
    const request: MatchRequest = {
      therapist_id: therapistId,
    };

    this.matchingService.requestMatch(request).subscribe({
      next: (res) => {
        console.log('Match successfully requested:', res);
        // Navigate to the scheduler calendar passing the therapist ID if needed
        this.router.navigate(['/feature/consultation/schedules']);
      },
      error: (err) => {
        console.error('Failed to request match', err);
      },
    });

    // Optimistically remove from list
    this.recommendedTherapists = this.recommendedTherapists.filter(
      (t) => t.profileId !== therapistId,
    );
  }

  onSwipePass(therapistId: string) {
    console.log('Passing on therapist:', therapistId);
    // Ideally call swipeService.swipe(therapistId, 'rejected') here
    // For now, optimistically log and hide from UI
    this.recommendedTherapists = this.recommendedTherapists.filter(
      (t) => t.profileId !== therapistId,
    );
  }
}
