// account.routes.ts
import { Routes } from '@angular/router';
import { UpdateAccountComponent } from './account';
// import { AccountLayoutComponent } from './account.component';
// import { AccountInfoComponent } from './pages/info/info.component';
// import { AccountPaymentComponent } from './pages/payment/payment.component';
// import { AccountSecurityComponent } from './pages/security/security.component';

export const accountRoutes: Routes = [
  {
    path: 'account',
    component: UpdateAccountComponent, // layout chứa sidebar + outlet
    children: [
      { path: 'info', component: UpdateAccountComponent },
      { path: 'security', component: UpdateAccountComponent },
      { path: 'payment', component: UpdateAccountComponent },
      { path: '', redirectTo: 'info', pathMatch: 'full' },
    ],
  },
];
