import { Routes } from '@angular/router';
import { DashboardComponent } from './dashboard/payment-dashboard';
import { Invoices } from './invoices/invoices';
import { Methods } from './methods/methods';
import { Overview } from './overview/overview';
import { PaymentManagementComponent } from './payment-management';
import { Refunds } from './refunds/refunds';
import { Reports } from './reports/reports';
import { Segments } from './segments/segments';
import { Settings } from './settings/settings';
import { Transactions } from './transactions/transactions';

export const paymentRoutes: Routes = [
  {
    path: '',
    component: PaymentManagementComponent,
    children: [
      { path: 'overview', component: Overview },
      { path: 'dashboard', component: DashboardComponent },
      { path: 'segments', component: Segments },
      { path: 'transactions', component: Transactions },
      { path: 'invoices', component: Invoices },
      { path: 'refunds', component: Refunds },
      { path: 'methods', component: Methods },
      { path: 'reports', component: Reports },
      { path: 'settings', component: Settings },
      { path: '', redirectTo: 'overview', pathMatch: 'full' },
    ],
  },
];
