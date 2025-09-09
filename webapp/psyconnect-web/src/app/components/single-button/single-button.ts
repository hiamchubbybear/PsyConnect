import { Component, EventEmitter, Input, Output } from '@angular/core';

@Component({
  selector: 'app-single-button',
  imports: [],
  templateUrl: './single-button.html',
  styleUrl: './single-button.scss',
})
export class SingleButton {
  @Input() primaryText = '';
  @Input() secondaryText = '';
  @Input() disabledPrimary = false;
  @Output() primaryClick = new EventEmitter<void>();
  @Output() secondaryClick = new EventEmitter<void>();

  onPrimary() {
    this.primaryClick.emit();
  }
}
