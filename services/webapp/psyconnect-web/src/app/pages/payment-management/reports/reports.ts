import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

interface ReportData {
  period: string;
  revenue: number;
  transactions: number;
  refunds: number;
  netRevenue: number;
}

@Component({
  selector: 'app-reports',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './reports.html',
  styleUrl: './reports.scss'
})
export class Reports implements OnInit {
  reports: ReportData[] = [];
  selectedPeriod: 'week' | 'month' | 'year' = 'month';

  summary = {
    totalRevenue: 0,
    totalTransactions: 0,
    totalRefunds: 0,
    netRevenue: 0
  };

  ngOnInit() {
    this.loadReports();
  }

  loadReports() {
    this.reports = [
      {
        period: 'January 2024',
        revenue: 12450.00,
        transactions: 45,
        refunds: 225.00,
        netRevenue: 12225.00
      },
      {
        period: 'February 2024',
        revenue: 15680.00,
        transactions: 52,
        refunds: 150.00,
        netRevenue: 15530.00
      },
      {
        period: 'March 2024',
        revenue: 18920.00,
        transactions: 63,
        refunds: 300.00,
        netRevenue: 18620.00
      },
      {
        period: 'April 2024',
        revenue: 14230.00,
        transactions: 48,
        refunds: 175.00,
        netRevenue: 14055.00
      }
    ];

    this.calculateSummary();
  }

  calculateSummary() {
    this.summary = {
      totalRevenue: this.reports.reduce((sum, r) => sum + r.revenue, 0),
      totalTransactions: this.reports.reduce((sum, r) => sum + r.transactions, 0),
      totalRefunds: this.reports.reduce((sum, r) => sum + r.refunds, 0),
      netRevenue: this.reports.reduce((sum, r) => sum + r.netRevenue, 0)
    };
  }

  exportReport() {
    console.log('Exporting report...');
  }
}
