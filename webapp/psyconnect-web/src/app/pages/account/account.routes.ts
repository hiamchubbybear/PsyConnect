// app.routes.ts (hoặc account.routes.ts nếu bạn tách riêng)
import { Routes } from '@angular/router';
import { UpdateAccountComponent } from './account';
export const routes: Routes = [
  {
    path: 'account',
    component:UpdateAccountComponent, // account.html có sidebar + router-outlet
    children: [
      {
        path: '',
        redirectTo: 'profile', // default
        pathMatch: 'full'
      },
      {
        path: 'profile',
        component: UpdateAccountComponent
      },
      {
        path: 'security',
        component: UpdateAccountComponent
      },
      {
        path: 'payment',
        component: UpdateAccountComponent
      }
    ]
  }
];
