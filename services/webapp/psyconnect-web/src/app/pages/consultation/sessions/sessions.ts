import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import { ConsultationSession } from '../../../models/consultation.model';
import { SessionService } from '../../../services/consultation/session.service';

type SessionTab = 'upcoming' | 'completed' | 'cancelled';

@Component({
  selector: 'consultation-sessions',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './sessions.html',
  styleUrl: './sessions.scss',
})
export class Sessions implements OnInit {
  sessions: ConsultationSession[] = [];
  isLoading = true;
  activeTab: SessionTab = 'upcoming';

  constructor(
    private sessionService: SessionService,
    private router: Router,
    private translateService: TranslateService,
  ) {}

  ngOnInit(): void {
    this.loadSessions();
  }

  loadSessions(): void {
    this.isLoading = true;
    this.sessionService.getAllSessions().subscribe({
      next: (data: any) => {
        
        this.sessions = Array.isArray(data)
          ? data
          : data?.data || data?.sessions || [];
        this.isLoading = false;
      },
      error: () => {
        this.sessions = [];
        this.isLoading = false;
      },
    });
  }

  setTab(tab: SessionTab) {
    this.activeTab = tab;
  }

  get filteredSessions(): ConsultationSession[] {
    const now = new Date();
    switch (this.activeTab) {
      case 'upcoming':
        return this.sessions.filter(
          (s) =>
            s.status === 'pending' ||
            s.status === 'active' ||
            (s.start_time && new Date(s.start_time) > now),
        );
      case 'completed':
        return this.sessions.filter((s) => s.status === 'completed');
      case 'cancelled':
        return this.sessions.filter((s) => s.status === 'cancelled');
      default:
        return this.sessions;
    }
  }

  get upcomingCount() {
    return this.sessions.filter(
      (s) => s.status === 'pending' || s.status === 'active',
    ).length;
  }

  get completedCount() {
    return this.sessions.filter((s) => s.status === 'completed').length;
  }

  get cancelledCount() {
    return this.sessions.filter((s) => s.status === 'cancelled').length;
  }

  getStatusClass(status: string): string {
    const map: Record<string, string> = {
      pending: 'status--pending',
      active: 'status--active',
      completed: 'status--completed',
      cancelled: 'status--cancelled',
    };
    return map[status] || '';
  }

  getModeIcon(mode: string): string {
    switch (mode) {
      case 'online':
        return '💻';
      case 'in_person':
        return '🏥';
      case 'phone':
        return '📞';
      case 'chat':
        return '💬';
      default:
        return '📋';
    }
  }

  formatDate(dateStr: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleDateString('vi-VN', {
      day: '2-digit',
      month: '2-digit',
      year: 'numeric',
    });
  }

  formatTime(dateStr: string): string {
    if (!dateStr) return '—';
    return new Date(dateStr).toLocaleTimeString('vi-VN', {
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  getPriceLabel(price: number | undefined): string {
    if (!price || price <= 0) {
      return this.translateService.instant('CONSULTATION.Booking.Free');
    }
    return `${new Intl.NumberFormat('vi-VN').format(price)} VND`;
  }

  joinSession(session: ConsultationSession): void {
    if (session.location_info?.link) {
      window.open(session.location_info.link, '_blank');
    }
  }

  getPaymentUrl(session: ConsultationSession): void {
    this.sessionService.getPaymentUrl(session.session_id || '').subscribe({
      next: (res) => window.open(res.payment_url, '_blank'),
      error: (err) => console.error('Payment URL error', err),
    });
  }

  cancelSession(session: ConsultationSession): void {
    if (
      !confirm(
        this.translateService.instant('CONSULTATION.Booking.CancelConfirm'),
      )
    ) {
      return;
    }
    this.sessionService.cancelSession(session.session_id || '').subscribe({
      next: () => this.loadSessions(),
      error: (err) => console.error('Cancel error', err),
    });
  }
}
