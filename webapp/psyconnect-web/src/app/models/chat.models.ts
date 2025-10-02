export interface User {
  id: string;
  name: string;
  avatar: string;
  status?: 'available' | 'busy' | 'offline';
}

export interface Message {
  id: string;
  userId: string;
  userName: string;
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
export const fakeMessages: Message[] = [
  {
    id: '1',
    userId: 'u1',
    userName: 'Alice',
    userAvatar: 'assets/avatars/alice.png',
    content: 'Hey, how are you?',
    timestamp: new Date(Date.now() - 1000 * 60 * 5),
    isMine: false,
  },
  {
    id: '2',
    userId: 'me',
    userName: 'You',
    userAvatar: 'assets/avatars/me.png',
    content: 'I’m good, just working on the Angular chat UI 😊',
    timestamp: new Date(Date.now() - 1000 * 60 * 3),
    isMine: true,
  },
  {
    id: '3',
    userId: 'u1',
    userName: 'Alice',
    userAvatar: 'assets/avatars/alice.png',
    content: 'Nice! Send me a screenshot later.',
    timestamp: new Date(Date.now() - 1000 * 60 * 1),
    isMine: false,
  },
];
