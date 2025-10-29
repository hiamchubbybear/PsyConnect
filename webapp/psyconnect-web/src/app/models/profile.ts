export interface ProfileResponse {
  code: number;
  message: string;
  data: {
    accountId: string;
    profileId: string;
    firstName: string;
    lastName: string;
    dob: string;
    address: string;
    gender: string;
    avatarUri: string;
  };
}
export interface ProfileUpdateResponse {
  code: number;
  message: string;
  data: {
    username: string;
    profileId: string;
    firstName: string;
    lastName: string;
    dob: string;
    address: string;
    gender: string;
    avatarUri: string;
  };
}
export interface UserProfileUpdateRequest {
  username: string;
  firstName: string;
  lastName: string;
  dob: string;
  address: string;
  gender: string;
  avatarUri: string;
}
export interface UserProfile {
  id: string;
  name: string;
  role: string;
  email: string;
  phone: string;
  city: string;
  avatarUrl?: string;
  coverUrl?: string;
  accountStatus: string;
  twoFactorAuth: string;
  userType: string;
  propertyAccess: string;
}
