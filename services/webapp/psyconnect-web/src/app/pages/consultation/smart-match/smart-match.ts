import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { MatDialog } from '@angular/material/dialog';
import { Router } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { Therapist } from '../../../models/swipe-card';
import { AuthService } from '../../../services/auth/auth.service';
import { MatchingService } from '../../../services/consultation/matching.service';
import { SwipeService } from '../../../services/swipe/swipe.service';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'consultation-smart-match',
  standalone: true,
  imports: [CommonModule, AvatarFallbackPipe, TranslateModule],
  templateUrl: './smart-match.html',
  styleUrls: ['./smart-match.scss'],
})
export class SmartMatchComponent implements OnInit {
  recommendedTherapists: Therapist[] = [];
  currentIndex = 0;
  isLoading = true;
  missingProfile = false;
  isTherapist = false;

  constructor(
    private swipeService: SwipeService,
    private matchingService: MatchingService,
    private authService: AuthService,
    private router: Router,
    private dialog: MatDialog,
  ) {}

  get currentTherapist(): Therapist | null {
    return this.recommendedTherapists[this.currentIndex] ?? null;
  }

  get totalCount(): number {
    return this.recommendedTherapists.length;
  }

  ngOnInit(): void {
    if (this.authService.isTherapist()) {
      this.isTherapist = true;
      this.isLoading = false;
      return;
    }
    this.fetchRecommendations();
  }

  fetchRecommendations() {
    this.isLoading = true;
    this.missingProfile = false;
    this.swipeService.getSwipeData().subscribe({
      next: (res) => {
        const therapists = res.data || [];
        if (therapists.length === 0) {
          this.swipeService.triggerUpdate().subscribe({
            next: () => {
              this.swipeService.getSwipeData().subscribe({
                next: (retryRes) => {
                  this.recommendedTherapists = retryRes.data || [];
                  this.currentIndex = 0;
                  this.isLoading = false;
                },
                error: () => {
                  this.isLoading = false;
                },
              });
            },
            error: (triggerErr) => {
              if ([503, 500, 404, 401].includes(triggerErr.status)) {
                this.missingProfile = true;
              }
              this.isLoading = false;
            },
          });
        } else {
          this.recommendedTherapists = therapists;
          this.currentIndex = 0;
          this.isLoading = false;
        }
      },
      error: (err) => {
        if ([503, 500, 404, 401].includes(err.status)) {
          this.missingProfile = true;
        }
        this.isLoading = false;
      },
    });
  }

  goToCreateProfile() {
    this.router.navigate(['/feature/consultation/create-profile']);
  }

  goToSessions() {
    this.router.navigate(['/feature/consultation/sessions']);
  }

  onPass() {
    if (!this.currentTherapist) return;

    this.recommendedTherapists.splice(this.currentIndex, 1);
    if (this.currentIndex >= this.recommendedTherapists.length) {
      this.currentIndex = 0;
    }
  }

  onBookSession() {
    if (!this.currentTherapist) return;
    const therapistId = this.currentTherapist.profileId;
    this.router.navigate(['/feature/consultation/book', therapistId]);
  }

  onMessage() {
    if (!this.currentTherapist) return;

    this.router.navigate(['/feature/chat'], {
      queryParams: { therapistId: this.currentTherapist.profileId },
    });
  }

  onViewProfile() {
    if (!this.currentTherapist) return;

    this.router.navigate([
      '/feature/consultation/therapist',
      this.currentTherapist.profileId,
    ]);
  }
}
