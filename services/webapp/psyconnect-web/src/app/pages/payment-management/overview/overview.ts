import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

interface PaymentStats {
  totalRevenue: number;
  monthlyRevenue: number;
  totalTransactions: number;
  pendingPayments: number;
}

@Component({
  selector: 'app-overview',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  template: `
    <div class="overview-container">
      <header class="page-header">
        <h1>{{ 'PAYMENT.Overview.Title' | translate }}</h1>
        <p class="subtitle">{{ 'PAYMENT.Overview.Subtitle' | translate }}</p>
      </header>

      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon revenue">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><line x1="12" y1="1" x2="12" y2="23"></line><path d="M17 5H9.5a3.5 3.5 0 0 0 0 7h5a3.5 3.5 0 0 1 0 7H6"></path></svg>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.TotalRevenue' | translate }}</span>
            <span class="stat-value">\${{ stats.totalRevenue.toLocaleString() }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon monthly">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><rect x="3" y="4" width="18" height="18" rx="2" ry="2"></rect><line x1="16" y1="2" x2="16" y2="6"></line><line x1="8" y1="2" x2="8" y2="6"></line><line x1="3" y1="10" x2="21" y2="10"></line></svg>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.MonthlyRevenue' | translate }}</span>
            <span class="stat-value">\${{ stats.monthlyRevenue.toLocaleString() }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon transactions">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="16 3 21 3 21 8"></polyline><line x1="4" y1="14" x2="21" y2="3"></line><polyline points="8 21 3 21 3 16"></polyline><line x1="20" y1="10" x2="3" y2="21"></line></svg>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.Transactions' | translate }}</span>
            <span class="stat-value">{{ stats.totalTransactions }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon pending">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><circle cx="12" cy="12" r="10"></circle><polyline points="12 6 12 12 16 14"></polyline></svg>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.Pending' | translate }}</span>
            <span class="stat-value">{{ stats.pendingPayments }}</span>
          </div>
        </div>
      </div>

      <div class="chart-placeholder">
        <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"><polyline points="22 12 18 12 15 21 9 3 6 12 2 12"></polyline></svg>
        <p>{{ 'PAYMENT.Overview.ChartComingSoon' | translate }}</p>
      </div>
    </div>
  `,
  styles: [`
    .overview-container {
      padding: var(--spacing-2xl);
      max-width: 1400px;
      margin: 0 auto;
    }
    .page-header h1 {
      font-size: var(--font-size-3xl);
      font-weight: 700;
      color: var(--color-text);
      margin: 0 0 var(--spacing-sm) 0;
    }
    .subtitle {
      color: var(--color-text-muted);
      margin-bottom: var(--spacing-2xl);
    }
    .stats-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: var(--spacing-lg);
      margin-bottom: var(--spacing-2xl);
    }
    .stat-card {
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: var(--border-radius-lg);
      padding: var(--spacing-xl);
      display: flex;
      align-items: center;
      gap: var(--spacing-lg);
    }
    .stat-icon {
      width: 56px;
      height: 56px;
      border-radius: var(--border-radius-md);
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 1.5rem;
    }
    .stat-icon.revenue {
      background: rgba(16, 185, 129, 0.1);
      color: #10b981;
    }
    .stat-icon.monthly {
      background: rgba(59, 130, 246, 0.1);
      color: #3b82f6;
    }
    .stat-icon.transactions {
      background: rgba(139, 92, 246, 0.1);
      color: #8b5cf6;
    }
    .stat-icon.pending {
      background: rgba(245, 158, 11, 0.1);
      color: #f59e0b;
    }
    .stat-content {
      display: flex;
      flex-direction: column;
      gap: var(--spacing-xs);
    }
    .stat-label {
      font-size: var(--font-size-sm);
      color: var(--color-text-muted);
      font-weight: 500;
    }
    .stat-value {
      font-size: var(--font-size-2xl);
      font-weight: 700;
      color: var(--color-text);
    }
    .chart-placeholder {
      background: var(--color-surface);
      border: 1px solid var(--color-border);
      border-radius: var(--border-radius-lg);
      padding: var(--spacing-3xl);
      text-align: center;
    }
    .chart-placeholder i {
      font-size: 4rem;
      color: var(--color-text-muted);
      opacity: 0.3;
      margin-bottom: var(--spacing-lg);
    }
    .chart-placeholder p {
      color: var(--color-text-muted);
    }
  `]
})
export class Overview implements OnInit {
  stats: PaymentStats = {
    totalRevenue: 0,
    monthlyRevenue: 0,
    totalTransactions: 0,
    pendingPayments: 0
  };

  ngOnInit() {
    this.loadStats();
  }

  loadStats() {
    
    this.stats = {
      totalRevenue: 45280,
      monthlyRevenue: 12450,
      totalTransactions: 156,
      pendingPayments: 8
    };
  }
}
