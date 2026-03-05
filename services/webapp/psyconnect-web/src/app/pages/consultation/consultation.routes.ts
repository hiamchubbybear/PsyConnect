import { Routes } from '@angular/router';
import { ConsultationComponent } from './consultation';
import { Discover } from './discover/discover';
import { Reviews } from './reviews/reviews';
import { SchedulerComponent } from './schedules/scheduler';
import { Sessions } from './sessions/sessions';
import { Settings } from './settings/settings';
import { SmartMatchComponent } from './smart-match/smart-match';

export const consultationRoutes: Routes = [
  {
    path: '',
    component: ConsultationComponent,
    children: [
      { path: 'smart-match', component: SmartMatchComponent },
      { path: 'discover', component: Discover },
      { path: 'sessions', component: Sessions },
      { path: 'schedule', component: SchedulerComponent },
      {
        path: 'therapist-profile',
        loadComponent: () =>
          import('./therapist-profile/therapist-profile').then(
            (m) => m.TherapistProfileComponent,
          ),
      },
      {
        path: 'my-profile',
        loadComponent: () =>
          import('./my-profile/my-profile').then(
            (m) => m.MyConsultationProfileComponent,
          ),
      },
      { path: 'history', component: History },
      { path: 'reviews', component: Reviews },
      {
        path: 'book/:therapistId',
        loadComponent: () =>
          import('./booking-page/booking-page').then(
            (m) => m.BookingPageComponent,
          ),
      },
      { path: 'settings', component: Settings },
      { path: '', redirectTo: 'discover', pathMatch: 'full' },
    ],
  },
];
