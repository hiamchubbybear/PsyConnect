import { CommonModule } from '@angular/common';
import { Component, ElementRef, HostListener, ViewChild } from '@angular/core';
import {
  CollapsibleSidebarComponent,
  SidebarItem,
} from '../../components/collapsible-sidebar/collapsible-sidebar';
import { ProfileSectionComponent } from './profile/profile-update';
import { SecuritySectionComponent } from './security/security-update';

@Component({
  selector: 'app-update-account',
  standalone: true,
  imports: [
    CommonModule,
    CollapsibleSidebarComponent,
    ProfileSectionComponent,
    SecuritySectionComponent,
  ],
  template: `
    <app-collapsible-sidebar
      [title]="'Account Management'"
      [titleTranslateKey]="'ACCOUNT_MANAGEMENT.Sidebar.Title'"
      [items]="sidebarItems"
      [width]="'240px'"
      [(collapsed)]="isSidebarCollapsed"
      (itemClick)="scrollToSection($event.route)"
    ></app-collapsible-sidebar>

    <div
      class="main-content"
      #contentRef
      [style.margin-left]="isSidebarCollapsed ? '-30px' : '240px'"
    >
      <section id="profile">
        <app-profile-section></app-profile-section>
      </section>
      <section id="security">
        <app-security-section></app-security-section>
      </section>
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
export class UpdateAccountComponent {
  @ViewChild('contentRef') contentRef!: ElementRef;
  isSidebarCollapsed = false;
  activeSection: string | null = null;

  sidebarItems: SidebarItem[] = [
    {
      label: 'Account Info',
      route: 'profile',
      translateKey: 'ACCOUNT_MANAGEMENT.Sidebar.Items.AccountInfo',
    },
    {
      label: 'Security',
      route: 'security',
      translateKey: 'ACCOUNT_MANAGEMENT.Sidebar.Items.Security',
    },
    {
      label: 'Payment',
      route: 'payment',
      translateKey: 'ACCOUNT_MANAGEMENT.Sidebar.Items.Payment',
    },
    {
      label: 'Consultation',
      route: 'consultation',
      translateKey: 'ACCOUNT_MANAGEMENT.Sidebar.Items.Consultation',
    },
  ];

  scrollToSection(sectionId: string) {
    if (!this.contentRef) return;
    const el = this.contentRef.nativeElement.querySelector('#' + sectionId);
    if (el) el.scrollIntoView({ behavior: 'smooth', block: 'start' });
  }

  @HostListener('window:scroll', [])
  onScroll() {
    if (!this.contentRef) return;
    const sections = this.contentRef.nativeElement.querySelectorAll('section');
    let current: string | null = null;
    sections.forEach((section: HTMLElement) => {
      const rect = section.getBoundingClientRect();
      if (rect.top <= 120 && rect.bottom >= 120) {
        current = section.id;
      }
    });
  }
}
