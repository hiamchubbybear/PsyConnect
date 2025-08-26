import { CommonModule } from '@angular/common';
import { Component, HostListener } from '@angular/core';

@Component({
  selector: 'account-sidebar',
  templateUrl: './account-sidebar.html',
  styleUrls: ['./account-sidebar.scss'],
  imports : [CommonModule] ,
  standalone : true
})
export class AccountSidebarComponent {

  sections = [
    { id: 'account-info', label: 'Thông tin tài khoản' },
    { id: 'personal-info', label: 'Thông tin cá nhân' },
    { id: 'security', label: 'Bảo mật' },
    { id: 'payment', label: 'Thanh toán' }
  ];

  activeSection: string = this.sections[0].id;

  scrollTo(sectionId: string) {
    const el = document.getElementById(sectionId);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
      this.activeSection = sectionId;
    }
  }
  @HostListener('window:scroll', [])
  onWindowScroll() {
    for (let sec of this.sections) {
      const el = document.getElementById(sec.id);
      if (el) {
        const rect = el.getBoundingClientRect();
        if (rect.top <= 100 && rect.bottom >= 100) {
          this.activeSection = sec.id;
          break;
        }
      }
    }
  }
}
