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
import { ConsultationComponent } from './pages/consultation/consultation';
import { FeedComponent } from './pages/feed/feed';
import { Login } from './pages/login/login';
import { Notfound } from './pages/notfound/notfound';
import { ResetPasswordComponent } from './pages/password-reset/password-reset';
import { PaymentManagementComponent } from './pages/payment-management/payment-management';
import { PostManagementComponent } from './pages/post-management/post-managements';
import { RequestResetComponent } from './pages/request-reset/request-reset';
import { RequestResetSuccessComponent } from './pages/request-reset/success/request-reset-success';
import { SearchComponent } from './pages/search/search';
import { MultiStepRegisterComponent } from './pages/signup/signup';
import { FriendsPage } from './pages/socialprofile/friends/friends-page';
import { SocialProfile } from './pages/socialprofile/socialprofile';
import { ProfilePageComponent } from './pages/socialprofile/wall/wall';
import { Start } from './pages/start/start';

export const routes: Routes = [
  { path: '', component: Start, canActivate: [guestGuard] },
  {
    path: 'auth/signup',
    component: MultiStepRegisterComponent,
    canActivate: [guestGuard],
  },
  { path: 'auth/login', component: Login, canActivate: [guestGuard] },
  { path: 'oauth2/callback', component: Login, canActivate: [guestGuard] },

  {
    path: 'profile/:id',
    component: ProfilePageComponent,
    canActivate: [authGuard],
  },
  { path: 'chat-page', component: ChatComponent, canActivate: [authGuard] },
  {
    path: 'chat/:id',
    component: ChatComponent,
    canActivate: [authGuard],
  },
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
  { path: 'feature/feed', component: FeedComponent, canActivate: [authGuard] },
  {
    path: 'feed/post/:id',
    loadComponent: () =>
      import('./pages/feed/post-detail/post-detail').then(
        (m) => m.PostDetailComponent
      ),
    canActivate: [authGuard],
  },
  { path: 'feature/search', component: SearchComponent },
  {
    path: 'feature/article',
    component: PostManagementComponent,
    children: [
      {
        path: 'all',
        loadComponent: () =>
          import('../app/pages/post-management/all/all-posts/all-posts').then(
            (m) => m.AllPosts
          ),
      },
      {
        path: 'create',
        loadComponent: () =>
          import('./pages/post-management/create/create-post').then(
            (m) => m.CreatePostComponent
          ),
      },
      {
        path: 'categories',
        loadComponent: () =>
          import(
            '../app/pages/post-management/categories/categories/categories'
          ).then((m) => m.Categories),
      },
      {
        path: 'reviews',
        loadComponent: () =>
          import('../app/pages/post-management/reviews/reviews/reviews').then(
            (m) => m.Reviews
          ),
      },
      {
        path: 'settings',
        loadComponent: () =>
          import(
            '../app/pages/post-management/settings/settings/settings'
          ).then((m) => m.Settings),
      },
      { path: '', redirectTo: 'all', pathMatch: 'full' },
    ],
  },
  {
    path: 'feature/consultation',
    component: ConsultationComponent,
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'discover', pathMatch: 'full' },
      {
        path: 'create-profile',
        loadComponent: () =>
          import(
            './pages/consultation/create-profile/create-profile.component'
          ).then((m) => m.CreateConsultationProfileComponent),
      },
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
        path: 'schedules',
        loadComponent: () =>
          import('./pages/consultation/schedules/scheduler').then(
            (m) => m.SchedulerComponent
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
          import('./pages/consultation/consultation/consultation').then(
            (m) => m.ProfileFormClientComponent
          ),
      },
    ],
  },
  { path: 'feature/chat', component: ChatComponent, canActivate: [authGuard] },
  {
    path: 'feature/social',
    component: SocialProfile,
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'wall', pathMatch: 'full' },
      { path: 'wall', component: ProfilePageComponent },
      { path: 'friends', component: FriendsPage },
    ],
  },

  { path: 'about-us', component: AboutUsComponent },
  {
    path: 'payment',
    component: PaymentManagementComponent,
    canActivate: [authGuard],
    children: [
      {
        path: 'dashboard',
        loadComponent: () =>
          import('./pages/payment-management/dashboard/payment-dashboard').then(
            (m) => m.DashboardComponent
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

  {
    path: 'schedule',
    redirectTo: 'feature/schedules',
    pathMatch: 'full',
  },
  {
    path: 'consultation',
    redirectTo: 'feature/consultation',
    pathMatch: 'prefix',
  },

  { path: 'start', component: Start },
  { path: '**', component: Notfound },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
