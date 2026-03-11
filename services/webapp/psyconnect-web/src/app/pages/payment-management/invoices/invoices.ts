import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';
import {
  MockDataService,
  MockInvoice,
} from '../../../services/mock-data.service';
import { TableComponent } from '../../../shared/ui-atoms/table/table.component';
import { IconComponent } from '../../../shared/ui-atoms/icon/icon.component';

@Component({
  selector: 'app-invoices',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule, TableComponent, IconComponent],
  templateUrl: './invoices.html',
  styleUrl: './invoices.scss',
})
export class Invoices implements OnInit {
  invoices: MockInvoice[] = [];
  filteredInvoices: MockInvoice[] = [];
  searchQuery = '';
  selectedFilter: 'all' | 'paid' | 'unpaid' | 'overdue' = 'all';

  constructor(private mockDataService: MockDataService) {}

  ngOnInit() {
    this.loadInvoices();
  }

  loadInvoices() {
    this.mockDataService.getInvoices().subscribe((invoices) => {
      this.invoices = invoices;
      this.applyFilters();
    });
  }

  applyFilters() {
    let filtered = [...this.invoices];

    if (this.selectedFilter !== 'all') {
      filtered = filtered.filter((i) => i.status === this.selectedFilter);
    }

    if (this.searchQuery) {
      const query = this.searchQuery.toLowerCase();
      filtered = filtered.filter(
        (i) =>
          i.invoiceNumber.toLowerCase().includes(query) ||
          i.customer.toLowerCase().includes(query),
      );
    }

    this.filteredInvoices = filtered;
  }

  getStatusClass(status: string): string {
    switch (status) {
      case 'paid':
        return 'status-paid';
      case 'unpaid':
        return 'status-unpaid';
      case 'overdue':
        return 'status-overdue';
      default:
        return '';
    }
  }

  formatDate(date: string): string {
    return new Date(date).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'short',
      day: 'numeric',
    });
  }

  downloadInvoice(invoice: MockInvoice) {
    console.log('Download invoice:', invoice.invoiceNumber);
  }
}
