import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {
    CollapsibleSidebarComponent,
    SidebarItem,
} from '../../components/collapsible-sidebar/collapsible-sidebar';

@Component({
  selector: 'app-consultation',
  standalone: true,
  imports: [CollapsibleSidebarComponent, CommonModule, RouterOutlet],
  template: `
    <app-collapsible-sidebar
      [title]="'Consultation'"
      [titleTranslateKey]="'CONSULTATION.Sidebar.Title'"
      [items]="sidebarItems"
      [width]="'240px'"
      [(collapsed)]="isSidebarCollapsed"
      (itemClick)="onSidebarItemClick($event)"
    />

    <div
      class="main-content"
      [style.margin-left]="isSidebarCollapsed ? '-30px' : '240px'"
    >
      <router-outlet></router-outlet>
    </div>
  `,
  styles: [
    `
      .main-content {
        margin-left: 360px;
        padding: 2rem;
        transition: margin-left 0.3s ease;
      }
    `,
  ],
})
export class ConsultationComponent {
  isSidebarCollapsed = false;

  sidebarItems: SidebarItem[] = [
    {
      label: 'Discover',
      route: 'discover',
      translateKey: 'CONSULTATION.Sidebar.Items.Discover',

      exact: true,
    },
    {
      label: 'Sessions',
      route: 'sessions',
      translateKey: 'CONSULTATION.Sidebar.Items.Sessions',
    },
    {
      label: 'Schedule',
      route: 'schedule',
      translateKey: 'CONSULTATION.Sidebar.Items.Schedule',
    },
    {
      label: 'History',
      route: 'history',
      translateKey: 'CONSULTATION.Sidebar.Items.History',
    },
    {
      label: 'Reviews',
      route: 'reviews',
      translateKey: 'CONSULTATION.Sidebar.Items.Reviews',
    },
    {
      label: 'Settings',
      route: 'settings',
      translateKey: 'CONSULTATION.Sidebar.Items.Settings',
    },
  ];

  onSidebarItemClick(item: SidebarItem) {
    console.log('Sidebar item clicked:', item);
  }
}
