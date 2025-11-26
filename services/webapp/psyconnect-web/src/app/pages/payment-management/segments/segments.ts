import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

interface CustomerSegment {
  id: string;
  name: string;
  customerCount: number;
  totalRevenue: number;
  avgTransactionValue: number;
  color: string;
}

@Component({
  selector: 'app-segments',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './segments.html',
  styleUrl: './segments.scss'
})
export class Segments implements OnInit {
  segments: CustomerSegment[] = [];

  ngOnInit() {
    this.loadSegments();
  }

  loadSegments() {
    this.segments = [
      {
        id: '1',
        name: 'Premium Clients',
        customerCount: 45,
        totalRevenue: 25680.00,
        avgTransactionValue: 570.67,
        color: '#8b5cf6'
      },
      {
        id: '2',
        name: 'Regular Clients',
        customerCount: 128,
        totalRevenue: 18920.00,
        avgTransactionValue: 147.81,
        color: '#3b82f6'
      },
      {
        id: '3',
        name: 'New Clients',
        customerCount: 67,
        totalRevenue: 8450.00,
        avgTransactionValue: 126.12,
        color: '#10b981'
      },
      {
        id: '4',
        name: 'Inactive Clients',
        customerCount: 34,
        totalRevenue: 2150.00,
        avgTransactionValue: 63.24,
        color: '#f59e0b'
      }
    ];
  }

  getPercentage(segment: CustomerSegment): number {
    const total = this.segments.reduce((sum, s) => sum + s.customerCount, 0);
    return (segment.customerCount / total) * 100;
  }
}
