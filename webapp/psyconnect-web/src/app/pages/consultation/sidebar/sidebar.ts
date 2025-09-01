import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    EventEmitter,
    HostListener,
    Input,
    Output,
    ViewChild,
} from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'consultation-sidebar',
  imports: [CommonModule, TranslateModule],
  templateUrl: './sidebar.html',
  styleUrl: './sidebar.scss',
})
export class ConsultationSidebar {
  @Output() navigate = new EventEmitter<string>();
  @ViewChild('content') contentRef!: ElementRef;
  @Input() activeSection: string | null = null;

  scrollToSection(sectionId: string) {
    const el = this.contentRef.nativeElement.querySelector('#' + sectionId);
    if (el) {
      el.scrollIntoView({ behavior: 'smooth', block: 'start' });
    }
  }
  onClick(section: string) {
    this.navigate.emit(section);
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
