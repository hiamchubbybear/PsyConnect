import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Output } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';

export type CallType = 'audio' | 'video' | 'consultation';

@Component({
  selector: 'app-call-options-menu',
  standalone: true,
  imports: [CommonModule, TranslateModule],
  templateUrl: './call-options-menu.component.html',
  styleUrls: ['./call-options-menu.component.scss'],
})
export class CallOptionsMenuComponent {
  @Output() callTypeSelected = new EventEmitter<CallType>();

  showMenu = false;

  toggleMenu() {
    this.showMenu = !this.showMenu;
  }

  selectCallType(type: CallType) {
    this.callTypeSelected.emit(type);
    this.showMenu = false;
  }

  
  closeMenu() {
    this.showMenu = false;
  }
}
