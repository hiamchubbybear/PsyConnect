import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterModule, RouterOutlet } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { AuthService } from '../../services/auth/auth.service';

@Component({
  selector: 'app-consultation',
  standalone: true,
  imports: [CommonModule, RouterModule, RouterOutlet, TranslateModule],
  templateUrl: './consultation.html',
  styleUrls: ['./consultation.scss'],
})
export class ConsultationComponent {
  navItems: any[] = [];

  constructor(private authService: AuthService) {
    this.setupNavigation();
  }

  private setupNavigation() {
    const role = this.authService.getRole();

    if (role === 'therapist') {
      this.navItems = [
        {
          label: 'Calendar',
          route: '/feature/consultation/schedules',
          translateKey: 'CONSULTATION.Tabs.Schedule',
        },
        {
          label: 'My Sessions',
          route: '/feature/consultation/sessions',
          translateKey: 'CONSULTATION.Tabs.Sessions',
        },
        {
          label: 'Overview',
          route: '/feature/consultation/discover',
          translateKey: 'CONSULTATION.Tabs.Overview',
        },
        {
          label: 'My Profile',
          route: '/feature/consultation/my-profile',
          translateKey: 'CONSULTATION.Tabs.MyProfile',
        },
      ];
    } else {
      this.navItems = [
        {
          label: 'Calendar',
          route: '/feature/consultation/schedules',
          translateKey: 'CONSULTATION.Tabs.Schedule',
        },
        {
          label: 'My Sessions',
          route: '/feature/consultation/sessions',
          translateKey: 'CONSULTATION.Tabs.Sessions',
        },
        {
          label: 'Discover',
          route: '/feature/consultation/discover',
          translateKey: 'CONSULTATION.Tabs.Overview',
        },
        {
          label: 'Consultators',
          route: '/feature/consultation/smart-match',
          translateKey: 'CONSULTATION.Tabs.Consultators',
        },
        {
          label: 'My Profile',
          route: '/feature/consultation/my-profile',
          translateKey: 'CONSULTATION.Tabs.MyProfile',
        },
      ];
    }
  }
}
