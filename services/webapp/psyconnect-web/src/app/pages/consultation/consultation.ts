import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-consultation',
  standalone: true,
  imports: [CommonModule, RouterModule, RouterOutlet, TranslateModule],
  templateUrl: './consultation.html',
  styleUrls: ['./consultation.scss'],
})
export class ConsultationComponent {
  navItems = [
    {
      label: 'Overview',
      route: '/feature/consultation/discover',
      translateKey: 'CONSULTATION.Tabs.Overview',
    },
    {
      label: 'My Sessions',
      route: '/feature/consultation/sessions',
      translateKey: 'CONSULTATION.Tabs.Sessions',
    },
    {
      label: 'Calendar',
      route: '/feature/consultation/schedules',
      translateKey: 'CONSULTATION.Tabs.Schedule',
    },
    {
      label: 'Consultators',
      route: '/feature/consultation/smart-match',
      translateKey: 'CONSULTATION.Tabs.Consultators',
    },
  ];
}
