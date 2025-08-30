import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'button-group',
  templateUrl: './button-group.html',
  styleUrls: ['./button-group.scss'],
})
export class ButtonGroupComponent {
  @Input() primaryText = '';
  @Input() secondaryText = '';

  @Output() primaryClick = new EventEmitter<void>();
  @Output() secondaryClick = new EventEmitter<void>();

  onPrimary() {
    this.primaryClick.emit();
  }

  onSecondary() {
    this.secondaryClick.emit();
  }
}
