import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-profile-field-dropdown',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './field-row.html',
  styleUrls: ['./field-row.scss'],
})
export class ProfileFieldDropdownComponent {
  @Input() label = '';
  @Input() displayValue = '';
  @Input() field = '';
  @Input() editable = true;

  @Input() copyable = false;

  @Output() editRequested = new EventEmitter<string>();
  @Output() copyRequested = new EventEmitter<string>();

  isExpanded = false;
  isEditing = false;

  toggleExpanded() {
    this.isExpanded = !this.isExpanded;
  }

  onCopyClick(event: Event) {
    event.stopPropagation();
    const valueToCopy = this.displayValue || '';
    navigator.clipboard
      .writeText(valueToCopy)
      .then(() => console.log('Copied value:', valueToCopy))
      .catch((err) => console.error('Copy failed:', err));

    this.copyRequested.emit(valueToCopy);
  }
  handleHeaderClick(event: Event) {
    event.stopPropagation();

    if (this.editable) {
      this.toggleExpanded();
    } else if (this.copyable) {
      const valueToCopy = this.displayValue || '';
      if (!valueToCopy) return;

      navigator.clipboard
        .writeText(valueToCopy)
        .then(() => console.log('Copied value:', valueToCopy))
        .catch((err) => console.error('Copy failed:', err));

      this.copyRequested.emit(valueToCopy);
    }
  }

  onEditAreaClick(event: Event) {
    if (!this.editable) return;
    event.stopPropagation();
    this.isEditing = !this.isEditing;
    this.editRequested.emit(this.field);
  }
}
