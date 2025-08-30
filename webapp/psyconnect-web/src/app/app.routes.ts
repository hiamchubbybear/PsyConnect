import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { authGuard } from './guards/auth-guard';
import { guestGuard } from './guards/guest-guard';
import { UpdateAccountComponent } from './pages/account/account-update';
import { PaymentSessionComponent } from './pages/account/payment/payment';
import { ProfileSectionComponent } from './pages/account/profile/profile-update';
import { SecuritySessionComponent } from './pages/account/security/security';
import { ChatComponent } from './pages/chat/chatpage/chatpage';
import { Homepage } from './pages/homepage/homepage';
import { Login } from './pages/login/login';
import { Notfound } from './pages/notfound/notfound';
import { MultiStepRegisterComponent } from './pages/signup/signup';
import { SwipeDeckComponent } from './pages/swipe/swipe-deck';

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
    path: 'account',
    component: UpdateAccountComponent,
    canActivate: [authGuard],
    children: [
      { path: '', redirectTo: 'profile', pathMatch: 'full' },
      { path: 'profile', component: ProfileSectionComponent },
      { path: 'security', component: SecuritySessionComponent },
      { path: 'payment', component: PaymentSessionComponent },
    ],
  },
  { path: 'feature/feed', component: Notfound, canActivate: [authGuard] },
  {
    path: 'feature/schedule',
    component: SwipeDeckComponent,
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
