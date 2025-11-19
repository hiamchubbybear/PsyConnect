import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { FormsModule } from '@angular/forms';

interface Invoice {
  invoiceNumber: string;
  vendor: string;
  billingDate: string;
  status: 'Paid' | 'Unpaid';
  amount: string;
}

@Component({
  selector: 'app-overview',
  standalone: true,
  imports: [CommonModule, FormsModule],
  templateUrl: './overview.html',
  styleUrls: ['./overview.scss'],
})
export class Overview {
  allInvoices: Invoice[] = [
    {
      invoiceNumber: '514684654865',
      vendor: 'Jane Cooper',
      billingDate: '2/19/21',
      status: 'Paid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '5467319467348',
      vendor: 'Wade Warren',
      billingDate: '5/7/16',
      status: 'Paid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '1345705945446',
      vendor: 'Esther Howard',
      billingDate: '9/18/16',
      status: 'Unpaid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '5440754979',
      vendor: 'Cameron Williamson',
      billingDate: '2/11/12',
      status: 'Paid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '1234567984543',
      vendor: 'Brooklyn Simmons',
      billingDate: '9/18/16',
      status: 'Unpaid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '8454134649707',
      vendor: 'Leslie Alexander',
      billingDate: '1/28/17',
      status: 'Unpaid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '2130164040451',
      vendor: 'Jenny Wilson',
      billingDate: '5/27/15',
      status: 'Paid',
      amount: '$500.00',
    },
    {
      invoiceNumber: '043910464504',
      vendor: 'Guy Hawkins',
      billingDate: '8/2/19',
      status: 'Paid',
      amount: '$500.00',
    },

    {
      invoiceNumber: '88937465231',
      vendor: 'Marvin McKinney',
      billingDate: '6/14/20',
      status: 'Paid',
      amount: '$400.00',
    },
    {
      invoiceNumber: '98234756234',
      vendor: 'Courtney Henry',
      billingDate: '3/3/21',
      status: 'Unpaid',
      amount: '$750.00',
    },
    {
      invoiceNumber: '12398745623',
      vendor: 'Devon Lane',
      billingDate: '4/25/18',
      status: 'Paid',
      amount: '$600.00',
    },
  ];

  displayedInvoices: Invoice[] = [];
  itemsToShow = 6;

  constructor() {
    this.displayedInvoices = this.allInvoices.slice(0, this.itemsToShow);
  }

  loadMore() {
    this.itemsToShow += 3;
    this.displayedInvoices = this.allInvoices.slice(0, this.itemsToShow);
  }

  canLoadMore(): boolean {
    return this.displayedInvoices.length < this.allInvoices.length;
  }
}
