export interface Reaction {
  id: string;
  post_id: string;
  user_id: string;
  reaction_type: ReactionType;
  created_at: string;
}

// Simplified to upvote/downvote only
export type ReactionType = 'upvote' | 'downvote';

// Reaction images using existing SVG assets
export const REACTION_IMAGES: Record<ReactionType, string> = {
  upvote: '/assets/reaction/arrow-big-up-dash.svg',
  downvote: '/assets/reaction/arrow-big-down-dash.svg'
};

// Action icons
export const ACTION_ICONS = {
  comment: '/assets/reaction/message-circle.svg',
  share: '/assets/reaction/share.svg'
};

// For backward compatibility
export const REACTION_EMOJIS: Record<ReactionType, string> = {
  upvote: '👍',
  downvote: '👎'
};

export const REACTION_TYPES = REACTION_EMOJIS;
