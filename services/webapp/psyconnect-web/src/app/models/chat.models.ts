export interface User {
  id: string;
  name: string;
  avatar: string;
  status?: 'available' | 'busy' | 'offline';
}
export interface ChatFromApi {
  id: string;
  senderId: string;
  conversationId: string;
  text: string;
  createdAt: string;
  updateAt: string;
}
export interface Friend {
  profileId: string;
  firstName: string;
  lastName: string;
  avatarUri: string;
  isOnline?: boolean;
  hasNewMessage?: boolean;
}
export interface Message {
  id: string;
  conversationId: string;
  userName: string;
  senderId: string;
  userAvatar: string;
  content: string;
  timestamp: Date;
  isMine: boolean;
}

export interface Chat {
  id: string;
  title: string;
  avatar: string;
  lastMessage: string;
  timestamp: string;
  unread?: number;
  isTyping?: boolean;
}
