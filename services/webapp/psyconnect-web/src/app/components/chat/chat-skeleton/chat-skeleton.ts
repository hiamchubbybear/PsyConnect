import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';

@Component({
  selector: 'app-chat-skeleton',
  standalone: true,
  imports: [CommonModule],
  template: `
    <div class="chat-skeleton">
      <div class="skeleton-item" *ngFor="let i of skeletonItems">
        <div class="skeleton-avatar"></div>
        <div class="skeleton-content">
          <div class="skeleton-line skeleton-name"></div>
          <div class="skeleton-line skeleton-message"></div>
        </div>
      </div>
    </div>
  `,
  styleUrls: ['./chat-skeleton.scss'],
})
export class ChatSkeletonComponent {
  @Input() count = 5;

  get skeletonItems() {
    return Array(this.count).fill(0);
  }
}
