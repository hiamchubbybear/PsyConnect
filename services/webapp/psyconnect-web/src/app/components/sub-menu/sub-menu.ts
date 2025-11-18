// src/app/submenu/sub-menu.component.ts
import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';

export type SubItem = {
  id: string;
  label: string;
  icon?: string; // can be emoji or svg/icon class
  route?: string;
};

@Component({
  selector: 'app-sub-menu',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './sub-menu.html',
  styleUrls: ['./sub-menu.scss'],
})
export class SubMenuComponent {
  @Input() items: SubItem[] = [];

  @Input() activeId?: string;

  @Output() select = new EventEmitter<SubItem>();

  onSelect(item: SubItem) {
    this.select.emit(item);
    this.activeId = item.id;
  }
}
