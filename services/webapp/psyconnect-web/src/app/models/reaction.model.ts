export interface Reaction {
  id: string;
  post_id: string;
  user_id: string;
  reaction_type: ReactionType;
  created_at: string;
}

export type ReactionType = 'like' | 'love' | 'laugh' | 'think' | 'sad' | 'angry';

export const REACTION_TYPES: Record<ReactionType, string> = {
  like: '👍',
  love: '❤️',
  laugh: '😄',
  think: '🤔',
  sad: '😢',
  angry: '😠'
};
