import { HttpClient, HttpParams } from '@angular/common/http';
import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { environment } from '../../../environments/environment';
import { Comment, CreateCommentRequest } from '../../models/comment.model';

@Injectable({
  providedIn: 'root',
})
export class CommentService {
  private apiUrl = `${environment.apiUrl}/v1/consultation`;

  constructor(private http: HttpClient) {}

  // Get comments tree for a post (with nested replies)
  getComments(
    postId: string,
    limit: number = 20,
    skip: number = 0,
    depth: number = 2
  ): Observable<{ comments: Comment[]; total: number }> {
    const params = new HttpParams()
      .set('limit', limit.toString())
      .set('skip', skip.toString())
      .set('depth', depth.toString());
    return this.http.get<{ comments: Comment[]; total: number }>(
      `${this.apiUrl}/posts/${postId}/comments`,
      { params }
    );
  }

  // Create comment
  createComment(
    postId: string,
    comment: CreateCommentRequest
  ): Observable<Comment> {
    return this.http.post<Comment>(
      `${this.apiUrl}/posts/${postId}/comments`,
      comment
    );
  }

  // Update comment
  updateComment(commentId: string, content: string): Observable<Comment> {
    return this.http.put<Comment>(`${this.apiUrl}/comments/${commentId}`, {
      content,
    });
  }

  // Delete comment
  deleteComment(commentId: string): Observable<void> {
    return this.http.delete<void>(`${this.apiUrl}/comments/${commentId}`);
  }

  // Get replies
  getReplies(commentId: string): Observable<Comment[]> {
    return this.http.get<Comment[]>(
      `${this.apiUrl}/comments/${commentId}/replies`
    );
  }
}
