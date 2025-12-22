export interface Reaction {
  id: string;
  post_id: string;
  user_id: string;
  reaction_type: ReactionType;
  created_at: string;
}

// Simplified to up/down only to match backend constants
export type ReactionType = 'up' | 'down';

// Reaction images using existing SVG assets
export const REACTION_IMAGES: Record<ReactionType, string> = {
  up: '/assets/reaction/arrow-big-up-dash.svg',
  down: '/assets/reaction/arrow-big-down-dash.svg',
};

// Action icons
export const ACTION_ICONS = {
  comment: '/assets/reaction/message-circle.svg',
  share: '/assets/reaction/share.svg',
};

// Fallback text (simple arrows instead of emoji)
export const REACTION_FALLBACK: Record<ReactionType, string> = {
  up: '↑',
  down: '↓',
};

export const REACTION_TYPES = REACTION_FALLBACK;
