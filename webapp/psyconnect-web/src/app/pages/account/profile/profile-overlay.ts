import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { TranslateModule } from '@ngx-translate/core';

@Component({
  selector: 'app-profile-overlay',
  standalone: true,
  imports: [CommonModule, FormsModule, TranslateModule],
  templateUrl: './profile-overlay.html',
  styleUrls: ['./profile-overlay.scss'],
})
export class ProfileOverlayComponent {
  @Input() editing: string | null = null;
  @Input() draftData: {
    firstName: string;
    middleName: string;
    lastName: string;
  } = {
    firstName: '',
    middleName: '',
    lastName: '',
  };
  @Input() editValue: any = '';

  @Output() saveName = new EventEmitter<void>();
  @Output() saveField = new EventEmitter<string>();
  @Output() cancel = new EventEmitter<void>();

  onSaveField(field: string) {
    this.saveField.emit(field);
  }
}
