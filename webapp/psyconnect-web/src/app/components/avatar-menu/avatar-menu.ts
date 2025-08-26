import { CommonModule } from '@angular/common';
import {
    Component,
    EventEmitter,
    HostListener,
    Input,
    Output,
} from '@angular/core';
import { RouterModule } from '@angular/router';
import { ThemeService } from '../../services/theme/theme-service';

@Component({
  selector: 'app-avatar-menu',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './avatar-menu.html',
  styleUrl: './avatar-menu.scss',
})
export class AvatarMenuComponent {
  isOpen = false;
  constructor(private themeService: ThemeService) {}
  toggleMenu() {
    this.isOpen = !this.isOpen;
  }

  @Input() userAvatarUrl: string = '';
  @Input() userName: string = '';
  @Input() userEmail: string = '';
  @Output() logoutClicked = new EventEmitter<void>();

  toggleTheme() {
    this.themeService.toggleTheme();
  }
  logout() {
    this.logoutClicked.emit();
  }
  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape' && this.isOpen) {
      this.isOpen = false;
    }
  }
}
