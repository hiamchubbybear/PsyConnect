import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-profile-header',
  standalone: true,
  imports: [CommonModule ,TranslateModule],
  templateUrl: './profile-header.html',
  styleUrls: ['./profile-header.scss'],
})
export class ProfileHeaderComponent {
  @Input() name: string = '';
  @Input() role: string = '';
  @Output() edit = new EventEmitter<void>();

  onEdit(): void {
    this.edit.emit();
  }
}
