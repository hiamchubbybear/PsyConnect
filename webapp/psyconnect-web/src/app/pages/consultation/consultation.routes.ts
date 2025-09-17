import { Routes } from '@angular/router';
import { Consultation } from './consultation';
import { Discover } from './discover/discover';
import { Reviews } from './reviews/reviews';
import { Schedules } from './schedules/schedules';
import { Sessions } from './sessions/sessions';
import { Settings } from './settings/settings';

export const consultationRoutes: Routes = [
  {
    path: '',
    component: Consultation,
    children: [
      { path: 'discover', component: Discover },
      { path: 'sessions', component: Sessions },
      { path: 'schedule', component: Schedules },
      { path: 'history', component: History },
      { path: 'reviews', component: Reviews },
      { path: 'settings', component: Settings },
      { path: '', redirectTo: 'discover', pathMatch: 'full' },
    ],
  },
];
