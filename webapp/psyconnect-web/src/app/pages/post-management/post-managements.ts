import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import {
    CollapsibleSidebarComponent,
    SidebarItem,
} from '../../components/collapsible-sidebar/collapsible-sidebar';

@Component({
  selector: 'app-post-management',
  standalone: true,
  imports: [CollapsibleSidebarComponent, CommonModule, RouterOutlet],
  template: `
    <app-collapsible-sidebar
      [title]="'Post Management'"
      [titleTranslateKey]="'POST.Sidebar.Title'"
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
        padding: 2rem;
        transition: margin-left 0.3s ease;
      }
    `,
  ],
})
export class PostManagementComponent {
  isSidebarCollapsed = false;

  sidebarItems: SidebarItem[] = [
    {
      label: 'All Posts',
      route: 'all',
      translateKey: 'POST.Sidebar.Items.AllPosts',
      exact: true,
    },
    {
      label: 'Create',
      route: 'create',
      translateKey: 'POST.Sidebar.Items.Create',
    },
    {
      label: 'Categories',
      route: 'categories',
      translateKey: 'POST.Sidebar.Items.Categories',
    },
    {
      label: 'Reviews',
      route: 'reviews',
      translateKey: 'POST.Sidebar.Items.Reviews',
    },
    {
      label: 'Settings',
      route: 'settings',
      translateKey: 'POST.Sidebar.Items.Settings',
    },
  ];

  onSidebarItemClick(item: SidebarItem) {
    console.log('Sidebar item clicked:', item);
  }
}
