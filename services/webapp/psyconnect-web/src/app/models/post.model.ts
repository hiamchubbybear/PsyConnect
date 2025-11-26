export interface Post {
  id: string;
  title: string;
  content: string;
  author_id: string;

  // Enhanced fields
  tags: string[];
  categories: string[];
  media: MediaAttachment[];
  mentions: string[];
  hashtags: string[];

  // Metadata
  visibility: 'public' | 'private' | 'followers';
  post_type: 'article' | 'question' | 'discussion' | 'resource';

  // Engagement metrics
  view_count: number;
  like_count: number;
  comment_count: number;
  share_count: number;

  // New vote fields
  upvote_count?: number;
  downvote_count?: number;
  user_vote?: 'upvote' | 'downvote' | null; // Current user's vote

  // Control
  is_deleted: boolean;

  created_at: string;
  updated_at: string;
  status?: 'published' | 'draft';
}

export interface MediaAttachment {
  type: 'image' | 'video' | 'document';
  url: string;
  caption?: string;
}

export interface CreatePostRequest {
  title: string;
  content: string;
  tags?: string[];
  categories?: string[];
  media?: MediaAttachment[];
  visibility?: 'public' | 'private' | 'followers';
  post_type?: 'article' | 'question' | 'discussion' | 'resource';
}
