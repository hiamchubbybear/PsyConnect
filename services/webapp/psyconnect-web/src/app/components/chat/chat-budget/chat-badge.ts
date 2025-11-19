import { CommonModule } from '@angular/common';
import { Component, input } from '@angular/core';

type Status = 'available' | 'busy' | 'offline' | 'away';

@Component({
  selector: 'app-status-badge',
  standalone: true,
  imports: [CommonModule],
  templateUrl: `./status-badge.html`,
  styleUrl: `./status-badge.scss`,
})
export class StatusBadgeComponent {
  status = input<Status>('available');
  showDropdown = input(true);
}
