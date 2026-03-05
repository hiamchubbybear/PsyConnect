export interface Comment {
  id: string;
  post_id: string;
  user_id: string;
  author_id: string; 
  content: string;
  parent_comment_id?: string;

  
  depth?: number;
  path?: string;
  reply_count?: number;
  replies?: Comment[];

  like_count: number;
  is_deleted: boolean;
  created_at: string;
  updated_at: string;
  isLiked?: boolean;

  
  author?: {
    id: string;
    first_name: string;
    last_name: string;
    avatar?: string;
  };
}

export interface CreateCommentRequest {
  content: string;
  parent_comment_id?: string;
}
