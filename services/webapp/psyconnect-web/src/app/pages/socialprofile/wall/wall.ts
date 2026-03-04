import { CommonModule } from '@angular/common';
import { Component, OnInit } from '@angular/core';
import { ActivatedRoute } from '@angular/router';
import { TranslateModule, TranslateService } from '@ngx-translate/core';
import {
  InfoBlockComponent,
  InfoField,
} from '../../../components/social/info-block/info-block';
import { ProfileCardComponent } from '../../../components/social/profile-card/profile-card';
import { Post } from '../../../models/post.model';
import { NewsfeedService } from '../../../services/newsfeed/newsfeed.service';
import { Profile } from '../../../services/profile/profile';
import { UserContextService } from '../../../services/profile/profile-service';
import { PostCardComponent } from '../../../shared/ui-atoms/post-card/post-card.component';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [
    CommonModule,
    ProfileCardComponent,
    InfoBlockComponent,
    PostCardComponent,
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
  userPosts: Post[] = [];
  isLoadingPosts = false;

  constructor(
    private route: ActivatedRoute,
    private profileService: Profile,
    private userContext: UserContextService,
    private newsfeedService: NewsfeedService,
    private translate: TranslateService,
  ) {}

  ngOnInit(): void {
    this.route.params.subscribe((params) => {
      const userId = params['id'];
      const currentUser = this.userContext.getUser();

      // If no ID or ID is current user, load own profile
      if (
        !userId ||
        userId === currentUser?.accountId ||
        userId === currentUser?.profileId ||
        userId === 'me'
      ) {
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
        const currentUser = this.userContext.getUser();
        const targetId = raw.accountId || currentUser?.accountId;
        if (targetId) {
          this.loadUserPosts(targetId);
        }
      },
    });
  }

  loadUserProfile(id: string) {
    this.profileService.getProfileById(id).subscribe({
      next: (res: any) => {
        const raw = res.data || res;
        this.userProfile = this.transformToCardFormat(raw);
        this.mapFields(raw);
        const targetId = raw.accountId;
        if (targetId) {
          this.loadUserPosts(targetId);
        }
      },
    });
  }

  loadUserPosts(userId: string) {
    this.isLoadingPosts = true;
    this.newsfeedService.getUserPosts(userId).subscribe({
      next: (posts) => {
        this.userPosts = posts || [];
        this.isLoadingPosts = false;
      },
      error: () => {
        this.isLoadingPosts = false;
        this.userPosts = [];
      },
    });
  }

  /**
   * Transforms ProfileResponse.data shape into the shape expected by profile-card.html
   * (user.name, user.avatarUrl, user.role, user.coverUrl)
   */
  private transformToCardFormat(profile: any): any {
    return {
      avatarUrl: profile.avatarUri || 'assets/images/default-avatar.png',
      name:
        `${profile.firstName || ''} ${profile.lastName || ''}`.trim() ||
        'Chưa cập nhật',
      username: profile.username || '',
      joinDate: profile.joinDate || new Date().toISOString(),
      followers: profile.followersCount || 0,
      following: profile.followingCount || 0,
      bio: profile.bio || profile.description || 'Chưa có thông tin',
    };
  }

  mapFields(raw: any) {
    if (!raw) return;

    this.staffInfoFields = [
      {
        label: 'SOCIAL.wall.fields.email',
        value: raw.email || 'N/A',
      },
      {
        label: 'SOCIAL.wall.fields.phone',
        value: raw.phone || 'N/A',
      },
      {
        label: 'SOCIAL.wall.fields.address',
        value: raw.address || 'N/A',
      },
    ];

    this.websiteFields = [
      {
        label: 'SOCIAL.wall.fields.gender',
        value: raw.gender || 'N/A',
      },
      {
        label: 'SOCIAL.wall.fields.dob',
        value: raw.dob || 'N/A',
      },
      {
        label: 'SOCIAL.wall.fields.bio',
        value: raw.bio || raw.description || 'N/A',
      },
    ];
  }

  onEditStaffInfo(): void {}
  onEditWebsite(): void {}
}
