import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { ChatComponent } from './pages/chat/chatpage/chatpage';
import { Homepage } from './pages/homepage/homepage';
import { Login } from './pages/login/login';
import { Notfound } from './pages/notfound/notfound';
import { MultiStepRegisterComponent } from './pages/signup/signup';

export const routes: Routes = [
  {
    path: '',
    component: Homepage,
  },
  {
    path: 'auth/signup',
    component: MultiStepRegisterComponent,
  },
  {
    path: 'auth/login',
    component: Login,
  },
  { path: 'oauth2/callback', component: Login },
  {
    path: 'chat-page',
    component: ChatComponent,
  },
  {
    path: '**',
    component: Notfound,
  },
  {
    path: 'feature/feed',
    component: Notfound,
  },
  {
    path: 'feature/schedules',
    component: Notfound,
  },
  {
    path: 'feature/chat',
    component: Notfound,
  },
  {
    path: 'feature/profile',
    component: Notfound,
  },
];
@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
