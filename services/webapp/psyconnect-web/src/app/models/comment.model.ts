export interface Comment {
  id: string;
  post_id: string;
  user_id: string;
  author_id: string; // Author ID for display
  content: string;
  parent_comment_id?: string;
  like_count: number;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
}

export interface CreateCommentRequest {
  content: string;
  parent_comment_id?: string;
}
