import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {
  CollapsibleSidebarComponent,
  SidebarItem,
} from '../../components/collapsible-sidebar/collapsible-sidebar';

@Component({
  selector: 'app-socialprofile',
  imports: [RouterOutlet, CollapsibleSidebarComponent],
  template: `<app-collapsible-sidebar
      [title]="'Social'"
      [titleTranslateKey]="'SOCIAL.sidebar.title'"
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
    </div>`,
  styles: `
    .main-content {
      margin-left: 360px;
      padding: 2rem;
      transition: margin-left 0.3s ease;
    }
  `,
})
export class SocialProfile {
  isSidebarCollapsed = false;

  sidebarItems: SidebarItem[] = [
    {
      label: 'Wall',
      route: 'wall',
      translateKey: 'SOCIAL.sidebar.items.wall',
      exact: true,
    },
    {
      label: 'Friends',
      route: 'friends',
      translateKey: 'SOCIAL.sidebar.items.friends',
    },
    {
      label: 'Discover',
      route: 'discover',
      translateKey: 'SOCIAL.sidebar.items.discover',
    },
    {
      label: 'Settings',
      route: 'settings',
      translateKey: 'SOCIAL.sidebar.items.settings',
    },
  ];

  onSidebarItemClick(item: SidebarItem) {
    console.log('Sidebar item clicked:', item);
  }
}
