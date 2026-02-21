import { CommonModule } from '@angular/common';
import { Component, ElementRef, Input, ViewChild } from '@angular/core';
import { LucideAngularModule } from 'lucide-angular';

@Component({
  selector: 'app-horizontal-list',
  standalone: true,
  imports: [CommonModule, LucideAngularModule],
  templateUrl: './horizontal-list.component.html',
  styleUrls: ['./horizontal-list.component.scss'],
})
export class HorizontalListComponent {
  @Input() title = '';
  @Input() isEmpty = false;
  @Input() emptyText = '';
  @Input() showNav = true;

  @ViewChild('scrollContainer') scrollContainer!: ElementRef<HTMLDivElement>;

  scrollLeft() {
    if (this.scrollContainer) {
      this.scrollContainer.nativeElement.scrollBy({
        left: -300,
        behavior: 'smooth',
      });
    }
  }

  scrollRight() {
    if (this.scrollContainer) {
      this.scrollContainer.nativeElement.scrollBy({
        left: 300,
        behavior: 'smooth',
      });
    }
  }
}
