import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {
    CollapsibleSidebarComponent,
    SidebarItem,
} from '../../components/collapsible-sidebar/collapsible-sidebar';

@Component({
  selector: 'app-payment-management',
  standalone: true,
  imports: [CollapsibleSidebarComponent, CommonModule, RouterOutlet],
  template: `
    <app-collapsible-sidebar
      [title]="'Payment Management'"
      [titleTranslateKey]="'PAYMENT.Sidebar.Title'"
      [items]="sidebarItems"
      [width]="'240px'"
      [(collapsed)]="isSidebarCollapsed"
      (itemClick)="onSidebarItemClick($event)"
    />

    <div
      class="main-content"
      [style.margin-left]="isSidebarCollapsed ? '-30px' : '240px'"
    >
      <router-outlet></router-outlet>
    </div>
  `,
  styles: [
    `
      .main-content {
        margin-left: 360px;
        padding: 2rem;
        transition: margin-left 0.3s ease;
      }
    `,
  ],
})
export class PaymentManagementComponent {
  isSidebarCollapsed = false;

  sidebarItems: SidebarItem[] = [
    {
      label: 'Overview',
      route: 'overview',
      translateKey: 'PAYMENT.Sidebar.Items.Overview',
      exact: true,
    },
    {
      label: 'Dashboard',
      route: 'dashboard',
      translateKey: 'PAYMENT.Sidebar.Items.Dashboard',
    },
    {
      label: 'Segments',
      route: 'segments',
      translateKey: 'PAYMENT.Sidebar.Items.Segments',
    },
    {
      label: 'Transactions',
      route: 'transactions',
      translateKey: 'PAYMENT.Sidebar.Items.Transactions',
    },
    {
      label: 'Invoices',
      route: 'invoices',
      translateKey: 'PAYMENT.Sidebar.Items.Invoices',
    },
    {
      label: 'Refunds',
      route: 'refunds',
      translateKey: 'PAYMENT.Sidebar.Items.Refunds',
    },
    {
      label: 'Methods',
      route: 'methods',
      translateKey: 'PAYMENT.Sidebar.Items.Methods',
    },
    {
      label: 'Reports',
      route: 'reports',
      translateKey: 'PAYMENT.Sidebar.Items.Reports',
    },
    {
      label: 'Settings',
      route: 'settings',
      translateKey: 'PAYMENT.Sidebar.Items.Settings',
    },
  ];

  onSidebarItemClick(item: SidebarItem) {
    console.log('Sidebar item clicked:', item);
  }
}
