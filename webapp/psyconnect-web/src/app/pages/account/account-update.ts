import { CommonModule } from '@angular/common';
import { Component, ElementRef, HostListener, ViewChild } from '@angular/core';
import { AppAccountSidebar } from './account-sidebar/account-sidebar';
import { PaymentSessionComponent } from './payment/payment';
import { ProfileSectionComponent } from './profile/profile-update';
import { SecuritySectionComponent } from "./security/security-update";

@Component({
  selector: 'app-update-account',
  standalone: true,
  imports: [
    CommonModule,
    AppAccountSidebar,
    PaymentSessionComponent,
    ProfileSectionComponent,
    SecuritySectionComponent
],
  templateUrl: './account-update.html',
  styleUrls: ['./account-update.scss'],
})
export class UpdateAccountComponent {
  @ViewChild('content') contentRef!: ElementRef;
  activeSection: string | null = null;

  scrollToSection(sectionId: string) {
    const el = this.contentRef.nativeElement.querySelector('#' + sectionId);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
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

    this.activeSection = current;
  }
}
