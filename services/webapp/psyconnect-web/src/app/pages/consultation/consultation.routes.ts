import { Routes } from '@angular/router';
import { ConsultationComponent } from './consultation';
import { Discover } from './discover/discover';
import { Reviews } from './reviews/reviews';
import { SchedulerComponent } from './schedules/scheduler';
import { Sessions } from './sessions/sessions';
import { Settings } from './settings/settings';

export const consultationRoutes: Routes = [
  {
    path: '',
    component: ConsultationComponent,
    children: [
      { path: 'discover', component: Discover },
      { path: 'sessions', component: Sessions },
      { path: 'schedule', component: SchedulerComponent },
      {
        path: 'therapist-profile',
        loadComponent: () =>
          import('./therapist-profile/therapist-profile').then(
            (m) => m.TherapistProfileComponent
          ),
      },
      { path: 'history', component: History },
      { path: 'reviews', component: Reviews },
      { path: 'settings', component: Settings },
      { path: '', redirectTo: 'discover', pathMatch: 'full' },
    ],
  },
];
