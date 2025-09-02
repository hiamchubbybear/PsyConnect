import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    HostListener,
    Input,
    ViewChild,
} from '@angular/core';
import { Discover } from './discover/discover';
import { ConsultationSidebar } from './sidebar/sidebar';

@Component({
  selector: 'app-consultation',
  imports: [ConsultationSidebar, CommonModule, Discover],
  templateUrl: './consultation.html',
  styleUrl: '../account/account-update.scss',
  standalone: true,
})
export class Consultation {
  @ViewChild('content') contentRef!: ElementRef;
  @Input() activeSection: string | null = null;

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
