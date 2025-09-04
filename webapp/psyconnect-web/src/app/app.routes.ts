import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from './guards/auth-guard';
import { guestGuard } from './guards/guest-guard';
import { PasswordResetGuard } from './guards/password-reset-guard';
import { RequestResetSuccessGuard } from './guards/request-reset-guard';
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
import { RequestResetComponent } from './pages/request-reset/request-reset';
import { RequestResetSuccessComponent } from './pages/request-reset/success/request-reset-success';
import { MultiStepRegisterComponent } from './pages/signup/signup';

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
  },
  { path: 'feature/chat', component: Notfound, canActivate: [authGuard] },
  { path: 'feature/profile', component: Notfound, canActivate: [authGuard] },

  { path: '**', component: Notfound },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
