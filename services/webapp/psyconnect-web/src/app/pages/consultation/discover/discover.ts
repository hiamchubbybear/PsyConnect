import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import {
  OverviewData,
  SessionService,
} from '../../../services/consultation/session.service';

@Component({
  selector: 'consultation-discover',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule],
  templateUrl: './discover.html',
  styleUrls: ['./discover.scss'],
})
export class Discover implements OnInit {
  overview: OverviewData | null = null;
  isLoading = true;

  constructor(private sessionService: SessionService) {}

  ngOnInit(): void {
    this.loadOverview();
  }

  loadOverview(): void {
    this.sessionService.getOverview().subscribe({
      next: (data) => {
        this.overview = data;
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      },
    });
  }

  formatDate(dateStr?: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('vi-VN', {
      weekday: 'short',
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  }

  formatTime(dateStr?: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleTimeString('vi-VN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  formatCurrency(num: number): string {
    return new Intl.NumberFormat('vi-VN').format(num) + ' VND';
  }

  get nextSessionDate(): string {
    return this.formatDate(this.overview?.next_session?.start_time);
  }

  get nextSessionTime(): string {
    return this.formatTime(this.overview?.next_session?.start_time);
  }
}
