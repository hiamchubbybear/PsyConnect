
import { CommonModule } from '@angular/common';
import { Component, Input } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import { UserProfile } from '../../../models/profile';
import { AvatarFallbackPipe } from '../../../shared/pipes/avatar-fallback.pipe';

@Component({
  selector: 'app-profile-card',
  standalone: true,
  imports: [CommonModule, TranslateModule, AvatarFallbackPipe],
  templateUrl: './profile-card.html',
  styleUrls: ['./profile-card.scss'],
})
export class ProfileCardComponent {
  @Input() user!: UserProfile;
  @Input() isOwnProfile: boolean = false;

  getInitials(name: string): string {
    return name
      .split(' ')
      .map((n) => n[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  }
}
