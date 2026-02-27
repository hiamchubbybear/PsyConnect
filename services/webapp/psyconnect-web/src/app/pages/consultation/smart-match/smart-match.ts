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
        const therapists = res.data || [];
        if (therapists.length === 0) {
          console.log('No recommendations found, triggering rescore...');
          this.swipeService.triggerUpdate().subscribe({
            next: () => {
              this.swipeService.getSwipeData().subscribe({
                next: (retryRes) => {
                  this.recommendedTherapists = retryRes.data || [];
                  this.isLoading = false;
                },
                error: (retryErr) => {
                  console.error('Failed to load recommendations after rescore', retryErr);
                  this.isLoading = false;
                }
              });
            },
            error: (triggerErr) => {
              console.error('Failed to trigger recommendation update', triggerErr);
              if (triggerErr.status === 503 || triggerErr.status === 500 || triggerErr.status === 404 || triggerErr.status === 401) {
                this.missingProfile = true;
              }
              this.isLoading = false;
            }
          });
        } else {
          this.recommendedTherapists = therapists;
          this.isLoading = false;
        }
      },
      error: (err) => {
        console.error('Failed to load recommendations', err);
        if (err.status === 503 || err.status === 500 || err.status === 404 || err.status === 401) {
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
