import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

interface Refund {
  id: string;
  transactionId: string;
  amount: number;
  reason: string;
  status: 'pending' | 'approved' | 'rejected';
  requestDate: string;
  customer: string;
}

@Component({
  selector: 'app-refunds',
  standalone: true,
  imports: [CommonModule, TranslateModule, FormsModule],
  templateUrl: './refunds.html',
  styleUrl: './refunds.scss'
})
export class Refunds implements OnInit {
  refunds: Refund[] = [];
  filteredRefunds: Refund[] = [];
  searchQuery = '';
  selectedFilter: 'all' | 'pending' | 'approved' | 'rejected' = 'all';

  ngOnInit() {
    this.loadRefunds();
  }

  loadRefunds() {
    this.refunds = [
      {
        id: 'ref1',
        transactionId: 'TXN-2024-001',
        amount: 150.00,
        reason: 'Session cancelled by therapist',
        status: 'approved',
        requestDate: new Date(Date.now() - 2 * 24 * 60 * 60 * 1000).toISOString(),
        customer: 'John Doe'
      },
      {
        id: 'ref2',
        transactionId: 'TXN-2024-002',
        amount: 75.00,
        reason: 'Duplicate payment',
        status: 'pending',
        requestDate: new Date(Date.now() - 1 * 24 * 60 * 60 * 1000).toISOString(),
        customer: 'Jane Smith'
      },
      {
        id: 'ref3',
        transactionId: 'TXN-2024-003',
        amount: 200.00,
        reason: 'Service not provided',
        status: 'rejected',
        requestDate: new Date(Date.now() - 5 * 24 * 60 * 60 * 1000).toISOString(),
        customer: 'Mike Johnson'
      },
      {
        id: 'ref4',
        transactionId: 'TXN-2024-004',
        amount: 100.00,
        reason: 'Customer request',
        status: 'pending',
        requestDate: new Date(Date.now() - 3 * 24 * 60 * 60 * 1000).toISOString(),
        customer: 'Sarah Williams'
      }
    ];
    this.applyFilters();
  }

  applyFilters() {
    let filtered = [...this.refunds];

    if (this.selectedFilter !== 'all') {
      filtered = filtered.filter(r => r.status === this.selectedFilter);
    }

    if (this.searchQuery) {
      const query = this.searchQuery.toLowerCase();
      filtered = filtered.filter(r =>
        r.transactionId.toLowerCase().includes(query) ||
        r.customer.toLowerCase().includes(query) ||
        r.reason.toLowerCase().includes(query)
      );
    }

    this.filteredRefunds = filtered;
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'approved': return 'status-approved';
      case 'pending': return 'status-pending';
      case 'rejected': return 'status-rejected';
      default: return '';
    }
  }

  approveRefund(refund: Refund) {
    refund.status = 'approved';
    console.log('Approved refund:', refund.id);
  }

  rejectRefund(refund: Refund) {
    refund.status = 'rejected';
    console.log('Rejected refund:', refund.id);
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    });
  }
}
