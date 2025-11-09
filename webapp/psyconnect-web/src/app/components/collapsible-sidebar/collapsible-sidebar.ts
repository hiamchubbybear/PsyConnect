import { CommonModule } from '@angular/common';
import { Component, EventEmitter, Input, Output } from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { SidebarService } from '../../services/sidebar/sidebar';

export interface SidebarItem {
  label: string;
  route: string;
  translateKey?: string;
  icon?: string;
  exact?: boolean;
}

@Component({
  selector: 'app-collapsible-sidebar',
  standalone: true,
  imports: [CommonModule, TranslateModule, RouterModule],
  templateUrl: './collapsible-sidebar.html',
  styleUrls: ['./collapsible-sidebar.scss'],
})
export class CollapsibleSidebarComponent {
  @Input() title: string = 'Menu';
  @Input() titleTranslateKey?: string;
  @Input() items: SidebarItem[] = [];
  @Input() width: string = '240px';
  @Input() collapsed: boolean = false;

  @Output() collapsedChange = new EventEmitter<boolean>();
  @Output() itemClick = new EventEmitter<SidebarItem>();
  constructor(private sidebarService: SidebarService) {}
  toggleSidebar() {
    this.collapsed = !this.collapsed;
    this.collapsedChange.emit(this.collapsed);
    this.sidebarService.toggleSidebar(this.collapsed);
  }

  onItemClick(item: SidebarItem) {
    this.itemClick.emit(item);
  }
}
