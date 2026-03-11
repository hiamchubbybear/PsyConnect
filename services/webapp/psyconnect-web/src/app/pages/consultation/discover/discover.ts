import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
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

  constructor(
    private sessionService: SessionService,
    private translateService: TranslateService,
  ) {}

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

  getWeekDay(dateStr?: string): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    const dayIndex = date.getDay() === 0 ? 6 : date.getDay() - 1;
    const days = [
      'CONSULTATION.SESSION.Labels.Weekday.Monday',
      'CONSULTATION.SESSION.Labels.Weekday.Tuesday',
      'CONSULTATION.SESSION.Labels.Weekday.Wednesday',
      'CONSULTATION.SESSION.Labels.Weekday.Thursday',
      'CONSULTATION.SESSION.Labels.Weekday.Friday',
      'CONSULTATION.SESSION.Labels.Weekday.Saturday',
      'CONSULTATION.SESSION.Labels.Weekday.Sunday',
    ];
    return days[dayIndex];
  }

  getDayNum(dateStr?: string): string {
    if (!dateStr) return '';
    return new Date(dateStr).getDate().toString().padStart(2, '0');
  }

  getMonthShort(dateStr?: string): string {
    if (!dateStr) return '';
    const date = new Date(dateStr);
    return date.toLocaleDateString(this.translateService.currentLang || 'vi', {
      month: 'short',
    });
  }

  formatTime(dateStr?: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleTimeString(
      this.translateService.currentLang || 'vi-VN',
      {
        hour: '2-digit',
        minute: '2-digit',
      },
    );
  }

  formatCurrency(num: number): string {
    return new Intl.NumberFormat(this.translateService.currentLang || 'vi-VN').format(num) + ' VND';
  }

  getModeLabel(mode: string): string {
    const key = `CONSULTATION.SESSION.Mode.${
      mode === 'in_person' ? 'InPerson' : mode.charAt(0).toUpperCase() + mode.slice(1)
    }`;
    return this.translateService.instant(key) || mode;
  }

  getStatusLabel(status: string): string {
    const keyMap: Record<string, string> = {
      pending: 'CONSULTATION.SESSION.Status.Pending',
      active: 'CONSULTATION.SESSION.Status.Active',
      completed: 'CONSULTATION.SESSION.Status.Completed',
      cancelled: 'CONSULTATION.SESSION.Status.Cancelled',
    };
    return this.translateService.instant(keyMap[status] || status) || status;
  }

  get nextSessionTime(): string {
    return this.formatTime(this.overview?.next_session?.start_time);
  }
}
