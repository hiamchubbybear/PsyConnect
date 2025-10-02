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
