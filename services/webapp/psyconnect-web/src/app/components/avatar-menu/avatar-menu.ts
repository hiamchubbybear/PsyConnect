import { CommonModule } from '@angular/common';
import {
    Component,
    ElementRef,
    EventEmitter,
    HostListener,
    Input,
    Output,
} from '@angular/core';
import { RouterModule } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import { SecureStorageService } from '../../encrypt/secure';
import { ThemeService } from '../../services/theme/theme-service';

@Component({
  selector: 'app-avatar-menu',
  standalone: true,
  imports: [CommonModule, RouterModule, TranslateModule],
  templateUrl: './avatar-menu.html',
  styleUrl: './avatar-menu.scss',
})
export class AvatarMenuComponent {
  isOpen = false;
  constructor(
    private themeService: ThemeService,
    private secureStorage: SecureStorageService,
    private elementRef: ElementRef
  ) {}
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
    this.secureStorage.clear();
  }
  @HostListener('document:click', ['$event'])
  onClickOutside(event: Event) {
    const target = event.target as HTMLElement;
    const hostElement = event.currentTarget as Document;

    if (this.isOpen && !this.elementRef.nativeElement.contains(target)) {
      this.isOpen = false;
    }
  }

  @HostListener('document:keydown', ['$event'])
  onEsc(event: KeyboardEvent) {
    if (event.key === 'Escape' && this.isOpen) {
      this.isOpen = false;
    }
  }
}
