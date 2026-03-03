import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateModule } from '@ngx-translate/core';
import {
  InfoBlockComponent,
  InfoField,
} from '../../../components/social/info-block/info-block';
import { ProfileCardComponent } from '../../../components/social/profile-card/profile-card';
import { Profile } from '../../../services/profile/profile';
import { UserContextService } from '../../../services/profile/profile-service';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [
    CommonModule,
    ProfileCardComponent,
    InfoBlockComponent,
    TranslateModule,
  ],
  templateUrl: './wall.html',
  styleUrls: ['./wall.scss'],
})
export class ProfilePageComponent implements OnInit {
  userProfile: any = null;
  isOwnProfile = false;
  staffInfoFields: InfoField[] = [];
  websiteFields: InfoField[] = [];

  constructor(
    private route: ActivatedRoute,
    private profileService: Profile,
    private userContext: UserContextService,
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      const userId = params['id'];
      const currentUser = this.userContext.getUser();

      if (!userId || userId === currentUser?.accountId) {
        this.isOwnProfile = true;
        this.loadMyProfile();
      } else {
        this.isOwnProfile = false;
        this.loadUserProfile(userId);
      }
    });
  }

  loadMyProfile() {
    this.profileService.getProfile().subscribe({
      next: (res: any) => {
        const raw = res.data || res;
        this.userProfile = this.transformToCardFormat(raw);
        this.mapFields(raw);
      },
    });
  }

  loadUserProfile(id: string) {
    this.profileService.getProfileById(id).subscribe({
      next: (res: any) => {
        const raw = res.data || res;
        this.userProfile = this.transformToCardFormat(raw);
        this.mapFields(raw);
      },
    });
  }

  /**
   * Transforms ProfileResponse.data shape into the shape expected by profile-card.html
   * (user.name, user.avatarUrl, user.role, user.coverUrl)
   */
  private transformToCardFormat(raw: any): any {
    const firstName = raw.firstName || '';
    const lastName = raw.lastName || '';
    return {
      ...raw,
      name: [firstName, lastName].filter(Boolean).join(' ') || 'N/A',
      avatarUrl: raw.avatarUri || raw.avatarUrl || '',
      role: raw.role || '',
      coverUrl: raw.coverUrl || raw.coverUri || null,
    };
  }

  mapFields(raw: any) {
    if (!raw) return;

    this.staffInfoFields = [
      { label: 'Email', value: raw.email || 'N/A' },
      { label: 'Phone', value: raw.phone || 'N/A' },
      { label: 'Address', value: raw.address || 'N/A' },
    ];

    this.websiteFields = [
      { label: 'Gender', value: raw.gender || 'N/A' },
      { label: 'DOB', value: raw.dob || 'N/A' },
      { label: 'Bio', value: raw.bio || raw.description || 'N/A' },
    ];
  }

  onEditStaffInfo(): void {}
  onEditWebsite(): void {}
}
