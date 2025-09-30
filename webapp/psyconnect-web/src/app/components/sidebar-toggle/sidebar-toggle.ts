import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'sidebar-toggle',
  standalone: true,
  imports: [CommonModule],
  template: `
    <button
      class="sidebar-toggle"
      [class.expanded]="!collapsed"
      (click)="toggle.emit()"
      [attr.aria-label]="collapsed ? 'Mở sidebar' : 'Đóng sidebar'"
    >
      <span class="toggle-icon">{{ collapsed ? '▲' : '▼' }}</span>
    </button>
  `,
  styles: [
    `
      .sidebar-toggle {
        position: fixed;
        bottom: 0;
        left: 0;
        width: 240px;
        height: 24px;
        background: linear-gradient(to top, var(--color-border), transparent);
        border: none;
        border-top: 1px solid var(--color-border);
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 998;
        transition: all 0.3s ease;
        opacity: 0.6;

        &:hover {
          opacity: 1;
          height: 28px;
          background: linear-gradient(to top, var(--color-bg), transparent);
        }

        .toggle-icon {
          font-size: 10px;
          color: var(--color-text);
          transition: transform 0.3s ease;
        }

        &.expanded .toggle-icon {
          transform: rotate(180deg);
        }
      }
    `,
  ],
})
export class SidebarToggle {
  @Input() collapsed = false;
  @Output() toggle = new EventEmitter<void>();
}
