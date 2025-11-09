import { UserProfile } from '../services/profile/profile-service';

export interface ChatUser {
  profileId: string;
  name: string;
  avatar: string;
  isOnline: boolean;
  status: string;
}

export function mapUserProfileToChatUser(
  profile: UserProfile | null
): ChatUser {
  if (!profile) {
    return {
      profileId: 'fallback',
      name: 'Guest User',
      avatar: 'https://i.pravatar.cc/150?img=1',
      isOnline: true,
      status: 'CHAT.EXAMPLE.status.available',
    };
  }
  console.log('Profile', {
    id: profile.profileId || 'unknown',
    name:
      `${profile.firstName ?? ''} ${profile.lastName ?? ''}`.trim() ||
      'No Name',
    avatar: profile.avatarUri || 'https://i.pravatar.cc/150?img=1',
    isOnline: true,
    status: 'CHAT.EXAMPLE.status.available',
  });

  return {
    profileId: profile.profileId || 'unknown',
    name:
      `${profile.firstName ?? ''} ${profile.lastName ?? ''}`.trim() ||
      'No Name',
    avatar: profile.avatarUri || 'https://i.pravatar.cc/150?img=1',
    isOnline: true,
    status: 'CHAT.EXAMPLE.status.available',
  };
}
