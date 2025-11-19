export interface Connection {
  id: string;
  name: string;
  avatarUrl: string;
  mutualCount?: number;
  mutualAvatars?: string[];
}

export interface FriendRequest extends Connection {
  requestedAt?: Date;
}

export interface FriendSuggestion extends Connection {
  mutualCount: number;
  mutualAvatars: string[];
}
