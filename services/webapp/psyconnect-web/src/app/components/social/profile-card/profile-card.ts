import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

import { Router } from '@angular/router';

@Component({
  selector: 'app-profile-card',
  standalone: true,
  imports: [CommonModule, TranslateModule, AvatarFallbackPipe],
  templateUrl: './profile-card.html',
  styleUrls: ['./profile-card.scss'],
})
export class ProfileCardComponent {
  @Input() user!: any;
  @Input() isOwnProfile: boolean = false;

  constructor(private router: Router) {}

  getInitials(name: string): string {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  }

  onMessage() {
    if (this.user && (this.user.profileId || this.user.accountId)) {
      const id = this.user.accountId || this.user.profileId;
      this.router.navigate(['/feature/chat', id]);
    }
  }

  onBookSession() {
    if (this.user && (this.user.profileId || this.user.accountId)) {
      const id = this.user.profileId || this.user.accountId;
      this.router.navigate(['/feature/consultation/book', id]);
    }
  }
}
