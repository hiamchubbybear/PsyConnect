import { ConsultationSession } from './consultation.model';

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
  sessionData?: ConsultationSession;
  isSystem?: boolean;
  createdAt: string;
  updateAt: string;
}
export interface Friend {
  profileId: string;
  firstName: string;
  lastName: string;
  avatarUri: string;
  role?: string; 
  isOnline?: boolean;
  hasNewMessage?: boolean;
  lastMessage?: string;
  lastMessageTime?: Date;
  isTyping?: boolean;
}
export interface Message {
  id: string;
  conversationId: string;
  userName: string;
  senderId: string;
  userAvatar: string;
  content: string;
  sessionData?: ConsultationSession;
  timestamp: Date;
  isMine: boolean;
  isSystem?: boolean; 
  
  status?: 'sending' | 'sent' | 'delivered' | 'failed';
  
  isFirstInGroup?: boolean;
  isLastInGroup?: boolean;
  showTimestamp?: boolean;
  showDateSeparator?: boolean;
  dateSeparatorText?: string;
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
