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
            <i class="fas fa-dollar-sign"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.TotalRevenue' | translate }}</span>
            <span class="stat-value">\${{ stats.totalRevenue.toLocaleString() }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon monthly">
            <i class="fas fa-calendar-alt"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.MonthlyRevenue' | translate }}</span>
            <span class="stat-value">\${{ stats.monthlyRevenue.toLocaleString() }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon transactions">
            <i class="fas fa-exchange-alt"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.Transactions' | translate }}</span>
            <span class="stat-value">{{ stats.totalTransactions }}</span>
          </div>
        </div>
        <div class="stat-card">
          <div class="stat-icon pending">
            <i class="fas fa-clock"></i>
          </div>
          <div class="stat-content">
            <span class="stat-label">{{ 'PAYMENT.Overview.Pending' | translate }}</span>
            <span class="stat-value">{{ stats.pendingPayments }}</span>
          </div>
        </div>
      </div>

      <div class="chart-placeholder">
        <i class="fas fa-chart-line"></i>
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
    // Mock data
    this.stats = {
      totalRevenue: 45280,
      monthlyRevenue: 12450,
      totalTransactions: 156,
      pendingPayments: 8
    };
  }
}
