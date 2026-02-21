import { CommonModule } from '@angular/common';
import { Component, Input, ViewEncapsulation } from '@angular/core';
import { PsyEmptyStateComponent } from '../empty-state/psy-empty-state.component';

@Component({
  selector: 'app-table',
  standalone: true,
  imports: [CommonModule, PsyEmptyStateComponent],
  templateUrl: './table.component.html',
  styleUrls: ['./table.component.scss'],
  encapsulation: ViewEncapsulation.None, // To apply styles to projected content
})
export class TableComponent {
  @Input() loading = false;
  @Input() isEmpty = false;
  @Input() emptyTitle = '';
  @Input() emptyDescription = '';
  @Input() emptyIcon = '';
}
