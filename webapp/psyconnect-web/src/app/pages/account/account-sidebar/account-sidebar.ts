import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'account-sidebar',
  templateUrl: './account-sidebar.html',
  styleUrls: ['./account-sidebar.scss'],
  standalone: true,
  imports: [CommonModule, RouterModule],
})
export class AppAccountSidebar {
  @Input() activeSection: string | null = null;
  @Output() navigate = new EventEmitter<string>();

  onClick(sectionId: string) {
    this.navigate.emit(sectionId);
  }
  onScroll(event: Event) {
    const target = event.target as HTMLElement;
    const sections = target.querySelectorAll('.settings-section');
    let current = '';

    sections.forEach((section: Element) => {
      const rect = section.getBoundingClientRect();
      if (rect.top <= 100 && rect.bottom >= 100) {
        current = section.id;
      }
    });

    if (current) {
      this.activeSection = current;
    }
  }
}
