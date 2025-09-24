import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from './guards/auth-guard';
import { guestGuard } from './guards/guest-guard';
import { PasswordResetGuard } from './guards/password-reset-guard';
import { RequestResetSuccessGuard } from './guards/request-reset-guard';
import { AboutUsComponent } from './pages/about-us/about-us';
import { UpdateAccountComponent } from './pages/account/account-update';
import { PaymentSessionComponent } from './pages/account/payment/payment';
import { ProfileSectionComponent } from './pages/account/profile/profile-update';
import { SecuritySectionComponent } from './pages/account/security/security-update';
import { ChatComponent } from './pages/chat/chatpage/chatpage';
import { Consultation } from './pages/consultation/consultation';
import { Homepage } from './pages/homepage/homepage';
import { Login } from './pages/login/login';
import { Notfound } from './pages/notfound/notfound';
import { ResetPasswordComponent } from './pages/password-reset/password-reset';
import { PaymentManagement } from './pages/payment-management/payment-management';
import { RequestResetComponent } from './pages/request-reset/request-reset';
import { RequestResetSuccessComponent } from './pages/request-reset/success/request-reset-success';
import { MultiStepRegisterComponent } from './pages/signup/signup';
import { Start } from './pages/start/start';

export const routes: Routes = [
  { path: '', component: Homepage, canActivate: [guestGuard] },
  {
    path: 'auth/signup',
    component: MultiStepRegisterComponent,
    canActivate: [guestGuard],
  },
  { path: 'auth/login', component: Login, canActivate: [guestGuard] },
  { path: 'oauth2/callback', component: Login, canActivate: [guestGuard] },

  { path: 'chat-page', component: ChatComponent, canActivate: [authGuard] },
  {
    path: 'auth/password-reset',
    component: ResetPasswordComponent,
    canActivate: [guestGuard, PasswordResetGuard],
  },
  {
    path: 'auth/request-reset',
    component: RequestResetComponent,
    canActivate: [guestGuard],
  },
  {
    path: 'auth/request-reset-success',
    component: RequestResetSuccessComponent,
    canActivate: [guestGuard, RequestResetSuccessGuard],
  },
  {
    path: 'account',
    component: UpdateAccountComponent,
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'profile', pathMatch: 'full' },
      { path: 'profile', component: ProfileSectionComponent },
      { path: 'security', component: SecuritySectionComponent },
      { path: 'payment', component: PaymentSessionComponent },
    ],
  },
  { path: 'feature/feed', component: Notfound, canActivate: [authGuard] },
  {
    path: 'feature/schedule',
    component: Consultation,
    canActivate: [authGuard],
    children: [
      {
        path: 'discover',
        loadComponent: () =>
          import('./pages/consultation/discover/discover').then(
            (m) => m.Discover
          ),
      },
      {
        path: 'sessions',
        loadComponent: () =>
          import('./pages/consultation/sessions/sessions').then(
            (m) => m.Sessions
          ),
      },
      {
        path: 'schedule',
        loadComponent: () =>
          import('./pages/consultation/schedules/schedules').then(
            (m) => m.Schedules
          ),
      },
      {
        path: 'history',
        loadComponent: () =>
          import('./pages/consultation/history/history').then((m) => m.History),
      },
      {
        path: 'reviews',
        loadComponent: () =>
          import('./pages/consultation/reviews/reviews').then((m) => m.Reviews),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/consultation/settings/settings').then(
            (m) => m.Settings
          ),
      },
      { path: '', redirectTo: 'discover', pathMatch: 'full' }, // default
    ],
  },
  { path: 'feature/chat', component: Notfound, canActivate: [authGuard] },
  { path: 'feature/profile', component: Notfound, canActivate: [authGuard] },
  { path: 'about-us', component: AboutUsComponent },
  {
    path: 'payment',
    component: PaymentManagement,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./pages/payment-management/dashboard/payment-dashboard').then(
            (m) => m.PaymentDashboardComponent
          ),
      },
      {
        path: 'overview',
        loadComponent: () =>
          import('./pages/payment-management/overview/overview').then(
            (m) => m.Overview
          ),
      },
      {
        path: 'invoices',
        loadComponent: () =>
          import('./pages/payment-management/invoices/invoices').then(
            (m) => m.Invoices
          ),
      },
      {
        path: 'methods',
        loadComponent: () =>
          import('./pages/payment-management/methods/methods').then(
            (m) => m.Methods
          ),
      },
      {
        path: 'refunds',
        loadComponent: () =>
          import('./pages/payment-management/refunds/refunds').then(
            (m) => m.Refunds
          ),
      },
      {
        path: 'reports',
        loadComponent: () =>
          import('./pages/payment-management/reports/reports').then(
            (m) => m.Reports
          ),
      },
      {
        path: 'segments',
        loadComponent: () =>
          import('./pages/payment-management/segments/segments').then(
            (m) => m.Segments
          ),
      },
      {
        path: 'transactions',
        loadComponent: () =>
          import('./pages/payment-management/transactions/transactions').then(
            (m) => m.Transactions
          ),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import('./pages/payment-management/settings/settings').then(
            (m) => m.Settings
          ),
      },
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
    ],
  },
  { path: 'start', component: Start },
  { path: '**', component: Notfound },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
