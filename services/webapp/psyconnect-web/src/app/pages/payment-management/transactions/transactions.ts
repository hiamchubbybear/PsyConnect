import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import {
  MockDataService,
  MockTransaction,
} from '../../../services/mock-data.service';
import { PsyEmptyStateComponent } from '../../../shared/ui-atoms/empty-state/psy-empty-state.component';
import { TableComponent } from '../../../shared/ui-atoms/table/table.component';
import { IconComponent } from '../../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-transactions',
  standalone: true,
  imports: [
    CommonModule,
    TranslateModule,
    FormsModule,
    TableComponent,
    PsyEmptyStateComponent,
    IconComponent,
  ],
  templateUrl: './transactions.html',
  styleUrl: './transactions.scss',
})
export class Transactions implements OnInit {
  transactions: MockTransaction[] = [];
  filteredTransactions: MockTransaction[] = [];
  searchQuery = '';
  selectedFilter: 'all' | 'completed' | 'pending' | 'failed' = 'all';
  selectedType: 'all' | 'payment' | 'refund' | 'subscription' = 'all';

  constructor(private mockDataService: MockDataService) {}

  ngOnInit() {
    this.loadTransactions();
  }

  loadTransactions() {
    this.mockDataService.getTransactions().subscribe((transactions) => {
      this.transactions = transactions;
      this.applyFilters();
    });
  }

  applyFilters() {
    let filtered = [...this.transactions];

    if (this.selectedFilter !== 'all') {
      filtered = filtered.filter((t) => t.status === this.selectedFilter);
    }

    if (this.selectedType !== 'all') {
      filtered = filtered.filter((t) => t.type === this.selectedType);
    }

    if (this.searchQuery) {
      const query = this.searchQuery.toLowerCase();
      filtered = filtered.filter(
        (t) =>
          t.description.toLowerCase().includes(query) ||
          t.customer.toLowerCase().includes(query),
      );
    }

    this.filteredTransactions = filtered;
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'completed':
        return 'status-completed';
      case 'pending':
        return 'status-pending';
      case 'failed':
        return 'status-failed';
      default:
        return '';
    }
  }

  getTypeIcon(type: string): string {
    switch (type) {
      case 'payment':
        return 'arrow-down';
      case 'refund':
        return 'undo';
      case 'subscription':
        return 'repeat';
      default:
        return 'help-circle';
    }
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }
}
