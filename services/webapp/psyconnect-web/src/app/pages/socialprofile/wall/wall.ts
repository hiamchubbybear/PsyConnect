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
    private userContext: UserContextService
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
        this.userProfile = res.data || res;
        this.mapFields();
      },
    });
  }

  loadUserProfile(id: string) {
    this.profileService.getProfileById(id).subscribe({
      next: (res: any) => {
        this.userProfile = res.data || res;
        this.mapFields();
      },
    });
  }

  mapFields() {
    if (!this.userProfile) return;

    this.staffInfoFields = [
      { label: 'Email', value: this.userProfile.email || 'N/A' },
      { label: 'Phone', value: this.userProfile.phone || 'N/A' },
      { label: 'Address', value: this.userProfile.address || 'N/A' },
    ];

    this.websiteFields = [
      { label: 'Gender', value: this.userProfile.gender || 'N/A' },
      { label: 'DOB', value: this.userProfile.dob || 'N/A' },
      { label: 'Bio', value: this.userProfile.bio || 'N/A' },
    ];
  }

  onEditStaffInfo(): void {}
  onEditWebsite(): void {}
}
