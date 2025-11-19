// src/app/features/profile/profile-page.component.ts
import { CommonModule } from '@angular/common';
import { Component } from '@angular/core';
import { TranslateModule } from '@ngx-translate/core';
import {
    InfoBlockComponent,
    InfoField,
} from '../../../components/social/info-block/info-block';
import { ProfileCardComponent } from '../../../components/social/profile-card/profile-card';
import { UserProfile } from '../../../models/profile';

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
export class ProfilePageComponent {
  mockUser: UserProfile = {
    id: '1',
    name: 'Danielle Pimentel',
    role: 'Leasing Agent',
    email: 'daniellepimentel@gmail.com',
    phone: '555-55-2261',
    city: 'Los Angeles',
    accountStatus: 'Account Created',
    twoFactorAuth: 'Not Set',
    userType: 'Staff Member',
    propertyAccess: 'All Property',
    avatarUrl:
      'https://i.pinimg.com/736x/09/e4/4f/09e44f60351b96df26652dc3bf775d98.jpg',
    coverUrl:
      'https://i.pinimg.com/736x/d0/5d/3a/d05d3a3f0ff18b9002bbd2bc19684fd7.jpg',
  };
  staffInfoFields: InfoField[] = [
    { label: 'Email', value: this.mockUser.email },
    { label: 'Phone', value: this.mockUser.phone },
    { label: 'City', value: this.mockUser.city },
  ];

  websiteFields: InfoField[] = [
    { label: 'Account Status', value: this.mockUser.accountStatus },
    { label: 'Two Factor Auth', value: this.mockUser.twoFactorAuth },
    { label: 'Role', value: this.mockUser.role },
    { label: 'User Type', value: this.mockUser.userType },
    { label: 'Property Access', value: this.mockUser.propertyAccess },
  ];

  onEditProfile(): void {
    console.log('Edit profile clicked');
  }

  onEditStaffInfo(): void {
    console.log('Edit staff info clicked');
  }

  onEditWebsite(): void {
    console.log('Edit website clicked');
  }
}
